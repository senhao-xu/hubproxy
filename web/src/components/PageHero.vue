<script setup lang="ts">
defineProps<{
  eyebrow: string
  title: string
  subtitle: string
  gradient?: boolean
}>()
</script>

<template>
  <header class="page-hero">
    <div class="page-hero-copy">
      <p class="eyebrow">{{ eyebrow }}</p>
      <h1 class="display-title" :class="{ 'gradient-text': gradient }">{{ title }}</h1>
      <p class="max-w-2xl text-sm leading-6 text-muted-foreground sm:text-base">
        {{ subtitle }}
      </p>
    </div>
    <div v-if="$slots.default" class="flex shrink-0 flex-wrap items-center gap-2 sm:justify-end">
      <slot />
    </div>
  </header>
</template>
