package main

import (
	"encoding/json"
	"errors"
	"fmt"
	htmlstd "html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/priivacy-ai/agent-log-analyzer/internal/app"
)

const optimizationPackPriceCents = 5000
const optimizationDownloadTTL = 24 * time.Hour

type optimizationCheckoutResponse struct {
	CheckoutURL string `json:"checkout_url"`
	SessionID   string `json:"session_id"`
}

type optimizationCheckoutRequest struct {
	SourceReportJobID string `json:"source_report_job_id"`
	SourceReportToken string `json:"source_report_token"`
	JobID             string `json:"job_id"`
	ReportToken       string `json:"report_token"`
}

type stripeCheckoutSession struct {
	ID              string            `json:"id"`
	URL             string            `json:"url"`
	Status          string            `json:"status"`
	PaymentStatus   string            `json:"payment_status"`
	Created         int64             `json:"created"`
	ClientReference string            `json:"client_reference_id"`
	Metadata        map[string]string `json:"metadata"`
}

func createOptimizationCheckoutHandler(store app.APIStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := parseOptimizationCheckoutRequest(r)
		if err != nil {
			writeErrorOrHTML(w, r, http.StatusBadRequest, err.Error())
			return
		}
		jobID := strings.TrimSpace(firstNonEmpty(request.SourceReportJobID, request.JobID))
		reportToken := strings.TrimSpace(firstNonEmpty(request.SourceReportToken, request.ReportToken))
		if jobID == "" || reportToken == "" {
			writeErrorOrHTML(w, r, http.StatusBadRequest, "source report is required")
			return
		}
		if _, _, err := authorizedReport(store, jobID, reportToken); err != nil {
			writeErrorOrHTML(w, r, http.StatusUnauthorized, err.Error())
			return
		}
		session, err := createStripeOptimizationSession(r, jobID, reportToken)
		if err != nil {
			writeErrorOrHTML(w, r, http.StatusPaymentRequired, err.Error())
			return
		}
		if wantsHTML(r) || !isJSONRequest(r) {
			http.Redirect(w, r, session.URL, http.StatusSeeOther)
			return
		}
		writeJSON(w, http.StatusCreated, optimizationCheckoutResponse{CheckoutURL: session.URL, SessionID: session.ID})
	}
}

func parseOptimizationCheckoutRequest(r *http.Request) (optimizationCheckoutRequest, error) {
	var request optimizationCheckoutRequest
	if isJSONRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return request, errors.New("invalid checkout request")
		}
		return request, nil
	}
	if err := r.ParseForm(); err != nil {
		return request, errors.New("invalid checkout form")
	}
	request.SourceReportJobID = r.Form.Get("source_report_job_id")
	request.SourceReportToken = r.Form.Get("source_report_token")
	request.JobID = r.Form.Get("job_id")
	request.ReportToken = r.Form.Get("report_token")
	return request, nil
}

func optimizationUnlockSuccessHandler(store app.APIStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
		if sessionID == "" {
			writeError(w, http.StatusBadRequest, "missing checkout session")
			return
		}
		session, err := retrievePaidStripeSession(sessionID)
		if err != nil {
			renderSimpleHTML(w, "Payment not confirmed", `<p>Stripe has not confirmed this optimization pack payment, or the download page has expired. Return to your report and try the unlock again.</p>`)
			return
		}
		jobID, reportToken, err := reportFromPaidSession(session)
		if err != nil {
			writeError(w, http.StatusForbidden, "invalid checkout session")
			return
		}
		if _, _, err := authorizedReport(store, jobID, reportToken); err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		artifactURL := publicBaseURL(r) + "/api/paid-artifacts/" + jobID + "/" + reportToken + "/" + url.PathEscape(session.ID) + "/plugin.zip"
		command := `PLUGIN_ZIP="/path/to/agent-analyzer-optimization-plugin.zip"
claude plugin install "$PLUGIN_ZIP"
/agent-analyzer-status`
		escapedCommand := htmlstd.EscapeString(command)
		body := fmt.Sprintf(
			`<p>Your one-time $50 payment is confirmed. Download the optimization pack generated from your analysis. This checkout-session download page expires after 24 hours.</p><p class="download-button-row"><a class="plugin-cta" href="%s">Download optimization pack</a><a class="plugin-cta" href="/r/%s/%s">Back to report</a></p><p>Included: deterministic hook pack, context compression helpers, slash-command coach, CLAUDE.md optimizer recommendations, retrieval recommendations, and statusline telemetry. The pack is generated from sanitized report JSON only.</p><p>Claude Code install:</p><div class="simple-command-copy"><pre><code>%s</code></pre><button type="button" class="copy-agents-line" data-copy="%s">Copy command</button></div><p>Choose your harness in <strong>INSTALL.md</strong>. Claude Code users should install persistently with <strong>claude plugin install</strong>; use <strong>claude --plugin-dir "$PLUGIN_ZIP"</strong> only for a one-session preview.</p>`,
			htmlstd.EscapeString(artifactURL),
			htmlstd.EscapeString(jobID),
			htmlstd.EscapeString(reportToken),
			escapedCommand,
			escapedCommand,
		)
		renderSimpleHTML(w, "Optimization pack unlocked", body)
	}
}

func getPaidArtifactHandler(store app.APIStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := retrievePaidStripeSession(r.PathValue("session"))
		if err != nil {
			writeError(w, http.StatusPaymentRequired, "payment not confirmed")
			return
		}
		jobID, reportToken, err := reportFromPaidSession(session)
		if err != nil || jobID != r.PathValue("id") || reportToken != r.PathValue("token") {
			writeError(w, http.StatusForbidden, "invalid checkout session")
			return
		}
		_, report, err := authorizedReport(store, jobID, reportToken)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err.Error())
			return
		}
		artifactBytes, err := renderPluginArtifactZip(report, publicBaseURL(r)+r.URL.Path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "could not generate plugin artifact")
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", `attachment; filename="agent-analyzer-optimization-plugin.zip"`)
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(artifactBytes)
	}
}

func createStripeOptimizationSession(r *http.Request, jobID, reportToken string) (stripeCheckoutSession, error) {
	secret := stripeSecretKey()
	if secret == "" {
		return stripeCheckoutSession{}, errors.New("stripe checkout is not configured")
	}
	reportURL := publicBaseURL(r) + "/r/" + jobID + "/" + reportToken
	values := url.Values{}
	values.Set("mode", "payment")
	values.Set("client_reference_id", jobID)
	values.Set("success_url", publicBaseURL(r)+"/optimization-unlock/success?session_id={CHECKOUT_SESSION_ID}")
	values.Set("cancel_url", reportURL+"#plugin-purchase")
	values.Set("line_items[0][quantity]", "1")
	values.Set("line_items[0][price_data][currency]", "usd")
	values.Set("line_items[0][price_data][unit_amount]", fmt.Sprintf("%d", optimizationPackPriceCents))
	values.Set("line_items[0][price_data][product_data][name]", "Agent Analyzer optimization pack")
	values.Set("line_items[0][price_data][product_data][description]", "Install the optimization pack generated from your analysis.")
	values.Set("metadata[report_token]", reportToken)
	values.Set("metadata[report_url]", reportURL)
	return stripePostCheckoutSession(secret, values)
}

func retrievePaidStripeSession(sessionID string) (stripeCheckoutSession, error) {
	session, err := stripeGetCheckoutSession(stripeSecretKey(), sessionID)
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	if session.PaymentStatus != "paid" {
		return stripeCheckoutSession{}, errors.New("checkout session is not paid")
	}
	if session.Created <= 0 {
		return stripeCheckoutSession{}, errors.New("checkout session missing created time")
	}
	if time.Now().UTC().After(time.Unix(session.Created, 0).UTC().Add(optimizationDownloadTTL)) {
		return stripeCheckoutSession{}, errors.New("checkout session download expired")
	}
	return session, nil
}

func stripePostCheckoutSession(secret string, values url.Values) (stripeCheckoutSession, error) {
	req, err := http.NewRequest(http.MethodPost, stripeAPIBase()+"/v1/checkout/sessions", strings.NewReader(values.Encode()))
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	req.SetBasicAuth(secret, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return doStripeSessionRequest(req)
}

func stripeGetCheckoutSession(secret, sessionID string) (stripeCheckoutSession, error) {
	if secret == "" {
		return stripeCheckoutSession{}, errors.New("stripe checkout is not configured")
	}
	req, err := http.NewRequest(http.MethodGet, stripeAPIBase()+"/v1/checkout/sessions/"+url.PathEscape(sessionID), nil)
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	req.SetBasicAuth(secret, "")
	return doStripeSessionRequest(req)
}

func doStripeSessionRequest(req *http.Request) (stripeCheckoutSession, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return stripeCheckoutSession{}, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return stripeCheckoutSession{}, fmt.Errorf("stripe checkout failed with status %d", res.StatusCode)
	}
	var session stripeCheckoutSession
	if err := json.Unmarshal(data, &session); err != nil {
		return stripeCheckoutSession{}, err
	}
	if session.ID == "" {
		return stripeCheckoutSession{}, errors.New("stripe checkout response missing session id")
	}
	return session, nil
}

func reportFromPaidSession(session stripeCheckoutSession) (string, string, error) {
	reportToken := session.Metadata["report_token"]
	if session.ClientReference == "" || reportToken == "" {
		return "", "", errors.New("checkout session missing report metadata")
	}
	return session.ClientReference, reportToken, nil
}

func stripeSecretKey() string {
	return strings.TrimSpace(os.Getenv("STRIPE_SECRET_KEY"))
}

func stripeAPIBase() string {
	if base := strings.TrimRight(strings.TrimSpace(os.Getenv("STRIPE_API_BASE")), "/"); base != "" {
		return base
	}
	return "https://api.stripe.com"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
