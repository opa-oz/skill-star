package skill_star

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/teacat/noire"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type Person struct {
	Name         string
	SkillsValues []int
}

type SkillsConfig struct {
	Skills []string // Skill's names. (ex: ["Strength", "Stamina", "Health", "Speed"])
	Depth  int      // Maximal available skill value (starting with 0)
}

type ImageConfig struct {
	Width    int
	Height   int
	NeedName bool

	BackgroundColor color.RGBA
	TextColor       color.RGBA
	PersonColor     color.RGBA
	StrokeColor     color.RGBA

	Radius int
}

// LinSpace returns evenly spaced numbers over a specified closed interval.
// @see {@link https://github.com/cpmech/gosl/blob/2f4609fa0c595209c66a30337acfb628e923ded2/utl/mylab.go}
func linSpace(start, stop float64, num int) (res []float64) {
	if num <= 0 {
		return []float64{}
	}
	if num == 1 {
		return []float64{start}
	}
	step := (stop - start) / float64(num-1)
	res = make([]float64, num)
	res[0] = start
	for i := 1; i < num; i++ {
		res[i] = start + float64(i)*step
	}
	res[num-1] = stop
	return
}

type poorPoint struct {
	x float64
	y float64
}

func makeCirclePoints(radius int, amount int, offsetPoint poorPoint) []poorPoint {
	points := make([]poorPoint, amount)
	intervals := linSpace(0, 2*math.Pi, amount)

	floatRadius := float64(radius)

	for i, item := range intervals {
		x := floatRadius * math.Sin(item)
		y := floatRadius * math.Cos(item)

		points[i] = poorPoint{x: offsetPoint.x + x, y: offsetPoint.y + y}
	}

	return points
}

// @see {@link https://github.com/teacat/noire}
func lightenDarkenColor(clr color.RGBA, amount float64) color.RGBA {
	r, g, b, a := clr.RGBA()

	c := noire.NewRGB(float64(r/255), float64(g/255), float64(b/255))
	newR, newG, newB := c.Lighten(amount).RGB()

	fmt.Println(r, g, b)
	fmt.Println(newR, newG, newB)

	return color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a)}
}

func drawCanvas(width int, height int, backgroundColor color.Color) *image.RGBA {
	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	draw.Draw(img, img.Bounds(), &image.Uniform{C: backgroundColor}, image.Point{}, draw.Src)

	return img
}

func GenerateSkillStar(cfg SkillsConfig, imageCfg ImageConfig, person Person) image.RGBA {
	img := drawCanvas(imageCfg.Width, imageCfg.Height, imageCfg.BackgroundColor)

	gc := draw2dimg.NewGraphicContext(img)

	gc.SetStrokeColor(imageCfg.StrokeColor)
	gc.SetLineWidth(2)

	middleX := float64(imageCfg.Width) / 2.0
	middleY := float64(imageCfg.Height) / 2.0

	skillsCount := len(cfg.Skills)
	radius := imageCfg.Radius

	depth := cfg.Depth + 1

	middlePoint := poorPoint{x: middleX, y: middleY}

	skillsOrbits := linSpace(0, float64(radius), depth)

	mainPoints := make([][]poorPoint, len(skillsOrbits))

	for i := len(skillsOrbits) - 1; i >= 1; i-- {
		points := makeCirclePoints(int(skillsOrbits[i]), skillsCount+1, middlePoint)

		mainPoints[i] = points

		gc.BeginPath()
		gc.MoveTo(points[0].x, points[0].y)

		for i := 1; i < len(points); i++ {
			gc.LineTo(points[i].x, points[i].y)
		}

		gc.Close()
		gc.Stroke()
	}

	for i := 0; i < len(mainPoints[0]); i++ {
		gc.BeginPath()
		gc.MoveTo(middleX, middleY)
		gc.LineTo(mainPoints[0][i].x, mainPoints[0][i].y)
		gc.Stroke()
	}

	shapeImg := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: imageCfg.Width, Y: imageCfg.Height}})

	gcc := draw2dimg.NewGraphicContext(shapeImg)

	lightenColor := lightenDarkenColor(imageCfg.PersonColor, 0.02)

	gcc.SetStrokeColor(imageCfg.PersonColor)
	gcc.SetFillColor(lightenColor)
	gcc.SetLineWidth(4)

	gcc.BeginPath()

	var shape []poorPoint

	for i, value := range person.SkillsValues {
		point := mainPoints[value][i]

		shape = append(shape, point)

		if i == 0 {
			gcc.MoveTo(point.x, point.y)
		} else {
			gcc.LineTo(point.x, point.y)
		}
	}

	gcc.Close()
	gcc.Stroke()

	maskImage := drawCanvas(imageCfg.Width, imageCfg.Height, color.RGBA{R: 10, G: 10, B: 10, A: 0xff})

	bounds := img.Bounds() //you have defined that both src and mask are same size, and maskImg is a grayscale of the src image. So we'll use that common size.
	mask := image.NewAlpha(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			//get one of r, g, b on the mask image ...
			r, _, _, _ := maskImage.At(x, y).RGBA()
			//... and set it as the alpha value on the mask.
			mask.SetAlpha(x, y, color.Alpha{A: uint8(255 - r)}) //Assuming that white is your transparency, subtract it from 255
		}
	}

	m := image.NewRGBA(bounds)
	draw.Draw(m, m.Bounds(), img, image.Point{}, draw.Src)

	draw.DrawMask(m, bounds, shapeImg, image.Point{}, mask, image.Point{}, draw.Over)

	shapeImg2 := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: imageCfg.Width, Y: imageCfg.Height}})
	gcc2 := draw2dimg.NewGraphicContext(shapeImg2)

	gcc2.SetFillColor(lightenColor)
	gcc2.SetStrokeColor(lightenColor)

	gcc2.BeginPath()

	for i, point := range shape {
		if i == 0 {
			gcc2.MoveTo(point.x, point.y)
		} else {
			gcc2.LineTo(point.x, point.y)
		}
	}

	gcc2.Close()
	gcc2.FillStroke()

	maskImage2 := drawCanvas(imageCfg.Width, imageCfg.Height, color.RGBA{R: 120, G: 120, B: 120, A: 0xff})

	mask2 := image.NewAlpha(bounds)
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			//get one of r, g, b on the mask image ...
			r, _, _, _ := maskImage2.At(x, y).RGBA()
			//... and set it as the alpha value on the mask.
			mask2.SetAlpha(x, y, color.Alpha{A: uint8(255 - r)}) //Assuming that white is your transparency, subtract it from 255
		}
	}

	m2 := image.NewRGBA(bounds)
	draw.Draw(m2, m2.Bounds(), m, image.Point{}, draw.Src)
	draw.DrawMask(m2, bounds, shapeImg2, image.Point{}, mask2, image.Point{}, draw.Over)

	return *m2
}
