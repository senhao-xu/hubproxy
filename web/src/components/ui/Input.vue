<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const model = defineModel<string>({ default: '' })

const props = defineProps<{
  class?: HTMLAttributes['class']
  type?: string
  id?: string
  placeholder?: string
  disabled?: boolean
}>()
</script>

<template>
  <input
    :id="id"
    v-model="model"
    :type="type || 'text'"
    :placeholder="placeholder"
    :disabled="disabled"
    :class="cn(
      'h-10 w-full min-w-0 rounded-md border border-input bg-surface-raised px-3 text-sm outline-none transition-[border-color,box-shadow] duration-150 placeholder:text-muted-foreground/80 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/35 disabled:cursor-not-allowed disabled:bg-muted disabled:opacity-60',
      props.class,
    )"
  >
</template>
