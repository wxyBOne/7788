package main

import (
	"log"
	"net/http"

	"seven-ai-backend/internal/config"
	"seven-ai-backend/internal/database"
	"seven-ai-backend/internal/handlers"
	"seven-ai-backend/internal/middleware"
	"seven-ai-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 初始化服务
	userService := services.NewUserService(db)
	characterService := services.NewCharacterService(db)
	companionService := services.NewCompanionService(db)
	conversationService := services.NewConversationService(db)
	aiService := services.NewAIService(
		cfg.AIAPIKey,
		cfg.AIBaseURL,
		cfg.ASRAPIKey,
		cfg.TTSAPIKey,
		cfg.VisionAPIKey,
	)

	// 初始化处理器
	userHandler := handlers.NewUserHandler(userService)
	characterHandler := handlers.NewCharacterHandler(characterService)
	companionHandler := handlers.NewCompanionHandler(companionService, aiService)
	conversationHandler := handlers.NewConversationHandler(conversationService, aiService)

	// 设置路由
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 路由组
	api := r.Group("/api/v1")
	{
		// 用户相关
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.GET("/profile", middleware.AuthRequired(), userHandler.GetProfile)
			users.PUT("/profile", middleware.AuthRequired(), userHandler.UpdateProfile)
		}

		// 预设角色相关
		characters := api.Group("/characters")
		{
			characters.GET("", characterHandler.GetAllCharacters)
			characters.GET("/search", characterHandler.SearchCharacters)
			characters.GET("/:id", characterHandler.GetCharacter)
		}

		// AI伙伴相关
		companions := api.Group("/companions")
		companions.Use(middleware.AuthRequired())
		{
			companions.POST("", companionHandler.CreateCompanion)
			companions.GET("", companionHandler.GetUserCompanions)
			companions.GET("/:id", companionHandler.GetCompanion)
			companions.PUT("/:id", companionHandler.UpdateCompanion)
			companions.DELETE("/:id", companionHandler.DeleteCompanion)
			companions.GET("/:id/growth", companionHandler.GetGrowthStatus)
			companions.GET("/:id/diary", companionHandler.GetDiary)
		}

		// 对话相关
		conversations := api.Group("/conversations")
		conversations.Use(middleware.AuthRequired())
		{
			conversations.POST("/chat", conversationHandler.Chat)
			conversations.POST("/voice-chat", conversationHandler.VoiceChat)
			conversations.POST("/image-chat", conversationHandler.ImageChat)
			conversations.GET("/history", conversationHandler.GetHistory)
			conversations.GET("/sessions/:sessionId", conversationHandler.GetSessionHistory)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 启动服务器
	log.Printf("服务器启动在端口 %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
