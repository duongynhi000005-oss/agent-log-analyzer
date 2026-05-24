# Published API Cost Translation

Research date: 2026-05-24

This note converts repeated benchmark token deltas into dollar estimates using published API inference rates. It does not replace harness-native cost fields. Claude Code and Codex can apply product-specific routing, credits, subscriptions, or internal model work that differs from direct API inference.

## Published Rates Used

Claude Code was run with `--model sonnet`. The exposed aggregate usage fields are repriced with Claude Sonnet 4.6 API rates:

- input: `$3.00` / MTok
- 1-hour cache write: `$6.00` / MTok
- cache read: `$0.30` / MTok
- output: `$15.00` / MTok

Source: https://platform.claude.com/docs/en/about-claude/pricing

The raw Claude stdout also reports internal model usage such as Haiku in some runs, but the comparison JSON exposes only aggregate usage buckets. Published API estimates therefore use Sonnet rates as an approximation and are reported separately from Claude Code native `total_cost_usd`.

Codex JSON usage is repriced under the published Standard rate for `gpt-5.3-codex`:

- input: `$1.75` / MTok
- cached input: `$0.175` / MTok
- output: `$14.00` / MTok

Sources:

- https://developers.openai.com/api/docs/pricing
- https://help.openai.com/en/articles/4936856-what-are-tokens-and-how-to-count-them

Because Codex exposes `reasoning_output_tokens` separately, this estimate bills reasoning tokens at the output-token rate. If a future exporter includes reasoning inside `output_tokens`, do not add it twice.

## Formula

Claude Sonnet estimate:

```text
cost =
  input_tokens * 3.00 / 1,000,000
+ cache_creation_input_tokens * 6.00 / 1,000,000
+ cache_read_input_tokens * 0.30 / 1,000,000
+ output_tokens * 15.00 / 1,000,000
```

Codex estimate:

```text
uncached_input_tokens = input_tokens - cached_input_tokens

cost =
  uncached_input_tokens * 1.75 / 1,000,000
+ cached_input_tokens * 0.175 / 1,000,000
+ output_tokens * 14.00 / 1,000,000
+ reasoning_output_tokens * 14.00 / 1,000,000
```

## Final 3x Mean Cost Results

All rows are means across three fresh baseline/optimized pairs.

| Suite | Harness/rate assumption | Baseline mean | Optimized mean | Published API delta mean | Native cost signal |
| --- | --- | ---: | ---: | ---: | --- |
| Agent Analyzer guided | Claude Sonnet 4.6 | `$0.246837` | `$0.187629` | `-$0.059207` | Claude Code native `-$0.044219` |
| claude-context limit 3 | Claude Sonnet 4.6 | `$0.222950` | `$0.280989` | `+$0.058038` | Claude Code native `+$0.048434` |
| context-mode batch | Claude Sonnet 4.6 | `$0.255157` | `$0.202981` | `-$0.052175` | Claude Code native `-$0.036390` |
| grepai path-constrained | Claude Sonnet 4.6 | `$0.259462` | `$0.221805` | `-$0.037657` | Claude Code native `-$0.017598` |
| claude-token-efficient | Claude Sonnet 4.6 | `$0.230198` | `$0.225991` | `-$0.004208` | Claude Code native `-$0.003828` |
| RTK explicit | Claude Sonnet 4.6 | `$0.244114` | `$0.199798` | `-$0.044316` | Claude Code native `-$0.031479` |
| Probe | Claude Sonnet 4.6 | `$0.230660` | `$0.268999` | `+$0.038340` | Claude Code native `+$0.038069` |
| Semble | Claude Sonnet 4.6 | `$0.275398` | `$0.161205` | `-$0.114194` | Claude Code native `-$0.089147` |
| Squeez | Claude Sonnet 4.6 | `$0.232566` | `$0.204342` | `-$0.028224` | Claude Code native `-$0.014049` |
| Agent Analyzer text guidance | GPT-5.3-Codex Standard | `$0.196060` | `$0.133667` | `-$0.062392` | uncached+output `-24,369`; reasoning `-45` |
| Caveman Claude | Claude Sonnet 4.6 | `$0.236776` | `$0.245987` | `+$0.009211` | Claude Code native `+$0.009919` |
| Caveman Codex | GPT-5.3-Codex Standard | `$0.186165` | `$0.152179` | `-$0.033986` | uncached+output `-4,739`; reasoning `-2` |

## Product Rule

Cost claims must name both the token category and the pricing surface:

- "Reduced tool-output tokens" is an Analyzer workflow claim.
- "Reduced visible output cost" is an output-token pricing claim.
- "Reduced published API cost" requires input/cache/output/reasoning repricing.
- "Reduced Claude Code/Codex cost" requires native harness billing/credit evidence.

Do not collapse those into one generic savings claim.
