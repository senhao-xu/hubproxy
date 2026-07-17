<script setup lang="ts">
import { computed, ref } from 'vue'
import { Check, Clipboard, Container, Link2, Rocket, Sparkles, Terminal } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import PageHero from '@/components/PageHero.vue'
import { copyText } from '@/lib/utils'

const input = ref('')
const output = ref('')
const error = ref('')
const copied = ref(false)

const host = computed(() => window.location.host)

const features = [
  { icon: Rocket, label: 'GitHub 加速' },
  { icon: Container, label: 'Docker 镜像' },
  { icon: Sparkles, label: 'Hugging Face' },
] as const

const dockerExamples = computed(() => [
  {
    id: 'official',
    label: '官方镜像',
    original: 'docker pull nginx',
    accelerated: `docker pull ${host.value}/nginx`,
  },
  {
    id: 'user',
    label: '用户镜像',
    original: 'docker pull user/app:tag',
    accelerated: `docker pull ${host.value}/user/app:tag`,
  },
  {
    id: 'ghcr',
    label: 'GHCR',
    original: 'docker pull ghcr.io/org/app',
    accelerated: `docker pull ${host.value}/ghcr.io/org/app`,
  },
])

const allowedHosts = [
  'github.com/',
  'raw.githubusercontent.com/',
  'gist.githubusercontent.com/',
  'huggingface.co/',
  'cdn-lfs.hf.co/',
]

function formatLink() {
  error.value = ''
  copied.value = false
  const link = input.value.trim()
  if (!link) {
    error.value = '请输入有效的链接'
    output.value = ''
    return
  }

  if (link.startsWith('https://') || link.startsWith('http://')) {
    output.value = `https://${host.value}/${link}`
    return
  }

  if (allowedHosts.some((prefix) => link.startsWith(prefix))) {
    output.value = `https://${host.value}/https://${link}`
    return
  }

  error.value = '请输入有效的 GitHub / Hugging Face 链接'
  output.value = ''
}

async function onCopy() {
  if (!output.value) return
  copied.value = await copyText(output.value)
}

function onOpen() {
  if (!output.value) return
  window.open(output.value, '_blank', 'noopener,noreferrer')
}
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <PageHero
      eyebrow="面向开发者和运维人员的加速服务"
      title="HubProxy"
      subtitle="GitHub 文件加速 · Docker 镜像加速 · Hugging Face 资源"
      gradient
    >
      <div class="flex flex-wrap gap-2">
        <span
          v-for="item in features"
          :key="item.label"
          class="feature-pill"
        >
          <component :is="item.icon" class="size-4" />
          {{ item.label }}
        </span>
      </div>
    </PageHero>

    <section class="workspace-panel">
      <div class="panel-header">
        <div class="panel-heading"><Link2 class="size-4 text-primary" /><span>链接生成器</span></div>
        <span class="status-badge font-mono">HTTPS PROXY</span>
      </div>
      <div class="panel-body field-block">
        <div>
          <label for="source-link" class="field-label">原始资源链接</label>
          <div class="flex flex-col gap-2 sm:flex-row">
            <Input
              id="source-link"
              v-model="input"
              class="sm:flex-1"
              placeholder="粘贴 GitHub / Hugging Face 原始链接"
              @keyup.enter="formatLink"
            />
            <Button @click="formatLink"><Rocket class="size-4" />获取加速链接</Button>
          </div>
          <span class="field-hint">支持 github.com、raw.githubusercontent.com、Hugging Face 和 Git LFS 资源。</span>
        </div>

        <Transition name="fade" mode="out-in">
          <p v-if="error" key="error" role="alert" class="status-message status-error">{{ error }}</p>
          <div v-else-if="output" key="output" class="space-y-3">
            <div role="status" class="status-message status-success items-center"><Check class="size-4 shrink-0" />加速链接已生成，可直接复制或打开。</div>
            <p class="code-block break-all">{{ output }}</p>
            <div class="flex flex-wrap gap-2">
              <Button variant="secondary" size="sm" @click="onCopy">
                <Clipboard class="size-4" />
                {{ copied ? '已复制' : '复制链接' }}
              </Button>
              <Button variant="outline" size="sm" @click="onOpen">
                <Link2 class="size-4" />
                打开链接
              </Button>
            </div>
          </div>
        </Transition>
      </div>
    </section>

    <section class="section-gap space-y-4">
      <div class="flex flex-col gap-1 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p class="eyebrow">Docker Pull</p>
          <h2 class="text-xl font-semibold">镜像加速命令</h2>
        </div>
        <p class="text-sm text-muted-foreground">在镜像名前加入当前站点域名。</p>
      </div>

      <div class="terminal-block">
        <div class="terminal-header">
          <span class="terminal-dot" />
          <span class="terminal-dot" />
          <span class="terminal-dot" />
          <Terminal class="ml-2 size-3.5 text-primary" />
          <span class="font-mono text-xs text-muted-foreground">shell / docker pull</span>
        </div>
        <div class="terminal-body">
          <div
            v-for="item in dockerExamples"
            :key="item.id"
            class="terminal-example"
          >
            <span class="example-tag">{{ item.label }}</span>
            <p class="break-all font-mono text-sm leading-relaxed">
              <span class="text-muted-foreground">$ </span>
              <span class="text-muted-foreground/70 line-through decoration-muted-foreground/40">{{ item.original }}</span>
            </p>
            <p class="break-all font-mono text-sm leading-relaxed">
              <span class="text-muted-foreground">$ </span>
              <span class="text-primary">{{ item.accelerated }}</span>
            </p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>
