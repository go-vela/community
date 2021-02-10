// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/go-vela/sdk-go/vela"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// server represents the information you need to
// interact with a Vela server
type server struct {
	// information about server
	Addr  string
	Key   string
	Token string

	// clients for interaction
	vela *vela.Client
}

func main() {
	// set the log level for the utility
	switch os.Getenv("VELA_LOG_LEVEL") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	s := server{
		Addr:  os.Getenv("VELA_ADDR"),
		Key:   os.Getenv("VELA_KEY"),
		Token: os.Getenv("VELA_TOKEN"),
	}

	logrus.Info("validating script setup")
	err := s.validate()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("establishing a connection with the server")
	err = s.login()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("modifying secrets with encryption")
	err = s.modify()
	if err != nil {
		logrus.Fatal(err)
	}
}

// validate the server information
func (s *server) validate() error {
	// check if the addr is set
	if len(s.Addr) == 0 {
		return fmt.Errorf("VELA_ADDR is not properly configured")
	}

	// check if the token is set
	if len(s.Token) == 0 {
		return fmt.Errorf("VELA_TOKEN is not properly configured")
	}

	// check if the key is set
	if len(s.Key) != 32 {
		return fmt.Errorf("VELA_KEY is not properly configured; invalid length specified: %d", len(s.Key))
	}

	return nil
}

// login creates an interactable client for the server
func (s *server) login() error {
	// Create a new vela client for interacting with server
	c, err := vela.NewClient(s.Addr, nil)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}

	// Set new token in existing client
	c.Authentication.SetTokenAuth(s.Token)

	// add the client to the server
	s.vela = c

	return nil
}

// modify reads the secrets in the database and encrypts the
// value then updates the secret to contain the encrypted value.
func (s *server) modify() error {
	logrus.Infof("retrieving all secrets from %s", s.Addr)

	// send API call to capture list of secrets from the database
	secrets, _, err := s.vela.Admin.Secret.GetAll(nil)
	if err != nil {
		return fmt.Errorf("unable to retrieve secrets: %w", err)
	}

	logrus.Infof("iterating through list of %d secrets", len(*secrets))

	// iterate through ll secrets stored in database
	for _, secret := range *secrets {
		// create output string for secret
		var output string

		// format the output string based off secret type
		switch secret.GetType() {
		case constants.SecretRepo:
			// syntax: repo/<org>/<repo>/<name>
			output = fmt.Sprintf("repo/%s/%s/%s", secret.GetOrg(), secret.GetRepo(), secret.GetName())
		case constants.SecretOrg:
			// syntax: org/<org>/*/<name>
			output = fmt.Sprintf("org/%s/*/%s", secret.GetOrg(), secret.GetName())
		case constants.SecretShared:
			// syntax: shared/<org>/<team>/<name>
			output = fmt.Sprintf("shared/%s/%s/%s", secret.GetOrg(), secret.GetTeam(), secret.GetName())
		}

		logrus.Debugf("encrypting secret %s", output)

		// encrypt secret value using provided key
		encVal, err := encrypt([]byte(secret.GetValue()), s.Key)
		if err != nil {
			return fmt.Errorf("unable to encrypt secret %s: %w", output, err)
		}

		// set new secret value that is encrypted
		secret.SetValue(encVal)

		logrus.Debugf("updating encrypted secret %s", output)

		// send API call to update secret in database
		_, _, err = s.vela.Admin.Secret.Update(&secret)
		if err != nil {
			return fmt.Errorf("unable to update secret %s: %w", output, err)
		}
	}

	return nil
}

// This func is a direct copy of Vela's encrypt implementation:
// https://github.com/go-vela/server/blob/master/secret/native/crypto.go#L22
func encrypt(data []byte, key string) (string, error) {
	// within the validate process we enforce a 64 bit key which
	// ensures all secrets are encrypted with AES-256:
	// https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// nonce is an arbitrary number used to to ensure that
	// old communications cannot be reused in replay attacks.
	// https://en.wikipedia.org/wiki/Cryptographic_nonce
	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	// encrypt the data with the randomly generated nonce
	encData := gcm.Seal(nonce, nonce, data, nil)

	// encode the encrypt data to make it network safe
	sEnc := base64.StdEncoding.EncodeToString(encData)

	return sEnc, nil
}
