package app

import (
	"cos-backend-com/src/account"
	"cos-backend-com/src/account/routers/oauth"
	"cos-backend-com/src/account/routers/users"
	"cos-backend-com/src/common/app"
	"cos-backend-com/src/common/providers"
	"cos-backend-com/src/common/providers/session"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/auth"
	"cos-backend-com/src/libs/filters"
	"cos-backend-com/src/libs/sdk/web3"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/mediocregopher/radix.v2/pool"
	s "github.com/wujiu2020/strip"
	"github.com/wujiu2020/strip/caches"
)

const (
	AppName = "account"
)

func AppInit(tea *s.Strip, confPath string, files ...string) *appConfig {
	app := &appConfig{app.New(tea, AppName), account.Env}
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
	Env *account.Config
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
	var rt http.RoundTripper
	if err := p.Injector().Find(&rt, ""); err != nil {
		panic(err)
	}

	var dbconn *sqlx.DB
	p.Injector().Find(&dbconn)

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

	cache, err := caches.NewRedisProvider(caches.RedisConfig{
		KeyPrefix: "cache:",
		Client:    redisPool,
	})
	if err != nil {
		p.Logger().Error("config caches:", err)
		os.Exit(1)
	}
	oauth2TokenURL := p.Env.Service.Account + "/oauth2/token"

	p.Provide(auth.AuthTransportProvider(oauth2TokenURL))
	p.ProvideAs(cache, (*caches.CacheProvider)(nil))
	p.Provide(web3.NewWeb3Service(account.Env.Service.Web3))
	p.Provide(providers.SessionLimiter(redisPool))
}

func (p *appConfig) ConfigFilters() {
	p.Filter(filters.Cors)
}

func (p *appConfig) ConfigRoutes() {
	p.Routers(util.VersionRouter())
	p.Routers(
		s.Router("/login",
			s.Post(users.Guest{}).Action("Login"),
		),
		s.Router("/logout",
			s.Get(users.Guest{}).Action("Logout"),
		),
		s.Router("/nonce",
			s.Post(users.Guest{}).Action("GetNonce"),
		),
		s.Router("/users/me",
			s.Filter(filters.LoginRequiredInner),
			s.Get(users.Users{}).Action("Me"),
			s.Router("/hunters",
				s.Put(users.Users{}).Action("UpdateHunter"),
			),
		),
		s.Router("/users/:userId",
			s.Get(users.Users{}).Action("Get"),
		),
		s.Router("/oauth2/token",
			s.Post(&oauth.Token{}).Action("GrantToken"),
		),
	)
}

func (p *appConfig) ConfigDone() {
}
