package schemas

type UserRegister struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EditRole struct {
	Role string `json:"role"`
}

type EditUserData struct {
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
