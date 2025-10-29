package listing

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Filter struct {
	Keyword     string
	Title       string
	ContentType string
	PaymentType string
}

type TblListing struct {
	Id               int               `gorm:"primaryKey;auto_increment;type:serial"`
	Title            string            `gorm:"type:character varying"`
	Slug             string            `gorm:"type:character varying"`
	Description      string            `gorm:"type:character varying"`
	ContentType      string            `gorm:"type:character varying"`
	ContentId        int               `gorm:"type:integer"`
	EntryId          int               `gorm:"type:integer"`
	IsDeleted        int               `gorm:"type:integer"`
	DeletedOn        time.Time         `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy        int               `gorm:"DEFAULT:NULL"`
	IsActive         int               `gorm:"type:integer"`
	CreatedOn        time.Time         `gorm:"type:timestamp without time zone"`
	CreatedBy        int               `gorm:"type:integer"`
	ModifiedOn       time.Time         `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy       int               `gorm:"DEFAULT:NULL;type:integer"`
	ImageName        string            `gorm:"type:character varying"`
	ImagePath        string            `gorm:"type:character varying"`
	Url              string            `gorm:"type:character varying"`
	PaymentType      string            `gorm:"type:character varying"`
	Price            int               `gorm:"type:integer"`
	MembershipId     int               `gorm:"type:integer"`
	MultiplePrice    datatypes.JSONMap `gorm:"type:json"`
	Featured         int               `gorm:"type:integer"`
	TenantId         string            `gorm:"type:character varying"`
	Tag              string            `gorm:"type:character varying"`
	MembershipLevel  string            `gorm:"-"`
	SubscriptionName string            `gorm:"-:migration;<-:false"`
	InitialPayment   string            `gorm:"-:migration;<-:false"`
	ChannelID        int               `gorm:"-:migration;<-:false"`
	EntryTitle       string            `gorm:"-:migration;<-:false"`
	ChannelName      string            `gorm:"-:migration;<-:false"`
	EntriesId        int               `gorm:"-:migration;<-:false"`
	ChannelSlug      string            `gorm:"-:migration;<-:false"`
}

type MultiplePrice struct {
	Buynow    int `json:"Buynow"`
	Integrate int `json:"Integrate"`
	Support   int `json:"Support"`
}

type ListingModel struct {
	Userid     int
	DataAccess int
}

var Listingmodels ListingModel

// UpdateListingStatus updates the is_active field of a listing by ID.

func (Listingmodel ListingModel) UpdateListingStatus(limit, offset int, filter Filter, tenantid string, DB *gorm.DB, id int, status int) error {
	// result := DB.Model(&TblListing{}).Where("id = ?", id).Update("featured", status)

	result := DB.Table("tbl_listings").
		Where("id = ? AND is_deleted = 0 AND tenant_id = ?", id, tenantid).
		Updates(map[string]interface{}{
			"featured": status,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no record found to update")
	}
	return nil
}

func (Listingmodel ListingModel) ListingList(limit, offset int, filter Filter, tenantid string, DB *gorm.DB) (results []TblListing, count int64, err error) {
	fmt.Println("ListingList Entry")

	query := DB.Table("tbl_listings AS l").
		Select("l.*, ce.channel_id, ce.title AS entry_title, c.channel_name, ce.id, c.slug_name").
		Joins("JOIN tbl_channel_entries AS ce ON ce.id = l.entry_id").
		Joins("JOIN tbl_channels AS c ON c.id = ce.channel_id").
		Where("l.is_deleted = 0 AND l.tenant_id = ?", tenantid).
		Order("l.created_on DESC")

	// Filters
	if filter.Keyword != "" {
		query = query.Where("LOWER(TRIM(l.title)) LIKE LOWER(TRIM(?))", "%"+filter.Keyword+"%")
	}
	if filter.Title != "" {
		query = query.Where("LOWER(TRIM(l.title)) LIKE LOWER(TRIM(?))", "%"+filter.Title+"%")
	}
	if filter.ContentType != "" {
		query = query.Where("l.content_type = ?", filter.ContentType)
	}
	if filter.PaymentType != "" {
		query = query.Where("l.payment_type = ?", filter.PaymentType)
	}

	// Count total before pagination
	if err = query.Count(&count).Error; err != nil {
		fmt.Println("Error at count:", err)
		return nil, 0, err
	}

	// Pagination
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}

	// Scan into custom struct
	if err = query.Scan(&results).Error; err != nil {
		fmt.Println("Error scanning results:", err)
		return nil, 0, err
	}

	fmt.Println("Listing list for table show:", results)
	return results, count, nil
}

func (Listingmodel ListingModel) CreateListing(listing TblListing, DB *gorm.DB) error {

	if err := DB.Table("tbl_listings").Create(&listing).Error; err != nil {

		return err
	}

	return nil

}

func (Listingmodel ListingModel) EditListing(id int, tenantid string, DB *gorm.DB) (listinglist TblListing, err error) {

	if err := DB.Table("tbl_listings").Where("id=? and tenant_id=? and is_deleted=0", id, tenantid).First(&listinglist).Error; err != nil {

		return TblListing{}, err
	}

	return listinglist, nil
}

func (Listingmodel ListingModel) UpdateListing(listing TblListing, DB *gorm.DB) error {

	if listing.ImageName != "" {
		fmt.Println("Update1::")
		if err := DB.Table("tbl_listings").Where("id=? and tenant_id=?", listing.Id, listing.TenantId).UpdateColumns(map[string]interface{}{"title": listing.Title, "slug": listing.Slug, "description": listing.Description, "content_type": listing.ContentType, "content_id": listing.ContentId, "entry_id": listing.EntryId, "modified_on": listing.ModifiedOn, "modified_by": listing.ModifiedBy, "image_name": listing.ImageName, "image_path": listing.ImagePath, "url": listing.Url, "payment_type": listing.PaymentType, "price": listing.Price, "membership_id": listing.MembershipId, "tag": listing.Tag}).Error; err != nil {

			return err
		}

	} else {

		if err := DB.Table("tbl_listings").Where("id=? and tenant_id=?", listing.Id, listing.TenantId).UpdateColumns(map[string]interface{}{"title": listing.Title, "slug": listing.Slug, "description": listing.Description, "content_type": listing.ContentType, "content_id": listing.ContentId, "entry_id": listing.EntryId, "modified_on": listing.ModifiedOn, "modified_by": listing.ModifiedBy, "url": listing.Url, "payment_type": listing.PaymentType, "price": listing.Price, "membership_id": listing.MembershipId, "tag": listing.Tag}).Error; err != nil {

			return err
		}
	}

	return nil
}

func (Listingmodel ListingModel) DeleteListing(id int, tenantid string, deletedby int, deletedon time.Time, DB *gorm.DB) error {

	if err := DB.Table("tbl_listings").Where("id=? and tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_by": deletedby, "deleted_on": deletedon}).Error; err != nil {

		return err
	}

	return nil
}

func (Listingmodel ListingModel) MultiSelectListingsDelete(listing *TblListing, id []int, DB *gorm.DB) error {

	if err := DB.Table("tbl_listings").Where("id in (?) and tenant_id=?", id, listing.TenantId).UpdateColumns(map[string]interface{}{"is_deleted": listing.IsDeleted, "deleted_on": listing.DeletedOn, "deleted_by": listing.DeletedBy}).Error; err != nil {

		return err
	}

	return nil

}

func (Listingmodel ListingModel) FetchListingsByIds(ids []string, tag string, tenantid string, DB *gorm.DB) (listing []TblListing, err error) {

	if tag == "" {

		if err := DB.Table("tbl_listings").
			Select("tbl_listings.*, tbl_mstr_membershiplevels.subscription_name as subscription_name, tbl_mstr_membershiplevels.initial_payment as initial_payment").
			Joins("LEFT JOIN tbl_mstr_membershiplevels ON tbl_mstr_membershiplevels.id = tbl_listings.membership_id").
			Where("tbl_listings.id IN (?) AND tbl_listings.tenant_id = ? and tbl_listings.is_deleted=0", ids, tenantid).
			Scan(&listing).Error; err != nil {

			return []TblListing{}, err
		}
	} else if tag != "" {

		if err := DB.Table("tbl_listings").
			Select("tbl_listings.*, tbl_mstr_membershiplevels.subscription_name as subscription_name, tbl_mstr_membershiplevels.initial_payment as initial_payment").
			Joins("LEFT JOIN tbl_mstr_membershiplevels ON tbl_mstr_membershiplevels.id = tbl_listings.membership_id").
			Where("tbl_listings.id IN (?) AND tbl_listings.tag=?  AND tbl_listings.tenant_id = ? and tbl_listings.is_deleted=0", ids, tag, tenantid).
			Scan(&listing).Error; err != nil {

			return []TblListing{}, err
		}
	}

	return listing, nil

}

func (Listingmodel ListingModel) FetchListingBySlugName(ids []string, slugname string, tenantid string, DB *gorm.DB) (listing TblListing, err error) {

	if err := DB.Table("tbl_listings").
		Select("tbl_listings.*, tbl_mstr_membershiplevels.subscription_name as subscription_name, tbl_mstr_membershiplevels.initial_payment as initial_payment").
		Joins("LEFT JOIN tbl_mstr_membershiplevels ON tbl_mstr_membershiplevels.id = tbl_listings.membership_id").
		Where("tbl_listings.id IN (?) AND tbl_listings.slug=? AND tbl_listings.tenant_id = ?", ids, slugname, tenantid).
		First(&listing).Error; err != nil {

		return TblListing{}, err
	}

	return listing, nil

}
