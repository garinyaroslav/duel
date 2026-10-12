package internal

//
// import (
// 	"fmt"
// 	"log"
// 	"net"
// )
//
// func main() {
// 	raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:1053")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	c, err := net.DialUDP("udp", nil, raddr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer c.Close()
//
// 	if _, err := c.Write([]byte("This is a UDP message")); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	buf := make([]byte, 2048)
// 	n, err := c.Read(buf)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(string(buf[:n]))
// }
