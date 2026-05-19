package main

import (
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/microsoft/ApplicationInsights-Go/appinsights"    
    "crm-mock/db"
    "crm-mock/handler"
)

type statusWriter struct {
    http.ResponseWriter
    status int
}

func (w *statusWriter) WriteHeader(status int) {
    w.status = status
    w.ResponseWriter.WriteHeader(status)
}

func appInsightsMiddleware(client appinsights.TelemetryClient) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
            start := time.Now()
            next.ServeHTTP(sw, r)
            duration := time.Since(start)

            reqTelemetry := appinsights.NewRequestTelemetry(r.Method, r.URL.String(), duration, strconv.Itoa(sw.status))
            reqTelemetry.Name = r.Method + " " + r.URL.Path
            reqTelemetry.Source = r.RemoteAddr
            reqTelemetry.Properties["host"] = r.Host
            reqTelemetry.Properties["path"] = r.URL.Path
            client.Track(reqTelemetry)
        })
    }
}

func main() {
    conn, err := db.NewConnection()
    if err != nil {
        log.Fatal("failed to connect db:", err)
    }
    defer conn.Close()

    for i := 0; i < 10; i++ {
        if err = conn.Ping(); err == nil { break }
        log.Printf("waiting for db... attempt %d", i+1)
        time.Sleep(2 * time.Second)
    }

    ikey := os.Getenv("APPINSIGHTS_INSTRUMENTATIONKEY")
    var aiClient appinsights.TelemetryClient
    if ikey != "" {
        aiClient = appinsights.NewTelemetryClient(ikey)
        defer aiClient.Channel().Flush()
    }
    
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.SetHeader("Content-Type", "application/json"))

    if aiClient != nil {
        r.Use(appInsightsMiddleware(aiClient))
    }
    
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok","service":"crm-mock"}`))
    })

    accountH     := handler.NewAccountHandler(conn)
    opportunityH := handler.NewOpportunityHandler(conn)
    leadH        := handler.NewLeadHandler(conn)
    contactH     := handler.NewContactHandler(conn)

    // CRM APIs (consumed by SAP)
    r.Route("/crm", func(r chi.Router) {
        // Accounts
        r.Get("/accounts/{id}",          accountH.GetByID)
        r.Post("/accounts",              accountH.Create)
        r.Put("/accounts/{id}",           accountH.Update)
        r.Delete("/accounts/{id}",        accountH.Delete)

        // Opportunities
        r.Post("/opportunities",          opportunityH.Create)
        r.Put("/opportunities/{id}",      opportunityH.Update)
        r.Delete("/opportunities/{id}",   opportunityH.Delete)

        // Leads
        r.Get("/leads/{id}",             leadH.GetByID)
        r.Post("/leads",                 leadH.Create)
        r.Patch("/leads/{id}",           leadH.Update)
        r.Delete("/leads/{id}",          leadH.Delete)

        // Contacts
        r.Get("/contacts/{id}",          contactH.GetByID)
        r.Post("/contacts",              contactH.Create)
        r.Delete("/contacts/{id}",       contactH.Delete)
    })

    port := os.Getenv("APP_PORT")
    if port == "" { port = "8082" }

    log.Printf("CRM Mock running on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}
