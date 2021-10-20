package main

import "fmt"

func main() {

	for i := 0; i < 10; i++ {
		pwdhash, _ := AesEncrypt([]byte("2019339964026"))
		fmt.Println(AesDecrypt(pwdhash))
	}
}
