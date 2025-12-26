package main

import (
	"fmt"
	"os"

	"github.com/Android12349/food_recomendation/profile_managment_service/config"
	"github.com/Android12349/food_recomendation/profile_managment_service/internal/bootstrap"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	profileStorage := bootstrap.InitPGStorage(cfg)
	menuGenerationProducer := bootstrap.InitMenuGenerationProducer(cfg)
	profileService := bootstrap.InitProfileService(profileStorage, menuGenerationProducer, cfg)
	profileApi := bootstrap.InitProfileManagementAPI(profileService)

	bootstrap.AppRun(*profileApi, cfg)
}
