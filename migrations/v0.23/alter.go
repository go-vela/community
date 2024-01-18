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
	AlterReposApproveBuild = `
  ALTER TABLE repos 
      ADD COLUMN IF NOT EXISTS approve_build VARCHAR(250)
      ;
`
	AlterReposAllowEvents = `
	ALTER TABLE repos
	  ADD COLUMN IF NOT EXISTS allow_events INTEGER
	  ;
`

	AlterSecretsAllowEvents = `
	ALTER TABLE secrets
	  ADD COLUMN IF NOT EXISTS allow_events INTEGER
	  ;
`

	AlterBuildsApprovedAt = `
	ALTER TABLE builds
	  ADD COLUMN IF NOT EXISTS approved_at INTEGER
	  ;
`

	AlterBuildsApprovedBy = `
	ALTER TABLE builds
	  ADD COLUMN IF NOT EXISTS approved_by INTEGER
	  ;
`
)

// Alter will run a series of ALTER statements against the
// database to modify the tables that require new or
// modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	logrus.Infof("altering %s table to add branch column", constants.TableRepo)
	// alter repos table to add branch column
	err := d.Gorm.Exec(AlterReposApproveBuild).Error
	if err != nil {
		return fmt.Errorf("unable to add approve_build column to %s table: %v", constants.TableRepo, err)
	}

	err = d.Gorm.Exec(AlterReposAllowEvents).Error
	if err != nil {
		return fmt.Errorf("unable to add allow_events column to %s table: %v", constants.TableRepo, err)
	}

	err = d.Gorm.Exec(AlterSecretsAllowEvents).Error
	if err != nil {
		return fmt.Errorf("unable to add allow_events column to %s table: %v", constants.TableSecret, err)
	}

	err = d.Gorm.Exec(AlterBuildsApprovedAt).Error
	if err != nil {
		return fmt.Errorf("unable to add approved_at column to %s table: %v", constants.TableBuild, err)
	}

	err = d.Gorm.Exec(AlterBuildsApprovedBy).Error
	if err != nil {
		return fmt.Errorf("unable to add approved_by column to %s table: %v", constants.TableBuild, err)
	}

	logrus.Infof("%s and %s tables successfully altered", constants.TableRepo, constants.TableBuild)

	return nil
}
