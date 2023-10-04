package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alecthomas/kong"
	"github.com/ph1048/ebpfdbg/ebpflog"
)

var CLI struct {
	Serve struct {
		Port    int    `help:"HTTP port to listen" default:"31337"`
		Address string `help:"Interface IP to listen on" default:"0.0.0.0"`
	} `cmd:"" help:"Serve HTML based interface via HTTP"`
	Generate struct {
		Path string `help:"Path to save HTML file" default:"out.html"`
	} `cmd:"" help:"Save HTML file"`
	Input string `required help:"Path to source eBPF verifier log data"`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "serve":
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			file, err := os.Open(CLI.Input)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			ctx := ebpflog.NewLogContext(file, w)
			err = ctx.Process()
			if err != nil {
				log.Fatal(err)
			}
		})
		listenPort := CLI.Serve.Address + ":" + strconv.Itoa(CLI.Serve.Port)
		fmt.Println("Serving on http://" + listenPort)
		fmt.Println("Press Ctrl+C to stop")
		err := http.ListenAndServe(listenPort, nil)
		if err != nil {
			log.Fatal(err)
		}
	case "generate":
		file, err := os.Open(CLI.Input)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		outfile, err := os.Create(CLI.Generate.Path)
		if err != nil {
			log.Fatal(err)
		}
		defer outfile.Close()
		ctx := ebpflog.NewLogContext(file, outfile)
		err = ctx.Process()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Written to", CLI.Input)
	}
}
