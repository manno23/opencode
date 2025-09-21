//go:build generate
// +build generate

package api

//go:generate sh -c "ogen -clean -config ./ogen.yml -target ./ogen -package api unified-openapi.json"

