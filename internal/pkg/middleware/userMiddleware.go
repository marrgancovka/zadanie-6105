package middleware

import (
	"net/http"
	"zadanie-6105/internal/myErrors"
	"zadanie-6105/internal/pkg/users"
	"zadanie-6105/internal/pkg/utils"
)

type UserMiddleware struct {
	r users.UserRepository
}

func NewMiddleware(r users.UserRepository) *UserMiddleware {
	return &UserMiddleware{r: r}
}

func (m *UserMiddleware) UserExistsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")

		if username != "" {
			ok, err := m.r.UserIsExists(username)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
			if !ok {
				utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
				return
			}
		} else {
			utils.WriteError(w, http.StatusUnauthorized, myErrors.ErrUserNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
