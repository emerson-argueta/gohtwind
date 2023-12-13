package cmds

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenSecretKey() {
	fmt.Println("Generating secret key...")
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	str := base64.StdEncoding.EncodeToString(key)
	fmt.Println("Secret key generated:")
	fmt.Println(str)
	fmt.Println("Please copy this key and paste it in your .env file as GOHTWIND_SECRET_KEY")
}
