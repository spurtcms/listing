package listing
 
import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
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
    Id                    int                   `gorm:"primaryKey;auto_increment;type:serial"`
    Title                 string                `gorm:"type:character varying"`
    Slug                  string                `gorm:"type:character varying"`
    Description           string                `gorm:"type:character varying"`
    ContentType           string                `gorm:"type:character varying"`
    ContentId             int                   `gorm:"type:integer"`
    EntryId               int                   `gorm:"type:integer"`
    IsDeleted             int                   `gorm:"type:integer"`
    DeletedOn             time.Time             `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
    DeletedBy             int                   `gorm:"DEFAULT:NULL"`
    IsActive              int                   `gorm:"type:integer"`
    CreatedOn             time.Time             `gorm:"type:timestamp without time zone"`
    CreatedBy             int                   `gorm:"type:integer"`
    ModifiedOn            time.Time             `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
    ModifiedBy            int                   `gorm:"DEFAULT:NULL;type:integer"`
    ImageName             string                `gorm:"type:character varying"`
    ImagePath             string                `gorm:"type:character varying"`
    VideoName             string                `gorm:"type:character varying"`
    VideoPath             string                `gorm:"type:character varying"`
    Url                   string                `gorm:"type:character varying"`
    PaymentType           string                `gorm:"type:character varying"`
    Price                 int                   `gorm:"type:integer"`
    MembershipId          int                   `gorm:"type:integer"`
    MultiplePrice         datatypes.JSON        `gorm:"type:jsonb"`
    Featured              int                   `gorm:"type:integer"`
    TenantId              string                `gorm:"type:character varying"`
    Tag                   string                `gorm:"type:character varying"`
    MembershipLevel       string                `gorm:"-"`
    SubscriptionName      string                `gorm:"-:migration;<-:false"`
    InitialPayment        string                `gorm:"-:migration;<-:false"`
    ChannelID             int                   `gorm:"-:migration;<-:false"`
    EntryTitle            string                `gorm:"-:migration;<-:false"`
    ChannelName           string                `gorm:"-:migration;<-:false"`
    EntriesId             int                   `gorm:"-:migration;<-:false"`
    ChannelSlug           string                `gorm:"-:migration;<-:false"`
    CategoryName          string                `gorm:"-:migration;<-:false"`
    CourseTitle           string                `gorm:"-:migration;<-:false"`
    MultiplePriceCategory MultiplePriceCategory `gorm:"-"`
}
type MultiplePriceCategory struct {
    BuyNow    int `json:"BuyNow"`
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
    fmt.Println("ListingList:", err)
 
    query := DB.Debug().Table("tbl_listings AS l").
        Select(`
        l.*,
        ce.channel_id,
        ce.title AS entry_title,
        c.channel_name,
        ce.id as entry_id,
        c.slug_name as channel_slug,
        ml.subscription_name,
        ml.initial_payment,
        l.multiple_price,
        cat.category_name,
        co.title as course_title
    `).
 
        //  If NOT course OR entry_id != 0 → join channel entries
        Joins(`LEFT JOIN tbl_channel_entries AS ce
           ON ce.id = l.entry_id
           AND l.content_type != 'course'`).
        //  Join membership only when membership selected
        Joins(`LEFT JOIN tbl_mstr_membershiplevels AS ml
           ON ml.id = l.membership_id
           AND l.membership_id != 0
           AND l.payment_type = 'membership'`).
        //  Join channel when ce exists
        Joins(`LEFT JOIN tbl_channels AS c
           ON c.id = ce.channel_id`).
        //  If course & entry_id = 0 → join courses
        Joins(`LEFT JOIN tbl_courses AS co
           ON l.content_type = 'course'
           AND l.entry_id = 0
           AND co.id = l.content_id`).
        //  join category for course
        Joins(`LEFT JOIN tbl_categories AS cat
           ON l.content_type = 'course'
           AND l.entry_id = 0
           AND cat.id = co.category_id`).
        Where("l.is_deleted = 0 AND l.tenant_id = ?", tenantid).
        Order("l.created_on DESC")
 
    // Keyword Filter
    if filter.Keyword != "" {
        like := "%" + strings.TrimSpace(filter.Keyword) + "%"
        query = query.Where("(LOWER(TRIM(l.title)) LIKE LOWER(TRIM(?)) OR LOWER(TRIM(ce.title)) LIKE LOWER(TRIM(?)))",
            like, like)
    }
 
    //  Title Filter
    if filter.Title != "" {
        like := "%" + strings.TrimSpace(filter.Title) + "%"
        query = query.Where("(LOWER(TRIM(l.title)) LIKE LOWER(TRIM(?)) OR LOWER(TRIM(ce.title)) LIKE LOWER(TRIM(?)))",
            like, like)
    }
 
    //  ContentType Filter
    if filter.ContentType != "" {
        query = query.Where("l.content_type = ?", filter.ContentType)
    }
    fmt.Println(filter.PaymentType, "filter.PaymentType")
    //  PaymentType Filter
    if filter.PaymentType != "" {
        query = query.Where("l.payment_type = ?", filter.PaymentType)
    }
 
    //  Count before pagination
    if err = query.Count(&count).Error; err != nil {
        fmt.Println("Error at count:", err)
        return nil, 0, err
    }
 
    //  Pagination
    if limit != 0 {
        query = query.Limit(limit).Offset(offset)
    }
 
    //  Result scan
    if err = query.Scan(&results).Error; err != nil {
        fmt.Println("Error scanning results:", err)
        return nil, 0, err
    }
    for i := range results {
        if len(results[i].MultiplePrice) > 0 {
            err = json.Unmarshal(results[i].MultiplePrice, &results[i].MultiplePriceCategory)
            if err != nil {
                fmt.Println("JSON Unmarshal Error:", err)
            }
            fmt.Println(results[i].CourseTitle, "results[i].CourseTitle")
        }
    }
 
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
        if err := DB.Table("tbl_listings").Where("id=? and tenant_id=?", listing.Id, listing.TenantId).UpdateColumns(map[string]interface{}{"title": listing.Title, "slug": listing.Slug, "description": listing.Description, "content_type": listing.ContentType, "content_id": listing.ContentId, "entry_id": listing.EntryId, "modified_on": listing.ModifiedOn, "modified_by": listing.ModifiedBy, "image_name": listing.ImageName, "image_path": listing.ImagePath, "video_name": listing.VideoName, "video_path": listing.VideoPath, "url": listing.Url, "payment_type": listing.PaymentType, "price": listing.Price, "membership_id": listing.MembershipId, "multiple_price": listing.MultiplePrice, "tag": listing.Tag}).Error; err != nil {
 
            return err
        }
 
    } else {
 
        if err := DB.Table("tbl_listings").Where("id=? and tenant_id=?", listing.Id, listing.TenantId).UpdateColumns(map[string]interface{}{"title": listing.Title, "slug": listing.Slug, "description": listing.Description, "content_type": listing.ContentType, "content_id": listing.ContentId, "entry_id": listing.EntryId, "modified_on": listing.ModifiedOn, "modified_by": listing.ModifiedBy, "video_name": listing.VideoName, "video_path": listing.VideoPath, "url": listing.Url, "payment_type": listing.PaymentType, "price": listing.Price, "membership_id": listing.MembershipId, "multiple_price": listing.MultiplePrice, "tag": listing.Tag}).Error; err != nil {
 
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