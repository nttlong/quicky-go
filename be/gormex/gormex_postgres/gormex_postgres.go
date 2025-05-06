package gormex_postgres

import (
	"sync"
	"vngom/gormex"
	"vngom/gormex/dbconfig"
)

type GormExPostgres struct {
	Config dbconfig.IDbConfig
}

func (g *GormExPostgres) GetDbConfig() dbconfig.IDbConfig {
	return g.Config
}

var ins *GormExPostgres
var once sync.Once

func NewGormEx(cfg dbconfig.IDbConfig) gormex.IGormEx {
	once.Do(func() {
		ins = &GormExPostgres{
			Config: cfg,
		}
	})
	return ins

}
