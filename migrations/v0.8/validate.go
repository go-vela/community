// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
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
	case d.Actions.EncryptUsers:
		fallthrough
	case d.Actions.SyncCounter:
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

	// check if the database concurrency limit is set
	if d.ConcurrencyLimit < 1 {
		return fmt.Errorf("VELA_CONCURRENCY_LIMIT is not properly configured")
	}

	// check if either the all or encrypt users action was provided
	if d.Actions.All || d.Actions.EncryptUsers {
		// enforce AES-256, so check explicitly for 32 bytes on the key
		//
		// nolint: gomnd // ignore magic number
		if len(d.EncryptionKey) != 32 {
			return fmt.Errorf("invalid length for VELA_DATABASE_ENCRYPTION_KEY provided: %d", len(d.EncryptionKey))
		}
	}

	return nil
}
