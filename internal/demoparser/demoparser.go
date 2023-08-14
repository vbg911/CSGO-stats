package demoparser

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/geo/r2"
	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/msg"
	"image"
	"io"
	"os"
	"strings"
)
import "CSGO-stats/internal/structures"

func ParseDemo(tournamentName string, matchName string, filename string, demoPath string) (structures.MapStats, error) {
	demoName, _ := strings.CutSuffix(filename, ".dem")

	f, err := os.Open(demoPath)
	checkError(err)

	p := dem.NewParser(f)

	header, err := p.ParseHeader()
	checkError(err)

	playersFootStep := make(structures.FootSteps, 11)
	playersDuckKill := make(structures.DuckKills, 10)
	playersFlashedKill := make(structures.FlashedKills, 10)
	playersAirborneKill := make(structures.AirborneKills, 10)
	playersWallbangKill := make(structures.WallbangKills, 10)
	playersSmokeKill := make(structures.SmokeKills, 10)
	playersNoScopeKill := make(structures.NoScopeKills, 10)
	playersWeaponShot := make(structures.WeaponShots, 11)
	playersWeaponReload := make(structures.WeaponReloads, 11)
	playersJump := make(structures.Jumps, 11)
	playersSmoke := make(structures.Smokes, 11)
	playersHEGrenade := make(structures.HEGrenades, 11)
	playersBombDrop := make(structures.BombDrops, 11)
	playersFlashbang := make(structures.Flashbangs, 11)
	playersIncendiaryGrenade := make(structures.IncendiaryGrenades, 11)
	playersMolotov := make(structures.Molotovs, 11)
	playersDecoyGrenade := make(structures.DecoyGrenades, 11)

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

	nadeTrajectories := make(structures.NadeTrajectories) // Trajectories of all destroyed nades
	p.RegisterEventHandler(func(e events.GrenadeProjectileDestroy) {
		gs := p.GameState()
		id := e.Projectile.UniqueID()
		nadeTrajectories[gs.TotalRoundsPlayed()] = make(map[int64]*common.GrenadeProjectile)
		nadeTrajectories[gs.TotalRoundsPlayed()][id] = e.Projectile
	})

	infernos := make(structures.Infernos)
	p.RegisterEventHandler(func(e events.InfernoExpired) {
		gs := p.GameState()
		id := e.Inferno.UniqueID()
		infernos[gs.TotalRoundsPlayed()] = make(map[int64]*common.Inferno)
		infernos[gs.TotalRoundsPlayed()][id] = e.Inferno
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

	return structures.MapStats{
		TournamentName:           tournamentName,
		MatchName:                matchName,
		DemoName:                 demoName,
		DemoHash:                 fileMD5(demoPath),
		DemoPath:                 demoPath,
		MapMetadata:              mapMetadata,
		MapRadarImg:              mapRadarImg,
		MapName:                  header.MapName,
		Players:                  stats,
		FirePoints:               firePoints,
		DeathPoints:              deathPoints,
		GrenadePoints:            GrenadePoints,
		NadesProjectiles:         nadeTrajectories,
		NadesInferno:             infernos,
		PlayersFootStep:          playersFootStep,
		PlayersWeaponShot:        playersWeaponShot,
		PlayersWeaponReload:      playersWeaponReload,
		PlayersJump:              playersJump,
		PlayersBombDrop:          playersBombDrop,
		PlayersSmoke:             playersSmoke,
		PlayersHEGrenade:         playersHEGrenade,
		PlayersMolotov:           playersMolotov,
		PlayersIncendiaryGrenade: playersIncendiaryGrenade,
		PlayersFlashbang:         playersFlashbang,
		PlayersDecoyGrenade:      playersDecoyGrenade,
	}, nil
}

func statsFor(p *common.Player, fs structures.FootSteps, dk structures.DuckKills, fk structures.FlashedKills, airk structures.AirborneKills, wbk structures.WallbangKills, sk structures.SmokeKills, ns structures.NoScopeKills, ws structures.WeaponShots, wr structures.WeaponReloads, pj structures.Jumps, bd structures.BombDrops, smoke structures.Smokes, he structures.HEGrenades, molotov structures.Molotovs, ctMoly structures.IncendiaryGrenades, flashbang structures.Flashbangs, decoy structures.DecoyGrenades) structures.PlayerStats {
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
