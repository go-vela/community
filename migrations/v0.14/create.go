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
	CreateIndexBuildsSource = `
CREATE INDEX CONCURRENTLY
IF NOT EXISTS
builds_source ON builds (source);
`
)

// Create will run a series of CREATE INDEX statements against
// the database to add new indexes for existing columns.
func (d *db) Create() error {
	logrus.Debug("executing create from provided configuration")

	logrus.Infof("creating index for %s table on created column", constants.TableBuild)
	// create an index on the source column for the builds table
	err := d.Gorm.Exec(CreateIndexBuildsSource).Error
	if err != nil {
		return fmt.Errorf("unable to create index for %s table: %v", constants.TableBuild, err)
	}

	return nil
}
