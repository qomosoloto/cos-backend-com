package app

import (
	"cos-backend-com/src/common/app"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/eth"
	"cos-backend-com/src/eth/processor"
	"cos-backend-com/src/libs/filters"

	"github.com/jmoiron/sqlx"

	t "github.com/wujiu2020/strip"
)

const (
	AppName = "eth"
)

func AppInit(tea *t.Strip, confPath string, files ...string) *appConfig {
	app := &appConfig{app.New(tea, AppName), eth.Env}
	app.ConfigLoad(app.Env, confPath, files...)
	app.ConfigRoutes()
	app.ConfigDB()
	app.ConfigFilters()
	return app
}

type appConfig struct {
	app.AppConfig
	Env *eth.Config
}

// 配置完毕，做一些运行时检查与初始化
func (p *appConfig) Start() {
	processor.InitEthClient()
}

func (p *appConfig) ConfigRoutes() {
	p.Routers(util.VersionRouter())
}

func (p *appConfig) ConfigDB() *sqlx.DB {
	return p.ConnectDB()
}

func (p *appConfig) ConfigFilters() {
	p.Filter(filters.Cors)
}
