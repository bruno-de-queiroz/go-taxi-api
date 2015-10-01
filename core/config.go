package core

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/olebedev/config"
	"os"
	"regexp"
)

type Enviroment int

const (
	PRODUCTION Enviroment = iota
	DEVELOPMENT
	TEST
)

var (
	enviroments = [3]string{"production", "development", "test"}
	regex = regexp.MustCompile(`(.*)\$([^\n\t\s]+)`)
)

func (e Enviroment) String() string {
	return enviroments[e]
}

type LogConfig struct {
	Rotate bool
	Level  string
	File   string
	Format string
}

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Sender   string
}

type DatabaseConfig struct {
	Host     string
	Name     string
	Username string
	Password string
}

type TokenConfig struct {
	Expiration string
	Secret     string
}

type Config struct {
	*config.Config
	Enviroment string
	Database   *DatabaseConfig
	Token      *TokenConfig
	Email      *EmailConfig
	Log        *LogConfig
	Name       string
	Version    string
	Port       string
	Mode       string
}

func (cfg *Config) ExtendWithFile(p string) (err error) {
	data, err := parseEnvs(p)
	if err != nil {
		return err
	}

	cc, err := config.ParseYaml(string(data))
	if err != nil {
		return err
	}

	cc, err = cc.Get(cfg.Enviroment)
	if err != nil {
		return err
	}

	reg := regexp.MustCompile(`([\w\d\_\-]+)\.yml$`)
	ps := reg.FindAllStringSubmatch(p, -1)

	return cfg.Set(ps[0][1], cc.Root)
}

func parseEnvs(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	data := make([]byte, 0)
	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		parsed := scanner.Text() + "\n"
		p := regex.FindAllStringSubmatch(parsed, -1)
		if len(p) > 0 {
			parsed = fmt.Sprintf("%s\"%s\"\n", p[0][1], os.Getenv(p[0][2]))
		}

		data = append(data, []byte(parsed)...)
	}

	return data, nil
}

//TODO Melhorar parser do config
func NewConfig(e string, p string) (c *Config, err error) {
	var mode string

	switch e {
	case PRODUCTION.String():
		mode = "release"
	case DEVELOPMENT.String(), TEST.String():
		mode = "debug"
	default:
		return nil, errors.New("No valid enviroment.")
	}

	data, err := parseEnvs(p)
	if err != nil {
		return nil, err
	}

	cfg, err := config.ParseYaml(string(data))
	if err != nil {
		return nil, err
	}

	lg, err := cfg.Get("log")
	if err != nil {
		return nil, err
	}

	em, err := cfg.Get("email")
	if err != nil {
		return nil, err
	}

	db, err := cfg.Get("database")
	if err != nil {
		return nil, err
	}

	tk, err := cfg.Get("token")
	if err != nil {
		return nil, err
	}

	c = &Config{
		cfg,
		e,
		&DatabaseConfig{db.UString("host"), db.UString("name"), db.UString("username"), db.UString("password")},
		&TokenConfig{tk.UString("expiration"), tk.UString("secret")},
		&EmailConfig{em.UString("host"), em.UString("port"), em.UString("username"), em.UString("password"), em.UString("sender")},
		&LogConfig{lg.UBool("rotate"), lg.UString("level"), lg.UString("file"), lg.UString("format")},
		cfg.UString("name"),
		cfg.UString("version"),
		cfg.UString("port"),
		mode,
	}

	return c, nil
}
