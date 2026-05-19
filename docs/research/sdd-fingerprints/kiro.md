# Kiro (kiro)

- Status: verified
- Category: AWS agentic IDE with spec-driven workflow
- Competitor priority: 4
- Official repository: (no fully open-source repo — Kiro is an Amazon product; the IDE downloads at https://kiro.dev/)
- Official docs: https://kiro.dev/docs
- Release / package source: https://kiro.dev/downloads (Mac / Windows / Linux installers)
- Aliases: ["Kiro", "kiro", "Kiro IDE", "AWS Kiro"]

## Markers (public-source only)

### config_dir
- `.kiro/` — per-project Kiro workspace directory. Documented at https://kiro.dev/docs/specs/concepts (the IDE writes a `.kiro/` folder per project containing specs and steering rules).
- `.kiro/specs/` — directory of spec files.
- `.kiro/steering/` — directory of steering rules (`*.md`).

### config_file
- `.kiro/specs/<feature>/requirements.md`
- `.kiro/specs/<feature>/design.md`
- `.kiro/specs/<feature>/tasks.md`

The "requirements → design → tasks" triplet inside `.kiro/specs/<feature>/` is Kiro's canonical spec layout. Documented at https://kiro.dev/docs/specs/best-practices.

### package_manifest
- (none — Kiro is distributed as a desktop IDE installer, not a published npm/PyPI package)

### command_name
- (the IDE binary is `Kiro` on macOS; not commonly on $PATH. No documented `kiro` CLI on $PATH as of cutoff)

### slash_command
- (none documented — Kiro uses chat panel, not slash commands shared with Claude Code)

### mcp_server_name
- (none documented)

### skill_name
- (none — Kiro is its own IDE, not a Claude Code extension)

### plugin_manifest
- (none documented)

### cli_binary
- (none reliably on $PATH; do NOT allowlist a `kiro` binary)

### cli_version_probe
- (not applicable)

## Negative-test markers (must NOT trigger this detector alone)
- Generic `requirements.md`, `design.md`, `tasks.md` outside a `.kiro/specs/<feature>/` directory — these names are extremely generic. Require the `.kiro/` parent directory.
- The Japanese name "Kiro" — many unrelated mentions; only count when paired with `.kiro/` artifact.
- "kiro.dev" mention in a URL is medium signal but not high without artifact.

## Confidence wiring
- **High**: `.kiro/` directory present **and** any of (`.kiro/specs/<feature>/requirements.md`, `.kiro/steering/<rule>.md`).
- **Medium**: free-text mention of `(?i)\bKiro IDE\b` or `kiro.dev` URL reference with no artifact.
- **Low**: bare `(?i)\bKiro\b` mention (very generic word; common name).

## Source references (citations)
- https://kiro.dev/ — official product home page.
- https://kiro.dev/docs/specs/concepts — describes the spec workflow and `.kiro/` workspace concept.
- https://kiro.dev/docs/specs/best-practices — describes the `requirements.md` → `design.md` → `tasks.md` triplet.
- https://kiro.dev/docs/steering — describes `.kiro/steering/` rule files.

## Open questions
- No public CLI binary name confirmed on $PATH; if a future Kiro release ships one, add a `cli_binary` and `cli_version_probe` entry.
