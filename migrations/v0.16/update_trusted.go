// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

const (
	UpdateReposTrusted = `
UPDATE repos
SET trusted = 'true'
WHERE id IN (
    SELECT id
    FROM repos r 
    INNER JOIN (
        SELECT 
            repo_id
        FROM steps 
        WHERE image LIKE ? AND 
        created > (
            SELECT EXTRACT(EPOCH FROM (NOW() - ? * INTERVAL '1 DAY'))
        )   
        GROUP BY repo_id
    ) t 
    ON r.id = t.repo_id 
    WHERE active = 't'
);
`

	UpdateReposTrustedNoPersonalOrgs = `
UPDATE repos
SET trusted = 'true'
WHERE id IN (
    SELECT rt.id FROM
        (SELECT id, org, user_id, full_name
        FROM repos r 
        INNER JOIN (
            SELECT 
                repo_id
            FROM steps 
            WHERE image LIKE ? AND 
            created > (
                SELECT EXTRACT(EPOCH FROM (NOW() - ? * INTERVAL '1 DAY'))
            )   
            GROUP BY repo_id
        ) t 
        ON r.id = t.repo_id) rt
    INNER JOIN users u
    ON rt.user_id = u.id
    WHERE rt.org != u.name
);
`
)

// UpdateReposTrusted will take a list of privileged images and no. of days
// and update repos that used one of those images within that amount of time
// to be trusted, including personal org repos.
func (d *db) UpdateReposTrusted(images []string, days string) error {
	for _, image := range images {
		pattern := fmt.Sprintf("%%%s%%", image)
		logrus.Infof("updating repos that have used the %s image in the past %s days to be trusted", image, days)
		err := d.Gorm.Exec(UpdateReposTrusted, pattern, days).Error
		if err != nil {
			return fmt.Errorf("unable to update trusted field in %s table: %v", constants.TableRepo, err)
		}
		logrus.Info("update successful")
	}

	return nil
}

// UpdateReposTrusted will take a list of privileged images and no. of days
// and update repos that used one of those images within that amount of time
// to be trusted, NOT including personal org repos.
func (d *db) UpdateReposTrustedNoPersonalOrgs(images []string, days string) error {
	for _, image := range images {
		pattern := fmt.Sprintf("%%%s%%", image)
		logrus.Infof("updating non-personal org repos that have used the %s image in the past %s days to be trusted", image, days)
		err := d.Gorm.Exec(UpdateReposTrustedNoPersonalOrgs, pattern, days).Error
		if err != nil {
			return fmt.Errorf("unable to update trusted field in %s table: %v", constants.TableRepo, err)
		}
		logrus.Info("update successful")
	}

	return nil
}
