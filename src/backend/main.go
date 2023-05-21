package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Record struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default()) // use default CORS middleware settings

	api, err := cloudflare.New(os.Getenv("CF_API_KEY"), os.Getenv("CF_API_EMAIL"))
	if err != nil {
		panic(err)
	}

	r.GET("/api/records/:zone", func(c *gin.Context) {
		zoneName := c.Param("zone")
		zoneID, err := api.ZoneIDByName(zoneName)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Assuming ResourceContainer is a struct that contains an Identifier
		zoneContainer := &cloudflare.ResourceContainer{
			Identifier: zoneID,
		}

		recFilter := cloudflare.ListDNSRecordsParams{}

		cfRecords, _, err := api.ListDNSRecords(context.Background(), zoneContainer, recFilter)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		records := make([]Record, len(cfRecords))
		for i, record := range cfRecords {
			records[i] = Record{ID: record.ID, Name: record.Name, Type: record.Type, Content: record.Content}
		}
		c.JSON(200, records)
	})

	r.GET("/api/types/:zone", func(c *gin.Context) {
		zoneName := c.Param("zone")
		zoneID, err := api.ZoneIDByName(zoneName)
		records, _, err := api.ListDNSRecords(context.Background(), &cloudflare.ResourceContainer{Identifier: zoneID}, cloudflare.ListDNSRecordsParams{})
		if err != nil {
			log.Println("Zone ID: ", zoneID)
			log.Fatal(err)
		}

		// Use a map to get unique record types
		typesMap := make(map[string]bool)
		for _, record := range records {
			typesMap[record.Type] = true
		}

		// Convert map keys to a slice
		types := make([]string, 0, len(typesMap))
		for t := range typesMap {
			types = append(types, t)
		}

		c.JSON(http.StatusOK, types)
	})

	r.DELETE("/api/record/:zone/:id", func(c *gin.Context) {
		zoneName := c.Param("zone")
		recordID := c.Param("id")
		zoneID, err := api.ZoneIDByName(zoneName)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Assuming ResourceContainer is a struct that contains an Identifier
		zoneContainer := &cloudflare.ResourceContainer{
			Identifier: zoneID,
		}

		err = api.DeleteDNSRecord(context.Background(), zoneContainer, recordID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "success"})
	})

	// This serves all static files in the build directory
	r.Static("/app/", "./build")

	r.NoRoute(func(c *gin.Context) {
		// If no route match, send back the main index.html file.
		c.File("./build/index.html")
	})

	r.Run()
}
