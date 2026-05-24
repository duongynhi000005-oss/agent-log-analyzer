# Token-Saving Tooling Matrix

This matrix is the starting allowlist for paid remediation recommendations. The plugin may recommend these tools only with explicit user approval and the waiver gate. Installation is never automatic.

## Token Category Definitions

Recommendations must name the token category they target. "Saves tokens" is not specific enough.

| Category | Meaning | What to watch |
| --- | --- | --- |
| Input/context tokens | Prompt text, project instructions, tool schemas, file reads, tool results, and prior conversation context sent back to the model. | Retrieval and compression tools can reduce this, but extra instructions, indexing summaries, MCP schemas, or more turns can increase it. |
| Tool-output tokens | The subset of input/context tokens created by shell, MCP, file, and search outputs. | This is the most actionable waste surface for Claude Analyzer because it can often be reduced with scoped commands and compact output. |
| Output tokens | Visible assistant text and tool-call JSON emitted by the model. | Terseness rules can reduce this, but lower visible output is not proof of lower full-session cost. |
| Reasoning tokens | Hidden model reasoning budget reported by some harnesses, such as Codex. | Most plugin guidance does not directly reduce this; it usually changes when the task is simpler, has fewer retries, or requires fewer planning/tool loops. |
| Cached input tokens | Reused prompt/context tokens that may be billed or quota-weighted differently from fresh input. | A run can have lower visible output and still higher uncached-plus-output usage if cache behavior changes. |
| Telemetry only | Usage visibility, cost accounting, statuslines, and dashboards. | Useful for measurement and behavior change, but not a direct reducer unless the user acts on the signal. |

## Tier 1: Bundle Or Strongly Recommend

| Tool | Source | Primary category | Does not directly reduce | Product Decision |
| --- | --- | --- | --- | --- |
| Context Mode | https://github.com/mksglu/context-mode | Tool-output and input/context tokens | Output or reasoning tokens unless it also avoids turns | Recommend only when analysis shows tool-output bloat or context growth spikes, and label current benchmark as not proven until the intended command path is exercised. |
| ccusage | https://github.com/ryoppippi/ccusage | Telemetry only | Input, output, reasoning, or tool-output tokens | Always recommend as the independent metrics layer and comparison source for our analyzer. Do not call it a direct reducer. |
| claude-context | https://github.com/zilliztech/claude-context | Input/context tokens through semantic retrieval | Output or reasoning tokens | Recommend for repeated file reads or large-repo workflows only when retrieval replaces broad reads enough to amortize indexing/MCP overhead. Flag external API/vector DB requirements. |
| grepai | https://github.com/yoanbernabeu/grepai | Input/context and tool-output tokens through compact local retrieval | Output or reasoning tokens | Recommend for repeated file reads when the user wants local-first retrieval, with small limits and path filters. Requires embedding provider setup. |
| Claude Code Hooks Mastery | https://github.com/disler/claude-code-hooks-mastery | Implementation reference | Any token category by itself | Reference only. Use to guide our own hook design; do not ask users to install as a runtime dependency. |
| claude-token-efficient | https://github.com/drona23/claude-token-efficient | Output tokens and future assistant verbosity | Tool-output, cached input, or reasoning tokens | Recommend as a minimal diff, not overwrite. Useful only when output verbosity dominates enough to offset persistent CLAUDE.md input-token cost. |

## Tier 2: Supporting Recommendations

| Tool | Source | Primary category | Does not directly reduce | Product Decision |
| --- | --- | --- | --- | --- |
| ccstatusline | https://github.com/sirmalloc/ccstatusline | Telemetry only | Input, output, reasoning, or tool-output tokens | Recommend only if it does not conflict with Context Mode or existing statusline config. |
| Claude Code Usage Monitor | https://github.com/Maciek-roboblog/Claude-Code-Usage-Monitor | Telemetry only | Input, output, reasoning, or tool-output tokens | Optional external monitor for power users. Keep outside the plugin runtime. |
| Claude Code Usage Tracker | https://github.com/LyndonWangWork/Claude-Code-Usage-Tracker | Telemetry only | Input, output, reasoning, or tool-output tokens | Mention as an alternative monitor, not a default recommendation. |
| awesome-claude-code | https://github.com/hesreallyhim/awesome-claude-code | Ecosystem discovery only | Any token category by itself | Monitor continuously for candidate tools; never install from it directly. |

## Tier 3: Research Candidates

| Tool | Source | Primary category | Does not directly reduce | Product Decision |
| --- | --- | --- | --- | --- |
| RTK | https://github.com/rtk-ai/rtk | Tool-output and input/context tokens | Output or reasoning tokens | Promising for severe shell-output bloat. Treat as advanced because it rewrites command execution through hooks. Prefer explicit commands before global hooks. |
| Caveman | https://github.com/JuliusBrussee/caveman | Output tokens only | Tool-output, input/context, cached input, or reasoning tokens | Research only. Benchmarks showed visible output fell but full-session Claude/Codex token usage worsened. |
| memsearch | https://github.com/zilliztech/memsearch | Input/context tokens through memory retrieval | Output or reasoning tokens | Research only. Promising for sparse memory, but too stateful for the initial low-risk paid pack. |

## Recommendation Mapping

| Analyzer Signal | Recommend | Token category claim |
| --- | --- | --- |
| `tool_output_bloat` | Context Mode, RTK, quiet verification commands | Reduces tool-output/input-context only when actually used in the command path. |
| `repeated_file_reads` | grepai, claude-context, language-server/code-intelligence plugins | Reduces input/context only if retrieval replaces duplicate reads. |
| `retry_loop` | hooks architecture reference, statusline awareness, session hygiene skill | Indirectly reduces input/output/reasoning only if it prevents failed attempts or extra turns. |
| `context_growth_spikes` | Context Mode, ccstatusline, session hygiene, minimal CLAUDE.md review | Targets input/context growth; ccstatusline is telemetry only. |
| `verbose_assistant_output` | claude-token-efficient or Caveman-style terse-response guidance | Targets output tokens only; not proof of lower full-session cost. |
| Any paid scan | ccusage, awesome-claude-code monitoring note | ccusage is telemetry only; awesome-claude-code is discovery only. |

## Cost Translation Rule

Dollar claims must be calculated from the published rate card for the exact model and token categories. Do not infer cost from a single token number.

- Claude Sonnet 4.6 API estimates use separate input, cache-write, cache-read, and output rates.
- Codex API estimates use input, cached-input, output, and separately reported reasoning-output tokens where available.
- Native Claude Code or Codex product billing can differ from direct API inference rates; report it as a separate surface.
- A tool may reduce tool-output tokens but increase published API cost if cache reads, cache writes, output, or reasoning grow.

## Guardrails

- Prefer plugin marketplace or package-manager installs over curl scripts.
- Curl install scripts are allowed only as reviewed fallback instructions, never silent defaults.
- External hosted retrieval tools must disclose API keys, data movement, and vendor dependency.
- Local semantic search tools must disclose indexing cost, storage location, and embedding provider.
- Tools that rewrite shell commands are advanced-only and must be separately approved.
- Generated artifacts must not include unknown private tool names from logs.
