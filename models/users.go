package models

const (
	STUDENT int = 0
	TEACHER int = 1
)

func ValidateUserType(userType int) string {
	if userType == TEACHER {
		return "teacher"
	} else if userType == STUDENT {
		return "student"
	} else {
		return ""
	}
}

type UserCreateSchema struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Type      string `json:"type"`
}

type UserResponseSchema struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Type      string `json:"type"`
}

func (data UserCreateSchema) Insert(uuidString string) {

}

func (data UserCreateSchema) Update(uuidString string) {

}

func (data UserCreateSchema) ConvertToMap(uuidString string) {

}
