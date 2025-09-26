//go:build generate
// +build generate

package opencode

//go:generate sh -c "ogen -clean -config ./ogen.yml -target ./ogen -package ogenapi openapi.yaml"
