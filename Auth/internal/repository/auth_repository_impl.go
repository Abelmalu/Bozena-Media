package repository

import (
	"context"
	"database/sql"
	"log"

	model "github.com/abelmalu/golang-posts/Auth/internal/models"
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

func (authRepo *AuthRepository) Register(ctx context.Context,user model.User)(*model.User,error){
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