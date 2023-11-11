package main

import (
	"encoding/json"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/glueops/vault-backup-validator/logger"
	"github.com/glueops/vault-backup-validator/vault"
	"github.com/hashicorp/vault/api"
)


func InitiateVaultSetup(c *gin.Context) (*api.Client, error){

	// Start the vault process with configs
	_, err := vault.SetupVault()
	if(err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return nil, err
	}

	// Store the vault unseal key and token to use later
	initialSecret := make(map[string]string)
	jsonData, err := os.ReadFile("secrets.json")
    if err != nil {
        logger.Logger.Error("Error reading secrets.json file: "+err.Error())
        return nil, err
    }

	json.Unmarshal(jsonData, &initialSecret)
	initialVaultSecret := vault.VaultSecrets{
		Keys: []string{initialSecret["key"]},
		Token: initialSecret["token"],
	}

	vault_addr := "http://localhost:8200"
	client, err := vault.NewVault(vault_addr,initialVaultSecret.Token)
	if(err != nil){
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return nil, err
	}
	
	vaultObj := vault.Vault{Client: client}
	// Unseal Vault 
	_, err = vaultObj.Unseal(&initialVaultSecret)
	if(err != nil){
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return client, err
	}

	return client, nil

}

func VerifyBackup(client *api.Client, requestBody vault.RestoreParams, c *gin.Context) (bool, error){

	vaultObj := vault.Vault{Client: client}
	// Restore snapshot in vault
	err := vault.RestoreSnapshotFromS3(client,requestBody)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return false, err
	}
	logger.Logger.Info("Vault restore successful..Proceeding to unseal..")
	
	secrets, err := vaultObj.ParseSecrets(requestBody.SourceKeysURL)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return false, err
	}

	// Unseal vault again after restore, with the backup keys
	_ , err = vaultObj.Unseal(secrets)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return false, err
	}
	logger.Logger.Info("Vault is unsealed..proceeding to verify the restore")

	// Verify restored values in vault
	verify_success, err := vault.VerifyRestore(vaultObj.Client, secrets, requestBody)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return verify_success, err
	}
	return verify_success, nil
}

func validateHandler(c *gin.Context) {

	// Cleanup vault incase of execution errors or in general
	defer vault.CleanupVault()

	// Unmarshal the input to vault.RestoreParams
	var requestBody vault.RestoreParams
	if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid JSON data, "+err.Error(),
        })
        return
    }

	// Validate inuput params
	if err := vault.ValidateResotreParams(requestBody); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
	}

	client, err := InitiateVaultSetup(c)
	if(err != nil){
		return
	}
	verify_success, err := VerifyBackup(client, requestBody, c)
	if(err != nil){
		return
	}
	
	if(verify_success){
		logger.Logger.Info("Backup is verified successfully and it is valid")
		response := gin.H{
			"status":  "success",
			"message": "Backup is valid",
		}
		c.JSON(http.StatusOK, response)
	}else{
		logger.Logger.Info("Backup is verified and it is Invalid")
		response := gin.H{
			"status":  "error",
			"message": "Backup is invalid",
		}
		c.JSON(http.StatusBadRequest, response)
	}

}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/api/v1/validate",validateHandler)
	return r
}

func main() {

	logger.InitLogger()
	r := setupRouter()
	r.Run(":8080")
}