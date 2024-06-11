package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var nUrl = flag.String("u", "nats://acc:acc@localhost:4222", "URL for Nats server to publish to")
	flag.Parse()

	// Use signal exit
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	exit := make(chan struct{}, 1)

	// Create an unauthenticated connection to NATS.
	log.Println("Connecting to NATS on:", *nUrl)
	nc, err := nats.Connect(*nUrl)
	if err != nil {
		log.Fatal(err)
	}

	ctx, ca := context.WithCancel(context.Background())
	var count int

	// Run publisher script
	go func() {
		defer close(exit)

		for {
			time.Sleep(time.Duration(rand.IntN(1000) * int(time.Millisecond)))
			select {
			case <-ctx.Done():
				return
			default:
				count++
				// Messages are published to subjects
				if err := nc.Publish("ts.pub-golang-client", []byte(fmt.Sprintf("hello from GO - %d", count))); err != nil {
					log.Printf("Publish error: %v\n", err)
				}
			}
		}
	}()

	// Run default prometheus metrics
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8888", nil)
	}()

	log.Println("awaiting Ctrl+C signal")
	<-done

	log.Println("exiting")
	ca()
	<-exit
}
