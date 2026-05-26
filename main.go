package main

import (
	"encoding/json"
	"fmt"
	"time"
)
var person map[string]string = map[string]string{

	"name":"abel",
}
func main() {
	// Convert seconds into a human-readable time object
	t := time.Unix(1700000000, 0)
   json.Unmarshal(person,t)
	fmt.Println(t) // Output: 2023-11-14 22:13:20 +0000 UTC (depending on your local timezone)
}
