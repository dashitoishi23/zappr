package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	usertransport "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/transports"
	"github.com/oklog/oklog/pkg/group"
)

func main() {
	fs := flag.NewFlagSet("main", flag.ExitOnError)
	var (
		httpAddr = fs.String("http-addr", ":8081", "HTTP listen address")
	)

	var (
		service = userservice.NewUserService()
		endpoint = userendpoint.New(service)
		httpHandler = usertransport.NewHttpHandler(endpoint)
	)

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil{
			os.Exit(1)
		}
		g.Add(func() error {
			return http.Serve(httpListener, httpHandler)
		}, func(error){
			httpListener.Close()
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <- c:
				return fmt.Errorf("Received signal %s", sig)
			case <- cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
}