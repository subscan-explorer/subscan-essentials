package configs

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	xtime "github.com/itering/subscan/pkg/time"
	"golang.org/x/exp/slog"

	"github.com/itering/subscan/util"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Bootstrap struct {
	Server   *Server   `json:"server,omitempty"`
	Database *Database `json:"database,omitempty"`
	Redis    *Redis    `json:"redis,omitempty"`
}

type Server struct {
	Http *ServerHttp `json:"http,omitempty"`
}

type ServerHttp struct {
	Network string `json:"network,omitempty"`
	Addr    string `json:"addr,omitempty"`
	Timeout string `json:"timeout,omitempty"`
}

type Database struct {
	Driver string `json:"-"` // Unused
	DSN    string `json:"-"`
	Api    string `json:"api"`
	Task   string `json:"task"`
	Test   string `json:"test"`
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
	conf := koanf.New(".")
	if err := conf.Load(file.Provider(configFilename), yaml.Parser()); err != nil {
		panic("load config file error: " + err.Error())
	}
	if err := conf.Unmarshal("", &Boot); err != nil {
		panic("unmarshal config file error: " + err.Error())
	}

	if Boot.Database == nil || Boot.Redis == nil {
		panic(fmt.Errorf("config.yaml not completed"))
	}

	slog.Info("config", "boot", fmt.Sprintf("%+v", Boot))

	Boot.Database.mergeEnvironment()
	Boot.Redis.mergeEnvironment()
}

func setVarDefaultValueStr(variable *string, defaultValue string) {
	if variable != nil && *variable == "" {
		*variable = defaultValue
	}
}

func (dc *Database) mergeEnvironment() {
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

	dsnStr := dsn.String()

	// for gorm
	if dsn.Scheme == "mysql" {
		dsnStr = fmt.Sprintf("%s@tcp(%s)%s?%s", getUserOfDSN(dsn), dsn.Host, dsn.Path, dsn.RawQuery)
	}
	dc.DSN = dsnStr
	dc.Driver = dsn.Scheme
}

func (dc *Database) getEnvDSN() *url.URL {
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

func (dc *Database) getYamlDSN() (*url.URL, error) {
	var (
		isTaskMode = os.Getenv("TASK_MOD") == "true"
		isTestMode = os.Getenv("TEST_MOD") == "true"
		err        error
		dsn        *url.URL
	)

	if isTaskMode {
		dsn, err = ParseDSN(dc.Task)
	} else if isTestMode {
		dsn, err = ParseDSN(dc.Test)
	} else {
		dsn, err = ParseDSN(dc.Api)
	}

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

func (*Database) mergeDefaultDSNs(a, b *url.URL) (*url.URL, error) {
	if a == nil && b == nil {
		return nil, errors.New("must have least one is non-nil")
	}
	emptyUrl := &url.URL{}
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
	dbPath := valueOrDefaultStr("subscan", getDBNameOfDSN(a), getDBNameOfDSN(b))
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
	start, end := 0, 0
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
