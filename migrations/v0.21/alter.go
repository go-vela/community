// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

const (
	AlterSchedules = `
  ALTER TABLE schedules 
      ADD COLUMN IF NOT EXISTS branch VARCHAR(250)
      ;
`
)

// Alter will run a series of ALTER statements against the
// database to modify the tables that require new or
// modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	logrus.Infof("altering %s table to add branch column", constants.TableSchedule)
	// alter schedules table to add branch column
	err := d.Gorm.Exec(AlterSchedules).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableSchedule, err)
	}

	logrus.Infof("%s table successfully altered", constants.TableSchedule)

	return nil
}
