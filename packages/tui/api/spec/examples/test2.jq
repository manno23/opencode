# Helpers
def opkeys: ["get","put","post","delete"];
def isobj: (type=="object");
def g(path; v): setpath(path; v);
def gp(path): getpath(path);

# Ensure components.schemas
(.components |= (. // {}))
| (.components.schemas |= (. // {}))

# Walk every path item and its HTTP methods
| (.paths |=
    ( . as $root
    | to_entries
    | map(
        if (.value|isobj) then
          { key: .key,
            value:
              ( .value
              | reduce opkeys[] as $m (.;
                  if has($m) and (.[ $m ]|isobj) and (.[ $m ].responses|isobj) then
                    # consider common success codes; adjust if you use others
                    (["200","201","202"] | .[]) as $sc
                    | if (.[ $m ].responses[$sc].content["application/json"].schema // null) then
                        # fetch operationId (fallback to path+method if missing)
                        (.[ $m ].operationId // (. + {opid: (. + "-" + $m)}).opid) as $opid
                        | (.[ $m ].responses[$sc].content["application/json"].schema) as $S

                        # Compute a unique wrapper name
                        | ("Resp_" + ($opid | gsub("[^A-Za-z0-9_]"; "_")) + "_" + $sc) as $wrapperName
                        | ("#/components/schemas/" + $wrapperName) as $wrapperRef

                        # Case 1: schema is a $ref → create a wrapper component with allOf
                        | if ($S|type)=="object" and ($S["$ref"] // "") | startswith("#/components/schemas/") then
                            # Create wrapper only if absent
                            | ( if ($root.components.schemas[$wrapperName]? | type) == "null"
                                then $root
                                     | g(["components","schemas",$wrapperName];
                                         { allOf: [ { "$ref": $S["$ref"] } ],
                                           "x-ogen-name": $wrapperName })
                                else $root end ) as $root2
                            # Switch op schema to wrapper
                            | ($root2
                               | g(["paths",.key,$m,"responses",$sc,"content","application/json","schema"];
                                    { "$ref": $wrapperRef }))
                          else
                            # Case 2: inline schema → hoist to a named component
                            ( $root
                              # put (or overwrite) the component
                              | g(["components","schemas",$wrapperName];
                                  ($S + { "x-ogen-name": $wrapperName }))
                              # and switch op schema to the new ref
                              | g(["paths",.key,$m,"responses",$sc,"content","application/json","schema"];
                                   { "$ref": $wrapperRef })
                            )
                          end
                      else .
                      end
                  )
                )
              )
          }
        else .
        end
      )
    | from_entries )
  )
