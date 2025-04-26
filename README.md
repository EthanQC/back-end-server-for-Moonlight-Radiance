# 月华

## 项目概述
《月华》是一款以**月相变化**为核心机制，结合了牌组构建、竞速、跑分、肉鸽（Rogue-like）等多种机制以及**中国传统元素**的单机游戏，[详情]()（steam 页面建设中）

目前《月华》处于早期开发和概念阶段，预计在五月初会发布第一版初始demo

**demo 不代表最终作品质量**

本仓库是《月华》的后端开发仓库，旨在记录我的独立开发过程

目录结构：

    ├── api/                    
    │   ├── http/              # HTTP API定义
    │   │    └── routes.go     # 路由定义
    │   └── proto/             # gRPC定义
    │   │    └── game.proto    # 游戏服务定义
    │   └── websocket/         # WebSocket API 定义
    │   │    └── status.go     # 实时状态同步
    ├── cmd/                   
    │   └── server/           
    │        └── main.go       # 主程序入口
    ├── configs/               
    │   ├── config.go          # 配置文件
    ├── docs/                  
    │   ├── architecture.md    # 服务器架构文档
    │   ├── features.md        # 服务器功能文档
    │   └── requirements.md    # 需求文档
    ├── internal/              
    │   ├── auth/              # 认证模块
    │   │   ├── middleware.go  # 认证中间件
    │   │   └── jwt.go         # 认证服务        
    │   ├── card/          # 卡牌模块
    │   │   ├── model.go 
    │   │   ├── handler.go  
    │   │   └── service.go
    │   ├── battlemap/     # 对战地图模块
    │   │   ├── model.go
    │   │   ├── handler.go 
    │   │   └── service.go 
    │   ├── racemap/       # 竞速地图模块
    │   │   ├── model.go
    │   │   ├── handler.go
    │   │   └── service.go 
    │   ├── room/          # 房间模块
    │   │   ├── model.go
    │   │   ├── hanlder.go
    │   │   └── service.go
    │   └── user/              # 用户模块
    │       ├── model.go    
    │       ├── hanlder.go     
    │       └── service.go   
    └── pkg/                  
        └── common/           
            ├── db.go          # 数据库连接
            ├── redis.go       # Redis连接
            └── logger.go      # 简易日志模块

## 快速开始
#### 环境要求

* Go 1.22+
* MySQL 8.0+
* Redis 6.0+
* Ubuntu 24.04

#### 安装步骤

* 配置环境变量：
    * `export DB_DSN="root:your_password@tcp(localhost:3306)/moonlight?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"`
    * `export REDIS_ADDR="localhost:6379"`
    * `export JWT_SECRET="your-secret-key"`
* 启动 MySQL、redis：
    * `sudo systemctl start mysql`
    * `sudo systemctl start redis`
* [更改 MySQL 权限](https://github.com/EthanQC/my-learning-record/blob/main/database/FAQ.md)
* 创建数据库表：
    * 登录：`mysql -u root -p`
    * 运行 SQL 脚本：
        * `source migrations/001_init.sql`
        * `source migrations/002_add_card_tables.sql`
        * `source migrations/003_insert_cards.sql`
        * `source migrations/004_add_room_tables.sql`
* 启动服务器：`go run cmd/server/main.go`

#### 测试账号
* Ethan Wu
    * 123456

## 文档索引
* [架构设计](https://github.com/EthanQC/back-end-server-for-Moonlight-Radiance/blob/main/docs/architecture.md)
* [功能说明](https://github.com/EthanQC/back-end-server-for-Moonlight-Radiance/blob/main/docs/features.md)
* [已实现模块详细说明](https://github.com/EthanQC/back-end-server-for-Moonlight-Radiance/blob/main/docs/implemented.md)
* [服务器游戏玩法和设计需求](https://github.com/EthanQC/back-end-server-for-Moonlight-Radiance/blob/main/docs/requirements.md)

## 开发进度
| 日期 | 当前进度 |
| ------- | ------- |
| 2025.2.15 | 创建了前后端仓库，完成了基本的框架架构设计 |
| 2025.2.16 | 目前正在完善基本玩法设计，预计1-2天即可完成，随后会马上转向pve玩法的开发；开始编写需求文档 |
| 2025.2.17 | 搭建了项目基本框架，开始写后端服务器的代码，需要学习大量新知识 |
| 2025.2.18 | 准备开始使用虚幻五引擎开发，玩法设计基本完成，只差对战地图设计 |
| 2025.2.19-2.22 | 在出去玩 |
| 2025.2.23 | 初步完成了后端的用户和认证模块的开发，感觉可能3月初做不出来demo了 |
| 2025.2.24 | 画了 pve 第六章的对战地图 |
| 2025.2.25-2.28 | 在上班，详情见碎碎念 |
| 2025.3.1-3.2 | 上班太累了，在出去玩 |
| 2025.3.3 | 把基本的代码框架弄完了，还没细看，不知道该怎么规划开发进度 |
| 2025.3.4-3.7 | 在上班和上学，累 |
| 2025.3.8 | 在改简历 |
| 2025.3.9 | 整理学习计划，准备大学特学一下，希望4月初能出 demo |
| 2025.3.10 | 上班，累，玩，思考目前的项目架构 |
| 2025.3.11 | 看虚幻五官方技术文档，画完了 pve 第四章的对战地图并上传 |
| 2025.3.12 | 上课，写作业，累 |
| 2025.3.13 | 看虚幻五官方文档，感觉到了虚幻五不小的学习成本 |
| 2025.3.14 | 上班，看 MySQL 底层八股，大幅度调整游戏的后端服务器仓库结构 |
| 2025.3.15 | 彻底梳理清楚后端服务器仓库的结构，添加了简易的前端页面，把服务器跑起来并测试了目前已经完成的用户注册、登录和认证模块，它们都没问题，后续可以再根据需求作修改 |
| 2025.3.16 | 完成了卡牌模块的第一版开发，包括分层结构和单元测试，完成后发现需要设计一下游戏联机版的玩法，之前只确定了单机版的，联机版还只是个概念，准备花时间好好构思一下玩法设计；开始规划单机版和联机版的开发日程 |
| 2025.3.17-3.20 | 上班上课，累 |
| 2025.3.21 | 初始化前端项目，准备开始写；把后端服务器的基本模板全写完了，但还要自己一个一个看一个一个调整修改才行 |
| 2025.3.22-3.23 | 增加SQL脚本，搭建前端代码框架 |
| 2025.3.24 | 调整后端路由，添加前端必需文件 |
| 2025.3.25-3.26 | 修改前端 UI 界面，实现了游戏 Home 页的 MVP 版本；修改了后端路由配置，实现了前后端用户以及卡牌模块的 http 通信；完善项目 readme 结构，修改文档模块，补充对已实现模块的详细解释，便于学习 |
| 2025.3.27 | 有点忘了干了啥了，好像是在看已经实现模块的代码然后看了下部署上云的东西 |
| 2025.3.28 | 看已经实现模块的代码，添加 SQL 脚本 |
| 2025.3.29 | 改动了 common 模块，实现了房间模块，开始修改卡牌模块，卡牌模块要改的地方好多啊啊啊，前端也好多要改的啊啊啊啊啊 |
| 2025.3.30-3.31 | 学习虚幻五引擎 |
| 2025.4.1-4.2 | 学习虚幻五引擎 |
| 2025.4.3 | 上课、做实验、累，看剧 |
| 2025.4.4-4.5 | 我受够了这样半吊子的代码，我要先看看别人是怎么做的，把这些技术栈搞懂之后再继续我的开发，我不能再这样忍受下去了 |
| 2025.4.6 | 阶段性反思总结，规划接下来的方向，决定先把重心放回后端开发与基础、底层的学习 |
| 2025.4.7 | 上班，看八股，改学习记录仓库 |
| 2025.4.8-4.17 | 上班，上学，看八股，刷算法题，更新学习记录仓库 |
| 2025.4.18 | 冷静了十天，重新开始继续开发项目，只是因为不能将所有空闲时间都花在项目上，还要看八股刷题和学其他的什么的，所以 demo 发布时间可能又又要延期了 👉👈 |
| 2025.4.19-4.25 | 准备面试 |
| 2025.4.26 | 重新评估后意识到这个项目目前在实际面试后并不讨巧，且开发时阻力较大，于是决定暂时搁置目前的联机版开发，游戏开发这块准备纯粹作为自己的兴趣来了，所以后续这个仓库可能要过一段时间才会继续维护了 |