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
	if err := authRepo.DB.QueryRowContext(ctx,query, user.Name, user.Username, user.Email, user.Password).Scan(&newUser.ID, &newUser.Role); err != nil {
		// Change *pq.Error to *pgconn.PgError
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return nil,pgErr
		}
		log.Printf("Registration DB Error %v", err)

		return nil,err

	}

	return &newUser,nil
}
func (authrepo *AuthRepository) Login(ctx context.Context,userName,password string)(*model.User,error){
	var user model.User
	query := `SELECT * FROM users WHERE username=$1`
	err := authrepo.DB.QueryRowContext(ctx,query, userName).Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role)

	if err != nil {

		log.Printf("Login DB Error %v", err)

		return nil,err
		
	}


  return &user,nil

}
func (authRepo *AuthRepository) Logout(ctx context.Context, tokenID string) ( error){

	

		query := `DELETE FROM refresh_tokens WHERE token_text=$1`
	
		result, err := authRepo.DB.ExecContext(ctx,query, tokenID)
	
		if err != nil {
	
			log.Fatalf("DB Exec error %v", err)
			return err
		}
	
		_, err = result.RowsAffected()
		if err != nil {
	
			log.Fatalf("db exec error %v", err)
			return err
		}
		return err
	
	


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


func (authRepo *AuthRepository)  RevokeRefreshToken(refreshToken string) error {

	query := `
	
	UPDATE refresh_tokens SET revoked=TRUE 
	WHERE token_text = $1 AND revoked=FALSE `

	result, err := authRepo.DB.Exec(query, refreshToken)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// detect reuse attempt
	if rowsAffected == 0 {
		// token was already revoked or doesn't exist
		log.Printf("refresh token already revoked or not found")
	}

	return err

}

func (authRepo *AuthRepository) GetRefreshToken(refreshToken string) (*model.RefreshToken, error) {

	var refreshRecord model.RefreshToken

	// hashing the token because stored tokens are hashed
	hashedrefreshToken := pkg.HashToken(refreshToken)

	query := `SELECT * FROM refresh_tokens where token_text = $1;`

	if err := authRepo.DB.QueryRow(query, hashedrefreshToken).Scan(&refreshRecord); err != nil {

		return nil, err
	}

	return &refreshRecord, nil
}


func (authRepo *AuthRepository)  GetUserByID(ID int) (*model.User, error) {
	var user model.User
	query := `SELECT * FROM users WHERE id=$1`

	err := authRepo.DB.QueryRow(query, ID).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}