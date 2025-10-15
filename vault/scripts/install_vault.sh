#!/bin/bash
echo "BUILT IN OPENBAO_VERSION: ${CACHED_OPENBAO_VERSION}"
PROVIDED_OPENBAO_VERSION="$1" # Get the Vault version from the first command-line argument

cd ~
if [ "$CACHED_OPENBAO_VERSION" != "$PROVIDED_VERSION" ]; then
    echo "Removing cached bao and installing: ${PROVIDED_VERSION}"
    rm -f /usr/bin/bao
    wget -O https://github.com/openbao/openbao/releases/download/v${PROVIDED_OPENBAO_VERSION}/bao_${PROVIDED_OPENBAO_VERSION}_Linux_x86_64.tar.gz /tmp/bao_${PROVIDED_OPENBAO_VERSION}_Linux_x86_64.tar.gz
    tar -xzvf /tmp/bao_${PROVIDED_OPENBAO_VERSION}_Linux_x86_64.tar.gz bao && mv bao /usr/bin/bao && rm /tmp/bao_${PROVIDED_OPENBAO_VERSION}_Linux_x86_64.tar.gz
fi


