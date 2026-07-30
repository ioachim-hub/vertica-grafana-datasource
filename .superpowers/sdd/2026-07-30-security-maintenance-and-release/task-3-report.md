# Task 3 Report: Frontend Production Dependencies

## Changes

- Aligned `@grafana/data`, `@grafana/runtime`, `@grafana/ui`, and
  `@grafana/e2e-selectors` at `13.1.1`.
- Updated the requested safe direct patches: React and React DOM `19.2.8`,
  react-use `17.6.1`, tslib `2.8.1`, and SWC core `1.15.47`.
- Kept `@grafana/e2e` at `10.4.12`, ESLint 8, Jest 29, and TypeScript 5.3.
- Added scoped Yarn resolutions for the remaining production advisory paths:
  - Grafana i18n descendants use brace-expansion `5.0.8`.
  - Grafana runtime/UI use Immutable `5.1.9`; Sass retains its compatible
    Immutable 4 release.
  - Slate and Slate React use Lodash `4.18.1`; the unrelated Grafana e2e
    development dependency retains Lodash `4.17.21`.
- Regenerated `yarn.lock`.

## Audit trace

The baseline production audit reported 83 findings, including 12 high findings:

- Lodash through `@grafana/ui > slate`.
- brace-expansion through
  `@grafana/data > @grafana/i18n > @typescript-eslint/... > minimatch`.
- Immutable through `@grafana/ui`.

After the aligned upgrades:

- `uuid` remains only under the pinned development-only `@grafana/e2e` tree.
- DOMPurify is `3.4.12`, owned by `@grafana/data`.
- React Router 5 remains behind Grafana UI's React Router DOM 5 compatibility
  path and React Router 6.30.3 remains under `react-router-dom-v5-compat`; no
  incompatible React Router 7 override was added.
- The scoped resolutions remove the remaining production high findings without
  changing incompatible majors for unrelated owners.
- Final audit summary: 18 moderate, 0 high, 0 critical. Yarn Classic returns
  exit code 4 because moderate findings remain even with `--level high`; no
  high or critical finding remains to accept.

## Verification

- `yarn install --frozen-lockfile`: passed.
- `yarn typecheck`: passed.
- `yarn test:ci`: passed, 4 suites and 19 tests.
- `yarn build`: passed.
- `yarn audit --groups dependencies --level high`: no high or critical
  findings; 18 moderate findings remain and produce Yarn Classic exit code 4.
- `yarn lint`: failed on pre-existing source-format issues outside this task's
  package/lockfile scope:
  - `src/QueryEditor.test.tsx`: five `react/react-in-jsx-scope` errors.
  - `src/module.js`: missing final newline.
  - `webpack.config.ts`: missing final newline.

## Concerns

- Grafana 13.1.1 declares React 18 peer ranges while this repository's existing
  source contract uses React 19; Yarn reports the same peer-range warnings
  during installation.
- The lint baseline must be repaired separately if a completely green aggregate
  frontend gate is required.
