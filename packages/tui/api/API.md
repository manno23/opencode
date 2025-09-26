# MANUAL SPEC UPDATES 

** HIGH IMPORTANCE FIXES **

1.  ### Ogen does not generate `type Session` for the Go client
    #### FIXED BY OVERLAY  openapi.overlay.json
  - The schema uses complex features like anyOf for unions (e.g., revert/share fields)
    and sum types with overlapping names,
  - see INFO logs from go generate:
  "sum types with same names not implemented", "complex anyOf not implemented")
  - This results in Opt* wrappers for Session-related req/res types but no base Session struct
  - Need to:
    1. Simplify the spec: Add explicit discriminators to the Session schema
       (e.g., discriminator: {propertyName: "type"} for anyOf branches)
       in openapi.overlay.json or the base spec.
    2. Regenerate: cd packages/tui/api && go generate ./....
    3. Verify: Grep for "type Session struct" in ogen/oas_schemas_gen.go
       – it should appear with fields like ID, Title, etc.
   - The Fix
    • Added the Session schema with a discriminator on the union field (e.g., revert as anyOf with x-ogen-discriminator mapping to Revert/Share types).
    • Used x-ogen-oneof-discriminator extension for the top-level Event union (already partially there, but enhanced for Session-related events).
    • Ensured optional fields use OptString etc. for compatibility.


2.  ### Missing Parameter in Property
    #### FIXED IN SERVER CHANGES , ALSO CAN BE FIXED BY SCRIPT  tui/api/fix-openapi.ts
  - What the fix did: Added location specifiers (query, path, header, cookie) to parameters.
  - Impact on request handling:
      Significant - Without in, the server literally doesn't know where to look for the parameter
      A parameter could be in the URL path, query string, headers, or cookies
      This affects actual request parsing and validation
  - Seriousness of ignoring: 🔴 HIGH (9/10)
  - This is a real issue - the OpenAPI spec is genuinely invalid without it.


3.  ### Missing Parameter schema Property
    #### FIXED IN SERVER CHANGES , ALSO CNA BE FIXED BY SCRIPT  tui/api/fix-openapi.ts
  - What the fix did: Added type definitions ({type: "string"}, etc.) to parameters.
  - Impact on request handling:
  - Moderate - Without schema, there's no type validation or coercion
  - Request handlers might receive strings when expecting numbers
  - No validation of format, length, patterns, etc.
  - Seriousness of ignoring: 🟡 MEDIUM (6/10)
  - Functionally works but loses type safety and validation. Your API becomes more fragile.


## LOW IMPORTANCE FIXES **

4.  ### Type Mismatch - exclusiveMinimum Number → Boolean
    #### FIXED BY SCRIPT  tui/api/fix-openapi.ts
  - What the fix did: Changed exclusiveMinimum: 0 to exclusiveMinimum: true, minimum: 0.
  - Impact on request handling:
  - Cosmetic for validation - Both mean "greater than 0"
  - The actual validation behavior is identical
  - Just different ways to express the same constraint
  - Seriousness of ignoring: 🟢 LOW (2/10)
  - This is purely a spec format issue - ogen is being pedantic about OpenAPI 3.1 vs 3.0 syntax. Your API works identically either way.


## IGNORED ERRORS

- These are all appearing on the final spec
- Probably alot here as seen by the linter

`$ redocly lint openapi.unified.json`

  #### IGNORED IN OGEN CONFIG   ogen.yml
  1. Sum types with same names
    - Not sure all the problems missing here
  2. complex anyOf
    - Also not sure here




## Build Steps

1. Generate base: Server routes → openapi.base.json.    Includes Fixes 
2. Merge openapi.base.json into our overlay with fixes, openapi.overlay.json

`bun run dev generate > openapi.base.json`
`yq eval-all 'select(fileIndex == 0) * select(fileIndex == 1)' openapi.base.json openapi.overlay.json > openapi.unified.json`





$ cd api && go generate ./..
  - openapi.unified.json:3161:47 -> cannot unmarshal !!int `0` into boo
"description": "Timeout in milliseconds for requests to this provider. Default is 300000 (5 minutes). Set to false disable timeout."
          3160 |                           "type": "integer"
        → 3161 |                           "exclusiveMinimum": 0
          3162 |                           "maximum": 900719925474099
          3163 |                         }
          3165 |                           "description": "Disable timeout for this provider entirely."

generation failed
    main.ru
        /home/jm/lib/go/pkg/mod/github.com/ogen-go/ogen@v1.14.0/cmd/ogen/main.go:37
generate.go:6: running "sh": exit status 
