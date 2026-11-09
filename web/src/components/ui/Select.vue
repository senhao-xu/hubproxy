<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const model = defineModel<string>({ default: '' })

const props = defineProps<{
  class?: HTMLAttributes['class']
  id?: string
  options: Array<{ value: string; label: string }>
  disabled?: boolean
  placeholder?: string
}>()
</script>

<template>
  <select
    :id="id"
    v-model="model"
    :disabled="disabled"
    :class="cn(
      'h-10 w-full min-w-0 rounded-md border border-input bg-surface-raised px-3 text-sm outline-none transition-[border-color,box-shadow] duration-150 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/35 disabled:cursor-not-allowed disabled:bg-muted disabled:opacity-60',
      props.class,
    )"
  >
    <option v-if="placeholder !== undefined" value="">{{ placeholder }}</option>
    <option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
  </select>
</template>
