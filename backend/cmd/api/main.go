package main

import (
	"errors"
	"fmt"
	"go-auth-postgres/backend/server"
	"net/http"
)

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %v", err))
	}
}
