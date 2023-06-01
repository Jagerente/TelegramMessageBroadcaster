package main

import (
	"DC_NewsSender/internal/db"
	"DC_NewsSender/internal/db/models"
	"DC_NewsSender/internal/db/repositories"
	"DC_NewsSender/internal/telegram"
	"DC_NewsSender/pkg/configuration"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Config struct {
	TgToken        string `mapstructure:"TG_TOKEN"`
	TgMasterId     int64  `mapstructure:"TG_MASTER_ID"`
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         uint16 `mapstructure:"POSTGRES_PORT"`
	Debug          bool   `mapstructure:"DEBUG"`
}

var (
	bot    *telegram.Core
	orm    *gorm.DB
	logger *zap.Logger
)

func configureMaster(provider *repositories.Provider, masterId int64) error {
	repo := provider.CreateAdminsRepo()
	_, err := repo.FindById(masterId)
	if err != nil {
		_, err := repo.Add(&models.Admin{Id: masterId, Name: "Master", IsMaster: true})
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	cfg, err := configuration.New[Config]()
	if err != nil {
		panic(err)
	}

	env := cfg.ENV
	logger = configuration.GetLogger(env.Debug)

	var dbConfig *db.PostgresDatabaseConfiguration = &db.PostgresDatabaseConfiguration{
		Host:         env.DBHost,
		UserName:     env.DBUserName,
		UserPassword: env.DBUserPassword,
		DatabaseName: env.DBName,
		Port:         env.DBPort,
	}
	logger.Sugar().Debug(dbConfig)
	orm, err = db.InitializePostgresDatabase(dbConfig)
	if err != nil {
		logger.Panic(err.Error())
	}

	provider := repositories.CreateProvider(orm)

	if err := configureMaster(provider, env.TgMasterId); err != nil {
		logger.Panic(err.Error())
	}

	bot, err = telegram.CreateBotCore(&telegram.BotConfig{
		Token:  env.TgToken,
		Db:     provider,
		Logger: logger,
		Debug:  env.Debug})
	if err != nil {
		logger.Panic(err.Error())
	}
}

func main() {
	defer db.CleanupConnection(orm)
	defer logger.Sync()

	bot.Run()
}
