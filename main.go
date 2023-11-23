package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type room struct {
	ID        string  `json:"id"`
	Class     string  `json:"class"`
	Capacity  int     `json:"capacity"`
	Facility  string  `json:"facility"`
	Price     float64 `json:"price"`
	RoomCount int     `json:"roomcount"`
	Available bool    `json:"available"`
}

type reservationPayload struct {
	RoomReserved int `json:"roomreserved"`
}

var rooms = []room{
	{
		ID:        "1",
		Class:     "MegaVip",
		Capacity:  10,
		Facility:  "Double Bed King size, Free Food, Big Luxury Bathub, Swimming Pool, Gym, Custom Request",
		Price:     25000000,
		RoomCount: 10,
		Available: true,
	},
	{
		ID:        "2",
		Class:     "Vip",
		Capacity:  20,
		Facility:  "Double Bed King size, Free Food, Big Luxury Bathub,Custom Request Charged",
		Price:     17000000,
		RoomCount: 15,
		Available: true,
	},
	{
		ID:        "3",
		Class:     "Golde",
		Capacity:  30,
		Facility:  "Double Bed King size, Free Food, Medium Luxury Bathub",
		Price:     15000000,
		RoomCount: 40,
		Available: true,
	},
	{
		ID:        "4",
		Class:     "Silver",
		Capacity:  7,
		Facility:  "Medium Bed size, Small Luxury Bathub",
		Price:     10000000,
		RoomCount: 50,
		Available: true,
	},
	{
		ID:        "5",
		Class:     "Bronze",
		Capacity:  4,
		Facility:  "Medium Bed size, Free Morning Meal",
		Price:     2000000,
		RoomCount: 60,
		Available: true,
	},
}

func main() {
	router := gin.Default()
	router.GET("/rooms", getRooms)
	router.GET("/rooms/:id", getRoomById)
	router.POST("/rooms", postRooms)
	router.PUT("/rooms/:id", editRoomById)
	router.DELETE("/rooms/:id", deleteRoomById)
	router.POST("/uploads-room", uploadRoomImage)

	router.POST("/reserve/:id", reserveRoom)
	router.Run("localhost:8083")

}

func getRooms(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, rooms)
}

func postRooms(c *gin.Context) {
	var newRoom room

	if err := c.BindJSON(&newRoom); err != nil {
		return
	}

	rooms = append(rooms, newRoom)
	c.IndentedJSON(http.StatusCreated, newRoom)
}

func getRoomById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range rooms {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})

}

func editRoomById(c *gin.Context) {
	id := c.Param("id")

	var updatedRoom room

	if err := c.BindJSON(&updatedRoom); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid Json provided"})
		return
	}

	for i, r := range rooms {
		if r.ID == id {
			rooms[i].Class = updatedRoom.Class
			rooms[i].Capacity = updatedRoom.Capacity
			rooms[i].Facility = updatedRoom.Facility
			rooms[i].Price = updatedRoom.Price
			rooms[i].RoomCount = updatedRoom.RoomCount
			rooms[i].Available = updatedRoom.Available

			c.IndentedJSON(http.StatusOK, rooms[i])
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})

}

func deleteRoomById(c *gin.Context) {
	id := c.Param("id")

	for i, r := range rooms {
		if r.ID == id {
			rooms = append(rooms[:i], rooms[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "Room Deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})
}

func reserveRoom(c *gin.Context) {
	id := c.Param("id")
	var roomReserved reservationPayload

	if err := c.BindJSON(&roomReserved); err != nil {
		fmt.Println("Error binding JSON:", err)
		fmt.Println(&roomReserved)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON provided"})
		return
	}

	fmt.Println(roomReserved)

	if roomReserved.RoomReserved < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "cant reserve room in 0 room"})
		return
	}

	for i, r := range rooms {
		if r.ID == id {
			if rooms[i].RoomCount < roomReserved.RoomReserved {
				rooms[i].Available = false
				c.IndentedJSON(http.StatusForbidden, gin.H{"message": "cant reserve over the limit"})
				return
			}
			rooms[i].RoomCount -= roomReserved.RoomReserved
			message := fmt.Sprintf("Room reserved %d", roomReserved)
			c.IndentedJSON(http.StatusOK, gin.H{"message": message})

			return

		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})

}

func uploadRoomImage(c *gin.Context) {

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing form"})
		return
	}
	defer file.Close()

	fmt.Printf("Received file: %+v\n", handler.Filename)

	// Create the "uploads" directory if it doesn't exist
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating upload directory"})
			return
		}
	}

	// Construct the path to save the file
	filePath := filepath.Join(uploadDir, handler.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating file"})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying file"})
		return
	}

	fmt.Println("File saved successfully")
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
