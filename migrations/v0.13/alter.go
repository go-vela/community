// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

const AlterBuildsErrorSize = `
ALTER TABLE builds
ALTER COLUMN error
TYPE varchar(1000);
`

// Alter will run a series of ALTER statements against the database
// to modify the tables that require new or modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	logrus.Infof("altering %s table to add error column", constants.TableBuild)
	// alter the builds table to modify the size of the error column
	// note: previous size was 500
	err := d.Gorm.Exec(AlterBuildsErrorSize).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableBuild, err)
	}

	return nil
}
