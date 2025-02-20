package main

import (
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"

	auth "github.com/abbot/go-http-auth"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

// Get    curl -i 'http://localhost:8080?foo=<strong>bar</strong>'
// Post:  curl -i -d "<em>Hi</em>" 'http://localhost:8080?foo=<strong>bar</strong>'
// admin:long-memorable-password
// Auth1:  curl -i 'http://localhost:8080/authenticated1' -H 'Authorization: Basic YWRtaW46bG9uZy1tZW1vcmFibGUtcGFzc3dvcmQ='
// Auth2:  curl -i 'http://localhost:8080/authenticated2' -H 'Authorization: Basic YWRtaW46bG9uZy1tZW1vcmFibGUtcGFzc3dvcmQ='

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "can't read the env: %s", err)
	}
	http.HandleFunc("/", getHandler)
	http.HandleFunc("/404", http.NotFound)
	http.HandleFunc("/500", internalServerError)
	http.HandleFunc("/authenticated1", authenticateHandler)

	authenticator := auth.NewBasicAuthenticator("localhost:8080", secret)
	http.HandleFunc("/authenticated2", authenticator.Wrap(authHandler))

	http.HandleFunc("/limited", limitHandler)

	fmt.Fprintf(os.Stdout, "Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "Port is not empty %s", err)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")

	bytes, err := os.ReadFile("./index.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't open the file: %s", err)
		return
	}
	htmlCode := strings.Split(string(bytes), "</body>")
	w.Write([]byte(htmlCode[0]))

	query := r.URL.Query()
	if len(query) > 0 {
		x := "<p>Query parameters:</p><ul>"
		for key, values := range query {
			escapeValue := html.EscapeString(values[0])
			x += fmt.Sprintf("<li>%s: [%s]</li>", key, escapeValue)
		}
		x += "</ul>"
		w.Write([]byte(x))
	}

	if r.Method == http.MethodPost {
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't read the body of request: %s", err)
		}
		escapeBody := html.EscapeString(string(requestBody))
		w.Write([]byte(escapeBody))
	}
	w.Write([]byte("</body>"))
	w.Write([]byte(htmlCode[1]))
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)

	w.Write([]byte("Internal server error"))
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	authValue := header["Authorization"]
	if len(authValue) == 0 {
		unauthorizedResponse(w)
		return
	}
	base4Credential := strings.Split(authValue[0], " ")[1]
	decodedCredential, err := base64.StdEncoding.DecodeString(base4Credential)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't decoding: %s", err)
		unauthorizedResponse(w)
		return
	}
	credential := strings.Split(string(decodedCredential), ":")
	userName := credential[0]
	password := credential[1]
	if userName == os.Getenv("AUTH_USERNAME") && password == os.Getenv("AUTH_PASSWORD") {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/plain")
		text := fmt.Sprintf("Hello, %s", userName)
		w.Write([]byte(text))
		return
	}
	unauthorizedResponse(w)
}

func unauthorizedResponse(w http.ResponseWriter) {
	w.WriteHeader(401)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Unauthorized"))
}

func secret(user, realm string) string {
	if user == os.Getenv("AUTH_USERNAME") {
		// https://pkg.go.dev/github.com/abbot/go-http-auth#section-readme
		// openssl passwd -1 -salt dkjfhsjk long-memorable-password   ($1$ means using MD5 algorithm)
		return "$1$dkjfhsjk$ChcNEj9N9NVhLYMuPPvQe1"
	}
	return ""
}

func authHandler(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	fmt.Fprintf(w, "Hello, %s!", r.Username)
}

// Why I don't get any Failed requests?
func limitHandler(w http.ResponseWriter, r *http.Request) {
	limiter := rate.NewLimiter(100, 30)
	if limiter.Allow() {
		w.WriteHeader(200)
		w.Write([]byte("Limited"))
		return
	}
	w.WriteHeader(503)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Service Unavailable"))
}
