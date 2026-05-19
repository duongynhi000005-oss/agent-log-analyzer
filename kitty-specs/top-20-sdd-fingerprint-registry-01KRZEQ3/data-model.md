# Phase 1 — Data Model

**Mission**: `top-20-sdd-fingerprint-registry-01KRZEQ3`
**Date**: 2026-05-19

This document defines the typed values introduced by this mission. Each entry shows the Go type, the JSON shape (where serialized), and the invariants the type carries.

## Enumerations

### `sdd.SourceClass`

Go: `type SourceClass string`

Permitted values (FR-006):

| Value | Meaning |
| --- | --- |
| `config_dir` | Tool-specific directory (e.g., `.kittify/`, `.openspec/`). |
| `config_file` | Tool-specific file (e.g., `openspec.yaml`). |
| `package_manifest` | A package entry in `package.json` / `pyproject.toml` / `go.mod` / similar. |
| `command_name` | Free-text mention of a tool's CLI name in scrubbed input. Low-confidence on its own. |
| `slash_command` | A registered slash command identifier (`/<name>`). |
| `mcp_server_name` | A `mcp__<name>__` server identifier. |
| `skill_name` | A registered skill identifier. |
| `plugin_manifest` | A plugin manifest filename or id. |
| `cli_binary` | An installable CLI binary name probed via `LookPath`. |
| `cli_version_probe` | A safe `--version` / `version` probe of a `cli_binary`. |

**Invariants.** Any field of type `SourceClass` may only hold one of the values above. The loader rejects unknown values with a panic at startup (registry is embedded; bad data is a build-time bug, not a runtime one).

### `sdd.Confidence`

Go: `type Confidence string`

Values (FR-007): `high`, `medium`, `low`. Loader-validated.

### `sdd.Status`

Go: `type Status string`

Values (FR-015): `verified`, `research_needed`, `blocked`. Loader-validated.

## Registry record

### `sdd.SDDDetector`

```go
type SDDDetector struct {
    ID                 string           `json:"id"`
    DisplayName        string           `json:"display_name"`
    Aliases            []string         `json:"aliases,omitempty"`
    Category           string           `json:"category"`
    CompetitorPriority int              `json:"competitor_priority"`
    Status             Status           `json:"status"`
    SourceReferences   []SourceRef      `json:"source_references"`
    Markers            []Marker         `json:"markers"`
    ConfidenceRules    []ConfidenceRule `json:"confidence_rules"`
}

type SourceRef struct {
    Kind string `json:"kind"` // e.g., "official_repo", "docs", "release"
    URL  string `json:"url"`
    Note string `json:"note,omitempty"`
}

type Marker struct {
    SourceClass SourceClass `json:"source_class"`
    Pattern     string      `json:"pattern,omitempty"` // regex, used by all classes except cli_*
    Binary      string      `json:"binary,omitempty"`  // for cli_binary / cli_version_probe
    VersionArgs []string    `json:"version_args,omitempty"` // for cli_version_probe; defaults to ["--version"]
    Negative    bool        `json:"negative,omitempty"` // if true, presence VETOES this detector for this fixture
    Note        string      `json:"note,omitempty"`
}

type ConfidenceRule struct {
    Confidence    Confidence    `json:"confidence"`
    Requires      []SourceClass `json:"requires_any_of,omitempty"`        // at least one of these source classes matched
    RequiresAll   []SourceClass `json:"requires_all_of,omitempty"`        // all of these source classes matched
    RequiresCount int           `json:"requires_distinct_classes,omitempty"` // at least N distinct source classes matched
}
```

**Invariants.**

- `ID` is an allowlisted public identifier. Never derived from user input.
- `ID` matches `^[a-z][a-z0-9_]{2,}$`.
- `Status == "verified"` implies `len(SourceReferences) >= 1`.
- A detector with status `research_needed` MUST NOT participate in evaluation (the evaluator filters them out before scoring). Per C-001 this mission ships zero detectors in that state, but the schema supports it.
- The loader compiles every regex `Pattern` at startup; bad patterns panic at startup.
- A detector with a `cli_binary` marker MUST also list the binary on the allowlist (the loader enforces this).
- `Marker.Negative == true` means: if this marker's pattern matches the input, the entire detector emits no fingerprint regardless of other matches. Used for cross-negative rules where one tool's marker forbids another's detection.

## Report-emitted record

### `analyzer.EcosystemFingerprint`

```go
type EcosystemFingerprint struct {
    ID            string   `json:"id"`
    Confidence    string   `json:"confidence"`
    Sources       []string `json:"sources"`
    EvidenceCount int      `json:"evidence_count"`
    Active        bool     `json:"active,omitempty"`
    Installed     bool     `json:"installed,omitempty"`
    VersionBucket string   `json:"version_bucket,omitempty"`
}
```

**Invariants.**

- `ID` is exactly the detector `ID` from the registry (allowlisted, public, lowercase snake).
- `Sources` is sorted, deduplicated, and contains only valid `SourceClass` string values.
- `EvidenceCount` ≥ `len(Sources)` (a single source class may match multiple markers).
- `Confidence` is one of `high`, `medium`, `low`.
- `VersionBucket` is the output of `normalizeVersionBucket` only. Never raw `--version` output.
- `Installed == true` requires that the CLI probe's `LookPath` returned true for the detector's `cli_binary` marker.
- `Active == true` requires `Confidence == "high"` AND ≥1 runtime-touch source (slash_command, mcp_server_name, or cli_binary).
- Serialization: the type contains only the seven primitive fields above. No `interface{}`, no `map[string]string`, no free-text "evidence" field. This bounded shape is what makes NFR-001 enforceable structurally.

## Modified existing types

### `analyzer.Ecosystem`

Add one field (preserving existing fields per C-004):

```go
type Ecosystem struct {
    // ... existing fields unchanged ...
    WorkflowFingerprints []EcosystemFingerprint `json:"workflow_fingerprints,omitempty"`
}
```

**Invariant.** `WorkflowFingerprints` may be `nil` (omitted) or a sorted, deduplicated slice by `ID`. Sort key: `(competitor_priority asc, ID asc)` so the highest-priority tools surface first.

## CLI probe abstraction

### `sdd.CLIProbe`

```go
type CLIProbe interface {
    LookPath(name string) (found bool)
    Version(ctx context.Context, name string, args []string) (rawOutput string, ok bool)
}
```

**Implementations.**

- `RealProbe` — wraps `os/exec`. Hard rules per research R-03.
- `FakeProbe` — table-driven map of `binary → (found bool, version string)`. Used by every test.

**Invariant.** Callers may only pass `name` values that appear as a `cli_binary` marker in some loaded `SDDDetector`. The evaluator enforces this; ad-hoc callers do not exist outside the `sdd` package.

## Storage shape (embedded JSON)

`internal/analyzer/signatures/sdd_detectors.json` is a JSON array of `SDDDetector` records. Schema lives at [contracts/sdd-detector.schema.json](contracts/sdd-detector.schema.json).

The file is embedded via the existing `signatureFS` (`//go:embed signatures/*.json`) so it ships in the binary. No on-disk read paths are added.

## Unknown-count buckets

Existing fields `Ecosystem.UnknownMCPServerCount`, `UnknownSkillCount`, `UnknownPluginCount` are reused unchanged. The new evaluator does not introduce a fourth unknown-count bucket; tools not in the registry are simply absent from `WorkflowFingerprints`.

## Forbidden state

The data model deliberately omits any field that could leak private content:

- No `Path string` on `EcosystemFingerprint` (despite the underlying probe knowing the resolved path).
- No `RawVersionOutput string`.
- No `MatchedText string`.
- No `EvidenceLocations []string`.
- No `map[string]string` keyed by user input.
- No free-text `Note` on report-emitted records.

The corresponding `Marker.Pattern` regex is **input** (registry data), not **output** (report data); patterns never appear in the serialized report.
