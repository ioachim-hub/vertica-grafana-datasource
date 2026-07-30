#!/usr/bin/env bash

set -euo pipefail

readonly plugin_id="vertica-grafana-datasource"
readonly script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
readonly repo_root="$(cd -- "${script_dir}/.." && pwd)"

cd "${repo_root}"

if (( $# > 1 )); then
  echo "Usage: $0 [v<package-version>]" >&2
  exit 2
fi

version="$(node -p "require('./package.json').version")"
readonly version
readonly expected_tag="v${version}"
readonly supplied_tag="${1:-}"

if [[ -n "${supplied_tag}" && "${supplied_tag}" != "${expected_tag}" ]]; then
  echo "Release tag '${supplied_tag}' does not match package.json version '${expected_tag}'." >&2
  exit 1
fi

yarn install --frozen-lockfile
yarn build
go run github.com/magefile/mage@v1.17.2 -v

plugin_version="$(node -p "require('./dist/plugin.json').info.version")"
if [[ "${plugin_version}" != "${version}" ]]; then
  echo "dist/plugin.json version '${plugin_version}' does not match package.json version '${version}'." >&2
  exit 1
fi

readonly backend_binary="dist/gpx_vertica-grafana-plugin_linux_amd64"
if [[ ! -x "${backend_binary}" ]]; then
  echo "Required executable '${backend_binary}' is missing or is not executable." >&2
  exit 1
fi

temporary_root="$(cd -- "${TMPDIR:-/tmp}" && pwd)"
readonly temporary_root
temporary_dir="$(mktemp -d "${temporary_root}/${plugin_id}.release.XXXXXXXX")"
readonly temporary_dir
cleanup() {
  case "${temporary_dir:-}" in
    "${temporary_root}/${plugin_id}.release."????????)
      [[ ! -d "${temporary_dir}" ]] || rm -rf -- "${temporary_dir}"
      ;;
    "")
      ;;
    *)
      echo "Refusing to clean unexpected temporary path '${temporary_dir}'." >&2
      ;;
  esac
}
trap cleanup EXIT

readonly staged_plugin="${temporary_dir}/${plugin_id}"
readonly temporary_archive="${temporary_dir}/${plugin_id}-${version}.zip"
readonly release_dir="${repo_root}/release"
readonly release_archive="${release_dir}/${plugin_id}-${version}.zip"

cp -R dist "${staged_plugin}"

# Normalize archive metadata so identical build output produces identical ZIPs.
source_date_epoch="${SOURCE_DATE_EPOCH:-$(git log -1 --format=%ct)}"
if [[ ! "${source_date_epoch}" =~ ^[0-9]+$ ]]; then
  echo "SOURCE_DATE_EPOCH must be a non-negative integer." >&2
  exit 1
fi
if (( source_date_epoch < 315532800 )); then
  source_date_epoch=315532800
fi
find "${staged_plugin}" -exec touch -h -d "@${source_date_epoch}" {} +

(
  cd "${temporary_dir}"
  find "${plugin_id}" -print | LC_ALL=C sort | zip -q -X "${temporary_archive}" -@
)

mkdir -p "${release_dir}"
mv -f -- "${temporary_archive}" "${release_archive}"

echo "Created ${release_archive}"
