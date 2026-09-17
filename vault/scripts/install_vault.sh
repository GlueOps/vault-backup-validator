#!/bin/bash
echo "BUILT IN OPENBAO_VERSION: ${CACHED_OPENBAO_VERSION}"
PROVIDED_OPENBAO_VERSION="$1" # Get the Vault version from the first command-line argument

# OpenBao renamed its release tarballs starting with v2.6.0:
#   <= 2.5.x : bao_<ver>_Linux_x86_64.tar.gz
#   >= 2.6.0 : openbao_<ver>_linux_amd64.tar.gz
# Backups can come from either era, so try the current name first and fall
# back to the legacy one.
download_bao() {
    local version="$1"
    local dest="$2"
    local base="https://github.com/openbao/openbao/releases/download/v${version}"
    local name
    for name in "openbao_${version}_linux_amd64.tar.gz" "bao_${version}_Linux_x86_64.tar.gz"; do
        echo "Trying ${base}/${name}"
        if wget -q -O "${dest}" "${base}/${name}"; then
            return 0
        fi
        rm -f "${dest}"
    done
    echo "ERROR: could not download OpenBao ${version} under either release asset name" >&2
    return 1
}

cd ~
if [ "$CACHED_OPENBAO_VERSION" != "$PROVIDED_OPENBAO_VERSION" ]; then
    echo "Removing cached bao and installing: ${PROVIDED_OPENBAO_VERSION}"
    rm -f /usr/bin/bao
    TARBALL="/tmp/openbao_${PROVIDED_OPENBAO_VERSION}.tar.gz"
    download_bao "${PROVIDED_OPENBAO_VERSION}" "${TARBALL}" || exit 1
    tar -xzvf "${TARBALL}" bao && mv bao /usr/bin/bao && rm -f "${TARBALL}" || exit 1
fi
