package filters

import (
	"cos-backend-com/src/libs/auth"
	"net/http"
	"strings"

	"github.com/wujiu2020/strip"
)

// CORS Headers
func Cors(ctx strip.Context, log strip.ReqLogger, rw http.ResponseWriter, req *http.Request, authTr auth.RoundTripper) {
	origin := req.Header.Get("Origin")
	if strings.HasPrefix(origin, "https://d.comunion.io") || strings.HasPrefix(origin, "https://d1.comunion.io") {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Language, Content-Type")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Max-Age", "86400")
	}
	if req.Method == "OPTIONS" {
		rw.WriteHeader(204)
	}
}
