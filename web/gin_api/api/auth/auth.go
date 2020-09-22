package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

type AccessToken struct {
	UUID  string `json:"uuid"`
	Token string `json:"token"`
}

type AccessDetails struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	AccessUUID  string `json:"access_uuid"`
	AccessToken string `json:"access_token"`
}

func GeneratePassword(plainText string, cost ...int) (string, error) {
	var rounds int
	if len(cost) == 0 {
		rounds = bcrypt.DefaultCost
	} else {
		rounds = cost[0]
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), rounds)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswords(inputPassword string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func GenerateToken(id int, email string) (AccessToken, error) {
	token := AccessToken{}

	tUUID, _ := uuid.NewV4()
	token.UUID = tUUID.String()
	aClaims := jwt.MapClaims{}
	aClaims["access_uuid"] = tUUID
	aClaims["user_id"] = id
	aClaims["user_email"] = email
	tWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, aClaims)
	tSigned, sErr := tWithClaims.SignedString([]byte(viper.GetString("JWT_ACCESS_SECRET")))
	if sErr != nil {
		return AccessToken{}, sErr
	}

	token.Token = tSigned
	return token, nil
}

func ValidateRequest(r *http.Request) error {
	token, err := ExtractTokenFromRequest(r)
	if err != nil {
		return err
	}
	return ValidateToken(token)
}

func ValidateToken(encodedToken string) error {
	token, err := VerifyToken(encodedToken)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return errors.New("token doesn't meet jwt claims")
	}
	return nil
}

// VerifyToken parses a jwt token, checks it's signature and returns it as a jwt.Token or an
// error if it isn't valid
func VerifyToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("JWT_ACCESS_SECRET")), nil
	})
}

func ExtractRequestTokenMetadata(r *http.Request) (AccessDetails, error) {
	tokenStr, err := ExtractTokenFromRequest(r)

	if err != nil {
		return AccessDetails{}, err
	}
	return ExtractTokenMetadata(tokenStr)
}

func ExtractTokenMetadata(encodedToken string) (AccessDetails, error) {
	token, err := VerifyToken(encodedToken)
	if err != nil {
		return AccessDetails{}, err
	}

	t := AccessDetails{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		t.AccessToken = encodedToken
		if t.AccessUUID, ok = claims["access_uuid"].(string); !ok {
			return AccessDetails{}, errors.New("invalid token")
		}
		if t.UserID, err = strconv.Atoi(fmt.Sprintf("%.f", claims["user_id"])); err != nil {
			return AccessDetails{}, errors.New("invalid token")
		}
	}
	return t, nil
}

func ExtractTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) == 2 {
		return tokenString[1], nil
	}
	return "", errors.New("no authorization token found")
}
