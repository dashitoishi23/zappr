package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	database "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/database"
	tenantendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/endpoints"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	tenanttransports "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/transports"
	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	usertransport "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/transports"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/state"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/util"
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
		httpAddr = fmt.Sprintf("0.0.0.0:%v", 9000)
		// httpDebugAddr = fmt.Sprintf("0.0.0.0:%v", 9001)
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	state.InitState()

	db, dbErr := database.OpenDBConnection(os.Getenv("POSTGRESQL_CONN_STRING"))

	var servers []commonmodels.HttpServerConfig

	if dbErr == nil {
		fmt.Print(db.Statement.Vars...)
		var (
			userService = userservice.NewUserService(db)
			userEndpoint = userendpoint.New(userService, logger)
			userServers = usertransport.NewHttpHandler(userEndpoint)
		)

		servers = append(servers, userServers...)

		var (
			tenantService = repository.NewBaseCRUD[tenantmodels.Tenant](db)
			tenantEndpoint = tenantendpoint.New(tenantService, logger)
			tenantServers = tenanttransports.NewHandler(tenantEndpoint, logger)
		)

		servers = append(servers, tenantServers...)

		httpHandler := util.RootHttpHandler(servers)

		httpListener, err := net.Listen("tcp", httpAddr)
		fmt.Println(httpListener.Addr().String(), err)

		var g group.Group
		{				
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