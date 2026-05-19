# Spec-Driven Develop (spec_driven_develop)

- Status: research_needed
- Category: spec-driven-development workflow toolkit
- Competitor priority: 9
- Official repository: (multiple repos use variants of this name; canonical one not confirmed as of cutoff)
- Official docs: (not confirmed)
- Release / package source: (not confirmed)
- Aliases: ["Spec-Driven Develop", "spec-driven-develop", "spec-driven-development"]

## Markers (public-source only)

### config_dir
- (not confirmed)

### config_file
- (not confirmed)

### package_manifest
- (not confirmed — candidate npm/PyPI names not verified)

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
- Generic phrase `spec-driven development` / `spec driven development` — this is the **category name** itself. It is mentioned in countless blog posts, docs of every tool in this registry, and academic papers. It MUST NOT match a tool detector on the phrase alone.
- Any mention inside an unrelated tool's docs that uses "spec-driven development" as a category descriptor.

## Confidence wiring
- **High**: cannot define until tool-specific markers are recorded.
- **Medium**: n/a.
- **Low**: n/a — refuse to emit low-confidence detection for this name because it overlaps with the category descriptor.

## Source references (citations)
- _None — this entry is `research_needed` pending disambiguation between the category descriptor and the specific tool named "Spec-Driven Develop"._

## Open questions / what's missing
1. The brief lists "Spec-Driven Develop" as a tool, but the phrase is also the umbrella category for the entire registry. The team needs to confirm the specific upstream repo intended.
2. Without disambiguation, any detector built on the name alone would be high-false-positive across every other SDD tool's docs.
3. **A-04 trigger**: This is a particularly strong candidate for a scope conversation — the name itself fights the detector design.
