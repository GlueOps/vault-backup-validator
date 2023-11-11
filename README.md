# Vault Backup Validatior

This API validates Vault backups using a POST request. It checks the integrity of the backup data to ensure it can be safely restored.

## Endpoint

- `POST /api/v1/validate`

## Request Parameters

Send a JSON object with the following fields:

- `source_backup_url` (string): URL of the source Vault backup.
- `source_keys_url` (string): URL of the source Vault unseal keys and token.
- `path_values_map` (json): key-pair values of secret_path and the expected key-values in that path.
- `vault_version` (string): version of vault where the backup is taken from.
## Response

- Status 200 OK: Validation successful.
- Status 400 Bad Request: Validation unsuccessful/Invalid or missing parameters.
- Status 500 Internal Server Error: Unexpected error.

Sample Request:

POST /api/v1/validate

```json
{
    "source_backup_url": "https://example.com/backup.snap",
    "source_keys_url": "https://example.com/keys.json",
    "path_values_map":{
        "secret/key-1-for-balaji": {
            "key1":"value1",
            "key2":"value-2"
        }
    },
    "vault_version":"1.15.0"
}
```

keys.json Expected format

```json
{
    {
  "keys": [
    "key_in_hexa"
  ],
  "keys_base64": [
    "key_in_base64"
  ],
  "root_token": "root_token_string"
}
}
```

Sample Response

```json
{
    "message":"Backup is valid",
    "status":"success"
}
```