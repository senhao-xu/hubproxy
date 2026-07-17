<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Container, Github, Menu, Rocket, Search, X, Zap } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import ThemeToggle from '@/components/ThemeToggle.vue'

const STORAGE_KEY = 'theme'
const route = useRoute()
const isDark = ref(false)
const menuOpen = ref(false)

const links = [
  { to: '/', label: 'GitHub 加速', icon: Rocket },
  { to: '/images', label: '离线镜像', icon: Container },
  { to: '/search', label: '镜像搜索', icon: Search },
] as const

const currentPath = computed(() => route.path)

function applyTheme(dark: boolean) {
  isDark.value = dark
  document.documentElement.classList.toggle('dark', dark)
  localStorage.setItem(STORAGE_KEY, dark ? 'dark' : 'light')
}

function toggleTheme() {
  applyTheme(!isDark.value)
}

function closeMenu() {
  menuOpen.value = false
}

onMounted(() => {
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'dark' || saved === 'light') {
    applyTheme(saved === 'dark')
  } else {
    applyTheme(window.matchMedia('(prefers-color-scheme: dark)').matches)
  }
})
</script>

<template>
  <div class="shell-atmosphere flex min-h-screen flex-col text-foreground">
    <header class="sticky top-0 z-50 border-b border-border bg-surface/95 backdrop-blur-md">
      <div class="mx-auto flex h-15 max-w-7xl items-center justify-between gap-3 px-4 sm:px-6 lg:px-8">
        <RouterLink
          to="/"
          class="flex items-center gap-2.5 rounded-md font-display text-base font-semibold transition-colors hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          @click="closeMenu"
        >
          <span class="brand-mark flex size-8 items-center justify-center rounded-md">
            <Zap class="size-4" />
          </span>
          <span>HubProxy</span>
          <span class="hidden font-mono text-[10px] font-medium text-muted-foreground sm:inline">CONSOLE</span>
        </RouterLink>

        <nav class="hidden items-center gap-1 md:flex">
          <RouterLink
            v-for="link in links"
            :key="link.to"
            :to="link.to"
            class="inline-flex h-9 items-center gap-1.5 rounded-md border px-3 text-sm font-medium transition-colors duration-150"
            :class="currentPath === link.to ? 'border-primary/30 bg-primary/10 text-primary' : 'border-transparent text-muted-foreground hover:bg-accent hover:text-foreground'"
          >
            <component :is="link.icon" class="size-4" />
            {{ link.label }}
          </RouterLink>
          <ThemeToggle :is-dark="isDark" button-class="ml-1" @toggle="toggleTheme" />
          <a
            href="https://github.com/sky22333/hubproxy"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="在 GitHub 查看 HubProxy"
            title="GitHub"
            class="ml-1 inline-flex size-10 items-center justify-center rounded-md border border-transparent text-muted-foreground transition-colors hover:bg-accent hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <Github class="size-4" />
          </a>
        </nav>

        <div class="flex items-center gap-0.5 md:hidden">
          <ThemeToggle :is-dark="isDark" @toggle="toggleTheme" />
          <Button
            variant="ghost"
            size="icon"
            :aria-label="menuOpen ? '关闭菜单' : '打开菜单'"
            :aria-expanded="menuOpen"
            aria-controls="mobile-navigation"
            @click="menuOpen = !menuOpen"
          >
            <Transition name="fade" mode="out-in">
              <X v-if="menuOpen" key="x" class="size-4" />
              <Menu v-else key="menu" class="size-4" />
            </Transition>
          </Button>
        </div>
      </div>

      <Transition name="menu">
        <div v-if="menuOpen" class="border-t border-border bg-surface px-4 py-2 md:hidden">
          <nav id="mobile-navigation" aria-label="移动端主导航" class="mx-auto flex max-w-7xl flex-col gap-1">
            <RouterLink
              v-for="link in links"
              :key="link.to"
              :to="link.to"
              class="inline-flex min-h-10 items-center gap-2 rounded-md border px-3 text-sm font-medium transition-colors"
              :class="currentPath === link.to ? 'border-primary/30 bg-primary/10 text-primary' : 'border-transparent text-muted-foreground hover:bg-accent'"
              @click="closeMenu"
            >
              <component :is="link.icon" class="size-4" />
              {{ link.label }}
            </RouterLink>
          </nav>
        </div>
      </Transition>
    </header>

    <main class="mx-auto w-full max-w-7xl flex-1 px-4 py-7 text-base sm:px-6 sm:py-9 lg:px-8">
      <slot />
    </main>

    <footer class="border-t border-border bg-surface/70 px-4 py-5">
      <div class="mx-auto flex max-w-7xl items-center justify-between gap-4 text-xs text-muted-foreground">
        <span class="font-mono">HubProxy / developer utility</span>
        <a
          href="https://github.com/sky22333/hubproxy"
          target="_blank"
          rel="noopener noreferrer"
          aria-label="在 GitHub 查看 HubProxy"
          class="inline-flex items-center gap-2 rounded-sm text-muted-foreground transition-colors duration-150 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <Github class="size-4" />
          <span>GitHub</span>
        </a>
      </div>
    </footer>
  </div>
</template>
