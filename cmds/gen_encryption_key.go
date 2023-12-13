package cmds

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenEncryptionKey() {
	fmt.Println("Generating encryption key...")
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	str := base64.StdEncoding.EncodeToString(key)
	fmt.Println("Encryption key generated:")
	fmt.Println(str)
	fmt.Println("Please copy this key and paste it in your .env file as GOHTWIND_32B_ENC_KEY")
}
