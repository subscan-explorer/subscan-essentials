package configs

import (
	"errors"
	"fmt"
	"github.com/itering/subscan/util"
	"net/url"
	"os"
	"strings"

	xtime "github.com/itering/subscan/pkg/time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

type Bootstrap struct {
	Server   *Server   `json:"server,omitempty"`
	Database *Database `json:"database,omitempty"`
	Redis    *Redis    `json:"redis,omitempty"`
	UI       *UI       `json:"ui,omitempty"`
}

type Server struct {
	Http *ServerHttp `json:"http,omitempty"`
	Grpc *ServerGrpc `json:"grpc,omitempty"`
}

type UI struct {
	EnableSubstrate bool `json:"enable_substrate"`
	EnableEvm       bool `json:"enable_evm"`
}

type ServerHttp struct {
	Network string `json:"network,omitempty"`
	Addr    string `json:"addr,omitempty"`
	Timeout string `json:"timeout,omitempty"`
}

type ServerGrpc struct {
	Addr string `json:"addr,omitempty"`
}

type IDatabase interface {
	mergeEnvironment()
	GetHost() string
}

type Mysql struct {
	Host     string
	DSN      string   `json:"-"`
	Multiple []string `json:"-"`
	Api      string   `json:"api"`
	Test     string   `json:"test"`
}

func (dc *Mysql) GetHost() string {
	return dc.Host
}

type Postgres struct {
	DSN string `json:"-"`
	Api string `json:"api"`

	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	Multiple []string
}

type Database struct {
	Driver   string    `json:"-"` // Unused
	Mysql    *Mysql    `json:"mysql"`
	Postgres *Postgres `json:"postgres"`
}

type Redis struct {
	Proto        string         `json:"proto"`
	Addr         string         `json:"addr"`
	DbName       int            `json:"db_name"`
	Password     string         `json:"password"`
	Idle         int            `json:"idle"`
	Active       int            `json:"active"`
	ReadTimeout  xtime.Duration `json:"read_timeout"`
	WriteTimeout xtime.Duration `json:"write_timeout"`
}

var Boot Bootstrap

func Init() {
	configFilename := fmt.Sprintf("%s/config.yaml", util.ConfDir)
	if !util.FileExists(configFilename) {
		panic(fmt.Errorf("config file %s is not exists", configFilename))
	}

	c := config.New(config.WithSource(file.NewSource(configFilename)))
	if err := c.Load(); err != nil {
		panic(err)
	}
	if err := c.Scan(&Boot); err != nil {
		panic(err)
	}

	if Boot.Database == nil || Boot.Redis == nil {
		panic(fmt.Errorf("config.yaml not completed"))
	}

	// db driver
	Boot.Database.Driver = util.GetEnv("DB_DRIVER", "mysql")
	if Boot.Database.Driver == "mysql" {
		Boot.Database.Mysql.mergeEnvironment()
	} else if Boot.Database.Driver == "postgres" {
		Boot.Database.Postgres.mergeEnvironment()
	} else {
		panic(fmt.Errorf("unsupported db driver: %s", Boot.Database.Driver))
	}

	Boot.Redis.mergeEnvironment()
}

func setVarDefaultValueStr(variable *string, defaultValue string) {
	if variable != nil && *variable == "" {
		*variable = defaultValue
	}
}

func (dc *Mysql) mergeEnvironment() {
	var (
		err                  error
		dsn, envDSN, fileDSN *url.URL
	)
	envDSN = dc.getEnvDSN()
	fileDSN, err = dc.getYamlDSN()
	if err != nil {
		panic(err)
	}

	dsn, err = dc.mergeDefaultDSNs(fileDSN, envDSN)
	if err != nil {
		panic(err)
	}

	var dsnStr = dsn.String()
	dc.Host = dsn.Host
	// for gorm
	if dsn.Scheme == "mysql" {
		dsnStr = fmt.Sprintf("%s@tcp(%s)%s?%s", getUserOfDSN(dsn), dsn.Host, dsn.Path, dsn.RawQuery)
	}

	dc.DSN = dsnStr
}

func (d *Database) GetHost() string {
	if Boot.Database.Driver == "mysql" {
		return Boot.Database.Mysql.Host
	} else if Boot.Database.Driver == "postgres" {
		return Boot.Database.Postgres.GetHost()
	}
	return ""

}

func (dc *Mysql) getEnvDSN() *url.URL {
	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbName := os.Getenv("MYSQL_DB")
	dbPort := os.Getenv("MYSQL_PORT")

	var user *url.Userinfo
	if dbUser != "" && dbPass != "" {
		user = url.UserPassword(dbUser, dbPass)
	} else {
		user = url.User(dbUser)
	}

	addr := fmt.Sprintf("%s:%s", dbHost, dbPort)
	if dbHost == "" && dbPort == "" {
		addr = ""
	}

	return &url.URL{
		Scheme: "mysql",
		User:   user,
		Host:   addr,
		Path:   fmt.Sprintf("/%s", dbName),
	}
}

func (dc *Mysql) getYamlDSN() (*url.URL, error) {
	var (
		err error
		dsn *url.URL
	)

	dsn, err = ParseDSN(dc.Api)

	if err != nil {
		return nil, err
	}

	return dsn, err
}

func valueOrDefaultStr(defaultValue string, values ...string) string {
	var final string
	for _, value := range values {
		if value != "" {
			final = value
		}
	}

	if final != "" {
		return final
	}

	return defaultValue
}

func getDBNameOfDSN(dsn *url.URL) string {
	return strings.TrimPrefix(dsn.Path, "/")
}

func getUserOfDSN(dsn *url.URL) string {
	user := dsn.User.Username()
	if password, ok := dsn.User.Password(); ok {
		user = fmt.Sprintf("%s:%s", dsn.User.Username(), password)
	}
	return user
}

func (_ *Mysql) mergeDefaultDSNs(a, b *url.URL) (*url.URL, error) {
	if a == nil && b == nil {
		return nil, errors.New("must have least one is non-nil")
	}
	var emptyUrl = &url.URL{}
	if a == nil {
		a = emptyUrl
	}
	if b == nil {
		b = emptyUrl
	}

	scheme := valueOrDefaultStr("mysql", a.Scheme, b.Scheme)
	user := valueOrDefaultStr("root", getUserOfDSN(a), getUserOfDSN(b))
	var userInfo *url.Userinfo
	if strings.Contains(user, ":") {
		s := strings.SplitN(user, ":", 2)
		userInfo = url.UserPassword(s[0], s[1])
	} else {
		userInfo = url.User(user)
	}
	host := valueOrDefaultStr("127.0.0.1", a.Hostname(), b.Hostname())
	port := valueOrDefaultStr("3306", a.Port(), b.Port())
	dbPath := valueOrDefaultStr("subscan-essentials", getDBNameOfDSN(a), getDBNameOfDSN(b))
	query := valueOrDefaultStr("", a.RawQuery, b.RawQuery)
	return &url.URL{
		Scheme:   scheme,
		User:     userInfo,
		Host:     fmt.Sprintf("%s:%s", host, port),
		Path:     fmt.Sprintf("/%s", dbPath),
		RawQuery: query,
	}, nil
}

func (rc *Redis) mergeEnvironment() {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost != "" {
		redisHost = fmt.Sprintf("%s:%s", redisHost, util.GetEnv("REDIS_PORT", "6379"))
	} else {
		redisHost = rc.Addr
	}
	setVarDefaultValueStr(&redisHost, "127.0.0.1:6379")

	rc.Addr = redisHost
	rc.DbName = util.StringToInt(util.GetEnv("REDIS_DATABASE", util.IntToString(rc.DbName)))
	rc.Password = util.GetEnv("REDIS_PASSWORD", rc.Password)
}

func ParseDSN(dsn string) (*url.URL, error) {
	foundKey := false
	extendScheme := ""
	var start, end = 0, 0
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '@' {
			foundKey = true
			end = i
			break
		}
	}

	if i := strings.Index(dsn[:end], "://"); i > 0 {
		start = i + 3
	}

	if foundKey {
		start = start + strings.Index(dsn[start:end], ":") + 1
		return url.Parse(extendScheme + dsn[:start] + url.QueryEscape(dsn[start:end]) + dsn[end:])
	}

	return url.Parse(dsn)
}

func (p Postgres) GetHost() string {
	return p.Host
}

func (p *Postgres) mergeEnvironment() {
	p.Host = valueOrDefaultStr("127.0.0.1", p.Host, os.Getenv("POSTGRES_HOST"))
	p.User = valueOrDefaultStr("gorm", p.User, os.Getenv("POSTGRES_USER"))
	p.Password = valueOrDefaultStr("gorm", p.Password, os.Getenv("POSTGRES_PASS"))
	p.DBName = valueOrDefaultStr("subscan-essentials", p.DBName, os.Getenv("POSTGRES_DB"))
	p.Port = valueOrDefaultStr("9920", p.Port, os.Getenv("POSTGRES_PORT"))
	p.SSLMode = valueOrDefaultStr("disable", p.SSLMode, os.Getenv("POSTGRES_SSL_MODE"))

	// host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disablea
	p.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", p.Host, p.User, p.Password, p.DBName, p.Port, p.SSLMode)
}

func EnvSandbox(fn func()) {
	originEnvsBackup := os.Environ()
	os.Clearenv()
	defer func() {
		os.Clearenv()
		for _, s := range originEnvsBackup {
			env := strings.SplitN(s, "=", 2)
			_ = os.Setenv(env[0], env[1])
		}
	}()
	fn()
}
