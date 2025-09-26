#!/usr/bin/env bash
set -euo pipefail
B=${1:-spec.bundle.json}
mkdir -p report

jq -r '.paths|to_entries[]|. as $p|.value|to_entries[]|select(.value.responses.default|not)|{path:$p.key,method:.key,operationId:.value.operationId,responses:(.value.responses|keys)}' "$B" > report/ops_missing_default.json

jq -r '.paths|to_entries[] as $p|$p.value|to_entries[]|select(.value.responses.default and ((.value.responses|keys|map(test("^[45]..$"))|any)))|{path:$p.key,method:.key,operationId:.value.operationId,error_statuses:(.value.responses|keys|map(select(test("^[45]..$"))))}' "$B" > report/ops_default_plus_4xx.json

#jq -r '.paths|to_entries[] as $p|$p.value|to_entries[] as $op|(($op.value.parameters//[])|map({in,name,loc:"op"})) as $opP|(($p.value.parameters//[])|map({in,name,loc:"path"})) as $pathP|($opP+$pathP) as $params|[group_by(.in+"|"+.name)[]|select(length>1)|{in:.[0].in,name:.[0].name,sources:map(.loc)}] as $dups|select(($dups|length)>0)|{path:$p.key,method:$op.key,operationId:$op.value.operationId,duplicates:$dups}' "$B" > report/dup_params.json
# jq -r '
# .paths
# | to_entries[] as $p
# | $p.value
# | to_entries[]
# | . as $op
# | ($op.value.responses // {}) as $resp
# | select(($resp | has("default")) and
#          ([$resp|keys[] | select(test("^[45][0-9][0-9]$"))] | length > 0))
# | {
#     path: $p.key,
#     method: $op.key,
#     operationId: $op.value.operationId,
#     error_statuses: ([$resp|keys[] | select(test("^[45][0-9][0-9]$"))])
#   }' "$B" > report/dup_params.json


#jq -r '.paths|to_entries[] as $p|$p.value|to_entries[]|select(.value.responses["200"].content["text/event-stream"])|{path:$p.key,method:.key,operationId:.value.operationId,has_default:(.value.responses.default|type?//"null"),other_errors:(.value.responses|keys|map(select(test("^[45]..$"))))}' "$B" > report/sse_ops.json

# jq -r '
# def params(x): (x // []) | map({in, name, __loc});
# .paths
# | to_entries[] as $p
# | ($p.value.parameters // []) as $pp
# | $p.value
# | to_entries[] as $op
# | ($op.value.parameters // []) as $opp
# | ( [ ($pp | map(. + {__loc:"path"})) , ($opp | map(. + {__loc:"op"})) ] | add ) as $all
# | ($all
#    | group_by(.in + "|" + .name)
#    | map(select(length > 1)
#          | {in: .[0].in, name: .[0].name, sources: (map(.__loc) | unique)})
#   ) as $dups
# | select(($dups | length) > 0)
# | {path:$p.key, method:$op.key, operationId:$op.value.operationId, duplicates:$dups}
# ' "$B" > report/sse_ops.json
#
# jq -r '.components.schemas|to_entries[]|select(.value.oneOf or .value.anyOf)|{schema:.key,kind:(if .value.oneOf then "oneOf" else "anyOf" end),branches:((.value.oneOf//.value.anyOf)|map(if has("$ref") then .["$ref"]|split("/")|last else (.title//.type//"inline") end))}' "$B" > report/union_candidates.json

jq -r 'def refs: ..|objects|select(has("$ref"))|.["$ref"]; def schema_keys:(.components.schemas|keys)//[]; [refs|select(startswith("#/components/schemas/"))|ltrimstr("#/components/schemas/")] as $r|schema_keys as $k|{missing_schema_refs:($r|unique|map(select(. as $x|$k|index($x)|not))),total_unique_refs:($r|unique|length),total_schemas:($k|length)}' "$B" > report/undefined_refs.json

jq -r '.components.schemas|to_entries[]|select(.value|(has("x-ogen-name") or has("x-ogen-sum-type") or has("x-ogen-map")))|{schema:.key,x_ogen:{name:.value["x-ogen-name"],sum:.value["x-ogen-sum-type"],map:.value["x-ogen-map"]}}' "$B" > report/x_ogen_usage.json

jq -r '[
  .paths|to_entries[] as $p
  | $p.value|to_entries[]
  | .value.responses|to_entries[]
  | .value.content?|keys[]?
]|group_by(.)|map({media_type:.[0],count:length})' "$B" > report/media_types.json

echo "Wrote: report/*.json"

