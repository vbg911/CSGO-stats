package demoparser

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/geo/r2"
	"github.com/llgcode/draw2d/draw2dimg"
	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/msg"
	"github.com/markus-wa/go-heatmap/v2"
	"github.com/markus-wa/go-heatmap/v2/schemes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"os"
	"strconv"
	"strings"
)
import "CSGO-stats/internal/structures"

const (
	dotSize = 20
	opacity = 128
)

func ParseDemo(tournamentName string, matchName string, filename string, demoPath string) (structures.MapStats, error) {
	demoName, _ := strings.CutSuffix(filename, ".dem")

	f, err := os.Open(demoPath)
	checkError(err)

	p := dem.NewParser(f)

	header, err := p.ParseHeader()
	checkError(err)

	playersFootStep := make(structures.PlayersFootSteps, 11)
	playersDuckKill := make(structures.PlayersDuckKills, 10)
	playersFlashedKill := make(structures.PlayersFlashedKills, 10)
	playersAirborneKill := make(structures.PlayersAirborneKills, 10)
	playersWallbangKill := make(structures.PlayersWallbangKills, 10)
	playersSmokeKill := make(structures.PlayerSmokeKills, 10)
	playersNoScopeKill := make(structures.PlayerNoScopeKills, 10)
	playersWeaponShot := make(structures.PlayerWeaponShots, 11)
	playersWeaponReload := make(structures.PlayerWeaponReloads, 11)
	playersJump := make(structures.PlayerJumps, 11)
	playersSmoke := make(structures.PlayerSmokes, 11)
	playersHEGrenade := make(structures.PlayerHEGrenades, 11)
	playersBombDrop := make(structures.PlayerBombDrops, 11)
	playersFlashbang := make(structures.PlayerFlashbangs, 11)
	playersIncendiaryGrenade := make(structures.PlayerIncendiaryGrenades, 11)
	playersMolotov := make(structures.PlayerMolotovs, 11)
	playersDecoyGrenade := make(structures.PlayerDecoyGrenades, 11)

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

		if e.Killer != nil {
			if e.Killer.IsDucking() {
				playersDuckKill[e.Killer.SteamID64] += 1
			}
		}

		if e.Killer != nil {
			if e.Killer.IsBlinded() {
				playersFlashedKill[e.Killer.SteamID64] += 1
			}
		}

		if e.Killer != nil {
			if e.Killer.IsAirborne() {
				playersAirborneKill[e.Killer.SteamID64] += 1
			}
		}

		if e.IsWallBang() {
			playersWallbangKill[e.Killer.SteamID64] += 1
		}

		if e.ThroughSmoke {
			playersSmokeKill[e.Killer.SteamID64] += 1
		}

		if e.NoScope {
			playersNoScopeKill[e.Killer.SteamID64] += 1
		}

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

	nadeTrajectories := make(map[int64]*common.GrenadeProjectile) // Trajectories of all destroyed nades
	p.RegisterEventHandler(func(e events.GrenadeProjectileDestroy) {
		id := e.Projectile.UniqueID()
		nadeTrajectories[id] = e.Projectile
	})

	infernos := make(map[int64]*common.Inferno)
	p.RegisterEventHandler(func(e events.InfernoExpired) {
		id := e.Inferno.UniqueID()
		infernos[id] = e.Inferno
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
		generateTrajectories(mapMetadata, mapRadarImg, nadeTrajectories, infernos, "GrenadeTrajectories\\random", demoName+" №"+strconv.Itoa(gs.TotalRoundsPlayed()+1)+".jpeg")

		for k := range nadeTrajectories {
			delete(nadeTrajectories, k)
		}
		for k := range infernos {
			delete(infernos, k)
		}
	})

	err = p.ParseToEnd()
	if err != nil {
		return structures.MapStats{}, err
	}

	players := p.GameState().Participants().Playing()
	var stats []structures.PlayerStats
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

	overallstats := make(map[string]int)
	overallstats["playersFootStep"] = playersFootStep[0]
	overallstats["playersWeaponShot"] = playersWeaponShot[0]
	overallstats["playersWeaponReload"] = playersWeaponReload[0]
	overallstats["playersJump"] = playersJump[0]
	overallstats["playersBombDrop"] = playersBombDrop[0]
	overallstats["playersSmoke"] = playersSmoke[0]
	overallstats["playersHEGrenade"] = playersHEGrenade[0]
	overallstats["playersMolotov"] = playersMolotov[0]
	overallstats["playersIncendiaryGrenade"] = playersIncendiaryGrenade[0]
	overallstats["playersFlashbang"] = playersFlashbang[0]
	overallstats["playersDecoyGrenade"] = playersDecoyGrenade[0]

	for _, player := range stats {
		fmt.Println(player.FormatString() + "\n")
	}
	generateHeatMap(firePoints, mapRadarImg, demoName+".jpeg", "WeaponFire")
	generateHeatMap(deathPoints, mapRadarImg, demoName+".jpeg", "PlayerDeath")
	generateHeatMap(GrenadePoints, mapRadarImg, demoName+".jpeg", "GrenadeThrow")
	println("\n")
	return structures.MapStats{
		TournamentName: tournamentName,
		DemoHash:       fileMD5(demoPath),
		DemoPath:       demoPath,
		MapName:        header.MapName,
		Players:        stats,
		FirePoints:     firePoints,
		DeathPoints:    deathPoints,
		GrenadePoints:  GrenadePoints,
		OverallStats:   overallstats,
	}, nil
}

func buildInfernoPath(mapMetadata ex.Map, gc *draw2dimg.GraphicContext, vertices []r2.Point) {
	xOrigin, yOrigin := mapMetadata.TranslateScale(vertices[0].X, vertices[0].Y)
	gc.MoveTo(xOrigin, yOrigin)

	for _, fire := range vertices[1:] {
		x, y := mapMetadata.TranslateScale(fire.X, fire.Y)
		gc.LineTo(x, y)
	}

	gc.LineTo(xOrigin, yOrigin)
}

func generateTrajectories(mapMetadata ex.Map, mapRadarImg image.Image, nadeMap map[int64]*common.GrenadeProjectile, infernos map[int64]*common.Inferno, folder string, name string) {
	var (
		colorFireNade    color.Color = color.RGBA{0xff, 0x00, 0x00, 0xff} // Red
		colorInferno     color.Color = color.RGBA{0xff, 0xa5, 0x00, 0xff} // Orange
		colorInfernoHull color.Color = color.RGBA{0xff, 0xff, 0x00, 0xff} // Yellow
		colorHE          color.Color = color.RGBA{0x00, 0xff, 0x00, 0xff} // Green
		colorFlash       color.Color = color.RGBA{0x00, 0x00, 0xff, 0xff} // Blue, because of the color on the nade
		colorSmoke       color.Color = color.RGBA{0xbe, 0xbe, 0xbe, 0xff} // Light gray
		colorDecoy       color.Color = color.RGBA{0x96, 0x4b, 0x00, 0xff} // Brown, because it's shit :)
	)

	// Create output canvas
	dest := image.NewRGBA(mapRadarImg.Bounds())

	// Draw image
	draw.Draw(dest, dest.Bounds(), mapRadarImg, image.Point{}, draw.Src)

	// Initialize the graphic context
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetFillColor(colorInferno)

	// Calculate hulls
	counter := 0
	hulls := make([][]r2.Point, len(infernos))
	for _, i := range infernos {
		hulls[counter] = i.Fires().ConvexHull2D()
		counter++
	}

	for _, hull := range hulls {
		buildInfernoPath(mapMetadata, gc, hull)
		gc.Fill()
	}

	// Then the outline
	gc.SetStrokeColor(colorInfernoHull)
	gc.SetLineWidth(1) // 1 px wide

	for _, hull := range hulls {
		buildInfernoPath(mapMetadata, gc, hull)
		gc.FillStroke()
	}

	gc.SetLineWidth(1)                      // 1 px lines
	gc.SetFillColor(color.RGBA{0, 0, 0, 0}) // No fill, alpha 0

	for _, nade := range nadeMap {
		// Set color
		switch nade.WeaponInstance.Type {
		case common.EqMolotov:
			fallthrough
		case common.EqIncendiary:
			gc.SetStrokeColor(colorFireNade)

		case common.EqHE:
			gc.SetStrokeColor(colorHE)

		case common.EqFlash:
			gc.SetStrokeColor(colorFlash)

		case common.EqSmoke:
			gc.SetStrokeColor(colorSmoke)

		case common.EqDecoy:
			gc.SetStrokeColor(colorDecoy)

		default:
			// Set alpha to 0 so we don't draw unknown stuff
			gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0x00})
		}

		// Draw path
		x, y := mapMetadata.TranslateScale(nade.Trajectory[0].X, nade.Trajectory[0].Y)
		gc.MoveTo(x, y) // Move to a position to start the new path

		for _, pos := range nade.Trajectory[1:] {
			x, y := mapMetadata.TranslateScale(pos.X, pos.Y)
			gc.LineTo(x, y)
		}
		gc.FillStroke()
	}
	err := os.Mkdir("img/"+folder, 0666)
	if err != nil && !os.IsExist(err) {
	}
	f, err := os.Create("img/" + folder + "/" + name)
	err = jpeg.Encode(f, dest, &jpeg.Options{
		Quality: 100,
	})
	checkError(err)
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
	err := os.Mkdir("img/"+folder, 0666)
	if err != nil && !os.IsExist(err) {
	}
	f, err := os.Create("img/" + folder + "/" + name)
	// Write to stdout
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	checkError(err)
}

func statsFor(p *common.Player, fs structures.PlayersFootSteps, dk structures.PlayersDuckKills, fk structures.PlayersFlashedKills, airk structures.PlayersAirborneKills, wbk structures.PlayersWallbangKills, sk structures.PlayerSmokeKills, ns structures.PlayerNoScopeKills, ws structures.PlayerWeaponShots, wr structures.PlayerWeaponReloads, pj structures.PlayerJumps, bd structures.PlayerBombDrops, smoke structures.PlayerSmokes, he structures.PlayerHEGrenades, molotov structures.PlayerMolotovs, ctMoly structures.PlayerIncendiaryGrenades, flashbang structures.PlayerFlashbangs, decoy structures.PlayerDecoyGrenades) structures.PlayerStats {
	return structures.PlayerStats{
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func fileMD5(path string) string {
	h := md5.New()
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(h, f)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
