package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	model "github.com/abelmalu/golang-posts/Auth/internal/models"
	"github.com/abelmalu/golang-posts/pkg"
	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository{

	return &AuthRepository{
		DB:db,
	}
}

func (authRepo *AuthRepository) Register(ctx context.Context,user *model.User)(*model.User,error){
	var newUser model.User


	query := `INSERT INTO users(name,username,email,password) VALUES($1,$2,$3,$4) RETURNING id,role`
	if err := authRepo.DB.QueryRow(query, user.Name, user.Username, user.Email, user.Password).Scan(&newUser.ID, &newUser.Role); err != nil {
		// Change *pq.Error to *pgconn.PgError
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return nil,pgErr
		}
		log.Printf("Registration DB Error %v", err)

		return nil,err

	}

	return &newUser,nil
}

func (authRepo *AuthRepository)StoreRefreshTokens(userID int, refreshToken string, expiresAt time.Time, clientType string) (sql.Result, error) {

	// hashing the token before inserting to a db
	refreshToken = pkg.HashToken(refreshToken)

	query := `INSERT INTO refresh_tokens (user_id,token_text,expires_at,client_type) VALUES($1,$2,$3,$4)`

	result, err := authRepo.DB.Exec(query, userID, refreshToken, expiresAt, clientType)
	if err != nil {

		return nil, err
	}

	return result, nil

}