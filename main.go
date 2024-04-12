package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func proxy(c *gin.Context) {
	remote, err := url.Parse("https://88.99.21.177:443")
	if err != nil {
		log.Panicln("Cannot connect to remote addr: ", err)
	}
	// cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("cert.pem")})
	// certificate, err := tls.X509KeyPair("cert.pem", "key.pem")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	// proxy.Transport = &http.Transport{
	// 	TLSClientConfig: &tls.Config{
	// 		Certificates:       []tls.Certificate{},
	// 		InsecureSkipVerify: true,
	// 	},
	// }
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	log.Println("Header: ", c.Request.Header)
	log.Println("Body: ", c.Request.Body)
	log.Println("Host: ", c.Request.Host)
	log.Println("URL: ", c.Request.URL)
	log.Println("Content lenngth: ", c.Request.ContentLength)
	log.Println(c.Writer.Header())
	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {
	route := gin.Default()
	route.Any("*path", proxy)
	// route.Run(":8080")
	go func() {
		// Run HTTP server on port 8080
		if err := http.ListenAndServe(":8080", route); err != nil {
			log.Fatalf("HTTP server failed: %s", err)
		}
	}()

	go func() {
		// Run HTTPS server on port 8443
		if err := http.ListenAndServeTLS(":8443", "server.cert", "server.key", route); err != nil {
			log.Fatalf("HTTPS server failed: %s", err)
		}
	}()

	// Block main goroutine to keep servers running
	select {}
}
