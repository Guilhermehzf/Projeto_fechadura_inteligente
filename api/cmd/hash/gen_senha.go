package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)


var senha string

func main() {
	gen_senha()
}
func gen_senha() {
	senha = ""

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erro ao gerar hash: %v", err)
	}

	fmt.Println("Senha:", senha)
	fmt.Println("Hash bcrypt:", string(hash))
}

