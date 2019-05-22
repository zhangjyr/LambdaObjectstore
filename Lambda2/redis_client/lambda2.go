package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gomodule/redigo/redis"
)

var c, ok = redis.Dial("tcp", "52.201.234.235:8080")

func HandleRequest() {
	fmt.Println("client", c)
	if ok != nil {
		// handle error
		fmt.Println("client err", ok)
	}
	n, err := c.Do("SET", "k1", []byte{233, 1})
	if err != nil {
		fmt.Println("Set error", err)
	}
	fmt.Println("set complete", n)
	//n, _ = redis.Int(c.Do("INCR", "k1"))
	//fmt.Printf("%#v\n", n)

	//defer c.Close()
}

func main() {
	lambda.Start(HandleRequest)
}
