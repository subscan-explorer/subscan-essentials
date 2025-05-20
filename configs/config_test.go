package configs

import (
	"fmt"
	"github.com/itering/subscan/util"
	"net/url"
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

const confExamplePath = "config.yaml.example"

func TestENVSandbox(t *testing.T) {
	_ = os.Setenv("fake_env", "1=1")
	defer func() {
		_ = os.Unsetenv("fake_env")
	}()
	EnvSandbox(func() {
		if os.Getenv("fake_env") == "1=1" {
			t.Fatal("sandbox clear env failed")
		}

		_ = os.Setenv("fake_env", "2")
	})
	if os.Getenv("fake_env") != "1=1" {
		t.Fatal("sandbox not working")
	}
}

func copyExampleConfig(confDir string) {
	confPath := fmt.Sprintf("%s/config.yaml", confDir)
	if err := os.Link(confExamplePath, confPath); err != nil {
		panic(err)
	}
}

func containsStr(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}

	return false
}

func loadEnvs(envs map[string]string, skipped []string) {
	for name, value := range envs {
		if containsStr(skipped, name) {
			continue
		}

		err := os.Setenv(name, value)
		if err != nil {
			panic(err)
		}
	}
}

func writeConfigWithoutInit(t *testing.T) string {
	confDir := t.TempDir()
	copyExampleConfig(confDir)
	err := os.Setenv("CONF_DIR", confDir)
	if err != nil {
		panic(err)
	}
	util.ConfDir = os.Getenv("CONF_DIR")
	return confDir
}

func writeConfigAndInit(t *testing.T) string {
	confDir := writeConfigWithoutInit(t)
	// test beginning
	Init()
	return confDir
}

var fakeEnvs = map[string]string{
	"MYSQL_HOST":           "111.111.111.111",
	"MYSQL_USER":           "fake_db_user",
	"MYSQL_PASS":           "fake_db_pass*123+456[789]",
	"MYSQL_DB":             "fake_db",
	"MYSQL_PORT":           "1111",
	"REDIS_HOST":           "222.222.222.222",
	"REDIS_PORT":           "2222",
	"REDIS_PASSWORD":       "fake_redis_password*123+456",
	"CACHE_REDIS_HOST":     "223.223.223.223",
	"CACHE_REDIS_PORT":     "3333",
	"CACHE_REDIS_PASSWORD": "fake_cache_redis_password*123+456",
	"SENTRY_DSN":           "https://sentry.fake/fakepath?fakeparam",
}

func getFakeDefaultEnvDSN() *url.URL {
	username := fakeEnvs["MYSQL_USER"]
	password := fakeEnvs["MYSQL_PASS"]
	var user = url.User(username)
	if password != "" {
		user = url.UserPassword(username, password)
	}
	return &url.URL{
		Scheme: "mysql",
		User:   user,
		Host:   fmt.Sprintf("%s:%s", fakeEnvs["MYSQL_HOST"], fakeEnvs["MYSQL_PORT"]),
		Path:   fmt.Sprintf("/%s", fakeEnvs["MYSQL_DB"]),
	}
}

func setFakeEnvs(skipped []string) {
	os.Clearenv()
	loadEnvs(fakeEnvs, skipped)
}

// Testing Init

func TestDefaultInit(t *testing.T) {
	EnvSandbox(func() {
		writeConfigAndInit(t)
		// check nil
		if Boot.Database == nil {
			t.Error("Boot.Database is nil")
			return
		}
		if Boot.Redis == nil {
			t.Error("Boot.Redis is nil")
			return
		}

	})
}

func TestFakeEnvInit(t *testing.T) {
	EnvSandbox(func() {
		setFakeEnvs(nil)
		writeConfigAndInit(t)
		// check value
		// THIS required example's value
		yamlDSN, err := ParseDSN(Boot.Database.Mysql.Api)
		if err != nil {
			panic(err)
		}
		rightDSN, err := (&Database{}).Mysql.mergeDefaultDSNs(getFakeDefaultEnvDSN(), yamlDSN)
		if err != nil {
			panic(err)
		}

		// gorm test
		if rightDSN.Scheme == "mysql" {
			user := getUserOfDSN(rightDSN)
			if err != nil {
				panic(err)
			}
			rightDSNMysql := fmt.Sprintf("%s@tcp(%s)%s?%s", user, rightDSN.Host, rightDSN.Path, rightDSN.RawQuery)
			t.Run("Gorm mysql dsn", func(t *testing.T) {
				if Boot.Database.Mysql.DSN != rightDSNMysql {
					t.Logf("Right DSN: %s", rightDSNMysql)
					t.Logf("Test DSN: %s", Boot.Database.Mysql.DSN)
					t.Fatal("unexpected mysql Database.DSN")
				}
			})
		} else {
			t.Run("Other DSN", func(t *testing.T) {
				if Boot.Database.Mysql.DSN != rightDSN.String() {
					t.Logf("Right DSN: %s", rightDSN.String())
					t.Logf("Test DSN: %s", Boot.Database.Mysql.DSN)
					t.Fatal("unexpected Database.DSN")
				}
			})
		}

		if Boot.Redis.Addr != "222.222.222.222:2222" {
			t.Logf("Boot.Redis.Addr is %s", Boot.Redis.Addr)
			t.Fatal("unexpected Redis.Addr")
		}
		if Boot.Redis.Password != "fake_redis_password*123+456" {
			t.Logf("Boot.Redis.Password is %s", Boot.Redis.Password)
			t.Fatal("unexpected Redis.Password")
		}
	})
}

func fakeBootstrap(t *testing.T) Bootstrap {
	confDir := writeConfigAndInit(t)
	var boot Bootstrap
	filename := fmt.Sprintf("%s/config.yaml", confDir)
	c := config.New(config.WithSource(file.NewSource(filename)))
	if err := c.Load(); err != nil {
		panic(err)
	}
	if err := c.Scan(&boot); err != nil {
		panic(err)
	}
	return boot
}

// Testing Database

func TestFakeDBGetEnvDSN(t *testing.T) {
	EnvSandbox(func() {
		setFakeEnvs(nil)
		boot := fakeBootstrap(t)
		env := boot.Database.Mysql.getEnvDSN()
		if env.String() != getFakeDefaultEnvDSN().String() {
			t.Fatalf("unexpected database dsn: %s", env.String())
		}
	})
}

func dsnParseAndString(dsn string) string {
	u, err := ParseDSN(dsn)
	if err != nil {
		panic(err)
	}

	return u.String()
}
func TestFakeDBGetYamlDSN(t *testing.T) {
	EnvSandbox(func() {
		writeConfigWithoutInit(t)
		boot := fakeBootstrap(t)
		env, err := boot.Database.Mysql.getYamlDSN()
		if err != nil {
			t.Fatal(err)
		}

		if env.String() != dsnParseAndString(boot.Database.Mysql.Api) {
			t.Fatalf("unexpected database api dsn: %s", env.String())
		}
		os.Clearenv()
	})
}

func TestMergeDefaultDSNsFunc(t *testing.T) {
	EnvSandbox(func() {
		db := &Database{}
		testEnv, err := ParseDSN("env://fake_user1:fake_pass1@fake_host:1111/fake_db1?fake_param=1")
		if err != nil {
			panic(err)
		}
		testYaml, err := ParseDSN("yaml://fake_user:fake_pass@fake_host:2222/fake_db?fake_param=0")
		if err != nil {
			panic(err)
		}

		t.Run("Test merge order Env final", func(t *testing.T) {
			ns, err := db.Mysql.mergeDefaultDSNs(testYaml, testEnv)
			if err != nil {
				t.Fatal(err)
			}

			// test order
			if ns.Scheme != testEnv.Scheme {
				t.Fatalf("unexpected scheme: %s", ns.Scheme)
			}
			if ns.User.String() != testEnv.User.String() {
				t.Fatalf("unexpected user: %s", ns.User.String())
			}
			if ns.Host != testEnv.Host {
				t.Fatalf("unexpected host: %s", ns.Host)
			}
			if ns.Path != testEnv.Path {
				t.Fatalf("unexpected path: %s", ns.Host)
			}
			if ns.RawQuery != testEnv.RawQuery {
				t.Fatalf("unexpected path: %s", ns.RawQuery)
			}
		})

		t.Run("Test merge order Yaml final", func(t *testing.T) {
			ns, err := db.Mysql.mergeDefaultDSNs(testEnv, testYaml)
			if err != nil {
				t.Fatal(err)
			}

			// test order
			if ns.Scheme != testYaml.Scheme {
				t.Fatalf("unexpected scheme: %s", ns.Scheme)
			}
			if ns.User.String() != testYaml.User.String() {
				t.Fatalf("unexpected user: %s", ns.User.String())
			}
			if ns.Host != testYaml.Host {
				t.Fatalf("unexpected host: %s", ns.Host)
			}
			if ns.Path != testYaml.Path {
				t.Fatalf("unexpected path: %s", ns.Host)
			}
			if ns.RawQuery != testYaml.RawQuery {
				t.Fatalf("unexpected path: %s", ns.RawQuery)
			}
		})

		t.Run("Test one nil param A", func(t *testing.T) {
			ns, err := db.Mysql.mergeDefaultDSNs(testEnv, nil)
			if err != nil {
				t.Fatal(err)
			}
			if ns.String() != testEnv.String() {
				t.Fatalf("unexpected value: %s", ns.String())
			}
		})

		t.Run("Test one nil param B", func(t *testing.T) {
			ns, err := db.Mysql.mergeDefaultDSNs(nil, testEnv)
			if err != nil {
				t.Fatal(err)
			}
			if ns.String() != testEnv.String() {
				t.Fatalf("unexpected value: %s", ns.String())
			}
		})

		t.Run("Test default override", func(t *testing.T) {
			emptyDSN, err := ParseDSN("")
			if err != nil {
				panic(err)
			}

			ns, err := db.Mysql.mergeDefaultDSNs(emptyDSN, nil)
			if err != nil {
				t.Fatal(err)
			}

			defaultDSN, err := ParseDSN("mysql://root@127.0.0.1:3306/subscan-essentials")
			if err != nil {
				panic(err)
			}
			if ns.String() != defaultDSN.String() {
				t.Fatalf("unexpected value: %s", ns.String())
			}
		})

		t.Run("Test missing envs", func(t *testing.T) {
			emptyDSN := &url.URL{}
			EnvSandbox(func() {
				// Mysql port missing
				if err := os.Setenv("MYSQL_HOST", "127.0.0.1"); err != nil {
					panic(err)
				}
				testDSN := &url.URL{}
				ns, err := db.Mysql.mergeDefaultDSNs(emptyDSN, testDSN)
				if err != nil {
					t.Fatal(err)
				}

				if ns.Port() != "3306" {
					if ns.Port() == "" {
						t.Fatal("unexpected value: <empty>")
					} else {
						t.Fatalf("unexpected value: %s", ns.Port())
					}

				}
			})
		})
	})
}

// Testing Redis

func TestRedisMergeNothing(t *testing.T) {
	EnvSandbox(func() {
		boot := fakeBootstrap(t)
		redis := boot.Redis
		backup := *boot.Redis
		os.Clearenv()
		redis.mergeEnvironment()
		if redis.Addr != backup.Addr {
			t.Fatalf("unexpected value: %s", redis.Addr)
		}
		if redis.DbName != backup.DbName {
			t.Fatalf("unexpected value: %d", redis.DbName)
		}
		if redis.Password != backup.Password {
			t.Fatalf("unexpected value: %s", redis.Password)
		}
	})
}

func TestRedisMergeEnv(t *testing.T) {
	EnvSandbox(func() {
		boot := fakeBootstrap(t)
		redis := boot.Redis
		setFakeEnvs(nil)
		redis.mergeEnvironment()
		if redis.Addr != fmt.Sprintf("%s:%s", fakeEnvs["REDIS_HOST"], fakeEnvs["REDIS_PORT"]) {
			t.Fatalf("unexpected value: %s", redis.Addr)
		}
		if redis.DbName != util.StringToInt(fakeEnvs["REDIS_DATABASE"]) {
			t.Fatalf("unexpected value: %d", redis.DbName)
		}
		if redis.Password != fakeEnvs["REDIS_PASSWORD"] {
			t.Fatalf("unexpected value: %s", redis.Password)
		}
	})
}

func TestRedisMissingEnv(t *testing.T) {
	EnvSandbox(func() {
		// Test missing REDIS_PORT
		err := os.Setenv("REDIS_HOST", "111.111.111.111")
		if err != nil {
			panic(err)
		}
		redis := &Redis{}
		redis.mergeEnvironment()
		if redis.Addr != "111.111.111.111:6379" {
			t.Fatalf("unexpected value: %s", redis.Addr)
		}
	})
}

func TestRedisMergeDefault(t *testing.T) {
	EnvSandbox(func() {
		emptyRedis := &Redis{}
		emptyRedis.mergeEnvironment()
		if emptyRedis.Addr != "127.0.0.1:6379" {
			t.Fatalf("unexpected value: %s", emptyRedis.Addr)
		}
		if emptyRedis.DbName != 0 {
			t.Fatalf("unexpected value: %d", emptyRedis.DbName)
		}
		if emptyRedis.Password != "" {
			t.Fatalf("unexpected value: %s", emptyRedis.Password)
		}
	})
}

func TestGetUserOfDSN(t *testing.T) {
	dsn := &url.URL{}
	dsn.User = url.User("test_user")
	if u := getUserOfDSN(dsn); u != "test_user" {
		t.Fatalf("unexpected value: %s", u)
	}

	dsn.User = url.UserPassword("test_user", "")
	if u := getUserOfDSN(dsn); u != "test_user:" {
		t.Fatalf("unexpected value: %s", u)
	}

	dsn.User = url.UserPassword("test_user", "test_pass+123*456")
	if u := getUserOfDSN(dsn); u != "test_user:test_pass+123*456" {
		t.Fatalf("unexpected value: %s", u)
	}
}

func TestParseDSN(t *testing.T) {
	t.Run("Normal dsn", func(t *testing.T) {
		dsn, err := ParseDSN("mysql://user:password@host:1111/dbname?param=1")
		if err != nil {
			t.Fatal(err)
		}

		if dsn.Scheme != "mysql" ||
			dsn.User.String() != "user:password" ||
			dsn.Hostname() != "host" || dsn.Port() != "1111" ||
			dsn.Path != "/dbname" ||
			dsn.RawQuery != "param=1" {
			t.Fatalf("unexpected value: %s", dsn.String())
		}
	})

	t.Run("Special chars password", func(t *testing.T) {
		dsn, err := ParseDSN("mysql://user:pass*1+2/3(4)[5]%6@7@#$%^&{}[];:\"<>,?/\\.!~@host:1111/dbname?param=1")
		if err != nil {
			t.Fatal(err)
		}

		password, ok := dsn.User.Password()
		if !ok {
			t.Fatal("unexpected password")
		}
		if dsn.Scheme != "mysql" ||
			dsn.User.Username() != "user" ||
			password != "pass*1+2/3(4)[5]%6@7@#$%^&{}[];:\"<>,?/\\.!~" ||
			dsn.Hostname() != "host" || dsn.Port() != "1111" ||
			dsn.Path != "/dbname" ||
			dsn.RawQuery != "param=1" {
			t.Fatalf("unexpected value: %s", dsn.String())
		}
	})

	t.Run("Only Params", func(t *testing.T) {
		dsn, err := ParseDSN("?parseTime=true&loc=Local&charset=utf8mb4,utf8")
		if err != nil {
			t.Fatal(err)
		}

		password, ok := dsn.User.Password()
		if ok {
			t.Fatal("unexpected password")
		}
		if dsn.Scheme != "" ||
			dsn.User.Username() != "" ||
			password != "" ||
			dsn.Hostname() != "" || dsn.Port() != "" ||
			dsn.Path != "" ||
			dsn.RawQuery != "parseTime=true&loc=Local&charset=utf8mb4,utf8" {
			t.Fatalf("unexpected value: %s", dsn.String())
		}
	})

	t.Run("Missing password", func(t *testing.T) {
		dsn, err := ParseDSN("mysql://user@host:1111/dbname?param=1")
		if err != nil {
			t.Fatal(err)
		}

		password, ok := dsn.User.Password()
		if ok {
			t.Fatal("unexpected password")
		}
		if dsn.Scheme != "mysql" ||
			dsn.User.Username() != "user" ||
			password != "" ||
			dsn.Hostname() != "host" || dsn.Port() != "1111" ||
			dsn.Path != "/dbname" ||
			dsn.RawQuery != "param=1" {
			t.Fatalf("unexpected value: %s", dsn.String())
		}
	})

	t.Run("Empty DSN", func(t *testing.T) {
		dsn, err := ParseDSN("")
		if err != nil {
			t.Fatal(err)
		}

		password, ok := dsn.User.Password()
		if ok {
			t.Fatal("unexpected password")
		}
		if dsn.Scheme != "" ||
			dsn.User.Username() != "" ||
			password != "" ||
			dsn.Hostname() != "" || dsn.Port() != "" ||
			dsn.Path != "" ||
			dsn.RawQuery != "" {
			t.Fatalf("unexpected value: %s", dsn.String())
		}
	})
}
