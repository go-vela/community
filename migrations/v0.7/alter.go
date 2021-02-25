// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	AlterBuilds = `
ALTER TABLE builds
ADD COLUMN IF NOT EXISTS
deploy_payload VARCHAR(2000);
`

	AlterUsers = `
ALTER TABLE users
ADD COLUMN IF NOT EXISTS
refresh_token VARCHAR(500);
`

	AlterWorkers = `
ALTER TABLE workers
ADD COLUMN IF NOT EXISTS
build_limit INTEGER;
`
)

// Alter will run a series of ALTER statements against the
// database to modify the tables that require new or
// modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	// create a fresh database client
	//
	// This will allow us to manually run SQL statements
	// against the database.
	_database, err := gorm.Open(d.Driver, d.Config)
	if err != nil {
		return err
	}
	defer _database.Close()

	logrus.Infof("altering %s table to add deploy_payload column", constants.TableBuild)
	// alter the builds table to add the deploy_payload column
	_, err = _database.DB().Exec(AlterBuilds)
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableBuild, err)
	}

	logrus.Infof("altering %s table to add refresh_token column", constants.TableUser)
	// alter the users table to add the refresh_token column
	_, err = _database.DB().Exec(AlterUsers)
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableUser, err)
	}

	logrus.Infof("altering %s table to add build_limit column", constants.TableWorker)
	// alter the workers table to add the build_limit column
	_, err = _database.DB().Exec(AlterWorkers)
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableWorker, err)
	}

	return nil
}
