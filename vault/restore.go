package vault

import (
	"context"
	"fmt"
	"net/http"
	"crypto/tls"

	govault "github.com/hashicorp/vault/api"
)

type RestoreParams struct {
	SourceBackupURL             string `json:"source_backup_url"`
	SourceKeysURL               string `json:"source_keys_url"`
	SourceTokenURL              string `json:"source_token_url"`
    DestinationVaultURL         string `json:"destination_vault_url"`
    DestinationVaultToken       string `json:"destination_vault_token"`
}

func RestoreSnapshotFromS3(v *govault.Client, p RestoreParams) error{

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(p.SourceBackupURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status: %s, verify if the url is still valid", resp.Status)
	}

	//call restore with resp.body and client with token
	if err := v.Sys().RaftSnapshotRestoreWithContext(context.Background(),resp.Body,true); err != nil {
		return err
	}

    return nil
}

func ValidateResotreParams(p RestoreParams) error{
    if(p.SourceBackupURL == "" || p.SourceKeysURL == "" || p.SourceTokenURL == ""){
        return fmt.Errorf("one or more input from the client are empty")
    }
    return nil
}

func VerifyRestore(v *govault.Client, secrets *VaultSecrets) (bool, error){

	v.SetToken(secrets.Token)
	secret, err := v.Logical().Read("secret/key-1-for-balaji")
	if(err != nil){
		return false, err
	}
	data := secret.Data
	for key, value := range(data){
		if(key == "key1" && value == "value1"){
			return true, nil
		}
	}
	return false, nil
}
