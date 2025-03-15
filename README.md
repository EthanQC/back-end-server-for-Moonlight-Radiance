# 月华
《月华》是一款以**月相变化**为核心机制，结合了牌组构建、竞速、跑分、肉鸽（Rogue-like）等多种机制以及**中国传统元素**的单机游戏，[详情]()（网站建设中）

目前《月华》处于早期开发和概念阶段，预计在3月初会发布第一版初始demo

本仓库是《月华》的后端开发仓库，旨在记录我的独立开发过程

## 如何运行
* 运行环境：
    * 
* 创建数据库表：
    * 
* 更改 MySQL 权限：
    * 
* 配置环境变量：
    * `export DB_DSN="root:your_password@tcp(localhost:3306)/moonlight?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"`
    * `export REDIS_ADDR="localhost:6379"`
    * `export JWT_SECRET="your-secret-key"`
* 启动 MySQL、redis：
    * `sudo systemctl start mysql`
    * `sudo systemctl start redis`
* 启动服务器：`go run cmd/server/main.go`

## 测试账号
* Ethan Wu
    * 123456
* Ethan
    * 111

## 服务器架构

本项目使用微服务架构，结合GO语言与C++，在Ubuntu环境下实现高并发支持和状态实时更新，集成注册中心，自动注册微服务元数据

    ├── api/                    
    │   ├── http/              # HTTP API定义
    │   │   ├── routes.go      # 路由定义
    │   │   ├── handlers/      # API处理器
    │   └── proto/             # gRPC定义
    │   │    └── game.proto    # 游戏服务定义
    │   └── ws/                # WebSocket API 定义
    ├── build/                 # 构建相关文件
    │   ├── Dockerfile        
    │   └── docker-compose.yml
    ├── cmd/                   
    │   └── server/           
    │       └── main.go       # 主程序入口
    ├── configs/               
    │   ├── app.yaml          # 应用配置
    │   └── db.yaml           # 数据库配置
    ├── docs/                  # 文档
    │   ├── api/             # API文档
    │   └── design/          # 设计文档
    ├── internal/              
    │   ├── auth/             # 认证模块
    │   │   ├── handler.go    # HTTP处理
    │   │   ├── middleware.go # 认证中间件
    │   │   └── service.go    # 认证服务
    │   ├── game/            
    │   │   ├── card/         # 卡牌系统
        │   │   ├── handler.go 
    │   │   │   ├── model.go  
    │   │   │   └── service.go
    │   │   ├── engine/       # 游戏引擎
    │   │   │   ├── battle.go
    │   │   │   └── score.go 
    │   │   ├── map/       # 游戏引擎
    │   │   │   ├── handler.go
    │   │   │   └── model.go 
    │   │   └── room/         # 房间管理
    │   │       ├── model.go
    │   │       └── service.go
    │   └── user/             # 用户模块
    │       ├── handler.go    
    │       ├── model.go     
    │       └── service.go   
    └── pkg/                  
        ├── common/           
        │   ├── db.go        # 数据库连接
        │   ├── redis.go     # Redis连接
        │   └── logger.go    # 日志工具
        └── middleware/      
            └── auth.go      # 通用认证中间件

* **Go 服务**：主服务器，用于处理用户请求和管理业务逻辑
    * 主程序：负责启动服务和初始化
        * 使用Go的标准库来搭建Web服务（net/http，gin，echo等框架可以选择）
        * 提供RESTful API来与前端和C++服务交互
        * 使用goroutine和channel来处理并发任务
    * 数据库管理：使用gorm库来与MySQL进行交互
    * 缓存管理：使用go-redis库来实现与Redis的连接
    * 高性能服务：如果需要计算密集型任务，Go通过调用C++的接口来将计算任务交给C++模块处理
    * API网关：负责对外暴露接口，并与后端各个服务进行路由转发
* **C++服务**：计算密集型的服务，通过gRPC或HTTP与Go服务通信
    * 用C++编写一些计算密集型的任务（如数据处理、图像处理等）
    * C++与Go通过gRPC进行通信：
        * 定义Protobuf接口，Go与C++服务分别实现该接口
        * 通过gRPC请求来进行跨语言的数据传递和函数调用
* API网关：负责前后端之间的请求路由，可以使用Nginx或者使用Go自定义开发
    * 通过RESTful API（或者gRPC）来设计后端的各个接口。
    * 例如：
        * 用户登录（POST /login）
        * 获取商品信息（GET /products）
        * 创建订单（POST /orders）
        * 获取用户订单（GET /orders/{userId}）
* 数据库：
    * MySQL：关系型数据存储
        * 设计数据表来存储用户信息、商品信息、订单信息等。
        * 需要用到gorm（Go的ORM框架）来简化数据库操作。
    * Redis：缓存
        * 用作缓存，提升数据访问的效率。例如可以缓存热门商品、用户信息等。
    * ElasticSearch：日志和数据索引检索
        * 用于商品的全文搜索以及日志的存储和查询。
* 容器化：Docker + Kubernetes（用于扩展和容器化部署）
    * Docker化：
        * 为每个服务（Go服务、C++服务、MySQL、Redis等）创建一个Docker镜像。
        * 使用docker-compose来管理多个容器的运行。
    * Kubernetes：
        * 部署在Kubernetes上，确保服务的高可用和可伸缩性。
        * 使用Kubernetes的Deployment和Service来管理不同服务的部署和负载均衡。
    * CI/CD：
        * 配置GitHub Actions或GitLab CI来实现自动化部署。
        * 在每次代码提交时自动执行测试、构建镜像并推送到Docker Registry。
* 消息队列：RabbitMQ/Kafka（用于服务之间的异步通信）
* 部署：完整服务迁移上云部署，可选常规服务器部署、云托管、Fass等方式，这里不做任何限制
* 测试：有比较良好的编码规范严格按照技术编码规范进行编码，与此同时，针对业务类需求编写了相对完善的单元测试用例，选取合适的单元测试框架
* 编码：侧重服务端实现，会提前定义好各个功能对应的接口（接口定义推荐使用Protobuf，但不做限制），按说明实现接口即可在客户端中看到运行效果
    * 服务端最基本的结构只需要服务端程序和数据库即可，服务端程序连接数据库，响应客户端请求完成对应功能。同时需要根据功能，设计合理的数据模型，并创建对应的数据表，其中日志文件等可以保存到本地

Web框架：Gin


微服务框架：Go Micro


数据库管理：GORM + MySQL
* 考虑将数据库进行分区和分表处理，避免单个表的查询过于繁重。

缓存管理：go-redis


微服务通信：gRPC/HTTP + Protocol Buffers
* 在某些服务间的通信上选择HTTP，gRPC则用于对性能要求较高的部分

消息队列：NATS/Kafka


日志与搜索：ElasticSearch
* 根据不同的需求调整索引的刷新频率，避免过于频繁的更新操作影响性能

容器化：Docker + Kubernetes
* 考虑集成一个多集群管理工具，如Istio，来进行服务网格管理，使得微服务之间的流量管理、路由、负载均衡更加高效与灵活
* 在CI/CD中，可以结合蓝绿部署策略来确保系统升级的无缝过渡，减少服务中断时间

实时通信：WebSocket
* WebSocket适用于实时通信，但可能需要额外考虑连接数的管理与负载均衡。在用户量大的情况下，WebSocket连接可能会成为瓶颈，可以结合Redis的Pub/Sub或Kafka来分担消息传递的压力，确保系统的高可用性

测试框架：Go Testing + Testify



## 服务器功能
* 用户管理模块：
    * 注册与登录：用户注册、登录功能，支持玩家的身份认证，可以使用JWT进行令牌验证。（steam SDK？）
    * 玩家状态管理：管理玩家的月光值、卡牌状态、地图进度等。
    * 登录认证
        * 可以使用第三方现成的登录验证框架（CasBin、Satoken等），对请求进行身份验证
        * 可配置的认证白名单，对于某些不需要认证的接口或路径，允许直接访问
        * 可配置的黑名单，对于某些异常的用户，直接进行封禁处理（可选）
        * JWT认证：非常适合无状态的API设计，支持高并发的用户请求处理。可以考虑加上Token的刷新机制（Refresh Token），保证用户登录状态不会因为长时间未活动而被强制注销。
        * OAuth2.0与第三方认证：如果你的游戏支持第三方登录（如Steam等），OAuth2.0是一个理想的选择，可以配合JWT使用，实现轻量级的认证和授权。
    * 权限认证（高级）
        * 根据用户的角色和权限，对请求进行授权检查，确保只有具有相应权限的用户能够访问特定的服务或接口。
        * 支持正则表达模式的权限匹配
        * 支持动态更新用户权限信息，当用户权限发生变化时，权限校验能够实时生效。
* 游戏引擎模块：
    * 卡牌机制引擎：根据游戏灵感中的卡牌机制（如“同相”、“反相”、“月相周期”）计算并更新玩家的状态。
    * 地图与格子管理：管理玩家的竞速地图，计算玩家的进度、归属权、月光消耗等。
    * 事件机制：实现全局事件，如“月光风暴”、“干扰区”等，根据游戏进度动态触发，并通知前端。
    * 卡牌机制与事件系统：建议尽早开始对卡牌机制、地图管理、事件机制的业务逻辑建模，并使用单元测试来验证各个模块的正确性。可以使用类似状态机（State Machine）的设计模式来管理不同的游戏状态。
    * 并发处理：在高并发的环境下，可能需要使用类似Redis的分布式锁来同步玩家的游戏状态，避免并发修改引发数据不一致。
* PVP与PVE机制：
    * PVP机制：管理玩家之间的对抗与月光消耗，判断玩家是否有资格前进。
    * PVE机制：根据章节的主题与BOSS设计，控制剧情的推进，利用月相牌解锁剧情。
* 游戏存档与重启：
    * 保存游戏的进度和玩家状态，支持断线重连和游戏重启。
    * 存储玩家的卡牌、月光进度、地图状态等数据。
    * 高可用存档服务：存档功能需要非常高的可用性，建议使用分布式存储（如Ceph或MinIO）来存储玩家的游戏进度，支持断点续传和快速重启。
* 推送与通知：
    * 后端通过WebSocket推送实时游戏状态更新给前端。
    * Push通知与消息队列的结合：WebSocket可以作为主要的实时推送手段，但对于离线用户的推送通知，考虑结合消息队列（如NATS或RabbitMQ）来管理通知的发送，确保在不同的网络环境下也能及时送达。
* 日志记录与监控
  * 对服务的运行状态和请求处理过程进行详细的日志记录，方便故障排查和性能分析。
  * 提供实时监控功能，能够及时发现和解决系统中的问题。
* 容错机制
  * 该服务应具备一定的容错能力，当出现部分下游服务不可用或网络故障时，能够自动切换到备用服务或进行降级处理。
  * 保证下游在异常情况下，系统的整体可用性不会受太大影响，且核心服务可用。
  * 服务应该具有一定的流量兜底措施，在服务流量激增时，应该给予一定的限流措施。

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