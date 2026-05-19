# GSD (gsd)

- Status: research_needed
- Category: spec-driven-development (Get Spec Done / Get Stuff Done — exact expansion unconfirmed)
- Competitor priority: 6
- Official repository: (not confirmed from public source as of cutoff)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["GSD", "Get Spec Done", "gsd"]

## Markers (public-source only)

### config_dir
- (not confirmed — candidate `.gsd/` based on naming convention but **not verified**)

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
- (not confirmed; do NOT allowlist until repo is identified)

### cli_version_probe
- (not applicable until binary confirmed)

## Negative-test markers (must NOT trigger this detector alone)
- Bare `GSD` — extremely common acronym ("Get Shit Done", productivity framework, unrelated tools named GSD on npm/PyPI). Must NOT match on the three-letter token alone.
- The unrelated `gsd` npm package (file-system watcher) — different product.
- Any productivity-app mention of "GSD" — different domain.

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: n/a.
- **Low**: n/a (refuse low-confidence detection on a three-letter acronym).

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. The brief at `start-here.md` lists "GSD" without a repo URL. Several SDD-adjacent repos use the acronym; the team needs to confirm which one the brief intends.
2. Without an official repo, this detector cannot be `verified` per FR-013/FR-014.
3. **A-04 trigger**: Recommend the reviewer / user clarify the upstream source before WP05 implements the detector. Until then this entry is `research_needed`.
