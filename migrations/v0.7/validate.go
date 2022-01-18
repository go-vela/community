// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
	case d.Actions.CompressLogs:
		fallthrough
	case d.Actions.DropIndexes:
		fallthrough
	case d.Actions.EncryptSecrets:
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

	// check if the database configuration is set
	if len(d.Config) == 0 {
		return fmt.Errorf("VELA_DATABASE_CONFIG is not properly configured")
	}

	// check if the database concurrency limit is set
	if d.ConcurrencyLimit < 1 {
		return fmt.Errorf("VELA_CONCURRENCY_LIMIT is not properly configured")
	}

	// check if either the all or compress logs action was provided
	if d.Actions.All || d.Actions.CompressLogs {
		// check if the database build limit is set
		if d.BuildLimit < 0 {
			return fmt.Errorf("VELA_BUILD_LIMIT is not properly configured")
		}

		// check if the compression level is valid
		switch d.CompressionLevel {
		case constants.CompressionNegOne:
			fallthrough
		case constants.CompressionZero:
			fallthrough
		case constants.CompressionOne:
			fallthrough
		case constants.CompressionTwo:
			fallthrough
		case constants.CompressionThree:
			fallthrough
		case constants.CompressionFour:
			fallthrough
		case constants.CompressionFive:
			fallthrough
		case constants.CompressionSix:
			fallthrough
		case constants.CompressionSeven:
			fallthrough
		case constants.CompressionEight:
			fallthrough
		case constants.CompressionNine:
			break
		default:
			return fmt.Errorf("invalid level for VELA_DATABASE_COMPRESSION_LEVEL provided: %d", d.CompressionLevel)
		}
	}

	// check if either the all or encrypt secrets action was provided
	if d.Actions.All || d.Actions.EncryptSecrets {
		// enforce AES-256, so check explicitly for 32 bytes on the key
		//
		// nolint: gomnd // ignore magic number
		if len(d.EncryptionKey) != 32 {
			return fmt.Errorf("invalid length for VELA_DATABASE_ENCRYPTION_KEY provided: %d", len(d.EncryptionKey))
		}

		// check if the database secret limit is set
		if d.SecretLimit < 0 {
			return fmt.Errorf("VELA_SECRET_LIMIT is not properly configured")
		}
	}

	return nil
}
