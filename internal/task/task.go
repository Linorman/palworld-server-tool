package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/config"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/system"

	"github.com/go-co-op/gocron/v2"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"
)

var s gocron.Scheduler

func BackupTask(db *bbolt.DB) {
	logger.Info("Scheduling backup for all servers...\n")

	servers := config.GetEnabledServers()
	for _, server := range servers {
		logger.Infof("Backing up server %s (%s)...\n", server.Name, server.Id)

		path, err := tool.BackupWithConfig(&server)
		if err != nil {
			logger.Errorf("Backup failed for server %s: %v\n", server.Id, err)
			continue
		}

		err = service.AddBackupByServer(db, server.Id, database.Backup{
			ServerId: server.Id,
			BackupId: uuid.New().String(),
			Path:     path,
			SaveTime: time.Now(),
		})
		if err != nil {
			logger.Errorf("Failed to save backup record for server %s: %v\n", server.Id, err)
			continue
		}

		logger.Infof("Auto backup for server %s to %s\n", server.Id, path)

		keepDays := server.Save.BackupKeepDays
		if keepDays == 0 {
			keepDays = 7
		}
		err = tool.CleanOldBackupsByServer(db, server.Id, keepDays)
		if err != nil {
			logger.Errorf("Failed to clean old backups for server %s: %v\n", server.Id, err)
		}
	}
}

func PlayerSync(db *bbolt.DB) {
	logger.Info("Scheduling Player sync for all servers...\n")

	servers := config.GetEnabledServers()
	for _, server := range servers {
		logger.Infof("Syncing players for server %s (%s)...\n", server.Name, server.Id)

		onlinePlayers, err := tool.ShowPlayersWithConfig(&server)
		if err != nil {
			logger.Errorf("Failed to get online players for server %s: %v\n", server.Id, err)
			continue
		}

		err = service.PutPlayersOnlineByServer(db, server.Id, onlinePlayers)
		if err != nil {
			logger.Errorf("Failed to save online players for server %s: %v\n", server.Id, err)
			continue
		}

		logger.Infof("Player sync done for server %s\n", server.Id)

		playerLogging := viper.GetBool("task.player_logging")
		if playerLogging {
			go PlayerLoggingByServer(&server, onlinePlayers)
		}

		kickInterval := viper.GetBool("manage.kick_non_whitelist")
		if kickInterval {
			go CheckAndKickPlayersByServer(db, &server, onlinePlayers)
		}
	}
}

func isPlayerWhitelisted(player database.OnlinePlayer, whitelist []database.PlayerW) bool {
	for _, whitelistedPlayer := range whitelist {
		if (player.PlayerUid != "" && player.PlayerUid == whitelistedPlayer.PlayerUID) ||
			(player.SteamId != "" && player.SteamId == whitelistedPlayer.SteamID) {
			return true
		}
	}
	return false
}

// Server-specific player caches
var playerCaches = make(map[string]map[string]string)
var firstPolls = make(map[string]bool)

func PlayerLoggingByServer(server *config.Server, players []database.OnlinePlayer) {
	loginMsg := viper.GetString("task.player_login_message")
	logoutMsg := viper.GetString("task.player_logout_message")

	tmp := make(map[string]string, len(players))
	for _, player := range players {
		if player.PlayerUid != "" {
			tmp[player.PlayerUid] = player.Nickname
		}
	}

	playerCache, exists := playerCaches[server.Id]
	if !exists {
		playerCache = make(map[string]string)
		playerCaches[server.Id] = playerCache
	}

	firstPoll, exists := firstPolls[server.Id]
	if !exists {
		firstPoll = true
		firstPolls[server.Id] = true
	}

	if !firstPoll {
		for id, name := range tmp {
			if _, ok := playerCache[id]; !ok {
				BroadcastVariableMessageByServer(server, loginMsg, name, len(players))
			}
		}
		for id, name := range playerCache {
			if _, ok := tmp[id]; !ok {
				BroadcastVariableMessageByServer(server, logoutMsg, name, len(players))
			}
		}
	}
	firstPolls[server.Id] = false
	playerCaches[server.Id] = tmp
}

func BroadcastVariableMessageByServer(server *config.Server, message string, username string, onlineNum int) {
	message = strings.ReplaceAll(message, "{username}", username)
	message = strings.ReplaceAll(message, "{online_num}", strconv.Itoa(onlineNum))
	message = strings.ReplaceAll(message, "{server_name}", server.Name)
	arr := strings.Split(message, "\n")
	for _, msg := range arr {
		err := tool.BroadcastWithConfig(server, msg)
		if err != nil {
			logger.Warnf("Broadcast fail for server %s, %s \n", server.Id, err)
		}
		// 连续发送不知道为啥行会错乱, 只能加点延迟
		time.Sleep(1000 * time.Millisecond)
	}
}

// Legacy function for backward compatibility
func BroadcastVariableMessage(message string, username string, onlineNum int) {
	message = strings.ReplaceAll(message, "{username}", username)
	message = strings.ReplaceAll(message, "{online_num}", strconv.Itoa(onlineNum))
	arr := strings.Split(message, "\n")
	for _, msg := range arr {
		err := tool.Broadcast(msg)
		if err != nil {
			logger.Warnf("Broadcast fail, %s \n", err)
		}
		// 连续发送不知道为啥行会错乱, 只能加点延迟
		time.Sleep(1000 * time.Millisecond)
	}
}

func CheckAndKickPlayersByServer(db *bbolt.DB, server *config.Server, players []database.OnlinePlayer) {
	whitelist, err := service.ListWhitelistByServer(db, server.Id)
	if err != nil {
		logger.Errorf("Failed to get whitelist for server %s: %v\n", server.Id, err)
		return
	}

	for _, player := range players {
		if !isPlayerWhitelisted(player, whitelist) {
			identifier := player.SteamId
			if identifier == "" {
				logger.Warnf("Kicked %s fail on server %s, SteamId is empty \n", player.Nickname, server.Id)
				continue
			}
			err := tool.KickPlayerWithConfig(server, fmt.Sprintf("steam_%s", identifier))
			if err != nil {
				logger.Warnf("Kicked %s fail on server %s, %s \n", player.Nickname, server.Id, err)
				continue
			}
			logger.Warnf("Kicked %s successful on server %s \n", player.Nickname, server.Id)
		}
	}
	logger.Infof("Check whitelist done for server %s\n", server.Id)
}

func SavSync() {
	logger.Info("Scheduling Sav sync for all servers...\n")

	servers := config.GetEnabledServers()
	for _, server := range servers {
		if server.Save.Path == "" {
			logger.Warnf("Save path not configured for server %s, skipping\n", server.Id)
			continue
		}

		logger.Infof("Syncing save for server %s (%s)...\n", server.Name, server.Id)

		err := tool.DecodeWithConfig(&server, server.Save.Path)
		if err != nil {
			logger.Errorf("Failed to decode save for server %s: %v\n", server.Id, err)
			continue
		}

		logger.Infof("Sav sync done for server %s\n", server.Id)
	}
}

func Schedule(db *bbolt.DB) {
	s := getScheduler()

	playerSyncInterval := time.Duration(viper.GetInt("task.sync_interval"))
	savSyncInterval := time.Duration(viper.GetInt("save.sync_interval"))
	backupInterval := time.Duration(viper.GetInt("save.backup_interval"))

	if playerSyncInterval > 0 {
		go PlayerSync(db)
		_, err := s.NewJob(
			gocron.DurationJob(playerSyncInterval*time.Second),
			gocron.NewTask(PlayerSync, db),
		)
		if err != nil {
			logger.Errorf("%v\n", err)
		}
	}

	if savSyncInterval > 0 {
		go SavSync()
		_, err := s.NewJob(
			gocron.DurationJob(savSyncInterval*time.Second),
			gocron.NewTask(SavSync),
		)
		if err != nil {
			logger.Errorf("%v\n", err)
		}
	}

	if backupInterval > 0 {
		go BackupTask(db)
		_, err := s.NewJob(
			gocron.DurationJob(backupInterval*time.Second),
			gocron.NewTask(BackupTask, db),
		)
		if err != nil {
			logger.Error(err)
		}
	}

	_, err := s.NewJob(
		gocron.DurationJob(300*time.Second),
		gocron.NewTask(system.LimitCacheDir, filepath.Join(os.TempDir(), "palworldsav-"), 5),
	)
	if err != nil {
		logger.Errorf("%v\n", err)
	}

	s.Start()
}

func Shutdown() {
	s := getScheduler()
	err := s.Shutdown()
	if err != nil {
		logger.Errorf("%v\n", err)
	}
}

func initScheduler() gocron.Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		logger.Errorf("%v\n", err)
	}
	return s
}

func getScheduler() gocron.Scheduler {
	if s == nil {
		return initScheduler()
	}
	return s
}

// PlayerSyncByServer synchronizes player data for a specific server
func PlayerSyncByServer(db *bbolt.DB, serverId string) {
	logger.Infof("Syncing players for server %s...\n", serverId)

	server, exists := config.GetServer(serverId)
	if !exists {
		logger.Errorf("Server %s not found\n", serverId)
		return
	}

	onlinePlayers, err := tool.ShowPlayersWithConfig(server)
	if err != nil {
		logger.Errorf("Failed to get online players for server %s: %v\n", serverId, err)
		return
	}

	err = service.PutPlayersOnlineByServer(db, serverId, onlinePlayers)
	if err != nil {
		logger.Errorf("Failed to save online players for server %s: %v\n", serverId, err)
		return
	}

	logger.Infof("Player sync done for server %s\n", serverId)

	playerLogging := viper.GetBool("task.player_logging")
	if playerLogging {
		go PlayerLoggingByServer(server, onlinePlayers)
	}

	kickInterval := viper.GetBool("manage.kick_non_whitelist")
	if kickInterval {
		go CheckAndKickPlayersByServer(db, server, onlinePlayers)
	}
}

// SavSyncByServer synchronizes save data for a specific server
func SavSyncByServer(serverId string) {
	logger.Infof("Syncing save for server %s...\n", serverId)

	server, exists := config.GetServer(serverId)
	if !exists {
		logger.Errorf("Server %s not found\n", serverId)
		return
	}

	if server.Save.Path == "" {
		logger.Warnf("Save path not configured for server %s, skipping\n", serverId)
		return
	}

	err := tool.DecodeWithConfig(server, server.Save.Path)
	if err != nil {
		logger.Errorf("Failed to decode save for server %s: %v\n", serverId, err)
		return
	}

	logger.Infof("Sav sync done for server %s\n", serverId)
}
