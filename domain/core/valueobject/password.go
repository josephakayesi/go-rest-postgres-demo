package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Password string

func (p Password) Hash() string {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), 10)

	if err != nil {
		fmt.Println("unable to hash password")
		panic(err)
	}

	return string(hash)
}

func (p Password) DoesPasswordMatch(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hash()), []byte(plainPassword))
	if err != nil {
		fmt.Println("unable to compare hash and password. error: ", err)
		return false
	}

	return true
}

func (p Password) String() string {
	return string(p)
}
