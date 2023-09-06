package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Product struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Type   string `json:"type"`
	Price  string `json:"price"`
	MadeIn string `json:"madeIn"`
}

var (
	db *gorm.DB
)

func main() {
	dsn := "host=localhost user=cruder password=jw8 dbname=crudapp port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return
	}

	err2 := db.AutoMigrate(&Product{})
	if err2 != nil {
		log.Fatal(err2)
	}

	router := gin.Default()
	err3 := router.SetTrustedProxies(nil)
	if err3 != nil {
		log.Fatal(err3)
		return
	}
	router.GET("/entity/:id", GetEntityByIDHandler)
	router.GET("/entities", GetAllEntitiesHandler)
	router.POST("/up", PostEntityHandler)
	router.PUT("/up/:id", PutEntityHandler)
	router.DELETE("/entity/:id", DeleteEntityByIDHandler)

	err4 := router.Run("localhost:8000")
	if err4 != nil {
		log.Fatal(err4)
		return
	}
}

func GetEntityByIDHandler(c *gin.Context) {
	var baseProduct Product
	result := db.Where("id = ?", c.Param("id")).First(&baseProduct)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, "NOT FOUND ENTITY")
		return
	} else if result.Error == nil {
		c.IndentedJSON(http.StatusOK, baseProduct)
		return
	}
	c.IndentedJSON(http.StatusBadRequest, result.Error)
}

func GetAllEntitiesHandler(c *gin.Context) {
	var products []Product
	result := db.Find(&products)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, result.Error)
		log.Print(result.Error)
		return
	}
	c.IndentedJSON(http.StatusOK, products)
}

func DeleteEntityByIDHandler(c *gin.Context) {
	del := Product{ID: c.Param("id")}
	result := db.Where("id = ?", c.Param("id")).Delete(&del)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, "accepted")
		log.Print(result.Error)
	} else if result.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, "NOT FOUND ENTITY")
	} else {
		c.IndentedJSON(http.StatusAccepted, "accepted")
	}
}

func PostEntityHandler(c *gin.Context) {
	var catchProduct Product
	if err := c.BindJSON(&catchProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	result := db.Create(&catchProduct)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		log.Print(result.Error)
	} else {
		c.JSON(http.StatusAccepted, gin.H{"Status": "accepted"})
	}
}

func PutEntityHandler(c *gin.Context) {
	var catchProduct, baseProduct Product
	if err := c.BindJSON(&catchProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	result := db.Where("id = ?", c.Param("id")).First(&baseProduct)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, "NOT FOUND ENTITY, entity doesn't exist")
		return
	}
	saved := db.Save(&catchProduct)
	if saved.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": saved.Error})
		return
	}
	c.JSON(http.StatusAccepted, "Accepted")
}
