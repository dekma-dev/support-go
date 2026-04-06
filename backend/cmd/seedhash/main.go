package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	for _, p := range []string{"admin123", "agent123", "client123"} {
		h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		fmt.Printf("%s -> %s\n", p, string(h))
	}
}
