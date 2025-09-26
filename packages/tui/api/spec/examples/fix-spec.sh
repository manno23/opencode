#!/usr/bin/env bash
set -euo pipefail

SRC=${1:-openapi.json}
TMP=.tmp.spec.json

redocly bundle ../openapi.base.json -o openapi.yaml

cp "$SRC" "$TMP"

echo "1) dedupe params... "
yq -i '
  def dedupe:
    (. // []) as $a
    | reduce $a[] as $p ({seen:{}, out:[]};
        ($p.in + "|" + $p.name) as $k
        | if .seen[$k] then . else .seen += {($k):true} | .out += [$p] end
      ) | .out;
  (.paths[] | select(has("parameters")) | .parameters) |= dedupe |
  (.paths[] | .[]? | select(has("parameters")) | .parameters) |= dedupe
' "$TMP"
echo "done.\n"

echo "2) remove common 4xx/5xx and add default... "
for s in 400 401 403 404 409 422 429 500 501 502 503; do
  yq -i ".paths[]|.*?|select(has(\"responses\") and .responses[\"$s\"]).responses |= del(.\"$s\")" "$TMP"
done
echo "done.\n"

yq -i '
  (.components.schemas //= {}) |
  .components.schemas.ErrorResponse //= {
    type: "object",
    x-ogen-name: "ErrorResponse",
    properties: { error: {type:"string"}, code:{type:"string"}, details:{type:"object","additionalProperties":true} },
    required: ["error"],
    additionalProperties: false
  } |
  (.paths[] | .[]? | select(has("responses")) | .responses) |=
    (. + { "default": { description: "Generic error",
                        content: { "application/json": { schema: {"$ref":"#/components/schemas/ErrorResponse"} } } } })
' "$TMP"
echo "done.\n"

# 3) fix known 3.1 exclusives → 3.0 form (example: Config.timeout); add your other paths as needed
echo "3) fix known 3.1 exclusives... "
yq -i '
  with(.components.schemas.Config.properties.provider.additionalProperties.properties.options.properties.timeout;
       . = {
         description: "Timeout in ms. Use `false` to disable.",
         oneOf: [
           {"type":"integer","minimum":0,"exclusiveMinimum":true,"maximum":9007199254740991},
           {"type":"boolean","enum":[false]}
         ],
         default: 300000
       })
' "$TMP"
echo "done.\n"

# 4) output final JSON for ogen
redocly bundle "$TMP" -o ../openapi.fixed.json --ext json

