<template>
  <div class="schema-form">
    <div v-for="(schema, key) in properties" :key="key" class="mb-4">
      <div class="flex items-center justify-between mb-1">
        <span class="text-sm font-bold text-regular">
          {{ key }}
          <span v-if="isRequired(String(key))" class="text-danger">*</span>
        </span>
        <span class="text-xs text-secondary">{{ schema.type }}</span>
      </div>
      <div class="text-xs text-placeholder mb-2" v-if="schema.description">
        {{ schema.description }}
      </div>

      <!-- String / Enum -->
      <template v-if="schema.type === 'string'">
        <el-select
          v-if="schema.enum"
          v-model="model[key]"
          class="w-full"
          clearable
          placeholder="Select"
        >
          <el-option v-for="opt in schema.enum" :key="opt" :label="opt" :value="opt" />
        </el-select>
        <el-input v-else v-model="model[key]" :placeholder="String(key)" />
      </template>

      <!-- Boolean -->
      <template v-else-if="schema.type === 'boolean'">
        <el-switch v-model="model[key]" />
      </template>

      <!-- Number / Integer -->
      <template v-else-if="schema.type === 'number' || schema.type === 'integer'">
        <el-input-number
          v-model="model[key]"
          class="!w-full"
          controls-position="right"
          :placeholder="String(key)"
        />
      </template>

      <!-- Object (Recursive) -->
      <template v-else-if="schema.type === 'object' && schema.properties">
        <div class="pl-4 border-l-2 border-stone-200 dark:border-stone-700 py-2">
          <!-- Initialize object if null/undefined is handled by parent or init logic?
               We need to make sure model[key] exists. -->
          <div v-if="model[key] && typeof model[key] === 'object'">
            <SchemaForm v-model="model[key]" :schema="schema" :root-schema="rootSchema" />
          </div>
          <div v-else class="text-xs text-secondary">
            Optional Object is null.
            <el-button link type="primary" size="small" @click="initObject(String(key))"
              >Initialize {{ key }}</el-button
            >
          </div>
        </div>
      </template>

      <!-- Array -->
      <template v-else-if="schema.type === 'array'">
        <div class="pl-4 border-l-2 border-stone-200 dark:border-stone-700 py-2">
          <template v-if="schema.items && schema.items.type === 'string'">
            <!-- Simple string array editor -->
            <el-select
              v-model="model[key]"
              multiple
              filterable
              allow-create
              default-first-option
              :reserve-keyword="false"
              placeholder="Enter strings"
              class="w-full"
            />
          </template>
          <template v-else>
            <el-input
              type="textarea"
              :model-value="JSON.stringify(model[key])"
              disabled
              placeholder="Complex array, please use JSON mode"
            />
          </template>
        </div>
      </template>

      <!-- Fallback -->
      <template v-else>
        <el-input
          type="textarea"
          :model-value="JSON.stringify(model[key])"
          disabled
          placeholder="Complex type, please use JSON mode"
        />
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { isPropertyRequired, resolveRef } from '@/utils/schemaUtils'

const props = defineProps<{
  modelValue: Record<string, any>
  schema: any
  rootSchema?: any
}>()

const emit = defineEmits(['update:modelValue'])

const model = computed({
  get: () => props.modelValue || {},
  set: (val) => emit('update:modelValue', val),
})

watch(
  () => props.modelValue,
  (val) => {
    emit('update:modelValue', val)
  },
  { deep: true },
)

const rootSchema = computed(() => props.rootSchema || props.schema)

const properties = computed(() => {
  if (!props.schema?.properties) return {}
  const propsMap: Record<string, any> = {}
  Object.entries(props.schema.properties).forEach(([key, val]: [string, any]) => {
    propsMap[key] = resolveRef(val, rootSchema.value)
  })
  return propsMap
})

const isRequired = (key: string) => {
  return isPropertyRequired(key, props.schema)
}

const initObject = (key: string) => {
  const newModel = { ...model.value }
  newModel[key] = {}
  emit('update:modelValue', newModel)
}
</script>

<script lang="ts">
export default {
  name: 'SchemaForm',
}
</script>

<style scoped>
.text-danger {
  color: var(--ep-color-danger);
}
</style>
