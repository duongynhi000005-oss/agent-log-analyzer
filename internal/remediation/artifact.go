package remediation

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/robertdouglass/claude-log-analyzer/internal/analyzer"
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
	SchemaVersion          string               `json:"schema_version"`
	Generator              string               `json:"generator"`
	PluginName             string               `json:"plugin_name"`
	PluginVersion          string               `json:"plugin_version"`
	GeneratedAt            time.Time            `json:"generated_at"`
	Source                 SourceSummary        `json:"source"`
	Customizations         []Customization      `json:"customizations"`
	VettedRecommendations  []ToolRecommendation `json:"vetted_recommendations"`
	RequiredAcknowledgment string               `json:"required_acknowledgment"`
	Files                  []File               `json:"files"`
	Install                Install              `json:"install"`
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

type ToolRecommendation struct {
	ID                string `json:"id"`
	Category          string `json:"category"`
	Why               string `json:"why"`
	InstallCommand    string `json:"install_command"`
	RequiredBinary    string `json:"required_binary,omitempty"`
	BinaryInstallHint string `json:"binary_install_hint,omitempty"`
	Source            string `json:"source"`
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
	recommendations := toolingRecommendations(report)
	acknowledgment := liabilityAcknowledgment()
	files := baseFiles(report, pluginName, recommendations, acknowledgment)
	customizations := customizationPlan(report)
	files = append(files, customizationFiles(customizations)...)
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	artifact := Artifact{
		SchemaVersion:          "2026-05-18",
		Generator:              "claude-log-analyzer/remediation@" + GeneratorVersion,
		PluginName:             pluginName,
		PluginVersion:          pluginVersion(report),
		GeneratedAt:            generatedAt,
		Source:                 sourceSummary(report),
		Customizations:         customizations,
		VettedRecommendations:  recommendations,
		RequiredAcknowledgment: acknowledgment,
		Files:                  files,
	}
	artifact.Install = installInstructions(pluginName, options.ArtifactURL)
	return artifact
}

func baseFiles(report analyzer.Report, pluginName string, recommendations []ToolRecommendation, acknowledgment string) []File {
	manifest := map[string]any{
		"$schema":     "https://json.schemastore.org/claude-code-plugin-manifest.json",
		"name":        pluginName,
		"description": "Benchmark-backed Claude Code token hygiene workflow generated from a Claude Analyzer report.",
		"version":     pluginVersion(report),
		"author": map[string]string{
			"name": "Claude Log Analyzer",
		},
		"keywords": []string{"claude-code", "tokens", "context", "profiler", "benchmark", "tool-output"},
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
			Path:    "WAIVER.md",
			Mode:    "0644",
			Content: waiverFile(acknowledgment),
		},
		{
			Path:    "commands/claude-analyzer-status.md",
			Mode:    "0644",
			Content: statusCommand(report),
		},
		{
			Path:    "commands/claude-analyzer-tooling.md",
			Mode:    "0644",
			Content: toolingCommand(recommendations),
		},
		{
			Path:    "commands/agent-analyzer-proof.md",
			Mode:    "0644",
			Content: proofCommand(report),
		},
		{
			Path:    "commands/agent-analyzer-review.md",
			Mode:    "0644",
			Content: reviewCommand(report),
		},
		{
			Path:    "commands/agent-analyzer-status.md",
			Mode:    "0644",
			Content: statusCommand(report),
		},
		{
			Path:    "commands/agent-analyzer-tooling.md",
			Mode:    "0644",
			Content: toolingCommand(recommendations),
		},
		{
			Path:    "agents/token-hygiene-reviewer.md",
			Mode:    "0644",
			Content: tokenHygieneReviewerAgent(report),
		},
		{
			Path:    "skills/codebase-navigation/SKILL.md",
			Mode:    "0644",
			Content: codebaseNavigationSkill(report),
		},
		{
			Path:    "skills/session-hygiene/SKILL.md",
			Mode:    "0644",
			Content: sessionHygieneSkill(report),
		},
		{
			Path:    "skills/tooling-setup/SKILL.md",
			Mode:    "0644",
			Content: toolingSetupSkill(recommendations),
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
			Files:    []string{"skills/output-budget/SKILL.md", "skills/tooling-setup/SKILL.md"},
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

func toolingRecommendations(report analyzer.Report) []ToolRecommendation {
	seen := map[string]bool{}
	var out []ToolRecommendation
	add := func(rec ToolRecommendation) {
		if rec.ID == "" || seen[rec.ID] {
			return
		}
		seen[rec.ID] = true
		out = append(out, rec)
	}
	findingIDs := map[string]bool{}
	for _, finding := range report.Findings {
		findingIDs[finding.ID] = true
	}

	add(ToolRecommendation{
		ID:             "agent-analyzer-workflow",
		Category:       "built_in_plugin_workflow",
		Why:            "Core recommendation: Agent Analyzer guidance reduced estimated tokens by 12,370, tool-output tokens by 12,698, Claude output tokens by 504, and published API-rate cost by 24.0% across three fresh noisy-repo runs while preserving the quality gate.",
		InstallCommand: "Included in this plugin. Start with /agent-analyzer-review, then apply the generated output-budget, retrieval-hygiene, and session-hygiene skills before installing third-party tools.",
		Source:         "docs/benchmarks/repeated-benchmark-suite.md",
	})

	if findingIDs["tool_output_bloat"] || findingIDs["context_growth_spikes"] {
		add(ToolRecommendation{
			ID:             "context-mode",
			Category:       "context_defense",
			Why:            "Conditional reducer: repeated runs cut estimated tokens by 12,359, tool-output tokens by 13,257, and published API-rate cost by 20.4%, but visible output rose on average. Use only for tool-output/context bloat, not as a generic output or reasoning-token reducer.",
			InstallCommand: "/plugin marketplace add mksglu/context-mode\n/plugin install context-mode@context-mode\n/reload-plugins\n/context-mode:ctx-doctor",
			Source:         "https://github.com/mksglu/context-mode",
		})
		add(ToolRecommendation{
			ID:                "rtk",
			Category:          "advanced_shell_compression",
			Why:               "Conditional reducer: explicit RTK runs cut estimated tokens by 12,446, tool-output tokens by 12,716, and published API-rate cost by 18.2%. Use explicit commands first; do not enable global hooks until the user accepts the higher-risk waiver.",
			InstallCommand:    "brew install rtk\n# Start with explicit commands such as: rtk go test ./...\n# Enable hooks with `rtk init -g` only after reviewing the waiver and confirming shell rewriting is acceptable.",
			RequiredBinary:    "rtk",
			BinaryInstallHint: "macOS: brew install rtk. Linux/macOS fallback: curl -fsSL https://raw.githubusercontent.com/rtk-ai/rtk/refs/heads/master/install.sh | sh",
			Source:            "https://github.com/rtk-ai/rtk",
		})
		add(ToolRecommendation{
			ID:                "squeez",
			Category:          "explicit_shell_compression",
			Why:               "Conditional reducer: repeated Squeez runs cut estimated tokens by 8,471, tool-output tokens by 8,917, and published API-rate cost by 12.1%. Use it for noisy shell/log tasks, not as a general reasoning or output-token reducer.",
			InstallCommand:    "Review the current Squeez installation instructions, then run noisy commands through Squeez explicitly before considering shell integration.",
			RequiredBinary:    "squeez",
			BinaryInstallHint: "Install only from the current upstream instructions after reviewing the source.",
			Source:            "https://github.com/claudioemmanuel/squeez",
		})
	}

	if findingIDs["repeated_file_reads"] {
		add(ToolRecommendation{
			ID:                "grepai",
			Category:          "local_semantic_retrieval",
			Why:               "Conditional reducer: path-constrained grepai runs cut estimated tokens by 14,567, tool-output tokens by 15,571, and published API-rate cost by 14.5%. Keep limits and path filters tight; it does not directly reduce visible output or reasoning tokens.",
			InstallCommand:    "brew install yoanbernabeu/tap/grepai\ngrepai init\ngrepai watch",
			RequiredBinary:    "grepai",
			BinaryInstallHint: "Requires an embedding provider such as Ollama; install with curl script only after reviewing the GitHub source.",
			Source:            "https://github.com/yoanbernabeu/grepai",
		})
		add(ToolRecommendation{
			ID:                "semble",
			Category:          "path_limited_semantic_retrieval",
			Why:               "Positive reducer on this fixture: repeated Semble runs cut estimated tokens by 16,301, tool-output tokens by 16,060, Claude output by 480, and published API-rate cost by 41.5%. Use for bounded code retrieval where indexing/search output replaces repeated file reads.",
			InstallCommand:    "Review the current Semble installation instructions, index the target repository, and run bounded path-limited searches before reading files.",
			RequiredBinary:    "semble",
			BinaryInstallHint: "Install only from the current upstream instructions after reviewing the source and local data/indexing behavior.",
			Source:            "https://github.com/MinishLab/semble",
		})
	}

	if len(out) == 0 {
		add(ToolRecommendation{
			ID:             "agent-analyzer-workflow",
			Category:       "built_in_plugin_workflow",
			Why:            "No high-confidence third-party reducer matched this report. Keep the built-in Agent Analyzer workflow and avoid unproven installs.",
			InstallCommand: "Included in this plugin. Run /agent-analyzer-review and use the generated hygiene skills before adding tools.",
			Source:         "Claude Analyzer deterministic fallback",
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
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

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
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
	prompt := "Install the generated Claude Analyzer optimization plugin for this session. Explain that it contains benchmark-backed workflow guidance and only conditional reducers that survived repeated runs. Summarize the waiver and ask for approval before executing it. Do not print plugin archive contents.\n\n```sh\n" + command + "\n```"
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

Use the included skills and commands to run the practices that survived repeated benchmarks: scoped reads, output-budgeted commands, retry breaks, and session hygiene. The paid pack no longer recommends telemetry-only tools or negative benchmark tools as reducers.

Start with /agent-analyzer-status, /agent-analyzer-review, and /agent-analyzer-tooling. For benchmark or proof work, run /agent-analyzer-proof before claiming savings.
`, report.AggregateEvent.ScoreBucket, report.AggregateEvent.WasteBucket)
}

func waiverFile(acknowledgment string) string {
	return `# Required Acknowledgment

` + acknowledgment + `

This remediation pack may ask Claude Code to recommend or install third-party software only when the matching benchmarked finding exists. Those tools can execute code or access local/project data according to their own permissions.

Before installing anything:

1. Review every recommended tool and its source.
2. Confirm the command matches your operating system and package manager.
3. Confirm you have backups or version control for any repository Claude may modify.
4. Approve each installation separately.
5. Stop if Claude proposes an unvetted source, a destructive command, or a credential change you do not understand.

Claude Analyzer is not responsible for damage, data loss, credential exposure, billing impact, or other consequences caused by Claude Code, recommended tools, package managers, language servers, plugins, MCP servers, or user-approved commands.
`
}

func toolingCommand(recommendations []ToolRecommendation) string {
	return `---
description: Review the generated benchmark-backed token-saving setup recommendations.
---

# Claude Analyzer Tooling Setup

Read WAIVER.md first. Do not install anything until the user explicitly acknowledges the waiver and approves each command.

Only install recommendations listed below. Do not add ccusage, ccstatusline, claude-context, Probe, Caveman, claude-rlm, or claude-token-efficient as token-saving reducers for this workflow; the benchmark either classified them as telemetry, negative, harness-specific, or too small for the default pack.

Recommended actions:

` + recommendationMarkdown(recommendations) + `

Procedure:

1. Inspect the repository language stack and confirm each recommendation still applies.
2. Confirm the matching finding exists before installing a conditional tool.
3. Ask before installing each binary or plugin.
4. After installing plugins, run /reload-plugins.
5. If any tool source differs from the recommendation, stop and ask the user.
`
}

func codebaseNavigationSkill(report analyzer.Report) string {
	return fmt.Sprintf(`---
description: Use when Claude needs to understand a large or unfamiliar codebase without wasting context.
---

# Codebase Navigation

This skill follows Anthropic's large-codebase guidance: make the codebase navigable before adding more automation.

Rules:

1. Keep root CLAUDE.md lean: architecture map, critical commands, and gotchas only.
2. Prefer subdirectory CLAUDE.md files for local build/test conventions.
3. Start Claude in the relevant subdirectory when the task has a clear scope.
4. Build or update a concise codebase map when top-level folders are not self-explanatory.
5. Prefer existing symbol navigation or targeted rg searches before broad grep/read loops.
6. Avoid adding MCP servers merely because they exist; extra schemas and retrieval calls can increase context unless they replace repeated reads.

Generated score bucket: %s
Generated waste bucket: %s
`, report.AggregateEvent.ScoreBucket, report.AggregateEvent.WasteBucket)
}

func toolingSetupSkill(recommendations []ToolRecommendation) string {
	return `---
description: Use when setting up benchmark-backed token-saving tools.
---

# Tooling Setup

Install only with explicit user approval.

` + recommendationMarkdown(recommendations) + `

Installation discipline:

1. Read WAIVER.md to the user in summary form and get explicit acceptance before proceeding.
2. Verify the current OS, package manager, and existing binaries.
3. Use explicit commands before enabling hooks or automatic shell integration.
4. Run /reload-plugins after plugin installation.
5. If a recommended binary is already installed, do not reinstall it.
6. If a repository has custom tooling, prefer its checked-in setup docs over generic install commands.
7. Treat telemetry-only tools as measurement aids, not as reducers.
`
}

func recommendationMarkdown(recommendations []ToolRecommendation) string {
	var b strings.Builder
	for _, rec := range recommendations {
		b.WriteString("- ")
		b.WriteString(rec.ID)
		b.WriteString(" (")
		b.WriteString(rec.Category)
		b.WriteString("): ")
		b.WriteString(rec.Why)
		b.WriteString("\n")
		b.WriteString("  Install: `")
		b.WriteString(rec.InstallCommand)
		b.WriteString("`\n")
		if rec.RequiredBinary != "" {
			b.WriteString("  Required binary: `")
			b.WriteString(rec.RequiredBinary)
			b.WriteString("`\n")
		}
		if rec.BinaryInstallHint != "" {
			b.WriteString("  Binary install hint: `")
			b.WriteString(rec.BinaryInstallHint)
			b.WriteString("`\n")
		}
		b.WriteString("  Source: ")
		b.WriteString(rec.Source)
		b.WriteString("\n")
	}
	return b.String()
}

func liabilityAcknowledgment() string {
	return "I understand that Claude Analyzer provides deterministic analysis and vetted setup recommendations, but any installation or code change is executed by Claude Code, my package manager, or third-party tools with my approval and at my own risk."
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

func reviewCommand(report analyzer.Report) string {
	findings := strings.Join(sourceSummary(report).FindingIDs, ", ")
	if findings == "" {
		findings = "baseline-hygiene"
	}
	return fmt.Sprintf(`---
description: Review the current Claude Code session for avoidable token waste before continuing.
---

# Agent Analyzer Review

Use the token-hygiene-reviewer agent if available. Keep the review short and evidence-based.

Checklist:

1. Name the dominant generated finding set: %s.
2. Identify any current repeated file reads, noisy tool output, retry loops, or context-growth pivots.
3. Recommend one next action that preserves task quality while reducing avoidable context.
4. If optional tools are suggested, route to /agent-analyzer-tooling and ask before installing anything.

Do not claim token savings from this session unless a before/after benchmark has been measured by Claude Analyzer.
`, findings)
}

func proofCommand(report analyzer.Report) string {
	return fmt.Sprintf(`---
description: Explain what proof is required before claiming the plugin reduces token waste.
---

# Agent Analyzer Proof

This plugin was generated from a sanitized Claude Analyzer report. It recommends the workflow and conditional tools that survived repeated benchmarking, but new savings claims still require a controlled before/after benchmark.

Do not claim token savings until Claude Analyzer has measured both baseline and optimized logs. Name the token category: input/context, tool-output, visible output, cached input, telemetry-only, or reasoning tokens when the harness exposes reasoning usage.

Current generated evidence:

- Efficiency score bucket: %s
- Waste bucket: %s
- Repeated Agent Analyzer benchmark: 12,370 fewer estimated tokens, 12,698 fewer tool-output tokens, 504 fewer Claude output tokens, and 24.0%% lower published API-rate cost across three fresh runs.
- Scale rule: for comparable work, monthly savings = baseline monthly API-equivalent spend * 23.986%%. A team spending $5,000/month on similar Claude Sonnet coding work would save about $1,199/month.
- Raw transcript included in plugin: no
- Unknown private ecosystem names included in plugin: no

Before making a public claim:

1. Use the same task prompt and same starting commit for baseline and optimized Claude Code -p runs.
2. Analyze both logs with Claude Analyzer.
3. Compare total estimated tokens, input/context movement where available, visible output tokens where available, reasoning tokens where available, cached input where available, avoidable waste range, rereads, noisy tool output, retry loops, context-growth spikes, and task quality.
4. Translate token deltas to cost only with the published rate card for the exact model and token category. Keep API-rate estimates separate from Claude Code or Codex native billing.
5. Publish only sanitized reports and methodology.
6. Explain a null result honestly if the optimized run does not improve measured waste, or if an output-only reduction is offset by higher context, reasoning, or published API cost.
`, report.AggregateEvent.ScoreBucket, report.AggregateEvent.WasteBucket)
}

func tokenHygieneReviewerAgent(report analyzer.Report) string {
	findings := strings.Join(sourceSummary(report).FindingIDs, ", ")
	if findings == "" {
		findings = "baseline-hygiene"
	}
	return fmt.Sprintf(`---
name: token-hygiene-reviewer
description: Reviews Claude Code plans and session traces for avoidable token waste while preserving task quality.
---

You are a token-hygiene reviewer for Claude Code sessions.

Generated Claude Analyzer finding set: %s.

Review rules:

1. Prioritize task completion quality over token reduction.
2. Call out repeated reads, large unbounded command output, retry loops, and context pivots with concrete evidence.
3. Recommend the smallest workflow change that reduces avoidable context.
4. Prefer built-in shell discipline and targeted reads before third-party tools.
5. Distinguish input/context, tool-output, visible output, cached input, telemetry-only, and reasoning-token effects. Terse prose is not proof of lower full-session cost.
6. Do not install software or edit project files.
7. Do not claim savings unless Claude Analyzer measured both baseline and optimized logs.

Return:

- quality risk
- avoidable-waste risk
- one recommended next action
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
2. Do not use Glob, find, tree, or broad directory listing when failing tests or exact terms already identify candidate paths.
3. If a package path is known, inspect that package directly instead of discovering the whole repository.
4. Read the narrowest range that can answer the question.
5. After reading a file, keep a short file-state summary before deciding to reread it.
6. If the same file appears again, state what new fact is needed before reading.
7. Avoid dumping entire files unless the file is small and central to the task.
`
}

func outputBudgetSkill() string {
	return `---
description: Use when shell, test, grep, build, or tool output may become large.
---

# Output Budget

Before running commands likely to print large output:

1. Prefer quiet flags, focused tests, and specific paths.
2. Run the narrowest relevant test target while debugging, then run the full suite once at the end.
3. For Go projects, start with package-scoped tests such as: go test ./internal/foo ./internal/bar.
4. For final Go verification, prefer: go test ./... >/tmp/go-test.log && echo "go test ./... passed" || { grep -A3 -B2 -E "FAIL|--- FAIL" /tmp/go-test.log; exit 1; }.
5. Pipe noisy commands through tail, head, rg, jq, or sed -n with a clear bound.
6. Capture full logs to a file only when needed, then inspect focused excerpts.
7. Never paste unbounded command output into context.
8. For repeated failing tests, stop after the second similar failure and inspect the invariant.
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
