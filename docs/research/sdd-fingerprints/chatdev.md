# ChatDev (chatdev)

- Status: verified
- Category: multi-agent software development framework
- Competitor priority: 11
- Official repository: https://github.com/OpenBMB/ChatDev
- Official docs: https://github.com/OpenBMB/ChatDev#readme
- Release / package source: source clone of https://github.com/OpenBMB/ChatDev (Python; no published PyPI package for the main framework — install is `git clone` + `pip install -r requirements.txt`)
- Aliases: ["ChatDev", "chatdev", "OpenBMB ChatDev"]

## Markers (public-source only)

### config_dir
- `WareHouse/` — output directory for ChatDev "software" projects. The `run.py` driver writes each company-name/project-name run under `WareHouse/<name>_<org>_<timestamp>/`. Documented in https://github.com/OpenBMB/ChatDev#-quickstart.
- `CompanyConfig/` — per-company config directory bundled in the repo (e.g., `CompanyConfig/Default/`, `CompanyConfig/Art/`).

### config_file
- `run.py` at repo root invoking ChatDev's `chatdev/` Python package — distinctive entry point. Documented in the README under "Quickstart" (`python3 run.py --task "..." --name "..."`).
- `ChatChainConfig.json`, `PhaseConfig.json`, `RoleConfig.json` — three JSON config files inside `CompanyConfig/<company>/`. Documented in https://github.com/OpenBMB/ChatDev/wiki/Customization.

### package_manifest
- (no published PyPI package for the main framework; ChatDev is typically run from a git clone. Do not rely on a package-manifest marker.)

### command_name
- (none — ChatDev is invoked as `python3 run.py ...`, not via a packaged CLI binary on $PATH)

### slash_command
- (none — ChatDev does not register Claude Code slash commands)

### mcp_server_name
- (none)

### skill_name
- (none)

### plugin_manifest
- (none)

### cli_binary
- (none on $PATH; do NOT allowlist a generic `chatdev` binary)

### cli_version_probe
- (not applicable)

## Negative-test markers (must NOT trigger this detector alone)
- The English phrase "chat dev" or "chat development" — too generic.
- The unrelated GitHub repo names containing "chatdev" but not under the `OpenBMB` org — require the trio `CompanyConfig/<company>/{ChatChainConfig,PhaseConfig,RoleConfig}.json` for High confidence.
- A `run.py` file alone — Python projects commonly have a `run.py`; require ChatDev-specific config filenames.

## Confidence wiring
- **High**: any two of (`WareHouse/`, `CompanyConfig/<company>/ChatChainConfig.json`, `CompanyConfig/<company>/PhaseConfig.json`, `CompanyConfig/<company>/RoleConfig.json`) present in the same repo.
- **Medium**: one of the JSON config filenames listed above present.
- **Low**: free-text mention of `(?i)\bChatDev\b` together with `OpenBMB`.

## Source references (citations)
- https://github.com/OpenBMB/ChatDev — official repository (README documents `run.py`, `WareHouse/`, `CompanyConfig/` directory structure, and the JSON config files).
- https://github.com/OpenBMB/ChatDev/wiki/Customization — wiki page documenting `ChatChainConfig.json`, `PhaseConfig.json`, `RoleConfig.json`.
- https://github.com/OpenBMB/ChatDev/blob/main/run.py — the driver script that is the canonical entry point.
