# Agentic Code (agentic_code)

- Status: research_needed
- Category: agentic-coding workflow (canonical product not confirmed)
- Competitor priority: 19
- Official repository: (not confirmed — "agentic code" is a category/buzzword as well as a name; the specific tool intended by the brief is unverified)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["Agentic Code", "agentic-code"]

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
- The phrase "agentic code" / "agentic coding" — used generically across the entire agent-tooling industry (Claude Code, Cursor, Cline, Aider all describe themselves as "agentic coding"). MUST NOT match on the phrase alone.
- The category descriptor "Agentic Code" in marketing copy or blog posts.

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: n/a.
- **Low**: n/a — refuse low-confidence detection on the category phrase.

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. Like "Spec-Driven Develop", the name "Agentic Code" overlaps with a popular category descriptor.
2. The brief lists this as a tool; without a specific upstream repo or vendor URL the detector cannot ship.
3. **A-04 trigger**: strong candidate for scope conversation — name clashes with the category itself.
