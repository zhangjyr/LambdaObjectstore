package Lambda1

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type DialConfig struct {
	Network string
	Address string
}

func main() {
	//c, err := redis.Dial("tcp", ":3333")
	c, err := redis.Dial("tcp", "52.201.234.235:8080")
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	n, _ := redis.Bytes(c.Do("GET", "k1"))
	fmt.Println(n)
	//n, _ = redis.Int(c.Do("INCR", "k1"))
	//fmt.Printf("%#v\n", n)

	defer c.Close()
}
