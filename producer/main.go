package main

import (
	accountcontrollers "producer/controllers/account"
	accountservice "producer/services/account"
	eventproducerservice "producer/services/producer"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func main() {
	producer, err := sarama.NewSyncProducer(viper.GetStringSlice("kafka.servers"), nil)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	eventProducer := eventproducerservice.NewEventProducer(producer)
	accountService := accountservice.NewAccountService(eventProducer)
	accountController := accountcontrollers.NewAccountController(accountService)

	app := fiber.New()

	app.Post("/openAccount", accountController.OpenAccount)
	app.Post("/depositFund", accountController.DepositFund)
	app.Post("/withdrawFund", accountController.WithdrawFund)
	app.Post("/closeAccount", accountController.CloseAccount)

	app.Listen(":8000")
}
