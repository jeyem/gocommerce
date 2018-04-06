package user

import (
	"github.com/asaskevich/govalidator"
)

type Form struct {
	Fullname string `form:"fullname" json:"fullname"`
	Call     string `form:"call" json:"call"`
	Email    string `form:"email" json:"email" valid:"required"`
	Password string `form:"password" json:"password" valid:"required"`
}

type CallForm struct {
	Call string `form:"call" json:"call" valid:"required"`
}

type MailForm struct {
	Email string `form:"email" json:"email" valid:"required"`
}

type SecureKeyForm struct {
	Identifier string `form:"identifier" json:"identifier" valid:"required"`
	Key        string `form:"key" json:"key" valid:"required"`
}

func (f Form) Login() (*User, error) {
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return nil, err
	}
	u := new(User)
	if err := u.Auth(f.Email, f.Password); err != nil {
		return nil, err
	}
	return u, nil

}

func (f Form) Register() (*User, error) {
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

func (f CallForm) StepsLogin() (*User, error) {
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return nil, err
	}
	u := new(User)
	if err := u.LoadByCall(f.Call); err != nil {
		u.Call = f.Call
	}
	u.GenerateSecureKey()
	if err := u.Save(); err != nil {
		return nil, err
	}
	return u, nil
}

func (f MailForm) StepsLogin() (*User, error) {
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return nil, err
	}
	u := new(User)
	f.Email = EmailFixer(f.Email)
	if err := u.LoadByMail(f.Email); err != nil {
		return nil, err
	}
	u.GenerateSecureKey()
	if err := u.Save(); err != nil {
		return nil, err
	}
	return u, nil
}

func (f SecureKeyForm) Login() (*User, error) {
	if _, err := govalidator.ValidateStruct(f); err != nil {
		return nil, err
	}
	u = new(User)
	if err := u.AuthByCallAndSecureKey(f.Identifier, f.Key); err == nil {
		return u, nil
	}
	f.Identifier = EmailFixer(f.Identifier)
	if err := u.AuthByMailAndSecureKey(f.Identifier, f.Key); err != nil {
		return nil, err
	}
	return u, nil
}
