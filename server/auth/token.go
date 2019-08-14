package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type authClaims struct {
	email string
	jwt.StandardClaims
}

type AuthToken struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

func NewAuth(secret string, alg string) *AuthToken {
	return &AuthToken{secret: []byte(secret), signingMethod: jwt.GetSigningMethod(alg)}
}

func (a *AuthToken) sign(email string) *jwt.Token {
	exp := time.Now().Add(24 * time.Hour).Unix()
	return jwt.NewWithClaims(a.signingMethod, authClaims{email, jwt.StandardClaims{ExpiresAt: exp}})
}

func (a *AuthToken) parseClaims(tokenString string) (*authClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return a.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*authClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("unexpected claims")
	}
}

func (a *AuthToken) TokenStringForUser(email string) (string, error) {
	return a.sign(email).SignedString(a.secret)
}

func (a *AuthToken) UserFromTokenString(tokenString string) (string, error) {
	claims, err := a.parseClaims(tokenString)
	if err != nil {
		return "", err
	} else {
		return claims.email, nil
	}
}
