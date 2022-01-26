// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

const (
	AlterReposBuildLimit = `
ALTER TABLE repos
ADD COLUMN IF NOT EXISTS
build_limit INTEGER
DEFAULT 10;
`

	AlterReposPreviousName = `
ALTER TABLE repos
ADD COLUMN IF NOT EXISTS
previous_name VARCHAR(100);
`

	AlterSecretsCreatedAt = `
ALTER TABLE secrets
ADD COLUMN IF NOT EXISTS
created_at INTEGER;
`

	AlterSecretsCreatedBy = `
ALTER TABLE secrets
ADD COLUMN IF NOT EXISTS
created_by VARCHAR(250);
`

	AlterSecretsUpdatedAt = `
ALTER TABLE secrets
ADD COLUMN IF NOT EXISTS
updated_at INTEGER;
`

	AlterSecretsUpdatedBy = `
ALTER TABLE secrets
ADD COLUMN IF NOT EXISTS
updated_by VARCHAR(250);
`
)

// Alter will run a series of ALTER statements against the database
// to modify the tables that require new or modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	logrus.Infof("altering %s table to add build_limit column", constants.TableRepo)
	// alter the repos table to add the build_limit column
	err := d.Gorm.Exec(AlterReposBuildLimit).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableRepo, err)
	}

	logrus.Infof("altering %s table to add previous_name column", constants.TableRepo)
	// alter the repos table to add the previous_name column
	err = d.Gorm.Exec(AlterReposPreviousName).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableRepo, err)
	}

	logrus.Infof("altering %s table to add created_at column", constants.TableSecret)
	// alter the secrets table to add the created_at column
	err = d.Gorm.Exec(AlterSecretsCreatedAt).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableSecret, err)
	}

	logrus.Infof("altering %s table to add created_by column", constants.TableSecret)
	// alter the secrets table to add the created_by column
	err = d.Gorm.Exec(AlterSecretsCreatedBy).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableSecret, err)
	}

	logrus.Infof("altering %s table to add updated_at column", constants.TableSecret)
	// alter the secrets table to add the updated_at column
	err = d.Gorm.Exec(AlterSecretsUpdatedAt).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableSecret, err)
	}

	logrus.Infof("altering %s table to add updated_by column", constants.TableSecret)
	// alter the secrets table to add the updated_by column
	err = d.Gorm.Exec(AlterSecretsUpdatedBy).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableSecret, err)
	}

	return nil
}
