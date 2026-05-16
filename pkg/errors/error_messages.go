package ierrors
 

type ErrorMessage string

const (
MSGInvalidRequestBody ErrorMessage = "Invalid Request Body"
MSGInternalServerError ErrorMessage = "Internal Server Error"
MSGUnauthorizedAccess ErrorMessage = "Unauthorized"
MSGSomethingWentWrong ErrorMessage = "Something Went Wrong"
MSGRefreshTokenNotFound ErrorMessage = "Refresh token not found"
MSGUsernameIsRequired ErrorMessage = "Username is Required"



)