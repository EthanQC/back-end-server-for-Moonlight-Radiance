package http

import (
	"time"

	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/api/websocket"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/auth"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/card"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/room"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/internal/user"
	"github.com/EthanQC/back-end-server-for-Moonlight-Radiance/pkg/common"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// 推荐使用自定义配置，而不是 cors.Default()
	// 确保 "Authorization" 被允许
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://127.0.0.1:5173", "http://localhost:5173"},
		// 如果你想在生产环境允许别的域名，也可以加进来
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.SetTrustedProxies([]string{"127.0.0.1"})

	// WebSocket handler
	wsHandler := websocket.NewHandler()
	hub := wsHandler.Hub()
	go hub.Run()

	// 用户相关API
	userGroup := router.Group("/api/user")
	{
		userGroup.POST("/register", user.RegisterHandler)
		userGroup.POST("/login", user.LoginHandler)
	}

	// 需要认证的API
	authorized := router.Group("/api")
	authorized.Use(auth.AuthMiddleware())
	{
		// WebSocket连接
		authorized.GET("/ws", wsHandler.HandleConnection)

		// 房间相关
		roomGroup := authorized.Group("/rooms")
		{
			roomHandler := room.NewRoomHandler(room.NewRoomService(common.DB))
			roomGroup.POST("/create", roomHandler.CreateRoomHandler)
			roomGroup.POST("/join", roomHandler.JoinRoomHandler)
			roomGroup.GET("/state", roomHandler.GetRoomStateHandler)
		}

		// 卡牌相关
		cardGroup := authorized.Group("/cards")
		{
			cardHandler := card.NewCardHandler(card.NewCardService(common.DB))
			cardGroup.POST("/init", cardHandler.InitDeckHandler)
			cardGroup.POST("/state", cardHandler.GetCardStateHandler)
			cardGroup.POST("/draw", cardHandler.DrawCardsHandler)
			cardGroup.POST("/play", cardHandler.PlayCardHandler)
			cardGroup.POST("/endTurn", cardHandler.EndTurnHandler)
		}

		// // 对战地图相关
		// mapGroup := authorized.Group("/battlemap")
		// {
		// 	mapHandler := battlemap.NewBattleMapHandler()
		// 	mapGroup.POST("/create", mapHandler.CreateMapHandler)
		// 	mapGroup.POST("/placeCard", mapHandler.PlaceCardHandler)
		// 	mapGroup.GET("/state", mapHandler.GetMapStateHandler)
		// }

		// // 竞速地图相关
		// raceGroup := authorized.Group("/racemap")
		// {
		// 	raceHandler := racemap.NewRaceMapHandler()
		// 	raceGroup.POST("/create", raceHandler.CreateMapHandler)
		// 	raceGroup.POST("/move", raceHandler.MoveForwardHandler)
		// 	raceGroup.GET("/state", raceHandler.GetPositionHandler)
		// }
	}

	return router
}
