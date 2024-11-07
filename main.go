package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	controller "github.com/ivannadomz/Practica7_Usuarios_Go/controler"
	"github.com/ivannadomz/Practica7_Usuarios_Go/database"
)

func main() {
	// Connect to the database
	db, err := database.NewDataBaseDriver() // Cambiado a NewDataBaseDriver
	if err != nil {
		fmt.Println("Database connection error")
		return
	}

	// Migrate database structure
	db.AutoMigrate(&database.User{})

	// Create a user controller
	userCtrl := controller.NewUserController(db)

	// Add some initial users
	addInitialUsers(userCtrl)

	// Get users
	users, err := userCtrl.GetUsers()
	if err != nil {
		fmt.Println("Error fetching users: ", err)
		return
	}

	// Update the third user if they exist
	if len(users) > 2 {
		updateUser(userCtrl, users[2])
	}

	// Delete the third user if they exist
	if len(users) > 2 {
		deleteUser(userCtrl, users[2].ID)
	}

	// Start the web server
	startServer(userCtrl)
}

func addInitialUsers(ctrl controller.UserController) {
	ctrl.CreateUser(database.User{
		Name:  "Israel Olivares",
		Email: "israelf04@gmail.com",
	})
	ctrl.CreateUser(database.User{
		Name:  "Natalia Castellano",
		Email: "nati24rod@gmail.com",
	})
	ctrl.CreateUser(database.User{
		Name:  "Sebastian Castro",
		Email: "sebbarbas@gmaiñ",
	})
	ctrl.CreateUser(database.User{
		Name:  "Melissa Ruiz",
		Email: "meliruiz04@gmail.com",
	})
}

func updateUser(ctrl controller.UserController, user database.User) {
	user.Name = "Mauricio Tellez"
	user.Email = "mauel45@gmail.com"
	err := ctrl.UpdateUser(user)
	if err != nil {
		fmt.Println("Error actualizando al usuario ", err)
	} else {
		fmt.Println("usuario actualizado")
	}
}

func deleteUser(ctrl controller.UserController, id uint) {
	err := ctrl.DeleteUser(id)
	if err != nil {
		fmt.Println("Error eliminando al usuario ", err)
	} else {
		fmt.Println("Usuario eliminado")
	}
}

func startServer(ctrl controller.UserController) {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/", func(c *gin.Context) {
		users, err := ctrl.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":      "Home",
			"user_count": len(users),
			"user_list":  users,
		})
	})

	router.GET("/api/users", func(c *gin.Context) {
		users, err := ctrl.GetUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener usuarios"})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	router.POST("/api/users", func(c *gin.Context) {
		var user database.User
		if err := c.ShouldBindJSON(&user); err == nil {
			if err := ctrl.CreateUser(user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando usuario"})
				return
			}
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalido payload"})
		}
	})

	router.DELETE("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalido ID"})
			return
		}

		if err := ctrl.DeleteUser(uint(idParsed)); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado"})
	})

	router.PUT("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		var user database.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		user.ID = uint(idParsed) // Assign ID to the user
		if err := ctrl.UpdateUser(user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	router.Run(":8001")
}
