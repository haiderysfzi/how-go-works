// Lesson 21: A tour of the standard library: JSON, time, HTTP, logging, files.
//
// WHAT YOU WILL LEARN
//   - Converting Go values to and from JSON, and how struct tags control it
//   - Go's unusual way of formatting dates
//   - Running an HTTP server and calling it with an HTTP client
//   - Structured logging with log/slog
//   - Reading and writing files
//
// BACKGROUND: Go's "standard library" ships with the language and covers most everyday
// needs (web servers, JSON, crypto, compression, testing) without third-party
// packages. It is stable, well documented at pkg.go.dev, and the best place to learn
// idiomatic Go by reading its source.
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

// STRUCT TAGS control JSON. The text in backticks is read at runtime by encoding/json
// (via reflection: inspecting types while the program runs):
//   - `json:"name"`            use "name" as the JSON key (instead of the field name "Name")
//   - `json:"email,omitempty"` leave the key out entirely when the value is empty/zero
//
// Only EXPORTED (capitalised) fields are visible to the json package, because reflection
// from another package cannot see private fields. So `secret` is silently skipped.
type User struct {
	Name    string    `json:"name"`
	Email   string    `json:"email,omitempty"` // omit when empty
	Age     int       `json:"age"`
	secret  string    // unexported: never serialised
	Created time.Time `json:"created"`
}

func main() {
	// ---- JSON ----
	// MARSHAL = Go value -> JSON bytes. Unmarshal = JSON bytes -> Go value.
	// (The "_" ignores the error to keep the demo short; real code must check it.)
	// time.Date(year, month, day, hour, min, sec, nanosec, location).
	u := User{Name: "Ada", Age: 36, secret: "x", Created: time.Date(2026, 11, 20, 0, 0, 0, 0, time.UTC)}
	data, _ := json.Marshal(u)
	fmt.Println(string(data)) // {"name":"Ada","age":36,"created":"2026-11-20T00:00:00Z"}  (no email, no secret)

	// Unmarshal needs a POINTER (&back) so it can fill in the struct. Keys that are
	// missing from the JSON (email, created) simply stay at their zero value.
	var back User
	if err := json.Unmarshal([]byte(`{"name":"Bob","age":5}`), &back); err == nil {
		fmt.Printf("%+v\n", back.Name)
	}
	// When you do not know the shape in advance, decode into map[string]any.
	// JSON objects become map[string]any, arrays become []any, and ALL numbers become float64.
	var generic map[string]any // unknown shape: numbers become float64
	json.Unmarshal([]byte(`{"a":1,"b":[true,"x"]}`), &generic)
	fmt.Println(generic) // map[a:1 b:[true x]]

	// ---- time ----
	// time.Duration is a span of time stored as a count of nanoseconds; constants like
	// time.Minute let you write 90 * time.Minute. Printing gives "1h30m0s".
	// AddDate(years, months, days) adds calendar units: Jan 31 + 1 month = "Feb 31",
	// which does not exist, so Go NORMALISES it forward to March 3. Date math has
	// many edge cases like this.
	d := 90 * time.Minute
	fmt.Println(d, d.Hours(), time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0).Format("2006-01-02"))
	// DATE FORMATTING: instead of codes like "%Y-%m-%d", Go uses an EXAMPLE date as the
	// template: Mon Jan 2 15:04:05 MST 2006. Written with numbers it reads 1 2 3 4 5 6 7
	// (month 1, day 2, hour 3 (15 = 3pm), minute 4, second 5, year 6, zone 7).
	// You show Go how THAT moment would look in your format, and it reformats the real
	// time to match. "2006-01-02" means year-month-day.
	// time.Parse(layout, text) reads a string; time.RFC3339 is a ready-made layout
	// for the common ISO timestamp.
	t, _ := time.Parse(time.RFC3339, "2026-11-20T10:30:00Z")
	fmt.Println(t.Weekday(), t.Format("Jan 2, 2006 at 3:04PM"))

	// ---- HTTP server and client (httptest runs a real server in-process) ----
	// An HTTP SERVER listens for requests and answers each with a response. In net/http:
	//   - a Handler is "something that handles a request" (a func(w, r) works)
	//   - a ServeMux (the "router") maps URL patterns to handlers
	// The pattern "GET /hello/{name}" (Go 1.22+) matches only GET requests and captures
	// the part of the path in braces; r.PathValue("name") reads it back.
	//   w (ResponseWriter) is where you write the response (headers, status, body)
	//   r (*Request) describes what the client sent
	// json.NewEncoder(w).Encode(...) writes JSON straight into the response body.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) { // Go 1.22 routing
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"hello": r.PathValue("name")})
	})
	// httptest.NewServer starts a REAL HTTP server on a random free port of this
	// machine and gives its URL, so the demo needs no fixed port and does not clash
	// with anything. `defer srv.Close()` shuts it down when main ends.
	srv := httptest.NewServer(mux)
	defer srv.Close()

	// The CLIENT side: http.Get sends a GET request and returns the response once the
	// headers arrive. The body is a stream that must be read and then CLOSED, or the
	// underlying network connection leaks and cannot be reused.
	resp, err := http.Get(srv.URL + "/hello/gopher")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()                                       // ALWAYS close response bodies
	body, _ := io.ReadAll(resp.Body)                              // read the whole stream into memory
	fmt.Println(resp.StatusCode, strings.TrimSpace(string(body))) // 200 {"hello":"gopher"}
	// A real server: http.ListenAndServe(":8080", mux). Each request runs in its own
	// goroutine (lesson 15), which is how one Go process handles thousands of
	// simultaneous connections. That also means handlers sharing data must guard it
	// (lesson 17).

	// ---- structured logging ----
	// log/slog writes log lines as key=value pairs (or JSON) instead of free-form text,
	// so tools can search and filter them: level=INFO msg="request handled" status=200 ...
	// The handler options below use ReplaceAttr to drop the timestamp attribute so the
	// output is identical on every run. Returning an empty slog.Attr deletes it.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{} // drop the timestamp for stable output
			}
			return a
		},
	}))
	// After the message come alternating key, value pairs.
	logger.Info("request handled", "status", 200, "path", "/hello")

	// ---- files and io ----
	// os.CreateTemp makes a uniquely named file in the system's temp folder and returns
	// it already open. Files implement io.Writer, so Fprintln can write into them
	// (the same interface idea as lesson 10). Close() flushes and releases the OS
	// handle; we close BEFORE re-reading it. os.ReadFile reads the entire file into a
	// []byte; %q prints it quoted so the "\n" at the end is visible.
	f, _ := os.CreateTemp("", "x")
	defer os.Remove(f.Name()) // delete the temp file when main ends
	fmt.Fprintln(f, "line one")
	f.Close()
	content, _ := os.ReadFile(f.Name())
	fmt.Printf("%q\n", content) // "line one\n"
}
