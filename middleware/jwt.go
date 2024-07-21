package middleware

import (
	"context"
	"fmt"
	"net/http"
)

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//// for now I am skipping token authentication on register and assuming the client is not controlled by me
		////and that the client doesn'thave a way to get a public token
		ctx := r.Context()
		fmt.Println(r.URL.Path)
		if r.URL.Path == "/api/v1/register" || r.URL.Path == "/api/v1/login" {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("Authorization")
		c, err := ValidateJWT(token)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := ParseClaims(c)
		ctx = context.WithValue(ctx, "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
