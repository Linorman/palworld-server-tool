package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zaigie/palworld-server-tool/internal/config"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"net/http"
)

type ServerListResponse struct {
	Servers []ServerStatus `json:"servers"`
}

type ServerStatus struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Status      string `json:"status"`
	OnlineCount int    `json:"online_count"`
	MaxPlayers  int    `json:"max_players"`
}

type ServerCreateRequest struct {
	Id          string                 `json:"id" binding:"required"`
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Config      map[string]interface{} `json:"config" binding:"required"`
}

type ServerUpdateRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     *bool                  `json:"enabled"`
	Config      map[string]interface{} `json:"config"`
}

// listServers godoc
//
//	@Summary		List all servers
//	@Description	Get list of all configured servers with their status
//	@Tags			Server Management
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ServerListResponse
//	@Router			/api/servers [get]
func listServers(c *gin.Context) {
	servers := config.GetConfig().Servers
	var serverStatuses []ServerStatus

	for _, server := range servers {
		status := "offline"
		onlineCount := 0
		maxPlayers := 0

		if server.Enabled {
			// Try to get server status
			if _, err := tool.InfoWithConfig(&server); err == nil {
				status = "online"
				if metrics, err := tool.MetricsWithConfig(&server); err == nil {
					if currentPlayers, ok := metrics["current_player_num"].(int); ok {
						onlineCount = currentPlayers
					}
					if maxPlayerNum, ok := metrics["max_player_num"].(int); ok {
						maxPlayers = maxPlayerNum
					}
				}
			}
		}

		serverStatuses = append(serverStatuses, ServerStatus{
			Id:          server.Id,
			Name:        server.Name,
			Description: server.Description,
			Enabled:     server.Enabled,
			Status:      status,
			OnlineCount: onlineCount,
			MaxPlayers:  maxPlayers,
		})
	}

	c.JSON(http.StatusOK, ServerListResponse{Servers: serverStatuses})
}

// getServer godoc
//
//	@Summary		Get server details
//	@Description	Get detailed information about a specific server
//	@Tags			Server Management
//	@Accept			json
//	@Produce		json
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{object}	ServerStatus
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id} [get]
func getServerDetails(c *gin.Context) {
	serverId := c.Param("server_id")
	server, exists := config.GetServer(serverId)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	status := "offline"
	onlineCount := 0
	maxPlayers := 0

	if server.Enabled {
		if _, err := tool.InfoWithConfig(server); err == nil {
			status = "online"
			if metrics, err := tool.MetricsWithConfig(server); err == nil {
				if currentPlayers, ok := metrics["current_player_num"].(int); ok {
					onlineCount = currentPlayers
				}
				if maxPlayerNum, ok := metrics["max_player_num"].(int); ok {
					maxPlayers = maxPlayerNum
				}
			}
		}
	}

	c.JSON(http.StatusOK, ServerStatus{
		Id:          server.Id,
		Name:        server.Name,
		Description: server.Description,
		Enabled:     server.Enabled,
		Status:      status,
		OnlineCount: onlineCount,
		MaxPlayers:  maxPlayers,
	})
}

// createServer godoc
//
//	@Summary		Create a new server
//	@Description	Add a new server configuration
//	@Tags			Server Management
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server	body		ServerCreateRequest	true	"Server configuration"
//	@Success		201		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		409		{object}	ErrorResponse
//	@Router			/api/servers [post]
func createServer(c *gin.Context) {
	var req ServerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if server ID already exists
	if _, exists := config.GetServer(req.Id); exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Server ID already exists"})
		return
	}

	// Create new server configuration
	newServer := config.Server{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
	}

	// Parse config into server structure
	if rconConfig, ok := req.Config["rcon"].(map[string]interface{}); ok {
		if addr, ok := rconConfig["address"].(string); ok {
			newServer.Rcon.Address = addr
		}
		if pass, ok := rconConfig["password"].(string); ok {
			newServer.Rcon.Password = pass
		}
		if timeout, ok := rconConfig["timeout"].(float64); ok {
			newServer.Rcon.Timeout = int(timeout)
		}
		if useBase64, ok := rconConfig["use_base64"].(bool); ok {
			newServer.Rcon.UseBase64 = useBase64
		}
	}

	if restConfig, ok := req.Config["rest"].(map[string]interface{}); ok {
		if addr, ok := restConfig["address"].(string); ok {
			newServer.Rest.Address = addr
		}
		if user, ok := restConfig["username"].(string); ok {
			newServer.Rest.Username = user
		}
		if pass, ok := restConfig["password"].(string); ok {
			newServer.Rest.Password = pass
		}
		if timeout, ok := restConfig["timeout"].(float64); ok {
			newServer.Rest.Timeout = int(timeout)
		}
	}

	if saveConfig, ok := req.Config["save"].(map[string]interface{}); ok {
		if path, ok := saveConfig["path"].(string); ok {
			newServer.Save.Path = path
		}
		if decodePath, ok := saveConfig["decode_path"].(string); ok {
			newServer.Save.DecodePath = decodePath
		}
		if syncInterval, ok := saveConfig["sync_interval"].(float64); ok {
			newServer.Save.SyncInterval = int(syncInterval)
		}
		if backupInterval, ok := saveConfig["backup_interval"].(float64); ok {
			newServer.Save.BackupInterval = int(backupInterval)
		}
		if backupKeepDays, ok := saveConfig["backup_keep_days"].(float64); ok {
			newServer.Save.BackupKeepDays = int(backupKeepDays)
		}
	}

	// Add to configuration
	globalConfig := config.GetConfig()
	globalConfig.Servers = append(globalConfig.Servers, newServer)

	c.JSON(http.StatusCreated, gin.H{"success": true})
}

// updateServer godoc
//
//	@Summary		Update server configuration
//	@Description	Update an existing server configuration
//	@Tags			Server Management
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string				true	"Server ID"
//	@Param			server		body		ServerUpdateRequest	true	"Server configuration"
//	@Success		200			{object}	SuccessResponse
//	@Failure		400			{object}	ErrorResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id} [put]
func updateServer(c *gin.Context) {
	serverId := c.Param("server_id")
	var req ServerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	globalConfig := config.GetConfig()
	serverIndex := -1
	for i, server := range globalConfig.Servers {
		if server.Id == serverId {
			serverIndex = i
			break
		}
	}

	if serverIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	// Update server configuration
	server := &globalConfig.Servers[serverIndex]
	if req.Name != "" {
		server.Name = req.Name
	}
	if req.Description != "" {
		server.Description = req.Description
	}
	if req.Enabled != nil {
		server.Enabled = *req.Enabled
	}

	// Update config if provided
	if req.Config != nil {
		if rconConfig, ok := req.Config["rcon"].(map[string]interface{}); ok {
			if addr, ok := rconConfig["address"].(string); ok {
				server.Rcon.Address = addr
			}
			if pass, ok := rconConfig["password"].(string); ok {
				server.Rcon.Password = pass
			}
			if timeout, ok := rconConfig["timeout"].(float64); ok {
				server.Rcon.Timeout = int(timeout)
			}
			if useBase64, ok := rconConfig["use_base64"].(bool); ok {
				server.Rcon.UseBase64 = useBase64
			}
		}

		if restConfig, ok := req.Config["rest"].(map[string]interface{}); ok {
			if addr, ok := restConfig["address"].(string); ok {
				server.Rest.Address = addr
			}
			if user, ok := restConfig["username"].(string); ok {
				server.Rest.Username = user
			}
			if pass, ok := restConfig["password"].(string); ok {
				server.Rest.Password = pass
			}
			if timeout, ok := restConfig["timeout"].(float64); ok {
				server.Rest.Timeout = int(timeout)
			}
		}

		if saveConfig, ok := req.Config["save"].(map[string]interface{}); ok {
			if path, ok := saveConfig["path"].(string); ok {
				server.Save.Path = path
			}
			if decodePath, ok := saveConfig["decode_path"].(string); ok {
				server.Save.DecodePath = decodePath
			}
			if syncInterval, ok := saveConfig["sync_interval"].(float64); ok {
				server.Save.SyncInterval = int(syncInterval)
			}
			if backupInterval, ok := saveConfig["backup_interval"].(float64); ok {
				server.Save.BackupInterval = int(backupInterval)
			}
			if backupKeepDays, ok := saveConfig["backup_keep_days"].(float64); ok {
				server.Save.BackupKeepDays = int(backupKeepDays)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// deleteServer godoc
//
//	@Summary		Delete server
//	@Description	Remove a server configuration
//	@Tags			Server Management
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			server_id	path		string	true	"Server ID"
//	@Success		200			{object}	SuccessResponse
//	@Failure		401			{object}	ErrorResponse
//	@Failure		404			{object}	ErrorResponse
//	@Router			/api/servers/{server_id} [delete]
func deleteServer(c *gin.Context) {
	serverId := c.Param("server_id")

	globalConfig := config.GetConfig()
	serverIndex := -1
	for i, server := range globalConfig.Servers {
		if server.Id == serverId {
			serverIndex = i
			break
		}
	}

	if serverIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	// Remove server from configuration
	globalConfig.Servers = append(globalConfig.Servers[:serverIndex], globalConfig.Servers[serverIndex+1:]...)

	// TODO: Clean up associated data from database

	c.JSON(http.StatusOK, gin.H{"success": true})
}
