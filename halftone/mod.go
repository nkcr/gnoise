package halftone

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
)

// Proc is a processing function that render points
type Proc func(w io.Writer, width, height, startRadius float64, ease Ease)

// Ease is an ease function to parametrize the radius evolution
type Ease func(float64) float64

// main options, would be a good idea to set them as flags
type Opts struct {
	Width  float64
	Height float64
	Radius float64
	Fill   string
	Proc   Proc
	Ease   Ease
}

func RenderFile(opts Opts, filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	defer f.Close()

	render(f, opts)

	return nil
}

func RenderStr(opts Opts) string {
	w := &bytes.Buffer{}

	render(w, opts)

	return w.String()
}

func render(w io.Writer, opts Opts) {
	fmt.Fprintf(w, `<svg viewBox="0 0 %.1f %.1f" xmlns="http://www.w3.org/2000/svg" id="svg">`, opts.Width, opts.Height)
	fmt.Fprintf(w, `<g fill="%s">`, opts.Fill)
	opts.Proc(w, opts.Width, opts.Height, opts.Radius, opts.Ease)
	fmt.Fprint(w, "</g>")
	fmt.Fprint(w, `</svg>`)
}

// GradientCompact stacks circles leaving the smallest gap possible.
func GradientCompact(w io.Writer, width, height, startRadius float64, ease Ease) {
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

// GradientSquared puts each circle on a square, making them touch only on a
// diagonal.
func GradientSquared(w io.Writer, width, height, startRadius float64, ease Ease) {
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

func EaseOutCirc(t float64) float64 {
	t--
	return math.Sqrt(1 - (t * t))
}

func EaseOutExpo(t float64) float64 {
	if t == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*t)
}

func EaseInExpo(t float64) float64 {
	if t == 0 {
		return 0
	}
	return math.Pow(2, 10*(t-1))
}

func EaseInQuad(t float64) float64 {
	return t * t
}

func EaseLinear(t float64) float64 {
	return t
}

func EaseInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	t = 2*t - 1
	return -0.5 * (t*(t-2) - 1)
}

func EaseInOutQuart(t float64) float64 {
	t *= 2
	if t < 1 {
		return 0.5 * t * t * t * t
	}
	t -= 2
	return -0.5 * (t*t*t*t - 2)
}
