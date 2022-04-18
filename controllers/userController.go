package controllers

import (
	"context"
	"github.com/eaaslan/go-jwt/database"
	"github.com/eaaslan/go-jwt/helpers"
	"github.com/eaaslan/go-jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		user := &models.User{}
		foundUser := &models.User{}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ok, err := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "password is not correct"})
			return
		}
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.UserID)
		helpers.UpdateAllToken(token, refreshToken, foundUser.UserID)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.UserID}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, foundUser)

	}
}

func VerifyPassword(userPass, foundUserPass string) (bool, error) {

	err := bcrypt.CompareHashAndPassword([]byte(userPass), []byte(foundUserPass))
	if err != nil {
		return false, err
	}
	return true, nil
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()
		user := &models.User{}
		if err := c.BindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": err})
			return
		}
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		//TODO try "if count>0" here too
		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, bson.M{"error ": "email or phone number already used"})
		}
		//TODO try different time format
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt = user.CreatedAt
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, *&user.UserID)
		user.Token = &token
		user.RefreshToken = &refreshToken
		//TODO try to insertone with &user
		insertedNum, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		defer cancel()
		c.JSON(http.StatusOK, insertedNum)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)

		user := &models.User{}
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		c.JSON(http.StatusOK, user)
	}
}
