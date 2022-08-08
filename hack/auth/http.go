package auth

import (
	"github.com/okcredit/go-common/errors"
	"github.com/okcredit/go-common/httpx"
	"log"
	"net/http"
	"strings"
)

var ErrUnauthorized = errors.From(401, "unauthorized")

func HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(HeaderAuthorization)
		tokens := strings.Split(authHeader, " ")
		if len(tokens) != 2 {
			log.Printf("auth: invalid auth header '%s'", authHeader)
			httpx.WriteError(w, ErrUnauthorized)
			return
		}
		tokenType := tokens[0]
		token := tokens[1]

		switch tokenType {

		case TokenTypeBearer:
			// merchant token (app to service authentication)
			merchantID, err := VerifyMerchantToken(token)
			if err != nil {
				log.Printf("auth: invalid merchant token: %v", err)
				httpx.WriteError(w, ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, WithMerchantContext(r, merchantID))

		case TokenTypeAdmin:
			// admin token (service to service authentication)
			if err := VerifyAdminToken(token, r); err != nil {
				log.Printf("auth: invalid admin token: %v", err)
				httpx.WriteError(w, ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, r)

		default:
			log.Printf("auth: invalid token type '%s'", tokenType)
			httpx.WriteError(w, ErrUnauthorized)
			return
		}
	})
}

func NewHttpClient() *http.Client {
	return &http.Client{Transport: &authorizer{}}
}

type authorizer struct{}

func (*authorizer) RoundTrip(r *http.Request) (*http.Response, error) {
	r = AuthorizeAdmin(r)
	return http.DefaultTransport.RoundTrip(r)
}
