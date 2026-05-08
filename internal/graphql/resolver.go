package graphql

import (
	"github.com/user/library-api/internal/db"
	"github.com/user/library-api/internal/openlibrary"
)

type Resolver struct {
	Store    db.StoreInterface
	OLClient openlibrary.Client
}
