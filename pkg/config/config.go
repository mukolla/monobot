package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken  string //must BindEnv
	MonoBankToken  string //must BindEnv
	TelegramBotUrl string `mapstructure:"bot_url"`
	DbPath         string `mapstructure:"db_file"`
	Message        Message
}

type Message struct {
	Errors   Errors
	Response Response
}

type Errors struct {
	Default                string `mapstructure:"default"`
	UnknownError           string `mapstructure:"unknownError"`
	AuthTokenNotFound      string `mapstructure:"authTokenNotFound"`
	Unauthorized           string `mapstructure:"unauthorized"`
	GetBalance             string `mapstructure:"getBalance"`
	GetTransactionList     string `mapstructure:"getTransactionList"`
	AccountNotFoundByToken string `mapstructure:"accountNotFoundByToken"`
}

type Response struct {
	Start             string `mapstructure:"start"`
	SavedSuccessfully string `mapstructure:"savedSuccessfully"`
	UnknownCommand    string `mapstructure:"unknown_command"`
	AlreadyUsed       string `mapstructure:"alreadyUsed"`
	ChoiceAccountUsed string `mapstructure:"choiceAccountUsed"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("main")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		//panic(fmt.Errorf("fatal error config file: %w", err))
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("message.response", &cfg.Message.Response); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("message.errors", &cfg.Message.Errors); err != nil {
		return nil, err
	}

	if err := parserEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parserEnv(cfg *Config) error {

	if err := viper.BindEnv("token"); err != nil {
		return err
	}

	if err := viper.BindEnv("mono_bank_token"); err != nil {
		return err
	}

	cfg.TelegramToken = viper.GetString("token")
	cfg.MonoBankToken = viper.GetString("mono_bank_token")
	return nil
}
