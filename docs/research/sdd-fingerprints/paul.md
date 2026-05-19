# PAUL (paul)

- Status: research_needed
- Category: spec-driven-development assistant (exact expansion unconfirmed)
- Competitor priority: 12
- Official repository: (not confirmed from public source as of cutoff)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["PAUL", "paul"]

## Markers (public-source only)

### config_dir
- (not confirmed)

### config_file
- (not confirmed)

### package_manifest
- (not confirmed)

### command_name
- (not confirmed — candidate `paul` binary, but **far too generic** to allowlist without disambiguation)

### slash_command
- (not confirmed)

### mcp_server_name
- (not confirmed)

### skill_name
- (not confirmed)

### plugin_manifest
- (not confirmed)

### cli_binary
- (not confirmed; do NOT allowlist a `paul` binary — name is a common given name and will collide with countless unrelated tools and personal scripts on user $PATH)

### cli_version_probe
- (not applicable)

## Negative-test markers (must NOT trigger this detector alone)
- The given name "Paul" — extremely common in commit messages, author fields, comments, and prose. MUST NOT match on the bare word.
- Acronym `PAUL` — could expand to many unrelated phrases.
- Any `paul` binary on $PATH could be a personal script.

## Confidence wiring
- **High**: cannot define until tool-specific markers (config dir/file or slash command) are recorded.
- **Medium**: n/a.
- **Low**: refuse low-confidence detection on the bare name "PAUL" / "paul".

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. The brief lists "PAUL" without expansion or repo URL.
2. The name is a strong false-positive risk; aggressive disambiguation will be required.
3. **A-04 trigger**: scope conversation needed before WP05 implements a detector.
