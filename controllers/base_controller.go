package controllers

import (
	models "github.com/sundialventure/golangcommons/models"

	"io/ioutil"
	"strings"
	"time"

	"github.com/sundialventure/golangcommons/sys"
	"github.com/sundialventure/golangcommons/utility"

	"github.com/gin-gonic/gin"
	log "gopkg.in/inconshreveable/log15.v2"
)

//BaseController ... The base controller struct from which all other controller inherit
type BaseController struct {
	DataStore    utility.DataStore
	HTTPUtilDunc utility.HTTPUtilityFunctions
	Logger       log.Logger
}

//SaveUDFValues ...
func (base BaseController) SaveUDFValues(udfvalues []models.UDFValues) ([]models.UDFValues, *utility.DataStoreError) {
	dberr := base.DataStore.SaveDatabaseObject(&udfvalues)
	if dberr != nil {
		base.Logger.Error(dberr.CustomErr.Msg)
		return nil, dberr
	}
	return udfvalues, nil
}

//RetrieveUDFValues ...
func (base BaseController) RetrieveUDFValues(entityid int, entityname string) (string, *utility.DataStoreError) {
	udfvalues := []models.UDFValues{}

	dberr := base.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"entity_id": entityid, "entity_name": entityname}, &udfvalues)
	if dberr != nil {
		return "", dberr
	}
	jsonString := "["
	for _, udfval := range udfvalues {
		jsonString = "{\"" + udfval.FldName + "\":" + udfval.FldValue + "},"
	}
	jsonString = strings.TrimRight(jsonString, ",")
	jsonString = jsonString + "]"
	return jsonString, nil
}

// CreateEntity ... Add new entity
func (base BaseController) CreateEntity(c *gin.Context, entity interface{}) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := base.DataStore.SaveDatabaseObject(entity)
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			base.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, entity)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// UpdateEntity ... Update entity
func (base BaseController) UpdateEntity(c *gin.Context, entity interface{}) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := base.DataStore.UpdateDatabaseObject(entity, map[string]interface{}{"updated_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			base.Logger.Error(dberr.Error())
			return
		}
		c.JSON(200, entity)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// DeleteEntity ... Add delete entity
func (base BaseController) DeleteEntity(c *gin.Context, entity interface{}) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err == nil {
		dberr := base.DataStore.UpdateDatabaseObject(&entity, map[string]interface{}{"status": "DELETE", "deleted_at": time.Now()})
		if dberr != nil {
			c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": body})
			base.Logger.Error(dberr.CustomErr.Msg)
			return
		}
		c.JSON(200, entity)
		return
	}
	c.JSON(500, gin.H{"message": sys.ERR_UNMARSHAL_BODY, "status": sys.FAIL, "payload": body})
}

// RetrieveEntityByID ... Add delete entity
func (base BaseController) RetrieveEntityByID(c *gin.Context, id int, entity interface{}) {
	err := base.DataStore.FetchFirstGenericObject(map[string]interface{}{"id": id, "status": "ACTIVE"}, &entity)
	if err == nil {
		c.JSON(200, entity)
		return
	}
	base.Logger.Error(err.CustomErr.Msg)
	c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": err.CustomErr.Msg})
}

// RetrieveEntityByFields ... Add delete entity
func (base BaseController) RetrieveEntityByFields(c *gin.Context, id int, entity interface{}, parameters map[string]interface{}, relatedent ...interface{}) {

	var dberr *utility.DataStoreError
	dberr = base.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"id": id, "status": "ACTIVE"}, entity, relatedent)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		base.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, entity)
	return

}

// RetrieveAllEntityByFields ... Add delete entity
func (base BaseController) RetrieveAllEntityByFields(c *gin.Context, entity interface{}, parameters map[string]interface{}) {
	err := base.DataStore.FetchAllGenericObjectWithParam(parameters, &entity)
	if err == nil {
		c.JSON(200, entity)
		return
	}
	base.Logger.Error(err.CustomErr.Msg)
	c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": err.CustomErr.Msg})
}

//RetrieveAllEntity ...
func (base BaseController) RetrieveAllEntity(c *gin.Context, entities interface{}, parameters map[string]interface{}) {
	var dberr *utility.DataStoreError
	dberr = base.DataStore.SearchAnyGenericObjectList(map[string]interface{}{"status": "ACTIVE"}, &entities)
	if dberr != nil {
		c.JSON(500, gin.H{"message": sys.FAIL, "status": sys.FAIL, "payload": dberr.CustomErr.Msg})
		base.Logger.Error(dberr.CustomErr.Msg)
		return
	}
	c.JSON(200, entities)
	return
}
