// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// Encrypt attempts to capture all users from the database. The
// function will then spawn the configured number of go routines
// to begin processing the list of users concurrently. For each
// user, the function will publish it to a channel all routines
// are listening on to encrypt the user fields.
func (d *db) Encrypt() error {
	logrus.Debug("executing encrypt from provided configuration")

	logrus.Info("capturing all users from the database")
	// capture all users from the database
	users, err := d.Client.GetUserList()
	if err != nil {
		return err
	}

	// create new error group to encrypt secret values concurrently
	group := new(errgroup.Group)
	// create new channel to process users concurrently
	channel := make(chan *library.User)

	// add set limit of routines to errgroup
	// and begin processing secrets
	for i := 0; i < d.ConcurrencyLimit; i++ {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := i

		logrus.Infof("spawning go routine %d to listen on user channel", tmp)

		// spawn a goroutine to begin encrypting user
		// fields that are published to the channel
		group.Go(func() error {
			return d.EncryptUsers(tmp, channel)
		})
	}

	// iterate through all users from the database
	for _, user := range users {
		logrus.Infof("publishing user %d to channel", user.GetID())

		// publish the user to the channel
		channel <- user

		logrus.Debugf("user %d published to channel", user.GetID())
	}

	logrus.Debug("closing channel for publishing users")

	// close channel to signal goroutines to stop processing
	close(channel)

	logrus.Debug("waiting for user go routines to complete")

	return group.Wait()
}

// EncryptUsers will iterate over all users published to the
// channel until the channel is closed. For each user published
// to the channel the function will update the user with an
// encrypted value.
func (d *db) EncryptUsers(index int, channel chan *library.User) error {
	logrus.Infof("go routine %d: listening on user channel", index)

	// iterate through all users published to the channel
	for u := range channel {
		logrus.Infof("go routine %d: encrypting the value for user %d", index, u.GetID())

		// update user with encryption in the database
		err := d.Client.UpdateUser(u)
		if err != nil {
			return err
		}

		logrus.Tracef("go routine %d: value encrypted for user %d", index, u.GetID())
	}

	logrus.Infof("go routine %d: shutting down on user channel", index)

	return nil
}
