// src/utils/schemaUtils.ts

/**
 * Helper function to check if a property is required in a schema
 */
export function isPropertyRequired(propertyName: string, schema: any): boolean {
  return schema.required?.includes(propertyName) ?? false
}

/**
 * Generates a default value based on a JSON schema type
 */
export function generateDefaultValue(schema: any, propertyName?: string, parentSchema?: any): any {
  if ('default' in schema && schema.default !== undefined) {
    return schema.default
  }

  // Check if this property is required in the parent schema
  const isRequired =
    propertyName && parentSchema ? isPropertyRequired(propertyName, parentSchema) : false

  switch (schema.type) {
    case 'string':
      return isRequired ? '' : undefined
    case 'number':
    case 'integer':
      return isRequired ? 0 : undefined
    case 'boolean':
      return isRequired ? false : undefined
    case 'array':
      return isRequired ? [] : undefined
    case 'object': {
      if (!schema.properties) return isRequired ? {} : undefined

      const obj: any = {}
      Object.entries(schema.properties).forEach(([key, prop]) => {
        if (isPropertyRequired(key, schema)) {
          const value = generateDefaultValue(prop, key, schema)
          if (value !== undefined) {
            obj[key] = value
          }
        }
      })
      return isRequired ? obj : Object.keys(obj).length > 0 ? obj : undefined
    }
    case 'null':
      return null
    default:
      return undefined
  }
}

/**
 * Resolves $ref references in JSON schema
 */
export function resolveRef(schema: any, rootSchema: any): any {
  if (!('$ref' in schema) || !schema.$ref) {
    return schema
  }

  const ref = schema.$ref

  if (ref.startsWith('#/')) {
    const path = ref.substring(2).split('/')
    let current: any = rootSchema

    for (const segment of path) {
      if (current && typeof current === 'object' && current !== null && segment in current) {
        current = current[segment]
      } else {
        console.warn(`Could not resolve $ref: ${ref}`)
        return schema
      }
    }

    return current
  }

  console.warn(`Unsupported $ref format: ${ref}`)
  return schema
}

/**
 * Normalizes union types to simple types for form rendering
 */
export function normalizeUnionType(schema: any): any {
  // Handle anyOf with exactly string and null
  if (
    schema.anyOf &&
    schema.anyOf.length === 2 &&
    schema.anyOf.some((t: any) => t.type === 'string') &&
    schema.anyOf.some((t: any) => t.type === 'null')
  ) {
    return { ...schema, type: 'string', anyOf: undefined, nullable: true }
  }

  // Handle anyOf with exactly boolean and null
  if (
    schema.anyOf &&
    schema.anyOf.length === 2 &&
    schema.anyOf.some((t: any) => t.type === 'boolean') &&
    schema.anyOf.some((t: any) => t.type === 'null')
  ) {
    return { ...schema, type: 'boolean', anyOf: undefined, nullable: true }
  }

  // Handle anyOf with exactly number/integer and null
  if (
    schema.anyOf &&
    schema.anyOf.length === 2 &&
    (schema.anyOf.some((t: any) => t.type === 'number') ||
      schema.anyOf.some((t: any) => t.type === 'integer')) &&
    schema.anyOf.some((t: any) => t.type === 'null')
  ) {
    return { ...schema, type: 'number', anyOf: undefined, nullable: true }
  }

  // Handle anyOf with exactly array and null
  if (
    schema.anyOf &&
    schema.anyOf.length === 2 &&
    schema.anyOf.some((t: any) => t.type === 'array') &&
    schema.anyOf.some((t: any) => t.type === 'null')
  ) {
    return { ...schema, type: 'array', anyOf: undefined, nullable: true }
  }

  // Handle array type with null types (e.g. type: ["string", "null"])
  if (Array.isArray(schema.type) && schema.type.includes('null')) {
    const types = schema.type.filter((t: string) => t !== 'null')
    if (types.length === 1) {
      return { ...schema, type: types[0], nullable: true }
    }
  }

  return schema
}
