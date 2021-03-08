// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Exec takes the provided configuration and attempts to
// run a series of different functions that will
// manipulate indexes, tables and columns.
func (d *db) Exec(c *cli.Context) error {
	logrus.Debug("executing workload from provided configuration")

	// check if either the all or alter tables action was provided
	if d.Actions.All || d.Actions.AlterTables {
		// alter required tables in the database
		err := d.Alter()
		if err != nil {
			return err
		}
	}

	// check if either the all or drop indexes action was provided
	if d.Actions.All || d.Actions.DropIndexes {
		// drop unused indexes in the database
		err := d.Drop()
		if err != nil {
			return err
		}
	}

	// create new database service
	err := d.New(c)
	if err != nil {
		return err
	}

	// check if either the all or compress logs action was provided
	if d.Actions.All || d.Actions.CompressLogs {
		// compress all log entries in the database
		err = d.Compress()
		if err != nil {
			return err
		}
	}

	// check if either the all or encrypt secrets action was provided
	if d.Actions.All || d.Actions.EncryptSecrets {
		// encrypt all secret values in the database
		err = d.Encrypt()
		if err != nil {
			return err
		}
	}

	return nil
}
