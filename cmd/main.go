package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	usertransport "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/transports"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oklog/oklog/pkg/group"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(err.Error())
	}
	
	var (
		httpAddr = net.JoinHostPort("localhost", "9000");
	)

	var (
		service = userservice.NewUserService()
		endpoint = userendpoint.New(service)
		httpHandler = usertransport.NewHttpHandler(endpoint)
	)

	var g group.Group
	{
		httpListener, err := net.Listen("tcp", httpAddr)
		fmt.Println(httpListener.Addr().String(), err)
		
		g.Add(func() error {
			fmt.Println(httpAddr)
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

	g.Run()
}