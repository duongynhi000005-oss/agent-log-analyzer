# Signature Research Tooling

The analyzer ships with a committed public-tool signature registry. The research crawler is a separate review tool for finding candidate additions from public sources.

Run:

```sh
./scripts/research-signatures.sh
```

Output is written to `.data/signature-candidates.json`, which is intentionally gitignored. Review candidates manually before editing `internal/analyzer/signatures/*.json`.

Set `GITHUB_TOKEN` before running if GitHub search rate limits anonymous requests. MCP Registry sources support `max_pages` cursor pagination in `docs/signature-research-sources.json`.

## Privacy Boundary

- Only public ecosystem sources are crawled.
- User uploads, reports, and aggregate events are never input to this tool.
- Unknown private MCP, skill, and plugin names remain count-only in production analytics.
- Candidate output contains public package/repository/registry names and source URLs, not user-log excerpts.

## Source Types

- `mcp_registry`: official MCP Registry server metadata.
- `npm_search`: npm package search results for Claude/MCP terms.
- `github_search_repositories`: GitHub repository search results for Claude Code, MCP, plugins, skills, and workflow frameworks.
- `text`: generic public text extraction for MCP tool names, `claude mcp add` commands, and `.claude/skills` paths.

## Review Workflow

1. Run the crawler and sort candidates by category and confidence.
2. Promote only public, stable, recognizable tooling names.
3. Prefer deterministic patterns from official docs, package names, registry IDs, or repository names.
4. Add or update analyzer tests for promoted signatures.
5. Never promote one-off private names found only in user uploads.

## Useful Primary Sources

- Anthropic Claude Code plugin and skill docs describe plugins as bundles of skills, commands, agents, hooks, MCP servers, LSP servers, and monitors.
- Claude Code skill docs define `.claude/skills/<name>/SKILL.md`, `.claude/commands/<name>.md`, and slash invocation names.
- The official MCP Registry exposes `GET /v0/servers` and provides canonical public MCP server metadata.
