package remediation

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/robertDouglass/claude-log-analyzer/internal/analyzer"
)

const GeneratorVersion = "0.1.0"

var safeValueRE = regexp.MustCompile(`^[a-z0-9_.:-]+$`)

var publicEcosystemIDs = map[string]map[string]bool{
	"agent": {
		"aider": true, "claude_code": true, "cline": true, "codex": true, "continue": true,
		"copilot": true, "cursor": true, "gemini_cli": true, "opencode": true, "roo": true, "windsurf": true,
	},
	"framework": {
		"agent_sessions": true, "aider": true, "bmad": true, "ccusage": true, "claude_code_hooks": true,
		"codex_transcript_viewer": true, "context_engineering": true, "openspec": true, "spec_kit": true,
		"spec_kitty": true, "superpowers": true,
	},
	"mcp": {
		"browser": true, "brave_search": true, "context7": true, "docker": true, "fetch": true,
		"figma": true, "filesystem": true, "git": true, "github": true, "gitlab": true, "gmail": true,
		"google_drive": true, "jira": true, "kubernetes": true, "linear": true, "memory": true,
		"notion": true, "playwright": true, "postgres": true, "puppeteer": true, "sentry": true,
		"sequential_thinking": true, "slack": true, "supabase": true,
	},
	"skill": {
		"benchmark": true, "design_review": true, "investigate": true, "plan_ceo_review": true,
		"plan_eng_review": true, "qa": true, "review": true, "security": true, "ship": true,
	},
	"plugin": {
		"browser": true, "canva": true, "figma": true, "github": true, "gmail": true,
		"google_calendar": true, "google_drive": true, "linear": true, "notion": true, "slack": true,
	},
	"package_manager": {
		"bun": true, "cargo": true, "composer": true, "go": true, "npm": true,
		"pip": true, "pnpm": true, "poetry": true, "uv": true, "yarn": true,
	},
}

type Options struct {
	ArtifactURL string
	GeneratedAt time.Time
}

type Artifact struct {
	SchemaVersion  string          `json:"schema_version"`
	Generator      string          `json:"generator"`
	PluginName     string          `json:"plugin_name"`
	PluginVersion  string          `json:"plugin_version"`
	GeneratedAt    time.Time       `json:"generated_at"`
	Source         SourceSummary   `json:"source"`
	Customizations []Customization `json:"customizations"`
	Files          []File          `json:"files"`
	Install        Install         `json:"install"`
}

type SourceSummary struct {
	AnalyzerVersion string   `json:"analyzer_version"`
	ScoreBucket     string   `json:"score_bucket"`
	WasteBucket     string   `json:"waste_bucket"`
	FindingIDs      []string `json:"finding_ids"`
	KnownEcosystem  []string `json:"known_ecosystem"`
}

type Customization struct {
	ID       string   `json:"id"`
	Reason   string   `json:"reason"`
	Evidence Evidence `json:"evidence"`
	Files    []string `json:"files"`
}

type Evidence struct {
	FindingID  string `json:"finding_id,omitempty"`
	Severity   string `json:"severity,omitempty"`
	CostImpact string `json:"cost_impact,omitempty"`
	Count      int    `json:"count,omitempty"`
	TokenShare int    `json:"token_share_pct,omitempty"`
}

type File struct {
	Path    string `json:"path"`
	Mode    string `json:"mode"`
	Content string `json:"content"`
}

type Install struct {
	Command          string   `json:"command"`
	ClaudePrompt     string   `json:"claude_prompt"`
	UninstallCommand string   `json:"uninstall_command"`
	Notes            []string `json:"notes"`
}

func Generate(report analyzer.Report, options Options) Artifact {
	generatedAt := options.GeneratedAt.UTC()
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	pluginName := "claude-analyzer-optimization"
	files := baseFiles(report, pluginName)
	customizations := customizationPlan(report)
	files = append(files, customizationFiles(customizations)...)
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	artifact := Artifact{
		SchemaVersion:  "2026-05-18",
		Generator:      "claude-log-analyzer/remediation@" + GeneratorVersion,
		PluginName:     pluginName,
		PluginVersion:  pluginVersion(report),
		GeneratedAt:    generatedAt,
		Source:         sourceSummary(report),
		Customizations: customizations,
		Files:          files,
	}
	artifact.Install = installInstructions(pluginName, options.ArtifactURL)
	return artifact
}

func baseFiles(report analyzer.Report, pluginName string) []File {
	manifest := map[string]any{
		"$schema":     "https://json.schemastore.org/claude-code-plugin-manifest.json",
		"name":        pluginName,
		"description": "Deterministic Claude Code workflow hygiene generated from a Claude Analyzer report.",
		"version":     pluginVersion(report),
		"author": map[string]string{
			"name": "Claude Log Analyzer",
		},
		"keywords": []string{"claude-code", "tokens", "context", "profiler", "workflow-hygiene"},
	}
	manifestJSON, _ := json.MarshalIndent(manifest, "", "  ")
	files := []File{
		{
			Path:    ".claude-plugin/plugin.json",
			Mode:    "0644",
			Content: string(manifestJSON) + "\n",
		},
		{
			Path:    "README.md",
			Mode:    "0644",
			Content: readme(report),
		},
		{
			Path:    "hooks/hooks.json",
			Mode:    "0644",
			Content: hooksJSON(),
		},
		{
			Path:    "scripts/claude-analyzer-hook.py",
			Mode:    "0755",
			Content: hookScript(),
		},
		{
			Path:    "commands/claude-analyzer-status.md",
			Mode:    "0644",
			Content: statusCommand(report),
		},
		{
			Path:    "skills/session-hygiene/SKILL.md",
			Mode:    "0644",
			Content: sessionHygieneSkill(report),
		},
	}
	return files
}

func customizationPlan(report analyzer.Report) []Customization {
	var out []Customization
	byID := map[string]analyzer.Finding{}
	for _, finding := range report.Findings {
		byID[finding.ID] = finding
	}
	if finding, ok := byID["repeated_file_reads"]; ok {
		out = append(out, Customization{
			ID:       "retrieval-hygiene",
			Reason:   "Repeated reads indicate the session is rereading files instead of preserving file state summaries.",
			Evidence: evidenceFor(finding),
			Files:    []string{"skills/retrieval-hygiene/SKILL.md"},
		})
	}
	if finding, ok := byID["tool_output_bloat"]; ok {
		out = append(out, Customization{
			ID:       "output-budget",
			Reason:   "Tool output consumed a high share of estimated context tokens.",
			Evidence: evidenceFor(finding),
			Files:    []string{"skills/output-budget/SKILL.md", "hooks/hooks.json", "scripts/claude-analyzer-hook.py"},
		})
	}
	if finding, ok := byID["retry_loop"]; ok {
		out = append(out, Customization{
			ID:       "retry-breaker",
			Reason:   "Retry-loop behavior suggests the session needs an explicit stop-and-reframe routine.",
			Evidence: evidenceFor(finding),
			Files:    []string{"skills/retry-breaker/SKILL.md"},
		})
	}
	if finding, ok := byID["context_growth_spikes"]; ok {
		out = append(out, Customization{
			ID:       "context-pivot",
			Reason:   "Context growth spikes indicate task pivots should trigger compaction or session splits.",
			Evidence: evidenceFor(finding),
			Files:    []string{"skills/session-hygiene/SKILL.md"},
		})
	}
	if len(out) == 0 {
		out = append(out, Customization{
			ID:       "baseline-hygiene",
			Reason:   "No high-confidence waste pattern dominated the scan; install the baseline session hygiene workflow.",
			Evidence: Evidence{},
			Files:    []string{"skills/session-hygiene/SKILL.md"},
		})
	}
	return out
}

func customizationFiles(customizations []Customization) []File {
	needed := map[string]bool{}
	for _, customization := range customizations {
		for _, path := range customization.Files {
			needed[path] = true
		}
	}
	var files []File
	if needed["skills/retrieval-hygiene/SKILL.md"] {
		files = append(files, File{
			Path:    "skills/retrieval-hygiene/SKILL.md",
			Mode:    "0644",
			Content: retrievalHygieneSkill(),
		})
	}
	if needed["skills/output-budget/SKILL.md"] {
		files = append(files, File{
			Path:    "skills/output-budget/SKILL.md",
			Mode:    "0644",
			Content: outputBudgetSkill(),
		})
	}
	if needed["skills/retry-breaker/SKILL.md"] {
		files = append(files, File{
			Path:    "skills/retry-breaker/SKILL.md",
			Mode:    "0644",
			Content: retryBreakerSkill(),
		})
	}
	return files
}

func evidenceFor(finding analyzer.Finding) Evidence {
	return Evidence{
		FindingID:  finding.ID,
		Severity:   finding.Severity,
		CostImpact: finding.CostImpact,
		Count:      finding.Evidence.Count,
		TokenShare: finding.Evidence.TokenShare,
	}
}

func sourceSummary(report analyzer.Report) SourceSummary {
	findingIDs := make([]string, 0, len(report.Findings))
	for _, finding := range report.Findings {
		if safeIdentifier(finding.ID) {
			findingIDs = append(findingIDs, finding.ID)
		}
	}
	sort.Strings(findingIDs)
	ecosystem := safeKnownEcosystem(report.Ecosystem)
	sort.Strings(ecosystem)
	return SourceSummary{
		AnalyzerVersion: report.Version,
		ScoreBucket:     report.AggregateEvent.ScoreBucket,
		WasteBucket:     report.AggregateEvent.WasteBucket,
		FindingIDs:      findingIDs,
		KnownEcosystem:  ecosystem,
	}
}

func safeKnownEcosystem(ecosystem analyzer.Ecosystem) []string {
	var out []string
	add := func(prefix string, values []string) {
		for _, value := range values {
			if safePublicID(prefix, value) {
				out = append(out, prefix+":"+value)
			}
		}
	}
	add("agent", ecosystem.CodingAgents)
	add("framework", ecosystem.WorkflowFrameworks)
	add("mcp", ecosystem.MCPServersKnown)
	add("skill", ecosystem.KnownSkills)
	add("plugin", ecosystem.KnownPlugins)
	add("package_manager", ecosystem.PackageManagers)
	for _, value := range []struct {
		prefix string
		raw    string
	}{
		{"client", ecosystem.Client},
		{"os", ecosystem.OperatingSystem},
		{"shell", ecosystem.Shell},
		{"vcs", ecosystem.VersionControl},
	} {
		if safeIdentifier(value.raw) {
			out = append(out, value.prefix+":"+value.raw)
		}
	}
	return out
}

func safePublicID(prefix, value string) bool {
	return safeIdentifier(value) && publicEcosystemIDs[prefix][value]
}

func safeIdentifier(value string) bool {
	return value != "" && safeValueRE.MatchString(value)
}

func pluginVersion(report analyzer.Report) string {
	if report.Version == "" {
		return "0.1.0"
	}
	return "0.1.0+" + strings.ReplaceAll(report.Version, ".", "-")
}

func installInstructions(pluginName, artifactURL string) Install {
	if artifactURL == "" {
		artifactURL = "<short-lived-plugin-zip-url>"
	}
	command := strings.Join([]string{
		`PLUGIN_URL="` + artifactURL + `"`,
		`PLUGIN_ZIP="$(mktemp -t claude-analyzer-plugin.XXXXXX.zip)"`,
		`curl -fsS "$PLUGIN_URL" -o "$PLUGIN_ZIP"`,
		`claude --plugin-dir "$PLUGIN_ZIP"`,
	}, "\n")
	prompt := "Install the generated Claude Analyzer optimization plugin for this session. Run the command below, explain what it installs, and ask for approval before executing it. Do not print plugin archive contents.\n\n```sh\n" + command + "\n```"
	return Install{
		Command:          command,
		ClaudePrompt:     prompt,
		UninstallCommand: "No persistent install is performed by the default command. Close the Claude Code session to unload " + pluginName + ".",
		Notes: []string{
			"The default command uses Claude Code's --plugin-dir support with a short-lived zip artifact.",
			"Marketplace installation can be added later once the plugin store is live.",
		},
	}
}

func readme(report analyzer.Report) string {
	return fmt.Sprintf(`# Claude Analyzer Optimization Plugin

Generated from deterministic Claude Analyzer metrics.

- Efficiency score bucket: %s
- Waste bucket: %s
- Raw transcript included: no
- Unknown private ecosystem names included: no

Use the included skills and hooks to keep Claude Code sessions scoped, reduce repeated reads, cap noisy output, and break retry loops.
`, report.AggregateEvent.ScoreBucket, report.AggregateEvent.WasteBucket)
}

func hooksJSON() string {
	return `{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "python3 \"${CLAUDE_PLUGIN_ROOT}/scripts/claude-analyzer-hook.py\"",
            "timeout": 5
          }
        ]
      }
    ]
  }
}
`
}

func hookScript() string {
	return `#!/usr/bin/env python3
import json
import re
import sys

try:
    payload = json.load(sys.stdin)
except Exception:
    sys.exit(0)

tool_input = payload.get("tool_input") or {}
command = tool_input.get("command") or ""
if not command:
    sys.exit(0)

warnings = []
if re.search(r"\b(cat|sed|nl|bat)\b", command) and not re.search(r"\b(head|tail|rg|grep|awk|jq)\b", command):
    warnings.append("This looks like a broad file read. Prefer targeted rg/head/tail output or summarize the file state once.")
if re.search(r"\b(find|ls -R|tree)\b", command):
    warnings.append("This command can create large output. Prefer rg --files with a narrow pattern.")
if re.search(r"\b(go test|npm test|pytest|cargo test)\b", command) and not re.search(r"\b(2>&1|tail|head|tee|--quiet|-q)\b", command):
    warnings.append("Test output can bloat context. Pipe to tail/head or run a narrower test first.")

if not warnings:
    sys.exit(0)

print(json.dumps({
    "hookSpecificOutput": {
        "hookEventName": "PreToolUse",
        "permissionDecision": "ask",
        "permissionDecisionReason": "Claude Analyzer: " + " ".join(warnings)
    }
}))
`
}

func statusCommand(report analyzer.Report) string {
	findings := strings.Join(sourceSummary(report).FindingIDs, ", ")
	if findings == "" {
		findings = "baseline-hygiene"
	}
	return fmt.Sprintf(`---
description: Show the Claude Analyzer session hygiene summary generated from the paid scan.
---

# Claude Analyzer Status

Report the current workflow hygiene posture in one terse line:

CTX discipline: watch | findings: %s | action: compact after pivots, cap shell output, avoid repeated reads.
`, findings)
}

func sessionHygieneSkill(report analyzer.Report) string {
	return fmt.Sprintf(`---
description: Use when a Claude Code session changes task type, grows context quickly, or needs a compact/clear decision.
---

# Session Hygiene

This plugin was generated from deterministic Claude Analyzer metrics.

Rules:

1. Keep debugging, architecture, and implementation in separate sessions when the task pivots.
2. Suggest /compact before a major subsystem pivot or after a long failed debugging branch.
3. Suggest /clear when the current context is dominated by stale assumptions.
4. Before rereading files, summarize what is already known and state the missing fact.
5. Keep advice short and operational.

Generated score bucket: %s
Generated waste bucket: %s
`, report.AggregateEvent.ScoreBucket, report.AggregateEvent.WasteBucket)
}

func retrievalHygieneSkill() string {
	return `---
description: Use when Claude Code is repeatedly reading the same files or running broad searches.
---

# Retrieval Hygiene

When inspecting code:

1. Prefer rg and targeted symbol/file searches before broad reads.
2. Read the narrowest range that can answer the question.
3. After reading a file, keep a short file-state summary before deciding to reread it.
4. If the same file appears again, state what new fact is needed before reading.
5. Avoid dumping entire files unless the file is small and central to the task.
`
}

func outputBudgetSkill() string {
	return `---
description: Use when shell, test, grep, build, or tool output may become large.
---

# Output Budget

Before running commands likely to print large output:

1. Prefer quiet flags, focused tests, and specific paths.
2. Pipe noisy commands through tail, head, rg, jq, or sed -n with a clear bound.
3. Capture full logs to a file only when needed, then inspect focused excerpts.
4. Never paste unbounded command output into context.
5. For repeated failing tests, stop after the second similar failure and inspect the invariant.
`
}

func retryBreakerSkill() string {
	return `---
description: Use after repeated failed commands, failed edits, or repeated attempts that do not change the failure mode.
---

# Retry Breaker

When the same failure repeats:

1. Stop after two similar failures.
2. State the invariant: what did not change across attempts.
3. Reduce the scope to the smallest failing command or file.
4. Re-read only the evidence needed for the next hypothesis.
5. Ask whether to compact or split the task if the session has drifted.
`
}
