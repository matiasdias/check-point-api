package security

import "golang.org/x/crypto/bcrypt"

func Hash(passWord string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(passWord), bcrypt.DefaultCost)
}

func VerifyPassWord(passWordHash, passWord string) error {
	return bcrypt.CompareHashAndPassword([]byte(passWordHash), []byte(passWord))
}
