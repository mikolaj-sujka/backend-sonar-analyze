package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"go_app/database"
	"go_app/models"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type UserController struct {
	Store *database.UserStore
}

func NewUserController(store *database.UserStore) *UserController {
	return &UserController{
		Store: store,
	}
}


var jwtKey = []byte("your_secret_key") 
func (uc *UserController) RegisterUser(c echo.Context) error {
    u := new(models.User)
    if err := c.Bind(u); err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid data provided"})
    }

    u.Email = strings.ToLower(u.Email) // Normalize email to lower case
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password"})
    }

    u.Password = string(hashedPassword)
    u.ID = uuid.New().String()

    uc.Store.AddUser(u)

    return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully", "user": u})
}

func (uc *UserController) LoginUser(c echo.Context) error {
    credentials := struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }{}

    if err := c.Bind(&credentials); err != nil {
        fmt.Println("Error binding request:", err)  
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
    }

    fmt.Println("Received credentials:", credentials) 

    emailLower := strings.ToLower(credentials.Email)
    user, exists := uc.Store.GetUserByEmail(emailLower)
    if !exists {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Email not found"})
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Password does not match"})
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &jwt.StandardClaims{
        Subject:   user.ID,
        ExpiresAt: expirationTime.Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create JWT token"})
    }

    return c.JSON(http.StatusOK, echo.Map{
        "message": "User logged in successfully",
        "token":   tokenString,
        "user": map[string]interface{}{
            "ID":       user.ID,
            "Username": user.Username,
            "Email":    user.Email,
        },
    })
}

func (uc *UserController) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, exists := uc.Store.GetUser(id)
	if !exists {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, echo.Map{"user": user})
}

func (uc *UserController) GoogleLogin(c echo.Context) error {
	url := oauthConf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (uc *UserController) GoogleCallback(c echo.Context) error {
	state := c.QueryParam("state")
    if state != "state-token" {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "State mismatch"})
    }

    code := c.QueryParam("code")
    if code == "" {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Authorization code not provided"})
    }

    oauthToken, err := oauthConf.Exchange(context.Background(), code)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Failed to exchange code for token", "detail": err.Error()})
    }

	client := oauthConf.Client(context.Background(), oauthToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Failed to retrieve user info"})
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to decode user info"})
	}

	var user *models.User
	existingUser, exists := uc.Store.GetUserByEmail(userInfo.Email)
	if exists {
		user = existingUser
	} else {
		user = &models.User{
			ID:       uuid.New().String(),
			Username: userInfo.Name,
			Email:    userInfo.Email,
			Password: "",
		}
		uc.Store.AddUser(user)
	}

	// Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   user.ID,
		ExpiresAt: expirationTime.Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(jwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create JWT token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "User logged in successfully",
		"token":   tokenString,
	})
}


var oauthConf = &oauth2.Config{}