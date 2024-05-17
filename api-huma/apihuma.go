package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/cobra"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI.
type Options struct {
	Port  int  `help:"Port to listen on" default:"8000"`
	Debug bool `debug:"Debug logging" default:false`
}

// GreetingInput represents the greeting operation request.
type GreetingInput struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	Q    string `query:"q" maxLength:"30" example:"q" doc:"Q"`
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"Author of the review"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message"`
	}
}

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func NewAuthMiddleware(api huma.API) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		auths := map[string]bool{
			"http":   false,
			"bearer": false,
			"apiKey": false,
		}
		for _, opScheme := range ctx.Operation().Security {
			for name := range opScheme {
				auths[name] = true
			}
		}

		slog.Info(">", "auths", auths)
		next(ctx)
		slog.Info("<")
	}
}

func main() {
	// Store the API so we can access it from other commands later.
	var api huma.API

	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		router := chi.NewMux()
		router.Use(middleware.RequestID)
		router.Use(middleware.Logger)
		// router.Use(middleware.RequestLogger(
		// 	&middleware.DefaultLogFormatter{
		// 		Logger:  log.New(logger, "", 0),
		// 		NoColor: true,
		// 	}))
		router.Use(middleware.Recoverer)
		router.Use(middleware.RealIP)
		router.Use(middleware.NoCache)
		router.Use(middleware.Heartbeat("/ping"))

		config := huma.DefaultConfig("My API", "1.0.0")
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"bearer": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
			"apiKey": {Type: "apiKey", In: "header", Name: "X-API-KEY"},
			"http":   {Type: "http", Scheme: "basic"},
		}
		api = humachi.New(router, config)

		api.UseMiddleware(NewAuthMiddleware(api))

		// apis
		huma.Register(api, huma.Operation{
			OperationID: "get-greeting",
			Summary:     "Get a greeting",
			Method:      http.MethodPost,
			Path:        "/api/greeting/{name}",
			Security:    []map[string][]string{{}, {"http": {}}, {"apiKey": {}}},
		}, func(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
			slog.Info("-")
			resp := &GreetingOutput{}
			resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
			return resp, nil
		})

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			slog.Info(fmt.Sprintf("Ready on port %[1]d: http://127.0.0.1:%[1]d/docs", options.Port))
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Add a command to print the OpenAPI spec.
	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := api.OpenAPI().YAML()
			fmt.Println(string(b))
		},
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
