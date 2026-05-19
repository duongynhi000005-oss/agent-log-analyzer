# Intent (intent)

- Status: research_needed
- Category: intent-spec / intent-driven development tool (canonical product not confirmed)
- Competitor priority: 15
- Official repository: (not confirmed — multiple unrelated repos use `intent` as a name)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["Intent", "intent"]

## Markers (public-source only)

### config_dir
- (not confirmed — candidate `.intent/` not verified)

### config_file
- (not confirmed)

### package_manifest
- (not confirmed)

### command_name
- (not confirmed; do NOT allowlist a bare `intent` binary — name is far too generic)

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
- Android `Intent` class / `Intent` filter — common Android development term; will match across millions of Android-developer transcripts.
- The English noun "intent" — common word.
- React/Web `intent` attributes on UI components.
- "User intent" / "intent classification" — common ML/NLP terminology.
- ANY bare mention of "intent" or "Intent" must NOT trigger this detector alone.

## Confidence wiring
- **High**: cannot define until a tool-specific config dir/file or slash command is recorded.
- **Medium**: n/a.
- **Low**: refuse low-confidence detection on the bare word.

## Source references (citations)
- _None — this entry is `research_needed` pending identification of the canonical repository._

## Open questions / what's missing
1. The brief lists "Intent" without a repo URL. The name collides with too many unrelated concepts.
2. Without a distinctive tool-specific marker (e.g., a `.intent/` directory with a known filename inside it), this detector cannot ship.
3. **A-04 trigger**: scope conversation needed — this is one of the highest-risk entries for false positives.
