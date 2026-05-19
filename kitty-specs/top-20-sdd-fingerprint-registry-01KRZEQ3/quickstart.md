# Quickstart — Adding or updating an SDD detector

This is the developer-facing flow for working on the SDD fingerprint registry.

## 1. Research the tool (public sources only)

Clone the tool's official repo into the workspace research area (never into the repo):

```sh
mkdir -p /Users/robert/code-analyzer-dev/claude-code-analyzer-20260519-082245-0QWuF7/research
cd /Users/robert/code-analyzer-dev/claude-code-analyzer-20260519-082245-0QWuF7/research
git clone <official-repo-url> <tool-id>
cd <tool-id>
# Read the docs, examine init/templates, observe generated files in a scratch dir.
```

Record fingerprints in a new file:

```sh
$EDITOR docs/research/sdd-fingerprints/<tool-id>.md
```

Use the template from `research.md` §R-01. Every marker MUST cite a public source.

## 2. Add or update the registry entry

Open `internal/analyzer/signatures/sdd_detectors.json` and add a `SDDDetector` object that matches `kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/contracts/sdd-detector.schema.json`.

Minimum fields:

- `id` — lowercase snake-case, e.g., `"github_spec_kit"`.
- `display_name` — human-friendly form, e.g., `"GitHub Spec Kit"`.
- `category` — e.g., `"sdd"`.
- `competitor_priority` — integer, 1 = highest.
- `status` — `"verified"` once at least one public source reference is recorded and at least one tool-specific marker exists.
- `source_references` — at least one for `verified`.
- `markers` — at least one. Use `source_class: "cli_binary"` with `binary: "<name>"` for CLI presence; use `source_class: "config_dir"` / `"config_file"` etc. with `pattern: "<regex>"` for textual matches.
- `confidence_rules` — at least one. Standard pattern:

```json
[
  { "confidence": "high",   "requires_distinct_classes": 2 },
  { "confidence": "medium", "requires_any_of": ["cli_binary", "slash_command", "mcp_server_name", "package_manifest"] },
  { "confidence": "low",    "requires_any_of": ["command_name"] }
]
```

## 3. Add a fixture and tests

Add `internal/analyzer/sdd/testdata/fixtures/<tool-id>.txt` carrying a scrubbed line or two that should trigger only this tool. Include in `evaluator_test.go`:

- A positive case: fixture → tool detected with expected confidence.
- A cross-negative case: fixture does NOT trigger any other detector.

For Spec Kitty / GitHub Spec Kit / OpenSpec specifically, ensure the 3×3 cross-negative matrix in `evaluator_test.go` covers the new entry per NFR-004.

## 4. Run the test suite

```sh
gofmt -w $(find . -name '*.go' -not -path './.git/*')
go test ./...
```

Then run the smoke (or document why blocked):

```sh
./scripts/smoke-local.sh
```

## 5. Verify privacy invariants

The serialization-leak test runs as part of `go test ./internal/analyzer/sdd/...` and asserts no forbidden string from `contracts/forbidden-strings.md` appears in any serialized `Report` or `AggregateEvent`.

If a new field is added to `EcosystemFingerprint`, the test must be extended with a corresponding canary value. The PR review checklist enforces this.

## 6. Verify CLI probe rules (if adding a `cli_binary` marker)

- Confirm the binary name is documented in the tool's public docs.
- Confirm `version_args` are safe (`--version`, `version`, or `-v` only; no flags involving paths, networks, auth, or servers).
- Confirm the binary's `--version` output is something `normalizeVersionBucket` can parse, or accept that `version_bucket` will be empty for this tool.
- Add the binary to the loader's allowlist check (automatic via the schema).

## 7. Commit and push

Follow the brief's commit / push instructions:

```sh
git switch -c codex/sdd-fingerprint-registry  # if not already on it
git status --short
git add internal/analyzer/sdd/ internal/analyzer/signatures/sdd_detectors.json \
        internal/analyzer/{ecosystem.go,types.go,registry.go} \
        docs/research/sdd-fingerprints/ docs/sdd-fingerprint-registry.md
git commit -m "Add SDD fingerprint registry"
git push -u origin codex/sdd-fingerprint-registry
```

## 8. Issue hygiene

Per the brief, comment on the relevant GitHub issues (#38, #42, #43, #44–#50, #66, #67) at the points named in the brief. Do not close any issue unless acceptance criteria are actually satisfied.

## What is OUT of scope for this quickstart

- Building dashboards or analytics UIs.
- Adding new third-party Go dependencies.
- Modifying existing legacy detection codepaths beyond the one `WorkflowFingerprints` field addition.
- Performing detection from private/unallowed CLI binary names.
