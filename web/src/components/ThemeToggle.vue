<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { Moon, Sun } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { cn } from '@/lib/utils'

defineProps<{
  isDark: boolean
  buttonClass?: HTMLAttributes['class']
}>()

const emit = defineEmits<{ toggle: [] }>()
</script>

<template>
  <Button
    variant="ghost"
    size="icon"
    :aria-label="isDark ? '切换到亮色主题' : '切换到暗色主题'"
    :title="isDark ? '切换到亮色主题' : '切换到暗色主题'"
    :class="cn(buttonClass)"
    @click="emit('toggle')"
  >
    <Transition name="fade" mode="out-in">
      <Sun v-if="isDark" key="sun" class="size-4" />
      <Moon v-else key="moon" class="size-4" />
    </Transition>
  </Button>
</template>
