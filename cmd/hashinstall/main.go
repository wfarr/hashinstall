package main

import (
	"flag"
	"log"
	"os"

	"github.com/wfarr/hashinstall"
)

var (
	name    = flag.String("name", "", "The name of the hashicorp tool")
	version = flag.String("version", "", "The version of the hashicorp tool")
	destdir = flag.String("destdir", "", "Where to extract the hashicorp tool")
)

func main() {
	flag.Parse()
	expect(*name, *version, *destdir)

	log.Printf("installing %v %v into %v\n", *name, *version, *destdir)

	info, debug := make(chan string), make(chan string)
	go func() {
		select {
		case msg := <-info:
			log.Println(msg)
		case msg := <-debug:
			log.Println("debug", msg)
		}
	}()

	err := hashinstall.Install(*name, *version, *destdir, info, debug)
	if err != nil {
		log.Fatal(err)
	}
}

func expect(args ...string) {
	for _, arg := range args {
		if arg == "" {
			flag.Usage()
			os.Exit(1)
		}
	}
}
