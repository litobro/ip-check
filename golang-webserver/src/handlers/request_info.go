package handlers

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type RequestInfo struct {
	IP            string `json:"ip"`
	UserAgent     string `json:"user_agent"`
	RequestMethod string `json:"request_method"`
	Referrer      string `json:"referrer"`
	XForwardedFor string `json:"x_forwarded_for"`
	Hostname      string `json:"hostname"`
	Url           string `json:"url"`
	PrettyJson    string `json:"-"`
}

//go:embed templates/request_info.html
var requestInfoTemplate string

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = "unknown"
	}
	userAgent := r.UserAgent()
	requestMethod := r.Method
	referrer := r.Referer()
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	hostname := r.Host

	requestInfo := RequestInfo{
		IP:            ip,
		UserAgent:     userAgent,
		RequestMethod: requestMethod,
		Referrer:      referrer,
		XForwardedFor: xForwardedFor,
		Hostname:      hostname,
		Url:           hostname + r.URL.Path,
	}

	// Log the request details to stdout
	log.Printf("IP: %s, UserAgent: %s, RequestMethod: %s, Referrer: %s, X-Forwarded-For: %s, Hostname: %s, URL: %s",
		requestInfo.IP, requestInfo.UserAgent, requestInfo.RequestMethod, requestInfo.Referrer, requestInfo.XForwardedFor, requestInfo.Hostname, requestInfo.Url)

	prettyJson, err := json.MarshalIndent(requestInfo, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestInfo.PrettyJson = string(prettyJson)

	if isBrowser(userAgent) && r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.New("request_info").Parse(requestInfoTemplate))
		tmpl.Execute(w, requestInfo)
		return
	}

	switch r.URL.Path {
	case "/":
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(ip))
	case "/ip":
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(ip))
	case "/ua":
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(userAgent))
	case "/json":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requestInfo)
	default:
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.New("request_info").Parse(requestInfoTemplate))
		tmpl.Execute(w, requestInfo)
	}
}

func isBrowser(userAgent string) bool {
	browsers := []string{"Mozilla", "Chrome", "Safari", "Firefox", "Edge"}
	for _, browser := range browsers {
		if strings.Contains(userAgent, browser) {
			return true
		}
	}
	return false
}
