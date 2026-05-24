# Claude Log Analyzer

Deterministic performance profiler for AI coding workflows. The product path is analysis first, download second: users measure their own Claude Code waste locally, then install a generated optimization pack only when the evidence supports it.

This repo starts with a Docker-local, end-to-end implementation:

- run the analyzer locally against one Claude Code JSONL log
- write a sanitized report JSON that the user can inspect before upload
- upload only the sanitized report JSON
- detect waste patterns and ecosystem fingerprints
- generate an ephemeral report JSON
- generate a waiver-gated optimization plugin from the paid aggregate report
- view the report in a static local web UI

The production target is CDN + local deterministic CLI + report-only upload + short-lived report storage. Local development intentionally avoids cloud dependencies so the complete flow can be tested before any infrastructure is provisioned.

There is intentionally no browser upload form. Claude Code logs live under `~/.claude`, which is awkward for Finder/browser upload flows. The public launch path is local-first:

1. The user installs the source-available CLI with `go install github.com/robertdouglass/claude-log-analyzer/cmd/claude-analyzer@v0.1.0`.
2. `claude-analyzer analyze --out ./claude-analyzer-report.json` finds the latest Claude Code JSONL log, parses and redacts it locally, and writes a sanitized report.
3. The user reviews the JSON with `jq . ./claude-analyzer-report.json`.
4. `claude-analyzer upload ./claude-analyzer-report.json` sends only the sanitized report to `POST /api/client-reports`.
5. The short-lived report is opened at `/r/{job_id}/{report_token}` and expires on the retention schedule.

Legacy raw-log token upload endpoints still exist for internal Docker smoke coverage while the paid scan is moved to the same local-first model. They are not the public onboarding path.

Paid delivery contract: [docs/remediation/plugin-artifacts.md](docs/remediation/plugin-artifacts.md).
Repeated benchmark suite: [docs/benchmarks/repeated-benchmark-suite.md](docs/benchmarks/repeated-benchmark-suite.md).
Primary sanitized benchmark recordings: [docs/benchmarks/primary-data/](docs/benchmarks/primary-data/).
Cost translation and scale-up math: [docs/benchmarks/api-cost-translation.md](docs/benchmarks/api-cost-translation.md).

Current product proof: three fresh Claude Code `-p` pairs on the noisy benchmark showed Agent Analyzer guidance reducing estimated tokens by `12,370`, tool-output tokens by `12,698`, visible Claude output by `504`, and published Claude Sonnet 4.6 API-rate cost by `23.986%` while preserving the `go test ./...` quality gate. At comparable workload scale, that is about `$1,199/month` on `$5,000/month` of Claude Sonnet API-equivalent coding usage.

The paid plugin allowlist is now intentionally narrow. It keeps the Agent Analyzer workflow and conditionally recommends only reducers that worked in repeated runs: Semble, context-mode, RTK, grepai, and Squeez. Telemetry tools such as ccusage/ccstatusline are not sold as reducers, and negative or too-small results such as claude-context, Probe, Caveman for Claude, claude-rlm, and claude-token-efficient are removed from default remediation advice.

## Local Runthrough

```bash
docker compose up --build
```

Open `http://localhost:8080`, click `Generate Local Commands`, and use the generated local analyze/review/upload flow. The smoke scripts still exercise the legacy token path with `testdata/fixtures/sample-claude.jsonl` for backend compatibility.

Smoke test:

```bash
./scripts/smoke-local.sh
```

If Docker Desktop is unavailable, the same API/worker path can be checked with:

```bash
./scripts/smoke-native.sh
```

Local load gate:

```bash
COMPOSE_PROJECT_NAME=claude-log-analyzer-load docker compose up --build -d
./scripts/load-local.sh 25
COMPOSE_PROJECT_NAME=claude-log-analyzer-load docker compose down -v
```

AWS-backend local smoke with LocalStack:

```bash
./scripts/smoke-aws-local.sh
```

## Development

```bash
go test ./...
go run ./cmd/api
go run ./cmd/worker
```

Useful local env vars:

- `CLAUDE_ANALYZER_DATA_DIR`, default `/tmp/claude-log-analyzer`
- `CLAUDE_ANALYZER_ADDR`, default `:8080`
- `CLAUDE_ANALYZER_WORKER_INTERVAL`, default `2s`

## Privacy Posture

Raw logs are treated as toxic. The launch UX parses and redacts locally, emits aggregate-safe ecosystem IDs only, and uploads only sanitized report JSON. Operational logs forbid raw prompt/tool text.

See [docs/data-retention-and-analytics.md](docs/data-retention-and-analytics.md).

Cloud launch checklist: [docs/cloud-launch-todo.md](docs/cloud-launch-todo.md).

## License

Claude Log Analyzer is source-available, not open source. You may inspect,
clone, and run the software for personal/internal evaluation and development
testing, but production, hosted, commercial, redistribution, and managed-service
uses require a separate written license.

See [LICENSE](LICENSE).
