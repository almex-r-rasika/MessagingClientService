package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)


func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Get("/", app.saveImportData)
	mux.Get("/v2/msgb/import_data/{page}", app.getImportDataHistory)
	mux.Delete("/v2/msgb/import_data/{importId}", app.deleteImportDataHistory)
    // TODO Implement GET API
	mux.Put("/v2/msgb/bulk_message", app.saveBulkMessage)
	mux.Get("/v2/msgb/bulk_message/{page}", app.getBulkMessages)
	mux.Get("/v2/msgb/bulk_message/{bulkMessageId}/{page}", app.getBulkMessage)

	mux.Get("/v2/msgb/message_template", app.getMessageTemplates)
	mux.Put("/v2/msgb/message_template/{templateType}", app.saveMessageTemplate)
	mux.Get("/v2/msgb/message_template/{templateType}/{messageTemplateId}", app.getMessageTemplate)
	mux.Post("/v2/msgb/message_template/{templateType}/{messageTemplateId}", app.updateMessageTemplate)
	mux.Delete("/v2/msgb/message_template/{templateType}/{messageTemplateId}", app.deleteMessageTemplate)
	
	return mux
}