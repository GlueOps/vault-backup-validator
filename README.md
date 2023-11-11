# Vault Backup Validatior

This API validates Vault backups using a POST request. It checks the integrity of the backup data to ensure it can be safely restored.

## Endpoint

- `POST /api/v1/validate`

## Request Parameters

Send a JSON object with the following fields:

- `source_backup_url` (string): URL of the source Vault backup.
- `source_keys_url` (string): URL of the source Vault unseal keys and token.
- `path_values_map` (json): key-pair values of secret_path and the expected key-values in that path.
## Response

- Status 200 OK: Validation successful.
- Status 400 Bad Request: Validation unsuccessful/Invalid or missing parameters.
- Status 500 Internal Server Error: Unexpected error.

Sample Request:

POST /api/v1/validate

```json
{
    "source_backup_url": "https://example.com/backup.zip",
    "source_keys_url": "https://example.com/keys.zip",
    "path_values_map":{
        "secret/key-1-for-balaji": {
            "key1":"value1",
            "key2":"value-2"
        }
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