// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/sirupsen/logrus"
)

// Compress attempts to capture all builds from the database.
// The function will then iterate through the list of builds
// and capture all logs for each build. Then, the function
// will iterate through each log for the build and update
// the log entry with compressed data.
func (d *db) Compress() error {
	logrus.Debug("executing compress from provided configuration")

	logrus.Info("capturing all builds from the database")
	// capture all builds from the database
	builds, err := d.Client.GetBuildList()
	if err != nil {
		return err
	}

	// iterate through all builds from the database
	for _, build := range builds {
		logrus.Infof("capturing all logs for build %d", build.GetID())

		// TODO: remove this hack
		//
		// this allows us to "ignore" the error messages
		// returned from GetBuildLogs()
		//
		// capture current log level
		currentLevel := logrus.GetLevel()
		// only output panic level logs
		logrus.SetLevel(logrus.PanicLevel)

		// capture all logs for the build from the database
		logs, err := d.Client.GetBuildLogs(build.GetID())
		if err != nil {
			return err
		}

		// TODO: remove this hack
		//
		// this allows us to "ignore" the error messages
		// returned from GetBuildLogs()
		//
		// output intended level of logs
		logrus.SetLevel(currentLevel)

		// iterate through all logs for the build from the database
		for _, log := range logs {
			// update log entry with compression in the database
			err = d.Client.UpdateLog(log)
			if err != nil {
				return err
			}
		}

		logrus.Debugf("all logs updated for build %d", build.GetID())
	}

	return nil
}
