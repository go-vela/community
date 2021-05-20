// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// Sync attempts to capture all repos from the database. The
// function will then spawn the configured number of go routines
// to begin processing the list of repos concurrently. For each
// repo, the function will publish it to a channel all routines
// are listening on to sync the repo counter value.
func (d *db) Sync() error {
	logrus.Debug("executing sync counter from provided configuration")

	logrus.Info("capturing all repos from the database")
	// capture all repos from the database
	repos, err := d.Client.GetRepoList()
	if err != nil {
		return err
	}

	// create new error group to encrypt repo values concurrently
	group := new(errgroup.Group)
	// create new channel to process repos concurrently
	channel := make(chan *library.Repo)

	// add set limit of routines to errgroup
	// and begin processing repos
	for i := 0; i < d.ConcurrencyLimit; i++ {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := i

		logrus.Infof("spawning go routine %d to listen on repo channel", tmp)

		// spawn a goroutine to begin syncing repo
		// counter values that are published to the channel
		group.Go(func() error {
			return d.SyncCounter(tmp, channel)
		})
	}

	// iterate through all repos from the database
	for _, repo := range repos {
		logrus.Infof("publishing repo %d to channel", repo.GetID())

		// publish the repo to the channel
		channel <- repo

		logrus.Debugf("repo %d published to channel", repo.GetID())
	}

	logrus.Debug("closing channel for publishing repos")

	// close channel to signal goroutines to stop processing
	close(channel)

	logrus.Debug("waiting for repo go routines to complete")

	return group.Wait()
}

// SyncCounter will iterate over all repos published to the
// channel until the channel is closed. For each repo published
// to the channel the function will get the latest build number
// and update the repo counter to that value.
func (d *db) SyncCounter(index int, channel chan *library.Repo) error {
	logrus.Infof("go routine %d: listening on repo channel", index)

	// iterate through all repos published to the channel
	for r := range channel {
		logrus.Infof("go routine %d: syncing counter the value for repo %d", index, r.GetID())

		// get the latest build ran
		lastBuild, err := d.Client.GetLastBuild(r)
		if err != nil {
			return err
		}

		r.SetCounter(lastBuild.GetNumber())

		// update repo with latest build number
		err = d.Client.UpdateRepo(r)
		if err != nil {
			return err
		}

		logrus.Tracef("go routine %d: value encrypted for repo %d", index, r.GetID())
	}

	logrus.Infof("go routine %d: shutting down on repo channel", index)

	return nil
}
