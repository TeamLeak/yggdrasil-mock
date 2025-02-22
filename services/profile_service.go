package services

import (
	"database/sql"
	"net/http"
	"yggdrasil/database"

	"github.com/gin-gonic/gin"
)

func ProfileHandler(c *gin.Context, db *sql.DB) {
	uuid := c.Param("uuid")
	//unsigned := c.Query("unsigned") == "false"

	character, err := database.GetCharacterByUUID(db, uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())

	texturesData := fmt.Sprintf(`{
		"timestamp": %s,
		"profileId": %s,
  		"profileName": %s,
		"textures": {
			"SKIN": {
				"url": "https://127.0.0.1:8080/textures/aeef16d337c7f0314683fea2464dfb23e9be4379",
				"metadata": {
					"model": "default"
				}
			},
			"CAPE": {
				"url": "https://127.0.0.1:8080/textures/e8ad33b989ccde06baefb6be595601bb14d97a85"
			}
		}
	}`, timestamp, character.UUID, character.Name)

	encodedTextures := base64.StdEncoding.EncodeToString([]byte(texturesData))

	response := gin.H{
		"id":   character.UUID,
		"name": character.Name,
		"properties": []gin.H{
			{
				"name":  "textures",
				"value": encodedTextures,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}
