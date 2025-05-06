package gormex_postgres_test

import (
	"testing"
	_ "vngom/gormex/dbconfig"
	"vngom/gormex/dbconfig/dbconfig_postgres"
	_ "vngom/gormex/dbconfig/dbconfig_postgres"

	"vngom/gormex/gormex_postgres"

	assert "github.com/stretchr/testify/assert"
)

var yamlFile = "E:/Docker/go/quicky-go/be/gormex/postgres.yaml"

func Test_New(t *testing.T) {
	cfg := dbconfig_postgres.PostgresDbConfig{}
	cfg.LoadFromYamlFile(yamlFile)

	g := gormex_postgres.NewGormEx(&cfg)
	g1 := gormex_postgres.NewGormEx(&cfg)
	assert.Equal(t, g, g1)

	err := g.GetDbConfig().PingDb()
	assert.Equal(t, nil, err)
}
