package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Counter struct {
	Id int `json:"id"`
	Value int `json:"value" gorm:"default:0"`
}

func Database() *gorm.DB {
	db, err := gorm.Open("sqlite3", "C:/sqlite/api.db")

	if err != nil {
		panic("Failed to connect to the database")
	}

	return db
}

func GetCounters(context *gin.Context) {
	var counters []Counter

	db := Database()
	db.Find(&counters)

	context.JSON(http.StatusOK, gin.H{"data": counters})
}

func GetCounter(context *gin.Context) {
	var counter Counter
	counterId := context.Param("id")

	db := Database()
	db.First(&counter, counterId)

	if counter.Id == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"error": 
				"The counter with id equals " + counterId + " wasn't found"})
	} else {
		context.JSON(http.StatusOK, gin.H{"data": counter})
	}
}

func CreateCounter(context *gin.Context) {
	var counter Counter

	db := Database()
	db.Save(&counter)

	context.JSON(http.StatusCreated, gin.H{"data": counter})
}

func UpdateCounter(context *gin.Context) {
	var counter Counter
	counterId := context.Param("id")

	db := Database()
	db.First(&counter, counterId)

	context.Bind(&counter)

	if counter.Id == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"error": 
				"The counter with id equals " + counterId + " wasn't found"})

		return
	}

	db.Save(&counter)

	context.JSON(http.StatusOK, gin.H{"data": counter})
}

func DeleteCounter(context *gin.Context) {
	var counter Counter
	counterId := context.Param("id")

	db := Database()
	db.First(&counter, counterId)	

	if counter.Id == 0 {
		context.JSON(http.StatusNotFound, gin.H{
			"error": 
				"The counter with id equals " + counterId + " wasn't found"})

		return
	}	

	db.Delete(&counter)

	context.JSON(http.StatusNoContent, gin.H{
		"data":
			"The counter with id equals " + counterId + "was successfully deleted"})
}

func main() {
	router := gin.Default()
	db := Database()
	db.AutoMigrate(&Counter{})
	api := router.Group("api/")
	{
		api.GET("/counters", GetCounters)
		api.GET("/counters/:id", GetCounter)
		api.POST("/counters", CreateCounter)
		api.PUT("/counters/:id", UpdateCounter)
		api.DELETE("/counters/:id", DeleteCounter)
	}

	router.Run()
}
