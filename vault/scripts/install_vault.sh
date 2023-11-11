#!/bin/bash

VAULT_VERSION="$1" # Get the Vault version from the first command-line argument
DOWNLOAD_URL="https://releases.hashicorp.com/vault/${VAULT_VERSION}/vault_${VAULT_VERSION}_linux_amd64.zip"

cd ~
# Check if Vault is installed  TODO: change it to /usr/local/bin
if [ /usr/bin/vault ]; then
    # Remove existing Vault 
    sudo rm /usr/bin/vault
fi

# Download Vault
wget -O vault.zip "${DOWNLOAD_URL}"

# Unzip the downloaded file without prompting
unzip -o vault.zip

# Make Vault executable
chmod +x vault

# Move Vault to a directory in your PATH without prompting
sudo mv -f vault /usr/bin/

# Clean up
rm vault.zip
