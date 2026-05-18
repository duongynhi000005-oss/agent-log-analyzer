package signatureresearch

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type fakeClient struct {
	body  string
	queue []string
	calls int
}

func (f *fakeClient) Do(req *http.Request) (*http.Response, error) {
	body := f.body
	if len(f.queue) > 0 {
		body = f.queue[0]
		f.queue = f.queue[1:]
	}
	f.calls++
	return &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}, nil
}

func TestCollectMCPRegistryCandidates(t *testing.T) {
	body := `{
		"servers": [
			{"server": {
				"name": "io.github.modelcontextprotocol/filesystem",
				"title": "Filesystem",
				"description": "Read and write local files",
				"packages": [{"registry":"npm","name":"@modelcontextprotocol/server-filesystem"}]
			}},
			{"server": {
				"name": "com.example/context7",
				"title": "Context7",
				"description": "Current docs for agents"
			}}
		]
	}`
	report, err := Collect(context.Background(), Config{Sources: []Source{{
		ID: "official-mcp", Kind: "mcp_registry", URL: "https://registry.example.test/v0/servers", Category: "mcp_servers",
	}}}, &fakeClient{body: body})
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	assertCandidate(t, report, "mcp_servers", "filesystem")
	assertCandidate(t, report, "mcp_servers", "context7")
	assertCandidate(t, report, "mcp_servers", "modelcontextprotocol_server_filesystem")
}

func TestTextExtractionFindsPublicFingerprintsWithoutPathNoise(t *testing.T) {
	body := `Use mcp__github__create_issue, claude mcp add playwright, and .claude/skills/security/SKILL.md.
Ignore /Users/robert/project and src/auth.ts.`
	report, err := Collect(context.Background(), Config{Sources: []Source{{
		ID: "docs", Kind: "text", URL: "https://example.test/docs", Category: "skills",
	}}}, &fakeClient{body: body})
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	assertCandidate(t, report, "mcp_servers", "github")
	assertCandidate(t, report, "mcp_servers", "playwright")
	assertCandidate(t, report, "skills", "security")
	if hasCandidate(report, "skills", "users") || hasCandidate(report, "skills", "auth") {
		t.Fatalf("path fragments should not become skill candidates: %#v", report.Candidates)
	}
}

func TestMCPRegistryPaginationUsesMaxPages(t *testing.T) {
	client := &fakeClient{queue: []string{
		`{"servers":[{"server":{"name":"io.github.example/first"}}],"metadata":{"nextCursor":"page-two"}}`,
		`{"servers":[{"server":{"name":"io.github.example/second"}}],"metadata":{}}`,
	}}
	report, err := Collect(context.Background(), Config{Sources: []Source{{
		ID: "official-mcp", Kind: "mcp_registry", URL: "https://registry.example.test/v0/servers?limit=1", Category: "mcp_servers", MaxPages: 2,
	}}}, client)
	if err != nil {
		t.Fatalf("Collect returned error: %v", err)
	}
	if client.calls != 2 {
		t.Fatalf("expected two paginated fetches, got %d", client.calls)
	}
	assertCandidate(t, report, "mcp_servers", "first")
	assertCandidate(t, report, "mcp_servers", "second")
}

func assertCandidate(t *testing.T, report Report, category, id string) {
	t.Helper()
	if !hasCandidate(report, category, id) {
		t.Fatalf("missing candidate %s/%s in %#v", category, id, report.Candidates)
	}
}

func hasCandidate(report Report, category, id string) bool {
	for _, candidate := range report.Candidates {
		if candidate.Category == category && candidate.ID == id {
			return true
		}
	}
	return false
}
