package listing

import (
	"fmt"
	"strings"
	"time"

	"github.com/spurtcms/auth/migration"
)

func ListingSetup(config Config) *Listing {

	migration.AutoMigration(config.DB, config.DataBaseType)

	return &Listing{
		AuthEnable:       config.AuthEnable,
		Permissions:      config.Permissions,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
		DB:               config.DB,
	}

}
func (listing *Listing) FeaturedListingsList(limit, offset int, filter Filter, tenantid string, id, featured int) (err error) {
	// (list []TblListing, Count int64, err error)
	var Autherr error

	if Autherr = AuthandPermission(listing); Autherr != nil {

		return Autherr
	}

	if filter.ContentType != "" {

		filter.ContentType = strings.ToLower(filter.ContentType)
	}

	if filter.PaymentType == "Price" {

		filter.PaymentType = strings.ToLower(filter.PaymentType)

	} else if filter.PaymentType == "Membership" {

		filter.PaymentType = strings.ToLower(filter.PaymentType)

	}

	err = Listingmodels.UpdateListingStatus(limit, offset, filter, tenantid, listing.DB, id, featured)
	if err != nil {
		return Autherr

	}
	return Autherr
}

func (listing *Listing) ListingsList(limit, offset int, filter Filter, tenantid string) (list []TblListing, Count int64, err error) {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return []TblListing{}, 0, Autherr
	}

	// if filter.ContentType == "Course" {

	//  filter.ContentType = strings.ToLower(filter.ContentType)

	// } else if filter.ContentType == "Channel" {

	//  filter.ContentType = strings.ToLower(filter.ContentType)

	// }

	if filter.ContentType != "" {

		filter.ContentType = strings.ToLower(filter.ContentType)
	}

	switch filter.PaymentType {
	case "Price":

		filter.PaymentType = strings.ToLower(filter.PaymentType)

	case "Membership":

		filter.PaymentType = strings.ToLower(filter.PaymentType)

	case "Multipleprice":

		filter.PaymentType = strings.ToLower(filter.PaymentType)

	}

	listinglist, _, _ := Listingmodels.ListingList(limit, offset, filter, tenantid, listing.DB)

	_, count, err := Listingmodels.ListingList(0, 0, filter, tenantid, listing.DB)
	if err != nil {

		return []TblListing{}, 0, err
	}

	return listinglist, count, nil

}

func (listing *Listing) CreateListing(create TblListing) error {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return Autherr
	}

	createdon, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	Create := TblListing{
		Title:         create.Title,
		Slug:          create.Slug,
		Description:   create.Description,
		ContentType:   create.ContentType,
		ContentId:     create.ContentId,
		EntryId:       create.EntryId,
		IsDeleted:     0,
		IsActive:      1,
		CreatedOn:     createdon,
		CreatedBy:     create.CreatedBy,
		ImagePath:     create.ImagePath,
		ImageName:     create.ImageName,
		VideoName:     create.VideoName,
		VideoPath:     create.VideoPath,
		Url:           create.Url,
		PaymentType:   create.PaymentType,
		Price:         create.Price,
		MembershipId:  create.MembershipId,
		MultiplePrice: create.MultiplePrice,
		Tag:           create.Tag,
		TenantId:      create.TenantId,
	}

	err := Listingmodels.CreateListing(Create, listing.DB)

	if err != nil {

		return err
	}

	return nil

}

func (listing *Listing) EditListings(id int, tenantid string) (list TblListing, err error) {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return TblListing{}, Autherr
	}

	list, err = Listingmodels.EditListing(id, tenantid, listing.DB)
	if err != nil {
		fmt.Println(err)
	}

	return list, nil

}

func (listing *Listing) UpdateListings(update TblListing) error {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return Autherr
	}

	modifiedon, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	Update := TblListing{
		Id:            update.Id,
		Title:         update.Title,
		Slug:          update.Slug,
		Description:   update.Description,
		ContentType:   update.ContentType,
		ContentId:     update.ContentId,
		EntryId:       update.EntryId,
		IsDeleted:     0,
		IsActive:      1,
		ModifiedOn:    modifiedon,
		ModifiedBy:    update.ModifiedBy,
		ImagePath:     update.ImagePath,
		ImageName:     update.ImageName,
		VideoName:     update.VideoName,
		VideoPath:     update.VideoPath,
		Url:           update.Url,
		PaymentType:   update.PaymentType,
		Price:         update.Price,
		MembershipId:  update.MembershipId,
		MultiplePrice: update.MultiplePrice,
		Tag:           update.Tag,
		TenantId:      update.TenantId,
	}

	err := Listingmodels.UpdateListing(Update, listing.DB)

	if err != nil {

		return err
	}

	return nil

}

func (listing *Listing) DeleteListing(id, userid int, tenantid string) error {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return Autherr
	}

	deletedon, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	deletedby := userid

	err := Listingmodels.DeleteListing(id, tenantid, deletedby, deletedon, listing.DB)

	if err != nil {

		return err
	}

	return nil

}

func (listing *Listing) MultiSelectDeleteListing(listingids []int, modifiedby int, tenantid string) error {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return Autherr
	}

	var Listing TblListing

	Listing.DeletedBy = modifiedby

	Listing.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	Listing.IsDeleted = 1

	Listing.TenantId = tenantid

	err := Listingmodels.MultiSelectListingsDelete(&Listing, listingids, listing.DB)
	if err != nil {

		return err

	}
	return nil
}

func (listing *Listing) GetListingsByIds(ids []string, tag string, tenantid string) (listings []TblListing, err error) {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return []TblListing{}, Autherr
	}

	listingslist, err := Listingmodels.FetchListingsByIds(ids, tag, tenantid, listing.DB)
	if err != nil {

		return []TblListing{}, err

	}
	return listingslist, nil
}

func (listing *Listing) GetListings(tag string, tenantid string) (listings []TblListing, err error) {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return []TblListing{}, Autherr
	}

	listingslist, err := Listingmodels.FetchListings(tag, tenantid, listing.DB)
	if err != nil {

		return []TblListing{}, err

	}
	return listingslist, nil
}

func (listing *Listing) GetListingBySlugName( slugname string, tenantid string) (listings TblListing, err error) {

	if Autherr := AuthandPermission(listing); Autherr != nil {

		return TblListing{}, Autherr
	}

	listingslist, err := Listingmodels.FetchListingBySlugName(slugname, tenantid, listing.DB)
	if err != nil {

		return TblListing{}, err

	}
	return listingslist, nil
}

// Check category name already exists
func (listing *Listing) CheckListingsName(listingid int, listingname string, tenantid string) (bool, error) {

	var listings TblListing

	err := Listingmodels.CheckListingsName(listings, listingid, listingname, tenantid, listing.DB)

	if err != nil {

		return false, err

	}

	return true, nil
}
