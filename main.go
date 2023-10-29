package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/glueops/vault-backup-validator/vault"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/api/v1/restore",restoreHandler)
	return r
}


func restoreHandler(c *gin.Context) {
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

	client, err := vault.NewVault(requestBody.DestinationVaultURL,requestBody.DestinationVaultToken)
	if(err != nil){
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
	}
	err = vault.RestoreSnapshotFromS3(client,requestBody)
	if(err!=nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
	}
	fmt.Println("Vault restore successful..Now proceeding to unseal..")
	
	vaultObj := vault.Vault{Client: client}
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
	if(verify_success){
		c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Invalid Backup, verification failed",
        })
        return
	}else{
		c.String(http.StatusOK, "Backup is Valid")
	}
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}