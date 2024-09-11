package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
}

type Password struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	StoredPassword string `json:"stored_password"`
}
type Pass struct {
	Password string `json:"password" validate:"required,min=8"`
}

type Parameter struct {
	Have8Char     bool `json:"have8char"`
	HaveNum       bool `json:"have_num"`
	SmallLetter   bool `json:"small_letter"`
	CapitalLetter bool `json:"capital_letter"`
	SpecialChar   bool `json:"special_char"`
	SuperStrong   bool `json:"super_strong"`
	Strong        bool `json:"strong"`
}
