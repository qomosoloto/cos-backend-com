package app

import (
	"cos-backend-com/src/libs/filters"
	notifications "cos-backend-com/src/notification"
	"net/http"
	"os"

	"cos-backend-com/src/common/app"
	"cos-backend-com/src/common/providers/session"
	"cos-backend-com/src/common/util"

	"github.com/jmoiron/sqlx"
	"github.com/mediocregopher/radix.v2/pool"
	t "github.com/wujiu2020/strip"
)

const (
	AppName = "notification"
)

func AppInit(tea *t.Strip, confPath string, files ...string) *appConfig {
	app := &appConfig{app.New(tea, AppName), notifications.Env}
	app.ConfigLoad(app.Env, confPath, files...)
	app.ConfigCheck()
	app.ConfigDB()
	app.ConfigProviders()
	app.ConfigFilters()
	app.ConfigRoutes()
	app.ConfigDone()
	return app
}

type appConfig struct {
	app.AppConfig
	Env *notifications.Config
}

// 配置完毕，做一些运行时检查与初始化
func (p *appConfig) Start() {
}

func (p *appConfig) ConfigCheck() {
}

func (p *appConfig) ConfigDB() *sqlx.DB {
	return p.ConnectDB()
}

func (p *appConfig) ConfigProviders() {

	redisPool, err := pool.NewCustom("tcp",
		p.Env.Redis.Addr,
		p.Env.Redis.PoolSize,
		util.RedisDialWithSecret(p.Env.Redis.Secret),
	)
	if err != nil {
		p.Logger().Error("redis pool.New:", err)
		os.Exit(1)
	}
	p.Provide(redisPool)

	if err := session.ConfigSessions(p.Strip, redisPool, p.Env.Session); err != nil {
		p.Logger().Error("config sessions:", err)
		os.Exit(1)
	}

	var rt http.RoundTripper
	if err := p.Injector().Find(&rt, ""); err != nil {
		panic(err)
	}
}

func (p *appConfig) ConfigFilters() {
	p.Filter(filters.Cors)
}

func (p *appConfig) ConfigRoutes() {
	p.Routers(util.VersionRouter())
}

func (p *appConfig) ConfigDone() {
}
