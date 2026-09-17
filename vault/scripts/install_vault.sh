#!/bin/bash
# Ensure the requested OpenBao version is the one installed at /usr/bin/bao.
#
# The decision is made by inspecting the installed binary (`bao version`), not a
# build-time env var: the validator is a long-running server and a previous
# request may already have swapped the binary for a different version.
set -euo pipefail

requested="${1:?OpenBao version required}"
requested="${requested#v}"

# Release assets are fetched through the GlueOps Nexus raw proxy for GitHub
# (same one GlueKube uses), with github.com as a fallback if the proxy is
# unreachable or does not have the asset.
#
# OpenBao renamed its release tarballs starting with v2.6.0:
#   <= 2.5.x : bao_<ver>_Linux_x86_64.tar.gz
#   >= 2.6.0 : openbao_<ver>_linux_amd64.tar.gz
# Backups can come from either era, so try the current name first and fall
# back to the legacy one.
GITHUB_PROXY_BASE="https://repo.gpkg.io/repository/raw-github"
GITHUB_BASE="https://github.com"

download_bao() {
    local version="$1"
    local dest="$2"
    local path="openbao/openbao/releases/download/v${version}"
    local name base
    # Exhaust the proxy (both names) before touching github.com.
    for base in "${GITHUB_PROXY_BASE}" "${GITHUB_BASE}"; do
        for name in "openbao_${version}_linux_amd64.tar.gz" "bao_${version}_Linux_x86_64.tar.gz"; do
            echo "Trying ${base}/${path}/${name}"
            if wget -q --timeout=30 --tries=2 -O "${dest}" "${base}/${path}/${name}"; then
                return 0
            fi
            rm -f "${dest}"
        done
    done
    echo "ERROR: could not download OpenBao ${version} under either release asset name from the proxy or github.com" >&2
    return 1
}

# `|| true` is required: under pipefail a missing or non-executable bao would
# abort the script instead of yielding an empty string.
installed="$(bao version 2>/dev/null | awk 'NR==1{print $2}' | sed 's/^v//' || true)"

if [ "${installed}" = "${requested}" ]; then
    echo "OpenBao ${installed} already installed, skipping download"
    exit 0
fi
echo "Installed OpenBao: ${installed:-none}. Installing: ${requested}"

tmpdir="$(mktemp -d)"
trap 'rm -rf "${tmpdir}"' EXIT

download_bao "${requested}" "${tmpdir}/bao.tar.gz" || exit 1
tar -xzf "${tmpdir}/bao.tar.gz" -C "${tmpdir}" bao
# mv (rename), never cp: a rename is safe even if a previous bao process is
# still exiting, whereas cp over a running binary fails with ETXTBSY.
mv "${tmpdir}/bao" /usr/bin/bao
chmod 0755 /usr/bin/bao
echo "Installed OpenBao $(bao version | awk 'NR==1{print $2}')"
