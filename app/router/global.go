package router

func initGlobalRoutes(config *Config) {
	globalApi := config.Server.Group("/api/v1")

	// File
	globalApiFile := globalApi.Group("/")
	globalApiFile.POST("/bucket/create", config.FileHandler.CreateBucket)
	globalApiFile.POST("/upload", config.FileHandler.UploadFile)
	globalApiFile.GET("/download", config.FileHandler.DownloadFile)
	globalApiFile.GET("/list", config.FileHandler.GetAllFile)
	globalApiFile.GET("/list/:folder", config.FileHandler.GetSpecificFile)
	globalApiFile.DELETE("/delete", config.FileHandler.DeleteFile)
}
