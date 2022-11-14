// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

const UpdateReposTableSetUntrusted = `
UPDATE repos SET trusted = 'false';
`

// UpdateReposUntrusted will change all repos in the database
// to have a trusted value of 'false'.
func (d *db) UpdateReposUntrusted() error {
	logrus.Info("updating all repos to untrusted")
	// set all repos to untrusted
	err := d.Gorm.Exec(UpdateReposTableSetUntrusted).Error
	if err != nil {
		return fmt.Errorf("unable to update %s table: %v", constants.TableRepo, err)
	}
	logrus.Info("successfully updated all repos to be not trusted")

	return nil
}
