package main

import (
	"first/chnl"
	"first/internal/httptransport"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// client = redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379", // where the redis server is (6379 is default port for redis)

	// })

	var port int
	flag.IntVar(&port, "port", 0, "Socket on")

	flag.Parse()

	ch := make(chan chnl.Collection)

	server := &http.Server{Handler: httptransport.NewHandler(chnl.NewInMem(ch))}

	go func() {

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

		if err != nil {
			log.Panicf("cannot create tpc listener: %v", err)
		}

		log.Printf("      starting http server on %q", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Panicf("cannot start http server: %v", err)
		}

	}()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	log.Printf("Got exit signal %q. Bye", <-sig)

	// http.Handle("/", r)
	// http.ListenAndServe(":8080", nil)
}
