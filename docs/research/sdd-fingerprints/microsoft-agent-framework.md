# Microsoft Agent Framework (microsoft_agent_framework)

- Status: verified
- Category: agent-orchestration framework from Microsoft (formerly Semantic Kernel's agent surface, evolving into the "Agent Framework")
- Competitor priority: 17
- Official repository: https://github.com/microsoft/agent-framework
- Official docs: https://learn.microsoft.com/agent-framework/
- Release / package source:
  - .NET / NuGet: `Microsoft.AgentFramework.*` packages on https://www.nuget.org/packages?q=Microsoft.AgentFramework
  - Python / PyPI: `agent-framework` (https://pypi.org/project/agent-framework/) and related `agent-framework-core`
- Aliases: ["Microsoft Agent Framework", "agent-framework", "Microsoft.AgentFramework"]

## Markers (public-source only)

### config_dir
- (no canonical `.agent-framework/` workspace directory documented; the framework is consumed as a library, not as an init-driven scaffold)

### config_file
- `.csproj` containing a `<PackageReference Include="Microsoft.AgentFramework.*" />` — strong NuGet manifest signal. Documented at https://learn.microsoft.com/agent-framework/quickstart.
- `pyproject.toml` or `requirements.txt` containing `agent-framework` — Python signal.

### package_manifest
- NuGet: `Microsoft.AgentFramework`, `Microsoft.AgentFramework.OpenAI`, `Microsoft.AgentFramework.AzureAI`, `Microsoft.AgentFramework.A2A`, `Microsoft.AgentFramework.Workflows`. Documented at https://www.nuget.org/packages?q=Microsoft.AgentFramework.
- PyPI: `agent-framework`, `agent-framework-core`, `agent-framework-azure-ai`. Documented at https://pypi.org/project/agent-framework/.

### command_name
- (no documented CLI binary on $PATH)

### slash_command
- (none in the Claude Code surface)

### mcp_server_name
- (Microsoft Agent Framework supports being an MCP **client**; it does not ship a canonical `mcp__microsoft-agent-framework__*` server)

### skill_name
- (none documented as Claude Code skills)

### plugin_manifest
- (none documented)

### cli_binary
- (none reliably on $PATH; do NOT allowlist)

### cli_version_probe
- (not applicable)

## Negative-test markers (must NOT trigger this detector alone)
- The phrase "agent framework" — extremely generic; matches dozens of unrelated agent frameworks (LangChain, AutoGen, CrewAI, etc.). MUST NOT match on the bare phrase.
- "Microsoft" alone — obvious.
- The unrelated `agent-framework` npm package squatters / placeholders.
- `Microsoft.SemanticKernel` package references — Semantic Kernel is the **predecessor**, not the same product; record it as Medium signal only if the team agrees Semantic Kernel rolls into this detector.

## Confidence wiring
- **High**: a `<PackageReference Include="Microsoft.AgentFramework.*" />` in a `.csproj` **or** an `agent-framework` (PyPI) dependency in `pyproject.toml`/`requirements.txt` **plus** any second corroborating mention.
- **Medium**: exactly one of the above package-manifest references without a second signal.
- **Low**: free-text mention of `(?i)\bMicrosoft Agent Framework\b`.

## Source references (citations)
- https://github.com/microsoft/agent-framework — official repository.
- https://learn.microsoft.com/agent-framework/ — official Microsoft Learn docs.
- https://www.nuget.org/packages?q=Microsoft.AgentFramework — NuGet package listings.
- https://pypi.org/project/agent-framework/ — Python package source.

## Open questions / what's missing
- Confirm whether `Microsoft.SemanticKernel.*` references should also feed this detector or be a separate entry. Out of scope for WP04; flag for WP05/WP09 review.
