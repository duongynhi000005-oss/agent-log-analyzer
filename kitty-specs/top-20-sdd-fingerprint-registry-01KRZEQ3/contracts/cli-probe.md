# Contract: `sdd.CLIProbe`

```go
package sdd

import "context"

type CLIProbe interface {
    LookPath(name string) (found bool)
    Version(ctx context.Context, name string, args []string) (rawOutput string, ok bool)
}
```

## Behavioral contract

### `LookPath(name)`

- MUST return only a boolean. The resolved executable path MUST NOT be returned, logged, returned via error, or stored on the receiver.
- MUST be safe to call on any allowlisted binary name; not safe to call on attacker-controlled input. The evaluator guards this by only invoking `LookPath` for names that appear as a `cli_binary` marker in a loaded `SDDDetector`.
- MUST NOT read project data, modify files, contact networks, or shell out.

### `Version(ctx, name, args)`

- MUST execute the binary identified by `name` with `args`, no shell, no stdin, with a sanitized empty `Env`.
- MUST honor `ctx` deadlines; a 2-second deadline is the upper bound (NFR-002). Timeout returns `("", false)`.
- MAY return a non-empty `rawOutput` only to in-package callers (the version-bucket normalizer). The `rawOutput` is consumed within the `sdd` package, normalized, and the **raw string is then discarded**. It MUST NOT escape into the report or the aggregate event.
- Returns `ok = true` iff the process exited with status 0 within the deadline.

## Allowed `args`

For every detector that declares a `cli_version_probe` marker, `args` MUST come from the marker's `version_args` field (default `["--version"]`). Forbidden patterns include any value resembling a flag that:

- references file paths (`--config <path>`),
- references network addresses (`--registry <url>`),
- requests auth tokens,
- writes files,
- spawns servers,
- enables verbose / debug output that could leak environmental info.

The registry loader rejects detector entries whose `version_args` contain any of: `"--config"`, `"--registry"`, `"--token"`, `"--server"`, `"--login"`, or any value containing `/`. (Implementation: a static deny-list in the loader; rejection panics at startup.)

## Implementations

### `sdd.RealProbe`

Production implementation. Wraps `exec.LookPath` and `exec.CommandContext`. Implementation MUST satisfy the rules above. Unit-tested via integration tests that probe a known-safe binary (`go` itself) and assert the path is not leaked.

### `sdd.FakeProbe`

Test fake. Constructed with a map:

```go
type FakeProbe struct {
    Installed map[string]bool
    Versions  map[string]string // rawOutput per binary; "" means probe fails
}
```

Every unit test in the `sdd` package uses `FakeProbe`. The real `os/exec` codepath is exercised only by the `RealProbe` integration test.

## Invariants asserted in tests

1. `RealProbe.LookPath` returning `true` for a binary on `PATH` does not put the resolved path into any byte buffer reachable outside the function.
2. `RealProbe.Version` returning successfully does not surface raw output to any caller; the caller only ever sees the bucket.
3. `Version` honors a 2-second deadline (NFR-002).
4. The evaluator never calls `LookPath` for a name absent from the registry's `cli_binary` allowlist.
5. Disabling all probes (using `FakeProbe{}`) yields identical fingerprint output minus the `installed` / `version_bucket` fields.
