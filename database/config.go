package database

import (
	"fmt"
	"strings"
)

type Config struct {
	User     string
	Password string
	Name     string
}

func (c Config) DSN() string {
	var params []string

	if c.User != "" {
		params = append(params, fmt.Sprintf("user=%s", c.User))
	}

	if c.Password != "" {
		params = append(params, fmt.Sprintf("password=%s", c.Password))
	}

	if c.Name != "" {
		params = append(params, fmt.Sprintf("dbname=%s", c.Name))
	}

	params = append(params, "sslmode=disable")

	return strings.Join(params, " ")
}
