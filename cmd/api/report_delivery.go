package main

import (
	"encoding/json"
	"errors"
	"fmt"
	htmlstd "html"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/priivacy-ai/agent-log-analyzer/internal/analyzer"
	"github.com/priivacy-ai/agent-log-analyzer/internal/app"
)

type reportDeliveryRequest struct {
	Email             string `json:"email"`
	MarketingOptIn    bool   `json:"marketing_opt_in"`
	SourceReportJobID string `json:"source_report_job_id"`
	SourceReportToken string `json:"source_report_token"`
}

type reportDeliveryResponse struct {
	DeliveryID string `json:"delivery_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	ReportURL  string `json:"report_url,omitempty"`
	PluginURL  string `json:"plugin_url,omitempty"`
}

func createReportDeliveryHandler(store app.APIStore, sender emailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unlockStore, ok := store.(app.EmailUnlockStore)
		if !ok {
			writeError(w, http.StatusNotImplemented, "report delivery unavailable")
			return
		}
		request, err := parseReportDeliveryRequest(r)
		if err != nil {
			writeErrorOrHTML(w, r, http.StatusBadRequest, err.Error())
			return
		}
		normalized, err := normalizeEmail(request.Email)
		if err != nil {
			writeErrorOrHTML(w, r, http.StatusBadRequest, err.Error())
			return
		}
		if request.SourceReportJobID == "" || request.SourceReportToken == "" {
			writeErrorOrHTML(w, r, http.StatusBadRequest, "source report is required")
			return
		}
		job, report, err := authorizedReport(store, request.SourceReportJobID, request.SourceReportToken)
		if err != nil {
			writeErrorOrHTML(w, r, http.StatusUnauthorized, err.Error())
			return
		}
		reportPack, err := renderDownloadPackage(job, report)
		if err != nil {
			writeErrorOrHTML(w, r, http.StatusInternalServerError, "could not generate report pack")
			return
		}
		reportURL := publicBaseURL(r) + "/api/public-reports/" + job.ID + "/" + request.SourceReportToken + "/download.zip"
		now := time.Now().UTC()
		delivery := app.EmailUnlock{
			ID:                           app.NewJobID(),
			Email:                        normalized,
			EmailHash:                    app.HashEmail(normalized),
			MarketingOptIn:               request.MarketingOptIn,
			SourceReportJobID:            job.ID,
			Status:                       app.EmailUnlockUsed,
			CreatedAt:                    now,
			UpdatedAt:                    now,
			ConfirmedAt:                  now,
			LastTransactionalEmailSentAt: now,
		}
		if err := unlockStore.CreateEmailUnlock(delivery); err != nil {
			writeErrorOrHTML(w, r, http.StatusInternalServerError, "could not record email request")
			return
		}
		message := emailMessage{
			To:      normalized,
			Subject: "Your Agent Analyzer report pack",
			Body:    reportDeliveryEmailBody(reportURL, publicBaseURL(r)+"/r/"+job.ID+"/"+request.SourceReportToken+"#plugin-purchase"),
			Attachments: []emailAttachment{
				{
					Filename:    "agent-analyzer-report-pack.zip",
					ContentType: "application/zip",
					Data:        reportPack,
				},
			},
		}
		if err := sender.Send(message); err != nil {
			if errors.As(err, &errEmailSuppressed{}) {
				writeErrorOrHTML(w, r, http.StatusConflict, "email address is suppressed for transactional delivery")
				return
			}
			slogEmailDeliveryFailure("report_delivery", delivery.ID, err)
			writeEmailDeliveryErrorOrHTML(w, r, err)
			return
		}
		slog.Info("report delivery sent", "delivery_id", delivery.ID, "email_hash", delivery.EmailHash, "marketing_opt_in", delivery.MarketingOptIn)
		if wantsHTML(r) {
			renderReportDeliverySentPage(w, normalized, reportURL, publicBaseURL(r)+"/r/"+job.ID+"/"+request.SourceReportToken+"#plugin-purchase")
			return
		}
		writeJSON(w, http.StatusAccepted, reportDeliveryResponse{
			DeliveryID: delivery.ID,
			Status:     string(delivery.Status),
			Message:    "report pack sent",
			ReportURL:  reportURL,
		})
	}
}

func parseReportDeliveryRequest(r *http.Request) (reportDeliveryRequest, error) {
	var request reportDeliveryRequest
	if isJSONRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return request, errors.New("invalid report delivery request")
		}
		return request, nil
	}
	if err := r.ParseForm(); err != nil {
		return request, errors.New("invalid report delivery form")
	}
	request.Email = r.Form.Get("email")
	request.MarketingOptIn = r.Form.Get("marketing_opt_in") == "1" || r.Form.Get("marketing_opt_in") == "true" || r.Form.Get("marketing_opt_in") == "on"
	request.SourceReportJobID = r.Form.Get("source_report_job_id")
	request.SourceReportToken = r.Form.Get("source_report_token")
	return request, nil
}

func authorizedReport(store app.APIStore, jobID, reportToken string) (app.Job, analyzer.Report, error) {
	job, err := store.GetJob(jobID)
	if err != nil {
		return app.Job{}, analyzer.Report{}, errors.New("source report not found")
	}
	if !tokenMatches(job.ReportTokenHash, reportToken) {
		return app.Job{}, analyzer.Report{}, errors.New("invalid source report token")
	}
	report, err := store.GetReport(job.ID)
	if err != nil {
		return app.Job{}, analyzer.Report{}, errors.New("source report not found")
	}
	return job, report, nil
}

func isJSONRequest(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "application/json")
}

func renderReportDeliverySentPage(w http.ResponseWriter, email, reportURL, purchaseURL string) {
	body := fmt.Sprintf(
		`<p>We recorded <strong>%s</strong> and sent the free report pack link to that address.</p><p class="download-button-row"><a class="plugin-cta" href="%s">Download report pack</a><a class="plugin-cta" href="%s">Unlock optimization pack</a></p><p>The email also reminds you about the Spec Kitty training voucher and links to the <a href="https://github.com/Priivacy-ai/spec-kitty" rel="noopener noreferrer">Spec Kitty GitHub repo</a>.</p><p>The free report is enough to understand the problems found in your logs. The paid optimization pack is a one-time $50 unlock for installable Claude Code fixes and deeper customization generated from the same sanitized report JSON.</p>`,
		htmlstd.EscapeString(email),
		htmlstd.EscapeString(reportURL),
		htmlstd.EscapeString(purchaseURL),
	)
	renderSimpleHTML(w, "Report pack sent", body)
}

func reportDeliveryEmailBody(reportURL, purchaseURL string) string {
	return fmt.Sprintf(`Your Agent Analyzer report pack is ready.

The safe local scan found token leaks in your own coding logs. The free report pack is enough to understand the rereads, retries, context bloat, and noisy tool output found in your sessions.

Downloads and unlock:
- Report pack: %s
- One-time $50 optimization pack unlock: %s

Attachments:
- agent-analyzer-report-pack.zip: branded PDF guide, personalized PDF report, sanitized report JSON, plugin preview, and partner voucher.

Spec Kitty training voucher:
- Your report pack includes the partner training voucher.
- Spec Kitty Teamspace is coming soon for teams that want shared agentic coding workflows.
- Spec Kitty GitHub repo: https://github.com/Priivacy-ai/spec-kitty

Paid optimization pack:
- Deterministic hook pack.
- Context compression helpers.
- Slash-command coach.
- CLAUDE.md optimizer recommendations.
- Retrieval recommendations.
- Statusline telemetry.

After Stripe confirms payment, the success page shows install instructions and a payment-confirmed artifact link. The artifact is generated from sanitized report JSON only.

Privacy boundary:
- Raw transcripts were not attached.
- Raw transcripts were not uploaded to Agent Analyzer.
- The report pack was generated from sanitized report JSON for your private report link.
`, reportURL, purchaseURL)
}
