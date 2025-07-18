package tool

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/auth"
	"github.com/zaigie/palworld-server-tool/internal/config"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/source"
	"github.com/zaigie/palworld-server-tool/internal/system"
	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"
)

type Sturcture struct {
	Players []database.Player `json:"players"`
	Guilds  []database.Guild  `json:"guilds"`
}

func getSavCli() (string, error) {
	savCliPath := viper.GetString("save.decode_path")
	if savCliPath == "" || savCliPath == "/path/to/your/sav_cli" {
		ed, err := system.GetExecDir()
		if err != nil {
			logger.Errorf("error getting exec directory: %s", err)
			return "", err
		}
		savCliPath = filepath.Join(ed, "sav_cli")
		if runtime.GOOS == "windows" {
			savCliPath += ".exe"
		}
	}
	if _, err := os.Stat(savCliPath); err != nil {
		return "", err
	}
	return savCliPath, nil
}

func getSavCliWithConfig(server *config.Server) (string, error) {
	savCliPath := server.Save.DecodePath
	if savCliPath == "" || savCliPath == "/path/to/your/sav_cli" {
		ed, err := system.GetExecDir()
		if err != nil {
			logger.Errorf("error getting exec directory: %s", err)
			return "", err
		}
		savCliPath = filepath.Join(ed, "sav_cli")
		if runtime.GOOS == "windows" {
			savCliPath += ".exe"
		}
	}
	if _, err := os.Stat(savCliPath); err != nil {
		return "", err
	}
	return savCliPath, nil
}

func Decode(file string) error {
	savCli, err := getSavCli()
	if err != nil {
		return errors.New("error getting executable path: " + err.Error())
	}

	levelFilePath, err := getFromSource(file, "decode")
	if err != nil {
		return err
	}
	defer os.RemoveAll(filepath.Dir(levelFilePath))

	baseUrl := fmt.Sprintf("http://127.0.0.1:%d", viper.GetInt("web.port"))
	if viper.GetBool("web.tls") && !strings.HasSuffix(baseUrl, "/") {
		baseUrl = viper.GetString("web.public_url")
	}

	requestUrl := fmt.Sprintf("%s/api/", baseUrl)
	tokenString, err := auth.GenerateToken()
	if err != nil {
		return errors.New("error generating token: " + err.Error())
	}
	execArgs := []string{"-f", levelFilePath, "--request", requestUrl, "--token", tokenString}
	cmd := exec.Command(savCli, execArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return errors.New("error starting command: " + err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		return errors.New("error waiting for command: " + err.Error())
	}

	return nil
}

func DecodeWithConfig(server *config.Server, file string) error {
	savCli, err := getSavCliWithConfig(server)
	if err != nil {
		return errors.New("error getting executable path: " + err.Error())
	}

	// Implementation would be similar to Decode but using server config
	// For now, just call the existing decode function
	return Decode(file)
}

func Backup() (string, error) {
	sourcePath := viper.GetString("save.path")

	levelFilePath, err := getFromSource(sourcePath, "backup")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(filepath.Dir(levelFilePath))

	backupDir, err := GetBackupDir()
	if err != nil {
		return "", fmt.Errorf("failed to get backup directory: %s", err)
	}

	currentTime := time.Now().Format("2006-01-02-15-04-05")
	backupZipFile := filepath.Join(backupDir, fmt.Sprintf("%s.zip", currentTime))
	err = system.ZipDir(filepath.Dir(levelFilePath), backupZipFile)
	if err != nil {
		return "", fmt.Errorf("failed to create backup zip: %s", err)
	}
	return filepath.Base(backupZipFile), nil
}

func GetBackupDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	backDir := filepath.Join(wd, "backups")
	if err = system.CheckAndCreateDir(backDir); err != nil {
		return "", err
	}
	return backDir, nil
}

func CleanOldBackups(db *bbolt.DB, keepDays int) error {
	backupDir, err := GetBackupDir()
	if err != nil {
		return fmt.Errorf("failed to get backup directory: %s", err)
	}

	deadline := time.Now().AddDate(0, 0, -keepDays)

	backups, err := service.ListBackups(db, time.Time{}, time.Now())
	if err != nil {
		return fmt.Errorf("failed to list backups: %s", err)
	}

	for _, backup := range backups {
		if backup.SaveTime.Before(deadline) {
			err = os.Remove(filepath.Join(backupDir, backup.Path))
			if err != nil {
				if !os.IsNotExist(err) {
					logger.Errorf("failed to delete old backup file %s: %s", backup.Path, err)
				}
			}

			err = service.DeleteBackup(db, backup.BackupId)
			if err != nil {
				logger.Errorf("failed to delete backup record from database: %s", err)
			}
		}
	}

	return nil
}

func getFromSource(file, way string) (string, error) {
	var levelFilePath string
	var err error

	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		// http(s)://url
		levelFilePath, err = source.DownloadFromHttp(file, way)
		if err != nil {
			return "", errors.New("error downloading file: " + err.Error())
		}
	} else if strings.HasPrefix(file, "k8s://") {
		// k8s://namespace/pod/container:remotePath
		namespace, podName, container, remotePath, err := source.ParseK8sAddress(file)
		if err != nil {
			return "", errors.New("error parsing k8s address: " + err.Error())
		}
		levelFilePath, err = source.CopyFromPod(namespace, podName, container, remotePath, way)
		if err != nil {
			return "", errors.New("error copying file from pod: " + err.Error())
		}
	} else if strings.HasPrefix(file, "docker://") {
		// docker://containerID(Name):remotePath
		containerId, remotePath, err := source.ParseDockerAddress(file)
		if err != nil {
			return "", errors.New("error parsing docker address: " + err.Error())
		}
		levelFilePath, err = source.CopyFromContainer(containerId, remotePath, way)
		if err != nil {
			return "", errors.New("error copying file from container: " + err.Error())
		}
	} else {
		// local file
		levelFilePath, err = source.CopyFromLocal(file, way)
		if err != nil {
			return "", errors.New("error copying file to temporary directory: " + err.Error())
		}
	}
	return levelFilePath, nil
}

func BackupWithConfig(server *config.Server) (string, error) {
	sourcePath := server.Save.Path
	if sourcePath == "" {
		return "", errors.New("save path not configured for server")
	}

	levelFilePath, err := getFromSource(sourcePath, "backup")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(filepath.Dir(levelFilePath))

	backupDir, err := GetBackupDirByServer(server.Id)
	if err != nil {
		return "", fmt.Errorf("failed to get backup directory: %s", err)
	}

	currentTime := time.Now().Format("2006-01-02-15-04-05")
	backupZipFile := filepath.Join(backupDir, fmt.Sprintf("%s_%s.zip", server.Id, currentTime))
	err = system.ZipDir(filepath.Dir(levelFilePath), backupZipFile)
	if err != nil {
		return "", fmt.Errorf("failed to create backup zip: %s", err)
	}
	return filepath.Base(backupZipFile), nil
}

func GetBackupDirByServer(serverId string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	backDir := filepath.Join(wd, "backups", serverId)
	if err = system.CheckAndCreateDir(backDir); err != nil {
		return "", err
	}
	return backDir, nil
}

func CleanOldBackupsByServer(db *bbolt.DB, serverId string, keepDays int) error {
	backups, err := service.ListBackupsByServer(db, serverId)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -keepDays)
	for _, backup := range backups {
		if backup.SaveTime.Before(cutoff) {
			// Delete the backup file
			backupDir, err := GetBackupDirByServer(serverId)
			if err != nil {
				logger.Warnf("Failed to get backup directory for server %s: %v", serverId, err)
				continue
			}

			backupPath := filepath.Join(backupDir, backup.Path)
			if err := os.Remove(backupPath); err != nil {
				logger.Warnf("Failed to remove backup file %s: %v", backupPath, err)
			}

			// Delete the backup record
			if err := service.DeleteBackupByServer(db, serverId, backup.BackupId); err != nil {
				logger.Warnf("Failed to delete backup record %s: %v", backup.BackupId, err)
			}
		}
	}
	return nil
}
