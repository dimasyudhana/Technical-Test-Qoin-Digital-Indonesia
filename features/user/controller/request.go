package controller

import "github.com/dimasyudhana/Qoin-Digital-Indonesia/features/user"

type registerRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type loginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func RequestToCore(data interface{}) user.UserCore {
	res := user.UserCore{}
	switch v := data.(type) {
	case registerRequest:
		res.Username = v.Username
		res.Email = v.Email
		res.Password = v.Password
	case loginRequest:
		res.Username = v.Username
		res.Password = v.Password
	default:
		return user.UserCore{}
	}
	return res
}
