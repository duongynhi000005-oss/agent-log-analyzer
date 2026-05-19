# fspec (fspec)

- Status: research_needed
- Category: feature-spec / functional-spec authoring tool
- Competitor priority: 13
- Official repository: (not confirmed — multiple unrelated repos use `fspec`; the brief's intended upstream is unverified as of cutoff)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["fspec", "f-spec", "feature-spec"]

## Markers (public-source only)

### config_dir
- (not confirmed)

### config_file
- (not confirmed)

### package_manifest
- (not confirmed)

### command_name
- (not confirmed)

### slash_command
- (not confirmed)

### mcp_server_name
- (not confirmed)

### skill_name
- (not confirmed)

### plugin_manifest
- (not confirmed)

### cli_binary
- (not confirmed; do NOT allowlist)

### cli_version_probe
- (not applicable)

## Negative-test markers (must NOT trigger this detector alone)
- File extension `.fspec` is used by unrelated tools (HDL specification, ad-hoc DSLs).
- Bare token `fspec` in source comments — could refer to a function-spec or filesystem-spec helper variable.
- `Fspec` Ruby class names — unrelated.

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: n/a.
- **Low**: n/a.

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. The brief lists "fspec" without a repo URL; several GitHub repos use the name.
2. **A-04 trigger**: scope conversation needed before WP05 implements a detector.
