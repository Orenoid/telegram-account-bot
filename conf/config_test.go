package conf

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestGetConfigFromEnv(t *testing.T) {
	envs := map[string]string{
		"MYSQL_DSN":     strconv.Itoa(rand.Int()),
		"TELEBOT_TOKEN": strconv.Itoa(rand.Int()),
	}
	for key, value := range envs {
		_ = os.Setenv(key, value)
	}
	conf, err := GetConfigFromEnv()
	assert.NoError(t, err)
	assert.Equal(t, envs["MYSQL_DSN"], conf.MysqlDSN)
	assert.Equal(t, envs["TELEBOT_TOKEN"], conf.TelebotToken)
}
