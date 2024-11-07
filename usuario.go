package controller

import (
	"fmt"

	"github.com/ivannadomz/Practica7_Usuarios_Go/database"
	"gorm.io/gorm"
)

// UserController defines available methods for user operations.
type UserController interface {
	CreateUser(user database.User) error
	GetUsers() ([]database.User, error)
	DeleteUser(userID uint) error
	UpdateUser(user database.User) error
}

// userController structure that holds the database connection.
type userController struct {
	db *gorm.DB
}

// Get all users.
func (c *userController) GetUsers() ([]database.User, error) {
	var users []database.User
	result := c.db.Find(&users)

	if result.Error != nil {
		fmt.Println("Error obteniendo usuario", result.Error)
		return nil, result.Error
	}
	fmt.Println(users)
	return users, nil
}

// Create a new user.
func (c *userController) CreateUser(user database.User) error {
	result := c.db.Create(&user)
	if result.Error != nil {
		fmt.Println("Error creando usuario", result.Error)
		return result.Error
	}
	return nil
}

// Delete a user.
func (c *userController) DeleteUser(ID uint) error {
	result := c.db.Delete(&database.User{}, ID)

	if result.Error != nil {
		fmt.Println("Error borrando usuario:", result.Error)
		return result.Error
	}

	fmt.Println("Usuario eliminado")
	return nil
}

// Update a user.
func (c *userController) UpdateUser(user database.User) error {
	// Find the user by ID and update the details.
	result := c.db.Model(&database.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		fmt.Println("Error actualizando usuario:", result.Error)
		return result.Error
	}

	fmt.Println("Usuario actualizado")
	return nil
}

// NewUserController creates a new user controller instance.
func NewUserController(db *gorm.DB) UserController {
	return &userController{db}
}
