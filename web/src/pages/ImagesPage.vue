<script setup lang="ts">
import { ref } from 'vue'
import { Archive, Boxes, Loader2 } from 'lucide-vue-next'
import {
  fetchImageInfo,
  prepareBatchDownload,
  prepareSingleDownload,
  triggerDownload,
} from '@/api'
import { errorMessage } from '@/lib/utils'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import PageHero from '@/components/PageHero.vue'
import Switch from '@/components/ui/Switch.vue'
import Textarea from '@/components/ui/Textarea.vue'

const singleImage = ref('')
const singlePlatform = ref('linux/amd64')
const singleCompressed = ref(true)
const singleStatus = ref('')
const singleError = ref('')
const singleLoading = ref(false)

const batchText = ref('')
const batchPlatform = ref('linux/amd64')
const batchCompressed = ref(true)
const batchStatus = ref('')
const batchError = ref('')
const batchLoading = ref(false)

async function preflight(images: string[]) {
  for (const image of [...new Set(images)]) {
    await fetchImageInfo(image)
  }
}

async function onSingleSubmit() {
  singleError.value = ''
  singleStatus.value = ''
  const image = singleImage.value.trim()
  if (!image) {
    singleError.value = '请输入镜像名称'
    return
  }

  singleLoading.value = true
  singleStatus.value = '正在准备下载...'
  try {
    await preflight([image])
    const data = await prepareSingleDownload({
      image,
      platform: singlePlatform.value,
      compressed: singleCompressed.value,
    })
    if (!data.download_url) throw new Error('下载地址生成失败')
    triggerDownload(data.download_url)
    const platformText = singlePlatform.value.trim()
      ? ` (${singlePlatform.value.trim()})`
      : ''
    singleStatus.value = `开始下载 ${image}${platformText}`
  } catch (e) {
    singleStatus.value = ''
    singleError.value = errorMessage(e, '下载失败')
  } finally {
    singleLoading.value = false
  }
}

async function onBatchSubmit() {
  batchError.value = ''
  batchStatus.value = ''
  const images = batchText.value
    .split('\n')
    .map((line) => line.trim())
    .filter((line) => line && !line.startsWith('#'))

  if (images.length === 0) {
    batchError.value = '请输入镜像列表'
    return
  }

  batchLoading.value = true
  batchStatus.value = '正在准备批量下载...'
  try {
    await preflight(images)
    const data = await prepareBatchDownload({
      images,
      platform: batchPlatform.value,
      useCompressedLayers: batchCompressed.value,
    })
    if (!data.download_url) throw new Error('下载地址生成失败')
    triggerDownload(data.download_url)
    batchStatus.value = `开始下载 ${images.length} 个镜像`
  } catch (e) {
    batchStatus.value = ''
    batchError.value = errorMessage(e, '下载失败')
  } finally {
    batchLoading.value = false
  }
}
</script>

<template>
  <div class="mx-auto max-w-6xl">
    <PageHero
      eyebrow="Offline Image"
      title="离线镜像"
      subtitle="流式下载，兼容 docker load，支持多架构。"
    />

    <div class="space-y-5">
      <section class="workspace-panel">
        <div class="panel-header">
          <div>
            <div class="panel-heading">
              <Archive class="size-4 text-primary" />
              单镜像下载
            </div>
            <p class="panel-description">下载一个可供 docker load 使用的离线镜像包。</p>
          </div>
          <span class="panel-index">01</span>
        </div>
        <div class="panel-body space-y-4">
          <Transition name="fade" mode="out-in">
            <p v-if="singleError" key="error" role="alert" class="status-message status-error">{{ singleError }}</p>
            <p v-else-if="singleStatus" key="status" role="status" class="status-message status-info">
              <Loader2 v-if="singleLoading" class="size-4 shrink-0 animate-spin" />
              {{ singleStatus }}
            </p>
          </Transition>
          <div class="grid min-w-0 items-end gap-4 md:grid-cols-2 lg:grid-cols-[minmax(0,2fr)_minmax(0,1fr)_minmax(11rem,0.8fr)_auto]">
            <label class="min-w-0">
              <span class="field-label">镜像名称</span>
              <Input v-model="singleImage" placeholder="nginx 或 user/app:tag" />
            </label>
            <label class="min-w-0">
              <span class="field-label">目标架构（可选）</span>
              <Input v-model="singlePlatform" placeholder="linux/amd64" />
            </label>
            <div class="min-w-0">
              <span class="field-label">压缩层</span>
              <div class="setting-row min-h-10 py-1.5">
                <p class="min-w-0 text-xs text-muted-foreground">减小下载包体积</p>
                <Switch v-model:checked="singleCompressed" aria-label="单镜像使用压缩层" />
              </div>
            </div>
            <Button class="w-full md:col-span-2 lg:col-span-1 lg:w-auto" :disabled="singleLoading" @click="onSingleSubmit">
              <Loader2 v-if="singleLoading" class="size-4 animate-spin" />
              {{ singleLoading ? '准备中...' : '立即下载' }}
            </Button>
          </div>
        </div>
      </section>

      <section class="workspace-panel">
        <div class="panel-header">
          <div>
            <div class="panel-heading">
              <Boxes class="size-4 text-info" />
              批量下载
            </div>
            <p class="panel-description">每行一个镜像，忽略以 # 开头的注释。</p>
          </div>
          <span class="panel-index">02</span>
        </div>
        <div class="panel-body space-y-4">
          <Transition name="fade" mode="out-in">
            <p v-if="batchError" key="error" role="alert" class="status-message status-error">{{ batchError }}</p>
            <p v-else-if="batchStatus" key="status" role="status" class="status-message status-info">
              <Loader2 v-if="batchLoading" class="size-4 shrink-0 animate-spin" />
              {{ batchStatus }}
            </p>
          </Transition>
          <div class="grid min-w-0 gap-5 lg:grid-cols-[minmax(0,1fr)_17rem] lg:items-start">
            <label class="min-w-0">
              <span class="field-label">镜像列表</span>
              <Textarea
                v-model="batchText"
                class="min-h-56 lg:min-h-64"
                placeholder="alpine&#10;redis:alpine&#10;user/app:1.0"
              />
            </label>
            <div class="field-block min-w-0">
              <label class="min-w-0">
                <span class="field-label">目标架构（可选）</span>
                <Input v-model="batchPlatform" placeholder="linux/amd64" />
              </label>
              <div class="min-w-0">
                <span class="field-label">压缩层</span>
                <div class="setting-row min-h-10 py-1.5">
                  <p class="min-w-0 text-xs text-muted-foreground">减小下载包体积</p>
                  <Switch v-model:checked="batchCompressed" aria-label="批量下载使用压缩层" />
                </div>
              </div>
              <Button class="w-full" :disabled="batchLoading" @click="onBatchSubmit">
                <Loader2 v-if="batchLoading" class="size-4 animate-spin" />
                {{ batchLoading ? '准备中...' : '批量下载' }}
              </Button>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>
