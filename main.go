package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	/**
	for insert database this is several step must be create
	this step creates from under to upper
	5. user input form
	4. create handler : create struct input from user input
	3. create service : create mapping from input struct to user struct
	2. create repository (user/repository.go)
	1. call database (main.go)
	*/

	//create connection from database
	dsn := "root:@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService() //call jwt service
	userHandler := handler.NewUserHandler(userService, authService)

	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images", "./images") // this route for access avatar endpoint. '/images' is path endpoint, mean while './images' is name of folder
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)                    //API register
	api.POST("/sessions", userHandler.Login)                        //API session
	api.POST("/email_checkers", userHandler.CheckEmailAvailability) //API for check availability email

	api.GET("/campaigns", campaignHandler.GetCampaigns)    //API for get campaigns
	api.GET("/campaigns/:id", campaignHandler.GetCampaign) //APi for detail campaign, :id is URI

	/**
	authMiddleware for 'validate jwt token'
	authMiddleware(auth.service, user.service) : we just parse auth middleware
	authMiddleware() : we get return from authMiddleware function
	in the middleware me must 'Parsing' function, not get return value
	*/
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar) //API for upload avatar image

	router.Run()
}

/**
step by step middleware :
'this middleware for validate jwt token'
1. get value authorization. example: bearer tokentokentoken
2. from header authorization, we take just token value
3. we validate token
4, if token valid we get user_id
5. get user from db appropriate 'user_id' through service
6. if user available we set context with 'user'
*/
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") { //if in header not contain 'Bearer'
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			/**
			using abort with status json, because for reject middleware will call in middle
			example: api.POST("/avatars",<MIDDLEWARE>, userHandler.UploadAvatar)
			*/
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// value of token is 'Bearer <RANDOM_TOKEN>'
		var tokenString string = ""
		tokenArray := strings.Split(authHeader, " ")
		if len(tokenArray) == 2 {
			tokenString = tokenArray[1]
		}

		// auth service
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		// userservice
		user, err := userService.GetUserById(userID)

		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user) //current user meaning user is logged in or access this application
	}
}
