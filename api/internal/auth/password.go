package auth

import (
	"fmt"
	"svitorz/url-shortner/internal/config"

	"golang.org/x/crypto/bcrypt"
)

var cost int

func init() {
	var err error
	cost, err = config.LoadHashCost()

	if err != nil {
		fmt.Println("Erro ao buscar custo de hash")
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes), err
}

func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println("VerifyPassword", err)
	return err == nil
}
