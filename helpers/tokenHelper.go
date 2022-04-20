package helpers

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/eaaslan/go-jwt/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	UserType  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, userType, uid string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		UID:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Fatalln(err)
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Fatalln(err)
	}

	return token, refreshToken, nil
}

func UpdateAllToken(signedToken, signedRefreshToken, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	updatedObj := primitive.D{}
	updatedObj = append(updatedObj, bson.E{Key: "token", Value: signedToken})
	updatedObj = append(updatedObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: UpdatedAt})
	upsert := true
	//TODO why upsert *boolean ?
	filter := bson.M{"user_id": userID}
	opt := options.UpdateOptions{Upsert: &upsert}
	_, err := userCollection.UpdateOne(ctx, filter, bson.D{
		{"$set", updatedObj},
	},
		&opt,
	)
	defer cancel()
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		err = errors.New("the token is invalid")
		return nil, err
	}
	return claims, nil
}
