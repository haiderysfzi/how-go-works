// Lesson 21: A tour of the standard library: JSON, time, HTTP, files, flags.
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"
)

type User struct {
	Name    string    `json:"name"`
	Email   string    `json:"email,omitempty"` // omit when empty
	Age     int       `json:"age"`
	secret  string    // unexported: never serialised
	Created time.Time `json:"created"`
}

func main() {
	// ---- JSON ----
	u := User{Name: "Ada", Age: 36, secret: "x", Created: time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)}
	data, _ := json.Marshal(u)
	fmt.Println(string(data))

	var back User
	if err := json.Unmarshal([]byte(`{"name":"Bob","age":5}`), &back); err == nil {
		fmt.Printf("%+v\n", back.Name)
	}
	var generic map[string]any // unknown shape: numbers become float64
	json.Unmarshal([]byte(`{"a":1,"b":[true,"x"]}`), &generic)
	fmt.Println(generic)

	// ---- time ----
	d := 90 * time.Minute
	fmt.Println(d, d.Hours(), time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0).Format("2006-01-02"))
	// Go formats with a REFERENCE date: Mon Jan 2 15:04:05 MST 2006 (1 2 3 4 5 6 7).
	t, _ := time.Parse(time.RFC3339, "2026-11-20T10:30:00Z")
	fmt.Println(t.Weekday(), t.Format("Jan 2, 2006 at 3:04PM"))

	// ---- HTTP server and client (httptest runs a real server in-process) ----
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) { // Go 1.22 routing
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"hello": r.PathValue("name")})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/hello/gopher")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // ALWAYS close response bodies
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode, strings.TrimSpace(string(body)))
	// A real server: http.ListenAndServe(":8080", mux). Each request runs in its own goroutine.

	// ---- structured logging ----
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // drop the timestamp for stable output
			}
			return a
		},
	}))
	logger.Info("request handled", "status", 200, "path", "/hello")

	// ---- files and io ----
	f, _ := os.CreateTemp("", "x")
	defer os.Remove(f.Name())
	fmt.Fprintln(f, "line one")
	f.Close()
	content, _ := os.ReadFile(f.Name())
	fmt.Printf("%q\n", content)
}
