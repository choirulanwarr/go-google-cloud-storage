package router

func initGlobalRoutes(config *Config) {
	globalApi := config.Server.Group("/api/v1")

	// File
	globalApiFile := globalApi.Group("/")
	globalApiFile.POST("/upload", config.FileHandler.UploadFile)
	globalApiFile.GET("/download", config.FileHandler.DownloadFile)
	globalApiFile.GET("/list/:folder", config.FileHandler.GetFile)
}
