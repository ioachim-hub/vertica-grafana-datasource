# Vertica Datasource Security Maintenance and Release Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Preserve upstream Vertica datasource behavior while removing confirmed dependency vulnerabilities, modernizing supported CI actions, documenting fork maintenance, and producing the unsigned `v3.2.2` release artifact.

**Architecture:** Treat `upstream/main` at `v3.2.1` as the application reference and keep the fork as a thin maintenance overlay. Remediate vulnerabilities through direct dependency owners first, add narrowly justified transitive constraints only when necessary, and automate validation and packaging without modifying datasource behavior.

**Tech Stack:** Go 1.26.3, Grafana plugin SDK, Vertica Go driver, Grafana 13/React 19, TypeScript, Yarn Classic 1.22.19, Webpack, Mage, GitHub Actions.

## Global Constraints

- The upstream repository remains the behavioral and structural reference.
- Fork-only changes are limited to security remediation, supported build infrastructure, release automation, and maintenance documentation.
- The next fork release is `v3.2.2`, based on upstream `v3.2.1`.
- `@grafana/data`, `@grafana/runtime`, and `@grafana/ui` must use the same version.
- Do not adopt ESLint 10, Jest 30, or TypeScript 7 unless a confirmed vulnerability cannot be fixed within the current major line.
- Use Yarn `resolutions` only when no compatible direct-owner upgrade selects a fixed transitive version.
- The release remains unsigned.
- Prefix repository shell commands with `rtk`.
- Do not commit, tag, push, or publish unless the user explicitly requests it.

---

### Task 1: Add Fork Maintenance Instructions

**Files:**
- Create: `AGENTS.md`

**Interfaces:**
- Consumes: repository remotes `origin` and `upstream`, existing package scripts, Mage Linux target, and the approved design.
- Produces: the authoritative instructions agents must follow when syncing, auditing, changing, validating, and releasing this fork.

- [ ] **Step 1: Write `AGENTS.md`**

Create a concise root document containing these exact operational rules and commands:

```markdown
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
```

- [ ] **Step 2: Verify the instructions against repository metadata**

Run:

```bash
rtk git remote -v
rtk node -p "require('./package.json').packageManager"
rtk node -p "require('./package.json').scripts"
rtk go run github.com/magefile/mage@v1.17.2 -l
```

Expected: `origin` and `upstream` match the documented repositories, Yarn is
`1.22.19`, documented scripts exist, and Mage lists the Linux/default build.

- [ ] **Step 3: Check the document**

Run:

```bash
rtk git diff --check -- AGENTS.md
rtk rg -n "upstream/main|govulncheck|frozen-lockfile|Do not commit" AGENTS.md
```

Expected: no whitespace errors and all four policy anchors are present.

---

### Task 2: Remediate the Go Dependency Graph

**Files:**
- Modify: `go.mod:5`
- Modify: `go.sum`

**Interfaces:**
- Consumes: current Go imports under `pkg/` and `Magefile.go`.
- Produces: a tidy Go graph with no reachable `govulncheck` findings and no application-code changes.

- [ ] **Step 1: Preserve the failing security baseline**

Run:

```bash
rtk go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Expected before remediation: exit status 3 with `GO-2026-6061` reachable through
`google.golang.org/grpc@v1.81.1`.

- [ ] **Step 2: Upgrade direct owners and the vulnerable module**

Run:

```bash
rtk go get github.com/grafana/grafana-plugin-sdk-go@v0.294.0
rtk go get github.com/vertica/vertica-sql-go@v1.3.8
rtk go get google.golang.org/grpc@v1.82.1
rtk go mod tidy
```

Expected: `go.mod` selects SDK `v0.294.0`, Vertica driver `v1.3.8`, and gRPC at
least `v1.82.1`; `go.sum` reflects the resulting minimal module graph.

- [ ] **Step 3: Verify graph ownership and deprecated compatibility modules**

Run:

```bash
rtk go mod why -m google.golang.org/grpc
rtk go mod graph | rtk rg 'github.com/golang/protobuf|golang.org/x/tools/go/(expect|packages/packagestest)'
```

Expected: gRPC has a valid dependency path. Record any remaining deprecated
modules as transitive-only; do not add unrelated direct requirements merely to
remove their names.

- [ ] **Step 4: Run the complete backend gate**

Run:

```bash
rtk go mod tidy -diff
rtk gofmt -l pkg/*.go Magefile.go
rtk go vet ./...
rtk go test ./...
rtk go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

Expected: tidy has no diff, formatting prints no files, vet passes, all backend
tests pass, and `govulncheck` reports no reachable vulnerability.

---

### Task 3: Remediate Frontend Production Dependencies

**Files:**
- Modify: `package.json:75`
- Modify: `yarn.lock`

**Interfaces:**
- Consumes: the existing Grafana 13/React 19 source contract.
- Produces: an aligned Grafana dependency set and a frozen lockfile with no unaccepted high or critical production finding.

- [ ] **Step 1: Capture the current production audit**

Run:

```bash
rtk yarn install --frozen-lockfile
rtk yarn audit --groups dependencies --level high
```

Expected before remediation: audit failure containing high-severity findings
from the existing lockfile.

- [ ] **Step 2: Upgrade the aligned runtime owners and safe direct patches**

Run:

```bash
rtk yarn add --exact @grafana/data@13.1.1 @grafana/runtime@13.1.1 @grafana/ui@13.1.1
rtk yarn add --exact react@19.2.8 react-dom@19.2.8 react-use@17.6.1 tslib@2.8.1
rtk yarn add --dev --exact @grafana/e2e-selectors@13.1.1 @swc/core@1.15.47
```

Keep `@grafana/e2e@10.4.12`, ESLint 8, Jest 29, and TypeScript 5.3 unless a
remaining confirmed finding proves one of those owners must change.

- [ ] **Step 3: Trace every remaining high or critical finding**

Run:

```bash
rtk yarn audit --groups dependencies --level high
rtk yarn why uuid
rtk yarn why dompurify
rtk yarn why react-router
```

Expected: each remaining high/critical advisory has a visible direct-owner path.
Fix it by upgrading that owner within its current major. If no compatible owner
version exists, add the narrowest Yarn resolution that reaches the advisory's
fixed version and document the advisory in the change description. Do not force
`react-router@7` into a dependency expecting React Router 6.

- [ ] **Step 4: Run the complete frontend gate**

Run:

```bash
rtk yarn install --frozen-lockfile
rtk yarn lint
rtk yarn typecheck
rtk yarn test:ci
rtk yarn build
rtk yarn audit --groups dependencies --level high
```

Expected: installation is reproducible; lint, typecheck, tests, and build pass;
the audit has no unaccepted high or critical production finding.

---

### Task 4: Align and Harden Continuous Integration

**Files:**
- Modify: `.github/workflows/build.yml`
- Modify: `.github/workflows/backend-unit-testing.yml`
- Modify: `.github/workflows/frontend-unit-testing.yml`
- Create: `.github/workflows/security.yml`

**Interfaces:**
- Consumes: Go version from `go.mod`, Node 20, Yarn lockfile, and package scripts.
- Produces: supported CI actions and repeatable build/test/security gates.

- [ ] **Step 1: Modernize backend build CI**

Replace obsolete actions in `build.yml` with `actions/checkout@v4` and
`actions/setup-go@v5`, configure `go-version-file: go.mod` and cache, and replace
the mutating format command with:

```yaml
- name: Check formatting
  run: test -z "$(gofmt -l pkg/*.go Magefile.go)"
- name: Vet
  run: go vet ./...
- name: Build Linux backend
  run: go run github.com/magefile/mage@v1.17.2 -v
```

- [ ] **Step 2: Modernize backend unit-test CI**

Use `ubuntu-latest`, `actions/checkout@v4`, and `actions/setup-go@v5` with
`go-version-file: go.mod`. Remove the stale GOPATH/debug block and run:

```yaml
- name: Test
  run: go test -v ./...
```

- [ ] **Step 3: Make frontend CI reproducible**

Remove the obsolete commented Node 14 job. Use `actions/checkout@v4` and
`actions/setup-node@v4` with Node 20 and Yarn caching. Run:

```yaml
- run: yarn install --frozen-lockfile
- run: yarn lint
- run: yarn typecheck
- run: yarn test:ci
- run: yarn build
```

- [ ] **Step 4: Add security CI**

Create `security.yml` for pull requests, pushes to `main`, and weekly
`schedule: cron: '0 6 * * 1'`. Give it:

- a Go job using checkout v4, setup-go v5 with `go-version-file`, and
  `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`;
- a frontend job using checkout v4, setup-node v4 with Node 20/Yarn cache,
  frozen installation, and
  `yarn audit --groups dependencies --level high`.

- [ ] **Step 5: Validate workflow syntax and local equivalents**

Run:

```bash
rtk npx --yes prettier@3.9.6 --check '.github/workflows/*.yml'
rtk go vet ./...
rtk go test ./...
rtk yarn lint
rtk yarn typecheck
rtk yarn test:ci
rtk yarn build
```

Expected: workflow YAML is formatted and all local CI equivalents pass.

---

### Task 5: Add Reproducible Release Packaging

**Files:**
- Create: `scripts/package-release.sh`
- Create: `.github/workflows/release.yml`

**Interfaces:**
- Consumes: `package.json` version, Git tag, Webpack output in `dist/`, Mage Linux AMD64 binary, and plugin ID `vertica-grafana-datasource`.
- Produces: `release/vertica-grafana-datasource-<version>.zip` with a single required top-level plugin directory.

- [ ] **Step 1: Add a fail-fast packaging script**

Create `scripts/package-release.sh` that:

1. reads the version from `package.json`;
2. accepts an optional tag argument and rejects it unless it equals `v$version`;
3. runs `yarn install --frozen-lockfile`, `yarn build`, and
   `go run github.com/magefile/mage@v1.17.2 -v`;
4. verifies `dist/plugin.json` contains the package version;
5. verifies `dist/gpx_vertica-grafana-plugin_linux_amd64` exists and is
   executable;
6. copies `dist/` to a temporary
   `vertica-grafana-datasource/` directory;
7. creates
   `release/vertica-grafana-datasource-$version.zip`;
8. removes only its temporary directory.

Use `mktemp -d` and a shell `trap` for cleanup. Do not remove `dist/` or any
broad path.

- [ ] **Step 2: Add a tag-triggered release workflow**

Create `release.yml` triggered by `v*.*.*` tags with:

- `contents: write`;
- checkout v4;
- setup-go v5 using `go.mod`;
- setup-node v4 using Node 20 and Yarn cache;
- `bash scripts/package-release.sh "${GITHUB_REF_NAME}"`;
- `gh release create "$GITHUB_REF_NAME" release/*.zip --generate-notes
  --title "$GITHUB_REF_NAME"` using `GH_TOKEN: ${{ github.token }}`.

- [ ] **Step 3: Test version mismatch failure**

Run:

```bash
rtk bash scripts/package-release.sh v0.0.0
```

Expected: non-zero exit before dependency installation or build, with a message
that the tag does not match `package.json`.

- [ ] **Step 4: Build and inspect the release artifact locally**

Run:

```bash
rtk bash scripts/package-release.sh v3.2.1
rtk unzip -l release/vertica-grafana-datasource-3.2.1.zip
```

Expected: one top-level `vertica-grafana-datasource/` directory containing
`plugin.json`, frontend assets, and
`gpx_vertica-grafana-plugin_linux_amd64`.

---

### Task 6: Prepare Version `3.2.2` and Final Release Evidence

**Files:**
- Modify: `package.json:3`
- Modify: `CHANGELOG.md`
- Modify: `README.md`

**Interfaces:**
- Consumes: completed security changes and upstream base `v3.2.1`.
- Produces: release-ready metadata and copy-ready evidence; it does not create the tag or GitHub release without separate user authorization.

- [ ] **Step 1: Set the release version**

Change `package.json` from `"version": "3.2.1"` to `"version": "3.2.2"` without
changing `packageManager`.

- [ ] **Step 2: Replace the stale changelog entry**

Add a `3.2.2` entry dated `2026-07-30` that states:

- based on upstream Vertica datasource `v3.2.1`;
- fixed reachable Go vulnerabilities;
- refreshed vulnerable frontend production dependencies;
- updated supported GitHub Actions and aligned tool versions;
- added reproducible unsigned Linux AMD64 release packaging.

- [ ] **Step 3: Update maintenance/release documentation**

Update README prerequisites and build instructions to match Go 1.26.3, Node 20,
Yarn 1.22.19, the Linux-only Mage target, unsigned fork behavior, and the local
packaging command:

```bash
rtk bash scripts/package-release.sh v3.2.2
```

- [ ] **Step 4: Run the full release gate**

Run:

```bash
rtk go mod tidy -diff
rtk gofmt -l pkg/*.go Magefile.go
rtk go vet ./...
rtk go test ./...
rtk go run golang.org/x/vuln/cmd/govulncheck@latest ./...
rtk yarn install --frozen-lockfile
rtk yarn lint
rtk yarn typecheck
rtk yarn test:ci
rtk yarn build
rtk yarn audit --groups dependencies --level high
rtk bash scripts/package-release.sh v3.2.2
rtk unzip -l release/vertica-grafana-datasource-3.2.2.zip
rtk git diff --check
rtk git diff --stat upstream/main
```

Expected: all gates pass, the ZIP layout is correct, and the upstream diff is
limited to intentional security, CI, release, and maintenance files.

- [ ] **Step 5: Stop before external publication**

Report the exact tested commit, audit results, test counts, ZIP checksum, and
artifact path. Ask for explicit authorization before committing, tagging,
pushing, or running the tag-triggered GitHub release.
