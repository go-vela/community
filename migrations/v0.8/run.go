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
		Driver:           c.String("database.driver"),
		Address:          c.String("database.addr"),
		ConcurrencyLimit: c.Int("concurrency.limit"),
		EncryptionKey:    c.String("database.encryption.key"),
		Actions: &actions{
			All:          c.Bool("action.all"),
			AlterTables:  c.Bool("alter.tables"),
			EncryptUsers: c.Bool("encrypt.users"),
			SyncCounter:  c.Bool("sync.counter"),
		},
		Connection: &connection{
			Idle: c.Int("database.connection.open"),
			Life: c.Duration("database.connection.idle"),
			Open: c.Int("database.connection.life"),
		},
	}

	// validate database configuration
	err := d.Validate()
	if err != nil {
		return err
	}

	// create new database service
	err = d.New(c)
	if err != nil {
		return err
	}

	return d.Exec(c)
}
