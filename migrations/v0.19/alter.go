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
	AlterReposTopics = `
ALTER TABLE repos
ADD COLUMN IF NOT EXISTS
topics VARCHAR(1020);
`
)

// Alter will run a series of ALTER statements against the
// database to modify the tables that require new or
// modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	logrus.Infof("altering %s table to add topics column", constants.TableRepo)
	// alter builds table to add event_action column
	err := d.Gorm.Exec(AlterReposTopics).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableRepo, err)
	}

	return nil
}
