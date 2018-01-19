package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	method     = "pbkdf2:sha256"
	saltLength = 8
	iterations = 50000
	charLength = 64
	saltChars  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func MakePassword(password string) string {
	salt := genSalt()
	hash := hashInternal(salt, password)
	return fmt.Sprintf("pbkdf2:sha256:%v$%s$%s", iterations, salt, hash)
}

func CheckPassword(password string, hash string) bool {
	if strings.Count(hash, "$") < 2 {
		return false
	}
	pwdHashList := strings.Split(hash, "$")
	return pwdHashList[2] == hashInternal(pwdHashList[1], password)
}

func EmailFixer(email string) string {
	email = strings.ToLower(email)
	if strings.Index(email, "gmail.com") < 0 {
		return email
	}
	splited := strings.Split(email, "@")
	if len(splited) < 2 {
		return email
	}
	username := splited[0]
	domain := splited[1]
	username = strings.Replace(username, ".", "", -1)
	return username + "@" + domain
}

func genSalt() string {
	var bytes = make([]byte, saltLength)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = saltChars[v%byte(len(saltChars))]
	}
	return string(bytes)
}

func hashInternal(salt string, password string) string {
	hash := pbkdf2.Key([]byte(password), []byte(salt), iterations, charLength/2, sha256.New)
	return hex.EncodeToString(hash)
}
