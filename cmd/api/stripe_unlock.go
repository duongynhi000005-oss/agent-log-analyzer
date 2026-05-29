package main

import (
	"encoding/json"
	"errors"
	"fmt"
	htmlstd "html"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/priivacy-ai/agent-log-analyzer/internal/app"
)

const optimizationPackPriceCents int64 = 5000
const optimizationPackCurrency = "usd"

type optimizationUnlockRequest struct {
	SourceReportJobID string `json:"source_report_job_id"`
	SourceReportToken string `json:"source_report_token"`
}

type optimizationUnlockResponse struct {
	CheckoutURL string `json:"checkout_url"`
	SessionID   string `json:"session_id"`
}

type stripeCheckoutSession struct {
	ID            string            `json:"id"`
	URL           string            `json:"url"`
	PaymentStatus string            `json:"payment_status"`
	Status        string            `json:"status"`
	ClientRefID   string            `json:"client_reference_id"`
	CustomerEmail string            `json:"customer_email"`
	CustomerDetails struct {
		Email string `json:"email"`
	} `json:"customer_details"`
	Metadata map[string]string `json:"metadata"`
}

func createOptimizationUnlockHandler(store app.APIStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentStore, ok := store.(app.PaymentUnlockStore)
		if !ok {
			writeErrorOrHTML(w, r, http.StatusNotImplemented, "optimization unlock unavailable")
			return
		}
		request, err := parseOptimizationUnlockRequest(r)
		if err != nil {
			writeErrorOrHTML(w, r, http.StatusBadRequest, err.Error())
			return
		}
		job, _, err := authorizedReport(store, request.SourceReportJobID, request.SourceReportToken)
		if err != nil {
			writeErrorOrHTML(w, r, http.StatusUnauthorized, err.Error())
			return
		}
		if job.Status != app.StatusCompleted {
			writeErrorOrHTML(w, r, http.StatusConflict, "source report is not completed")
			return
		}
		session, err := createStripeCheckoutSession(r, job.ID, request.SourceReportToken)
		if err != nil {
			slog.Warn("stripe checkout create failed", "error_category", "stripe_checkout_create")
			writeErrorOrHTML(w, r, http.StatusServiceUnavailable, "paid optimization unlock is not configured")
			return
		}
		now := time.Now().UTC()
		unlock := app.PaymentUnlock{
			ID:                       app.NewJobID(),
			StripeSessionID:          session.ID,
			SourceReportJobID:        job.ID,
			SourceReportTokenHash:    tokenHash(request.SourceReportToken),
			AmountCents:              optimizationPackPriceCents,
			Currency:                 optimizationPackCurrency,
			Status:                   app.PaymentUnlockPending,
			CreatedAt:                now,
			UpdatedAt:                now,
			LastStripePaymentStatus:   session.PaymentStatus,
			LastStripeCheckoutStatus:  session.Status,
		}
		if err := paymentStore.CreatePaymentUnlock(unlock); err != nil {
			writeErrorOrHTML(w, r, http.StatusInternalServerError, "could not record optimization unlock")
			return
		}
		if wantsHTML(r) {
			http.Redirect(w, r, session.URL, http.StatusSeeOther)
			return
		}
		writeJSON(w, http.StatusCreated, optimizationUnlockResponse{CheckoutURL: session.URL, SessionID: session.ID})
	}
}

func optimizationUnlockSuccessHandler(store app.APIStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentStore, ok := store.(app.PaymentUnlockStore)
		if !ok {
			writeError(w, http.StatusNotImplemented, "optimization unlock unavailable")
			return
		}
		sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
		if sessionID == "" {
			writeError(w, http.StatusBadRequest, "checkout session is required")
			return
		}
		unlock, err := paymentStore.GetPaymentUnlockByStripeSessionID(sessionID)
		if err != nil {
			writeError(w, http.StatusNotFound, "optimization unlock not found")
			return
		}
		session, err := retrieveStripeCheckoutSession(sessionID)
		if err != nil {
			slog.Warn("stripe checkout retrieve failed", "error_category", "stripe_checkout_retrieve")
			writeError(w, http.StatusServiceUnavailable, "could not confirm payment")
			return
		}
		unlock.LastStripePaymentStatus = session.PaymentStatus
		unlock.LastStripeCheckoutStatus = session.Status
		if email := stripeSessionEmail(session); email != "" {
			unlock.LastStripeCustomerEmailHash = app.HashEmail(strings.ToLower(email))
		}
		if session.ClientRefID != "" && session.ClientRefID != unlock.SourceReportJobID {
			_ = paymentStore.UpdatePaymentUnlock(unlock)
			writeError(w, http.StatusUnauthorized, "checkout session does not match report")
			return
		}
		if session.Metadata != nil && session.Metadata["source_report_token_hash"] != "" && session.Metadata["source_report_token_hash"] != unlock.SourceReportTokenHash {
			_ = paymentStore.UpdatePaymentUnlock(unlock)
			writeError(w, http.StatusUnauthorized, "checkout session does not match report token")
			return
		}
		if session.PaymentStatus != "paid" {
			_ = paymentStore.UpdatePaymentUnlock(unlock)
			renderPaymentPendingPage(w, unlock.SourceReportJobID)
			return
		}
		token, err := newToken()
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not create download token")
			return
		}
		now := time.Now().UTC()
		unlock.Status = app.PaymentUnlockPaid
		unlock.DownloadTokenHash = tokenHash(token)
		unlock.DownloadTokenExpiresAt = now.Add(optimizationDownloadTokenTTL())
		if unlock.PaidAt.IsZero() {
			unlock.PaidAt = now
		}
		if err := paymentStore.UpdatePaymentUnlock(unlock); err != nil {
			writeError(w, http.StatusInternalServerError, "could not unlock optimization pack")
			return
		}
		renderOptimizationInstallPage(w, publicBaseURL(r)+"/api/paid-artifacts/"+unlock.SourceReportJobID+"/"+token+"/plugin.zip", unlock.DownloadTokenExpiresAt)
	}
}

func getPaidArtifactHandler(store app.APIStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paymentStore, ok := store.(app.PaymentUnlockStore)
		if !ok {
			writeError(w, http.StatusNotImplemented, "paid artifact unavailable")
			return
		}
		unlock, err := paymentStore.GetPaymentUnlockByDownloadTokenHash(tokenHash(r.PathValue("token")))
		if err != nil || unlock.SourceReportJobID != r.PathValue("id") {
			writeError(w, http.StatusUnauthorized, "invalid paid artifact token")
			return
		}
		if unlock.Status != app.PaymentUnlockPaid || time.Now().UTC().After(unlock.DownloadTokenExpiresAt) {
			writeError(w, http.StatusGone, "paid artifact token expired")
			return
		}
		job, err := store.GetJob(unlock.SourceReportJobID)
		if err != nil {
			writeError(w, http.StatusNotFound, "source report not found")
			return
		}
		if job.ReportTokenHash != unlock.SourceReportTokenHash {
			writeError(w, http.StatusUnauthorized, "source report token mismatch")
			return
		}
		report, err := store.GetReport(job.ID)
		if err != nil {
			writeError(w, http.StatusNotFound, "report not found")
			return
		}
		artifactBytes, err := renderPluginArtifactZip(report, publicBaseURL(r)+r.URL.Path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not generate plugin artifact")
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", `attachment; filename="agent-analyzer-optimization-pack.zip"`)
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(artifactBytes)
	}
}

func parseOptimizationUnlockRequest(r *http.Request) (optimizationUnlockRequest, error) {
	var request optimizationUnlockRequest
	if isJSONRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return request, errors.New("invalid optimization unlock request")
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return request, errors.New("invalid optimization unlock form")
		}
		request.SourceReportJobID = r.Form.Get("source_report_job_id")
		request.SourceReportToken = r.Form.Get("source_report_token")
	}
	if request.SourceReportJobID == "" || request.SourceReportToken == "" {
		return request, errors.New("source report is required")
	}
	return request, nil
}

func createStripeCheckoutSession(r *http.Request, jobID, reportToken string) (stripeCheckoutSession, error) {
	secret := strings.TrimSpace(os.Getenv("STRIPE_SECRET_KEY"))
	if secret == "" {
		return stripeCheckoutSession{}, errors.New("missing stripe secret key")
	}
	baseURL := publicBaseURL(r)
	form := url.Values{}
	form.Set("mode", "payment")
	form.Set("client_reference_id", jobID)
	form.Set("success_url", baseURL+"/optimization-unlock/success?session_id={CHECKOUT_SESSION_ID}")
	form.Set("cancel_url", baseURL+"/r/"+url.PathEscape(jobID)+"/"+url.PathEscape(reportToken)+"?payment=cancelled#plugin-purchase")
	form.Set("line_items[0][quantity]", "1")
	form.Set("line_items[0][price_data][currency]", optimizationPackCurrency)
	form.Set("line_items[0][price_data][unit_amount]", fmt.Sprintf("%d", optimizationPackPriceCents))
	form.Set("line_items[0][price_data][product_data][name]", "Agent Analyzer optimization pack")
	form.Set("line_items[0][price_data][product_data][description]", "Installable remediation/plugin pack generated from your analysis.")
	form.Set("metadata[source_report_job_id]", jobID)
	form.Set("metadata[source_report_token_hash]", tokenHash(reportToken))
	return callStripe(http.MethodPost, "/v1/checkout/sessions", form, secret)
}

func retrieveStripeCheckoutSession(sessionID string) (stripeCheckoutSession, error) {
	secret := strings.TrimSpace(os.Getenv("STRIPE_SECRET_KEY"))
	if secret == "" {
		return stripeCheckoutSession{}, errors.New("missing stripe secret key")
	}
	return callStripe(http.MethodGet, "/v1/checkout/sessions/"+url.PathEscape(sessionID), nil, secret)
}

func callStripe(method, path string, form url.Values, secret string) (stripeCheckoutSession, error) {
	base := strings.TrimRight(getenv("CLAUDE_ANALYZER_STRIPE_API_BASE", "https://api.stripe.com"), "/")
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequest(method, base+path, body)
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	req.SetBasicAuth(secret, "")
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return stripeCheckoutSession{}, fmt.Errorf("stripe status %d", resp.StatusCode)
	}
	var session stripeCheckoutSession
	if err := json.Unmarshal(data, &session); err != nil {
		return session, err
	}
	if session.ID == "" {
		return session, errors.New("stripe response missing session id")
	}
	return session, nil
}

func renderPaymentPendingPage(w http.ResponseWriter, jobID string) {
	body := fmt.Sprintf(`<p>Stripe has not confirmed payment for this checkout session yet.</p><p>Return to your private report link and try again after checkout completes.</p><p class="command-note">Report ID: %s</p>`, htmlstd.EscapeString(jobID))
	renderSimpleHTML(w, "Payment not confirmed", body)
}

func renderOptimizationInstallPage(w http.ResponseWriter, artifactURL string, expiresAt time.Time) {
	command := `PLUGIN_ZIP="/path/to/agent-analyzer-optimization-pack.zip"
claude plugin install "$PLUGIN_ZIP"
/agent-analyzer-status`
	escapedCommand := htmlstd.EscapeString(command)
	body := fmt.Sprintf(
		`<p>Your payment is confirmed. This short-lived page contains the generated optimization pack from your sanitized analysis.</p><p class="download-button-row"><a class="plugin-cta" href="%s">Download optimization pack</a></p><p>Download expires at <strong>%s UTC</strong>.</p><p>For Claude Code, unzip or save the pack locally, review <strong>INSTALL.md</strong>, then install the Claude plugin:</p><div class="simple-command-copy"><pre><code>%s</code></pre><button type="button" class="copy-agents-line" data-copy="%s">Copy command</button></div><p>The pack includes deterministic hooks, context compression helpers, slash-command coach guidance, CLAUDE.md optimizer recommendations, retrieval recommendations, and statusline telemetry. Raw transcripts are not included.</p>`,
		htmlstd.EscapeString(artifactURL),
		htmlstd.EscapeString(expiresAt.Format("2006-01-02 15:04")),
		escapedCommand,
		escapedCommand,
	)
	renderSimpleHTML(w, "Optimization pack unlocked", body)
}

func stripeSessionEmail(session stripeCheckoutSession) string {
	if session.CustomerDetails.Email != "" {
		return session.CustomerDetails.Email
	}
	return session.CustomerEmail
}

func optimizationDownloadTokenTTL() time.Duration {
	raw := getenv("CLAUDE_ANALYZER_OPTIMIZATION_DOWNLOAD_TTL", "30m")
	ttl, err := time.ParseDuration(raw)
	if err != nil || ttl <= 0 {
		return 30 * time.Minute
	}
	return ttl
}
