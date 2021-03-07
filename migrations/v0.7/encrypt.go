// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// Encrypt attempts to capture all secrets from the database. The
// function will then spawn the configured number of go routines
// to begin processing the list of secrets concurrently. For each
// secret, the function will publish it to a channel all routines
// are listening on to encrypt the secret value.
func (d *db) Encrypt() error {
	logrus.Debug("executing encrypt from provided configuration")

	logrus.Info("capturing all secrets from the database")
	// capture all secrets from the database
	secrets, err := d.Client.GetSecretList()
	if err != nil {
		return err
	}

	// create new error group to encrypt secret values concurrently
	group := new(errgroup.Group)
	// create new channel to process secrets concurrently
	channel := make(chan *library.Secret)

	// add set limit of routines to errgroup
	// and begin processing secrets
	for i := 0; i < d.ConcurrencyLimit; i++ {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := i

		logrus.Infof("spawning go routine %d to listen on secret channel", tmp)

		// spawn a goroutine to begin encrypting secret
		// values that are published to the channel
		group.Go(func() error {
			return d.EncryptSecrets(tmp, channel)
		})
	}

	// iterate through all secrets from the database
	for _, secret := range secrets {
		// handle the secret based off the id provided
		if d.SecretLimit > 0 && secret.GetID() > int64(d.SecretLimit) {
			logrus.Tracef("secret %d is greater than limit %d - skipping", secret.GetID(), d.SecretLimit)

			continue
		}

		logrus.Infof("publishing secret %d to channel", secret.GetID())

		// publish the secret to the channel
		channel <- secret

		logrus.Debugf("secret %d published to channel", secret.GetID())
	}

	logrus.Debug("closing channel for publishing secrets")

	// close channel to signal goroutines to stop processing
	close(channel)

	logrus.Debug("waiting for secret go routines to complete")

	return group.Wait()
}

// EncryptSecrets will iterate over all secrets published to the
// channel until the channel is closed. For each secret published
// to the channel the function will update the secret with an
// encrypted value.
func (d *db) EncryptSecrets(index int, channel chan *library.Secret) error {
	logrus.Infof("go routine %d: listening on secret channel", index)

	// iterate through all secrets published to the channel
	for s := range channel {
		logrus.Infof("go routine %d: encrypting the value for secret %d", index, s.GetID())

		// update secret with encryption in the database
		err := d.Client.UpdateSecret(s)
		if err != nil {
			return err
		}

		logrus.Tracef("go routine %d: value encrypted for secret %d", index, s.GetID())
	}

	logrus.Infof("go routine %d: shutting down on secret channel", index)

	return nil
}
