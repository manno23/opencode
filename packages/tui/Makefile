.PHONY: spec merge sdk build clean

SPEC_DIR := api
SPEC_BASE := $(SPEC_DIR)/base-openapi.json
SPEC_YAML := $(SPEC_DIR)/openapi.yaml
SPEC_UNIFIED := $(SPEC_DIR)/unified-openapi.yaml

spec:
	# Generate base spec from server (Hono/Zod). Try bun route first, fallback to opencode CLI.
	bun run dev --silent generate > $(SPEC_BASE) || opencode generate > $(SPEC_BASE)

merge:
	@if [ ! -f "$(SPEC_BASE)" ]; then echo "Missing $(SPEC_BASE). Run 'make spec' first."; exit 1; fi
	@if [ ! -f "$(SPEC_YAML)" ]; then echo "Missing $(SPEC_YAML). Create your delta YAML."; exit 1; fi
	# Deep-merge base (JSON) with delta (YAML) using yq v4; delta overrides base
	yq4 -p json -p yaml -o json eval-all 'select(fileIndex==0) *+ select(fileIndex==1)' "api/base-openapi.json" "api/openapi.yaml" > "api/unified-openapi.json"
	@echo "Wrote $(SPEC_UNIFIED)"

sdk:
	cd $(SPEC_DIR) && go generate ./...

build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$$(git describe --tags --always 2>/dev/null || echo dev)" -o tui cmd/opencode/main.go

clean:
	rm -rf api/ogen tui $(SPEC_BASE) $(SPEC_UNIFIED)
