package entity

import "database/sql"

type Product struct {
	ID      uint64
	Name    string
	Comment sql.NullString
}
