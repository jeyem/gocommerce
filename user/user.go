package user

import (
	"time"

	"github.com/labstack/echo"

	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID                  bson.ObjectId `bson:"_id,omitempty"`
	Email               string        `bson:"email"`
	Sites               []Site        `bson:"sites"`
	Fullname            string        `bson:"fullname"`
	Password            string        `bson:"password"`
	OrginEmail          string        `bson:"orgin_email"`
	SecureKey           string        `bson:"secure_key"`
	SecureKeyExpireTime time.Time     `bson:"secure_key_expire_time"`
	Gender              string        `bson:"gender"`
	Birthday            time.Time     `bson:"birthday"`
	Call                string        `bson:"call"`
	Avatar              string        `bson:"avatar"`
	CreatedAt           time.Time     `bson:"created_at"`
	LastModified        time.Time     `bson:"last_modified"`
	LastLogin           time.Time     `bson:"last_login"`
	GoogleID            string        `bson:"google_id"`
	Keywords            []string      `bson:"keywords"`
}

type Site struct {
	Site        bson.ObjectId `bson:"site"`
	IsOwner     bool          `bson:"is_owner" `
	Permissions []string      `bson:"permissions"`
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

func (u *User) LoadWithIDAndSite(id, site bson.ObjectId) error {
	return db.Where(bson.M{
		"_id":  id,
		"site": site,
	}).Find(u)
}
func (u *User) LoadWithEmailAndSite(email string, site bson.ObjectId) error {
	return db.Where(bson.M{
		"email": email,
		"site":  site,
	}).Find(u)
}

func (u User) checkDuplicate() error {
	duplicateuser := new(User)
	if err := db.Where(bson.M{"$or": []bson.M{
		{"email": u.Email}, {"call": u.Call}},
	}); err == nil {
		return duplicateuser
	}
	return nil
}

func (u *User) Update() error {
	u.setKeywords()
	return db.Update(u)
}

func (u *User) LoadWithMail(email string) error {
	email = EmailFixer(email)
	return db.Where(bson.M{
		"email": email,
	}).Find(u)
}

func (u *User) LoadWithCellphone(cellphone string) error {
	return db.Where(bson.M{
		"call": cellphone,
	}).Find(u)
}

func (u *User) Auth(email, password string) error {
	if err := u.LoadWithMail(email); err != nil {
		return ErrorUserPass
	}
	if ok := CheckPassword(password, u.Password); !ok {
		return ErrorUserPass
	}
	u.LastLogin = time.Now()
	return u.Update()
}

func (u *User) AuthWithSite(email, password string, site bson.ObjectId) error {
	if err := u.LoadWithEmailAndSite(email, site); err != nil {
		return ErrorUserPass
	}
	if ok := CheckPassword(password, u.Password); !ok {
		return ErrorUserPass
	}
	u.LastLogin = time.Now()
	return u.Update()
}

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
