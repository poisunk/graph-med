# Graph-Med 医疗知识图谱应用

## 项目介绍

Graph-Med 是一个基于知识图谱的医疗信息查询和智能问答系统，旨在提供专业、准确的医疗知识服务。系统整合了疾病、症状、治疗方法等医疗知识，通过图数据库构建知识网络，并结合大语言模型提供智能化的医疗咨询服务。

## 后端技术栈

- **编程语言**：Go 1.23
- **微服务框架**：go-zero
- **数据库**：
  - MySQL：存储用户数据等关系型数据
  - MongoDB：存储聊天记录、会话信息等文档型数据
- **缓存**：Redis
- **权限管理**：Casbin
- **依赖注入**：Wire
- **MCP服务**：mcp-go
- **认证授权**：JWT
- **异步任务队列**：Asynq

## 大语言模型集成

- 火山引擎 ARK Runtime API 集成
- MCP (Model Control Protocol) 工具集成，用于模型控制和功能扩展

## 运行截图

![截图1](./img/img.png)
![截图2](./img/img_1.png)

## 核心功能

### 1. 用户认证与授权

- 用户注册、登录功能
- 基于JWT的身份验证
- 基于Casbin的RBAC权限控制
- 验证码与邮箱验证

### 2. 医疗知识图谱

- 疾病信息查询
- 多标签分类浏览
- 知识图谱可视化
- 节点关系查询

### 3. 智能问答系统

- 基于大语言模型的医疗问答
- 会话管理与历史记录
- 用户反馈收集
- 实时流式响应
- 知识图谱增强的医疗问答

### 4. MCP服务功能

- 知识图谱工具调用
- 医疗专业提示词管理
- 疾病子图查询
- 大语言模型能力增强

## 项目结构

```
├── app                 # 应用服务
│   ├── captcha         # 验证码服务
│   │   ├── api        # HTTP API
│   │   └── rpc        # RPC 服务
│   ├── chat           # 聊天服务
│   │   ├── api        # HTTP API
│   │   ├── model      # 数据模型
│   │   └── rpc        # RPC 服务
│   ├── mqueue         # 消息队列服务
│   │   ├── job        # 任务处理
│   │   └── scheduler  # 任务调度
│   └── usercenter     # 用户中心服务
│       ├── api        # HTTP API
│       ├── model      # 数据模型
│       └── rpc        # RPC 服务
├── data               # 数据文件
├── deploy             # 部署相关
│   ├── goctl          # 代码生成工具
│   ├── nginx          # Nginx配置
│   ├── script         # 部署脚本
│   └── sql            # SQL脚本
├── img                # 图片资源
├── pkg                # 公共包
│   ├── ctxdata        # 上下文数据
│   ├── interceptor    # 拦截器
│   ├── result         # 响应结果
│   ├── tool           # 工具函数
│   └── xerr           # 错误处理
├── docker-compose.yml # Docker编排配置
└── Makefile           # 构建脚本
```
