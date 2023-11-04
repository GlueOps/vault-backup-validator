package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"os/exec"

	govault "github.com/hashicorp/vault/api"
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
        return nil, err
    }
	client.SetToken(token)
	return client, nil
}

func SetupVault() ([]byte, error){
	scriptPath := "vault/setup_vault.sh"

    cmd := exec.Command("bash", scriptPath)
    out, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("error starting vault: %v\n check vault.log", err)
    }
	return out, nil
}

func CleanupVault() ([]byte, error){
	scriptPath := "vault/cleanup_vault.sh"

    cmd := exec.Command("bash", scriptPath)
    out, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("error cleaning up vault: %v", err)
    }
	return out, nil
}

func (v Vault) ParseSecrets(keys_url string) (*VaultSecrets,error){

	resp, err := http.Get(keys_url)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil,fmt.Errorf("http request failed with status: %s, verify if the url is still valid", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil,fmt.Errorf("error reading response body: %s", err)	
	}
	
	var data VaultSecrets
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil,fmt.Errorf("error unmarshaling json: %s", err)
	}
	
	return &data, nil
}

func (v Vault) Unseal(secrets *VaultSecrets) (*govault.SealStatusResponse, error){

	sys := v.Client.Sys()
	var res *govault.SealStatusResponse
	for _, key := range(secrets.Keys){
		res, err := sys.Unseal(key)
		if (err != nil){
			return res, err
		}
	}
	return res, nil
}