# ----- helpers --------------------------------------------------------------

def opkeys: ["get","put","post","delete","patch","head","options","trace"];
def isobj: (type == "object");
def isarr: (type == "array");

# Keep the first occurrence of each (in,name) pair
def dedupe_params:
  (. // [])
  | if (type!="array") then .
    else reduce .[] as $p (
           {seen:{}, out:[]};
           (($p.in // "") + "|" + ($p.name // "")) as $k
           | if .seen[$k] then . else .seen[$k]=true | .out += [$p] end
         ) | .out
  end;

# Remove a fixed set of 4xx/5xx from responses (object-only)
def strip_error_statuses:
  if isobj then
    with_entries(
      select(
        (["400","401","403","404","409","422","429","500","501","502","503"] | index(.key)) | not
      )
    )
  else . end;

# Ensure a node exists at path (set to init if null / missing)
def ensure(path; init):
  if (getpath(path) | type) == "null" then setpath(path; init) else . end;

# ----- 1) Ensure canonical ErrorResponse schema ----------------------------
(.components |= (. // {}))
| (.components.schemas |= (. // {}))
| (.components.schemas.ErrorResponse =
  {
    type: "object",
    "x-ogen-name": "ErrorResponse",
    properties: {
      error:   { type: "string" },
      code:    { type: "string" },
      details: { type: "object", additionalProperties: true }
    },
    required: ["error"],
    additionalProperties: false
  })

# ----- 2) De-dupe parameters at path- and op-level -------------------------
| (.paths |=
    if isobj then
      ( . | to_entries
        | map(
            if (.value | isobj) then
              { key: .key,
                value:
                  ( .value
                    # path-level parameters
                    | if (has("parameters") and (.parameters|isarr))
                      then .parameters |= dedupe_params else . end
                    # op-level parameters on real HTTP methods
                    | ( reduce opkeys[] as $m (.;
                          if (has($m) and (.[ $m ]|isobj) and (.[ $m ].parameters|isarr))
                          then .[ $m ].parameters |= dedupe_params
                          else . end
                        )
                      )
                  )
              }
            else .
            end
          )
        | from_entries )
    else .
    end
  )

# ---- 3) Responses: drop 4xx/5xx, inject shared default (bulletproof) ------
| (.paths |=
    if isobj then
      ( . | to_entries
        | map(
            if (.value | isobj) then
              { key: .key,
                value:
                  ( .value
                    | reduce opkeys[] as $m (.;
                        if (has($m) and (.[ $m ]|isobj)) then
                          # 3a) force responses to be an object ({} if missing/array/null/other)
                          .[ $m ].responses = (.[ $m ].responses
                                                | if type=="object" then . else {} end)
                          # 3b) drop common 4xx/5xx without with_entries
                          | (.[ $m ].responses |= del(
                               .["400","401","403","404","409","422","429",
                                 "500","501","502","503"]
                             ))
                          # 3c) set/overwrite default
                          | .[ $m ].responses.default = {
                              description: "Generic error",
                              content: {
                                "application/json": {
                                  schema: { "$ref": "#/components/schemas/ErrorResponse" }
                                }
                              }
                            }
                        else . end)
                  )
              }
            else .
            end
          )
        | from_entries )
    else .
    end
  )

# ----- 4) Replace Config.provider.*.options.timeout with int-or-false union -
| ( ensure(["components"]; {})
  | ensure(["components","schemas"]; {})
  | ensure(["components","schemas","Config"]; {"type":"object","properties":{}})
  | ensure(["components","schemas","Config","properties"]; {})
  | ensure(["components","schemas","Config","properties","provider"]; {})
  | ensure(["components","schemas","Config","properties","provider","additionalProperties"]; {"type":"object","properties":{}})
  | ensure(["components","schemas","Config","properties","provider","additionalProperties","properties"]; {})
  | ensure(["components","schemas","Config","properties","provider","additionalProperties","properties","options"]; {"type":"object","properties":{}})
  | ensure(["components","schemas","Config","properties","provider","additionalProperties","properties","options","properties"]; {}) )
| setpath(
    ["components","schemas","Config","properties","provider","additionalProperties","properties","options","properties","timeout"];
    {
      description: "Timeout in ms. Use `false` to disable.",
      oneOf: [
        { type: "integer", minimum: 0, exclusiveMinimum: true, maximum: 9007199254740991 },
        { type: "boolean", enum: [false] }
      ],
      default: 300000
    }
  )
