package main

import (
	"fmt"
	"syscall/js"

	"github.com/nkcr/go-halftone/halftone"
)

func main() {
	fmt.Println("Go assembly test")
	js.Global().Set("golangRender", renderStrWrapper())
	<-make(chan struct{})
}

func renderStrWrapper() js.Func {
	renderStr := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 5 {
			return "Invalid number of arguments passed"
		}

		width := args[0].Float()
		height := args[1].Float()
		radius := args[2].Float()
		color := args[3].String()
		rendering := args[4].String()

		var proc halftone.Proc = halftone.GradientSquared
		if rendering == "compact" {
			proc = halftone.GradientCompact
		}

		opts := halftone.Opts{
			Width:  width,
			Height: height,
			Radius: radius,
			Fill:   color,
			Proc:   proc,
			Ease:   halftone.EaseLinear,
		}

		return halftone.RenderStr(opts)
	})

	return renderStr
}
