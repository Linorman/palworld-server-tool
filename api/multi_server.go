package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/config"
	"github.com/zaigie/palworld-server-tool/internal/database"
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

// TODO: 需要添加更多的多服务器API函数，包括：
// - putPlayersByServer
// - putGuildsByServer
// - syncDataByServer
// - listWhiteByServer
// - addWhiteByServer
// - removeWhiteByServer
// - putWhiteByServer
// - listRconCommandByServer
// - addRconCommandByServer
// - importRconCommandsByServer
// - sendRconCommandByServer
// - putRconCommandByServer
// - removeRconCommandByServer
// - listBackupsByServer
// - downloadBackupByServer
// - deleteBackupByServer
