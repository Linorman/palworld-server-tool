package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/config"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/task"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
)

// getServerInfo godoc
//
//	@Summary		Get Server Info by ID
//	@Description	Get Server Info by ID
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{object}	ServerInfo
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/info [get]
func getServerInfo(c *gin.Context) {
	serverId := c.Param("server_id")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	info, err := tool.InfoWithConfig(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &ServerInfo{info["version"], info["name"]})
}

// getServerMetricsById godoc
//
//	@Summary		Get Server Metrics by ID
//	@Description	Get Server Metrics by ID
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{object}	ServerMetrics
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/metrics [get]
func getServerMetricsById(c *gin.Context) {
	serverId := c.Param("server_id")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	metrics, err := tool.MetricsWithConfig(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &ServerMetrics{
		ServerFps:        metrics["server_fps"].(int),
		CurrentPlayerNum: metrics["current_player_num"].(int),
		ServerFrameTime:  metrics["server_frame_time"].(float64),
		MaxPlayerNum:     metrics["max_player_num"].(int),
		Uptime:           metrics["uptime"].(int),
		Days:             metrics["days"].(int),
	})
}

// listPlayersByServer godoc
//
//	@Summary		List Players by Server
//	@Description	List Players by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{array}		database.TersePlayer
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/players [get]
func listPlayersByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	if _, exists := config.GetServer(serverId); !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	db := database.GetDB()
	players, err := service.ListPlayersByServer(db, serverId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}

// getPlayerByServer godoc
//
//	@Summary		Get Player by Server
//	@Description	Get Player by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			player_uid	path		string	true	"Player UID"
//	@Success		200			{object}	database.Player
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/players/{player_uid} [get]
func getPlayerByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	playerUid := c.Param("player_uid")
	if _, exists := config.GetServer(serverId); !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	db := database.GetDB()
	player, err := service.GetPlayerByServer(db, serverId, playerUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

// listOnlinePlayersByServer godoc
//
//	@Summary		List Online Players by Server
//	@Description	List Online Players by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{array}		database.OnlinePlayer
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/online_players [get]
func listOnlinePlayersByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	onlinePlayers, err := tool.ShowPlayersWithConfig(server)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, onlinePlayers)
}

// listGuildsByServer godoc
//
//	@Summary		List Guilds by Server
//	@Description	List Guilds by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{array}		database.Guild
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/guilds [get]
func listGuildsByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	if _, exists := config.GetServer(serverId); !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	db := database.GetDB()
	guilds, err := service.ListGuildsByServer(db, serverId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, guilds)
}

// getGuildByServer godoc
//
//	@Summary		Get Guild by Server
//	@Description	Get Guild by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Param			server_id			path		string	true	"Server ID"
//	@Param			admin_player_uid	path		string	true	"Admin Player UID"
//	@Success		200					{object}	database.Guild
//	@Failure		400					{object}	ErrorResponse
//	@Failure		404					{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/guilds/{admin_player_uid} [get]
func getGuildByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	adminPlayerUid := c.Param("admin_player_uid")
	if _, exists := config.GetServer(serverId); !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	db := database.GetDB()
	guild, err := service.GetGuildByServer(db, serverId, adminPlayerUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, guild)
}

// publishBroadcastByServer godoc
//
//	@Summary		Publish Broadcast by Server
//	@Description	Publish Broadcast by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string				true	"Server ID"
//	@Param			broadcast	body		BroadcastRequest	true	"Broadcast"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/broadcast [post]
func publishBroadcastByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMessage(req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := tool.BroadcastWithConfig(server, req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// shutdownServerByServer godoc
//
//	@Summary		Shutdown Server by Server
//	@Description	Shutdown Server by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string				true	"Server ID"
//	@Param			shutdown	body		ShutdownRequest		true	"Shutdown"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/shutdown [post]
func shutdownServerByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var req ShutdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateMessage(req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Seconds == 0 {
		req.Seconds = 60
	}
	if err := tool.ShutdownWithConfig(server, req.Seconds, req.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// kickPlayerByServer godoc
//
//	@Summary		Kick Player by Server
//	@Description	Kick Player by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			player_uid	path		string	true	"Player UID"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/players/{player_uid}/kick [post]
func kickPlayerByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	playerUid := c.Param("player_uid")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	db := database.GetDB()
	player, err := service.GetPlayerByServer(db, serverId, playerUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if player.SteamId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SteamId is empty"})
		return
	}
	if err := tool.KickPlayerWithConfig(server, "steam_"+player.SteamId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// banPlayerByServer godoc
//
//	@Summary		Ban Player by Server
//	@Description	Ban Player by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			player_uid	path		string	true	"Player UID"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/players/{player_uid}/ban [post]
func banPlayerByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	playerUid := c.Param("player_uid")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	db := database.GetDB()
	player, err := service.GetPlayerByServer(db, serverId, playerUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if player.SteamId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SteamId is empty"})
		return
	}
	if err := tool.BanPlayerWithConfig(server, "steam_"+player.SteamId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// unbanPlayerByServer godoc
//
//	@Summary		Unban Player by Server
//	@Description	Unban Player by Server
//	@Tags			Multi-Server
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			player_uid	path		string	true	"Player UID"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/players/{player_uid}/unban [post]
func unbanPlayerByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	playerUid := c.Param("player_uid")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	db := database.GetDB()
	player, err := service.GetPlayerByServer(db, serverId, playerUid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if player.SteamId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SteamId is empty"})
		return
	}
	if err := tool.UnBanPlayerWithConfig(server, "steam_"+player.SteamId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// putPlayersByServer godoc
//
//	@Summary		Put Players By Server
//	@Description	Put Players By Server Only For SavSync
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string				true	"Server ID"
//	@Param			players		body		[]database.Player	true	"Players"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/players [put]
func putPlayersByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var players []database.Player
	if err := c.ShouldBindJSON(&players); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.PutPlayersByServer(database.GetDB(), serverId, players); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// putGuildsByServer godoc
//
//	@Summary		Put Guilds By Server
//	@Description	Put Guilds By Server Only For SavSync
//	@Tags			Guild
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string				true	"Server ID"
//	@Param			guilds		body		[]database.Guild	true	"Guilds"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/guilds [put]
func putGuildsByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var guilds []database.Guild
	if err := c.ShouldBindJSON(&guilds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.PutGuildsByServer(database.GetDB(), serverId, guilds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// syncDataByServer godoc
//
//	@Summary		Sync Data By Server
//	@Description	Sync Data By Server
//	@Tags			Sync
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			from		query		From	true	"from"	enum(rest,sav)
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/sync [post]
func syncDataByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	from := c.Query("from")
	if from == "rest" {
		go task.PlayerSyncByServer(database.GetDB(), serverId)
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	} else if from == "sav" {
		go task.SavSyncByServer(serverId)
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from parameter"})
}

// listWhiteByServer godoc
//
//	@Summary		List White List By Server
//	@Description	List White List By Server
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{object}	[]database.PlayerW
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/whitelist [get]
func listWhiteByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	players, err := service.ListWhitelistByServer(database.GetDB(), serverId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}

// addWhiteByServer godoc
//
//	@Summary		Add White List By Server
//	@Description	Add White List By Server
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string				true	"Server ID"
//	@Param			player		body		database.PlayerW	true	"Player"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/whitelist [post]
func addWhiteByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var player database.PlayerW
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.AddWhitelistByServer(database.GetDB(), serverId, player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// removeWhiteByServer godoc
//
//	@Summary		Remove White List By Server
//	@Description	Remove White List By Server
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string				true	"Server ID"
//	@Param			player		body		database.PlayerW	true	"Player"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/whitelist [delete]
func removeWhiteByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var player database.PlayerW
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用 SteamId 或 PlayerId 作为标识符
	identifier := player.SteamID
	if identifier == "" {
		identifier = player.PlayerUID
	}

	if err := service.RemoveWhitelistByServer(database.GetDB(), serverId, identifier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// putWhiteByServer godoc
//
//	@Summary		Put White List By Server
//	@Description	Put White List By Server
//	@Tags			Player
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string					true	"Server ID"
//	@Param			players		body		[]database.PlayerW		true	"Players"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/whitelist [put]
func putWhiteByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var players []database.PlayerW
	if err := c.ShouldBindJSON(&players); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.PutWhitelistByServer(database.GetDB(), serverId, players); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// listRconCommandByServer godoc
//
//	@Summary		List Rcon Commands By Server
//	@Description	List Rcon Commands By Server
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{object}	[]database.RconCommandList
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/rcon [get]
func listRconCommandByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	rcons, err := service.ListRconCommandsByServer(database.GetDB(), serverId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rcons)
}

// addRconCommandByServer godoc
//
//	@Summary		Add Rcon Command By Server
//	@Description	Add Rcon Command By Server
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string						true	"Server ID"
//	@Param			command		body		database.RconCommandList	true	"Rcon Command"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/rcon [post]
func addRconCommandByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var rcon database.RconCommandList
	if err := c.ShouldBindJSON(&rcon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.AddRconCommandByServer(database.GetDB(), serverId, rcon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// importRconCommandsByServer godoc
//
//	@Summary		Import Rcon Commands By Server
//	@Description	Import Rcon Commands By Server from a TXT file
//	@Tags			Rcon
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			file		formData	file	true	"Upload txt file"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/rcon/import [post]
func importRconCommandsByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
			return
		}
		placeholder := ""
		if len(parts) >= 3 {
			placeholder = parts[2]
		}
		//rconCommand := database.RconCommandList{
		//	Command:     parts[0],
		//	Remark:      parts[1],
		//	Placeholder: placeholder,
		//}
		rconCommand := database.RconCommandList{
			ServerId: serverId,
			UUID:     "",
			RconCommand: database.RconCommand{
				Command:     parts[0],
				Remark:      parts[1],
				Placeholder: placeholder,
			},
		}
		if err := service.AddRconCommandByServer(database.GetDB(), serverId, rconCommand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	if err := scanner.Err(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// sendRconCommandByServer godoc
//
//	@Summary		Send Rcon Command By Server
//	@Description	Send Rcon Command By Server
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string					true	"Server ID"
//	@Param			command		body		SendRconCommandRequest	true	"Rcon Command"
//	@Success		200			{object}	MessageResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/rcon/send [post]
func sendRconCommandByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var req SendRconCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据UUID从数据库中查找命令
	rcons, err := service.ListRconCommandsByServer(database.GetDB(), serverId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var targetCommand string
	for _, rcon := range rcons {
		if rcon.UUID == req.UUID {
			targetCommand = rcon.Command
			break
		}
	}

	if targetCommand == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rcon command not found"})
		return
	}

	execCommand := fmt.Sprintf("%s %s", targetCommand, req.Content)
	response, err := tool.CustomCommandWithConfig(server, execCommand)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": response})
}

// putRconCommandByServer godoc
//
//	@Summary		Put Rcon Command By Server
//	@Description	Put Rcon Command By Server
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string						true	"Server ID"
//	@Param			uuid		path		string						true	"UUID"
//	@Param			command		body		database.RconCommandList	true	"Rcon Command"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/rcon/{uuid} [put]
func putRconCommandByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	uuid := c.Param("uuid")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var rcon database.RconCommandList
	if err := c.ShouldBindJSON(&rcon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.PutRconCommandByServer(database.GetDB(), serverId, uuid, rcon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// removeRconCommandByServer godoc
//
//	@Summary		Remove Rcon Command By Server
//	@Description	Remove Rcon Command By Server
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			uuid		path		string	true	"UUID"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/rcon/{uuid} [delete]
func removeRconCommandByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	uuid := c.Param("uuid")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	err := service.RemoveRconCommandByServer(database.GetDB(), serverId, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// listBackupsByServer godoc
//
//	@Summary		List backups By Server within a specified time range
//	@Description	List all backups By Server or backups within a specific time range.
//	@Tags			backup
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			startTime	query		int		false	"Start time of the backup range in timestamp"
//	@Param			endTime		query		int		false	"End time of the backup range in timestamp"
//	@Success		200			{array}		database.Backup
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/backups [get]
func listBackupsByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	var startTimestamp, endTimestamp int64
	var startTime, endTime time.Time
	var err error

	startTimeStr, endTimeStr := c.Query("startTime"), c.Query("endTime")

	if startTimeStr != "" {
		startTimestamp, err = strconv.ParseInt(startTimeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time"})
			return
		}
		startTime = time.Unix(0, startTimestamp*int64(time.Millisecond))
	}

	if endTimeStr != "" {
		endTimestamp, err = strconv.ParseInt(endTimeStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time"})
			return
		}
		endTime = time.Unix(0, endTimestamp*int64(time.Millisecond))
	}

	backups, err := service.ListBackupsByServerWithTimeRange(database.GetDB(), serverId, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, backups)
}

// downloadBackupByServer godoc
//
//	@Summary		Download Backup By Server
//	@Description	Download a backup By Server
//	@Tags			backup
//	@Accept			json
//	@Produce		application/octet-stream
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			backup_id	path		string	true	"Backup ID"
//	@Success		200			{file}		"Backupfile"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Failure		500			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/backups/{backup_id} [get]
func downloadBackupByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	backupId := c.Param("backup_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	backup, err := service.GetBackupByServer(database.GetDB(), serverId, backupId)
	if err != nil {
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	backupDir, err := tool.GetBackupDir()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", backup.Path))
	c.File(filepath.Join(backupDir, backup.Path))
}

// deleteBackupByServer godoc
//
//	@Summary		Delete Backup By Server
//	@Description	Delete a backup By Server
//	@Tags			backup
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Param			backup_id	path		string	true	"Backup ID"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id}/backups/{backup_id} [delete]
func deleteBackupByServer(c *gin.Context) {
	serverId := c.Param("server_id")
	backupId := c.Param("backup_id")
	_, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	backup, err := service.GetBackupByServer(database.GetDB(), serverId, backupId)
	if err != nil {
		if err == service.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.DeleteBackupByServer(database.GetDB(), serverId, backupId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	backupDir, err := tool.GetBackupDir()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = os.Remove(filepath.Join(backupDir, backup.Path))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
