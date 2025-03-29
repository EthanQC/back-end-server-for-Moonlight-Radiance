## 服务器架构

本项目使用微服务架构，结合GO语言与C++，在Ubuntu环境下实现高并发支持和状态实时更新，集成注册中心，自动注册微服务元数据

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