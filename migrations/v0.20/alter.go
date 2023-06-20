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
	AlterWorkers = `
  ALTER TABLE workers 
      ADD COLUMN IF NOT EXISTS status VARCHAR(50),
      ADD COLUMN IF NOT EXISTS last_status_update_at INTEGER,
      ADD COLUMN IF NOT EXISTS running_build_ids VARCHAR(500),
      ADD COLUMN IF NOT EXISTS last_build_started_at INTEGER,
      ADD COLUMN IF NOT EXISTS last_build_finished_at INTEGER
      ;
`
)

// Alter will run a series of ALTER statements against the
// database to modify the tables that require new or
// modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	logrus.Infof("altering %s table to add visibility columns", constants.TableWorker)
	// alter builds table to add event_action column
	err := d.Gorm.Exec(AlterWorkers).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableWorker, err)
	}

	return nil
}
