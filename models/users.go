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
