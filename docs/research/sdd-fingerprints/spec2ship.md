# spec2ship (spec2ship)

- Status: research_needed
- Category: spec-driven-development ship workflow
- Competitor priority: 10
- Official repository: (not confirmed from public source as of cutoff)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["spec2ship", "spec-to-ship", "Spec2Ship"]

## Markers (public-source only)

### config_dir
- (not confirmed — candidate `.spec2ship/` based on naming convention; **not verified**)

### config_file
- (not confirmed)

### package_manifest
- (not confirmed)

### command_name
- (not confirmed — candidate `spec2ship` not allowlisted yet)

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
- The two-word phrase "spec to ship" — common English in shipping/logistics contexts.
- Mentions of "spec2ship" in unrelated CI/CD blog posts.

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: n/a.
- **Low**: free-text mention of `(?i)\bspec2ship\b` is a reasonable Low signal once the repo is identified, because the spelling is distinctive.

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. The brief lists "spec2ship" without a repo URL. Need to confirm.
2. **A-04 trigger**: scope conversation needed before WP05 implements a detector.
