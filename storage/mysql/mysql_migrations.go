// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package mysql

import "github.com/GuiaBolso/darwin"

// migrations migration will run every app start to make sure tables exists
func migrations() []darwin.Migration {
	return []darwin.Migration{
		{
			Version:     1,
			Description: "CREATE PASSWORD TABLE IF NOT EXISTS",
			Script: `
			CREATE TABLE IF NOT EXISTS password (
				id 			VARCHAR(20) NOT NULL COMMENT 'https://github.com/rs/xid',
				type		VARCHAR(12) NOT NULL,
				hash		VARCHAR(128) NOT NULL,
				secret		VARCHAR(128) NOT NULL,
				created_at	DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

				PRIMARY KEY (id),
				KEY type (type ASC)
		  	) ENGINE = InnoDB
		  	`,
		},
	}
}
