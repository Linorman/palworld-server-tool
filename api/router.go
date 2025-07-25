package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zaigie/palworld-server-tool/internal/auth"
)

type SuccessResponse struct {
	Success bool `json:"success"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EmptyResponse struct{}

func ignoreLogPrefix(path string) bool {
	prefixes := []string{"/swagger/", "/assets/", "/favicon.ico", "/map"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if !ignoreLogPrefix(param.Path) {
			statusColor := param.StatusCodeColor()
			methodColor := param.MethodColor()
			resetColor := param.ResetColor()
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.ErrorMessage,
			)
		}
		return ""
	})
}

func RegisterRouter(r *gin.Engine) {
	r.Use(Logger(), gin.Recovery())

	r.POST("/api/login", loginHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := r.Group("/api")

	// Server Management APIs
	serverGroup := apiGroup.Group("/servers")
	{
		serverGroup.GET("", listServers)
		serverGroup.GET("/:server_id", getServerDetails)
	}

	// Authenticated Server Management APIs
	authServerGroup := apiGroup.Group("/servers")
	authServerGroup.Use(auth.JWTAuthMiddleware())
	{
		authServerGroup.POST("", createServer)
		authServerGroup.PUT("/:server_id", updateServer)
		authServerGroup.DELETE("/:server_id", deleteServer)
	}

	anonymousGroup := apiGroup.Group("")
	{
		// Legacy single server APIs (backward compatibility)
		anonymousGroup.GET("/server", getServer)
		anonymousGroup.GET("/server/tool", getServerTool)
		anonymousGroup.GET("/server/metrics", getServerMetrics)
		anonymousGroup.GET("/player", listPlayers)
		anonymousGroup.GET("/player/:player_uid", getPlayer)
		anonymousGroup.GET("/online_player", listOnlinePlayers)
		anonymousGroup.GET("/guild", listGuilds)
		anonymousGroup.GET("/guild/:admin_player_uid", getGuild)

		// Multi-server APIs with optional server_id parameter
		anonymousGroup.GET("/servers/:server_id/info", getServerInfo)
		anonymousGroup.GET("/servers/:server_id/metrics", getServerMetricsById)
		anonymousGroup.GET("/servers/:server_id/players", listPlayersByServer)
		anonymousGroup.GET("/servers/:server_id/players/:player_uid", getPlayerByServer)
		anonymousGroup.GET("/servers/:server_id/online_players", listOnlinePlayersByServer)
		anonymousGroup.GET("/servers/:server_id/guilds", listGuildsByServer)
		anonymousGroup.GET("/servers/:server_id/guilds/:admin_player_uid", getGuildByServer)
	}

	authGroup := apiGroup.Group("")
	authGroup.Use(auth.JWTAuthMiddleware())
	{
		// Legacy single server APIs (backward compatibility)
		authGroup.POST("/server/broadcast", publishBroadcast)
		authGroup.POST("/server/shutdown", shutdownServer)
		authGroup.PUT("/player", putPlayers)
		authGroup.POST("/player/:player_uid/kick", kickPlayer)
		authGroup.POST("/player/:player_uid/ban", banPlayer)
		authGroup.POST("/player/:player_uid/unban", unbanPlayer)
		authGroup.PUT("/guild", putGuilds)
		authGroup.POST("/sync", syncData)
		authGroup.GET("/whitelist", listWhite)
		authGroup.POST("/whitelist", addWhite)
		authGroup.DELETE("/whitelist", removeWhite)
		authGroup.PUT("/whitelist", putWhite)
		authGroup.GET("/rcon", listRconCommand)
		authGroup.POST("/rcon", addRconCommand)
		authGroup.POST("/rcon/import", importRconCommands)
		authGroup.POST("/rcon/send", sendRconCommand)
		authGroup.PUT("/rcon/:uuid", putRconCommand)
		authGroup.DELETE("/rcon/:uuid", removeRconCommand)
		authGroup.GET("/backup", listBackups)
		authGroup.GET("/backup/:backup_id", downloadBackup)
		authGroup.DELETE("/backup/:backup_id", deleteBackup)

		// Multi-server APIs with server_id parameter
		authGroup.POST("/servers/:server_id/broadcast", publishBroadcastByServer)
		authGroup.POST("/servers/:server_id/shutdown", shutdownServerByServer)
		authGroup.PUT("/servers/:server_id/players", putPlayersByServer)
		authGroup.POST("/servers/:server_id/players/:player_uid/kick", kickPlayerByServer)
		authGroup.POST("/servers/:server_id/players/:player_uid/ban", banPlayerByServer)
		authGroup.POST("/servers/:server_id/players/:player_uid/unban", unbanPlayerByServer)
		authGroup.PUT("/servers/:server_id/guilds", putGuildsByServer)
		authGroup.POST("/servers/:server_id/sync", syncDataByServer)
		authGroup.GET("/servers/:server_id/whitelist", listWhiteByServer)
		authGroup.POST("/servers/:server_id/whitelist", addWhiteByServer)
		authGroup.DELETE("/servers/:server_id/whitelist", removeWhiteByServer)
		authGroup.PUT("/servers/:server_id/whitelist", putWhiteByServer)
		authGroup.GET("/servers/:server_id/rcon", listRconCommandByServer)
		authGroup.POST("/servers/:server_id/rcon", addRconCommandByServer)
		authGroup.POST("/servers/:server_id/rcon/import", importRconCommandsByServer)
		authGroup.POST("/servers/:server_id/rcon/send", sendRconCommandByServer)
		authGroup.PUT("/servers/:server_id/rcon/:uuid", putRconCommandByServer)
		authGroup.DELETE("/servers/:server_id/rcon/:uuid", removeRconCommandByServer)
		authGroup.GET("/servers/:server_id/backups", listBackupsByServer)
		authGroup.GET("/servers/:server_id/backups/:backup_id", downloadBackupByServer)
		authGroup.DELETE("/servers/:server_id/backups/:backup_id", deleteBackupByServer)
	}
}
