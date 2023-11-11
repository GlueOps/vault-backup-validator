package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"os/exec"

	govault "github.com/hashicorp/vault/api"
	"github.com/glueops/vault-backup-validator/logger"
)

type Vault struct{
	Client *govault.Client
}

type VaultSecrets struct{
	Keys []string `json:"keys"`
	Token string  `json:"root_token"`
}

func NewVault(url string,token string) (*govault.Client, error) {
	
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	client, err := govault.NewClient(&govault.Config{Address: url,HttpClient: httpClient})
	if err != nil {
		logger.Logger.Error(err.Error())
        return nil, err
    }
	client.SetToken(token)
	return client, nil
}

func InstallVault(vault_version string) ([]byte, error){

	logger.Logger.Info("Installing vault "+vault_version+"...")
	scriptPath := "vault/scripts/install_vault.sh"
    cmd := exec.Command("bash", scriptPath, vault_version)
    out, err := cmd.CombinedOutput()
    if err != nil {
		logger.Logger.Error(err.Error())
        return out, fmt.Errorf("error installing vault")
    }
	logger.Logger.Info("Vault installation done")
	return out, nil
}

func SetupVault() ([]byte, error){

	logger.Logger.Info("Setting up test vault server...")
	scriptPath := "vault/scripts/setup_vault.sh"

    cmd := exec.Command("bash", scriptPath)
    out, err := cmd.CombinedOutput()
    if err != nil {
		logger.Logger.Error(err.Error())
        return out, fmt.Errorf("error starting vault: %v\n check vault.log", err)
    }
	logger.Logger.Info("Vault setup done")
	return out, nil
}

func CleanupVault() ([]byte, error){

	logger.Logger.Info("Cleaning up vault")
	scriptPath := "vault/scripts/cleanup_vault.sh"

    cmd := exec.Command("bash", scriptPath)
    out, err := cmd.CombinedOutput()
    if err != nil {
		logger.Logger.Error(err.Error())
        return nil, fmt.Errorf("error cleaning up vault: %v", err)
    }
	logger.Logger.Info("Vault Cleanup done")
	return out, nil
}

func (v Vault) ParseSecrets(keys_url string) (*VaultSecrets,error){

	logger.Logger.Info("Parsing secrets from BackupKeys Endpoint...")
	resp, err := http.Get(keys_url)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil,err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil,fmt.Errorf("http request failed with status: %s, verify if the url is still valid", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil,fmt.Errorf("error reading response body: %s", err)	
	}
	
	var data VaultSecrets
	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil,fmt.Errorf("error unmarshaling json: %s", err)
	}
	
	return &data, nil
}

func (v Vault) Unseal(secrets *VaultSecrets) (*govault.SealStatusResponse, error){

	logger.Logger.Info("Unsealing vault...")
	sys := v.Client.Sys()
	var res *govault.SealStatusResponse
	for _, key := range(secrets.Keys){
		res, err := sys.Unseal(key)
		if (err != nil){
			logger.Logger.Error(err.Error())
			return res, err
		}
	}
	logger.Logger.Info("Vault unseal done")
	return res, nil
}