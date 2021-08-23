// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"time"
)

// actions represents all potential actions
// this utility can invoke with the database.
type actions struct {
	All         bool
	AlterTables bool
}

// connection represents all connection related
// information used to communicate with the
// provided database.
type connection struct {
	Idle int
	Life time.Duration
	Open int
}

// db represents all database related
// information used to communicate
// with the database.
type db struct {
	Driver  string
	Address string

	Actions    *actions
	Connection *connection
}
