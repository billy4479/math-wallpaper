package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"

	catppuccin "github.com/catppuccin/go"
	"github.com/gen2brain/go-fitz"
	"golang.org/x/image/draw"
)

func processImage(img *image.RGBA, col color.Color) {
	r, g, b, _ := col.RGBA()
	for x := range img.Rect.Size().X {
		for y := range img.Rect.Size().Y {
			originalColor := img.RGBAAt(x, y)
			a := 0x00ff - uint16(originalColor.R)
			if a > 0x0100 {
				a = 0
			}
			finalColor := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
			img.Set(x, y, finalColor)
		}
	}
}

var colors = []color.Color{
	catppuccin.Frappe.Rosewater(),
	catppuccin.Frappe.Flamingo(),
	catppuccin.Frappe.Pink(),
	catppuccin.Frappe.Mauve(),
	catppuccin.Frappe.Red(),
	catppuccin.Frappe.Maroon(),
	catppuccin.Frappe.Peach(),
	catppuccin.Frappe.Yellow(),
	catppuccin.Frappe.Green(),
	catppuccin.Frappe.Teal(),
	catppuccin.Frappe.Sky(),
	catppuccin.Frappe.Sapphire(),
	catppuccin.Frappe.Blue(),
	catppuccin.Frappe.Lavender(),
	catppuccin.Frappe.Text(),
}

func getImgs(n int) []*image.RGBA {
	pdfs, err := os.ReadDir("build")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result := []*image.RGBA{}

	for range n {
		pdf := pdfs[rand.Intn(len(pdfs))].Name()
		doc, err := fitz.New("build/" + pdf)
		if err != nil {
			fmt.Println(pdf, err)
			continue
		}

		img, err := doc.ImageDPI(0, 900)
		if err != nil {
			fmt.Println(pdf, err)
			continue
		}

		processImage(img, colors[rand.Intn(len(colors))])
		result = append(result, img)
		fmt.Println(pdf)
	}

	return result
}

func render() {
	formulas := getImgs(400)

	w := 2560.
	h := 1440.
	maxSize := 0.4
	minSize := 0.075
	maxR := max(w, h) / 2.5

	radial := true

	out := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	// https://github.com/catppuccin/go/issues/29
	r, g, b, _ := catppuccin.Frappe.Crust().RGBA()
	bg := image.NewUniform(color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 0xff})

	draw.Copy(out, image.Point{}, bg, out.Bounds(), draw.Src, nil)

	for _, formula := range formulas {
		var centerX, centerY, r float64
		if radial {
			angle := float64(rand.Intn(360))
			r = float64(rand.Intn(int(maxR)))

			centerX = math.Cos(angle)*r + w/2
			centerY = math.Sin(angle)*r + h/2

		} else {
			centerX = float64(rand.Intn(int(w)))
			centerY = float64(rand.Intn(int(h)))

			r = math.Sqrt(math.Pow(w/2-centerX, 2) + math.Pow(h/2-centerY, 2))
		}

		scale := (maxSize-minSize)*math.Pow(r/float64(maxR), 1.5) + minSize

		originalSize := formula.Bounds().Size()
		originalW := float64(originalSize.X)
		originalH := float64(originalSize.Y)

		newW := originalW * scale
		newH := originalH * scale
		x := centerX - float64(newW)/2
		y := centerY - float64(newH)/2

		dstRect := image.Rect(
			int(x),
			int(y),
			int(x+newW),
			int(y+newH),
		)

		draw.BiLinear.Scale(out, dstRect, formula, formula.Bounds(), draw.Over, nil)
	}

	f, err := os.Create("out.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer f.Close()

	err = png.Encode(f, out)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
