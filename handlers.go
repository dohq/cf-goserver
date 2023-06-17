package main

import (
	"net/http"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

func (c *Client) GetEnv(w http.ResponseWriter, r *http.Request) {
	if cfenv.IsRunningOnCF() {
		appEnv, err := cfenv.Current()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.Logger.Error("failed get env", zap.Error(err))
			return
		}

		mysqlService, err := appEnv.Services.WithTag("mysql")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.Logger.Error("failed get mysql service", zap.Error(err))
			return
		}

		bytes, err := json.Marshal(mysqlService)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}
