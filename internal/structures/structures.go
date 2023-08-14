package structures

import (
	"github.com/golang/geo/r2"
	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	"image"
	"strconv"
)

type NadeTrajectories map[int]map[int64]*common.GrenadeProjectile
type Infernos map[int]map[int64]*common.Inferno
type FootSteps map[uint64]int
type DuckKills map[uint64]int
type FlashedKills map[uint64]int
type AirborneKills map[uint64]int
type WallbangKills map[uint64]int
type SmokeKills map[uint64]int
type NoScopeKills map[uint64]int
type WeaponShots map[uint64]int
type WeaponReloads map[uint64]int
type Jumps map[uint64]int
type Smokes map[uint64]int
type HEGrenades map[uint64]int
type Flashbangs map[uint64]int
type BombDrops map[uint64]int
type IncendiaryGrenades map[uint64]int
type Molotovs map[uint64]int
type DecoyGrenades map[uint64]int

type MatchStats struct {
	TournamentName string     `json:"tournamentName"`
	MatchName      string     `json:"matchName"`
	Maps           []MapStats `json:"maps"`
}

type MapStats struct {
	TournamentName           string             `json:"tournamentName"`
	MatchName                string             `json:"matchName"`
	DemoName                 string             `json:"demoName"`
	DemoHash                 string             `json:"demoHash"`
	DemoPath                 string             `json:"demoPath"`
	MapMetadata              ex.Map             `json:"mapMetadata"`
	MapRadarImg              image.Image        `json:"mapRadarImg"`
	MapName                  string             `json:"mapName"`
	Players                  []PlayerStats      `json:"players"`
	FirePoints               []r2.Point         `json:"firePoints"`
	DeathPoints              []r2.Point         `json:"deathPoints"`
	GrenadePoints            []r2.Point         `json:"grenadePoints"`
	NadesProjectiles         NadeTrajectories   `json:"nadesProjectiles"`
	NadesInferno             Infernos           `json:"nadesInferno"`
	PlayersFootStep          FootSteps          `json:"playersFootStep"`
	PlayersWeaponShot        WeaponShots        `json:"playersWeaponShot"`
	PlayersWeaponReload      WeaponReloads      `json:"playersWeaponReload"`
	PlayersJump              Jumps              `json:"playersJump"`
	PlayersBombDrop          BombDrops          `json:"playersBombDrop"`
	PlayersSmoke             Smokes             `json:"playersSmoke"`
	PlayersHEGrenade         HEGrenades         `json:"playersHEGrenade"`
	PlayersMolotov           Molotovs           `json:"playersMolotov"`
	PlayersIncendiaryGrenade IncendiaryGrenades `json:"playersIncendiaryGrenade"`
	PlayersFlashbang         Flashbangs         `json:"playersFlashbang"`
	PlayersDecoyGrenade      DecoyGrenades      `json:"playersDecoyGrenade"`
}

type PlayerStats struct {
	SteamID64     uint64 `json:"steamID64"`
	Name          string `json:"name"`
	TeamName      string `json:"teamName"`
	CrosshairCode string `json:"crosshairCode"`
	MVP           int    `json:"mvp"`
	Kills         int    `json:"kills"`
	Deaths        int    `json:"deaths"`
	Assists       int    `json:"assists"`
	FootSteps     int    `json:"footSteps"`
	DuckKills     int    `json:"duckKills"`
	FlashedKills  int    `json:"flashKills"`
	AirborneKills int    `json:"airborneKills"`
	WallbangKills int    `json:"wallbangKills"`
	SmokeKills    int    `json:"smokeKills"`
	NoScopeKills  int    `json:"noScopeKills"`
	WeaponShots   int    `json:"weaponShots"`
	WeaponReloads int    `json:"weaponReloads"`
	PlayerJumps   int    `json:"playerJumps"`
	BombDrops     int    `json:"bombDrops"`
	Smokes        int    `json:"smokes"`
	HEnades       int    `json:"HEnades"`
	Molotov       int    `json:"molotov"`
	CTmoly        int    `json:"CTmoly"`
	Flasbang      int    `json:"flasbang"`
	Decoy         int    `json:"decoy"`
}

func (s PlayerStats) FormatString() string {
	return "Player name:    " + s.Name +
		"\nTeamName:       " + s.TeamName +
		"\nSteam id64:     " + strconv.FormatUint(s.SteamID64, 10) +
		"\nCrosshairCode   " + s.CrosshairCode +
		"\nMVPs:           " + strconv.Itoa(s.MVP) +
		"\nKills:          " + strconv.Itoa(s.Kills) +
		"\nDeath:          " + strconv.Itoa(s.Deaths) +
		"\nAssists:        " + strconv.Itoa(s.Assists) +
		"\nFootSteps:      " + strconv.Itoa(s.FootSteps) +
		"\nDuckKills:      " + strconv.Itoa(s.DuckKills) +
		"\nFlashedKills:   " + strconv.Itoa(s.FlashedKills) +
		"\nAirborneKills:  " + strconv.Itoa(s.AirborneKills) +
		"\nWallbangKills:  " + strconv.Itoa(s.WallbangKills) +
		"\nSmokeKills:     " + strconv.Itoa(s.SmokeKills) +
		"\nNoScopeKills:   " + strconv.Itoa(s.NoScopeKills) +
		"\nWeaponShots:    " + strconv.Itoa(s.WeaponShots) +
		"\nWeaponReloads:  " + strconv.Itoa(s.WeaponReloads) +
		"\nJumps:          " + strconv.Itoa(s.PlayerJumps) +
		"\nBombDrops       " + strconv.Itoa(s.BombDrops) +
		"\nSmokes          " + strconv.Itoa(s.Smokes) +
		"\nHE Grenade      " + strconv.Itoa(s.HEnades) +
		"\nMolotov         " + strconv.Itoa(s.Molotov) +
		"\nCT molotov      " + strconv.Itoa(s.CTmoly) +
		"\nFlashbang       " + strconv.Itoa(s.Flasbang) +
		"\nDecoy Grenade   " + strconv.Itoa(s.Decoy)
}
