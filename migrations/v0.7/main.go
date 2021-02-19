// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Utility Information

	app.Name = "vela-migration"
	app.HelpName = "vela-migration"
	app.Usage = "Vela utility used to migrate from v0.6.x to v0.7.x"
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

		// Database Flags

		&cli.StringFlag{
			EnvVars: []string{"VELA_DATABASE_DRIVER", "DATABASE_DRIVER"},
			Name:    "database.driver",
			Usage:   "sets the driver to be used for the database",
			Value:   "sqlite3",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_DATABASE_CONFIG", "DATABASE_CONFIG"},
			Name:    "database.config",
			Usage:   "sets the configuration string to be used for the database",
			Value:   "vela.sqlite",
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_OPEN", "DATABASE_CONNECTION_OPEN"},
			Name:    "database.connection.open",
			Usage:   "sets the number of open connections to the database",
			Value:   0,
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_IDLE", "DATABASE_CONNECTION_IDLE"},
			Name:    "database.connection.idle",
			Usage:   "sets the number of idle connections to the database",
			Value:   2,
		},
		&cli.DurationFlag{
			EnvVars: []string{"VELA_DATABASE_CONNECTION_LIFE", "DATABASE_CONNECTION_LIFE"},
			Name:    "database.connection.life",
			Usage:   "sets the amount of time a connection may be reused for the database",
			Value:   30 * time.Minute,
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

	// Utility Start

	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
