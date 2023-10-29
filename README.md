# Vault Backup Validatior

This API validates Vault backups using a POST request. It checks the integrity of the backup data to ensure it can be safely restored.

## Endpoint

- `POST /api/v1/restore`

## Request Parameters

Send a JSON object with the following fields:

- `source_backup_url` (string): URL of the source Vault backup.
- `source_keys_url` (string): URL of the source Vault encryption keys (if applicable).
- `source_token_url` (string): URL of the source Vault token (if applicable).
- `destination_vault_url` (string): URL of the destination Vault.
- `destination_vault_token` (string): Token for the destination Vault.

## Response

- Status 200 OK: Validation successful.
- Status 400 Bad Request: Invalid or missing parameters.
- Status 403 Unauthorized: Permission denied.
- Status 500 Internal Server Error: Unexpected error.

Sample Request:

```json
POST /api/validate-vault-backup
{
    "source_backup_url": "https://example.com/backup.zip",
    "source_keys_url": "https://example.com/keys.zip",
    "source_token_url": "https://example.com/token.txt",
    "destination_vault_url": "https://destination-vault-url",
    "destination_vault_token": "destination-token"
}
```

