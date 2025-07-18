package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zaigie/palworld-server-tool/internal/database"
	"go.etcd.io/bbolt"
)

// ListPlayersByServer returns all players for a specific server
func ListPlayersByServer(db *bbolt.DB, serverId string) ([]database.TersePlayer, error) {
	var players []database.TersePlayer
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var player database.Player
			if err := json.Unmarshal(v, &player); err != nil {
				return err
			}
			// Filter by server ID
			if player.ServerId == serverId {
				players = append(players, player.TersePlayer)
			}
			return nil
		})
	})
	return players, err
}

// GetPlayerByServer returns a specific player from a specific server
func GetPlayerByServer(db *bbolt.DB, serverId, playerUid string) (*database.Player, error) {
	var player database.Player
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		if b == nil {
			return ErrNoRecord
		}

		v := b.Get([]byte(fmt.Sprintf("%s_%s", serverId, playerUid)))
		if v == nil {
			return ErrNoRecord
		}

		return json.Unmarshal(v, &player)
	})
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// ListGuildsByServer returns all guilds for a specific server
func ListGuildsByServer(db *bbolt.DB, serverId string) ([]database.Guild, error) {
	var guilds []database.Guild
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var guild database.Guild
			if err := json.Unmarshal(v, &guild); err != nil {
				return err
			}
			// Filter by server ID
			if guild.ServerId == serverId {
				guilds = append(guilds, guild)
			}
			return nil
		})
	})
	return guilds, err
}

// GetGuildByServer returns a specific guild from a specific server
func GetGuildByServer(db *bbolt.DB, serverId, adminPlayerUid string) (*database.Guild, error) {
	var guild database.Guild
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		if b == nil {
			return ErrNoRecord
		}

		v := b.Get([]byte(fmt.Sprintf("%s_%s", serverId, adminPlayerUid)))
		if v == nil {
			return ErrNoRecord
		}

		return json.Unmarshal(v, &guild)
	})
	if err != nil {
		return nil, err
	}
	return &guild, nil
}

// ListWhitelistByServer returns all whitelist entries for a specific server
func ListWhitelistByServer(db *bbolt.DB, serverId string) ([]database.PlayerW, error) {
	var whitelist []database.PlayerW
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("whitelist"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var player database.PlayerW
			if err := json.Unmarshal(v, &player); err != nil {
				return err
			}
			// Filter by server ID
			if player.ServerId == serverId {
				whitelist = append(whitelist, player)
			}
			return nil
		})
	})
	return whitelist, err
}

// AddWhitelistByServer adds a player to whitelist for a specific server
func AddWhitelistByServer(db *bbolt.DB, serverId string, player database.PlayerW) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("whitelist"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("whitelist"))
			if err != nil {
				return err
			}
		}

		player.ServerId = serverId
		data, err := json.Marshal(player)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("%s_%s", serverId, player.PlayerUID)
		if player.PlayerUID == "" {
			key = fmt.Sprintf("%s_%s", serverId, player.SteamID)
		}

		return b.Put([]byte(key), data)
	})
}

// RemoveWhitelistByServer removes a player from whitelist for a specific server
func RemoveWhitelistByServer(db *bbolt.DB, serverId string, identifier string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("whitelist"))
		if b == nil {
			return nil
		}

		key := fmt.Sprintf("%s_%s", serverId, identifier)
		return b.Delete([]byte(key))
	})
}

// ListRconCommandsByServer returns all RCON commands for a specific server
func ListRconCommandsByServer(db *bbolt.DB, serverId string) ([]database.RconCommandList, error) {
	var commands []database.RconCommandList
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcon_commands"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var command database.RconCommandList
			if err := json.Unmarshal(v, &command); err != nil {
				return err
			}
			// Filter by server ID
			if command.ServerId == serverId {
				commands = append(commands, command)
			}
			return nil
		})
	})
	return commands, err
}

// AddRconCommandByServer adds an RCON command for a specific server
func AddRconCommandByServer(db *bbolt.DB, serverId string, command database.RconCommandList) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcon_commands"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("rcon_commands"))
			if err != nil {
				return err
			}
		}

		command.ServerId = serverId
		data, err := json.Marshal(command)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("%s_%s", serverId, command.UUID)
		return b.Put([]byte(key), data)
	})
}

// RemoveRconCommandByServer removes an RCON command for a specific server
func RemoveRconCommandByServer(db *bbolt.DB, serverId, uuid string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcon_commands"))
		if b == nil {
			return nil
		}

		key := fmt.Sprintf("%s_%s", serverId, uuid)
		return b.Delete([]byte(key))
	})
}

// ListBackupsByServer returns all backups for a specific server
func ListBackupsByServer(db *bbolt.DB, serverId string) ([]database.Backup, error) {
	var backups []database.Backup
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var backup database.Backup
			if err := json.Unmarshal(v, &backup); err != nil {
				return err
			}
			// Filter by server ID
			if backup.ServerId == serverId {
				backups = append(backups, backup)
			}
			return nil
		})
	})
	return backups, err
}

// GetBackupByServer returns a specific backup for a specific server
func GetBackupByServer(db *bbolt.DB, serverId, backupId string) (*database.Backup, error) {
	var backup database.Backup
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		if b == nil {
			return ErrNoRecord
		}

		v := b.Get([]byte(fmt.Sprintf("%s_%s", serverId, backupId)))
		if v == nil {
			return ErrNoRecord
		}

		return json.Unmarshal(v, &backup)
	})
	if err != nil {
		return nil, err
	}
	return &backup, nil
}

// DeleteBackupByServer deletes a specific backup for a specific server
func DeleteBackupByServer(db *bbolt.DB, serverId, backupId string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		if b == nil {
			return nil
		}

		key := fmt.Sprintf("%s_%s", serverId, backupId)
		return b.Delete([]byte(key))
	})
}

// PutPlayersByServer stores players for a specific server
func PutPlayersByServer(db *bbolt.DB, serverId string, players []database.Player) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("players"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("players"))
			if err != nil {
				return err
			}
		}

		for _, player := range players {
			player.ServerId = serverId
			data, err := json.Marshal(player)
			if err != nil {
				return err
			}

			key := fmt.Sprintf("%s_%s", serverId, player.PlayerUid)
			if err := b.Put([]byte(key), data); err != nil {
				return err
			}
		}

		return nil
	})
}

// PutGuildsByServer stores guilds for a specific server
func PutGuildsByServer(db *bbolt.DB, serverId string, guilds []database.Guild) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("guilds"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("guilds"))
			if err != nil {
				return err
			}
		}

		for _, guild := range guilds {
			guild.ServerId = serverId
			data, err := json.Marshal(guild)
			if err != nil {
				return err
			}

			key := fmt.Sprintf("%s_%s", serverId, guild.AdminPlayerUid)
			if err := b.Put([]byte(key), data); err != nil {
				return err
			}
		}

		return nil
	})
}

// PutPlayersOnlineByServer stores online players for a specific server
func PutPlayersOnlineByServer(db *bbolt.DB, serverId string, players []database.OnlinePlayer) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("online_players"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("online_players"))
			if err != nil {
				return err
			}
		}

		// Clear existing online players for this server
		c := b.Cursor()
		prefix := []byte(fmt.Sprintf("%s_", serverId))
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			c.Delete()
		}

		// Add new online players
		for _, player := range players {
			player.ServerId = serverId
			data, err := json.Marshal(player)
			if err != nil {
				return err
			}

			key := fmt.Sprintf("%s_%s", serverId, player.PlayerUid)
			if err := b.Put([]byte(key), data); err != nil {
				return err
			}
		}

		return nil
	})
}

// AddBackupByServer adds a backup record for a specific server
func AddBackupByServer(db *bbolt.DB, serverId string, backup database.Backup) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("backups"))
			if err != nil {
				return err
			}
		}

		backup.ServerId = serverId
		data, err := json.Marshal(backup)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("%s_%s", serverId, backup.BackupId)
		return b.Put([]byte(key), data)
	})
}

// PutWhitelistByServer stores whitelist for a specific server
func PutWhitelistByServer(db *bbolt.DB, serverId string, players []database.PlayerW) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("whitelist"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("whitelist"))
			if err != nil {
				return err
			}
		}

		// Clear existing whitelist for this server
		c := b.Cursor()
		prefix := []byte(fmt.Sprintf("%s_", serverId))
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			c.Delete()
		}

		// Add new whitelist players
		for _, player := range players {
			player.ServerId = serverId
			data, err := json.Marshal(player)
			if err != nil {
				return err
			}

			key := fmt.Sprintf("%s_%s", serverId, player.PlayerUID)
			if player.PlayerUID == "" {
				key = fmt.Sprintf("%s_%s", serverId, player.SteamID)
			}
			if err := b.Put([]byte(key), data); err != nil {
				return err
			}
		}

		return nil
	})
}

// PutRconCommandByServer updates an RCON command for a specific server
func PutRconCommandByServer(db *bbolt.DB, serverId, uuid string, command database.RconCommandList) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("rcon_commands"))
		if b == nil {
			var err error
			b, err = tx.CreateBucket([]byte("rcon_commands"))
			if err != nil {
				return err
			}
		}

		command.ServerId = serverId
		command.UUID = uuid
		data, err := json.Marshal(command)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("%s_%s", serverId, uuid)
		return b.Put([]byte(key), data)
	})
}

// ListBackupsByServerWithTimeRange returns backups for a specific server within a time range
func ListBackupsByServerWithTimeRange(db *bbolt.DB, serverId string, startTime, endTime time.Time) ([]database.Backup, error) {
	var backups []database.Backup
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("backups"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var backup database.Backup
			if err := json.Unmarshal(v, &backup); err != nil {
				return err
			}
			// Filter by server ID and time range
			if backup.ServerId == serverId {
				backupTime := backup.SaveTime
				if (startTime.IsZero() || backupTime.After(startTime)) &&
					(endTime.IsZero() || backupTime.Before(endTime)) {
					backups = append(backups, backup)
				}
			}
			return nil
		})
	})
	return backups, err
}
