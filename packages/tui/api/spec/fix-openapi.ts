#!/usr/bin/env bun

import SwaggerParser from '@apidevtools/swagger-parser';
import { type OpenAPIV3_1 } from 'openapi-types';
import { writeFileSync } from 'fs';
import { resolve } from 'path';

/**
 * Comprehensive OpenAPI 3.1.1 spec fixer for ogen compatibility
 * 
 * Fixes common issues:
 * - exclusiveMinimum/exclusiveMaximum numeric values → boolean
 * - Missing parameter 'in' property
 * - Missing parameter 'schema' or 'content' property
 * - Invalid parameter locations
 * - Missing required properties in schemas
 */

type ParameterObject = OpenAPIV3_1.ParameterObject;
type SchemaObject = OpenAPIV3_1.SchemaObject;
type OpenAPIDocument = OpenAPIV3_1.Document;

interface FixerOptions {
  inputFile: string;
  outputFile: string;
  defaultParameterLocation: 'query' | 'header' | 'path' | 'cookie';
  verbose: boolean;
}

class OpenAPIFixer {
  private options: FixerOptions;
  private fixes: string[] = [];

  constructor(options: Partial<FixerOptions> = {}) {
    this.options = {
      inputFile: 'openapi.bundled.json',
      outputFile: 'openapi.fixed.json',
      defaultParameterLocation: 'query',
      verbose: true,
      ...options
    };
  }

  private log(message: string): void {
    if (this.options.verbose) {
      console.log(`[OpenAPI Fixer] ${message}`);
    }
  }

  private addFix(fix: string): void {
    this.fixes.push(fix);
    this.log(fix);
  }

  /**
   * Infer schema type based on parameter name patterns
   */
  private inferSchemaFromParameterName(name: string, location: string): SchemaObject {
    const lowerName = name.toLowerCase();

    // UUID/ID patterns
    if (lowerName.includes('id') || lowerName.includes('uuid')) {
      return {
        type: 'string',
        format: location === 'path' ? 'uuid' : undefined,
        description: `${name} identifier`
      };
    }

    // Pagination patterns
    if (['limit', 'offset', 'page', 'size', 'count', 'per_page'].includes(lowerName)) {
      return {
        type: 'integer',
        minimum: 0,
        maximum: lowerName === 'limit' ? 1000 : undefined,
        description: `${name} for pagination`
      };
    }

    // Sorting patterns
    if (['sort', 'order', 'direction', 'sort_by', 'order_by'].includes(lowerName)) {
      return {
        type: 'string',
        enum: ['asc', 'desc'],
        description: `Sort ${name}`
      };
    }

    // Date/time patterns
    if (lowerName.includes('date') || lowerName.includes('time') || 
        ['created_at', 'updated_at', 'deleted_at'].includes(lowerName)) {
      return {
        type: 'string',
        format: 'date-time',
        description: `${name} timestamp`
      };
    }

    // Email patterns
    if (lowerName.includes('email')) {
      return {
        type: 'string',
        format: 'email',
        description: `${name} address`
      };
    }

    // Boolean patterns
    if (['active', 'enabled', 'disabled', 'deleted', 'published'].includes(lowerName) ||
        lowerName.startsWith('is_') || lowerName.startsWith('has_')) {
      return {
        type: 'boolean',
        description: `${name} flag`
      };
    }

    // Search/filter patterns
    if (['search', 'query', 'q', 'filter', 'term'].includes(lowerName)) {
      return {
        type: 'string',
        minLength: 1,
        description: `${name} term`
      };
    }

    // Default to string
    return {
      type: 'string',
      description: `${name} parameter`
    };
  }

  /**
   * Infer parameter location based on name patterns
   */
  private inferParameterLocation(name: string): 'query' | 'path' | 'header' {
    const lowerName = name.toLowerCase();

    // Path parameters
    if (lowerName === 'id' || lowerName.endsWith('_id') || 
        lowerName === 'uuid' || lowerName.endsWith('_uuid') ||
        ['slug', 'username', 'handle'].includes(lowerName)) {
      return 'path';
    }

    // Header parameters
    if (lowerName.startsWith('x-') || 
        ['authorization', 'content-type', 'accept', 'user-agent'].includes(lowerName)) {
      return 'header';
    }

    // Default to query
    return 'query';
  }

  /**
   * Fix exclusiveMinimum and exclusiveMaximum numeric values
   */
  private fixExclusiveMinMax(obj: any): any {
    if (typeof obj !== 'object' || obj === null) {
      return obj;
    }

    if (Array.isArray(obj)) {
      return obj.map(item => this.fixExclusiveMinMax(item));
    }

    const result = { ...obj };

    // Fix exclusiveMinimum
    if ('exclusiveMinimum' in result && typeof result.exclusiveMinimum === 'number') {
      const oldValue = result.exclusiveMinimum;
      result.exclusiveMinimum = oldValue !== 0;
      this.addFix(`Fixed exclusiveMinimum: ${oldValue} → ${result.exclusiveMinimum}`);
    }

    // Fix exclusiveMaximum
    if ('exclusiveMaximum' in result && typeof result.exclusiveMaximum === 'number') {
      const oldValue = result.exclusiveMaximum;
      result.exclusiveMaximum = oldValue !== 0;
      this.addFix(`Fixed exclusiveMaximum: ${oldValue} → ${result.exclusiveMaximum}`);
    }

    // Recursively fix nested objects
    for (const [key, value] of Object.entries(result)) {
      if (typeof value === 'object' && value !== null) {
        result[key] = this.fixExclusiveMinMax(value);
      }
    }

    return result;
  }

  /**
   * Fix parameter objects
   */
  private fixParameter(param: any): ParameterObject {
    const result = { ...param };

    // Fix missing 'name' (should not happen in valid specs)
    if (!result.name) {
      result.name = 'unnamed_param';
      this.addFix(`Added missing parameter name: unnamed_param`);
    }

    // Fix missing 'in' property
    if (!result.in) {
      result.in = this.inferParameterLocation(result.name);
      this.addFix(`Added missing 'in' property: ${result.name} → ${result.in}`);
    }

    // Validate 'in' property
    const validLocations = ['query', 'header', 'path', 'cookie'];
    if (!validLocations.includes(result.in)) {
      const oldLocation = result.in;
      result.in = this.inferParameterLocation(result.name);
      this.addFix(`Fixed invalid location: ${result.name} ${oldLocation} → ${result.in}`);
    }

    // Fix missing schema/content
    if (!result.schema && !result.content) {
      result.schema = this.inferSchemaFromParameterName(result.name, result.in);
      this.addFix(`Added missing schema for parameter: ${result.name}`);
    }

    // Fix path parameters (must be required)
    if (result.in === 'path' && result.required !== true) {
      result.required = true;
      this.addFix(`Fixed path parameter required: ${result.name} → required: true`);
    }

    return result as ParameterObject;
  }

  /**
   * Walk through the OpenAPI document and fix issues
   */
  private fixDocument(doc: OpenAPIDocument): OpenAPIDocument {

    // First, fix exclusiveMin/Max throughout the document
    const fixedDoc = this.fixExclusiveMinMax(doc) as OpenAPIDocument;

    // Fix parameters in paths
    if (fixedDoc.paths) {
      for (const [pathName, pathItem] of Object.entries(fixedDoc.paths)) {
        if (!pathItem || typeof pathItem !== 'object') continue;

        for (const [method, operation] of Object.entries(pathItem)) {
          if (!operation || typeof operation !== 'object' || 
              !['get', 'post', 'put', 'delete', 'patch', 'head', 'options', 'trace'].includes(method)) {
            continue;
          }

          // Fix operation parameters
          if (operation.parameters && Array.isArray(operation.parameters)) {
            operation.parameters = operation.parameters.map((param: any) => {
              if (param && typeof param === 'object' && !param.$ref) {
                return this.fixParameter(param);
              }
              return param;
            });
          }
        }

        // Fix path-level parameters
        if (pathItem.parameters && Array.isArray(pathItem.parameters)) {
          pathItem.parameters = pathItem.parameters.map((param: any) => {
            if (param && typeof param === 'object' && !param.$ref) {
              return this.fixParameter(param);
            }
            return param;
          });
        }
      }
    }

    return fixedDoc;
  }

  /**
   * Main fix method
   */
  async fix(): Promise<void> {
    try {
      this.log(`Reading OpenAPI spec from: ${this.options.inputFile}`);
      
      // Read and parse the OpenAPI document
      const inputPath = resolve(this.options.inputFile);
      const api = await SwaggerParser.parse(inputPath) as OpenAPIDocument;
      
      this.log(`Parsed OpenAPI ${api.openapi} spec successfully`);
      this.log(`Title: ${api.info?.title || 'Unknown'}`);
      this.log(`Version: ${api.info?.version || 'Unknown'}`);

      // Apply fixes
      this.log('Applying fixes...');
      const fixedDoc = this.fixDocument(api);

      // Validate the fixed document
      this.log('Validating fixed document...');
      await SwaggerParser.validate(fixedDoc);
      this.log('✅ Fixed document is valid!');

      // Write the fixed document
      const outputPath = resolve(this.options.outputFile);
      writeFileSync(outputPath, JSON.stringify(fixedDoc, null, 2), 'utf8');
      
      this.log(`✅ Fixed OpenAPI spec written to: ${this.options.outputFile}`);
      
      if (this.fixes.length > 0) {
        this.log('\n📋 Summary of fixes applied:');
        this.fixes.forEach((fix, index) => {
          this.log(`  ${index + 1}. ${fix}`);
        });
      } else {
        this.log('ℹ️  No fixes were needed');
      }

    } catch (error) {
      console.error('❌ Error fixing OpenAPI spec:', error);
      process.exit(1);
    }
  }
}

// CLI Usage
async function main() {
  const args = process.argv.slice(2);
  
  if (args.includes('--help') || args.includes('-h')) {
    console.log(`
OpenAPI Spec Fixer for ogen compatibility

Usage: tsx fix-openapi.ts [options]

Options:
  --input <file>     Input OpenAPI file (default: openapi.bundled.json)
  --output <file>    Output file (default: openapi.fixed.json)
  --location <loc>   Default parameter location (default: query)
  --quiet           Suppress verbose output
  --help            Show this help message

Examples:
  tsx fix-openapi.ts
  tsx fix-openapi.ts --input api.json --output api-fixed.json
  tsx fix-openapi.ts --quiet --location path
`);
    process.exit(0);
  }

  const options: Partial<FixerOptions> = {};

  for (let i = 0; i < args.length; i++) {
    switch (args[i]) {
      case '--input':
        options.inputFile = args[++i];
        break;
      case '--output':
        options.outputFile = args[++i];
        break;
      case '--location':
        const loc = args[++i] as any;
        if (['query', 'header', 'path', 'cookie'].includes(loc)) {
          options.defaultParameterLocation = loc;
        }
        break;
      case '--quiet':
        options.verbose = false;
        break;
    }
  }

  const fixer = new OpenAPIFixer(options);
  await fixer.fix();
}

// Run if called directly
if (require.main === module) {
  main().catch(console.error);
}

export { OpenAPIFixer };
