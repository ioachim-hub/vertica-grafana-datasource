# Security Maintenance and Release Design

## Purpose

Maintain `ioachim-hub/vertica-grafana-datasource` as a small security overlay on
the original Vertica Grafana data source. The upstream repository remains the
behavioral and structural reference. Fork-only changes must be limited to
security remediation, supported build infrastructure, release automation, and
maintenance documentation.

The next fork release will be `v3.2.2`, based on upstream `v3.2.1`.

## Upstream Contract

- `upstream` points to `https://github.com/vertica/vertica-grafana-datasource.git`.
- `origin` points to `https://github.com/ioachim-hub/vertica-grafana-datasource.git`.
- Upstream changes are integrated into the fork with a merge commit so published
  fork history is not rewritten.
- Application behavior and source structure follow `upstream/main`.
- Security changes should remain easy to identify with focused commits.
- Feature work and unrelated refactoring are out of scope.

## Confirmed Maintenance Findings

The baseline audit after syncing upstream `v3.2.1` found:

- `govulncheck ./...` reports reachable `GO-2026-6061` through
  `google.golang.org/grpc@v1.81.1`; the fixed version is `v1.82.1`.
- The Go graph contains deprecated compatibility modules, including
  `github.com/golang/protobuf` and deprecated `golang.org/x/tools` helper
  modules. These are transitive dependencies and should be removed by upgrading
  their direct owners where possible, not by editing indirect requirements
  blindly.
- `yarn audit --groups dependencies` reports 83 findings: 12 low, 59 moderate,
  and 12 high. Most are transitive dependencies owned by Grafana runtime/UI or
  older build tooling.
- `yarn outdated` shows patch/minor updates within current dependency ranges and
  major toolchain migrations. Major migrations are not automatically in scope.
- Backend workflows declare Go 1.21 while `go.mod` and the README declare Go
  1.26.3.
- GitHub workflows still use deprecated major versions such as
  `actions/checkout@v2`, `actions/setup-go@v2`, and `actions/setup-node@v1`.
- Releases are currently created manually from lightweight tags with manually
  uploaded ZIP files. No release workflow exists.

## Dependency Remediation Strategy

Dependency work uses the smallest upstream-compatible change that clears
confirmed findings:

1. Upgrade direct Go dependencies first, beginning with the Grafana plugin SDK
   and Vertica driver, then run `go mod tidy`.
2. Require `google.golang.org/grpc >= v1.82.1` if the direct-owner upgrades do
   not select a fixed version.
3. Use `govulncheck` call-graph results as the Go security gate. Vulnerabilities
   present only in unused transitive code are recorded separately and are not a
   reason for unrelated source changes.
4. Upgrade the aligned Grafana package set together. `@grafana/data`,
   `@grafana/runtime`, and `@grafana/ui` must use the same version.
5. Apply patch/minor updates to direct frontend dependencies and build tools
   when compatible with the existing Grafana 13 and React 19 application.
6. Use Yarn `resolutions` only when a vulnerable transitive dependency has no
   compatible fixed direct owner. Every security resolution must name the
   advisory or vulnerable package in the commit description.
7. Do not adopt unrelated major versions such as ESLint 10, Jest 30, or
   TypeScript 7 unless a remaining confirmed vulnerability cannot be fixed
   within the current major line.

## Continuous Integration

CI will be updated without changing application behavior:

- Use current supported major versions of official GitHub Actions.
- Read the Go version from `go.mod` rather than duplicating a stale version.
- Use Node 20 and Yarn 1.22.19, matching the project metadata.
- Install frontend dependencies with the frozen lockfile.
- Run formatting checks, `go vet`, `go test`, frontend lint, type checking,
  frontend tests, and production builds.
- Run `govulncheck` as a blocking Go security check.
- Run `yarn audit --groups dependencies --level high` as the blocking frontend
  security check. High and critical findings have no default exception; any
  exception requires an explicit repository note naming the advisory, exposure,
  mitigation, and review date.

## Release Workflow

A tag-driven GitHub Actions workflow will build reproducible fork releases:

1. A release commit updates `package.json` to `3.2.2`; the Webpack build injects
   this value into `dist/plugin.json`.
2. The release is tagged `v3.2.2`.
3. The workflow checks that the tag and package version match.
4. It installs the frozen Yarn lockfile and builds the frontend.
5. It builds the Linux AMD64 backend executable through the repository's Mage
   target.
6. It assembles the plugin under a top-level
   `vertica-grafana-datasource/` directory.
7. It creates `vertica-grafana-datasource-3.2.2.zip`.
8. It publishes a GitHub release and attaches the ZIP artifact.

The fork continues to distribute an unsigned plugin unless Grafana signing
credentials are deliberately configured later. The release notes identify the
upstream base tag and summarize only fork-specific security and maintenance
changes.

## AGENTS.md

The repository root will contain `AGENTS.md` with:

- the fork purpose and upstream-as-reference rule;
- the required `rtk` command prefix;
- safe upstream fetch, comparison, and merge commands;
- dependency ownership rules and the ban on unrelated modernization;
- exact Go and frontend audit/test/build commands;
- the version locations and tag convention;
- the release artifact layout and unsigned-plugin constraint;
- a rule not to commit, tag, push, or publish unless the user explicitly asks.

## Verification and Acceptance

The maintenance release is ready when:

- `go mod tidy -diff` is empty;
- `gofmt` and `go vet ./...` pass;
- `go test ./...` passes;
- `govulncheck ./...` reports no reachable vulnerabilities;
- frontend lint, type checking, tests, and production build pass;
- the production dependency audit has no unaccepted high or critical finding;
- the release workflow can produce the expected ZIP from a test tag or local
  equivalent;
- the ZIP contains the generated frontend files, plugin metadata, and Linux
  backend executable under the required top-level plugin directory;
- the diff from upstream contains only the intentional security, CI, release,
  and maintenance-documentation changes.

## Non-Goals

- New datasource features or UI changes.
- Refactoring upstream source code for style.
- Migrating package managers.
- Upgrading every dependency merely because a newer version exists.
- Publishing a signed Grafana plugin without signing credentials and an
  explicit request.
