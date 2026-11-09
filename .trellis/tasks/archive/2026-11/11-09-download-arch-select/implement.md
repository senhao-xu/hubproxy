# 执行计划：离线下载架构下拉选择

## Checklist

1. [ ] 后端：`src/handlers/imagetar.go` `handleImageInfo` platforms 输出含 variant、去重、排序
2. [ ] 后端测试：为 platforms 构建补一个单测或表驱动用例（多 variant 去重、排序、非 index 镜像为空）
3. [ ] 前端类型：`web/src/api.ts` 扩展 `ImageInfoResponse.data`（platforms/multiArch 等），`fetchImageInfo` 支持可选 `signal`
4. [ ] 前端原语：新增 `web/src/components/ui/Select.vue`（native select、defineModel、fallthrough、cn 合并样式）
5. [ ] 页面改造：`web/src/pages/ImagesPage.vue` 单镜像 + 批量接入下拉框（防抖查询、abort、并集、失败降级、选中值重置）
6. [ ] 移除原"目标架构（可选）"两处 `Input`，标签改为"目标架构"，保留"自动选择"默认项

## 验证命令

```bash
go build ./...
go test ./src/handlers/
cd web && npm run build
```

## 审查门

- 步骤 1-2 完成后：`curl "/api/image/info?image=nginx"`（或本地启动）确认返回去重排序后的 platforms，含 `linux/arm/v7` 类三段平台。
- 步骤 4-6 完成后：`npm run dev` 手动验证：输入 nginx → 下拉出现真实平台；输入乱码镜像 → 仅"自动选择"+提示；批量多镜像 → 并集。

## 回滚点

- 每个步骤独立可回退；步骤 5、6 前端可整体还原 `ImagesPage.vue`。
