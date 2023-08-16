package main

import (
	"CSGO-stats/internal/analyze"
	"CSGO-stats/internal/demoparser"
	"CSGO-stats/internal/structures"
	"CSGO-stats/internal/visualization"
	"fmt"
	"os"
)

func main() {
	demoFolder := "./demos"
	var (
		tournament string
		match      string
	)

	entries, err := os.ReadDir(demoFolder)
	checkError(err)

	var matchMaps []structures.MapStats
	var tournamentMatches []structures.MatchStats
	var tournaments []structures.Tournament
	for _, e := range entries {
		if e.IsDir() {
			tournament = e.Name()
			matches, err := os.ReadDir(demoFolder + "/" + tournament)
			checkError(err)
			for _, e := range matches {
				if e.IsDir() {
					match = e.Name()
					maps, err := os.ReadDir(demoFolder + "/" + tournament + "/" + match)
					checkError(err)
					for _, e := range maps {
						fmt.Println("Tournament: " + tournament + " match: " + match + " file: " + e.Name())
						pathToDemo := demoFolder + "/" + tournament + "/" + match + "/" + e.Name()
						mapStats, err := demoparser.ParseDemo(tournament, match, e.Name(), pathToDemo)

						matchMaps = append(matchMaps, mapStats)

						//fmt.Println("Все игроки вместе сделали: ", mapStats.PlayersFootStep[0], " шагов")
						//fmt.Println("Все игроки вместе сделали: ", mapStats.PlayersWeaponShot[0], " выстрелов")
						//fmt.Println("Все игроки вместе сделали: ", mapStats.PlayersWeaponReload[0], " перезарядок")
						//fmt.Println("Все игроки вместе сделали: ", mapStats.PlayersJump[0], " прыжков")
						//fmt.Println("Все игроки в сумме дропнули бомбу: ", mapStats.PlayersBombDrop[0], " раз(а)")
						//fmt.Println("Все игроки в сумме кинули Smoke: ", mapStats.PlayersSmoke[0], " раз(а)")
						//fmt.Println("Все игроки в сумме кинули HE Grenade: ", mapStats.PlayersHEGrenade[0], " раз(а)")
						//fmt.Println("Все игроки в сумме кинули Molotov: ", mapStats.PlayersMolotov[0], " раз(а)")
						//fmt.Println("Все игроки в сумме кинули Incendiary Grenade: ", mapStats.PlayersIncendiaryGrenade[0], " раз(а)")
						//fmt.Println("Все игроки в сумме кинули Flashbang: ", mapStats.PlayersFlashbang[0], " раз(а)")
						//fmt.Println("Все игроки в сумме кинули Decoy: ", mapStats.PlayersDecoyGrenade[0], " раз(а)")

						visualization.GenerateHeatMap(mapStats.FirePoints, mapStats.MapRadarImg, mapStats.DemoName+".jpeg", "WeaponFire")
						visualization.GenerateHeatMap(mapStats.DeathPoints, mapStats.MapRadarImg, mapStats.DemoName+".jpeg", "PlayerDeath")
						visualization.GenerateHeatMap(mapStats.GrenadePoints, mapStats.MapRadarImg, mapStats.DemoName+".jpeg", "GrenadeThrow")

						visualization.GenerateTrajectories(mapStats.MapMetadata, mapStats.MapRadarImg, mapStats.NadesProjectiles, mapStats.NadesInferno, "GrenadeTrajectories\\"+mapStats.TournamentName, mapStats.DemoName+".jpeg")
						checkError(err)
					}
				}
				tournamentMatches = append(tournamentMatches, structures.MatchStats{
					TournamentName: tournament,
					MatchName:      match,
					Maps:           matchMaps,
				})
				matchMaps = nil
			}
			tournaments = append(tournaments, structures.Tournament{
				TournamentName: tournament,
				Matches:        tournamentMatches,
			})
			tournamentMatches = nil
		}
	}

	for _, i := range tournaments {
		analyze.AnalyzeTournament(i)
		fmt.Println("турнир: " + i.TournamentName)
		for _, j := range i.Matches {
			fmt.Println("	матч: " + j.MatchName)
			for _, k := range j.Maps {
				fmt.Println("		карта: " + k.MapName)
			}
		}
	}

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
