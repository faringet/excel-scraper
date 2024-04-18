package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/receiveEmployee", receiveEmployeeHandler)

	if err := router.Run("localhost:3004"); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}

func receiveEmployeeHandler(c *gin.Context) {
	var employeeData interface{}
	if err := c.BindJSON(&employeeData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	switch data := employeeData.(type) {
	case []interface{}:
		fmt.Println("Received employee data:")
		for _, obj := range data {
			printEmployee(obj)
		}
	default:
		fmt.Println("Received employee data:")
		printEmployee(data)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data received successfully"})
}

func printEmployee(data interface{}) {
	employeeData, ok := data.(map[string]interface{})
	if !ok {
		fmt.Println("Invalid employee data")
		return
	}

	for key, value := range employeeData {
		fmt.Printf("%s: %v\n", key, value)
	}
	fmt.Println()
}
