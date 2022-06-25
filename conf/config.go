package conf

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reflect"
)

type config struct {
	MysqlDSN     string `mapstructure:"MYSQL_DSN"`
	TelebotToken string `mapstructure:"TELEBOT_TOKEN"`
}

func (c config) tags() []string {
	var tags []string
	t := reflect.TypeOf(c)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tags = append(tags, field.Tag.Get("mapstructure"))
	}
	return tags
}

func GetConfigFromEnv() (config, error) {
	var c config
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Info("The .env file has not been found in the current directory")
			for _, tag := range c.tags() {
				err = viper.BindEnv(tag)
				if err != nil {
					return config{}, errors.Wrapf(err, "failed to bind env: %s", tag)
				}
			}
		} else {
			return config{}, err
		}
	}

	if err := viper.Unmarshal(&c); err != nil {
		return config{}, errors.Wrapf(err, "failed to unmarshal config to struct")
	}

	return c, nil
}
