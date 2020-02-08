package controllers

import (
	"fmt"

	"sundialventure.com/common/models"
	"sundialventure.com/common/sys"
	"sundialventure.com/common/utility"

	//entities "sundialventure.com/payschool/entities"

	//"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "gopkg.in/inconshreveable/log15.v2"
)

//MaintenanceController ... Controller for all
type MaintenanceController struct {
	DataStore    utility.DataStore
	HTTPUtilDunc utility.HTTPUtilityFunctions
	Logger       log.Logger
}

/*----------------------------------------------------*/
/*	Country Methods
/*
/*
/*
------------------------------------------------------*/

// CreateCountry ... Add new member
func (man MaintenanceController) CreateCountry(c *gin.Context) {
	country := new(models.Country)
	err := c.BindJSON(country)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.SaveDatabaseObject(&country)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, country)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

//UpdateCountry ...
func (man MaintenanceController) UpdateCountry(c *gin.Context) {
	country := new(models.Country)
	err := c.BindJSON(country)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&country, map[string]interface{}{"updated_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, country)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// DeleteCountry ... Add new member
func (man MaintenanceController) DeleteCountry(c *gin.Context) {
	country := new(models.Country)
	err := c.BindJSON(country)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&country, map[string]interface{}{"status": "DELETE", "deleted_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.CustomErr.Msg)
			return
		}
		c.JSON(200, country)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// RetrieveAllCountry ... Add new member
func (man MaintenanceController) RetrieveAllCountry(c *gin.Context) {
	countries := []models.Country{}
	dberr := man.DataStore.FetchAllGenericObject(&countries)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		man.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, countries)
	return
}

//RetrieveCountryByID ...Retrieve Country By ID
func (man MaintenanceController) RetrieveCountryByID(c *gin.Context) {
	countryid, _ := strconv.Atoi(c.Params.ByName("id"))

	country := models.Country{}
	//regionState := new(models.RegionState)
	var dberr *utility.DataStoreError
	dberr = man.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"id": countryid, "status": "ACTIVE"}, &country)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		man.Logger.Error(dberr.CustomErr.Msg)
		return
	}

	c.JSON(200, country)
	return
}

// SetDefaultCountry ... Set a country as default country. Change the default field to Y and set all other to N
// UPDATE country SET default_country = 'Y' where id = ?;
// UPDATE country SET default_country = null;
func (man MaintenanceController) SetDefaultCountry(c *gin.Context) {
	country := new(models.Country)
	err := c.BindJSON(country)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&country, map[string]interface{}{"updated_at": time.Now(), "default_country": "Y"})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		man.DataStore.RDBMS.DB.Model(&country).Where("id <> ?", country.ID).Update("default_country", nil)
		c.JSON(200, country)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

//GetDefaultCountry ...
func (man MaintenanceController) GetDefaultCountry(c *gin.Context) {
	country := new(models.Country)
	err := man.DataStore.FetchFirstGenericObject(map[string]interface{}{"default_country": "Y", "status": "ACTIVE"}, &country)
	if err == nil {
		c.JSON(200, country)
		return
	}
	man.Logger.Error(err.CustomErr.Msg)
	c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": err.CustomErr.Msg})
}

/*
visitid, _ := strconv.Atoi(c.Params.ByName("id"))
	visitor := new(models.Visitor)
	err := visit.DataStore.FetchFirstGenericObject(map[string]interface{}{"id": visitid, "status": "ACTIVE"}, &visitor)
	if err == nil {
		c.JSON(200, visitor)
		return
	}
	visit.Logger.Error(err.CustomErr.Msg)
	c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": err.CustomErr.Msg})
*/

/*----------------------------------------------------*/
/*	Language Methods
/*
/*
/*
------------------------------------------------------*/

// CreateLanguage ... Add new member
func (man MaintenanceController) CreateLanguage(c *gin.Context) {
	language := new(models.Language)
	err := c.BindJSON(language)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.SaveDatabaseObject(&language)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, language)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

//UpdateLanguage ...
func (man MaintenanceController) UpdateLanguage(c *gin.Context) {
	language := new(models.Language)
	err := c.BindJSON(language)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&language, map[string]interface{}{"updated_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, language)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// DeleteLanguage ... Add new member
func (man MaintenanceController) DeleteLanguage(c *gin.Context) {
	language := new(models.Language)
	err := c.BindJSON(language)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&language, map[string]interface{}{"status": "DELETE", "deleted_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.CustomErr.Msg)
			return
		}
		c.JSON(200, language)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// RetrieveAllLanguage ... Add new member
func (man MaintenanceController) RetrieveAllLanguage(c *gin.Context) {
	languages := []models.Language{}
	dberr := man.DataStore.FetchAllGenericObject(&languages)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		man.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, languages)
	return
}

//RetrieveLanguageByID ...Retrieve Country By ID
func (man MaintenanceController) RetrieveLanguageByID(c *gin.Context) {
	languageid, _ := strconv.Atoi(c.Params.ByName("id"))
	language := models.Language{}
	//regionState := new(models.RegionState)
	var dberr *utility.DataStoreError
	dberr = man.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"id": languageid, "status": "ACTIVE"}, &language)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		man.Logger.Error(dberr.CustomErr.Msg)
		return
	}

	c.JSON(200, language)
	return
}

/*----------------------------------------------------*/
/*	RegionState Methods
/*
/*
/*
------------------------------------------------------*/

// CreateRegionState ... Add new member
func (man MaintenanceController) CreateRegionState(c *gin.Context) {
	regionState := new(models.RegionState)
	err := c.BindJSON(regionState)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.SaveDatabaseObject(&regionState)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, regionState)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

//UpdateRegionState ...
func (man MaintenanceController) UpdateRegionState(c *gin.Context) {
	regionState := new(models.RegionState)
	err := c.BindJSON(regionState)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&regionState, map[string]interface{}{"updated_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, regionState)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// DeleteRegionState ... Add new member
func (man MaintenanceController) DeleteRegionState(c *gin.Context) {
	regionState := new(models.RegionState)
	err := c.BindJSON(regionState)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&regionState, map[string]interface{}{"status": "DELETE", "deleted_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.CustomErr.Msg)
			return
		}
		c.JSON(200, regionState)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// RetrieveAllRegionStateByCountry ... Add new member
func (man MaintenanceController) RetrieveAllRegionStateByCountry(c *gin.Context) {
	regionstates := []models.RegionState{}
	//regionState := new(models.RegionState)
	var dberr *utility.DataStoreError
	countryid, _ := strconv.Atoi(c.Params.ByName("country_id"))
	dberr = man.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"country_id": countryid, "status": "ACTIVE"}, &regionstates)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		man.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, regionstates)
	return
}

// RetrieveAllRegionStateByID ... Add new member
func (man MaintenanceController) RetrieveAllRegionStateByID(c *gin.Context) {
	regionstate := models.RegionState{}
	//regionState := new(models.RegionState)
	var dberr *utility.DataStoreError
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	dberr = man.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"id": id, "status": "ACTIVE"}, &regionstate)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		man.Logger.Error(dberr.CustomErr.Msg)
		return
	}

	c.JSON(200, regionstate)
	return
}

/*----------------------------------------------------*/
/*	IDType Methods
/*
/*
/*
------------------------------------------------------*/

// CreateIDType ... Add new member
func (man MaintenanceController) CreateIDType(c *gin.Context) {
	idtype := new(models.IDType)
	err := c.BindJSON(idtype)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.SaveDatabaseObject(&idtype)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, idtype)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

//UpdateIDType ...
func (man MaintenanceController) UpdateIDType(c *gin.Context) {
	idtype := new(models.IDType)
	err := c.BindJSON(idtype)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&idtype, map[string]interface{}{"updated_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, idtype)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// DeleteIDType ... Add new member
func (man MaintenanceController) DeleteIDType(c *gin.Context) {
	idtype := new(models.IDType)
	err := c.BindJSON(idtype)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.UpdateDatabaseObject(&idtype, map[string]interface{}{"status": "DELETE", "deleted_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.CustomErr.Msg)
			return
		}
		c.JSON(200, idtype)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// RetrieveAllIDType ... Add new member
func (man MaintenanceController) RetrieveAllIDType(c *gin.Context) {
	idtypes := []models.IDType{}
	//idtype := new(models.IDType)
	var dberr *utility.DataStoreError
	dberr = man.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"status": "ACTIVE"}, &idtypes)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		man.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, idtypes)
	return
}

//ParameterController ...
type ParameterController struct {
	//DataStore    utility.DataStore
	//HTTPUtilDunc utility.HTTPUtilityFunctions
	//Logger       log.Logger
	BaseController
}

// CreateParameter ... Add new parish
func (param ParameterController) CreateParameter(c *gin.Context) {
	parameter := new(models.ParametersValues)
	err := c.BindJSON(parameter)
	if err == nil {
		param.BaseController.CreateEntity(c, &parameter)
	}
}

// UpdateParameter ... Update parish
func (param ParameterController) UpdateParameter(c *gin.Context) {
	parameter := new(models.ParametersValues)
	err := c.BindJSON(parameter)
	if err == nil {
		param.UpdateEntity(c, &parameter)
	}

}

// DeleteParameter ... delete parish
func (param ParameterController) DeleteParameter(c *gin.Context) {
	parameter := new(models.ParametersValues)
	err := c.BindJSON(parameter)
	if err == nil {
		param.DeleteEntity(c, &parameter)
	}
}

//RetrieveParameterByEntity ...Retrieve Parish By ID
func (param ParameterController) RetrieveParameterByEntity(c *gin.Context) {
	entityid, _ := strconv.Atoi(c.Params.ByName("entity_id"))
	entity := c.Params.ByName("entity")
	parameter := models.ParametersValues{}
	paramdict := map[string]interface{}{"status": "ACTIVE", "entity": entity, "entity_id": entityid}
	param.RetrieveEntityByFields(c, entityid, parameter, paramdict)
}

//RetrieveParameterDefaults ...Retrieve Parish By ID
func (param ParameterController) RetrieveParameterDefaults(c *gin.Context) {
	parameters := []models.ParametersValues{}
	//idtype := new(models.IDType)
	var dberr *utility.DataStoreError
	dberr = param.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"status": "ACTIVE"}, &parameters)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		param.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, parameters)
	return
}

//AgeGroups Methods

// CreateAgeGroup ... Add new parish
func (param ParameterController) CreateAgeGroup(c *gin.Context) {
	agegrp := new(models.AgeGroup)
	err := c.BindJSON(agegrp)

	if err == nil {
		param.BaseController.CreateEntity(c, &agegrp)
	}
}

// UpdateAgeGroup ... Update parish
func (param ParameterController) UpdateAgeGroup(c *gin.Context) {
	agegrp := new(models.AgeGroup)
	err := c.BindJSON(agegrp)
	if err == nil {
		param.UpdateEntity(c, &agegrp)
	}

}

// DeleteAgeGroup ... delete parish
func (param ParameterController) DeleteAgeGroup(c *gin.Context) {
	agegrp := new(models.AgeGroup)
	err := c.BindJSON(agegrp)
	if err == nil {
		param.DeleteEntity(c, &agegrp)
	}
}

//RetrieveAgeGroupByID ...Retrieve Parish By ID
func (param ParameterController) RetrieveAgeGroupByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	agegrp := models.AgeGroup{}
	paramdict := map[string]interface{}{"status": "ACTIVE", "id": id}
	param.RetrieveEntityByFields(c, id, agegrp, paramdict)
}

//RetrieveAgeGroup ...Retrieve Parish By ID
func (param ParameterController) RetrieveAgeGroup(c *gin.Context) {
	agegrps := []models.AgeGroup{}
	var dberr *utility.DataStoreError
	dberr = param.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"status": "ACTIVE"}, &agegrps)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		param.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, agegrps)
	return
}

//Currency Methods

// CreateCurrency ... Add new Currency
func (param ParameterController) CreateCurrency(c *gin.Context) {
	curr := new(models.Currency)
	err := c.BindJSON(curr)
	fmt.Println(err)
	if err == nil {
		param.BaseController.CreateEntity(c, &curr)
		return
	}
	if err != nil {
		param.Logger.Error(err.Error())
	}

}

// UpdateCurrency ... Update Currency
func (param ParameterController) UpdateCurrency(c *gin.Context) {
	curr := new(models.Currency)
	err := c.BindJSON(curr)
	if err == nil {
		param.UpdateEntity(c, &curr)
	}
}

// DeleteCurrency ... delete Currency
func (param ParameterController) DeleteCurrency(c *gin.Context) {
	curr := new(models.Currency)
	err := c.BindJSON(curr)
	if err == nil {
		param.DeleteEntity(c, &curr)
	}
}

//RetrieveCurrencyByID ...Retrieve Currency By ID
func (param ParameterController) RetrieveCurrencyByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	curr := models.Currency{}
	paramdict := map[string]interface{}{"status": "ACTIVE", "id": id}
	param.RetrieveEntityByFields(c, id, &curr, paramdict)
}

//RetrieveCurrencies ...Retrieve Currency By ID
func (param ParameterController) RetrieveCurrencies(c *gin.Context) {
	currs := []models.Currency{}
	var dberr *utility.DataStoreError
	dberr = param.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"status": "ACTIVE"}, &currs)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		param.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, currs)
	return
}

/*----------------------------------------------------*/
/*	Address Methods
/*
/*
/*
------------------------------------------------------*/

// CreateAddress ... Add new address
func (man MaintenanceController) CreateAddress(c *gin.Context) {
	address := new(models.Address)
	err := c.BindJSON(address)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.SaveDatabaseObject(&address)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, address)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

//RetrieveAddressByEntity ...Retrieve Addresses By Entity
func (man MaintenanceController) RetrieveAddressByEntity(c *gin.Context) {
	entity := c.Params.ByName("entity")
	personid, _ := strconv.Atoi(c.Params.ByName("person_id"))
	addresses := []models.Address{}
	//regionState := new(models.RegionState)
	var dberr *utility.DataStoreError
	dberr = man.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"entity": entity, "person_id": personid, "status": "ACTIVE"}, &addresses)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		man.Logger.Error(dberr.CustomErr.Msg)
		return
	}

	c.JSON(200, addresses)
	return
}

/*----------------------------------------------------*/
/*	Contact Methods
/*
/*
/*
------------------------------------------------------*/

// CreateContact ... Add new contact
func (man MaintenanceController) CreateContact(c *gin.Context) {
	contact := new(models.Contact)
	err := c.BindJSON(contact)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.SaveDatabaseObject(&contact)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, contact)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

/*----------------------------------------------------*/
/*	BAnk Account Methods
/*
/*
/*
------------------------------------------------------*/

// CreateBankAccount ... Add new contact
/*func (man MaintenanceController) CreateBankAccount(c *gin.Context) {
	bankaccounts := new(entities.BankAccounts)
	err := c.BindJSON(bankaccounts)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.SaveDatabaseObject(&bankaccounts)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, bankaccounts)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}*/

/*----------------------------------------------------*/
/*	Wallet Account Methods
/*
/*
/*
------------------------------------------------------*/

// CreateWalletAccount ... Add new wallet account
/*func (man MaintenanceController) CreateWalletAccount(c *gin.Context) {
	walletaccounts := new(entities.WalletAccounts)
	err := c.BindJSON(walletaccounts)
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := man.DataStore.SaveDatabaseObject(&walletaccounts)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			man.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, walletaccounts)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}*/
