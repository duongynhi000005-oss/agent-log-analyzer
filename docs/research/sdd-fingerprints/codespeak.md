# CodeSpeak (codespeak)

- Status: research_needed
- Category: natural-language-to-code spec workflow (canonical product not confirmed)
- Competitor priority: 20
- Official repository: (not confirmed from public source as of cutoff — several products use this name)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["CodeSpeak", "codespeak", "code-speak"]

## Markers (public-source only)

### config_dir
- (not confirmed — candidate `.codespeak/` not verified)

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
- The phrase "code speak" — could refer to programming language pedagogy, generic prose.
- Unrelated apps named CodeSpeak (there are voice-coding/accessibility tools and unrelated GitHub repos sharing the name).

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: n/a.
- **Low**: bare `(?i)\bcodespeak\b` (one token) once the canonical repo is identified — distinctive enough to be a reasonable Low signal.

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. The brief lists "CodeSpeak" without a repo URL.
2. **A-04 trigger**: scope conversation needed before WP05 implements a detector.
