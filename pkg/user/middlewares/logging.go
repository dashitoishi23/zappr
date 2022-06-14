package userservicemiddlewares

import (
	"context"
	"time"

	models "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/models"
	userservice "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/service"
	"github.com/go-kit/log"
)

type Middleware func(userservice.UserService) userservice.UserService

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next userservice.UserService) userservice.UserService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next userservice.UserService
}

func (l loggingMiddleware) GenerateJWTToken(ctx context.Context, userEmail string) (string, error) {
	func(begin time.Time){
		l.logger.Log("endpoint", "GenerateJWTToken",
	"userEmail", userEmail,
	"took", time.Since(begin).Milliseconds())
	}(time.Now())

	return l.next.GenerateJWTToken(ctx, userEmail)
}

func (l loggingMiddleware) ValidateLogin(ctx context.Context, jwtToken string) bool {
	func(begin time.Time){
		l.logger.Log("endpoint", "ValidateLogin",
					"took", time.Since(begin).Milliseconds(),
	)}(time.Now())
	
	return l.next.ValidateLogin(ctx, jwtToken)

}

func (l loggingMiddleware) SignupUser(ctx context.Context, newUser models.User) (models.User, error) {
	func(begin time.Time){
		l.logger.Log("endpoint", "SignupUser",
					"took", time.Since(begin).Milliseconds(),	
	)}(time.Now())

	return l.next.SignupUser(ctx, newUser)
}	

