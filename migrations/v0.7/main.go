// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os"

	"github.com/go-vela/sdk-go/vela"
	"github.com/sirupsen/logrus"
)

// server represents the information you need to
// interact with a Vela server
type server struct {
	// information about server
	Addr  string
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
		Token: os.Getenv("VELA_TOKEN"),
	}

	logrus.Info("validating server configuration")
	// validate provided server configuration
	err := s.validate()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("creating Vela client from server configuration")
	// create new client from provided server configuration
	err = s.new()
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

	return nil
}

// new creates an interactable client for the server
func (s *server) new() error {
	// Create a new vela client for interacting with server
	c, err := vela.NewClient(s.Addr, "", nil)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}

	// Set new token in existing client
	c.Authentication.SetTokenAuth(s.Token)

	// add the client to the server
	s.vela = c

	return nil
}
