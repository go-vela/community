// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/go-vela/server/database"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Utility Information

	app.Name = "vela-migration"
	app.HelpName = "vela-migration"
	app.Usage = "Vela utility used to migrate from v0.7.x to v0.8.x"
	app.Copyright = "Copyright (c) 2021 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Utility Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = "v1.0.0"

	// Utility Flags

	app.Flags = []cli.Flag{

		// Action Flags

		&cli.BoolFlag{
			EnvVars: []string{"VELA_ACTION_ALL", "ACTION_ALL"},
			Name:    "action.all",
			Usage:   "enables running all actions for v0.8.x",
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_ALTER_TABLES", "ALTER_TABLES"},
			Name:    "alter.tables",
			Usage:   "enables altering the table configuration for v0.8.x",
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_ENCRYPT_USERS", "ENCRYPT_USERS"},
			Name:    "encrypt.users",
			Usage:   "enables encrypting user fields for v0.8.x",
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_SYNC_COUNTER", "SYNC_COUNTER"},
			Name:    "sync.counter",
			Usage:   "enables encrypting secret values for v0.8.x",
		},

		// Limit Flags

		&cli.IntFlag{
			EnvVars: []string{"VELA_CONCURRENCY_LIMIT", "CONCURRENCY_LIMIT"},
			Name:    "concurrency.limit",
			Usage:   "sets the number of concurrent processes running",
			Value:   4,
		},

		// Logger Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_LOG_FORMAT", "LOG_FORMAT"},
			Name:    "log.format",
			Usage:   "set log format for the utility",
			Value:   "json",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_LOG_LEVEL", "LOG_LEVEL"},
			Name:    "log.level",
			Usage:   "set log level for the utility",
			Value:   "info",
		},
	}

	// Database Flags

	app.Flags = append(app.Flags, database.Flags...)

	// Utility Start

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
