package pkg



import (
    "github.com/golang-jwt/jwt/v5"
    "time"
)

var jwtSecret = []byte("your_secret_key")

func GenerateAcessToken(userID int) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "type":"access",
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//validating the signing method is HS256 / HS384 / HS512,
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrTokenSignatureInvalid
        }
        return jwtSecret, nil 
    })
}

func GenerateRefreshToken(userID int)(string,error){

    claims := jwt.MapClaims{
        "user_id":userID,
        "type":"refresh",
        "exp": time.Now().Add( 24*30 * time.Hour ).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

    return token.SignedString(jwtSecret)


}
