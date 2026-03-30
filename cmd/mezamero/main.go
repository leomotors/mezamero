package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/leomotors/mezamero/internal/config"
	"github.com/leomotors/mezamero/internal/wol"
)

//go:embed static
var staticFS embed.FS

// deviceDTO is the JSON shape for the UI.
type deviceDTO struct {
	MAC             string `json:"mac"`
	IP              string `json:"ip,omitempty"`
	Name            string `json:"name"`
	NameOriginal    string `json:"name_original,omitempty"`
	Description     string `json:"description,omitempty"`
	Spec            string `json:"spec,omitempty"`
	Image           string `json:"image"`
	BackgroundColor string `json:"background_color"`
	ForegroundColor string `json:"foreground_color"`
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config.yaml")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	root, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/devices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		out := make([]deviceDTO, 0, len(root.Devices))
		for i := range root.Devices {
			d := &root.Devices[i]
			out = append(out, deviceDTO{
				MAC:             d.MAC,
				IP:              d.IP,
				Name:            d.Name,
				NameOriginal:    d.NameOriginal,
				Description:     d.Description,
				Spec:            d.Spec,
				Image:           d.Image,
				BackgroundColor: d.BackgroundColor,
				ForegroundColor: d.ForegroundColor,
			})
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		enc.SetEscapeHTML(true)
		if err := enc.Encode(out); err != nil {
			http.Error(w, "encode", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/api/wake", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var body struct {
			MAC string `json:"mac"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		body.MAC = strings.TrimSpace(body.MAC)
		if body.MAC == "" {
			http.Error(w, "mac required", http.StatusBadRequest)
			return
		}
		var found *config.Device
		for i := range root.Devices {
			if normalizeMAC(root.Devices[i].MAC) == normalizeMAC(body.MAC) {
				found = &root.Devices[i]
				break
			}
		}
		if found == nil {
			http.Error(w, "unknown device", http.StatusNotFound)
			return
		}
		if err := wol.Send(found.MAC); err != nil {
			http.Error(w, fmt.Sprintf("wake failed: %v", err), http.StatusBadGateway)
			return
		}
		log.Printf("wake sent: device=%q mac=%s time=%s", found.Name, normalizeMAC(found.MAC), time.Now().UTC().Format(time.RFC3339Nano))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	staticRoot, err := fs.Sub(staticFS, "static")
	if err != nil {
		fmt.Fprintf(os.Stderr, "static: %v\n", err)
		os.Exit(1)
	}
	mux.Handle("/", spaHandler(staticRoot))

	hp := browserHostPort(*addr)
	const ansiBold, ansiReset = "\x1b[1m", "\x1b[0m"
	log.Printf("mezamero server started: listening on %s; open http://%s%s%s", *addr, ansiBold, hp, ansiReset)

	if err := http.ListenAndServe(*addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "server: %v\n", err)
		os.Exit(1)
	}
}

func normalizeMAC(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), "-", ":"))
}

// browserHostPort returns host:port suitable for a local browser (e.g. localhost:8080).
func browserHostPort(addr string) string {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		if strings.HasPrefix(addr, ":") {
			return "localhost" + addr
		}
		return "localhost:8080"
	}
	switch host {
	case "", "0.0.0.0", "::", "[::]":
		return net.JoinHostPort("localhost", port)
	default:
		return net.JoinHostPort(host, port)
	}
}

// spaHandler serves index.html for / and static files under the same root.
func spaHandler(root fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" || path == "index.html" {
			data, err := fs.ReadFile(root, "index.html")
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write(data)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}
