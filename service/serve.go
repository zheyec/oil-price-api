package service

import (
	"context"
	"fmt"
	"net/http"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("content-type", "application/json") 
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

var middleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[/]=>remote=>%s host=>%s   url=>%s   method=>%s\n", r.RemoteAddr, r.Host, r.URL, r.Method)
		next.ServeHTTP(w, r)
	})
}

//Run HTTP Server
func Run(ctx context.Context, port int) {

	//http router
	http.Handle("/get_oil_price", http.HandlerFunc(oilPricesHandle))

	ip := fmt.Sprintf(":%v", port)
	fmt.Printf("listen on:%s\n", ip)
	http.ListenAndServe(ip, nil)
}
