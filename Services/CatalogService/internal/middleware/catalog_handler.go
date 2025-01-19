package middleware

import (
	"CatalogService/internal/models"
	"CatalogService/internal/service"
	"encoding/json"
	"net/http"
)

func VerifyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyReq, err := http.NewRequest(http.MethodGet, "http://localhost:9999/validate/", nil)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}

		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		for _, cookie := range r.Cookies() {
			proxyReq.AddCookie(cookie)
		}
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, "Failed to connect to service", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			http.Redirect(w, r, "http://localhost:9997/login/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func VerifyAuthorMiddleware(next http.Handler, catalogService service.CatalogService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := service.NewProxyRequestStruct()
		respBody, _, err := p.ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)
		if err == p.ForbiddenError {
			http.Redirect(w, r, "http://localhost:9997/not-an-author/", http.StatusFound)
			return
		}
		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		var claims models.Claims

		if err := json.Unmarshal(respBody, &claims); err != nil {
			http.Error(w, "Status Forbidden", http.StatusForbidden)
			return
		}

		_, err = catalogService.GetAuthorByUserId(claims.UserId)
		if err != nil {
			http.Redirect(w, r, "http://localhost:9997/not-an-author/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func VerifyNotAuthorMiddleware(next http.Handler, catalogService service.CatalogService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := service.NewProxyRequestStruct()
		respBody, _, err := p.ProxyRequest(r, "http://localhost:9999/get-claims/", nil, http.MethodGet)
		if err == p.ForbiddenError {
			http.Redirect(w, r, "http://localhost:9997/not-an-author/", http.StatusFound)
			return
		}
		if err != nil {
			http.Error(w, "Server Internal Error", http.StatusInternalServerError)
			return
		}

		var claims models.Claims

		if err := json.Unmarshal(respBody, &claims); err != nil {
			http.Error(w, "Status Forbidden", http.StatusForbidden)
			return
		}

		_, err = catalogService.GetAuthorByUserId(claims.UserId)
		if err == nil {
			return
		}

		next.ServeHTTP(w, r)
	})
}
