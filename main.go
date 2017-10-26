package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/machinebox/sdk-go/facebox"
	"github.com/machinebox/sdk-go/x/boxutil"
)

func main() {
	var (
		addr  = flag.String("addr", ":9000", "address")
		state = flag.String("state", "./faces.machinebox.facebox", "facebox state file")
	)

	flag.Parse()
	facebox := facebox.New("http://localhost:8080")
	fmt.Println(`Web Face ID by Machine Box - https://machinebox.io/`)

	fmt.Println("Waiting for Facebox to be ready...")
	boxutil.WaitForReady(context.Background(), facebox)
	fmt.Println("Done!")

	fmt.Println("Setup facebox state")

	f, err := os.Open(*state)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = facebox.PostState(f)
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println("Done!")
	f.Close()

	fmt.Println("Go to:", *addr+"...")

	srv := NewServer("./assets", facebox)
	if err := http.ListenAndServe(*addr, srv); err != nil {
		log.Fatalln(err)
	}
}
