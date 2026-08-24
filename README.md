# azmetrics-demo

A minimal Go service that lists Azure Monitor **metric alert rules** for a subscription,
using the public `armmonitor` SDK. It exists as a **demo target** for the
*Licence to Patch* dependency-update trust agent.

## The landmine it demonstrates

The service is pinned to `armmonitor v0.12.0`. That version sends its metric-alert
requests with `api-version=2024-03-01-preview`.

Bumping to **`armmonitor v0.13.0`** silently changes the baked-in api-version to
**`2026-01-01`**, which Azure Resource Manager rejects at runtime with
`404 InvalidResourceType` for `microsoft.insights` metric-alert resource types.

Critically:

- **It compiles clean** — no API signature changed.
- **The test suite stays green** — `alerts/alerts_test.go` uses a mock transport that
  never reaches real ARM, so it cannot observe the api-version on the wire. (Verified:
  the v0.13.0 client sends `api-version=2026-01-01`, but the test asserts only success.)
- **Reachability tools see nothing** — no removed or changed symbol; the api-version is
  a constant inside the dependency, not a call site in this repo.
- **The changelog is silent** — `armmonitor`'s v0.13.0 release notes do not mention the
  api-version change.

So the only way to catch this before production is to **diff the dependency's source
between v0.12.0 and v0.13.0** and notice the changed api-version constant
(`version20240301Preview` → `version20260101`). That is what the agent does.

## Layout

- `alerts/alerts.go` — `ListAlertNames`, the call site depending on the metric-alerts api-version.
- `alerts/alerts_test.go` — a green unit test that is blind to the regression (the point).

## Run

```
go test ./...   # green on both v0.12.0 and v0.13.0
```

> Rebuilt clean-room from the public Azure SDK for Go. Not derived from any private codebase.
