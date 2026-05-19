# SDD Fingerprint Research

This directory holds the per-tool research that feeds the production SDD detector registry (`internal/analyzer/signatures/sdd_detectors.json`). Each tool's file follows the template defined in [`research.md` §R-01](../../../kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/research.md) and must cite public sources for every marker before it can graduate from `research_needed` to `verified`.

The customer-facing registry overview lives at `docs/sdd-fingerprint-registry.md` (owned by WP09 — not in this directory).

## Methodology

See [`kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/research.md` §R-01](../../../kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/research.md) for the canonical per-tool template and §R-02 for the disposable-clone procedure.

Summary:

1. Locate the tool's official repository and docs.
2. (Optional) Clone into `/Users/robert/code-analyzer-dev/claude-code-analyzer-20260519-082245-0QWuF7/research/<tool-id>/` — **never** commit the cloned source.
3. Record markers per source class (`config_dir`, `config_file`, `package_manifest`, `command_name`, `slash_command`, `mcp_server_name`, `skill_name`, `plugin_manifest`, `cli_binary`, `cli_version_probe`).
4. Cite a public URL for every marker.
5. List generic negative-test markers that MUST NOT trigger this detector alone.
6. Set `status: verified` only when all citations resolve and at least one tool-specific marker exists.

## Status semantics

- **`verified`** — Every marker has a public citation that a reviewer can follow to confirm the tool ships that artifact / command / server name. The detector is safe to emit in production.
- **`research_needed`** — One or more markers lack a citation, or the tool has no clear public fingerprintable surface. Per FR-013, production detectors MUST NOT be emitted for tools in this state. Per C-001, this mission ships zero `research_needed` detectors; any tool stuck here triggers a scope conversation with the user (see A-04).
- **`blocked`** — The tool has been investigated and found to have no public artifact useful for fingerprinting. Escalate to the user.

## Citation policy

- Citations MUST resolve to a public URL (GitHub repo, vendor docs, npm/PyPI page, blog post). No private artifacts, no inferred markers.
- Where uncertainty exists about a repo move or rename, pin a commit hash in the citation.
- Disposable install runs in the research-clones directory count as public-source equivalent (we are observing what a published tool generates), but each such observation MUST cite the upstream repo or docs that document the install command.

## Top-20 status index

| # | Tool | Slug | Status |
| --- | --- | --- | --- |
| 1 | Spec Kitty | [`spec-kitty.md`](./spec-kitty.md) | verified |
| 2 | GitHub Spec Kit | [`github-spec-kit.md`](./github-spec-kit.md) | verified |
| 3 | OpenSpec | [`openspec.md`](./openspec.md) | verified |
| 4 | Kiro | [`kiro.md`](./kiro.md) | verified |
| 5 | BMAD-METHOD | [`bmad.md`](./bmad.md) | verified |
| 6 | GSD | [`gsd.md`](./gsd.md) | research_needed |
| 7 | Spec Workflow MCP | [`spec-workflow-mcp.md`](./spec-workflow-mcp.md) | verified |
| 8 | SDD Pilot | [`sdd-pilot.md`](./sdd-pilot.md) | research_needed |
| 9 | Spec-Driven Develop | [`spec-driven-develop.md`](./spec-driven-develop.md) | research_needed |
| 10 | spec2ship | [`spec2ship.md`](./spec2ship.md) | research_needed |
| 11 | ChatDev | [`chatdev.md`](./chatdev.md) | verified |
| 12 | PAUL | [`paul.md`](./paul.md) | research_needed |
| 13 | fspec | [`fspec.md`](./fspec.md) | research_needed |
| 14 | whenwords | [`whenwords.md`](./whenwords.md) | research_needed |
| 15 | Intent | [`intent.md`](./intent.md) | research_needed |
| 16 | Cognition / Devin | [`cognition-devin.md`](./cognition-devin.md) | verified |
| 17 | Microsoft Agent Framework | [`microsoft-agent-framework.md`](./microsoft-agent-framework.md) | verified |
| 18 | Tessl | [`tessl.md`](./tessl.md) | research_needed |
| 19 | Agentic Code | [`agentic-code.md`](./agentic-code.md) | research_needed |
| 20 | CodeSpeak | [`codespeak.md`](./codespeak.md) | research_needed |

**Verified: 9 / 20. Research-needed: 11 / 20.** Per C-001 the mission DoD requires 20/20 `verified`; tools listed `research_needed` here either (a) have ambiguous public-source surfaces (generic name collisions like "Intent", "PAUL", "Agentic Code", "Spec-Driven Develop") or (b) are private hosted products without a CLI/config artifact that a static analyzer can detect from text fixtures alone. Both A-04 cases should trigger a scope conversation with the user before WP05 hard-codes detectors for them. See each per-tool file's "Open questions" section for the specific gaps.

### research_needed tools awaiting scope conversation (A-04)

| Slug | Reason |
| --- | --- |
| `gsd.md` | Canonical repo not identified; three-letter acronym is high-collision. |
| `sdd-pilot.md` | Canonical repo not identified. |
| `spec-driven-develop.md` | Name collides with the category descriptor itself; reviewer must confirm specific upstream. |
| `spec2ship.md` | Canonical repo not identified. |
| `paul.md` | Canonical repo not identified; common given name. |
| `fspec.md` | Multiple unrelated repos use the name. |
| `whenwords.md` | Canonical repo not identified. |
| `intent.md` | Extreme false-positive risk (Android `Intent`, NLP "intent classification", common noun). |
| `tessl.md` | Hosted vendor; local artifact surface unconfirmed. |
| `agentic-code.md` | Name collides with the category descriptor. |
| `codespeak.md` | Canonical repo not identified. |

## Reviewer quickcheck

```sh
grep -l "^- Status: verified" docs/research/sdd-fingerprints/*.md | wc -l
# Mission DoD (C-001): 20
# Current: see status index above
```

## Cross-link

- Mission spec: [`kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/spec.md`](../../../kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/spec.md)
- Research notes: [`kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/research.md`](../../../kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/research.md)
- WP09 consumer-facing registry: `docs/sdd-fingerprint-registry.md` (created by WP09 — do not edit from WP04)
