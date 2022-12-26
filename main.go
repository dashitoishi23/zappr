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
	initsetupendpoints "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/endpoints"
	initsetupmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/models"
	initsetupservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/service"
	initsetuptransports "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/initsetup/transports"
	masterroleendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/endpoints"
	masterrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/models"
	masterroletransports "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/role/transports"
	tenantendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/endpoints"
	tenantmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/models"
	tenantservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/service"
	tenanttransports "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/tenant/transports"
	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	usermodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	usertransport "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/transports"
	usermetadataendpoints "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/endpoints"
	usermetadatamodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/models"
	usermetadatatransports "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/usermetadata/transports"
	userrolemodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/userrole/models"
	redisutil "dev.azure.com/technovert-vso/Zappr/_git/Zappr/redis"
	"dev.azure.com/technovert-vso/Zappr/_git/Zappr/repository"
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

	//Initialise Zappr Config userDB
	db, dbErr := database.OpenZapprConfigDBConnection()

	redisPool := &redisutil.RedisPool{}

	redisPool.NewPool()

	client := redisPool.GetPool().Get()

	defer client.Close()

	var servers []commonmodels.HttpServerConfig

	if dbErr == nil {
		fmt.Print("Database connected!")
		fmt.Print(db.Statement.Vars...)

		var (
			initRepository = repository.Repository[initsetupmodels.DBConfig](db)
			initsetupService = initsetupservice.NewInitSetupService(initRepository)
			initSetupEndpoints = initsetupendpoints.New(initsetupService, logger)
			initSetupServers = initsetuptransports.NewHttpHandler(initSetupEndpoints)
		)

		servers = append(servers, initSetupServers...)

		userDB, err := database.OpenDBConnection(initRepository)

		if err == nil {

		var (
			userService = userservice.NewUserService(repository.Repository[usermodels.User](userDB), 
			repository.Repository[masterrolemodels.Role](userDB), repository.Repository[userrolemodels.UserRole](userDB), 
			repository.Repository[tenantmodels.Tenant](userDB), repository.Repository[usermodels.APIKey](userDB), 
			repository.Repository[usermodels.OAuthProvider](userDB))

			userEndpoint = userendpoint.New(userService, logger)
			userServers = usertransport.NewHttpHandler(userEndpoint)
		)

		servers = append(servers, userServers...)

		var (
			tenantService = repository.NewBaseCRUD(repository.Repository[tenantmodels.Tenant](userDB))
			tenantSpecificService = tenantservice.NewService(repository.Repository[tenantmodels.Tenant](userDB), 
			repository.Repository[masterrolemodels.Role](userDB))
			tenantEndpoint = tenantendpoint.New(tenantService, tenantSpecificService, logger, client)
			tenantServers = tenanttransports.NewHandler(tenantEndpoint, logger)
		)

		servers = append(servers, tenantServers...)

		var (
			roleService = repository.NewBaseCRUD(repository.Repository[masterrolemodels.Role](userDB))
			roleEndpoint = masterroleendpoint.New(roleService, logger)
			roleServers = masterroletransports.NewHandler(roleEndpoint, logger)
		)

		servers = append(servers, roleServers...)

		var (
			userMetadataService = repository.NewBaseCRUD(repository.Repository[usermetadatamodels.UserMetadata](userDB))
			userMetadataEndpoint = usermetadataendpoints.New(userMetadataService, logger)
			userMetadataServers = usermetadatatransports.NewHandler(userMetadataEndpoint, logger)
		)

		servers = append(servers, userMetadataServers...)

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
	}
	
	fmt.Print(dbErr.Error())
}