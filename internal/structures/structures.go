package structures

import (
	"github.com/golang/geo/r2"
	"strconv"
)

type PlayersFootSteps map[uint64]int
type PlayersDuckKills map[uint64]int
type PlayersFlashedKills map[uint64]int
type PlayersAirborneKills map[uint64]int
type PlayersWallbangKills map[uint64]int
type PlayerSmokeKills map[uint64]int
type PlayerNoScopeKills map[uint64]int
type PlayerWeaponShots map[uint64]int
type PlayerWeaponReloads map[uint64]int
type PlayerJumps map[uint64]int
type PlayerSmokes map[uint64]int
type PlayerHEGrenades map[uint64]int
type PlayerFlashbangs map[uint64]int
type PlayerBombDrops map[uint64]int
type PlayerIncendiaryGrenades map[uint64]int
type PlayerMolotovs map[uint64]int
type PlayerDecoyGrenades map[uint64]int

type MatchStats struct {
	TournamentName string     `json:"tournamentName"`
	MatchName      string     `json:"matchName"`
	Maps           []MapStats `json:"maps"`
}

type MapStats struct {
	TournamentName string         `json:"tournamentName"`
	DemoHash       string         `json:"demoHash"`
	DemoPath       string         `json:"demoPath"`
	MapName        int            `json:"mapName"`
	Players        []PlayerStats  `json:"players"`
	FirePoints     []r2.Point     `json:"firePoints"`
	DeathPoints    []r2.Point     `json:"deathPoints"`
	GrenadePoints  []r2.Point     `json:"grenadePoints"`
	OverallStats   map[string]int `json:"overallStats"`
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
