# SDD Pilot (sdd_pilot)

- Status: research_needed
- Category: spec-driven-development pilot/assistant
- Competitor priority: 8
- Official repository: (not confirmed from public source as of cutoff — the brief lists "SDD Pilot" but no canonical repo URL has been verified)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["SDD Pilot", "sdd-pilot", "sdd_pilot"]

## Markers (public-source only)

### config_dir
- (not confirmed — candidate `.sdd-pilot/` based on naming convention; **not verified**)

### config_file
- (not confirmed)

### package_manifest
- (not confirmed)

### command_name
- (not confirmed — candidate `sdd-pilot` not allowlisted yet)

### slash_command
- (not confirmed)

### mcp_server_name
- (not confirmed)

### skill_name
- (not confirmed)

### plugin_manifest
- (not confirmed)

### cli_binary
- (not confirmed; do NOT allowlist until repo is identified)

### cli_version_probe
- (not applicable)

## Negative-test markers (must NOT trigger this detector alone)
- Bare phrase "SDD pilot" — could be a generic descriptor ("a spec-driven-development pilot project") rather than this tool.
- Acronym `SDD` alone — far too generic (software design document, sudden death disease, etc.).

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: n/a.
- **Low**: n/a.

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. The brief at `start-here.md` names "SDD Pilot" without a repo URL.
2. Several GitHub repos share variants of this name; the team needs to confirm which one is intended.
3. Without an official repo, this detector cannot be `verified` per FR-013/FR-014.
4. **A-04 trigger**: scope conversation needed with the user before WP05 implements a detector.
