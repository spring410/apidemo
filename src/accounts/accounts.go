package accounts

type Account struct {
	ID         int    `json:"id"`
	Name       string `json:"name,omitempty"`
	Sex        int    `json:"sex,omitempty"`
	Age        int    `json:"age,omitempty"`
	Email      string `json:"email,omitempty"`
	Phone      string `json:"phone,omitempty"`
	CreateDate int64  `json:"createdate,omitempty"` //the number of seconds elapsed since January 1, 1970 UTC
}
