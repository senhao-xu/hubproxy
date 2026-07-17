<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDown, ArrowLeft, ChevronLeft, ChevronRight, Copy, Loader2, Search, Star, Tags } from 'lucide-vue-next'
import {
  fetchTags,
  searchImages,
  type Repository,
  type TagInfo,
} from '@/api'
import { copyText, errorMessage, formatArchs, formatNumber, formatSize, formatTimeAgo } from '@/lib/utils'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import PageHero from '@/components/PageHero.vue'

interface RepoView {
  raw: Repository
  displayName: string
  namespace: string
  name: string
  fullRepoName: string
}

const route = useRoute()
const router = useRouter()

const query = ref('')
const searching = ref(false)
const searchError = ref('')
const results = ref<RepoView[]>([])
const resultCount = ref(0)
const resultsPage = ref(1)
const pageSize = 25

const selected = ref<RepoView | null>(null)
const tagsLoading = ref(false)
const tagsError = ref('')
const tags = ref<TagInfo[]>([])
const tagFilter = ref('')
const tagsPage = ref(1)
const tagsHasMore = ref(false)
const copyHint = ref('')

const host = computed(() => window.location.host)

const hasResults = computed(() => results.value.length > 0)
const totalPages = computed(() => Math.max(1, Math.ceil(resultCount.value / pageSize)))
const hasMoreResults = computed(() => resultsPage.value < totalPages.value)

const filteredTags = computed(() => {
  const q = tagFilter.value.trim().toLowerCase()
  if (!q) return tags.value
  const exact: TagInfo[] = []
  const starts: TagInfo[] = []
  const includes: TagInfo[] = []
  for (const tag of tags.value) {
    const name = tag.name.toLowerCase()
    if (name === q) exact.push(tag)
    else if (name.startsWith(q)) starts.push(tag)
    else if (name.includes(q)) includes.push(tag)
  }
  return [...exact, ...starts, ...includes]
})

const displayTags = computed(() =>
  filteredTags.value.map((tag) => ({
    tag,
    archs: formatArchs(tag.images),
    size: formatSize(tag.full_size),
  })),
)

function toRepoView(item: Repository): RepoView | null {
  const rawName = item.repo_name || ''
  const namespace =
    item.namespace ||
    (item.is_official ? 'library' : rawName.includes('/') ? rawName.split('/')[0] : '')
  const name = rawName.replace(/^library\//, '').includes('/')
    ? rawName.split('/').pop() || ''
    : rawName.replace(/^library\//, '')

  if (!namespace || !name) return null

  const displayName = item.is_official
    ? name
    : item.namespace
      ? `${item.namespace}/${name}`
      : rawName.includes('/')
        ? rawName
        : `${namespace}/${name}`

  return {
    raw: item,
    displayName,
    namespace,
    name,
    fullRepoName: item.is_official ? name : `${namespace}/${name}`,
  }
}

async function runSearch(q: string, page = 1) {
  const trimmed = q.trim()
  if (!trimmed) {
    searchError.value = '请输入搜索关键词'
    return
  }

  searching.value = true
  searchError.value = ''
  results.value = []
  selected.value = null
  tags.value = []
  tagsError.value = ''
  tagFilter.value = ''

  try {
    let searchQuery = trimmed
    let targetRepo = ''
    if (trimmed.includes('/')) {
      const [ns] = trimmed.split('/')
      searchQuery = ns
      targetRepo = trimmed.toLowerCase()
    }

    const data = await searchImages(searchQuery, page, pageSize)
    const views = (data.results || [])
      .map(toRepoView)
      .filter((v): v is RepoView => v !== null)

    views.sort((a, b) => {
      if (targetRepo) {
        const aMatch =
          a.displayName.toLowerCase() === targetRepo ||
          a.fullRepoName.toLowerCase() === targetRepo
        const bMatch =
          b.displayName.toLowerCase() === targetRepo ||
          b.fullRepoName.toLowerCase() === targetRepo
        if (aMatch && !bMatch) return -1
        if (!aMatch && bMatch) return 1
      }
      if (!!a.raw.is_official !== !!b.raw.is_official) {
        return a.raw.is_official ? -1 : 1
      }
      return (b.raw.pull_count || 0) - (a.raw.pull_count || 0)
    })

    results.value = views
    resultCount.value = data.count ?? views.length
    resultsPage.value = page
    if (views.length === 0) searchError.value = '未找到相关镜像'
  } catch (e) {
    searchError.value = errorMessage(e, '搜索失败')
  } finally {
    searching.value = false
  }
}

async function onSearch() {
  const q = query.value.trim()
  await router.replace({ path: '/search', query: q ? { q } : {} })
  await runSearch(q, 1)
}

async function loadResultsPage(page: number) {
  if (page < 1 || page > totalPages.value || searching.value) return
  await runSearch(query.value, page)
  window.scrollTo(0, 0)
}

async function loadTagPage(repo: RepoView, page: number) {
  selected.value = repo
  tagsLoading.value = true
  tagsError.value = ''
  tagsPage.value = page
  if (page === 1) tagFilter.value = ''

  try {
    const data = await fetchTags(repo.namespace, repo.name, page, 100)
    tags.value = data.tags || []
    tagsHasMore.value = !!data.has_more
    if (tags.value.length === 0) tagsError.value = '该镜像暂无可用标签'
  } catch (e) {
    tags.value = []
    tagsError.value = errorMessage(e, '加载标签失败')
  } finally {
    tagsLoading.value = false
  }
}

function backToResults() {
  selected.value = null
  tags.value = []
  tagsError.value = ''
  tagFilter.value = ''
}

async function copyPull(tagName?: string) {
  if (!selected.value) return
  const image = tagName
    ? `${host.value}/${selected.value.fullRepoName}:${tagName}`
    : `${host.value}/${selected.value.fullRepoName}`
  const refName = `docker pull ${image}`
  const ok = await copyText(refName)
  copyHint.value = ok ? `已复制 ${refName}` : '复制失败'
  setTimeout(() => {
    if (copyHint.value.includes(refName) || copyHint.value === '复制失败') copyHint.value = ''
  }, 2000)
}

watch(
  () => route.query.q,
  async (q) => {
    const next = typeof q === 'string' ? q : ''
    if (next === query.value) return
    query.value = next
    if (next) await runSearch(next, 1)
  },
  { immediate: true },
)
</script>

<template>
  <div class="mx-auto max-w-6xl">
    <PageHero
      eyebrow="Docker Hub"
      title="镜像搜索"
      subtitle="检索官方与社区镜像，查看标签与架构，一键复制拉取命令。"
    />

    <Transition name="fade" mode="out-in">
      <div v-if="!selected" key="search" class="space-y-5">
        <section class="workspace-panel">
          <div class="panel-header">
            <div class="panel-heading"><Search class="size-4 text-primary" />镜像仓库检索</div>
            <span class="status-badge font-mono">DOCKER HUB</span>
          </div>
          <div class="panel-body">
            <div class="flex flex-col gap-2 sm:flex-row">
              <Input
                v-model="query"
                class="sm:flex-1"
                aria-label="搜索 Docker Hub 镜像"
                placeholder="例如 nginx、redis、library/ubuntu"
                @keydown.enter="onSearch"
              />
              <Button :disabled="searching" @click="onSearch">
                <Loader2 v-if="searching" class="size-4 animate-spin" />
                <Search v-else class="size-4" />
                {{ searching ? '搜索中...' : '搜索' }}
              </Button>
            </div>
          </div>
        </section>

        <p v-if="searchError" role="alert" class="status-message status-error">{{ searchError }}</p>

        <div v-if="searching" class="space-y-2">
          <div v-for="i in 3" :key="i" class="skeleton-row" />
        </div>

        <section v-else-if="hasResults" class="space-y-2">
          <div aria-live="polite" class="flex flex-wrap items-center justify-between gap-2 px-1 text-sm text-muted-foreground">
            <p>共 <span class="font-semibold text-foreground">{{ resultCount }}</span> 条结果 <template v-if="totalPages > 1">/ 第 {{ resultsPage }} / {{ totalPages }} 页</template></p>
            <p class="font-mono text-xs">SELECT A REPOSITORY TO INSPECT TAGS</p>
          </div>
          <div class="data-list">
            <button
              v-for="item in results"
              :key="`${item.namespace}/${item.name}`"
              type="button"
              class="data-row data-row-action w-full text-left"
              @click="loadTagPage(item, 1)"
            >
              <div class="flex min-w-0 flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                <div class="min-w-0">
                  <div class="mb-1 flex flex-wrap items-center gap-2">
                    <span class="break-all font-mono text-sm font-semibold text-foreground">{{ item.displayName }}</span>
                    <span
                      v-if="item.raw.is_official"
                      class="status-badge min-h-0 border-primary/25 bg-primary/8 py-0.5 text-[11px] text-primary"
                    >官方</span>
                  </div>
                  <p class="line-clamp-2 [overflow-wrap:anywhere] text-sm leading-6 text-muted-foreground">{{ item.raw.short_description || '暂无描述' }}</p>
                </div>
                <div class="flex shrink-0 gap-3 font-mono text-xs text-muted-foreground">
                  <span v-if="item.raw.star_count" class="inline-flex items-center gap-1"><Star class="size-3" />{{ formatNumber(item.raw.star_count) }}</span>
                  <span v-if="item.raw.pull_count" class="inline-flex items-center gap-1"><ArrowDown class="size-3" />{{ formatNumber(item.raw.pull_count) }}</span>
                </div>
              </div>
            </button>
          </div>
          <div v-if="totalPages > 1" class="flex items-center justify-center gap-1.5 pt-2">
            <Button
              variant="outline"
              size="sm"
              aria-label="上一页搜索结果"
              title="上一页"
              :disabled="searching || resultsPage <= 1"
              @click="loadResultsPage(resultsPage - 1)"
            >
              <ChevronLeft class="size-4" />
            </Button>
            <span class="min-w-14 text-center text-muted-foreground">第 {{ resultsPage }} 页</span>
            <Button
              variant="outline"
              size="sm"
              aria-label="下一页搜索结果"
              title="下一页"
              :disabled="searching || !hasMoreResults"
              @click="loadResultsPage(resultsPage + 1)"
            >
              <ChevronRight class="size-4" />
            </Button>
          </div>
        </section>
      </div>

      <div v-else key="tags" class="space-y-5">
        <button
          type="button"
          class="inline-flex items-center gap-1.5 text-sm font-medium text-muted-foreground transition-colors hover:text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          @click="backToResults"
        >
          <ArrowLeft class="size-4" />返回搜索结果
        </button>

        <section class="workspace-panel">
          <div class="panel-header">
            <div class="min-w-0 space-y-1">
              <div class="panel-heading">
                <Tags class="size-4 shrink-0 text-primary" />
                <h2 class="break-all font-mono text-base sm:text-lg">{{ selected.fullRepoName }}</h2>
                <span
                  v-if="selected.raw.is_official"
                  class="status-badge min-h-0 border-primary/25 bg-primary/8 py-0.5 text-[11px] text-primary"
                >官方</span>
              </div>
              <p class="panel-description">{{ selected.raw.short_description || '暂无描述' }}</p>
            </div>
            <Button variant="outline" size="sm" @click="copyPull()"><Copy class="size-4" />复制拉取命令</Button>
          </div>
          <div class="panel-body space-y-4">
            <Transition name="fade">
              <p
                v-if="copyHint"
                :role="copyHint === '复制失败' ? 'alert' : 'status'"
                class="status-message break-all"
                :class="copyHint === '复制失败' ? 'status-error' : 'status-success'"
              >
                {{ copyHint }}
              </p>
            </Transition>
            <div class="flex flex-col gap-2 lg:flex-row lg:items-center">
              <Input
                v-model="tagFilter"
                class="sm:flex-1"
                aria-label="筛选当前页标签"
                placeholder="筛选当前页标签..."
              />
              <div class="flex flex-wrap items-center gap-1.5">
                <Button
                  variant="outline"
                  size="sm"
                  aria-label="上一页标签"
                  title="上一页"
                  :disabled="tagsLoading || tagsPage <= 1"
                  @click="loadTagPage(selected, tagsPage - 1)"
                >
                  <ChevronLeft class="size-4" />
                </Button>
                <span class="min-w-14 text-center text-muted-foreground">第 {{ tagsPage }} 页</span>
                <Button
                  variant="outline"
                  size="sm"
                  aria-label="下一页标签"
                  title="下一页"
                  :disabled="tagsLoading || !tagsHasMore"
                  @click="loadTagPage(selected, tagsPage + 1)"
                >
                  <ChevronRight class="size-4" />
                </Button>
              </div>
            </div>
          </div>
        </section>

        <p v-if="tagsError" role="alert" class="status-message status-error">{{ tagsError }}</p>
        <div v-else-if="tagsLoading" class="space-y-3">
          <div v-for="i in 5" :key="i" class="skeleton-row" />
        </div>
        <p v-else-if="displayTags.length === 0" class="status-message status-info">没有匹配的标签</p>
        <div v-else class="data-list">
          <div
            v-for="{ tag, archs, size } in displayTags"
            :key="tag.name"
            class="data-row flex items-start justify-between gap-3"
          >
            <div class="min-w-0 space-y-1.5">
              <p class="break-all font-mono text-sm font-semibold">{{ tag.name }}</p>
              <p class="text-xs text-muted-foreground">
                <template v-if="size">{{ size }} · </template>
                {{ formatTimeAgo(tag.last_updated) }}
              </p>
              <div v-if="archs.length" class="flex flex-wrap gap-1.5">
                <span
                  v-for="arch in archs"
                  :key="arch"
                  class="rounded-sm border border-primary/20 bg-primary/8 px-2 py-0.5 font-mono text-[11px] text-primary"
                >{{ arch }}</span>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              class="shrink-0"
              :aria-label="`复制 ${tag.name} 拉取命令`"
              @click="copyPull(tag.name)"
            >
              <Copy class="size-4" />
              复制
            </Button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>
