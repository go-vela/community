// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "gorm.io/gorm"

// actions represents all potential actions
// this utility can invoke with the database.
type actions struct {
	All         bool
	AlterTables bool
}

// db represents all database related
// information used to communicate
// with the database.
type db struct {
	Driver  string
	Address string

	Actions *actions
	Gorm    *gorm.DB
}
