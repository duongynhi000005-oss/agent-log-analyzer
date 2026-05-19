# Data Retention And Analytics

## Retention Classes

```text
raw uploaded logs:
  launch CLI flow: never received by the server
  legacy internal token flow: deleted by claude-analyzer-sweeper
  analytics: never

intermediate parsed transcript:
  launch CLI flow: local process memory only
  legacy internal token flow: worker memory or short-lived encrypted temp
  analytics: never

sanitized report JSON:
  local MVP: stored under /data/reports after explicit upload
  production: 15 minutes free, 24 hours paid artifact
  analytics: never as raw JSON

job metadata:
  local MVP: job JSON files
  production: short-lived metadata with TTL
  analytics: status/timing/error category only

aggregate analytics:
  local MVP: not exported
  production: allowlisted numeric/categorical events
  analytics: yes, no raw strings
```

## Operational Logging Allowlist

Allowed:

- `job_id`
- `request_id`
- `file_size_bucket`
- `parser_type`
- `analyzer_version`
- `status`
- `duration_ms`
- `queue_wait_ms`
- `worker_exit_reason`
- `error_category`
- `redaction_counts_by_type`

Forbidden:

- raw transcript text
- prompts
- tool output
- file contents
- secret values
- command arguments
- raw unknown MCP/plugin/skill names
- repo names
- usernames
- hostnames
- full file paths

## Aggregate Ecosystem Intelligence

Collect by default:

- known public workflow framework IDs
- known public MCP server IDs
- known public plugin IDs
- known public skill IDs
- OS category
- shell category
- package manager category
- counts and buckets

Unknown private names are counted, not stored:

```json
{
  "unknown_mcp_server_count": 3,
  "unknown_skill_count": 4,
  "unknown_plugin_count": 1
}
```

Exact unknown names require explicit opt-in.

## Upload Scope

Free scan analyzes exactly one Claude Code JSONL session selected by the local CLI. The server receives only the generated sanitized report JSON after the user has had a chance to inspect it.

Paid scan must use the same local-first model for at most the 100 most recent Claude Code JSONL sessions. Aggregate analytics from paid scans must still use the same allowlist: known public ecosystem IDs, counts, buckets, timing, parser status, and redaction totals. Raw logs, raw paths, unknown private names, and report JSON are not retained as analytics.

## SDD Fingerprint Privacy

Spec-driven-development (SDD) tool detection follows the same retention and
allowlist discipline as the rest of ecosystem analytics. See
[`sdd-fingerprint-registry.md`](sdd-fingerprint-registry.md) for the
maintainer-facing overview.

- **Bounded record shape.** Each entry in `Ecosystem.WorkflowFingerprints`
  is an `EcosystemFingerprint` carrying only seven primitive fields: `id`,
  `confidence`, `sources`, `evidence_count`, `active`, `installed`,
  `version_bucket`. No `map[string]string`, no free-text `note`, no raw
  path, no raw `--version` output, no `evidence_locations` slice. See
  [data-model.md](../kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/data-model.md#analyzerecosystemfingerprint)
  for the exact type and invariants.
- **`id` is allowlisted.** It is exactly the detector id from the registry
  (lowercase snake-case, public, never derived from user input).
- **Unknown names remain counts only (FR-011).** Unknown or private MCP /
  skill / plugin / tool names are counted in
  `Ecosystem.UnknownMCPServerCount`, `UnknownSkillCount`, and
  `UnknownPluginCount`. Their names, identifiers, and any derivable hashes
  do not appear in aggregate output or aggregate events. Tools not in the
  registry are simply absent from `WorkflowFingerprints`.
- **Version buckets, not version strings.** `version_bucket` is the output
  of `normalizeVersionBucket` only (e.g., `"1.2"`, `"unknown"`, or absent).
  The raw `--version` stdout/stderr is consumed in-package and discarded
  before serialization.
- **Build-time enforcement (NFR-001).** A serialization-leak test asserts
  that none of the 16 forbidden raw-string categories appears in any fully
  populated `Report` or `AggregateEvent` payload. Adding a new field to
  `EcosystemFingerprint` requires extending the leak test with a
  corresponding canary in the same PR.

### Forbidden raw-string categories (high level)

The 16 categories the leak test covers (full canary fixtures live in
[`contracts/forbidden-strings.md`](../kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/contracts/forbidden-strings.md)):

1. user prompts
2. task descriptions
3. raw transcript excerpts
4. raw tool inputs
5. raw tool outputs
6. raw file paths
7. repo URLs
8. branch names
9. usernames
10. hostnames
11. emails
12. session IDs
13. transcript paths
14. private MCP / skill / plugin names
15. raw `LookPath` / `which` paths
16. raw `--version` output / stable hashes of private strings
