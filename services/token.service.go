package services

import (
	db "Hdfc_Assignment/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateToken create a new token record
func CreateToken(mail string, id primitive.ObjectID, tokenType string, expiresAt time.Time) (*db.Token, error) {
	claims := &db.UserClaims{
		Email: mail,
		Type:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   id.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.JWTSecretKey))
	if err != nil {
		return nil, errors.New("cannot create access token")
	}

	tokenModel := db.NewToken(id, tokenString, tokenType, expiresAt)
	err = mgm.Coll(tokenModel).Create(tokenModel)
	if err != nil {
		return nil, errors.New("cannot save access token to db")
	}

	return tokenModel, nil
}

// DeleteTokenById delete token with id
func DeleteTokenById(tokenId primitive.ObjectID) error {
	ctx := mgm.Ctx()
	deleteResult, err := mgm.Coll(&db.Token{}).DeleteOne(ctx, bson.M{field.ID: tokenId})
	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete token")
	}

	return nil
}

func GenerateTokens(mail string, userId primitive.ObjectID) (*db.Token, *db.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
	refreshExpiresAt := time.Now().Add(time.Duration(Config.JWTRefreshExpirationDays) * time.Hour * 24)

	accessToken, err := CreateToken(mail, userId, db.TokenTypeAccess, accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := CreateToken(mail, userId, db.TokenTypeRefresh, refreshExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

// // GenerateAccessTokens generates "access" and "refresh" token for user
// func GenerateUserTokens(user *db.User) (*db.Token, *db.Token, error) {
// 	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
// 	refreshExpiresAt := time.Now().Add(time.Duration(Config.JWTRefreshExpirationDays) * time.Hour * 24)

// 	accessToken, err := CreateToken(user.Email,user.ID,db.TokenTypeAccess, accessExpiresAt)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	refreshToken, err := CreateToken(user.Email,user.ID,db.TokenTypeRefresh, refreshExpiresAt)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return accessToken, refreshToken, nil
// }

// // GenerateAccessTokens generates "access" and "refresh" token for user
// func GenerateAdminTokens(admin *db.Admin) (*db.Token, *db.Token, error) {
// 	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
// 	refreshExpiresAt := time.Now().Add(time.Duration(Config.JWTRefreshExpirationDays) * time.Hour * 24)

// 	accessToken, err := CreateToken(admin.Email,admin.ID,db.TokenTypeAccess, accessExpiresAt)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	refreshToken, err := CreateToken(admin.Email,admin.ID,db.TokenTypeRefresh, refreshExpiresAt)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return accessToken, refreshToken, nil
// }

// // GenerateAccessTokens generates "access" and "refresh" token for user
// func GenerateEmployeeTokens(employee *db.Employee) (*db.Token, *db.Token, error) {
// 	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
// 	refreshExpiresAt := time.Now().Add(time.Duration(Config.JWTRefreshExpirationDays) * time.Hour * 24)

// 	accessToken, err := CreateToken(employee.Email,employee.ID,db.TokenTypeAccess, accessExpiresAt)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	refreshToken, err := CreateToken(employee.Email,employee.ID,db.TokenTypeRefresh, refreshExpiresAt)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return accessToken, refreshToken, nil
// }

// VerifyToken checks jwt validity, expire date, blacklisted
func VerifyToken(token string, tokenType string) (*db.Token, error) {
	claims := &db.UserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWTSecretKey), nil
	})

	if err != nil || claims.Type != tokenType {
		return nil, errors.New("not valid token")
	}

	if time.Now().Sub(claims.ExpiresAt.Time) > 10*time.Second {
		return nil, errors.New("token is expired")
	}

	tokenModel := &db.Token{}
	userId, _ := primitive.ObjectIDFromHex(claims.Subject)
	err = mgm.Coll(tokenModel).First(
		bson.M{"type": tokenType, "user": userId, "blacklisted": false},
		tokenModel,
	)
	if err != nil {
		return nil, errors.New("cannot find token")
	}

	return tokenModel, nil
}
