package user

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/voice-of-colombo/service/database"
	"github.com/voice-of-colombo/service/graph/model"
	"golang.org/x/crypto/bcrypt"
)

const maxAge = 60 * 60 * 4

func validateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		jwtSecret := os.Getenv("JWT_SERCET")
		return []byte(jwtSecret), nil
	})
	return token, err
}

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := c.Cookie("Authorization")

		if err != nil || authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		const BEARER_SCHEMA = "Bearer "
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, _ := validateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println(claims)
			return
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func LogoutUser(c *gin.Context) {
	c.SetCookie("Authorization", "", 0, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}

func LoginUser(c *gin.Context) {
	email, password := c.PostForm("email"), c.PostForm("password")

	var user User
	result := database.Db.First(&user, "email = ?", email)
	if result.Error != nil {
		c.AbortWithError(401, result.Error)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		c.AbortWithError(401, errors.Wrapf(err, "Error comparing passwords"))
		return
	}

	// Clubs should always be set for a user
	clubs, err := GetClubsForUser(user.ID)
	if err != nil || len(clubs) == 0 {
		c.AbortWithError(500, errors.Wrapf(err, "Unable to load clubs for login user"))
		return
	}

	claimDetails := map[string]interface{}{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"id":        user.ID,
		"isAdmin":   user.IsAdmin,
		"clubs":     clubs,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claimDetails))

	// Sign and get the complete encoded token as a string using the secret
	jwtSecret := os.Getenv("JWT_SERCET")
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.SetCookie("Authorization", fmt.Sprintf("Bearer %s", tokenString), maxAge, "/", "", false, true)
	c.JSON(http.StatusOK, claimDetails)
}

func SearchUsers(searchCriteria model.UserSearchCriteria) ([]*model.User, error) {
	var users []*User
	queryBuider := database.Db.Limit(*searchCriteria.Limit).Offset(*searchCriteria.Offset)

	if searchCriteria.ID != nil {
		queryBuider = queryBuider.Where("id = ?", searchCriteria.ID)
	}

	queryBuider = queryBuider.Find(&users)

	if queryBuider.Error != nil {
		return nil, queryBuider.Error
	}
	return UserListToGqraphQL(users), nil
}

func GetClubsForUser(id string) ([]*model.Club, error) {
	// Todo: hardcoded limit and offset since unlikely
	var clubMappings []*ClubUserMapping
	mappingResult := database.Db.
		Limit(1000).
		Offset(0).
		Where("user_id = ?", id).
		Find(&clubMappings)

	if mappingResult.Error != nil {
		return nil, errors.Wrapf(mappingResult.Error, "Error loading mappings for user id %s", id)
	}

	// Get club Ids
	getClubIds := func() []string {
		clubIds := make([]string, len(clubMappings))
		for _, clubMapping := range clubMappings {
			clubIds = append(clubIds, clubMapping.ClubId)
		}
		return clubIds
	}

	clubIds := getClubIds()

	var clubs []*Club
	clubResult := database.Db.
		Model(&model.Club{}).
		Select("id", "club_name").
		Where("id IN ?", clubIds).
		Find(&clubs)

	if clubResult.Error != nil {
		return nil, errors.Wrapf(clubResult.Error, "Error loading clubs for clubIds %v", clubIds)
	}

	return ClubListToGqraphQL(clubs), nil

}
