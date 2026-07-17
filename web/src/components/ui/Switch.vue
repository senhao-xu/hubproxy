<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const checked = defineModel<boolean>('checked', { default: false })

const props = defineProps<{
  id?: string
  class?: HTMLAttributes['class']
  disabled?: boolean
}>()

function toggle() {
  if (props.disabled) return
  checked.value = !checked.value
}
</script>

<template>
  <button
    :id="id"
    type="button"
    role="switch"
    :aria-checked="checked"
    :disabled="disabled"
    :class="cn(
      'relative inline-flex h-6 w-10 shrink-0 items-center rounded-full border transition-colors duration-150 outline-none focus-visible:ring-2 focus-visible:ring-ring/60 focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:opacity-50',
      checked ? 'border-primary bg-primary' : 'border-input bg-muted',
      props.class,
    )"
    @click="toggle"
  >
    <span
      :class="cn(
        'pointer-events-none block size-[18px] rounded-full bg-surface-raised shadow-sm transition-transform duration-150 ease-out',
        checked ? 'translate-x-[18px]' : 'translate-x-0.5',
      )"
    />
  </button>
</template>
