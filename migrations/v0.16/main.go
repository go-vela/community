// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
	app.Usage = "Vela utility used to migrate from v0.15.x to v0.16.x"
	app.Copyright = "Copyright (c) 2022 Target Brands, Inc. All rights reserved."
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
			Usage:   "enables running all actions for v0.16.x",
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_ACTION_SET_UNTRUSTED", "ACTION_SET_UNTRUSTED"},
			Name:    "action.untrusted",
			Usage:   "enables altering repos to untrusted for v0.16.x",
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_ACTION_UPDATE_PRIV_REPOS", "ACTION_UPDATE_PRIV_REPOS"},
			Name:    "action.update-trusted",
			Usage:   "enables altering repos to trusted using privileged images for v0.16.x",
		},

		// Trusted Update Flags

		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_RUNTIME_PRIVILEGED_IMAGES", "RUNTIME_PRIVILEGED_IMAGES"},
			Name:    "trusted-update.privileged-images",
			Usage:   "provide the privileged images to use to grant trusted status to repos that use them",
			Value:   cli.NewStringSlice("target/vela-docker"),
		},
		&cli.BoolFlag{
			EnvVars: []string{"VELA_ALLOW_PERSONAL_ORGS", "ALLOW_PERSONAL_ORGS"},
			Name:    "trusted-update.allow-personal-orgs",
			Usage:   "provide whether or not repos from personal orgs can be granted trusted status",
			Value:   true,
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_DAYS_BACK", "DAYS_BACK"},
			Name:    "trusted-update.days-back",
			Usage:   "provide how many days back to read for repos using privileged images",
			Value:   "90",
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
