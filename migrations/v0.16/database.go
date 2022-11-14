// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "gorm.io/gorm"

// actions represents all potential actions
// this utility can invoke with the database.
type actions struct {
	All             bool
	SetUntrusted    bool
	UpdatePrivRepos bool
}

// trustedOptions represents all the options
// for updating the trusted field for repos in the database.
type trustedOptions struct {
	Images       []string
	PersonalOrgs bool
	Days         string
}

// db represents all database related
// information used to communicate
// with the database.
type db struct {
	Driver  string
	Address string

	Actions        *actions
	TrustedOptions *trustedOptions
	Gorm           *gorm.DB
}
