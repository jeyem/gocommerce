package user

import (
	"time"

	"github.com/jeyem/gocommerce/util/random"

	"github.com/labstack/echo"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID                  bson.ObjectId `bson:"_id,omitempty"`
	Email               string        `bson:"email"`
	Fullname            string        `bson:"fullname"`
	Password            string        `bson:"password"`
	OrginEmail          string        `bson:"orgin_email"`
	SecureKey           string        `bson:"secure_key"`
	SecureKeyExpireTime time.Time     `bson:"secure_key_expire_time"`
	Gender              string        `bson:"gender"`
	Birthday            time.Time     `bson:"birthday"`
	Role                string        `bson:"role"`
	Call                string        `bson:"call"`
	Avatar              string        `bson:"avatar"`
	CreatedAt           time.Time     `bson:"created_at"`
	LastModified        time.Time     `bson:"last_modified"`
	LastLogin           time.Time     `bson:"last_login"`
	GoogleID            string        `bson:"google_id"`
	Keywords            []string      `bson:"keywords"`
}

func (u *User) Update() error {
	u.setKeywords()
	return db.Update(u)
}

func (u *User) Save() error {
	u.LastModified = time.Now()
	if u.ID.Valid() {
		return u.Update()
	}
	if err := u.checkDuplicate(); err != nil {
		return err
	}
	u.CreatedAt = time.Now()
	u.setKeywords()
	return db.Create(u)
}
func (u *User) setKeywords() {
	u.Keywords = []string{u.OrginEmail, u.Email, u.Fullname, u.Call}
}

func (u *User) Load(id bson.ObjectId) error {
	return db.Get(u, id)
}

func (u *User) LoadByMail(email string) error {
	return db.Where(bson.M{
		"email": email,
	}).Find(u)
}

func (u User) checkDuplicate() error {
	duplicateuser := new(User)
	if err := db.Where(bson.M{"$or": []bson.M{
		{"email": u.Email}, {"call": u.Call}},
	}); err == nil {
		return ErrorDuplicateUser
	}
	return nil
}

func (u *User) LoadByCall(cellphone string) error {
	return db.Where(bson.M{
		"call": cellphone,
	}).Find(u)
}

func (u *User) AuthByMail(email, password string) error {
	if err := u.LoadByMail(email); err != nil {
		return ErrorUserPass
	}
	if ok := CheckPassword(password, u.Password); !ok {
		return ErrorUserPass
	}
	u.LastLogin = time.Now()
	return u.Update()
}

func (u *User) AuthByCall(call, password string) error {
	if err := u.LoadbyCall(call); err != nil {
		return ErrorUserPass
	}
	if ok := CheckPassword(password, u.Password); !ok {
		return ErrorUserPass
	}
	u.LastLogin = time.Now()
	return u.Update()
}

func (u *User) AuthByCallAndSecureKey(call, key string) error {
	if err := db.Where(bson.M{
		"call":                   call,
		"secure_key":             key,
		"secure_key_expire_time": bson.M{"gt": time.Now()},
	}).Find(u); err != nil {
		return err
	}
	u.LastLogin = time.Now()
	return u.Update()
}

func (u *User) AuthByMailAndSecureKey(email, key string) error {
	if err := db.Where(bson.M{
		"email":                  email,
		"secure_key":             key,
		"secure_key_expire_time": bson.M{"gt": time.Now()},
	}).Find(u); err != nil {
		return err
	}
	u.LastLogin = time.Now()
	return u.Update()
}

func (u *User) GenerateSecureKey() string {
	u.SecureKey = random.Rand(5)
	u.SecureKeyExpireTime = time.Now().Add(time.Hour)
	return u.SecureKey
}

// minimal map responses
func (u *User) Rest() echo.Map {
	email := u.OrginEmail
	return echo.Map{
		"email":    email,
		"fullname": u.Fullname,
		"call":     u.Call,
		"avatar":   u.Avatar,
	}
}

func (u User) Temp() echo.Map {
	return echo.Map{
		"email":    u.OrginEmail,
		"fullname": u.Fullname,
		"call":     u.Call,
		"avatar":   u.Avatar,
	}
}
