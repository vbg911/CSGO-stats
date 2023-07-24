package main

import (
	"fmt"
	"github.com/golang/geo/r2"
	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/msg"
	heatmap "github.com/markus-wa/go-heatmap/v2"
	schemes "github.com/markus-wa/go-heatmap/v2/schemes"
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

		// Register handler for WeaponFire, triggered every time a shot is fired
		var points []r2.Point

		p.RegisterEventHandler(func(e events.WeaponFire) {
			// Translate positions from in-game coordinates to radar overview image pixels
			x, y := mapMetadata.TranslateScale(e.Shooter.Position().X, e.Shooter.Position().Y)

			points = append(points, r2.Point{X: x, Y: y})
		})

		p.RegisterEventHandler(func(events.DataTablesParsed) {

		})

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
		name, _ := strings.CutSuffix(e.Name(), ".dem")
		generateHeatMap(points, mapRadarImg, name+".jpeg")
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

func generateHeatMap(points []r2.Point, mapRadarImg image.Image, name string) {
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

	//
	// Drawing the image
	//

	// Create output canvas and use map overview image as base
	img := image.NewRGBA(mapRadarImg.Bounds())
	draw.Draw(img, mapRadarImg.Bounds(), mapRadarImg, image.Point{}, draw.Over)

	// Generate and draw heatmap overlay on top of the overview
	imgHeatmap := heatmap.Heatmap(image.Rect(0, 0, bounds.Dx(), bounds.Dy()), data, dotSize, opacity, schemes.AlphaFire)
	draw.Draw(img, bounds, imgHeatmap, image.Point{}, draw.Over)
	f, err := os.Create("img/AlphaFire_" + name)
	// Write to stdout
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpegQuality})
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
