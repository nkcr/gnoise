package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
)

// processing function that render points
type proc func(w io.Writer, width, height, startRadius float64, ease ease)

// ease function to parametrize the radius evolution
type ease func(float64) float64

// main options, would be a good idea to set them as flags
type opts struct {
	width  float64
	height float64
	radius float64
	fill   string
	output string
	proc   proc
	ease   ease
}

func main() {
	opts := opts{
		width:  300.0,
		height: 300.0,
		radius: 3.0,
		fill:   "#44FF44",
		output: "result.svg",
		proc:   gradientCompact,
		ease:   easeLinear,
	}

	render(opts)
}

func render(opts opts) {
	w := &bytes.Buffer{}

	fmt.Fprintf(w, `<svg viewBox="0 0 %.1f %f.1" xmlns="http://www.w3.org/2000/svg">`, opts.width, opts.height)
	fmt.Fprintf(w, `<g fill="%s">`, opts.fill)
	opts.proc(w, opts.width, opts.height, opts.radius, opts.ease)
	fmt.Fprint(w, "</g>")
	fmt.Fprint(w, `</svg>`)

	err := os.WriteFile(opts.output, w.Bytes(), 0744)
	if err != nil {
		panic(err)
	}
}

// gradientCompact stacks circles leaving the smallest gap possible.
func gradientCompact(w io.Writer, width, height, startRadius float64, ease ease) {
	// linHeight is the height of stack of 3 congruent semicircles
	lineHeight := startRadius * math.Sqrt(3)

	numLines := height / lineHeight
	numCols := width / startRadius

	// make the radius a bit bigger to have a nice overlay
	radiusFactor := 1.3

	for y := 0; y < int(numLines); y++ {
		for x := 0; x < int(numCols); x++ {
			xPos, yPos := float64(x)*startRadius*2, float64(y)*lineHeight

			if y%2 == 0 {
				xPos += startRadius
			}

			// we go from [0,startRadius]

			radius := (startRadius * (float64(y) / numLines))
			coeff := ease(radius / (startRadius))
			radius *= coeff
			radius *= radiusFactor

			fmt.Fprintf(w, "<circle cx=\"%.2f\" cy=\"%.2f\" r=\"%.2f\"/>", xPos, yPos, radius)
		}
	}
}

// gradientSquared puts each circle on a square, making them touch only on a
// diagonal.
func gradientSquared(w io.Writer, width, height, startRadius float64, ease ease) {
	// r*sqrt(2) is the size of a square that contains the circle
	numLines := height / (startRadius * math.Sqrt2)
	numCols := width / (startRadius * math.Sqrt2)

	// make the radius a bit bigger to have a nice overlay
	radiusFactor := 2.0

	for y := 0; y < int(numLines); y++ {
		for x := 0; x < int(numCols); x++ {
			xPos, yPos := float64(x)*startRadius*math.Sqrt2*2, float64(y)*startRadius*math.Sqrt2

			if y%2 == 0 {
				xPos += startRadius * math.Sqrt2
			}

			// we go from [0,startRadius]

			radius := (startRadius * (float64(y) / numLines))
			coeff := ease(radius / (startRadius))
			radius *= coeff
			radius *= radiusFactor
			//radius = startRadius

			fmt.Fprintf(w, "<circle cx=\"%.2f\" cy=\"%.2f\" r=\"%.2f\"/>", xPos, yPos, radius)
		}
	}
}

// Ease functions to adjust the radius size evolution. Taken from
// https://github.com/fogleman/ease.

func easeOutCirc(t float64) float64 {
	t--
	return math.Sqrt(1 - (t * t))
}

func easeOutExpo(t float64) float64 {
	if t == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}

func easeInExpo(t float64) float64 {
	if t == 0 {
		return 0
	}
	return math.Pow(2, 10*(t-1))
}

func easeInQuad(t float64) float64 {
	return t * t
}

func easeLinear(t float64) float64 {
	return t
}

func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	t = 2*t - 1
	return -0.5 * (t*(t-2) - 1)
}

func easeInOutQuart(t float64) float64 {
	t *= 2
	if t < 1 {
		return 0.5 * t * t * t * t
	}
	t -= 2
	return -0.5 * (t*t*t*t - 2)
}
