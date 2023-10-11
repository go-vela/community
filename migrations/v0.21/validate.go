// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// Validate verifies the provided database is configured properly.
func (d *db) Validate() error {
	logrus.Debug("validating provided database configuration")

	// check if an action was provided
	switch {
	case d.Actions.AlterTables:
		fallthrough
	case d.Actions.UpdateSchedules:
		if !d.Actions.AlterTables {
			logrus.Warning("cannot execute Update Schedules action without altering schedules table to include branch field")
		}
		fallthrough
	case d.Actions.All:
		break
	default:
		logrus.Warning("no vela-migration actions provided")
	}

	// check if the database driver is set
	if len(d.Driver) == 0 {
		return fmt.Errorf("VELA_DATABASE_DRIVER is not properly configured")
	}

	switch d.Driver {
	case constants.DriverPostgres:
		fallthrough
	case constants.DriverSqlite:
		break
	default:
		return fmt.Errorf("invalid VELA_DATABASE_DRIVER provided: %s", d.Driver)
	}

	// check if the database address is set
	if len(d.Address) == 0 {
		return fmt.Errorf("VELA_DATABASE_ADDR is not properly configured")
	}

	return nil
}
