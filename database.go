package main

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/dohq/go-cfserver/ent"

	_ "github.com/xiaoqidun/entps"
)

type MySQLService []struct {
	Name        string   `json:"Name"`
	Label       string   `json:"Label"`
	Tags        []string `json:"Tags"`
	Plan        string   `json:"Plan"`
	Credentials struct {
		Hostname string `json:"hostname"`
		JdbcURL  string `json:"jdbcUrl"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		URI      string `json:"uri"`
		Username string `json:"username"`
	} `json:"Credentials"`
	VolumeMounts []any `json:"VolumeMounts"`
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

		dsn, ok := service[0].CredentialString("uri")
		if ok != true {
			return nil, fmt.Errorf("failed get dsn from environment: %v", err)
		}

		drv, err := sql.Open(dialect.MySQL, dsn)
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
