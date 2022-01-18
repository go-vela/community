// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	DropBuildsIndex = `
DROP INDEX
IF EXISTS
builds_repo_id_number;
`

	DropHooksIndex = `
DROP INDEX
IF EXISTS
hooks_repo_id_number;
`

	DropLogsStepIndex = `
DROP INDEX
IF EXISTS
logs_step_id;
`

	DropLogsServiceIndex = `
DROP INDEX
IF EXISTS
logs_service_id;
`

	DropReposIndex = `
DROP INDEX
IF EXISTS
repos_full_name;
`

	DropSecretsIndex = `
DROP INDEX
IF EXISTS
secrets_type;
`

	DropServicesIndex = `
DROP INDEX
IF EXISTS
services_build_id_number;
`

	DropStepsIndex = `
DROP INDEX
IF EXISTS
steps_build_id_number;
`

	DropUsersIndex = `
DROP INDEX
IF EXISTS
users_name;
`
)

// Drop will run a series of DROP statements against the
// database to remove unused indexes.
func (d *db) Drop() error {
	logrus.Debug("executing drop from provided configuration")

	// create a fresh database client
	//
	// This will allow us to manually run SQL statements
	// against the database.
	_database, err := gorm.Open(d.Driver, d.Config)
	if err != nil {
		return err
	}
	defer _database.Close()

	logrus.Info("dropping builds_repo_id_number index")
	// drop the builds_repo_id_number index
	_, err = _database.DB().Exec(DropBuildsIndex)
	if err != nil {
		return fmt.Errorf("unable to drop builds_repo_id_number index: %v", err)
	}

	logrus.Info("dropping hooks_repo_id_number index")
	// drop the hooks_repo_id_number index
	_, err = _database.DB().Exec(DropHooksIndex)
	if err != nil {
		return fmt.Errorf("unable to drop hooks_repo_id_number index: %v", err)
	}

	logrus.Info("dropping logs_step_id index")
	// drop the logs_step_id index
	_, err = _database.DB().Exec(DropLogsStepIndex)
	if err != nil {
		return fmt.Errorf("unable to drop logs_step_id index: %v", err)
	}

	logrus.Info("dropping logs_service_id index")
	// drop the logs_service_id index
	_, err = _database.DB().Exec(DropLogsServiceIndex)
	if err != nil {
		return fmt.Errorf("unable to drop logs_service_id index: %v", err)
	}

	logrus.Info("dropping repos_full_name index")
	// drop the repos_full_name index
	_, err = _database.DB().Exec(DropReposIndex)
	if err != nil {
		return fmt.Errorf("unable to drop repos_full_name index: %v", err)
	}

	logrus.Info("dropping secrets_type index")
	// drop the secrets_type index
	_, err = _database.DB().Exec(DropSecretsIndex)
	if err != nil {
		return fmt.Errorf("unable to drop secrets_type index: %v", err)
	}

	logrus.Info("dropping services_build_id_number index")
	// drop the services_build_id_number index
	_, err = _database.DB().Exec(DropServicesIndex)
	if err != nil {
		return fmt.Errorf("unable to drop services_build_id_number index: %v", err)
	}

	logrus.Info("dropping steps_build_id_number index")
	// drop the steps_build_id_number index
	_, err = _database.DB().Exec(DropStepsIndex)
	if err != nil {
		return fmt.Errorf("unable to drop steps_build_id_number index: %v", err)
	}

	logrus.Info("dropping users_name index")
	// drop the users_name index
	_, err = _database.DB().Exec(DropUsersIndex)
	if err != nil {
		return fmt.Errorf("unable to drop users_name index: %v", err)
	}

	return nil
}
