package security

import "golang.org/x/crypto/bcrypt"

// Hash responsável por criptografar a senha
func Hash(passWord string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(passWord), bcrypt.DefaultCost)
}

// VerifyPassWord responsǽvel por verificar se a senha é correta
func VerifyPassWord(passWordHash, passWord string) error {
	return bcrypt.CompareHashAndPassword([]byte(passWordHash), []byte(passWord))
}
