package opencode

func main() {
	readSpec()
}

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --package api  --target api --clean openapi.json
