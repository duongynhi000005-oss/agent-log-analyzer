# Ecosystem Signatures

## See also

The [SDD fingerprint registry](sdd-fingerprint-registry.md) is the canonical
source for spec-driven-development tool detection (Spec Kitty, GitHub Spec
Kit, OpenSpec, Kiro, BMAD-METHOD, and the long tail). It defines a typed
`SDDDetector` schema with deterministic confidence rules, a privacy-safe CLI
probe contract, and a bounded `EcosystemFingerprint` report record. New SDD
detection lives in `Ecosystem.WorkflowFingerprints` (a `[]EcosystemFingerprint`
slice).

The legacy `WorkflowFrameworks []string` field documented below is preserved
unchanged per constraint
[C-004](../kitty-specs/top-20-sdd-fingerprint-registry-01KRZEQ3/spec.md#constraints)
so existing report consumers (web UI, golden tests) keep working. This file
continues to document the broader ecosystem-signature pattern that drives
that legacy field.

## Registries

Ecosystem intelligence is driven by versioned JSON registries embedded into the analyzer:

```text
internal/analyzer/signatures/
  coding_agents.json
  frameworks.json
  mcp_servers.json
  package_managers.json
  plugins.json
  skills.json
```

Rules:

- Known public tools emit stable IDs.
- Unknown MCP/skill/plugin names are counted, not stored.
- Patterns must not capture or emit private raw strings.
- Every new signature needs a fixture or detector test when practical.

The current registry is intentionally public-tool-only. Private company tool names remain outside aggregate analytics unless an explicit opt-in path is added later.

Candidate discovery tooling lives in `cmd/signature-research` and `scripts/research-signatures.sh`. It crawls public registries/search APIs and writes review-only output to `.data/signature-candidates.json`; it does not update the production registry automatically.
