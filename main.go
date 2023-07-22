package main

import (
	"fmt"
	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	"log"
	"os"
	"strconv"
)

type playersFootSteps map[uint64]int
type playersDuckKills map[uint64]int
type playersFlashedKills map[uint64]int
type playersAirborneKills map[uint64]int
type playersWallbangKills map[uint64]int

func main() {
	entries, err := os.ReadDir("./demos")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println(e.Name())
		f, err := os.Open("./demos/" + e.Name())
		if err != nil {
			log.Panic("failed to open demo file: ", err)
		}
		defer f.Close()

		p := dem.NewParser(f)
		defer p.Close()

		playersFootStep := make(map[uint64]int, 10)
		playersDuckKill := make(map[uint64]int, 10)
		playersFlashedKill := make(map[uint64]int, 10)
		playersAirborneKill := make(map[uint64]int, 10)
		playersWallbangKill := make(map[uint64]int, 10)
		p.RegisterEventHandler(func(e events.Footstep) { handleFootstep(e, playersFootStep) })
		p.RegisterEventHandler(func(e events.Kill) {
			handleKill(e, playersDuckKill, playersFlashedKill, playersAirborneKill, playersWallbangKill)
		})
		p.RegisterEventHandler(func(e events.RoundStart) {})

		err = p.ParseToEnd()
		if err != nil {
			log.Panic("failed to parse demo: ", err)
		}

		players := p.GameState().Participants().Playing()
		var stats []playerStats
		for _, p := range players {
			stats = append(stats, statsFor(p, playersFootStep, playersDuckKill, playersFlashedKill, playersAirborneKill, playersWallbangKill))
		}

		fmt.Println("Все игроки вместе сделали: ", playersFootStep[0])

		for _, player := range stats {
			fmt.Println(player.formatString() + "\n")
		}
		println("\n")
	}
}

func handleFootstep(e events.Footstep, footSteps playersFootSteps) {
	footSteps[e.Player.SteamID64] += 1
	footSteps[0] += 1
}

func handleKill(e events.Kill, duckKills playersDuckKills, flashKills playersFlashedKills, airborneKills playersFlashedKills, wallbangKills playersWallbangKills) {
	if e.Killer.IsDucking() {
		duckKills[e.Killer.SteamID64] += 1
	}

	if e.Killer.IsBlinded() {
		flashKills[e.Killer.SteamID64] += 1
	}

	if e.Killer.IsAirborne() {
		airborneKills[e.Killer.SteamID64] += 1
	}

	if e.IsWallBang() {
		wallbangKills[e.Killer.SteamID64] += 1
	}
	//fmt.Println(parser.GameState().IngameTick())
}

type playerStats struct {
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
}

func (s playerStats) formatString() string {
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
		"\nAirborneKills   " + strconv.Itoa(s.AirborneKills) +
		"\nWallbangKills   " + strconv.Itoa(s.WallbangKills)
}

func statsFor(p *common.Player, fs playersFootSteps, dk playersDuckKills, fk playersFlashedKills, airk playersAirborneKills, wbk playersWallbangKills) playerStats {
	return playerStats{
		SteamID64:     p.SteamID64,
		Name:          p.Name,
		TeamName:      p.TeamState.ClanName(),
		CrosshairCode: p.CrosshairCode(),
		MVP:           p.MVPs(),
		Kills:         p.Kills(),
		Deaths:        p.Deaths(),
		Assists:       p.Assists(),
		FootSteps:     fs[p.SteamID64],
		DuckKills:     dk[p.SteamID64],
		FlashedKills:  fk[p.SteamID64],
		AirborneKills: airk[p.SteamID64],
		WallbangKills: wbk[p.SteamID64],
	}
}
