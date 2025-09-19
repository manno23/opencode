.PHONY: spec merge sdk build clean

SPEC_DIR := api
SPEC_BASE := $(SPEC_DIR)/base-openapi.json
SPEC_YAML := $(SPEC_DIR)/openapi.yaml
SPEC_UNIFIED := $(SPEC_DIR)/unified-openapi.yaml

spec:
	# Generate base spec from server (Hono/Zod). Try bun route first, fallback to opencode CLI.
	bun run --silent dev generate > $(SPEC_BASE) || opencode generate > $(SPEC_BASE)

merge:
	# Merge base JSON with custom YAML delta. Requires yq v4.
	yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' $(SPEC_BASE) $(SPEC_YAML) > $(SPEC_UNIFIED)

sdk:
	cd $(SPEC_DIR) && go generate ./...

build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=$$(git describe --tags --always 2>/dev/null || echo dev)" -o tui cmd/opencode/main.go

clean:
	rm -rf api/ogen tui $(SPEC_BASE) $(SPEC_UNIFIED)
