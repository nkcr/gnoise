.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o ./docs/assets/halftone.wasm ./wasm/