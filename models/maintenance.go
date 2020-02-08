package models

//AgeGroup ...
type AgeGroup struct {
	BaseModel
	Code        string `json:"code" gorm:"not null;unique"`
	Description string `json:"description" gorm:"not null;unique"`
	MinAge      int    `json:"min_age"`
	MaxAge      int    `json:"max_age"`
}

// TableName ...
func (ageGroup AgeGroup) TableName() string {
	return "age_group"
}

//IDType ...
type IDType struct {
	BaseModel
	Name      string `json:"name" gorm:"not null;unique"`
	Code      string `json:"code" gorm:"not null;unique"`
	CountryID int    `json:"country_id"`
	IssuedBy  string `json:"issued_by"`
}

// TableName ...
func (title IDType) TableName() string {
	return "id_type"
}

//UserDefinedFlds ... User Defined Fields that are attached to key entities in the system such as
//Person, Members and other user defined forms
type UserDefinedFlds struct {
	BaseModel
	Entity   UDFEntity `json:"entity"`
	EntityID int       `json:"entity_id"`
	Name     string    `json:"name"`
	Type     string    `json:"type"`
}

// TableName ...
func (udf UserDefinedFlds) TableName() string {
	return "udf_table"
}

//UDFEntity ...
type UDFEntity struct {
	BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TableName ...
func (udfentity UDFEntity) TableName() string {
	return "udf_entity"
}

//UDFValues ...
type UDFValues struct {
	BaseModel
	FldName    string    `json:"fld_name"`
	FldValue   string    `json:"fld_value"`
	EntityName string    `json:"entity_name"`
	EntityID   string    `json:"entity_id"`
	Entity     UDFEntity `json:"entity"`
}

//TableName ...
func (udfvalue UDFValues) TableName() string {
	return "udf_value"
}

//ParametersValues ...
type ParametersValues struct {
	BaseModel
	Parameter string `json:"parameter"`
	Entity    string `json:"entity"`
	EntityID  int    `json:"entity_id"`
}

//TableName ...
func (paramval ParametersValues) TableName() string {
	return "parameter_values"
}

//Currency ...
type Currency struct {
	BaseModel
	Name         string `json:"name"  gorm:"not null;unique"`
	Description  string `json:"description"`
	ISOCode      string `json:"iso_code"  gorm:"not null;unique"`
	ISOCodeInt   int    `json:"iso_code_int"  gorm:"not null;unique"`
	DecimalPoint int    `json:"decimal_point"`
}

//TableName ...
func (curr Currency) TableName() string {
	return "currencies"
}
