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

	return nil
}
