package logger

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type logkey string

const key logkey = "logger"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		l := zap.NewExample()
		l = l.With(zap.Namespace("hometic"))
		l.Info(formatRequest(r))

		c := context.WithValue(r.Context(), key, l)
		newR := r.WithContext(c)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, newR)
	})
}

func L(ctx context.Context) *zap.Logger {
	val := ctx.Value(key)
	if val == nil {
		return zap.NewExample()
	}

	l, ok := val.(*zap.Logger)
	if ok {
		return l
	}

	return zap.NewExample()
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}
