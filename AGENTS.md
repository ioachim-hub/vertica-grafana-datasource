# Repository Maintenance Guide

## Purpose

This repository is a security-maintained fork of
`vertica/vertica-grafana-datasource`. Treat `upstream/main` as the behavioral
and structural reference. Keep fork-only changes limited to security fixes,
supported CI/release infrastructure, and maintenance documentation. Do not add
features or refactor upstream application code unless explicitly requested.

## Commands

Prefix shell commands with `rtk` as required by
`/home/ioachimlihor/.codex/RTK.md`.

## Syncing Upstream

1. Require a clean working tree: `rtk git status --short --branch`.
2. Fetch both remotes: `rtk git fetch --prune origin` and
   `rtk git fetch --prune upstream`.
3. Compare histories with
   `rtk git rev-list --left-right --count main...upstream/main`.
4. Inspect fork-only commits with `rtk git log --oneline upstream/main..main`.
5. Merge `upstream/main` into `main`; do not rebase published fork commits.
6. Resolve conflicts in favor of newer upstream behavior, retaining only
   still-relevant fork security/release changes.
7. Validate before any requested push.

## Dependency Policy

- Upgrade direct dependency owners before constraining transitive dependencies.
- Keep `@grafana/data`, `@grafana/runtime`, and `@grafana/ui` aligned.
- Do not perform unrelated major toolchain upgrades.
- Every Yarn resolution must document the vulnerable package/advisory and why a
  direct-owner upgrade was insufficient.
- Do not hide audit findings by changing audit thresholds.

### Yarn Resolution Register

| Resolution | Owner and rationale |
| --- | --- |
| `**/@grafana/i18n/**/brace-expansion` | Fixes the production `brace-expansion` advisory path owned by Grafana i18n; the aligned Grafana 13.1.1 packages do not select the fixed version. |
| `**/@grafana/runtime/immutable` | Fixes the production `immutable` advisory path owned by Grafana runtime while leaving unrelated major versions unchanged. |
| `**/@grafana/ui/immutable` | Fixes the production `immutable` advisory path owned by Grafana UI while leaving Sass on its compatible Immutable 4 release. |
| `**/slate/lodash` | Fixes the production `lodash` advisory path owned by Grafana UI's Slate dependency without overriding unrelated owners. |
| `**/slate-react/lodash` | Fixes the production `lodash` advisory path owned by Grafana UI's Slate React dependency without changing the development-only Grafana e2e tree. |
| `flatted` | Inherited from upstream `v3.2.1`; retains upstream's constraint for the ESLint/flat-cache development dependency path. |
| `js-cookie` | Inherited from upstream `v3.2.1`; retains upstream's constraint for the `react-use` dependency path. |
| `path-to-regexp` | Inherited from upstream `v3.2.1`; retains upstream's constraint for Grafana UI's React Router 5 compatibility path. |
| `protobufjs` | Inherited from upstream `v3.2.1` as a compatibility guard; no `protobufjs` package is currently selected in the installed dependency graph. |

## Validation

Backend:

`rtk go mod tidy -diff`
`rtk gofmt -l pkg/*.go Magefile.go`
`rtk go vet ./...`
`rtk go test ./...`
`rtk go run golang.org/x/vuln/cmd/govulncheck@latest ./...`

Frontend:

`rtk yarn install --frozen-lockfile`
`rtk yarn lint`
`rtk yarn typecheck`
`rtk yarn test:ci`
`rtk yarn build`
`rtk yarn audit --groups dependencies --level high`

## Version and Release

- `package.json` is the release version source.
- Webpack injects that version into `dist/plugin.json`.
- Tags use `vMAJOR.MINOR.PATCH`; the tag and package version must match.
- Mage builds the Linux AMD64 executable into `dist/`.
- Release ZIPs contain a top-level `vertica-grafana-datasource/` directory and
  are named `vertica-grafana-datasource-MAJOR.MINOR.PATCH.zip`.
- Releases are unsigned unless signing is explicitly requested and configured.
- Release notes identify the upstream base and list only fork-specific changes.

Do not commit, tag, push, create a GitHub release, or publish artifacts unless
the user explicitly asks.
