package infrastructure

import (
	"bytes"
	"fmt"
	"text/template"
	"user-management/config"
)

const MySQLURITemplate = "{{.Username}}:{{.Password}}@tcp({{.Server}})/{{.Database}}"

type DBUri struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Server   string `json:"server"`
	Database string `json:"database"`
}

func NewMySQLURI() DBUri {
	config := config.GetConfig()
	return DBUri{
		Username: config.MySQLUsername,
		Password: config.MySQLPassword,
		Server:   config.DBContainer,
		Database: config.MySQLDatabase,
	}
}

func (r DBUri) newTemplate() (*template.Template, error) {
	tmpl, err := template.New("mysql_db_uri").Parse(MySQLURITemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mysql db uri template: %w", err)
	}
	return tmpl, nil
}

func (r DBUri) Generate() (string, error) {
	template, err := r.newTemplate()
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = template.Execute(&buffer, r)
	if err != nil {
		return "", fmt.Errorf("failed to execute weather mysql db uri template: %w", err)
	}
	return buffer.String(), nil
}
