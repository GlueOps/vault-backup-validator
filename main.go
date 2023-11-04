package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/glueops/vault-backup-validator/vault"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/api/v1/validate",validateHandler)
	return r
}


func validateHandler(c *gin.Context) {

	var requestBody vault.RestoreParams

	if err := c.BindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid JSON data",
        })
        return
    }
	if err := vault.ValidateResotreParams(requestBody); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
	}

	_, err := vault.SetupVault()
	if(err != nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
	}

	initialSecret := make(map[string]string)
	jsonData, err := os.ReadFile("secrets.json")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
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
        return
	}
	
	vaultObj := vault.Vault{Client: client}
	vaultObj.Unseal(&initialVaultSecret)
	
	err = vault.RestoreSnapshotFromS3(client,requestBody)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
	}
	fmt.Println("Vault restore successful..Now proceeding to unseal..")
	
	
	secrets, err := vaultObj.ParseSecrets(requestBody.SourceKeysURL)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
	}
	_ , err = vaultObj.Unseal(secrets)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
	}
	fmt.Printf("vault is unsealed..moving to verify")
	verify_success, err := vault.VerifyRestore(vaultObj.Client,secrets)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
	}

	vault.CleanupVault()

	if(verify_success){
		c.String(http.StatusOK, "Backup is Invalid")
	}else{
		c.String(http.StatusOK, "Backup is Valid")
	}

}

func main() {
	r := setupRouter()
	r.Run(":8080")
}