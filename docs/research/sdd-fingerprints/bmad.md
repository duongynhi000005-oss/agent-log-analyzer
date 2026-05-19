# BMAD-METHOD (bmad)

- Status: verified
- Category: agent-team SDD framework (Breakthrough Method for Agile AI Driven Development)
- Competitor priority: 5
- Official repository: https://github.com/bmad-code-org/bmad-method
- Official docs: https://github.com/bmad-code-org/bmad-method#readme
- Release / package source: https://www.npmjs.com/package/bmad-method (installed via `npx bmad-method install`)
- Aliases: ["BMAD-METHOD", "BMAD", "bmad", "bmad-method"]

## Markers (public-source only)

### config_dir
- `.bmad-core/` — directory created by `npx bmad-method install`, containing agents, templates, checklists, and tasks. Documented at https://github.com/bmad-code-org/bmad-method#installation.
- `.bmad-core/agents/` — per-role agent definitions (analyst, pm, architect, dev, qa, sm, po).
- `.bmad-core/templates/`, `.bmad-core/checklists/`, `.bmad-core/tasks/` — bundled artifacts.

### config_file
- `.bmad-core/core-config.yaml` — primary BMAD configuration. Documented in the BMAD installation guide.
- `.bmad-core/agents/<role>.md` — agent role definitions (e.g., `analyst.md`, `dev.md`, `qa.md`, `sm.md`).

### package_manifest
- npm `bmad-method` (https://www.npmjs.com/package/bmad-method).

### command_name
- `bmad-method` — CLI from the npm package. Invoked as `npx bmad-method install` / `npx bmad-method update`. Documented in the README.

### slash_command
- (BMAD operates via agent prompt files, not slash commands in the Claude Code surface; documented prompt files include `*analyst`, `*pm`, `*architect`, `*sm`, `*dev`, `*qa` invocations inside the chat — these are textual conventions, not Claude Code slash commands.)

### mcp_server_name
- (none documented)

### skill_name
- (none documented as Claude Code skills)

### plugin_manifest
- (none documented)

### cli_binary
- `bmad-method` — typically run via `npx`; not installed globally by default. Allowlist with low priority.

### cli_version_probe
- args: `["--version"]`
- expected output pattern: `bmad[^\d]*(\d+\.\d+)` → bucket `MAJOR.MINOR`.

## Negative-test markers (must NOT trigger this detector alone)
- Bare `(?i)\bBMAD\b` mention without `.bmad-core/` artifact — could be unrelated acronym. The existing `frameworks.json` `bmad` regex `(?i)\bBMAD\b` is high-recall and known-noisy; the new detector should require the `.bmad-core/` directory or `bmad-method` package install for High confidence.
- The Polish/Spanish word "bmad" — n/a (not a common word, but be conservative).

## Confidence wiring
- **High**: `.bmad-core/` directory present **and** any of (`.bmad-core/core-config.yaml`, `.bmad-core/agents/dev.md`, `bmad-method` in `package.json` devDependencies).
- **Medium**: `bmad-method` listed in a `package.json` without `.bmad-core/` corroboration, OR `npx bmad-method install` text in transcript.
- **Low**: bare `(?i)\bBMAD\b` or `(?i)bmad-method` mention.

## Source references (citations)
- https://github.com/bmad-code-org/bmad-method — official repository (README documents `.bmad-core/` layout, agent roles, install command).
- https://www.npmjs.com/package/bmad-method — npm package.
- https://github.com/bmad-code-org/bmad-method/blob/main/docs/installation.md — installation guide describing `.bmad-core/` contents.

## Note on ID reuse
- The existing `internal/analyzer/signatures/frameworks.json` already carries a `bmad` entry. Per spec C-004, the new `sdd_detectors.json` keeps the same canonical ID `bmad` (alias `BMAD-METHOD`).
