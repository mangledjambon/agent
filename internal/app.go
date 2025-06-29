// Copyright (c) F5, Inc.
//
// This source code is licensed under the Apache License, Version 2.0 license found in the
// LICENSE file in the root directory of this source tree.

package internal

import (
	"context"
	"log/slog"

	"github.com/nginx/agent/v3/internal/bus"
	"github.com/nginx/agent/v3/internal/config"
	"github.com/nginx/agent/v3/internal/plugin"
	"github.com/spf13/cobra"
)

const (
	defaultMessagePipeChannelSize = 100
	defaultQueueSize              = 100
)

type App struct {
	commit  string
	version string
}

func NewApp(commit, version string) *App {
	return &App{commit, version}
}

func (a *App) Run(ctx context.Context) error {
	config.Init(a.version, a.commit)

	config.RegisterRunner(func(_ *cobra.Command, _ []string) {
		err := config.RegisterConfigFile()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to load configuration file", "error", err)
			return
		}

		agentConfig, err := config.ResolveConfig()
		if err != nil {
			slog.ErrorContext(ctx, "Invalid config", "error", err)
			return
		}

		slog.InfoContext(ctx, "Starting NGINX Agent",
			slog.String("version", a.version),
			slog.String("commit", a.commit),
		)

		messagePipe := bus.NewMessagePipe(defaultMessagePipeChannelSize)
		err = messagePipe.Register(defaultQueueSize, plugin.LoadPlugins(ctx, agentConfig))
		if err != nil {
			slog.ErrorContext(ctx, "Failed to register plugins", "error", err)
			return
		}

		messagePipe.Run(ctx)
	})
	err := config.Execute(ctx)
	if err != nil {
		return err
	}

	return nil
}
