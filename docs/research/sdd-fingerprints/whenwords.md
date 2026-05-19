# whenwords (whenwords)

- Status: research_needed
- Category: behavior-spec / scenario-driven workflow (BDD-style "When ... " phrasing implied)
- Competitor priority: 14
- Official repository: (not confirmed from public source as of cutoff)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["whenwords", "WhenWords", "when-words"]

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
- BDD "When ... " Gherkin-style step text — common in every Cucumber / Behat / Behave / pytest-bdd project. MUST NOT match.
- The word `whenwords` is distinctive once written as one token, which helps; the negative case is when prose contains "when words" with a space, which is generic English.

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: n/a.
- **Low**: bare `(?i)\bwhenwords\b` (one token) — could be a reasonable Low signal once the canonical repo is identified.

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. The brief lists "whenwords" without a repo URL.
2. **A-04 trigger**: scope conversation needed before WP05 implements a detector.
