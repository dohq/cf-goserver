package main

import (
	"net/http"
	"strconv"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/dohq/go-cfserver/ent"
	"github.com/dohq/go-cfserver/ent/user"
	"github.com/go-chi/render"
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

func (c *Client) UserAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	var payload ent.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.Logger.Error("failed parse body", zap.Error(err))
		return
	}

	u, err := c.Db.User.Create().
		SetName(payload.Name).
		SetAge(payload.Age).
		SetUID(payload.UID).
		Save(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.Logger.Error("failed insert user record", zap.Error(err))
		return
	}
	c.Logger.Info("insert ok", zap.Int("id", u.ID), zap.String("name", payload.Name))
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, u)
}

func (c *Client) UserGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.Logger.Error("failed convert age", zap.Error(err))
		return
	}

	u, err := c.Db.User.Query().
		Where(user.ID(id)).
		Only(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.Logger.Error("could not find user", zap.Error(err))
		return
	}

	c.Logger.Info("select ok", zap.Int("id", u.ID), zap.String("name", u.Name))
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, u)
}
