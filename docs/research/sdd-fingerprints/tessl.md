# Tessl (tessl)

- Status: research_needed
- Category: spec-driven-development platform (Tessl is a hosted vendor; CLI / artifact surface not fully confirmed from public source as of cutoff)
- Competitor priority: 18
- Official repository: (Tessl is a hosted commercial product; public open-source repo not confirmed)
- Official docs: https://www.tessl.io/ (vendor home; documentation linked from there but exact CLI / config layout is not exhaustively published)
- Release / package source: (not confirmed)
- Aliases: ["Tessl", "tessl", "Tessl.io"]

## Markers (public-source only)

### config_dir
- (not confirmed — candidate `.tessl/` not verified from public source)

### config_file
- (not confirmed)

### package_manifest
- (not confirmed — Tessl may publish npm/PyPI packages but exact names are not confirmed from this research pass)

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
- The vendor name "Tessl" is reasonably distinctive (low collision risk), but bare mentions still warrant Low confidence only.
- Avoid matching `tessellation` / `tessellate` substrings.

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: free-text mention of `(?i)\btessl\.io\b` URL.
- **Low**: bare `(?i)\bTessl\b` mention.

## Source references (citations)
- https://www.tessl.io/ — vendor home page (confirmed exists; deeper artifact docs not yet pulled into this research pass).

## Open questions / what's missing
1. Tessl appears to be a hosted product; need to confirm whether there is any local repo artifact (CLI, config file, npm package) that the static analyzer can detect.
2. If Tessl is fully hosted with no local artifact, treat like Cognition/Devin: text-mention markers only, and surface to the user as A-04 candidate.
3. **A-04 trigger**: scope conversation needed before WP05 implements a detector.
