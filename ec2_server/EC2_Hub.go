package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/tidwall/redcon"
	"net"
	"strings"
	"time"
)

type lambdaInstance struct {
	name  string
	conn  *net.TCPConn
	alive bool
}

var items = make(map[string][]byte)

var itemHandler = make(chan []byte, 100)

func handler(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {

	case "set":
		if len(cmd.Args) != 3 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
		items[string(cmd.Args[1])] = cmd.Args[2]
		conn.WriteString("OK")
		fmt.Println(cmd.Args[2])
		itemHandler <- cmd.Args[2]

	case "get":

		lambdaTrigger("Lambda2SmallJPG")
		time.Sleep(1 * time.Second)
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		value := <-itemHandler
		fmt.Println(value)
		conn.WriteBulk(value)
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")

	}
}
func accept(conn redcon.Conn) bool {
	fmt.Println("accept:", conn.RemoteAddr())
	return true
}
func close(conn redcon.Conn, err error) {

}

func lambdaTrigger(name string) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String("us-east-1")})

	_, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(name)})
	fmt.Println(name)
	if err != nil {
		fmt.Println("Error calling LambdaFunction")
	}
}

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080") // create TcpAddr
	if err != nil {
		fmt.Println(err)
		return
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr) //Start listen
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("start listening", tcpAddr)
	err = redcon.Serve(tcpListener, handler, accept, close)
	if err != nil {
		fmt.Println("err is", err)
	}

}
