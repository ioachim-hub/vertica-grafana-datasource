# Changelog

## 3.2.5 (2026-10-12)

- Based on upstream Vertica datasource `v3.2.1`.
- Build with Go 1.26.9 to include the October 8 standard library and toolchain security fixes.
- Upgrade `golang.org/x/net` to v0.60.0 to fix the reported HTTP/2 flow-control, SETTINGS, response-header, trailer, and HPACK vulnerabilities.

## 3.2.4 (2026-10-07)

- Based on upstream Vertica datasource `v3.2.1`.
- Built with Go 1.26.8 to fix standard library vulnerabilities GO-2026-6090 and GO-2026-5972.
- Bumped `google.golang.org/grpc` to v1.83.2 (GO-2026-6443, GO-2026-6348) and OpenTelemetry to v1.47.0 (GO-2026-6505).
- Bumped `grafana-plugin-sdk-go` to v0.296.5 and refreshed other Go dependencies.

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
