# Vault Backup Validatior

This API validates OpenBao backups using a POST request. It checks the integrity of the backup data to ensure it can be safely restored.

## Endpoint

- `POST /api/v1/validate`

## Request Parameters

Send a JSON object with the following fields:

- `source_backup_url` (string): URL of the source OpenBao backup.
- `source_keys_url` (string): URL of the source OpenBao unseal keys and token.
- `path_values_map` (json): key-pair values of secret_path and the expected key-values in that path.
- `vault_version` (string): version of OpenBao where the backup is taken from.
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
    "vault_version":"2.4.0"
}
```

keys.json Expected format

```json
{
  "keys": [
    "key_in_hexa"
  ],
  "keys_base64": [
    "key_in_base64"
  ],
  "root_token": "root_token_string"
}
```

Sample Response

```json
{
    "message":"Backup is valid",
    "status":"success"
}
```

To dev locally:
- Easiest to just create presigned urls using your existing tenant (e.g. tenant-foobar) and then finding the the required files keys.json + the *.snap backup file you want to restore/validate with.

Sample request using presigned urls while running locally:

```bash
curl -X POST http://localhost:8080/api/v1/validate \
-H "Content-Type: application/json" \
-d '{
    "source_backup_url": "https://XXXXXXXXXXXXXXXXXXXXXXX.s3.us-west-2.amazonaws.com/foobar/backups_with_expiration_enabled/hashicorp-vault-backups/2025-10-19/vault_2025-10-19T00%3A15%3A06.snap?response-content-disposition=inline&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Security-Token=IQoJb3JpZ2luX2VjECMaCXVzLXdlc3QtMiJGMEQCIAnJW3ERFB7%2F1eVUEe48Anu%2BJQ6Fhm%2FQyxPi1WLz51OwAiBvVOU8KcXLkD2IMU5Bwa%2FO%2BppYam5TmBxMOn6Rg4SXLyraBAjM%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAEaDDIyMjk5Mzc3MDU1NCIMs5Udgcwf1OZhXv99Kq4EhBZux3HuF3syyZU65AU4%2BrnvFHPgBydbdCgOBANy41h8CA%2B4O9h6Gw42LxFDKR66l1d4vrfk%2B39J2X%2BYmpGtr0nsuxsh4pwzblg4e7V4EDjv8g%2BfR%2FEuLsX6Hkb5mrAeNiExUxRwVS63i3vadLlqRfid73NbJHAJPVZUCdWNTR9XDdINHB3VKOqZHspFav4vxdoU4HYh3f%2BO68z3Zhpj9iFYX%2Bdb%2BMry5J1UfUm6C0JPR1h%2BNyoaeXm1lBHNZ3XWVUjfaLgp3sei03L7jlK8DODs%2BlmCf06SutVwCdFEqEiwPjtcR86hue6HElvbt8Oqmk3rAGI3sZNkA%2Fb5XbZ7TkinzRTFNwXKFEsApzHy8Zd1u9D%2B26gF5xlt7txdcKi6Dd%2FZKXnpVd4fVuhujArJTw2VGu%2F9KWNTkercW2BDyh963LE18PGjAQZSV6pIIPAo7DA0dVs%2BhJT1XLuB8%2FJDe4rceVOo1v5IkUI48BSN3bAqE0QUEyRWMOgo9MySNNlckQSgcK4481A1KgF9WTq7N3dXQIOzct0aSY31s3N1XUhpyNMWnGhFxQlEKbXTsLZz%2BcNKnUigppOctheiP1kw27M7RImRlyKHEvy%2BzVboi3RtxxXoZsqL4ZUxV1uM5f53jbf9aIVUM4eJSaJxefQiVBBTbP7OMZ9Mt%2BkIpze6p4vy%2FPBRUApdbA5vG4CWDIT7SMDNr2n%2FYpNho35NtrmO3SuEnSjc4nxS8UkIFJdNMKqs0ccGOsYCcskPTQUihgbWYriNQKXo7TUGaOpYtPz2S1vy4ss3nFnGZeT5ph8mp9x6XKJIYRhCniulkg8YuMP8SP29gWRLrA26qy39yXUTbmBqJqxC%2Bi1UCWPIIsvcdHqtoTXySNusNp%2FgxPjCm%2FAJjGfj40HzDcTjsfT0hYbH5CDNGLrwm9zhamQsXwKHh6UiQ8k0W3KuGF0CR3BVaSKkEYhX%2BKi6STqZJmrddkLixsZU8Bg9NV62rq2%2FX5n4PZDbLR5%2BVF25UuSeABKNgo0gnA8yaiNE5ssIBGFIR2Ds%2BvACi3rPVPg6%2F5zkMMESVHd7EjwqagPognDXfSrflo66N4V7S3BWMOkH%2Bl%2FNis0oqq37jcf0AQZyTMpn5ptTQqVa8s%2FCcncG84AwRp5MuKywfB8wMFY3zot%2FD0LiGVtlELJiVBaP9TVYzZ4A058%3D&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIATH23W5A5PGJDBNHK%2F20251019%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20251019T031129Z&X-Amz-Expires=43200&X-Amz-SignedHeaders=host&X-Amz-Signature=eb38c6447af9687e5931121bdb71387923e060cd333522885a2b5a0582c47f8c",
    "source_keys_url": "https://XXXXXXXXXXXXXXXXXXXXXXX.s3.us-west-2.amazonaws.com/foobar/hashicorp-vault-init/vault_access.json?response-content-disposition=inline&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Security-Token=IQoJb3JpZ2luX2VjECMaCXVzLXdlc3QtMiJGMEQCIAnJW3ERFB7%2F1eVUEe48Anu%2BJQ6Fhm%2FQyxPi1WLz51OwAiBvVOU8KcXLkD2IMU5Bwa%2FO%2BppYam5TmBxMOn6Rg4SXLyraBAjM%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAEaDDIyMjk5Mzc3MDU1NCIMs5Udgcwf1OZhXv99Kq4EhBZux3HuF3syyZU65AU4%2BrnvFHPgBydbdCgOBANy41h8CA%2B4O9h6Gw42LxFDKR66l1d4vrfk%2B39J2X%2BYmpGtr0nsuxsh4pwzblg4e7V4EDjv8g%2BfR%2FEuLsX6Hkb5mrAeNiExUxRwVS63i3vadLlqRfid73NbJHAJPVZUCdWNTR9XDdINHB3VKOqZHspFav4vxdoU4HYh3f%2BO68z3Zhpj9iFYX%2Bdb%2BMry5J1UfUm6C0JPR1h%2BNyoaeXm1lBHNZ3XWVUjfaLgp3sei03L7jlK8DODs%2BlmCf06SutVwCdFEqEiwPjtcR86hue6HElvbt8Oqmk3rAGI3sZNkA%2Fb5XbZ7TkinzRTFNwXKFEsApzHy8Zd1u9D%2B26gF5xlt7txdcKi6Dd%2FZKXnpVd4fVuhujArJTw2VGu%2F9KWNTkercW2BDyh963LE18PGjAQZSV6pIIPAo7DA0dVs%2BhJT1XLuB8%2FJDe4rceVOo1v5IkUI48BSN3bAqE0QUEyRWMOgo9MySNNlckQSgcK4481A1KgF9WTq7N3dXQIOzct0aSY31s3N1XUhpyNMWnGhFxQlEKbXTsLZz%2BcNKnUigppOctheiP1kw27M7RImRlyKHEvy%2BzVboi3RtxxXoZsqL4ZUxV1uM5f53jbf9aIVUM4eJSaJxefQiVBBTbP7OMZ9Mt%2BkIpze6p4vy%2FPBRUApdbA5vG4CWDIT7SMDNr2n%2FYpNho35NtrmO3SuEnSjc4nxS8UkIFJdNMKqs0ccGOsYCcskPTQUihgbWYriNQKXo7TUGaOpYtPz2S1vy4ss3nFnGZeT5ph8mp9x6XKJIYRhCniulkg8YuMP8SP29gWRLrA26qy39yXUTbmBqJqxC%2Bi1UCWPIIsvcdHqtoTXySNusNp%2FgxPjCm%2FAJjGfj40HzDcTjsfT0hYbH5CDNGLrwm9zhamQsXwKHh6UiQ8k0W3KuGF0CR3BVaSKkEYhX%2BKi6STqZJmrddkLixsZU8Bg9NV62rq2%2FX5n4PZDbLR5%2BVF25UuSeABKNgo0gnA8yaiNE5ssIBGFIR2Ds%2BvACi3rPVPg6%2F5zkMMESVHd7EjwqagPognDXfSrflo66N4V7S3BWMOkH%2Bl%2FNis0oqq37jcf0AQZyTMpn5ptTQqVa8s%2FCcncG84AwRp5MuKywfB8wMFY3zot%2FD0LiGVtlELJiVBaP9TVYzZ4A058%3D&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ASIATH23W5A5PGJDBNHK%2F20251019%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20251019T031200Z&X-Amz-Expires=43200&X-Amz-SignedHeaders=host&X-Amz-Signature=4e1fb4e6cb0ac658bef44f170a3224622f86d5d079c3c1ef390308d4d4275d21",
    "path_values_map":{
        "secret/login-service": {
            "password":"supersecretpassword"
        }
    },
    "vault_version":"2.4.1"
}'
```

- To run the app, as usual create a fresh cloud development environment and then install golang:
```bash
wget https://go.dev/dl/go1.25.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.25.3.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.zshrc
source ~/.zshrc
go version
go install github.com/go-delve/delve/cmd/dlv@latest
```
- Install this extension for vscode: https://marketplace.visualstudio.com/items?itemName=golang.Go
- To start the debugger you can create a launch.json configuration using vscode or place this file in: `/workspaces/glueops/.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
    
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}"
        }
    ]
}
```
- After creating that `launch.json`, go to `main.go` and just set a break point and kick off the debugger.

**NOTES:** 
- When running the app locally you may run into permission issues and might need to update the scripts/commands under vault-backup-validator/vault/scripts/* to use `sudo`.
- There are a number of references to `vault` instead of `openbao` as at one point this repository was for Hashicorp Vault but are now using `Openbao` instead. So it's expected you will continue to see both references throughout our repositories.

