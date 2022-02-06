package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/neonxp/geezer"

	"github.com/neonxp/geezer-example/services/hello"
)

var listen string

func main() {
	flag.StringVar(&listen, "listen", ":3000", "Host and port to listen (ex: '0.0.0.0:3000')")
	flag.Parse()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	app := geezer.NewHttpKernel()

	hello.RegisterHooks(app)
	_ = app.Register(hello.ServiceName, &hello.Service{})

	log.Printf("Started on %s\n", listen)
	srv := http.Server{Addr: listen, Handler: app}
	go func() {
		<-ctx.Done()
		srv.Close()
	}()
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
