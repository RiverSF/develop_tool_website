# MakeTransparent 使用说明

自动识别图片背景色，将该颜色（及接近色）转为透明，裁切有效内容后，按指定尺寸居中输出 PNG。

适用于纯色背景 Logo（黑、白、灰、或其他任意纯色底）。

## 目录结构

```
image_processing_tool/
├── MakeTransparent.exe   # 可执行工具
├── MakeTransparent.cs    # 源码（可重新编译）
├── README.md             # 本说明
└── output/               # 处理结果
    ├── river_logo_1800x600.png
    ├── river_logo_1800x800.png
    ├── river_logo_1920x1080.png
    ├── river_logo_2100x700.png
    └── river_logo_2100x800.png
```

## 依赖

- Windows
- .NET Framework（自带 `System.Drawing`，一般无需额外安装）

## 用法

```text
MakeTransparent <input> <output> [outW] [outH] [colorThresh]
```

| 参数 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `input` | 是 | 输入图片路径（建议 PNG/JPG） | — |
| `output` | 是 | 输出 PNG 路径 | — |
| `outW` | 否 | 输出宽度（像素） | `1920` |
| `outH` | 否 | 输出高度（像素） | `1080` |
| `colorThresh` | 否 | 与背景色的曼哈顿距离阈值 `|ΔR|+|ΔG|+|ΔB|` | `28` |

## 处理逻辑（简要）

1. 从四角附近采样，取中位 RGB 作为背景色
2. 将与背景色距离 ≤ 阈值的像素转为透明；边缘做软过渡，减轻锯齿色边
3. 按不透明内容计算包围盒并裁切
4. 等比缩放，四周留约 8% 边距，居中放到目标画布
5. 保存为带 Alpha 的 PNG

## 示例

在 `image_processing_tool` 目录下执行：

```powershell
# 默认输出 1920x1080（自动识别背景色）
.\MakeTransparent.exe .\input.png .\output\logo_1920x1080.png

# 指定尺寸
.\MakeTransparent.exe .\input.png .\output\logo_1800x600.png 1800 600
.\MakeTransparent.exe .\input.png .\output\logo_2100x800.png 2100 800

# 背景去得不干净时，可适当提高阈值（例如 40）
.\MakeTransparent.exe .\input.png .\output\logo.png 1920 1080 40

# 主体颜色接近背景、被误伤时，可降低阈值（例如 12）
.\MakeTransparent.exe .\input.png .\output\logo.png 1920 1080 12
```

## 重新编译

若修改了 `MakeTransparent.cs`：

```powershell
C:\Windows\Microsoft.NET\Framework64\v4.0.30319\csc.exe /nologo /t:exe /out:MakeTransparent.exe MakeTransparent.cs
```

（路径以本机已安装的 .NET Framework 为准。）

## 注意

- 适合**纯色背景**（黑/白/灰/其他单色底）。渐变、纹理、多色背景无法可靠去除。
- 主体颜色若与背景非常接近，可能被误判为透明；可调低 `colorThresh`。
- 部分看图软件会把透明区域显示成黑色，属预览行为，文件本身已是透明 PNG。
- 输出尺寸与原图比例差异较大时，Logo 会等比缩放并居中，不会拉伸变形。
