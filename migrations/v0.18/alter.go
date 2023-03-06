// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

const AlterUsersTokenSize = `
ALTER TABLE users
ALTER COLUMN token
TYPE varchar(1000);
`

// Alter will run a series of ALTER statements against the database
// to modify the tables that require new or modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	logrus.Infof("altering %s table to increase size of token column", constants.TableUser)
	// alter the users table to modify the size of the token column
	// note: previous size was 500
	err := d.Gorm.Exec(AlterUsersTokenSize).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableUser, err)
	}

	logrus.Info("successfully altered table")

	return nil
}
