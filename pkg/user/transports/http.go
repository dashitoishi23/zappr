package usertransport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	userendpoint "dev.azure.com/technovert-vso/Zappr/_git/Zappr/pkg/user/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(endpoints userendpoint.Set) http.Handler {
	m:= http.NewServeMux()

	m.Handle("/token", httptransport.NewServer(
		endpoints.GenerateToken,
		decodeGenerateTokenHTTPRequest,
		encodeHTTPGenericResponse,
	))

	return m
}

func decodeGenerateTokenHTTPRequest(_ context.Context, r *http.Request) (interface{}, error){
	var req userendpoint.GenerateTokenRequest
	fmt.Println(req)

	err := json.NewDecoder(r.Body).Decode(&req)

	if err == io.EOF{
		return req, nil	
	}
	
	return req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(error); ok && e != nil {
		fmt.Print(e)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}