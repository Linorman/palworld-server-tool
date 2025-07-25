package database

import "time"

type Pal struct {
	Level          int32    `json:"level"`
	Exp            int64    `json:"exp"`
	Hp             int64    `json:"hp"`
	MaxHp          int64    `json:"max_hp"`
	Type           string   `json:"type"`
	Gender         string   `json:"gender"`
	Nickname       string   `json:"nickname"`
	IsLucky        bool     `json:"is_lucky"`
	IsBoss         bool     `json:"is_boss"`
	IsTower        bool     `json:"is_tower"`
	Workspeed      int32    `json:"workspeed"`
	Melee          int32    `json:"melee"`
	Ranged         int32    `json:"ranged"`
	Defense        int32    `json:"defense"`
	Rank           int32    `json:"rank"`
	RankAttack     int32    `json:"rank_attack"`
	RankDefence    int32    `json:"rank_defence"`
	RankCraftspeed int32    `json:"rank_craftspeed"`
	Skills         []string `json:"skills"`
}

type OnlinePlayer struct {
	ServerId   string    `json:"server_id"`
	PlayerUid  string    `json:"player_uid"`
	SteamId    string    `json:"steam_id"`
	Nickname   string    `json:"nickname"`
	Ip         string    `json:"ip"`
	Ping       float64   `json:"ping"`
	LocationX  float64   `json:"location_x"`
	LocationY  float64   `json:"location_y"`
	Level      int32     `json:"level"`
	LastOnline time.Time `json:"last_online"`
}

type GuildPlayer struct {
	PlayerUid string `json:"player_uid"`
	Nickname  string `json:"nickname"`
}

type TersePlayer struct {
	ServerId       string           `json:"server_id"`
	PlayerUid      string           `json:"player_uid"`
	Nickname       string           `json:"nickname"`
	Level          int32            `json:"level"`
	Exp            int64            `json:"exp"`
	Hp             int64            `json:"hp"`
	MaxHp          int64            `json:"max_hp"`
	ShieldHp       int64            `json:"shield_hp"`
	ShieldMaxHp    int64            `json:"shield_max_hp"`
	MaxStatusPoint int32            `json:"max_status_point"`
	StatusPoint    map[string]int32 `json:"status_point"`
	FullStomach    float64          `json:"full_stomach"`
	SaveLastOnline string           `json:"save_last_online"`
	OnlinePlayer
}

type Player struct {
	TersePlayer
	Pals  []*Pal `json:"pals"`
	Items *Items `json:"items"`
}

type BaseCamp struct {
	Id        string  `json:"id"`
	Area      float64 `json:"area"`
	LocationX float64 `json:"location_x"`
	LocationY float64 `json:"location_y"`
}

type Guild struct {
	ServerId       string         `json:"server_id"`
	Name           string         `json:"name"`
	BaseCampLevel  int32          `json:"base_camp_level"`
	AdminPlayerUid string         `json:"admin_player_uid"`
	Players        []*GuildPlayer `json:"players"`
	BaseCamp       []BaseCamp     `json:"base_camp"`
}

type PlayerW struct {
	ServerId  string `json:"server_id"`
	Name      string `json:"name"`
	SteamID   string `json:"steam_id"`
	PlayerUID string `json:"player_uid"`
}

type RconCommand struct {
	Command     string `json:"command"`
	Placeholder string `json:"placeholder"`
	Remark      string `json:"remark"`
}

type RconCommandList struct {
	ServerId string `json:"server_id"`
	UUID     string `json:"uuid"`
	RconCommand
}

type Items struct {
	CommonContainerId           []*Item `json:"CommonContainerId"`
	DropSlotContainerId         []*Item `json:"DropSlotContainerId"`
	EssentialContainerId        []*Item `json:"EssentialContainerId"`
	FoodEquipContainerId        []*Item `json:"FoodEquipContainerId"`
	PlayerEquipArmorContainerId []*Item `json:"PlayerEquipArmorContainerId"`
	WeaponLoadOutContainerId    []*Item `json:"WeaponLoadOutContainerId"`
}

type Item struct {
	SlotIndex  int32  `json:"SlotIndex"`
	ItemId     string `json:"ItemId"`
	StackCount int32  `json:"StackCount"`
}

type Backup struct {
	ServerId string    `json:"server_id"`
	BackupId string    `json:"backup_id"`
	SaveTime time.Time `json:"save_time"`
	Path     string    `json:"path"`
}

// ServerInfo represents stored server configuration and status
type ServerInfo struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Enabled     bool                   `json:"enabled"`
	Status      string                 `json:"status"` // online, offline, error
	LastCheck   time.Time              `json:"last_check"`
	Config      map[string]interface{} `json:"config"`
}
