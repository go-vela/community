// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
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

	// create new wait group to compress build logs concurrently
	group := new(errgroup.Group)
	// create new channel to process build logs concurrently
	buildChannel := make(chan *library.Build)

	// add set limit of routines to errgroup
	// and begin processing build logs
	for i := 0; i < d.ConcurrencyLimit; i++ {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := i

		// spawn a goroutine to begin compressing build
		// logs that are published to the channel
		group.Go(func() error {
			return d.CompressBuildLogs(tmp, buildChannel)
		})
	}

	// iterate through all builds from the database
	for _, build := range builds {
		// handle the build based off the id provided
		if d.BuildLimit > 0 && build.GetID() > int64(d.BuildLimit) {
			logrus.Tracef("build %d is greater than threshold %d - skipping", build.GetID(), d.BuildLimit)

			continue
		}

		// handle the build based off the status provided
		switch build.GetStatus() {
		// build is in a pending state
		case constants.StatusPending:
			fallthrough
		// build is in a running state
		case constants.StatusRunning:
			logrus.Tracef("build %d has a status of %s - skipping", build.GetID(), build.GetStatus())

			continue
		}

		logrus.Infof("publishing build %d to channel", build.GetID())

		// publish the build to the channel
		buildChannel <- build

		logrus.Debugf("build %d published to channel", build.GetID())
	}

	logrus.Debug("closing channel for publishing builds")

	// close channel to signal goroutines to stop processing
	close(buildChannel)

	logrus.Debug("waiting for goroutines to complete")

	return group.Wait()
}

func (d *db) CompressBuildLogs(index int, buildChannel chan *library.Build) error {
	logrus.Infof("thread %d: listening on build channel", index)

	// iterate through all builds published to the channel
	for b := range buildChannel {
		logrus.Infof("thread %d: capturing all logs for build %d", index, b.GetID())

		// capture all logs for the build from the database
		logs, err := d.Client.GetBuildLogs(b.GetID())
		if err != nil {
			return err
		}

		logrus.Infof("thread %d: compressing all logs for build %d", index, b.GetID())

		// iterate through all logs for the build from the database
		for _, log := range logs {
			// update log entry with compression in the database
			err = d.Client.UpdateLog(log)
			if err != nil {
				return err
			}
		}

		logrus.Debugf("thread %d: all logs compressed for build %d", index, b.GetID())
	}

	logrus.Infof("thread %d: shutting down", index)

	return nil
}
