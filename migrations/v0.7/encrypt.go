// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"

	"golang.org/x/sync/errgroup"
)

// Encrypt does stuff...
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
	secretChannel := make(chan *library.Secret)

	// add set limit of routines to errgroup
	// and begin processing secrets
	for i := 0; i < d.ConcurrencyLimit; i++ {
		// https://golang.org/doc/faq#closures_and_goroutines
		tmp := i

		// spawn a goroutine to begin encrypting secret
		// values that are published to the channel
		group.Go(func() error {
			return d.EncryptSecrets(tmp, secretChannel)
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
		secretChannel <- secret

		logrus.Debugf("secret %d published to channel", secret.GetID())
	}

	logrus.Debug("closing channel for publishing secrets")

	// close channel to signal goroutines to stop processing
	close(secretChannel)

	logrus.Debug("waiting for goroutines to complete")

	return group.Wait()
}

func (d *db) EncryptSecrets(index int, secretChannel chan *library.Secret) error {
	logrus.Infof("thread %d: listening on secret channel", index)

	// iterate through all secrets published to the channel
	for s := range secretChannel {
		logrus.Infof("thread %d: encrypting the value for secret %d", index, s.GetID())

		// update secret with encryption in the database
		err := d.Client.UpdateSecret(s)
		if err != nil {
			return err
		}

		logrus.Debugf("thread %d: value encrypted for secret %d", index, s.GetID())
	}

	logrus.Infof("thread %d: shutting down on secret channel", index)

	return nil
}
