package ierrors 

type ErrorMessage string
const (

MSGNameIsRequired ErrorMessage = "name is Required"
MSGUsernameIsRequired ErrorMessage = "Username is Required"
MSGPasswordIsRequired ErrorMessage = "password is Required"
MSGEmailIsRequired ErrorMessage = "email is Required"
MSGTimeoutReached ErrorMessage = "time out reached"
MSGRequestCancelled ErrorMessage = "request cancelled"
MSGUsenameAlreadyExists ErrorMessage = "username already exists"
MSGEmailAlreadyExists ErrorMessage = "email already exists"
MSGSomethingWentWrong  ErrorMessage = "something went wrong"
MSGDatabaseError       ErrorMessage = "database error"
MSGNotFound            ErrorMessage = "not found"




)