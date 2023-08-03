package main

import (
	"fmt"
	"github.com/golang/geo/r2"
	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/msg"
	"github.com/markus-wa/go-heatmap/v2"
	"github.com/markus-wa/go-heatmap/v2/schemes"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"
)

type playersFootSteps map[uint64]int
type playersDuckKills map[uint64]int
type playersFlashedKills map[uint64]int
type playersAirborneKills map[uint64]int
type playersWallbangKills map[uint64]int
type playerSmokeKills map[uint64]int
type playerNoScopeKills map[uint64]int
type playerWeaponShots map[uint64]int
type playerWeaponReloads map[uint64]int
type playerJumps map[uint64]int
type playerSmokes map[uint64]int
type playerHEGrenades map[uint64]int
type playerFlashbangs map[uint64]int
type playerBombDrops map[uint64]int
type playerIncendiaryGrenades map[uint64]int
type playerMolotovs map[uint64]int
type playerDecoyGrenades map[uint64]int

const (
	dotSize     = 15
	opacity     = 128
	jpegQuality = 100
)

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

		header, err := p.ParseHeader()
		checkError(err)

		playersFootStep := make(playersFootSteps, 11)
		playersDuckKill := make(playersDuckKills, 10)
		playersFlashedKill := make(playersFlashedKills, 10)
		playersAirborneKill := make(playersAirborneKills, 10)
		playersWallbangKill := make(playersWallbangKills, 10)
		playersSmokeKill := make(playerSmokeKills, 10)
		playersNoScopeKill := make(playerNoScopeKills, 10)
		playersWeaponShot := make(playerWeaponShots, 11)
		playersWeaponReload := make(playerWeaponReloads, 11)
		playersJump := make(playerJumps, 11)
		playersSmoke := make(playerSmokes, 11)
		playersHEGrenade := make(playerHEGrenades, 11)
		playersBombDrop := make(playerBombDrops, 11)
		playersFlashbang := make(playerFlashbangs, 11)
		playersIncendiaryGrenade := make(playerIncendiaryGrenades, 11)
		playersMolotov := make(playerMolotovs, 11)
		playersDecoyGrenade := make(playerDecoyGrenades, 11)

		p.RegisterEventHandler(func(e events.Footstep) {
			playersFootStep[e.Player.SteamID64] += 1
			playersFootStep[0] += 1
		})

		var (
			mapMetadata ex.Map
			mapRadarImg image.Image
		)

		p.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_ServerInfo) {
			// Get metadata for the map that the game was played on for coordinate translations
			mapMetadata = ex.GetMapMetadata(header.MapName, msg.GetMapCrc())
			// Load map overview image
			mapRadarImg = ex.GetMapRadar(header.MapName, msg.GetMapCrc())
		})

		var deathPoints []r2.Point
		p.RegisterEventHandler(func(e events.Kill) {
			if e.Victim != nil {
				x, y := mapMetadata.TranslateScale(e.Victim.Position().X, e.Victim.Position().Y)
				deathPoints = append(deathPoints, r2.Point{X: x, Y: y})
			}
			handleKill(e, playersDuckKill, playersFlashedKill, playersAirborneKill, playersWallbangKill, playersSmokeKill, playersNoScopeKill)
		})

		// Register handler for WeaponFire, triggered every time a shot is fired
		var firePoints []r2.Point
		p.RegisterEventHandler(func(e events.WeaponFire) {
			// Translate positions from in-game coordinates to radar overview image pixels
			x, y := mapMetadata.TranslateScale(e.Shooter.Position().X, e.Shooter.Position().Y)
			firePoints = append(firePoints, r2.Point{X: x, Y: y})
			playersWeaponShot[e.Shooter.SteamID64] += 1
			playersWeaponShot[0] += 1
		})

		p.RegisterEventHandler(func(e events.PlayerJump) {
			playersJump[e.Player.SteamID64] += 1
			playersJump[0] += 1
		})

		p.RegisterEventHandler(func(e events.WeaponReload) {
			playersWeaponReload[e.Player.SteamID64] += 1
			playersWeaponReload[0] += 1
		})

		var GrenadePoints []r2.Point
		p.RegisterEventHandler(func(e events.GrenadeProjectileThrow) {
			x, y := mapMetadata.TranslateScale(e.Projectile.Position().X, e.Projectile.Position().Y)
			GrenadePoints = append(GrenadePoints, r2.Point{X: x, Y: y})

			if e.Projectile.WeaponInstance.String() == "Smoke Grenade" {
				playersSmoke[e.Projectile.Thrower.SteamID64] += 1
				playersSmoke[0] += 1
			}

			if e.Projectile.WeaponInstance.String() == "HE Grenade" {
				playersHEGrenade[e.Projectile.Thrower.SteamID64] += 1
				playersHEGrenade[0] += 1
			}

			if e.Projectile.WeaponInstance.String() == "Flashbang" {
				playersFlashbang[e.Projectile.Thrower.SteamID64] += 1
				playersFlashbang[0] += 1
			}

			if e.Projectile.WeaponInstance.String() == "Incendiary Grenade" {
				playersIncendiaryGrenade[e.Projectile.Thrower.SteamID64] += 1
				playersIncendiaryGrenade[0] += 1
			}

			if e.Projectile.WeaponInstance.String() == "Molotov" {
				playersMolotov[e.Projectile.Thrower.SteamID64] += 1
				playersMolotov[0] += 1
			}

			if e.Projectile.WeaponInstance.String() == "Decoy Grenade" {
				playersDecoyGrenade[e.Projectile.Thrower.SteamID64] += 1
				playersDecoyGrenade[0] += 1
			}
		})

		p.RegisterEventHandler(func(e events.BombDropped) {
			playersBombDrop[e.Player.SteamID64] += 1
			playersBombDrop[0] += 1
		})

		p.RegisterEventHandler(func(e events.RoundEnd) {
			gs := p.GameState()
			switch e.Winner {
			case common.TeamTerrorists:
				// Winner's score + 1 because it hasn't actually been updated yet
				fmt.Printf("Round finished: winnerSide=T  ; score=%d:%d\n", gs.TeamTerrorists().Score()+1, gs.TeamCounterTerrorists().Score())
			case common.TeamCounterTerrorists:
				fmt.Printf("Round finished: winnerSide=CT ; score=%d:%d\n", gs.TeamCounterTerrorists().Score()+1, gs.TeamTerrorists().Score())
			default:
				// Probably match medic or something similar
				fmt.Println("Round finished: No winner (tie)")
			}
			// Copy nade paths
		})

		err = p.ParseToEnd()
		if err != nil {
			log.Panic("failed to parse demo: ", err)
		}

		players := p.GameState().Participants().Playing()
		var stats []playerStats
		for _, p := range players {
			stats = append(stats, statsFor(p, playersFootStep, playersDuckKill, playersFlashedKill, playersAirborneKill, playersWallbangKill, playersSmokeKill, playersNoScopeKill, playersWeaponShot, playersWeaponReload, playersJump, playersBombDrop, playersSmoke, playersHEGrenade, playersMolotov, playersIncendiaryGrenade, playersFlashbang, playersDecoyGrenade))
		}

		fmt.Println("Все игроки вместе сделали: ", playersFootStep[0], " шагов")
		fmt.Println("Все игроки вместе сделали: ", playersWeaponShot[0], " выстрелов")
		fmt.Println("Все игроки вместе сделали: ", playersWeaponReload[0], " перезарядок")
		fmt.Println("Все игроки вместе сделали: ", playersJump[0], " прыжков")
		fmt.Println("Все игроки в сумме дропнули бомбу: ", playersBombDrop[0], " раз(а)")
		fmt.Println("Все игроки в сумме кинули Smoke: ", playersSmoke[0], " раз(а)")
		fmt.Println("Все игроки в сумме кинули HE Grenade: ", playersHEGrenade[0], " раз(а)")
		fmt.Println("Все игроки в сумме кинули Molotov: ", playersMolotov[0], " раз(а)")
		fmt.Println("Все игроки в сумме кинули Incendiary Grenade: ", playersIncendiaryGrenade[0], " раз(а)")
		fmt.Println("Все игроки в сумме кинули Flashbang: ", playersFlashbang[0], " раз(а)")
		fmt.Println("Все игроки в сумме кинули Decoy: ", playersDecoyGrenade[0], " раз(а)")

		for _, player := range stats {
			fmt.Println(player.formatString() + "\n")
		}
		name, _ := strings.CutSuffix(e.Name(), ".dem")
		generateHeatMap(firePoints, mapRadarImg, name+".jpeg", "WeaponFire")
		generateHeatMap(deathPoints, mapRadarImg, name+".jpeg", "PlayerDeath")
		generateHeatMap(GrenadePoints, mapRadarImg, name+".jpeg", "GrenadeThrow")
		println("\n")
	}
}

func handleKill(e events.Kill, duckKills playersDuckKills, flashKills playersFlashedKills, airborneKills playersAirborneKills, wallbangKills playersWallbangKills, smokeKills playerSmokeKills, noScopeKills playerNoScopeKills) {
	if e.Killer != nil {
		if e.Killer.IsDucking() {
			duckKills[e.Killer.SteamID64] += 1
		}
	}

	if e.Killer != nil {
		if e.Killer.IsBlinded() {
			flashKills[e.Killer.SteamID64] += 1
		}
	}

	if e.Killer != nil {
		if e.Killer.IsAirborne() {
			airborneKills[e.Killer.SteamID64] += 1
		}
	}

	if e.IsWallBang() {
		wallbangKills[e.Killer.SteamID64] += 1
	}

	if e.ThroughSmoke {
		smokeKills[e.Killer.SteamID64] += 1
	}

	if e.NoScope {
		noScopeKills[e.Killer.SteamID64] += 1
	}

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

func statsFor(p *common.Player, fs playersFootSteps, dk playersDuckKills, fk playersFlashedKills, airk playersAirborneKills, wbk playersWallbangKills, sk playerSmokeKills, ns playerNoScopeKills, ws playerWeaponShots, wr playerWeaponReloads, pj playerJumps, bd playerBombDrops, smoke playerSmokes, he playerHEGrenades, molotov playerMolotovs, ctMoly playerIncendiaryGrenades, flashbang playerFlashbangs, decoy playerDecoyGrenades) playerStats {
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
		SmokeKills:    sk[p.SteamID64],
		NoScopeKills:  ns[p.SteamID64],
		WeaponShots:   ws[p.SteamID64],
		WeaponReloads: wr[p.SteamID64],
		PlayerJumps:   pj[p.SteamID64],
		BombDrops:     bd[p.SteamID64],
		Smokes:        smoke[p.SteamID64],
		HEnades:       he[p.SteamID64],
		Molotov:       molotov[p.SteamID64],
		CTmoly:        ctMoly[p.SteamID64],
		Flasbang:      flashbang[p.SteamID64],
		Decoy:         decoy[p.SteamID64],
	}
}

func generateHeatMap(points []r2.Point, mapRadarImg image.Image, name string, folder string) {
	r2Bounds := r2.RectFromPoints(points...)
	padding := float64(dotSize) / 2.0 // Calculating padding amount to avoid shrinkage by the heatmap library
	bounds := image.Rectangle{
		Min: image.Point{X: int(r2Bounds.X.Lo - padding), Y: int(r2Bounds.Y.Lo - padding)},
		Max: image.Point{X: int(r2Bounds.X.Hi + padding), Y: int(r2Bounds.Y.Hi + padding)},
	}

	// Transform r2.Points into heatmap.DataPoints
	data := make([]heatmap.DataPoint, 0, len(points))

	for _, p := range points[1:] {
		// Invert Y since go-heatmap expects data to be ordered from bottom to top
		data = append(data, heatmap.P(p.X, p.Y*-1))
	}

	// Create output canvas and use map overview image as base
	img := image.NewRGBA(mapRadarImg.Bounds())
	draw.Draw(img, mapRadarImg.Bounds(), mapRadarImg, image.Point{}, draw.Over)

	// Generate and draw heatmap overlay on top of the overview
	imgHeatmap := heatmap.Heatmap(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), data, dotSize, opacity, schemes.AlphaFire)
	draw.Draw(img, bounds, imgHeatmap, image.Point{}, draw.Over)
	f, err := os.Create("img/" + folder + "/" + name)
	// Write to stdout
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpegQuality})
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
