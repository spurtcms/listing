package listing

import (
	"fmt"

	"github.com/spurtcms/auth/migration"
)

func ListingSetup(config Config) *Listing {

	fmt.Println("Heloo")

	migration.AutoMigration(config.DB, config.DataBaseType)

	return &Listing{
		AuthEnable:       config.AuthEnable,
		Permissions:      config.Permissions,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
		DB:               config.DB,
	}

}
