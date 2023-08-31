package config

import "github.com/spf13/viper"

type Config struct {
	PortHttp           string `mapstructure:"HTTP_PORT"`
	PgUser             string `mapstructure:"POSTGRES_USER"`
	PgPassword         string `mapstructure:"POSTGRES_Password"`
	PgHost             string `mapstructure:"POSTGRES_HOST"`
	PgPort             string `mapstructure:"POSTGRES_PORT"`
	PgDB               string `mapstructure:"POSTGRES_DB"`
	PgUrl              string `mapstructure:"POSTGRES_URL"`
	GDriveJSONFilePath string `mapstructure:"GOOGLE_DRIVE_JSON_FILE_PATH"`
}

// LoadConfig Конструктор для создания Config, который содержит считанные из .env файла данные.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
