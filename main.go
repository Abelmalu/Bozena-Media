package main

import (
	"encoding/base64"
	"fmt"
	"strconv"
)

func main() {
	num := 123

	// 1. Convert the integer to its string format ("123")
	numStr := strconv.Itoa(num)

	// 2. Encode the string bytes to Base64
	encoded := base64.StdEncoding.EncodeToString([]byte(numStr))

	fmt.Println(encoded) // Output: MTIz
}
