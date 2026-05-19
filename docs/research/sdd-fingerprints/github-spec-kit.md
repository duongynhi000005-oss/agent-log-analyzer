# GitHub Spec Kit (github_spec_kit)

- Status: verified
- Category: spec-driven-development scaffolding kit
- Competitor priority: 2
- Official repository: https://github.com/github/spec-kit
- Official docs: https://github.com/github/spec-kit#readme
- Release / package source: https://pypi.org/project/specify-cli/ (the kit installs via `uvx --from git+https://github.com/github/spec-kit.git specify init`) and GitHub releases at https://github.com/github/spec-kit/releases
- Aliases: ["Spec Kit", "spec-kit", "Specify CLI", "specify"]

## Markers (public-source only)

### config_dir
- `.specify/` — directory created by `specify init` containing scripts, memory, and templates. Documented in `github/spec-kit` README and in the project's quickstart (https://github.com/github/spec-kit/blob/main/README.md).
- `.specify/memory/` — persisted constitution and project memory created during init.
- `.specify/templates/` — Markdown templates used by the `/speckit.*` slash commands.

### config_file
- `.specify/memory/constitution.md` — project constitution file produced by `/speckit.constitution`. Documented at https://github.com/github/spec-kit#workflow.
- `.specify/scripts/bash/setup-plan.sh` and `.specify/scripts/powershell/setup-plan.ps1` — generated runner scripts (per the repo's `scripts/` layout).

### package_manifest
- PyPI `specify-cli` (https://pypi.org/project/specify-cli/) — Python CLI distributed via `uvx`/`pipx`.

### command_name
- `specify` — primary CLI binary (`specify init`, `specify check`). Documented in the README under "Quickstart".

### slash_command
- `/speckit.constitution`
- `/speckit.specify`
- `/speckit.clarify`
- `/speckit.plan`
- `/speckit.tasks`
- `/speckit.analyze`
- `/speckit.implement`
- `/speckit.checklist`

Namespace is `/speckit.<verb>` (one word, no hyphen, distinct from Spec Kitty's `/spec-kitty.<verb>`). Documented at https://github.com/github/spec-kit#commands.

### mcp_server_name
- (none — Spec Kit operates via slash commands and a CLI; no MCP server is shipped)

### skill_name
- (no first-party Claude Code skill manifest known)

### plugin_manifest
- (none documented)

### cli_binary
- `specify` — allowlisted binary.

### cli_version_probe
- args: `["--version"]`
- expected output pattern: `specify[^\d]*(\d+\.\d+)` → bucket `MAJOR.MINOR`.

## Negative-test markers (must NOT trigger this detector alone)
- `Spec Kitty` / `spec-kitty` / `.kittify/` (different product — see `spec-kitty.md`).
- `OpenSpec` / `openspec` / `openspec/` (different product — see `openspec.md`).
- Generic `specs/`, `plan.md`, `tasks.md` — these names are produced under `.specify/` by Spec Kit but are also produced by Spec Kitty and OpenSpec; alone they are ambiguous. Require `.specify/` directory presence or a `/speckit.*` slash command to disambiguate.
- The unrelated `specify` Python typing library (different package, no `.specify/` dir).

## Confidence wiring
- **High**: `.specify/` config_dir present **and** any of (`.specify/memory/constitution.md`, a `/speckit.*` slash command in transcript, `specify` CLI installed).
- **Medium**: exactly one of (`/speckit.*` slash command in transcript, `specify` CLI installed) without `.specify/` corroboration.
- **Low**: free-text mention of `(?i)\bgithub spec kit\b` or `(?i)\bspec[\s-]?kit\b` without artifact or CLI evidence.

## Source references (citations)
- https://github.com/github/spec-kit — official repository (README documents `specify init`, `.specify/` layout, and the `/speckit.*` slash commands).
- https://github.com/github/spec-kit/blob/main/README.md — `Quickstart` and `Commands` sections list all eight `/speckit.*` slash commands.
- https://pypi.org/project/specify-cli/ — package source.
- https://github.com/github/spec-kit/tree/main/templates — template files emitted into `.specify/templates/` by `specify init`.
