# Spec Workflow MCP (spec_workflow_mcp)

- Status: verified
- Category: MCP server providing spec-driven workflow tools (requirements → design → tasks)
- Competitor priority: 7
- Official repository: https://github.com/Pimzino/spec-workflow-mcp
- Official docs: https://github.com/Pimzino/spec-workflow-mcp#readme
- Release / package source: https://www.npmjs.com/package/@pimzino/spec-workflow-mcp (installed via `npx @pimzino/spec-workflow-mcp`)
- Aliases: ["Spec Workflow MCP", "spec-workflow-mcp", "@pimzino/spec-workflow-mcp"]

## Markers (public-source only)

### config_dir
- `.spec-workflow/` — workspace directory created by the server (contains specs, steering, and the dashboard state). Documented in https://github.com/Pimzino/spec-workflow-mcp#features.
- `.spec-workflow/specs/<feature>/` — per-feature workspace.
- `.spec-workflow/steering/` — steering rules.

### config_file
- `.spec-workflow/specs/<feature>/requirements.md`
- `.spec-workflow/specs/<feature>/design.md`
- `.spec-workflow/specs/<feature>/tasks.md`

The `requirements → design → tasks` triple inside `.spec-workflow/` is documented at https://github.com/Pimzino/spec-workflow-mcp#documents.

### package_manifest
- npm `@pimzino/spec-workflow-mcp` (https://www.npmjs.com/package/@pimzino/spec-workflow-mcp).

### command_name
- (server is invoked by the MCP client, not as a standalone CLI on $PATH)

### slash_command
- (none — Spec Workflow MCP is invoked through MCP tool calls, not slash commands)

### mcp_server_name
- `spec-workflow` — the MCP server's documented namespace. Tool calls take the form `mcp__spec-workflow__<tool>` (or `mcp__spec_workflow__<tool>` depending on the client's normalization). This is **the primary marker** for this tool.
- Tool names documented in the README include: `create-spec-doc`, `list-specs`, `get-spec-context`, `manage-tasks`, `request-approval`, `get-approval-status`, `create-steering-doc`, `get-steering-context`.

### skill_name
- (none)

### plugin_manifest
- (none)

### cli_binary
- (none on $PATH; `npx @pimzino/spec-workflow-mcp` is the documented runner)

### cli_version_probe
- (not applicable)

## Negative-test markers (must NOT trigger this detector alone)
- Generic `requirements.md`, `design.md`, `tasks.md` outside `.spec-workflow/specs/<feature>/` — these names overlap with Kiro and other SDD tools. Require the `.spec-workflow/` parent or an `mcp__spec-workflow__*` tool call to disambiguate.
- Bare phrase "spec workflow" — too generic.

## Confidence wiring
- **High**: any `mcp__spec[_-]workflow__*` tool call in the transcript **or** `.spec-workflow/` directory present + any one of (`@pimzino/spec-workflow-mcp` in `package.json`, `.spec-workflow/specs/<feature>/requirements.md`).
- **Medium**: `@pimzino/spec-workflow-mcp` listed in a `package.json` without `.spec-workflow/` or MCP call corroboration.
- **Low**: bare mention of `spec-workflow-mcp`.

## Source references (citations)
- https://github.com/Pimzino/spec-workflow-mcp — official repository (README documents the `.spec-workflow/` workspace, MCP tool names, and installation).
- https://www.npmjs.com/package/@pimzino/spec-workflow-mcp — npm package source.
