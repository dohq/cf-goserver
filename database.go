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
