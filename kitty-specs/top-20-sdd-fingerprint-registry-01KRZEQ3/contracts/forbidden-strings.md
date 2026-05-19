# Contract: Forbidden raw-string categories

This document is the canonical list used by the `sdd.leak_test` serialization-leak test (NFR-001). The test builds a fully populated `Report` carrying one canary value per category, serializes the report and its `AggregateEvent`, and asserts every canary is absent from both serialized outputs.

| # | Category | Canary value used in test |
| --- | --- | --- |
| 1 | user prompts | `FORBIDDEN-PROMPT-CANARY` |
| 2 | task descriptions | `FORBIDDEN-TASK-CANARY` |
| 3 | raw transcript excerpts | `FORBIDDEN-TRANSCRIPT-CANARY` |
| 4 | raw tool inputs | `FORBIDDEN-TOOLIN-CANARY` |
| 5 | raw tool outputs | `FORBIDDEN-TOOLOUT-CANARY` |
| 6 | raw file paths | `/Users/robert/private-project` |
| 7 | repo URLs | `git@github.com:private-org/private-repo.git` |
| 8 | branch names | `customer-private-branch` |
| 9 | usernames | `jdoe-private` |
| 10 | hostnames | `corp-internal-host.local` |
| 11 | emails | `private.user@example.internal` |
| 12 | session IDs | `sess_PRIVATESESSION` |
| 13 | transcript paths | `/Users/robert/.claude/projects/PRIVATE/transcript.jsonl` |
| 14 | private MCP / skill / plugin names | `mcp__private__internal`, `skill__private__internal`, `plugin__private__internal` |
| 15 | raw `LookPath` / `which` paths | `/opt/homebrew/bin/openspec`, `/usr/local/bin/spec-kitty` |
| 16 | raw `--version` output / stable hashes of private strings | `openspec 1.2.3 built /private/path`, `sha256-of-private-name-DO-NOT-LEAK` |

## Test harness sketch

```go
func TestReportSerializationContainsNoForbiddenStrings(t *testing.T) {
    canaries := []string{ /* all values from the table */ }

    rep := buildFullyPopulatedReportWithCanaries(t, canaries)

    for _, target := range []any{rep, rep.AggregateEvent} {
        out, err := json.Marshal(target)
        if err != nil {
            t.Fatal(err)
        }
        s := string(out)
        for _, c := range canaries {
            if strings.Contains(s, c) {
                t.Fatalf("forbidden canary %q leaked into %T JSON", c, target)
            }
        }
    }
}
```

## Failure mode

A failing assertion blocks the build. The category should be named in the failure message so reviewers can localize the leak quickly.

## Out-of-test policy

This canary list is the **floor**, not the ceiling. The schema contracts (`sdd-detector.schema.json`, `ecosystem-fingerprint.schema.json`) define the **structural** invariants that make the leak test redundant in the happy path; the leak test exists to catch schema drift and accidental new string fields.
