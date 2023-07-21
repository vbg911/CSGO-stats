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

var steps = 0

type playersFootSteps map[uint64]int
type playersDuckKills map[uint64]int

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

		p.RegisterEventHandler(func(e events.Footstep) { handleFootstep(e, playersFootStep) })
		p.RegisterEventHandler(func(e events.Kill) { handleKill(e, playersDuckKill) })
		p.RegisterEventHandler(handleRoundStart)

		err = p.ParseToEnd()
		if err != nil {
			log.Panic("failed to parse demo: ", err)
		}

		fmt.Println(p.GameState().TeamCounterTerrorists().ClanName())
		fmt.Println(p.GameState().TeamTerrorists().ClanName())
		fmt.Println(p.GameState().TeamCounterTerrorists().Members())
		fmt.Println(p.GameState().TeamTerrorists().Members())

		fmt.Println(p.GameState().Participants().TeamMembers(common.TeamCounterTerrorists))
		fmt.Println(p.GameState().Participants().TeamMembers(common.TeamTerrorists))

		players := p.GameState().Participants().Playing()
		var stats []playerStats
		for _, p := range players {
			stats = append(stats, statsFor(p, playersFootStep, playersDuckKill))
		}

		//fmt.Println(playersFootStep)
		//fmt.Println(playersDuckKill)
		//fmt.Println("Все игроки вместе сделали: ", steps)

		//for _, player := range stats {
		//	fmt.Println(player.formatString() + "\n")
		//}
		steps = 0
	}
}

func handleFootstep(e events.Footstep, footSteps playersFootSteps) {
	steps += 1
	footSteps[e.Player.SteamID64] += 1
	footSteps[0] += 1
	//fmt.Println("сделал шаг", e.Player.Name)
}

func handleKill(e events.Kill, duckKills playersDuckKills) {
	//fmt.Println(e.Killer.Name, " убил ", e.Victim.Name, " с помощью ", e.Weapon.String(), "e.Killer.IsDucking()", e.Killer.IsDucking())
	if e.Killer.IsDucking() {
		duckKills[e.Killer.SteamID64] += 1
	}
}

func handleRoundStart(e events.RoundStart) {
	//fmt.Println("Раунд начался")
}

type playerStats struct {
	SteamID64 uint64 `json:"steamID64"`
	Name      string `json:"name"`
	MVP       int    `json:"mvp"`
	Kills     int    `json:"kills"`
	Deaths    int    `json:"deaths"`
	Assists   int    `json:"assists"`
	FootSteps int    `json:"footSteps"`
	DuckKills int    `json:"duckKills"`
}

func (s playerStats) formatString() string {
	return "Player name:  " + s.Name +
		"\nSteam id64:   " + strconv.FormatUint(s.SteamID64, 10) +
		"\nMVPs:         " + strconv.Itoa(s.MVP) +
		"\nKills:        " + strconv.Itoa(s.Kills) +
		"\nDeath:        " + strconv.Itoa(s.Deaths) +
		"\nAssists:      " + strconv.Itoa(s.Assists) +
		"\nFootSteps:    " + strconv.Itoa(s.FootSteps) +
		"\nDuckKills:    " + strconv.Itoa(s.DuckKills)
}

func statsFor(p *common.Player, fs playersFootSteps, dk playersDuckKills) playerStats {
	return playerStats{
		SteamID64: p.SteamID64,
		Name:      p.Name,
		MVP:       p.MVPs(),
		Kills:     p.Kills(),
		Deaths:    p.Deaths(),
		Assists:   p.Assists(),
		FootSteps: fs[p.SteamID64],
		DuckKills: dk[p.SteamID64],
	}
}
