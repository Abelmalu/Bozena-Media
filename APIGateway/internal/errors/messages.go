package errors 

type ErrorMessage string
const (
MSGInvalidRequestBody ErrorMessage = "Invalid Request Body"
MSGInternalServerError ErrorMessage = "Internal Server Error"
MSGUnauthorizedAccess ErrorMessage = "Unauthorized"



)