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
	UpdateSchedules = `
  UPDATE schedules 
      SET branch = r.branch 
      FROM (SELECT id, branch FROM repos) r 
      WHERE schedules.repo_id = r.id
	  AND schedules.branch IS NULL
      ; 
`
)

// Update will run an UPDATE statement against the
// database to modify the tables that require new or
// modified data for the migration.
func (d *db) Update() error {
	logrus.Info("updating existing schedules to use repo default branch")

	err := d.Gorm.Exec(UpdateSchedules).Error
	if err != nil {
		return fmt.Errorf("unable to update schedules: %w", err)
	}

	logrus.Infof("%s table successfully updated", constants.TableSchedule)

	return nil
}
