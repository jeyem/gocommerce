package user

import (
	"errors"

	"gopkg.in/mgo.v2/bson"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	Fullname string `form:"fullname" json:"fullname"`
	Email    string `form:"email" json:"email" valid:"required"`
	Password string `form:"password" json:"password" valid:"required"`
}

type ForgetPasswordForm struct {
	Email     string `form:"email" json:"email"`
	Cellphone string `form:"cellphone" json:"cellphone"`
}

type ResetPasswordForm struct {
	Password        string `form:"password" json:"password" valid:"required" `
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" valid:"required"`
}

type ForgetPasswordKeyForm struct {
	Key string `form:"key" json:"key"`
}

func (f Form) Load(site bson.ObjectId) (*User, error) {
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return nil, err
	}
	u := new(User)
	if id.Valid() {
		return u, u.AuthWithSite(f.Email, f.Password, site)
	}
	return u, u.Auth(f.Email, f.Password)

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
	return u, u.Save()

}

func (f ResetPasswordForm) Validate() error {
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return err
	}

	if f.Password != f.ConfirmPassword {
		return errors.New("Passwords do not match")
	}

	return nil

}
