package skill_star

import (
	gI "github.com/opa-oz/golden-image"
	"image/color"
	"testing"
)

var imageCfg = ImageConfig{
	Width:  600,
	Height: 600,

	NeedName: false,

	BackgroundColor: color.RGBA{R: 255, G: 255, B: 255, A: 0xff},
	TextColor:       color.RGBA{R: 0, G: 0, B: 0, A: 0xff},
	StrokeColor:     color.RGBA{R: 0, G: 0, B: 0, A: 0xff},
	PersonColor:     color.RGBA{R: 255, G: 192, B: 203, A: 0xff},

	Radius: 250,
}

const treshold = 0.02

func TestGenerateBaseSkillStar(t *testing.T) {
	cfg := SkillsConfig{
		Skills: []string{"Strength", "Speed", "Health", "Regeneration", "Hydration"},
		Depth:  5,
	}

	person := Person{
		Name:         "Vladimir",
		SkillsValues: []int{3, 2, 1, 4, 3},
	}

	skillStar := GenerateSkillStar(cfg, imageCfg, person)

	gI.ToGildImage(t, treshold, skillStar)
}

func TestGenerateSmallestSkillStar(t *testing.T) {
	cfg := SkillsConfig{
		Skills: []string{"Fish", "Human", "Cow"},
		Depth:  2,
	}

	person := Person{
		Name:         "Beaver",
		SkillsValues: []int{1, 1, 2},
	}

	skillStar := GenerateSkillStar(cfg, imageCfg, person)

	gI.ToGildImage(t, treshold, skillStar)
}

func TestGenerateOctoSkillStar(t *testing.T) {
	cfg := SkillsConfig{
		Skills: []string{"One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight"},
		Depth:  8,
	}

	person := Person{
		Name:         "Beaver",
		SkillsValues: []int{1, 2, 3, 4, 5, 6, 7, 8},
	}

	skillStar := GenerateSkillStar(cfg, imageCfg, person)

	gI.ToGildImage(t, treshold, skillStar)
}
