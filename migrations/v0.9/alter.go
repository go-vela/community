// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
	AlterReposPipelineType = `
ALTER TABLE repos
ADD COLUMN IF NOT EXISTS
pipeline_type TEXT
DEFAULT 'yaml';
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
	_database, err := gorm.Open(d.Driver, d.Address)
	if err != nil {
		return err
	}
	defer _database.Close()

	logrus.Infof("altering %s table to add pipeline_type column", constants.TableRepo)
	// alter the repo table to add the pipeline_type column
	_, err = _database.DB().Exec(AlterReposPipelineType)
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableRepo, err)
	}

	return nil
}
