package visualization

import (
	"CSGO-stats/internal/structures"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/golang/geo/r2"
	"github.com/llgcode/draw2d/draw2dimg"
	ex "github.com/markus-wa/demoinfocs-golang/v3/examples"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
	"github.com/markus-wa/go-heatmap/v2"
	"github.com/markus-wa/go-heatmap/v2/schemes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"os"
	"strconv"
)

const (
	dotSize = 20
	opacity = 128
)

func ESKills(data []structures.SummaryStatistics, name string, amount int, url string) *charts.EffectScatter {
	es := charts.NewEffectScatter()
	es.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Top 10 FRAGGERS",
			Subtitle: "Based on the analysis of " + strconv.Itoa(amount) + " demos from \"" + name + "\"",
			SubtitleStyle: &opts.TextStyle{
				FontSize: 15,
			},
			SubLink: url,
			Right:   "40%",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: name + " analysis",
			Width:     "1200px",
			Height:    "600px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "10%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Name:  "top 10 fraggers " + name,
					Title: "Download Chart",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "Show Data",
					Lang:  []string{"data view", "close", "refresh"},
				},
			}}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "87%"}),
		charts.WithXAxisOpts(opts.XAxis{
			Name:        "",
			Type:        "",
			Show:        true,
			Data:        nil,
			SplitNumber: 0,
			Scale:       false,
			Min:         nil,
			Max:         nil,
			MinInterval: 0,
			MaxInterval: 0,
			GridIndex:   0,
			SplitArea:   nil,
			SplitLine:   nil,
			AxisLabel: &opts.AxisLabel{
				Show:            true,
				Interval:        "0",
				Inside:          false,
				Rotate:          0,
				Margin:          0,
				Formatter:       "",
				ShowMinLabel:    true,
				ShowMaxLabel:    true,
				Color:           "",
				FontStyle:       "",
				FontWeight:      "",
				FontFamily:      "",
				FontSize:        "",
				Align:           "",
				VerticalAlign:   "",
				LineHeight:      "",
				BackgroundColor: "",
			},
			AxisTick:    nil,
			AxisPointer: nil,
		}),
	)

	var players []string
	var kills []opts.EffectScatterData
	for _, j := range data {
		players = append(players, j.Name)
		kills = append(kills, opts.EffectScatterData{
			Name:  j.Name,
			Value: j.Kills,
		})
	}

	es.SetXAxis(players[:10]).AddSeries("Kills", kills[:10], charts.WithRippleEffectOpts(opts.RippleEffect{
		Period:    4,
		Scale:     4,
		BrushType: "fill",
	}))

	return es
}

func ESDeath(data []structures.SummaryStatistics, name string, amount int, url string) *charts.EffectScatter {
	es := charts.NewEffectScatter()
	es.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Top 10 players with the most deaths",
			Subtitle: "Based on the analysis of " + strconv.Itoa(amount) + " demos from \"" + name + "\"",
			SubtitleStyle: &opts.TextStyle{
				FontSize: 15,
			},
			SubLink: url,
			Right:   "40%",
		}),
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: name + " analysis",
			Width:     "1200px",
			Height:    "600px",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "10%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Name:  "top 10 deaths " + name,
					Title: "Download Chart",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "Show Data",
					Lang:  []string{"data view", "close", "refresh"},
				},
			}}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Right: "87%"}),
		charts.WithXAxisOpts(opts.XAxis{
			Name:        "",
			Type:        "",
			Show:        true,
			Data:        nil,
			SplitNumber: 0,
			Scale:       false,
			Min:         nil,
			Max:         nil,
			MinInterval: 0,
			MaxInterval: 0,
			GridIndex:   0,
			SplitArea:   nil,
			SplitLine:   nil,
			AxisLabel: &opts.AxisLabel{
				Show:            true,
				Interval:        "0",
				Inside:          false,
				Rotate:          0,
				Margin:          0,
				Formatter:       "",
				ShowMinLabel:    true,
				ShowMaxLabel:    true,
				Color:           "",
				FontStyle:       "",
				FontWeight:      "",
				FontFamily:      "",
				FontSize:        "",
				Align:           "",
				VerticalAlign:   "",
				LineHeight:      "",
				BackgroundColor: "",
			},
			AxisTick:    nil,
			AxisPointer: nil,
		}),
	)

	var players []string
	var deaths []opts.EffectScatterData
	for _, j := range data {
		players = append(players, j.Name)
		deaths = append(deaths, opts.EffectScatterData{
			Name:  j.Name,
			Value: j.Deaths,
		})
	}

	es.SetXAxis(players[:10]).AddSeries("Deaths", deaths[:10], charts.WithRippleEffectOpts(opts.RippleEffect{
		Period:    4,
		Scale:     4,
		BrushType: "fill",
	}))

	return es
}

func GenerateCharts(data structures.ChartData, tournamentName string, demoAmount int, url string) {
	page := components.NewPage()
	page.PageTitle = tournamentName
	page.AddCharts(
		ESKills(data["SortedByKills"], tournamentName, demoAmount, url),
		ESDeath(data["SortedByDeath"], tournamentName, demoAmount, url),
	)
	f, err := os.Create("html/charts.html")
	if err != nil {
		panic(err)
	}
	err = page.Render(io.MultiWriter(f))
	if err != nil {
		panic(err)
	}
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

func GenerateTrajectories(mapMetadata ex.Map, mapRadarImg image.Image, matchNades structures.NadeTrajectories, infernos structures.Infernos, folder string, name string) {
	var (
		colorFireNade    color.Color = color.RGBA{0xff, 0x00, 0x00, 0xff} // Red
		colorInferno     color.Color = color.RGBA{0xff, 0xa5, 0x00, 0xff} // Orange
		colorInfernoHull color.Color = color.RGBA{0xff, 0xff, 0x00, 0xff} // Yellow
		colorHE          color.Color = color.RGBA{0x00, 0xff, 0x00, 0xff} // Green
		colorFlash       color.Color = color.RGBA{0x00, 0x00, 0xff, 0xff} // Blue, because of the color on the nade
		colorSmoke       color.Color = color.RGBA{0xbe, 0xbe, 0xbe, 0xff} // Light gray
		colorDecoy       color.Color = color.RGBA{0x96, 0x4b, 0x00, 0xff} // Brown, because it's shit :)
	)

	for i, round := range matchNades {
		// Create output canvas
		dest := image.NewRGBA(mapRadarImg.Bounds())
		// Draw image
		draw.Draw(dest, dest.Bounds(), mapRadarImg, image.Point{}, draw.Src)
		// Initialize the graphic context
		gc := draw2dimg.NewGraphicContext(dest)

		gc.SetFillColor(colorInferno)

		// Calculate hulls
		hulls := make([][]r2.Point, len(infernos[i]))
		counter := 0
		for _, moly := range infernos[i] {
			hulls[counter] = moly.Fires().ConvexHull2D()
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

		gc.SetLineWidth(1.5)                    // 1 px lines
		gc.SetFillColor(color.RGBA{0, 0, 0, 0}) // No fill, alpha 0

		for _, nade := range round {
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
				gc.SetStrokeColor(color.RGBA{0x00, 0x00, 0x00, 0x00})
				fmt.Println("Unknown grenade type")
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
		f, err := os.Create("img/" + folder + "/" + strconv.Itoa(i) + "-" + name)
		err = jpeg.Encode(f, dest, &jpeg.Options{
			Quality: 100,
		})
		gc.Clear()
		checkError(err)
	}
}

func GenerateHeatMap(points []r2.Point, mapRadarImg image.Image, name string, folder string) {
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
