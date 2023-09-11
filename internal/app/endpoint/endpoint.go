package endpoint

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Endpoint struct {
	db *gorm.DB
}

func New() *Endpoint {
	endP := &Endpoint{}

	dsn := "host=localhost user=cruder password=jw8 dbname=crudapp port=5432 sslmode=disable"
	var err error
	endP.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = endP.db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal(err)
	}
	return endP
}

type Product struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Type   string `json:"type"`
	Price  string `json:"price"`
	MadeIn string `json:"madeIn"`
}

func (endP *Endpoint) GetEntityByID(c *gin.Context) {
	var baseProduct Product
	result := endP.db.Where("id = ?", c.Param("id")).First(&baseProduct)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, "NOT FOUND ENTITY")
		log.Print(result.Error)
		return
	} else if result.Error == nil {
		c.IndentedJSON(http.StatusOK, baseProduct)
		return
	}
	c.IndentedJSON(http.StatusBadRequest, result.Error)
	log.Print(result.Error)
}

func (endP *Endpoint) GetAllEntities(c *gin.Context) {
	var products []Product
	result := endP.db.Find(&products)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, result.Error)
		log.Print(result.Error)
		return
	}
	c.IndentedJSON(http.StatusOK, products)
}

func (endP *Endpoint) DeleteEntityByID(c *gin.Context) {
	del := Product{ID: c.Param("id")}
	result := endP.db.Where("id = ?", c.Param("id")).Delete(&del)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, "accepted")
		log.Print(result.Error)
	} else if result.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, "NOT FOUND ENTITY")
	} else {
		c.IndentedJSON(http.StatusAccepted, "accepted")
	}
}

func (endP *Endpoint) PostEntity(c *gin.Context) {
	var catchProduct Product
	if err := c.BindJSON(&catchProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	result := endP.db.Create(&catchProduct)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		log.Print(result.Error)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"Status": "accepted"})
}

func (endP *Endpoint) PutEntity(c *gin.Context) {
	var catchProduct, baseProduct Product
	if err := c.BindJSON(&catchProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	result := endP.db.Where("id = ?", c.Param("id")).First(&baseProduct)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, "NOT FOUND ENTITY, entity doesn't exist")
		log.Print(result.Error)
		return
	}
	saved := endP.db.Save(&catchProduct)
	if saved.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": saved.Error})
		log.Print(saved.Error)
		return
	}
	c.JSON(http.StatusAccepted, "Accepted")
}
