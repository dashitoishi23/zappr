package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	database "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/database"
	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	userservicemiddlewares "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/middlewares"
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	usertransport "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/transports"
	"github.com/go-kit/log"
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

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	db, dbErr := database.OpenDBConnection(os.Getenv("POSTGRESQL_CONN_STRING"))

	if dbErr == nil {
		fmt.Print(db.Statement.Vars...)
		var (
			userService = userservice.NewUserService(db)
			endpoint = userendpoint.New(userService, logger)
			userHttpHandler = usertransport.NewHttpHandler(endpoint)
		)
		userservicemiddlewares.LoggingMiddleware(logger)(userService)

		// var (
		// 	tenantService = repository.NewBaseCRUD[tenantmodels.Tenant](db)
		// 	tenantEndpoint = tenantendpoint.New(tenantService, logger)
		// 	tenantHttpHandler = tenanttransports.NewHandler(tenantEndpoint)
		// )

		// fmt.Print(tenantHttpHandler)

		httpListener, err := net.Listen("tcp", httpAddr)
		fmt.Println(httpListener.Addr().String(), err)

		http.Serve(httpListener, userHttpHandler)
		// http.Serve(httpListener, tenantHttpHandler)

		var g group.Group
		// {	
		// 	g.Add(func() error {
		// 		fmt.Println(httpAddr)
		// 		return http.Serve(httpListener, userHttpHandler)
		// 	})
		// }
		// {
		// 	// httpListener, err := net.Listen("tcp", httpAddr)
		// 	// fmt.Println(err)
			
		// 	g.Add(func() error {
		// 		fmt.Println(httpAddr)
		// 		return http.Serve(httpListener, tenantHttpHandler)
		// 	}, func(error){
		// 		httpListener.Close()
		// 	})
		// }
		{
			cancelInterrupt := make(chan struct{})
			g.Add(func() error {
				c := make(chan os.Signal, 1)
				signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
				select {
				case sig := <- c:
					return fmt.Errorf("received signal %s", sig)
				case <- cancelInterrupt:
					return nil
				}
			}, func(error) {
				close(cancelInterrupt)
			})
		}
	
		g.Run()
	}
	
	fmt.Print(dbErr.Error())
}