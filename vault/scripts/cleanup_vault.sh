#!/bin/bash

# Cleanup vault data
rm -r /data/raft
rm /data/vault.db

# Kill Vault process
vault_pid=$(pgrep -x bao)
kill -s 15 $vault_pid
