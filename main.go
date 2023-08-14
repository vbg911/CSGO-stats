package main

import (
	"CSGO-stats/internal/demoparser"
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
						//fmt.Println(mapStats)

						visualization.GenerateHeatMap(mapStats.FirePoints, mapStats.MapRadarImg, mapStats.DemoName+".jpeg", "WeaponFire")
						visualization.GenerateHeatMap(mapStats.DeathPoints, mapStats.MapRadarImg, mapStats.DemoName+".jpeg", "PlayerDeath")
						visualization.GenerateHeatMap(mapStats.GrenadePoints, mapStats.MapRadarImg, mapStats.DemoName+".jpeg", "GrenadeThrow")

						visualization.GenerateTrajectories(mapStats.MapMetadata, mapStats.MapRadarImg, mapStats.NadesProjectiles, mapStats.NadesInferno, "GrenadeTrajectories\\test4", mapStats.DemoName+".jpeg")

						checkError(err)
					}
				}
			}
		}
	}

}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
