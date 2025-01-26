package color

import (
	"fmt"
	"github.com/deanishe/awgo"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	ic "image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Color struct {
	values     []string
	optCommand string
}

func NewColor(values []string) *Color {
	cleanedValues := make([]string, 0, len(values))
	var optCommand string

	for _, v := range values {
		if optCommand == "" {
			if v == "hsl" || v == "hsv" {
				optCommand = v
				continue
			}
		}
		cleanedValues = append(cleanedValues, strings.Trim(strings.TrimSpace(v), ","))
	}

	return &Color{values: cleanedValues, optCommand: optCommand}
}

func (c *Color) Key() string {
	return "c"
}

func (c *Color) Do(wf *aw.Workflow) {
	log.Printf("[Color] \tprocessing..., values: %v\n", c.values)

	var color colorful.Color
	var err error

	if len(c.values) == 0 {
		c.setResult(wf, colorful.FastHappyColor())
		return
	}

	for _, value := range c.values {
		if value == "" {
			c.setResult(wf, colorful.FastHappyColor())
			return
		}
	}

	switch {
	case len(c.values) == 1 && c.values[0][0] == '#':
		color, err = colorful.Hex(c.values[0])
		if err != nil {
			c.setInvalidValue(wf)
			return
		}
	case len(c.values) == 3:
		// Try to parse as RGB
		r, err1 := strconv.ParseUint(c.values[0], 10, 8)
		g, err2 := strconv.ParseUint(c.values[1], 10, 8)
		b, err3 := strconv.ParseUint(c.values[2], 10, 8)
		if err1 == nil && err2 == nil && err3 == nil {
			color = colorful.Color{R: float64(r) / 255, G: float64(g) / 255, B: float64(b) / 255}
		} else {
			c.setInvalidValue(wf)
		}
	default:
		c.setInvalidValue(wf)
		return
	}

	c.setResult(wf, color)
}

func (c *Color) setResult(wf *aw.Workflow, color colorful.Color) {
	typeToValue := make(map[string]string)

	icon := c.createIcon(wf, color)
	if c.optCommand == "hsl" {
		h, s, l := color.Hsl()
		typeToValue["hsl"] = fmt.Sprintf("%.2f, %.2f, %.2f", h, s, l)
	} else if c.optCommand == "hsv" {
		h, s, v := color.Hsv()
		typeToValue["hsv"] = fmt.Sprintf("%.2f, %.2f, %.2f", h, s, v)
	} else {
		typeToValue["hex"] = color.Hex()

		r, g, b := color.RGB255()
		typeToValue["rgb"] = fmt.Sprintf("%d, %d, %d", r, g, b)
	}

	for tp, value := range typeToValue {
		wf.NewItem(value).
			Subtitle(fmt.Sprintf("Get color in %s", tp)).
			Arg(value).
			Copytext(value).
			Quicklook(value).
			Valid(true).
			Autocomplete(fmt.Sprintf("%s %s ", c.Key(), value)).
			Icon(&icon)
	}
}

func (c *Color) createIcon(wf *aw.Workflow, color colorful.Color) aw.Icon {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	r, g, b := color.RGB255()
	fillColor := ic.RGBA{R: r, G: g, B: b, A: 255}
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			img.Set(x, y, fillColor)
		}
	}

	// Create a new image with shadow
	shadowImg := image.NewRGBA(image.Rect(0, 0, 80, 80))
	draw.Draw(shadowImg, shadowImg.Bounds(), &image.Uniform{C: ic.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)

	// Draw shadow with blur effect
	gc := draw2dimg.NewGraphicContext(shadowImg)
	gc.SetFillColor(ic.RGBA{0, 0, 0, 128})
	draw2dkit.RoundedRectangle(gc, 8, 8, 72, 72, 10, 10)
	gc.Fill()

	// Draw the original image on top of the shadow
	draw.Draw(shadowImg, image.Rect(8, 8, 72, 72), img, image.Point{}, draw.Over)

	// Get the data directory path
	dataDir := wf.DataDir()
	iconPath := filepath.Join(dataDir, fmt.Sprintf("%s.png", color.Hex()))

	// Create a file to save the image
	file, err := os.Create(iconPath)
	if err != nil {
		log.Printf("Failed to create icon file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close icon file: %v", err)
		}
	}(file)

	// Encode the image to PNG format and save it to the file
	if err := png.Encode(file, shadowImg); err != nil {
		log.Printf("Failed to encode image to PNG: %v", err)
	}

	log.Printf("Icon saved to %s", iconPath)

	return aw.Icon{Value: iconPath, Type: aw.IconTypeImage}
}
func (c *Color) setInvalidValue(wf *aw.Workflow) {
	wf.NewItem(fmt.Sprintf("Invalid value: %v", c.values)).
		Valid(false)
}

//
//func (c *Color) Key() string {
//	return "c"
//}
//
//func (c *Color) Do(wf *aw.Workflow) {
//	log.Printf("[Color] \tprocessing..., values: %v\n", c.values)
//
//	color, err := colorful.Hex(c.values[0])
//	if err != nil {
//		color.RGB255()
//		return
//	}
//	if len(c.values) != 3 {
//		return
//	}
//
//	first := strconv.ParseFloat(c.values[0], 64)
//	second := strconv.ParseFloat(c.values[1], 64)
//	third := strconv.ParseFloat(c.values[2], 64)
//
//	colorful.Hcl(first, second, thrid)
//	colorful.Hsv(first, second, thrid)
//	colorful.LinearRgb(first, second, thrid)
//
//	result := "a"
//
//	wf.NewItem(result).
//		Subtitle("Generate a universal unique identifier (UUID).").
//		Arg(result).
//		Copytext(result).
//		Quicklook(result).
//		Valid(true).
//		Autocomplete(u.Key())
//}
