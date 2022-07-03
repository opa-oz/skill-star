package main

import (
	"github.com/opa-oz/skill-star"
	"image/color"
	"image/png"
	"os"
)

func main() {
	cfg := skill_star.SkillsConfig{
		Skills: []string{"Strength", "Speed", "Health", "Regeneration", "Hydration"},
		Depth:  5,
	}

	person := skill_star.Person{
		Name:         "Vladimir",
		SkillsValues: []int{3, 2, 1, 4, 3},
	}

	imageCfg := skill_star.ImageConfig{
		Width:  600,
		Height: 600,

		NeedName: false,

		BackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 0xff},
		TextColor:       color.RGBA{R: 0, G: 0, B: 0, A: 0xff},
		StrokeColor:     color.RGBA{R: 0, G: 0, B: 0, A: 0xff},
		PersonColor:     color.RGBA{R: 255, G: 192, B: 203, A: 0xff},

		Radius: 250,
	}

	img := skill_star.GenerateSkillStar(cfg, imageCfg, person)

	f, _ := os.Create("image.png")
	err := png.Encode(f, &img)
	if err != nil {
		return
	}
}
