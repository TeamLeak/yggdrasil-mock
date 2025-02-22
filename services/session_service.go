package services

import (
	"database/sql"
	"net/http"
	"strconv"
	"yggdrasil/database"

	"github.com/gin-gonic/gin"
)

func JoinServerHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		AccessToken     string `json:"accessToken"`
		SelectedProfile string `json:"selectedProfile"`
		ServerID        string `json:"serverId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := database.GetToken(db, req.AccessToken)
	if err != nil || token == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	if strconv.Itoa(token.CharacterID) != req.SelectedProfile {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid profile"})
		return
	}

	c.Status(http.StatusNoContent)
}

func HasJoinedHandler(c *gin.Context, db *sql.DB) {
	serverID := c.Query("serverId")
	username := c.Query("username")
	if serverID == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	character, err := database.FindCharacterByServerAndName(db, username, serverID)
	if err != nil {
		c.Status(http.StatusNoContent)
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
