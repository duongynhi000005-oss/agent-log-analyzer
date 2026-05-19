# OpenSpec (openspec)

- Status: verified
- Category: spec-driven-development change/proposal workflow
- Competitor priority: 3
- Official repository: https://github.com/Fission-AI/OpenSpec
- Official docs: https://github.com/Fission-AI/OpenSpec#readme
- Release / package source: https://www.npmjs.com/package/@fission-ai/openspec (installed via `npx @fission-ai/openspec` or `pnpm dlx`)
- Aliases: ["OpenSpec", "openspec", "@fission-ai/openspec"]

## Markers (public-source only)

### config_dir
- `openspec/` — top-level directory at repo root that holds the OpenSpec workspace (created by `openspec init`). Documented in https://github.com/Fission-AI/OpenSpec#getting-started.
- `openspec/changes/` — per-change proposal directory.
- `openspec/specs/` — capability-spec directory (canonical specs).

### config_file
- `openspec/project.md` — OpenSpec project context document. Documented in the upstream README.
- `openspec/AGENTS.md` — agent instructions written by `openspec init`. Documented in the upstream README.
- `openspec/changes/<change-id>/proposal.md` and `tasks.md` — per-change artifacts (path pattern).

### package_manifest
- npm `@fission-ai/openspec` (https://www.npmjs.com/package/@fission-ai/openspec).

### command_name
- `openspec` — CLI from the npm package (`openspec init`, `openspec validate`, `openspec list`, `openspec show`, `openspec archive`). Documented in the README's "Commands" section.

### slash_command
- `/openspec-proposal` (or `/openspec.<verb>` family in some agent integrations) — documented in agent prompts emitted by `openspec init`; precise namespace varies by AI client. Cite https://github.com/Fission-AI/OpenSpec#ai-agent-integration.

### mcp_server_name
- (none — OpenSpec operates via the `openspec` CLI and agent prompts)

### skill_name
- (no first-party Claude Code skill manifest documented)

### plugin_manifest
- (none documented)

### cli_binary
- `openspec` — allowlisted binary.

### cli_version_probe
- args: `["--version"]`
- expected output pattern: `openspec[^\d]*(\d+\.\d+)` → bucket `MAJOR.MINOR`.

## Negative-test markers (must NOT trigger this detector alone)
- `Spec Kit` / `spec-kit` / `.specify/` (different product — see `github-spec-kit.md`).
- `Spec Kitty` / `spec-kitty` / `.kittify/` (different product — see `spec-kitty.md`).
- Generic `specs/`, `tasks.md`, `proposal.md` — these names exist inside `openspec/` but also inside Spec Kit's `.specify/`; require `openspec/` directory at repo root to disambiguate.
- The unrelated "OpenAPI Spec" abbreviation `OpenAPI spec` — distinct standard.

## Confidence wiring
- **High**: `openspec/` top-level directory present **and** any of (`openspec/AGENTS.md`, `openspec/project.md`, `openspec` CLI installed, `/openspec*` slash command in transcript).
- **Medium**: exactly one of (`openspec` CLI installed, `/openspec*` slash command in transcript) without `openspec/` corroboration.
- **Low**: free-text mention of `(?i)\bopenspec\b` without artifact or CLI evidence.

## Source references (citations)
- https://github.com/Fission-AI/OpenSpec — official repository (README documents `openspec init`, the `openspec/` directory layout, `project.md`, `AGENTS.md`, and the `openspec` CLI verbs).
- https://www.npmjs.com/package/@fission-ai/openspec — published npm package source.
- https://github.com/Fission-AI/OpenSpec/tree/main/templates — bundled templates that `openspec init` writes into the user's `openspec/` directory.
