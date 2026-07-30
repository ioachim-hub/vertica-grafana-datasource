# Changelog

## 3.2.3 (2026-07-30)

- Based on upstream Vertica datasource `v3.2.1`.
- Rejected restored datasource credentials containing invalid UTF-8 before building connection strings.
- Sanitized backend errors and health messages before Grafana protobuf serialization.
- Encoded binary query results as hexadecimal and repaired malformed textual results.

## 3.2.2 (2026-07-30)

- Based on upstream Vertica datasource `v3.2.1`.
- Fixed reachable Go vulnerabilities.
- Refreshed vulnerable frontend production dependencies.
- Updated supported GitHub Actions and aligned tool versions.
- Added reproducible unsigned Linux AMD64 release packaging.
