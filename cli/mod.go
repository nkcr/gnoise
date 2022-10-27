package main

import (
	"github.com/nkcr/go-halftone/halftone"
)

func main() {
	filename := "result.svg"

	opts := halftone.Opts{
		Width:  50.0,
		Height: 50.0,
		Radius: 0.5,
		Fill:   "#44FF44",
		Proc:   halftone.GradientCompact,
		Ease:   halftone.EaseLinear,
	}

	err := halftone.RenderFile(opts, filename)
	if err != nil {
		panic(err)
	}
}
