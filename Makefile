.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o ./web/assets/halftone.wasm ./wasm/