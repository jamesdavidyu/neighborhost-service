package auth

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/db"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/config"
	"github.com/jamesdavidyu/neighborhost-service/utils"
)

type contextKey string

const NeighborKey contextKey = "neighborId"

func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenFromRequest(r)

		token, err := validateToken(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["neighborId"].(string)

		neighborId, err := strconv.Atoi(str)
		if err != nil {
			permissionDenied(w)
			return
		}

		neighbor, err := GetNeighborById(neighborId)
		if err != nil {
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, NeighborKey, neighbor.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, neighborId int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"neighborId": strconv.Itoa(neighborId),
		"expiredAt":  time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func scanRowIntoNeighbor(rows *sql.Rows) (*types.Neighbors, error) {
	neighbor := new(types.Neighbors)

	err := rows.Scan(&neighbor.ID, &neighbor.Email, &neighbor.Username, &neighbor.Zipcode, &neighbor.Verified, &neighbor.NeighborhoodID)
	if err != nil {
		return nil, err
	}

	return neighbor, nil
}

func GetNeighborById(id int) (*types.Neighbors, error) {
	db, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, email, username, zipcode, verified, neighborhood_id FROM neighbors WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	neighbor := new(types.Neighbors)
	for rows.Next() {
		neighbor, err = scanRowIntoNeighbor(rows)
		if err != nil {
			return nil, err
		}
	}

	if neighbor.ID == 0 {
		return nil, fmt.Errorf("error")
	}

	return neighbor, nil
}

func GetNeighborIDFromContext(ctx context.Context) int {
	neighborId, ok := ctx.Value(NeighborKey).(int)
	if !ok {
		return -1
	}

	return neighborId
}
