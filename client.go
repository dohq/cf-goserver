package main

import (
	"os"

	"github.com/dohq/go-cfserver/ent"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Client struct {
	Logger *zap.Logger
	Db     *ent.Client
}

func NewClient(isCF bool, ws *os.File, loglevel zapcore.LevelEnabler) (*Client, error) {
	ec := ecszap.EncoderConfig{
		EnableName:       false,
		EnableStackTrace: true,
		EnableCaller:     false,
		EncodeName:       zapcore.FullNameEncoder,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeDuration:   zapcore.NanosDurationEncoder,
		EncodeCaller:     ecszap.ShortCallerEncoder,
	}

	core := ecszap.NewCore(ec, ws, loglevel)
	logger := zap.New(core, zap.AddCaller()).Named("cf-goserver")

	db, err := newDBClient(isCF)
	if err != nil {
		return nil, err
	}

	return &Client{
		Logger: logger,
		Db:     db,
	}, nil
}
