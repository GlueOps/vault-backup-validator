#!/bin/bash

# Directory paths
RAFT_DIR="/data/raft"
CONFIG_FILE="vault/config.hcl"
VAULT_TOKEN_FILE="secrets.json"

# Check if /data/raft directory exists, and create it if it doesn't
if [ ! -d "$RAFT_DIR" ]; then
  mkdir -p "$RAFT_DIR"
fi

# Create peers.json in /data/raft directory
cat > "$RAFT_DIR/peers.json" << EOF
[
  {
    "id": "node1",
    "address": "127.0.0.1:8201",
    "non_voter": false
  }
]
EOF

# Start Vault server with the specified config file
nohup vault server -config="$CONFIG_FILE" > vault.log 2>&1 &

# Wait for Vault to start (adjust sleep duration as needed)
sleep 30

# Set VAULT_ADDR for API interactions
export VAULT_ADDR=http://127.0.0.1:8200

# Initialize Vault and store the unseal token and root token
vault_data=$(vault operator init -key-shares=1 -key-threshold=1 --format=json)
echo $vault_data
root_token=$(echo "$vault_data" | jq -r .root_token)
unseal_token=$(echo "$vault_data" | jq -r .unseal_keys_b64[0])


# Store the unseal and root token in a file
json_object=$(jq -n --arg key "$unseal_token" --arg token "$root_token" '{ "key": $key , "token": $token }')
echo "$json_object" > "$VAULT_TOKEN_FILE"
