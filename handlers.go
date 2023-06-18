package main

import (
	"net/http"
	"strconv"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/dohq/go-cfserver/ent/user"
	"github.com/go-chi/render"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
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

	name := r.URL.Query().Get("name")

	uid, err := uuid.Parse(r.URL.Query().Get("uid"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.Logger.Error("failed parse uid", zap.Error(err))
		return
	}

	age, err := strconv.Atoi(r.URL.Query().Get("age"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.Logger.Error("failed convert age", zap.Error(err))
		return
	}

	u, err := c.Db.User.Create().
		SetName(name).
		SetAge(age).
		SetUID(uid).
		Save(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.Logger.Error("failed insert user record", zap.Error(err))
		return
	}
	c.Logger.Info("insert ok", zap.Int("id", u.ID), zap.String("username", name))
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

	c.Logger.Info("select ok", zap.Int("id", u.ID), zap.String("username", u.Name))
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, u)
}
