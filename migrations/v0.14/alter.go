// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

const AlterBuildsTableAddEventAction = `
ALTER TABLE builds
ADD COLUMN IF NOT EXISTS
event_action VARCHAR(250);
`

const AlterBuildsTableAddPipelineID = `
ALTER TABLE builds
ADD COLUMN IF NOT EXISTS
pipeline_id INTEGER;
`

const AlterHooksTableAddEventAction = `
ALTER TABLE hooks
ADD COLUMN IF NOT EXISTS
event_action VARCHAR(250);
`

const AlterHooksTableAddWebhookID = `
ALTER TABLE hooks
ADD COLUMN IF NOT EXISTS
webhook_id INTEGER;
`

// Alter will run a series of ALTER statements against the database
// to modify the tables that require new or modified columns.
func (d *db) Alter() error {
	logrus.Debug("executing alter from provided configuration")

	logrus.Infof("altering %s table to add event_action column", constants.TableBuild)
	// alter builds table to add event_action column
	err := d.Gorm.Exec(AlterBuildsTableAddEventAction).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableBuild, err)
	}

	logrus.Infof("altering %s table to add pipeline_id column", constants.TableBuild)
	// alter builds table to add pipeline_id column
	err = d.Gorm.Exec(AlterBuildsTableAddPipelineID).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableBuild, err)
	}

	logrus.Infof("altering %s table to add webhook_id column", constants.TableHook)
	// alter hooks table to add webhook_id column
	err = d.Gorm.Exec(AlterHooksTableAddWebhookID).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableHook, err)
	}

	logrus.Infof("altering %s table to add event_action column", constants.TableHook)
	// alter hooks table to add event_action column
	err = d.Gorm.Exec(AlterHooksTableAddEventAction).Error
	if err != nil {
		return fmt.Errorf("unable to alter %s table: %v", constants.TableHook, err)
	}

	return nil
}
