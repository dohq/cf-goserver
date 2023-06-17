package main

import (
	"net/http"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
)

func (c *Client) GetEnv(w http.ResponseWriter, r *http.Request) {
	if cfenv.IsRunningOnCF() {
		appEnv := cfenv.CurrentEnv()

		bytes, err := json.Marshal(appEnv)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			c.Logger.Error("failed marshal env", zap.Error(err))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}
