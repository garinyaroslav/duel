package internal

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"
)

func startServer() {
	pc, err := net.ListenPacket("udp", ":1053")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	slog.Debug("Server started", "address", pc.LocalAddr().String())

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			log.Println("read:", err)
			continue
		}
		go respond(pc, addr, buf[:n])
		// go func() {
		// 	ts := time.Now().Format(time.ANSIC)
		// 	msg := fmt.Sprintf("time received: %s. Your message: %s!", ts, string(buf[:n]))
		// 	if _, err := pc.WriteTo([]byte(msg), addr); err != nil {
		// 		log.Println("write:", err)
		// 	}
		// }()
	}
}

func respond(pc net.PacketConn, addr net.Addr, payload []byte) {
	ts := time.Now().Format(time.ANSIC)
	msg := fmt.Sprintf("time received: %s. Your message: %s!", ts, string(payload))
	if _, err := pc.WriteTo([]byte(msg), addr); err != nil {
		log.Println("write:", err)
	}
}
