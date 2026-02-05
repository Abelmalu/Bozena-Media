package main

import (
	"context"
	"time"
)

func main() {
	make 

	ctx,cancel := context.WithTimeout(context.Background(),time.Second*2)
}