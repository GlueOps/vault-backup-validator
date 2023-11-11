package vault

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
	"strings"
	govault "github.com/hashicorp/vault/api"
	"github.com/glueops/vault-backup-validator/logger"
)

type RestoreParams struct {
	SourceBackupURL             string                `json:"source_backup_url"`
	SourceKeysURL               string                `json:"source_keys_url"`
	PathValuesMap               map[string]interface{}`json:"path_values_map"`
	VaultVersion                string                `json:"vault_version"`
}

func RestoreSnapshotFromS3(v *govault.Client, p RestoreParams) error{

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	logger.Logger.Info("Downloading backup...")
	resp, err := http.Get(p.SourceBackupURL)
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}
	logger.Logger.Info("Finished downloading backup")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Error("HTTP request failed with status: "+resp.Status+ ", verify if the backup url is still valid")
		return fmt.Errorf("HTTP request failed with status: %s, verify if the backup url is still valid", resp.Status)
	}

	//call restore with resp.body and client with token
	logger.Logger.Info("starting to restore...")
	if err := v.Sys().RaftSnapshotRestoreWithContext(context.Background(),resp.Body,true); err != nil {
		logger.Logger.Error(err.Error())
		return err
	}
	logger.Logger.Info("Restoring backup success")
    return nil
}

func ValidateResotreParams(p RestoreParams) error{
    if p.SourceBackupURL == "" || p.SourceKeysURL == "" || p.PathValuesMap == nil || p.VaultVersion == "" {
		logger.Logger.Error("one or more input from the client are empty")
        return fmt.Errorf("one or more input from the client are empty")
    }
	logger.Logger.Info("Inputs have been validated successfully")
    return nil
}

func VerifyRestore(v *govault.Client, secrets *VaultSecrets, restoreParams RestoreParams) (bool, error){

	logger.Logger.Info("Starting to verify the restore..")
	maxRetries := 3
    retryDelay := 2 * time.Second
	v.SetToken(secrets.Token)
	for path, values := range(restoreParams.PathValuesMap){
		parts := strings.Split(path, "/")
		path = strings.Join(parts[:1], "/") + "/data/" + strings.Join(parts[1:], "/")
		for retry := 0; retry < maxRetries; retry++ {
			content, err := v.Logical().Read(path)
			if err == nil {
				data := content.Data
				data = data["data"].(map[string]interface{})
				for key, value := range values.(map[string]interface{}) {
					if (value.(string) != data[key]){
						return false, nil
					}
				}
			}else{
				logger.Logger.Error(err.Error())
				logger.Logger.Info("Retrying to verify restore...")
				if retry < maxRetries-1 {
					// Backoff before retry
					time.Sleep(retryDelay)
				}
			}
		}
	}
    return true, nil
}
