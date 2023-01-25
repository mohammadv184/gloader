package mysql

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Options  url.Values
	Protocol string
}

const (
	TCPProtocol    = "tcp"
	SocketProtocol = "unix"
)

func (c *Config) String() string {
	var str strings.Builder
	str.WriteString(c.Username)
	if c.Password != "" {
		str.WriteString(":")
		str.WriteString(c.Password)
	}
	str.WriteString("@")
	switch c.Protocol {
	case TCPProtocol:
		str.WriteString("tcp(")
		str.WriteString(c.Host)
		str.WriteString(":")
		str.WriteString(strconv.Itoa(c.Port))
		str.WriteString(")")
	case SocketProtocol:
		str.WriteString("unix(")
		str.WriteString(c.Host)
		str.WriteString(")")
	}
	str.WriteString("/")
	str.WriteString(c.Database)
	if len(c.Options) > 0 {
		str.WriteString("?")
		str.WriteString(c.Options.Encode())
	}
	return str.String()
}

// parseConfig parses a DSN string into a Config.
// example: user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
func parseConfig(name string) (*Config, error) {
	config := &Config{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "",
		Database: "",
		Options:  make(url.Values),
		Protocol: TCPProtocol,
	}

	regex := regexp.MustCompile(`^(?P<user>[^:]+)(:(?P<password>[^@]+))?@((?P<protocol>[^()]+)?\()?(?P<host>[^:]+)(:(?P<port>[0-9]+))?\)?(/(?P<database>[^?]+)?(\?(?P<options>.+))?)?`)
	match := regex.FindStringSubmatch(name)
	result := make(map[string]string)
	for i, n := range regex.SubexpNames() {
		if i != 0 && n != "" && match[i] != "" {
			result[n] = match[i]
		}
	}

	if value, ok := result["user"]; ok {
		config.Username = value
	}
	if value, ok := result["password"]; ok {
		config.Password = value
	}
	if value, ok := result["protocol"]; ok {
		config.Protocol = value
	}
	if value, ok := result["host"]; ok {
		config.Host = value
	}
	if value, ok := result["port"]; ok {
		port, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %v", err)
		}
		config.Port = int(port)
	}
	if value, ok := result["database"]; ok {
		config.Database = value
	}
	if value, ok := result["options"]; ok {
		options, err := url.ParseQuery(value)
		if err != nil {
			return nil, fmt.Errorf("invalid options: %v", err)
		}
		config.Options = options
	}
	return config, nil
}
