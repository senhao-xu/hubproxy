<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { cn } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    variant?: 'default' | 'secondary' | 'outline' | 'ghost'
    size?: 'default' | 'sm' | 'icon'
    disabled?: boolean
    class?: HTMLAttributes['class']
  }>(),
  {
    variant: 'default',
    size: 'default',
  },
)

const variants: Record<NonNullable<typeof props.variant>, string> = {
  default: 'border border-primary bg-primary text-primary-foreground shadow-sm shadow-primary/15 hover:bg-primary/90',
  secondary: 'border border-border bg-secondary text-secondary-foreground hover:border-input hover:bg-secondary/75',
  outline: 'border border-input bg-surface text-foreground hover:border-primary/45 hover:bg-accent hover:text-accent-foreground',
  ghost: 'border border-transparent text-muted-foreground hover:bg-accent hover:text-accent-foreground',
}

const sizes: Record<NonNullable<typeof props.size>, string> = {
  default: 'h-10 px-4 text-sm',
  sm: 'h-8 px-3 text-xs',
  icon: 'size-10',
}
</script>

<template>
  <button
    type="button"
    :disabled="disabled"
    :class="cn(
      'inline-flex shrink-0 items-center justify-center gap-1.5 rounded-md font-semibold outline-none transition-[border-color,background-color,color,box-shadow] duration-150 focus-visible:ring-2 focus-visible:ring-ring/60 focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:pointer-events-none disabled:opacity-45',
      variants[variant],
      sizes[size],
      props.class,
    )"
  >
    <slot />
  </button>
</template>
