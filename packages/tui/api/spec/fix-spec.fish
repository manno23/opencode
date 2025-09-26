#!/usr/bin/env fish

functions clean

set -g DIR (dirname (status --current-filename))
pushd $DIR

jq -f fix-spec.jq openapi.base.json > openapi.fixed.json
redocly bundle openapi.fixed.json -o openapi.bundle.yaml
uvx oas-patch overlay openapi.bundle.yaml examples/openapi.ogen-overlay.yaml -o openapi.fixed.yaml

yq -pjson -oyaml openapi.fixed.json > ../openapi.yaml

rm openapi.fixed.json

popd

# jq '
#   def dedupe:
#     (. // [])
#     | reduce .[] as $p ({seen:{}, out:[]};
#         ($p.in + "|" + $p.name) as $k
#         | if .seen[$k] then . else .seen[$k]=true | .out += [$p] end
#       ) | .out;
#
#   .paths |=
#     ( with_entries( if .value.parameters then .value.parameters |= dedupe else . end )
#     | with_entries( .value |= with_entries( if .value.parameters then .value.parameters |= dedupe else . end ) )
#     )
# ' openapi.json > openapi.out.json && mv openapi.out.json openapi.json
#
# jq '
#   # Replace the entire timeout schema
#   .components.schemas.Config.properties.provider.additionalProperties.properties.options.properties.timeout =
#   {
#     description: "Timeout in ms. Use `false` to disable.",
#     oneOf: [
#       { "type":"integer", "minimum":0, "exclusiveMinimum":true, "maximum":9007199254740991 },
#       { "type":"boolean", "enum":[false] }
#     ],
#     default: 300000
#   }
#   |
#   # Ensure a canonical DirectoryQuery parameter
#   (.components.parameters //= {}) |
#   .components.parameters.DirectoryQuery = {
#     "name":"directory",
#     "in":"query",
#     "description":"Root directory / workspace context",
#     "schema":{"type":"string"},
#     "required":false
#   }
# ' openapi.json > openapi.out.json && mv openapi.out.json openapi.json
#
# jq '
#   .components.parameters |= (. // {} | .DirectoryQuery = {
#     "name":"directory","in":"query",
#     "description":"Root directory / workspace context",
#     "schema":{"type":"string"},"required":false
#   }) |
#   .paths |= (to_entries
#     | map(
#         .value as $path
#         | (.value.parameters // []) as $pp
#         # remove path-level
#         | .value.parameters = ($pp | map(select(.in!="query" or .name!="directory")))
#         # scrub & re-add at op-level
#         | .value = (.value | with_entries(
#             if .key == "parameters" then .
#             else
#               .value.parameters = (
#                 ((.value.parameters // [])
#                  | map(select(.in!="query" or .name!="directory")))
#                 + [{"$ref":"#/components/parameters/DirectoryQuery"}]
#               )
#             end))
#       )
#     | from_entries)
# ' openapi.json > openapi.out.json && mv openapi.out.json openapi.json

# Why this is safer: path-level + op-level parameters are unioned by generators. If you just “dedupe arrays” separately, you can still get a duplicate via inheritance. Nuking both and re-adding one $ref at the op level avoids that.
#
# If a few endpoints should not have ?directory, do a second pass to remove the $ref for those ops only, e.g.:
# yq -i '
#   del(.paths["/auth/token"].get.parameters[]
#       | select(."$ref" == "#/components/parameters/DirectoryQuery"))
# ' openapi.yaml

# uvx oas-patch overlay openapi.json openapi.ogen-overlay.yaml -o ../openapi.fixed.json
