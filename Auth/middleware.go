package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

//var UserCtxKey = &core.ContextKey{Name:"User"}

func UserContext(next http.Handler) http.Handler{
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authorizationH := r.Header.Get("Authorization")
		if authorizationH != "" {
			_, token, err := ParseAuthorizationHeader(authorizationH)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, map[string]string{"error": "Invalid request " + fmt.Sprintf("%v", err)})
				return 
			}
			user, err := GetUserFromToken(token)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, map[string]string{"error": "Invalid authentication token " + fmt.Sprintf("%v", err)})
				return
			}
			ctx = context.WithValue(ctx, "User", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}else{
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Authorization header not found"})
		}
	}
	return http.HandlerFunc(fn)
}
