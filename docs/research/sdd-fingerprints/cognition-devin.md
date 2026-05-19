# Cognition / Devin (cognition_devin)

- Status: verified
- Category: hosted autonomous software-engineering agent (Cognition's "Devin" product)
- Competitor priority: 16
- Official repository: (Devin is a hosted product — no open-source repo; vendor is Cognition AI)
- Official docs: https://docs.devin.ai/
- Release / package source: https://app.devin.ai (web app) — there is also the `@devin-ai/cli` style integration documented at https://docs.devin.ai/api-reference (subject to evolution)
- Aliases: ["Devin", "Devin AI", "Cognition", "Cognition Devin"]

## Markers (public-source only)

### config_dir
- (none on the user's local filesystem — Devin runs in Cognition's hosted sandbox, not the user's repo)

### config_file
- `.devin/` workspace files MAY be written when a user uses Devin's "Workspace" or "Repo" linking feature; documented under https://docs.devin.ai/repo-integration (subject to product evolution). Treat presence of `.devin/` as suggestive but not conclusive — confirm with a second marker.

### package_manifest
- (no PyPI/npm package known for Devin itself; the official Slack / IDE integrations are vendor-distributed)

### command_name
- (none on $PATH; do NOT allowlist a `devin` binary)

### slash_command
- (none in the Claude Code surface — Devin is invoked via its own chat UI / Slack app)

### mcp_server_name
- (none documented as `mcp__devin__*` in the public Cognition docs as of cutoff)

### skill_name
- (none — Devin is not packaged as a Claude Code skill)

### plugin_manifest
- (none)

### cli_binary
- (none reliably on $PATH; do NOT allowlist)

### cli_version_probe
- (not applicable)

## Negative-test markers (must NOT trigger this detector alone)
- The given name "Devin" — common first name; will collide with author fields and comments.
- "Cognition" as a generic word — common in psychology, ML, neuroscience contexts.
- Unrelated `devin` npm packages / vanity projects.

## Confidence wiring
- **High**: cannot reliably be High from a static analyzer because Devin runs server-side. Reserve High for the case where a user explicitly references Devin's web URL pattern (`https://app.devin.ai/sessions/<id>`) in committed text **plus** at least one second-class signal.
- **Medium**: free-text mention of `(?i)\bDevin AI\b` or `(?i)\bCognition Devin\b` or a `app.devin.ai` URL.
- **Low**: bare `(?i)\bDevin\b` mention — too noisy on its own.

## Source references (citations)
- https://www.cognition.ai/ — Cognition AI corporate site.
- https://docs.devin.ai/ — official product documentation.
- https://app.devin.ai/ — hosted application (URL pattern is the most fingerprintable signal).

## Open questions / what's missing
1. As a hosted product, Devin has no local CLI / config dir to fingerprint. Per FR-013, this might be the strongest A-04 candidate — the brief includes Devin but a static analyzer of a local repo cannot reliably detect a server-side product.
2. The current `verified` status reflects that the public-source citations (docs URL, app URL) are confirmed, but the detector itself will be limited to text-mention markers. WP05 implementers should not expect a config-dir / CLI High-confidence signal for this tool.
