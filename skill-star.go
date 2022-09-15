package skill_star

import (
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

	//fmt.Println(r, g, b)
	//fmt.Println(newR, newG, newB)

	return color.RGBA{R: uint8(newR), G: uint8(newG), B: uint8(newB), A: uint8(a)}
}

func drawCanvas(width int, height int, backgroundColor color.Color) *image.RGBA {
	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	draw.Draw(img, img.Bounds(), &image.Uniform{C: backgroundColor}, image.Point{}, draw.Src)

	return img
}

type Mask struct {
	bounds image.Rectangle
	width  int
	height int
}

func (c Mask) Prepare() (*draw2dimg.GraphicContext, image.RGBA) {
	shapeImg := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: c.width, Y: c.height}})
	gc := draw2dimg.NewGraphicContext(shapeImg)

	return gc, *shapeImg
}

func (c Mask) Draw(targetImage *image.RGBA, gc *draw2dimg.GraphicContext, img *image.RGBA, transparency uint8) *image.RGBA {
	maskImage := drawCanvas(c.width, c.height, color.RGBA{R: transparency, G: transparency, B: transparency, A: 0xff})

	bounds := c.bounds
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
	draw.Draw(m, m.Bounds(), targetImage, image.Point{}, draw.Src)
	draw.DrawMask(m, bounds, img, image.Point{}, mask, image.Point{}, draw.Over)

	return m
}

func drawSkillsShape(img *image.RGBA, shape []poorPoint, lightenColor color.Color, imageCfg ImageConfig) *image.RGBA {
	mask := Mask{bounds: img.Bounds(), width: imageCfg.Width, height: imageCfg.Height}
	gcc, shapeImg := mask.Prepare()

	gcc.SetStrokeColor(imageCfg.PersonColor)
	gcc.SetFillColor(lightenColor)
	gcc.SetLineWidth(4)

	gcc.BeginPath()

	for i, point := range shape {
		if i == 0 {
			gcc.MoveTo(point.x, point.y)
		} else {
			gcc.LineTo(point.x, point.y)
		}
	}

	gcc.Close()
	gcc.Stroke()

	img = mask.Draw(img, gcc, &shapeImg, 10)

	return img
}

func drawSkillsStroke(img *image.RGBA, shape []poorPoint, lightenColor color.Color, imageCfg ImageConfig) *image.RGBA {
	mask := Mask{bounds: img.Bounds(), width: imageCfg.Width, height: imageCfg.Height}
	gcc, shapeImg := mask.Prepare()

	gcc.SetFillColor(lightenColor)
	gcc.SetStrokeColor(lightenColor)

	gcc.BeginPath()

	for i, point := range shape {
		if i == 0 {
			gcc.MoveTo(point.x, point.y)
		} else {
			gcc.LineTo(point.x, point.y)
		}
	}

	gcc.Close()
	gcc.FillStroke()

	img = mask.Draw(img, gcc, &shapeImg, 120)
	return img
}

func GenerateSkillStar(cfg SkillsConfig, imageCfg ImageConfig, person Person) image.Image {
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

	target := len(mainPoints) - 1

	for i := 0; i < len(mainPoints[target]); i++ {
		gc.BeginPath()
		gc.MoveTo(middleX, middleY)
		gc.LineTo(mainPoints[target][i].x, mainPoints[target][i].y)
		gc.Stroke()
	}

	var shape []poorPoint

	for i, value := range person.SkillsValues {
		point := mainPoints[value][i]
		shape = append(shape, point)
	}

	lightenColor := lightenDarkenColor(imageCfg.PersonColor, 0.02)

	img = drawSkillsShape(img, shape, lightenColor, imageCfg)
	img = drawSkillsStroke(img, shape, lightenColor, imageCfg)

	return img
}
