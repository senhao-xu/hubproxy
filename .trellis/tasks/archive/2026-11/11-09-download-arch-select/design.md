# 技术设计：离线下载架构下拉选择

## 总体方案

复用 `/api/image/info` 作为平台列表唯一来源，后端修正 platforms 输出，前端新增 `Select` UI 原语并在 `ImagesPage.vue` 组合"输入 → 防抖查询 → 下拉渲染"逻辑。不新增接口、不改下载流程。

## 后端（src/handlers/imagetar.go）

`handleImageInfo` 中 platforms 构建逻辑（imagetar.go:980-987）调整：

- 平台串包含 variant：variant 非空时输出 `os/arch/variant`，否则 `os/arch`。
- 用 `map[string]struct{}` 去重后排序（`sort.Strings`），保证选项稳定无重复。
- 响应字段名保持 `platforms` / `multiArch` 不变，旧响应结构 `{success, data}` 不变 —— 前端向后兼容。

与下载侧匹配逻辑的兼容性：`selectPlatformImage`（imagetar.go:654-666）按 `/` 拆分为 os/arch/variant 三段匹配，三段平台串可直接透传；`prepareSingleDownload` 前端已对空值跳过 `platform` 参数（api.ts:90）。

## 前端

### 类型（web/src/api.ts）

```ts
export interface ImageInfoData {
  name?: string
  mediaType?: string
  digest?: string
  size?: number
  multiArch?: boolean
  platforms?: string[]
}
export interface ImageInfoResponse {
  success: boolean
  data?: ImageInfoData
}
```

`getJSON` 已支持透传 `init`，`fetchImageInfo` 增加可选 `signal`，供竞态取消。

### 新 UI 原语（web/src/components/ui/Select.vue）

遵循 component-guidelines：

- 原生 `<select>` + `defineModel<string>`，单根元素保留 attribute fallthrough（`aria-label` 由调用方传）。
- 样式对齐 `Input.vue`（读其 class 作为基线），用 `cn()` 合并；使用 style.css 语义 token，浅深色都要可见。
- Props：`options: Array<{ value: string; label: string }>`、`disabled?: boolean`、`placeholder?: string`（对应 value 为 `''` 的首选项）。
- 保持原生元素与 `focus-visible` 样式，不自绘下拉面板。

### 页面组合（web/src/pages/ImagesPage.vue）

业务状态与请求逻辑留在页面组件（spec 要求 shared 组件只管展示）：

- `singlePlatform` / `batchPlatform` 语义不变：`''` 表示自动选择；提交时为空即不传 platform。
- 单镜像：`watch(singleImage)` + 500ms 防抖 + `AbortController` 取消旧请求；成功后存 `singlePlatforms`，失败清空并置 `singleArchError` 提示。
- 批量：`watch(batchText)` 防抖后对去重镜像列表 `Promise.allSettled` 并发请求，成功结果取并集去重排序存 `batchPlatforms`；失败镜像忽略（提交时 preflight 仍会兜底校验）。
- 列表更新时若当前选中值不在列表且非 `''`，重置为 `''`（自动选择）。
- 请求期间下拉框 `disabled`，提示文案"正在获取架构..."；无列表时仅渲染"自动选择"一项。
- `preflight` 保持不变（提交前二次校验存在性与权限）。

## 边界与取舍

- **竞态**：以 abort + 请求序号丢弃过期响应，避免旧镜像的平台列表污染新下拉框。
- **防抖**：500ms，输入停顿即查询；批量在每次文本变更后整体重查（镜像数通常少，成本可接受）。
- **失败降级**：不阻塞下载，仅收窄为"自动选择"，与现状（后端自动选择）等价。
- **不缓存**：镜像 tag 可能更新，每次输入变更重新拉取，保持平台列表新鲜。

## 回滚

改动面为：后端单函数内逻辑、api.ts 类型、新增 Select.vue、ImagesPage.vue。回滚即还原对应文件并删除 Select.vue，无数据/接口兼容性负担。
