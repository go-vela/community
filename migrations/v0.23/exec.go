// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Exec takes the provided configuration and attempts to
// run a series of different functions that will
// manipulate indexes, tables and columns.
func (d *db) Exec(c *cli.Context) error {
	logrus.Debug("executing utility from provided configuration")

	var err error

	switch d.Driver {
	case constants.DriverSqlite:
		// create the new Postgres database client
		//
		// https://pkg.go.dev/gorm.io/gorm#Open
		d.Gorm, err = gorm.Open(sqlite.Open(d.Address), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("unable to connect to %s database: %v", constants.DriverSqlite, err)
		}
	case constants.DriverPostgres:
		// create the new Postgres database client
		//
		// https://pkg.go.dev/gorm.io/gorm#Open
		d.Gorm, err = gorm.Open(postgres.Open(d.Address), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("unable to connect to %s database: %v", constants.DriverPostgres, err)
		}
	}

	defer func() { _sql, _ := d.Gorm.DB(); _sql.Close() }()

	// check if either the all or alter tables flag was set
	if d.Actions.All || d.Actions.AlterTables {
		// alter required tables in the database
		err = d.Alter()
		if err != nil {
			return err
		}
	}

	// check if either the all or update.repos flag was set
	if d.Actions.All || d.Actions.UpdateRepos {
		// update repos table for allow_events and approve_build
		err = d.UpdateRepos()
		if err != nil {
			return err
		}
	}

	// check if either the all or update.secrets flag was set
	if d.Actions.All || d.Actions.UpdateSecrets {
		// update secrets table for allow_events
		err = d.UpdateSecrets()
		if err != nil {
			return err
		}
	}

	return nil
}
