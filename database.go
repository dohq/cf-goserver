package main

import (
	"context"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/dohq/go-cfserver/ent"
	"github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"

	_ "github.com/xiaoqidun/entps"
)

type ServiceCredentials struct {
	Hostname string      `json:"hostname,omitempty"`
	Port     json.Number `json:"port,omitempty"`
	Name     string      `json:"name,omitempty"`
	Username string      `json:"username,omitempty"`
	Password string      `json:"password,omitempty"`
	URI      string      `json:"uri,omitempty"`
	JdbcURL  string      `json:"jdbcUrl,omitempty"`
}

func newDBClient(isCF bool) (*ent.Client, error) {
	var c *ent.Client
	var err error

	if isCF {
		appEnv, err := cfenv.Current()
		if err != nil {
			return nil, err
		}

		service, err := appEnv.Services.WithTag("mysql")
		if err != nil {
			return nil, err
		}

		cds, err := json.Marshal(service[0].Credentials)
		if err != nil {
			return nil, err
		}

		var creds ServiceCredentials
		if err := json.Unmarshal(cds, &creds); err != nil {
			return nil, err
		}

		dsn := mysql.Config{
			DBName:               creds.Name,
			Addr:                 creds.Hostname + ":" + creds.Port.String(),
			User:                 creds.Username,
			Passwd:               creds.Password,
			Net:                  "tcp",
			ParseTime:            true,
			Collation:            "utf8mb4_unicode_ci",
			AllowNativePasswords: true,
		}

		drv, err := sql.Open(dialect.MySQL, dsn.FormatDSN())
		if err != nil {
			return nil, err
		}

		db := drv.DB()
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		c = ent.NewClient(ent.Driver(drv))

	} else {
		c, err = ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
		if err != nil {
			return nil, err
		}
	}

	if err = c.Schema.Create(context.Background(), schema.WithAtlas(true)); err != nil {
		return nil, err
	}

	return c, nil
}
