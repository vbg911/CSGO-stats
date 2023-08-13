package visualization

import (
	"CSGO-stats/internal/structures"
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
	"os"
	"strconv"
)

const (
	dotSize = 20
	opacity = 128
)

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

	// Create output canvas
	dest := image.NewRGBA(mapRadarImg.Bounds())

	// Draw image
	draw.Draw(dest, dest.Bounds(), mapRadarImg, image.Point{}, draw.Src)

	// Initialize the graphic context
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetFillColor(colorInferno)

	// Calculate hulls

	for i, round := range matchNades {
		// Set color

		counter := 0
		hulls := make([][]r2.Point, len(infernos))
		for _, fires := range infernos {
			for _, fire := range fires {
				hulls[counter] = fire.Fires().ConvexHull2D()
				counter++
			}
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
		f, err := os.Create("img/" + folder + "/" + strconv.Itoa(i) + name)
		err = jpeg.Encode(f, dest, &jpeg.Options{
			Quality: 100,
		})
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
