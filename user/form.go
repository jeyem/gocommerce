package user

import (
	"github.com/asaskevich/govalidator"
)

type Form struct {
	Fullname string `form:"fullname" json:"fullname"`
	Email    string `form:"email" json:"email" valid:"required"`
	Password string `form:"password" json:"password" valid:"required"`
}

func (f Form) Load() (*User, error) {
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return nil, err
	}
	u := new(User)
	if err := u.Auth(f.Email, f.Password); err != nil {
		return nil, err
	}
	return u, nil

}

func (f Form) Create() (*User, error) {
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return nil, err
	}
	u := &User{
		Fullname:   f.Fullname,
		Email:      EmailFixer(f.Email),
		OrginEmail: f.Email,
		Password:   MakePassword(f.Password),
	}
	if err := u.Save(); err != nil {
		return nil, err
	}
	return u, nil

}
