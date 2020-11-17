package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/go-vela/sdk-go/vela"
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

	logrus.Info("establish a connection with the server")
	err = s.login()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("execute the update to encrypt secrets")
	err = s.modify()
	if err != nil {
		logrus.Fatal(err)
	}
}

// validate the server information
func (s *server) validate() error {

	// check the addr is set
	if len(s.Addr) == 0 {
		return fmt.Errorf("addr is not properly configured")
	}

	// check the key is set
	if len(s.Key) != 32 {
		return fmt.Errorf("key is not properly configured; invalid length specified: %d", len(s.Key))
	}

	// check the token is set
	if len(s.Token) == 0 {
		return fmt.Errorf("username is not properly configured")
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
	// get the list of secrets from the database
	secrets, _, err := s.vela.Admin.Secret.GetAll(nil)
	if err != nil {
		return fmt.Errorf("unable to retrieve secrets: %w", err)
	}

	// update each secret value to be encrypted
	for _, secret := range secrets {
		encVal, err := encrypt([]byte(secret.GetValue()), s.Key)
		if err != nil {
			return fmt.Errorf("unable to encrypt secret %s: %w", secret.GetName(), err)
		}

		secret.SetName(encVal)

		// get the list of secrets from the database
		_, _, err = s.vela.Admin.Secret.Update(secret)
		if err != nil {
			return fmt.Errorf("unable to update secret %s: %w", secret.GetName(), err)
		}
	}

	return nil
}

// This func is a direct copy of Vela's encrypt implementation:
// TODO: Add link
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

	return string(gcm.Seal(nonce, nonce, data, nil)), nil
}
