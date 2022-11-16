// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// run executes the utility based
// off the configuration provided.
func run(c *cli.Context) error {
	// set log format for the utility
	switch c.String("log.format") {
	case "t", "text", "Text", "TEXT":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "j", "json", "Json", "JSON":
		fallthrough
	default:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// set log level for the utility
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// create database object
	d := &db{
		Driver:  c.String("database.driver"),
		Address: c.String("database.addr"),
		Actions: &actions{
			All:             c.Bool("action.all"),
			SetUntrusted:    c.Bool("action.untrusted"),
			UpdatePrivRepos: c.Bool("action.update-trusted"),
		},
		TrustedOptions: &trustedOptions{
			Images:       c.StringSlice("trusted-update.privileged-images"),
			PersonalOrgs: c.Bool("trusted-update.allow-personal-orgs"),
			Days:         c.String("trusted-update.days-back"),
		},
	}

	// validate database configuration
	err := d.Validate()
	if err != nil {
		return err
	}

	return d.Exec(c)
}
