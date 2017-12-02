// Copyright (c) Sarjono Mukti Aji <me@simukti.net>
// See LICENSE for details.

package memory

import (
	"github.com/hashicorp/go-memdb"
)

func migrations() *memdb.DBSchema {
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"password": {
				Name: "password",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:   "id",
						Unique: true,
						Indexer: &memdb.StringFieldIndex{
							// PasswordModel.ID
							Field: "ID",
						},
					},
				},
			},
		},
	}
}
