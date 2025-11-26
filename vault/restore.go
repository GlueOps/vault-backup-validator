package vault

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
	govault "github.com/hashicorp/vault/api"
	"github.com/glueops/vault-backup-validator/logger"
)


// normalizeWhitespace collapses all whitespace sequences to a single space
func normalizeWhitespace(s string) string {
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(s, " "))
}

// compareValues compares two values, normalizing JSON strings for semantic comparison
func compareValues(expected, actual interface{}) bool {
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(actual)
	
	// Direct match
	if string(expectedJSON) == string(actualJSON) {
		return true
	}
	
	// If both are strings, try parsing as JSON and compare normalized
	expectedStr, expectedIsStr := expected.(string)
	actualStr, actualIsStr := actual.(string)
	if expectedIsStr && actualIsStr {
		var expectedObj, actualObj interface{}
		if json.Unmarshal([]byte(expectedStr), &expectedObj) == nil &&
		   json.Unmarshal([]byte(actualStr), &actualObj) == nil {
			// Both are valid JSON strings, compare parsed/normalized
			e, _ := json.Marshal(expectedObj)
			a, _ := json.Marshal(actualObj)
			return string(e) == string(a)
		}
		
		// Not valid JSON - compare with normalized whitespace
		return normalizeWhitespace(expectedStr) == normalizeWhitespace(actualStr)
	}
	
	return false
}

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
	
	logger.Logger.Debug(fmt.Sprintf("Vault client address: %s", v.Address()))
	
	for path, values := range(restoreParams.PathValuesMap){
		originalPath := path
		parts := strings.Split(path, "/")
		path = strings.Join(parts[:1], "/") + "/data/" + strings.Join(parts[1:], "/")
		
		logger.Logger.Debug(fmt.Sprintf("Verifying path: %s (original: %s)", path, originalPath))
		
		for retry := 0; retry < maxRetries; retry++ {
			logger.Logger.Debug(fmt.Sprintf("Attempt %d/%d for path %s", retry+1, maxRetries, path))
			
			content, err := v.Logical().Read(path)
			if err == nil {
				if(content == nil){
					logger.Logger.Debug(fmt.Sprintf("Path %s returned nil content", path))
					logger.Logger.Error("Path does not exist or has no data")
					return false, fmt.Errorf("no values in the given path or the given path does not exist")
				}
				data := content.Data
				logger.Logger.Debug(fmt.Sprintf("Raw content.Data: %+v", data))
				data = data["data"].(map[string]interface{})
				logger.Logger.Debug(fmt.Sprintf("Extracted data: %+v", data))
				
				for key, value := range values.(map[string]interface{}) {
					expectedJSON, _ := json.Marshal(value)
					actualJSON, _ := json.Marshal(data[key])
					logger.Logger.Debug(fmt.Sprintf("Comparing key '%s':", key))
					logger.Logger.Debug(fmt.Sprintf("  Expected (type %T): %s", value, string(expectedJSON)))
					logger.Logger.Debug(fmt.Sprintf("  Actual   (type %T): %s", data[key], string(actualJSON)))
					
					if !compareValues(value, data[key]) {
						logger.Logger.Debug(fmt.Sprintf("Mismatch for key '%s' - expected: %s, got: %s", key, string(expectedJSON), string(actualJSON)))
						logger.Logger.Error("Value mismatch detected")
						return false, nil
					}
					logger.Logger.Debug(fmt.Sprintf("Key '%s' matches!", key))
				}
				break // Success, move to next path
			}else{
				logger.Logger.Debug(fmt.Sprintf("Error reading path %s: %s", path, err.Error()))
				logger.Logger.Info("Retrying to verify restore...")
				if retry < maxRetries-1 {
					time.Sleep(retryDelay)
				}
			}
		}
	}
	logger.Logger.Info("All paths verified successfully")
    return true, nil
}
