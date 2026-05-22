package dto



type ToggleLikeRequest struct {

	State bool `json:"state"`
}


type ToggleLikeResponse struct {

	Message string `json:"message"`
}