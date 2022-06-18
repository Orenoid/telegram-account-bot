package conf

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetConfigFromEnv(t *testing.T) {
	mysqlDSN := "dsn string"
	_ = os.Setenv("MYSQL_DSN", mysqlDSN)
	conf, err := GetConfigFromEnv()
	assert.NoError(t, err)
	assert.Equal(t, mysqlDSN, conf.MysqlDSN)
}
