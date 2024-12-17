# 网络工具箱需求文档

## 项目概述

网络工具箱是一个多功能的网络诊断和测试工具，主要提供网络测试、分析及结果输出功能。项目目标是开发一款易于使用的工具，帮助用户快速完成网络相关任务，并结合 AI 分析结果，提供智能化的建议与总结。

## 项目目标

- 提供高效、便捷的网络测试工具，集成多种常见网络操作。
- 支持扫描结果通过 AI 引擎进行分析，生成诊断报告和优化建议。
- 提供清晰的 CLI 交互体验，同时支持输出和存储扫描结果。
- 可扩展架构，支持后续功能增强和模块集成。

## 功能列表

### 1. Ping 检测工具

**功能描述**：检测目标主机是否在线，返回响应时间和丢包率。

**输入**：目标 IP 或域名。

**输出**：
- 响应时间（ms）
- 丢包率
- 连接成功/失败状态

**示例命令**：
```sh
network-toolbox ping --host=google.com
```

### 2. 端口扫描工具

**功能描述**：扫描目标主机的端口状态，返回开放/关闭端口。

**输入**：
- 目标 IP 地址或域名。
- 扫描端口范围（如 1-65535）。

**输出**：
- 开放端口列表
- 关闭端口列表
- 扫描耗时

**示例命令**：
```sh
network-toolbox portscan --host=192.168.1.1 --ports=1-1000
```

**附加功能**：
- 批量扫描多个目标。
- 扫描结果导出至文件（JSON、CSV）。

### 3. HTTP 请求工具

**功能描述**：模拟 HTTP 请求，支持 GET/POST 方法，自定义 Headers 和 Body。

**输入**：
- 请求方法（GET/POST）。
- 目标 URL。
- 请求头（Headers）和请求体（Body）。

**输出**：
- HTTP 响应状态码
- 响应时间
- 响应头和响应体

**示例命令**：
```sh
network-toolbox http --url="https://example.com" --method=POST --headers="Content-Type:application/json" --body='{"key":"value"}'
```

### 4. DNS 查询工具

**功能描述**：查询域名对应的 IP 地址、CNAME、MX 等记录。

**输入**：目标域名。

**输出**：
- 域名解析结果（IP、CNAME、MX 记录）。

**示例命令**：
```sh
network-toolbox dns --domain=example.com
```

### 5. 网络速度测试工具

**功能描述**：通过下载/上传测试文件，计算网络带宽速度。

**输入**：无（或指定测试服务器）。

**输出**：
- 下载速度（Mbps）
- 上传速度（Mbps）
- 网络延迟

**示例命令**：
```sh
network-toolbox speedtest
```

### 6. IP 地址查询工具

**功能描述**：查询本地 IP 地址，或通过第三方服务查询外部 IP 及地理位置信息。

**输入**：无。

**输出**：
- 本机内网 IP 地址
- 公网 IP 地址及地理位置信息

**示例命令**：
```sh
network-toolbox ip
```

### 7. Traceroute 工具

**功能描述**：显示数据包到目标主机的网络路径和延迟情况。

**输入**：目标 IP 地址或域名。

**输出**：
- 网络路径中的所有节点
- 每个节点的延迟时间

**示例命令**：
```sh
network-toolbox traceroute --host=google.com
```

### 8. AI 分析功能

**功能描述**：对工具箱生成的结果进行智能化分析，并提供诊断建议。

**输入**：
- 扫描结果（端口扫描、Ping、HTTP 请求等）。

**输出**：
- 异常检测（如开放端口风险、网络延迟问题）。
- 网络性能评估。
- 优化建议（安全加固、配置优化等）。
- 诊断总结报告。

**示例命令**：
```sh
network-toolbox portscan --host=192.168.1.1 --ai-analysis
```

**AI 数据来源**：
- 接入外部 AI 引擎（如 OpenAI API、Hugging Face）。
- 预训练规则模型（简单的本地逻辑）。

## 输出与存储需求

### CLI 直接输出

所有工具执行结果均可在 CLI 中直接展示。

### 结果导出

支持导出扫描和分析结果至 JSON、CSV、Markdown 文件。

### 日志记录

所有操作支持记录日志，方便后续排查问题。

### Web 支持（可选）

提供简单的 Web 界面展示和管理扫描结果。

## 非功能性需求

### 性能

支持高并发执行端口扫描和 HTTP 请求测试。

### 可扩展性

设计模块化架构，便于后续功能扩展。

### 易用性

提供清晰的命令行交互与帮助文档。

### 可靠性

处理网络异常情况（如超时、拒绝连接）并给出反馈。

## 技术栈

- **编程语言**：Go
- **框架与库**：
    - CLI 交互：Cobra
    - 网络功能：Go 标准库（net、http）
    - 并发处理：Goroutines 和 Channels
    - 数据存储：JSON / SQLite
    - 日志记录：log 包
    - AI 分析：外部 API（如 OpenAI API）或本地模型
- **工具**：
    - 数据可视化：go-chart 库
    - 测试：Go 原生测试框架

## 开发阶段与里程碑

### 阶段一：基础功能实现

- 完成 Ping、端口扫描、HTTP 请求、DNS 查询等基础功能。
- 实现 CLI 交互与日志记录。

### 阶段二：高级功能实现

- 添加网络速度测试、Traceroute 和 IP 查询功能。
- 支持结果导出（JSON/CSV/Markdown）。

### 阶段三：AI 分析集成

- 接入 AI 引擎，对扫描结果进行智能化分析。
- 输出优化建议和诊断总结。

### 阶段四：功能优化与扩展

- 添加批量操作、Web 界面支持。
- 性能优化与并发处理增强。

## 附录：示例用户操作

### 示例命令

**Ping 检测**：
```sh
network-toolbox ping --host=google.com
```

**端口扫描（带AI分析）**：
```sh
network-toolbox portscan --host=192.168.1.1 --ports=1-65535 --ai-analysis
```

**HTTP 请求测试**：
```sh
network-toolbox http --url="https://example.com" --method=POST --headers="Content-Type:application/json" --body='{"key":"value"}'
```

**输出结果到文件**：
```sh
network-toolbox portscan --host=example.com --ports=80,443 --output=result.json
```
