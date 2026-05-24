package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)
func main(){
  // var abel uint8 = 9
	b := make([]byte,4)
	_,_ = rand.Read(b)
	fmt.Println(b)

	fmt.Println(string(b))
    age := 252 

	fmt.Println(string(age))
generateRandomJTI()
}


func generateRandomJTI() {
	b := make([]byte, 16)
	_, _ = rand.Read(b) 
	fmt.Println(b)
	bStr := hex.EncodeToString(b)
	fmt.Println(bStr) 
} 