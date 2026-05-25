package main

import (
	"encoding/base64"
	"fmt"
)
func main() {
    // 1. Define your data (bits/bytes)
    data := []byte("hello world")
	
    var name string = "abel"
   // 2. Encode to Base64 string
	fmt.Println(string(data))
    encoded := base64.StdEncoding.EncodeToString(data)

    fmt.Println(encoded) // Output: aGVsbG8gd29ybGQ=

	sentences := "don't get frustrAtedሀ"

	
	var age rune = 2 

	fmt.Println(name,string(age))
	

	for _,sentence := range sentences {

		if sentence == 65 {

			fmt.Println(string(sentence))
			

		}

		if sentence == 4608{

			fmt.Println(string(sentence))

			


		}

		
	}



}
