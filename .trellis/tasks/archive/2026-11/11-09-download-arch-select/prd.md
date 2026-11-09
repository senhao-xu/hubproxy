# 离线下载架构改为真实架构下拉选择

## Goal

离线镜像页（`web/src/pages/ImagesPage.vue`）的"目标架构"由手填文本框改为下拉选择，且下拉框选项必须来自镜像真实存在的平台列表（registry manifest index），杜绝用户填入不存在的架构后被后端静默回退到第一个 manifest（`src/handlers/imagetar.go:677`）的问题。

## Background

- 当前单镜像与批量下载的"目标架构（可选）"均为自由文本 `Input`，默认 `linux/amd64`，用户可输入任意不存在的架构。
- 后端 `handleImageInfo`（`src/handlers/imagetar.go:936`）已返回 `data.platforms` 与 `multiArch`，但前端 `ImageInfoResponse` 类型只声明了 `{ success }`，未消费该数据。
- 后端 platforms 输出缺少 variant（如 `linux/arm/v7` 会被截断为 `linux/arm` 且重复），需一并修正。

## Requirements

1. **单镜像下载**：镜像名输入后（防抖）自动调用 `/api/image/info` 获取平台列表；"目标架构"渲染为下拉框，选项为 `自动选择`（空值）+ 接口返回的真实平台。
2. **真实性约束**：下拉框只出现镜像真实存在的架构；信息获取失败时仅保留"自动选择"并展示失败提示，不提供任何虚构架构选项。
3. **选中值校验**：镜像变更后若当前选中平台不在新列表中，自动重置为"自动选择"。
4. **variant 完整**：平台选项保留 variant 三段格式（如 `linux/arm/v7`），与下载侧 `selectPlatformImage` 的匹配逻辑（os/arch/variant 拆分）兼容。
5. **批量下载**：对所有镜像的平台列表取**并集**作为选项；单个镜像获取失败不影响其它镜像；并集为空时仅保留"自动选择"。
6. **行为兼容**：选"自动选择"时不传 `platform` 参数（与现状一致）；压缩层开关、prepare/token 下载流程不变。
7. 批量下载原有手填文本框**移除**，同样使用下拉框。

## Non-Goals

- 不新增独立后端接口（复用并增强 `/api/image/info`）。
- 不改动 `formatArchs`（Tags 页数据来源不同，互不影响）。
- 不做平台选项的多语言/收藏等增强。

## Acceptance Criteria

- [ ] 输入 `nginx` 后，单镜像下拉框出现该镜像真实平台（含 `linux/amd64` 等），无重复项。
- [ ] 选择含 variant 平台（如 `linux/arm/v7`）下载成功，且产物 manifest 平台与所选一致。
- [ ] 输入不存在/无权限镜像时，下拉框仅剩"自动选择"，页面出现失败提示，提交走原有错误路径。
- [ ] 批量下载输入多镜像时，下拉框选项为各镜像平台并集。
- [ ] 选"自动选择"时请求不带 `platform` 参数；选中具体平台时请求带完整平台串。
- [ ] `go build ./...`、`go test ./src/handlers/`、`cd web && npm run build`（含 vue-tsc 类型检查）全部通过。
