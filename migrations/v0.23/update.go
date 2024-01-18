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
	UpdateReposApproveBuild = `
  UPDATE repos 
      SET approve_build = 'fork-always'
	  ;
`

	UpdateReposAllowEvents = `
UPDATE repos
    SET allow_events =
        (CASE WHEN allow_push = true THEN 1 ELSE 0 END) |
        (CASE WHEN allow_pull = true THEN 4 | 16 | 1024 ELSE 0 END) |
        (CASE WHEN allow_tag = true THEN 2 ELSE 0 END) |
        (CASE WHEN allow_deploy = true THEN 8192 ELSE 0 END) |
        (CASE WHEN allow_comment = true THEN 16384 | 32768 ELSE 0 END)
;
`
	UpdateSecretsAllowEvents = `
UPDATE secrets
	SET allow_events =
	    (CASE WHEN events LIKE '%push%' THEN 1 ELSE 0 END) |
        (CASE WHEN events LIKE '%pull_request%' THEN 4 | 16 | 1024 ELSE 0 END) |
        (CASE WHEN events LIKE '%tag%' THEN 2 ELSE 0 END) |
        (CASE WHEN events LIKE '%deployment%' THEN 8192 ELSE 0 END) |
        (CASE WHEN events LIKE '%comment%' THEN 16384 | 32768 ELSE 0 END) |
		(CASE WHEN events LIKE '%schedule%' THEN 65536 ELSE 0 END)
;
`
)

// UpdateRepos will run an UPDATE statement against the
// database repos table to modify the data for the migration.
func (d *db) UpdateRepos() error {
	logrus.Info("updating existing repos to a default approve_build policy of 'fork-always'")

	err := d.Gorm.Exec(UpdateReposApproveBuild).Error
	if err != nil {
		return fmt.Errorf("unable to update repos approve_build: %w", err)
	}

	logrus.Infof("%s table approve_build successfully updated", constants.TableRepo)

	logrus.Info("migrating existing allow fields to allow_events")

	err = d.Gorm.Exec(UpdateReposAllowEvents).Error
	if err != nil {
		return fmt.Errorf("unable to update repos allow_events: %w", err)
	}

	logrus.Infof("%s table allow_events column successfully populated", constants.TableRepo)

	logrus.Info("migrating existing secrets events list to allow_events")

	err = d.Gorm.Exec(UpdateSecretsAllowEvents).Error
	if err != nil {
		return fmt.Errorf("unable to update secrets allow_events: %w", err)
	}

	logrus.Infof("%s table allow_events column successfully populated", constants.TableSecret)

	return nil
}

// UpdateSecrets will run an UPDATE statement against the
// database secrets table to modify the data for the migration.
func (d *db) UpdateSecrets() error {
	logrus.Info("migrating existing secrets events list to allow_events")

	err := d.Gorm.Exec(UpdateSecretsAllowEvents).Error
	if err != nil {
		return fmt.Errorf("unable to update secrets allow_events: %w", err)
	}

	logrus.Infof("%s table allow_events column successfully populated", constants.TableSecret)

	return nil
}
