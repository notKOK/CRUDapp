package endpoint

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Endpoint struct {
	db      *gorm.DB
	mycache *cache.Cache
}

type Product struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Type   string `json:"type"`
	Price  string `json:"price"`
	MadeIn string `json:"madeIn"`
}

func New() *Endpoint {
	endP := &Endpoint{}

	dsn := "host=127.0.0.1 user=cruder password=jw8 dbname=crudapp port=5432 sslmode=disable"
	var err error
	endP.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = endP.db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal(err)
	}
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"localhost": ":6379",
			"server2":   ":6380",
		},
	})

	endP.mycache = cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})
	return endP
}

func (endP *Endpoint) GetEntityByID(c *gin.Context) {

	ctx := context.TODO()

	var wanted Product
	sql := endP.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", c.Param("id")).First(&wanted)
	})

	if err := endP.mycache.Get(ctx, sql, &wanted); err == nil {
		c.IndentedJSON(http.StatusOK, wanted)
		return
	}

	result := endP.db.Where("id = ?", c.Param("id")).First(&wanted)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, "NOT FOUND ENTITY")
		log.Print(result.Error)
		return
	} else if result.Error == nil {
		c.IndentedJSON(http.StatusOK, wanted)
		if err := endP.mycache.Set(&cache.Item{
			Ctx:   ctx,
			Key:   sql,
			Value: wanted,
			TTL:   time.Minute,
		}); err != nil {
			log.Fatal(err)
		}
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
