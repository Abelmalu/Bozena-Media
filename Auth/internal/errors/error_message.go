package ierrors 

type ErrorMessage string
const (

MSGNameIsRequired ErrorMessage = "name is Required"
MSGUsernameIsRequired ErrorMessage = "Username is Required"
MSGPasswordIsRequired ErrorMessage = "password is Required"
MSGEmailIsRequired ErrorMessage = "email is Required"
MSGTimeoutReached ErrorMessage = "time out reached"
MSGRequestCanceled ErrorMessage = "request canceled"
MSGUsenameAlreadyExists ErrorMessage = "username already exists"
MSGEmailAlreadyExists ErrorMessage = "email already exists"
MSGSomethingWentWrong  ErrorMessage = "something went wrong"
MSGDatabaseError       ErrorMessage = "database error"
MSGNotFound            ErrorMessage = "not found"
MSGUnkownDevice		   ErrorMessage = "Unknown device type"
MSGFailedToValidateToken ErrorMessage = "unexpected signing method"
MSGRefreshTokenIsRequired ErrorMessage = "refresh token is required"
MSGUnauthorizedAccess ErrorMessage = "Unauthorized"
MSGUserNotFound 	  ErrorMessage = "Invalid username or password"
MSGBadRequest ErrorMessage = "invalid request"

)