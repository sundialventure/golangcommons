package models

import "time"

//User ...
type User struct {
	BaseModel
	UserName          string `json:"user_name" gorm:"not null;unique"`
	Salt              string `json:"salt"`
	Provider          string `json:"provider"`
	VerificationToken string `json:"verification_token"`
	Email             string `json:"email" gorm:"not null;unique"`
	Password          string `json:"password"`
	PhoneNo           string `json:"phone_no"`
	Active            bool   `json:"active"`
	UUID              string
	UserGroup         UserGroup
	UserGroupID       int64  `json:"user_group_id" sql:"type:bigint REFERENCES user_group(id)"`
	Person            Person `json:"person"`
	PersonID          *int64 `json:"person_id" sql:"type:bigint REFERENCES person(id)"`
}

//TableName ...
func (user User) TableName() string {
	return "users"
}

//UserGroup ...
type UserGroup struct {
	BaseModel
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description" gorm:"not null;unique"`
}

// TableName ...
func (ind *UserGroup) TableName() string {
	return "user_group"
}

//BaseModel ...
type BaseModel struct {
	ID        int64      `json:"id" gorm:"primary_key"`
	CreatedAt *time.Time `json:"created_at"` //sql:"DEFAULT: GETDATE()"
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Status    string     `json:"status" sql:"DEFAULT:'ACTIVE'"`
	RecStatus string     `json:"rec_status" sql:"DEFAULT:'A'"`
}

// TokenAuthentication ...
type TokenAuthentication struct {
	Token string `json:"token" form:"token"`
}

//Contact ...
type Contact struct {
	BaseModel
	PersonID   int64  `json:"person_id" gorm:"not null;unique" sql:"type:bigint REFERENCES person(id)"`
	Entity     string `json:"entity"`
	EntityID   int64  `json:"entity_id"`
	PhoneNo    string `json:"phone_no"`
	PhoneNo2   string `json:"phone_no2"`
	PhoneNo3   string `json:"phone_no3"`
	Email      string `json:"email"`
	OtherMail  string `json:"other_mail"`
	WebSite    string `json:"web_site"`
	FacebookID string `json:"facebook_id"`
	GoogleID   string `json:"google_id"`
	//UDFValues  []UDFValues `json:"user_defined_fields"`
}

// TableName ...
func (contact Contact) TableName() string {
	return "contact"
}

//Address ...
type Address struct {
	BaseModel
	PersonID    int64  `json:"person_id" gorm:"not null" sql:"type:bigint REFERENCES person(id)"`
	FamilyID    int64  `json:"family_id"`
	Entity      string `json:"entity"`
	EntityID    int    `json:"entity_id"`
	AddressType string `json:"address_type" sql:"DEFAULT:'HOUSE'"`
	HouseNo     string `json:"house_no"`
	Street      string `json:"street"`
	Area        string `json:"area"`
	CityTown    string `json:"city_town"`
	LandMark    string `json:"land_mark"`
	Coordinate  string `json:"coordinate"`
	RegionState int64  `json:"region_state" sql:"type:bigint REFERENCES region_state(id)"`
	Country     int64  `json:"country"`
	Address     string `json:"address" gorm:"not null"`
	Address2    string `json:"address2"`
	//RegionState   RegionState `gorm:"foreignkey:UserName"`
	//UDFValues   []UDFValues `json:"user_defined_fields"`
}

// TableName ...
func (address Address) TableName() string {
	return "address"
}

//Person ...
type Person struct {
	BaseModel
	//gorm.Model
	//ParishID           int       `json:"parish_id"`
	FirstName            string    `json:"first_name" gorm:"not null" gorm:"unique_index:idx_fname_lname_oname"`
	OtherName            string    `json:"other_name" gorm:"unique_index:idx_fname_lname_oname"`
	LastName             string    `json:"last_name" gorm:"not null" gorm:"unique_index:idx_fname_lname_oname"`
	Gender               string    `json:"gender"`
	DateOfBirth          time.Time `json:"date_of_birth"`
	DateJoined           time.Time `json:"date_joined" sql:"timestamp with time zone DEFAULT ('now'::text)::date"`
	AgeGroup             *int64    `json:"age_group"`
	Title                *int64    `json:"title"`
	CountryOfOriginID    *int64    `json:"country_of_origin_id" sql:"type:bigint REFERENCES country(id)"`
	CountryOfResidenceID *int64    `json:"country_of_residence_id" sql:"type:bigint REFERENCES country(id)"`
	RegionStateID        *int64    `json:"region_state_id" sql:"type:bigint REFERENCES region_state(id)"`
	PictureFile          string    `json:"picture_file"`
	MaritalStatus        string    `json:"marital_status" sql:"DEFAULT:'SINGLE'"`
	LanguageID           *int64    `json:"language_id" sql:"type:bigint REFERENCES language(id) DEFAULT 1 "`
	//UDFValues          []UDFValues `json:"user_defined_fields"`
	Addresses []Address `json:"addresses"`
	Contacts  []Contact `json:"contacts"`
	Image     string    `json:"image"`
}

// TableName ...
func (person Person) TableName() string {
	return "person"
}

//PersonIDType ...
type PersonIDType struct {
	BaseModel
	IDTypeID  int64      `json:"id_type_id" gorm:"not null;unique"`
	PersonID  int64      `json:"person_id" gorm:"not null;unique"`
	IDNumber  string     `json:"id_number" gorm:"not null"`
	ValidFrom *time.Time `json:"valid_from"`
	ValidTo   *time.Time `json:"valid_to"`
}

// TableName ...
func (personid PersonIDType) TableName() string {
	return "person_id_type"
}

//Language ...
type Language struct {
	BaseModel
	Name    string `json:"name" gorm:"not null;unique"`
	ISOCode string `json:"iso_code" gorm:"not null;unique"`
	//UserDefinedFlds []UserDefinedFlds `json:"user_defined_fields"`
}

// TableName ...
func (language Language) TableName() string {
	return "language"
}

//Country ...
type Country struct {
	BaseModel
	Name           string `json:"name" gorm:"not null;unique"`
	DefaultCountry string `json:"default_country"`
	ISOCode        string `json:"iso_code" gorm:"not null;unique"`
	//UserDefinedFlds []UserDefinedFlds `json:"user_defined_fields"`
	Languages    []Language `json:"languages"`
	CurrencyCode string     `json:"currency_code"`
	LanguageID   int64      `json:"language_id"`
	DialCode     string     `json:"dial_code"`
}

// TableName ...
func (ctry Country) TableName() string {
	return "country"
}

//RegionState ...
type RegionState struct {
	BaseModel
	Name         string `json:"name" gorm:"not null;unique"`
	CountryID    int64  `json:"country_id"`
	ISOCode      string `json:"iso_code"`
	Description  string `json:"description"`
	MainCityTown string `json:"main_city_town"`
	//UserDefinedFlds []UserDefinedFlds `json:"user_defined_fields"`
	Languages []Language `json:"languages"`
}

// TableName ...
func (rstate RegionState) TableName() string {
	return "region_state"
}
