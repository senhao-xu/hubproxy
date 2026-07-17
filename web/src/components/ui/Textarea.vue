<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const model = defineModel<string>({ default: '' })

const props = defineProps<{
  class?: HTMLAttributes['class']
  id?: string
  placeholder?: string
  disabled?: boolean
}>()
</script>

<template>
  <textarea
    :id="id"
    v-model="model"
    :placeholder="placeholder"
    :disabled="disabled"
    :class="cn(
      'min-h-36 w-full min-w-0 resize-y rounded-md border border-input bg-surface-raised px-3 py-2.5 font-mono text-sm leading-6 outline-none transition-[border-color,box-shadow] duration-150 placeholder:text-muted-foreground/80 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/35 disabled:cursor-not-allowed disabled:bg-muted disabled:opacity-60',
      props.class,
    )"
  />
</template>
