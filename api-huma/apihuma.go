package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	// _ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI.
type Options struct {
	Port  int  `help:"Port to listen on" default:"8000"`
	Debug bool `help:"Debug logging" default:"false"`
}

type PaginationParams struct {
	Page     int `query:"page" minimum:"1" default:"1"`
	PageSize int `query:"page-size" minimum:"10" maximum:"100" default:"50"`
}

type PaginationReponseParams struct {
	Page int `json:"page" minimum:"1" default:"1"`
}

type GetMeInput struct {
}

type GetMeOutput struct {
	Body struct {
		UserId   uint64 `json:"user_id" minimum:"1"  doc:"Current user id"`
		Username string `json:"username" doc:"Current username"`
	}
}

type GetLoginInput struct {
	Body struct {
		Username string `doc:"Username" maxLength:"150"`
		Password string `doc:"Password" maxLength:"150"`
	}
}

type GetLoginOutput struct {
	SetCookie []*http.Cookie `header:"Set-Cookie"`
}

type GetItemsInput struct {
	PaginationParams
	Search string `query:"q" maxLength:"200"`
}

type Item struct {
	ItemId uint64 `json:"id" minimum:"1"`
	Name   string `json:"name"`
}

type GetItemsOutput struct {
	Body struct {
		PaginationReponseParams
		Items []Item `json:"items"`
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

func getMe(ctx context.Context, input *GetMeInput) (*GetMeOutput, error) {
	slog.Info("-")
	resp := &GetMeOutput{}
	resp.Body.UserId = 1
	resp.Body.Username = "admin"
	return resp, nil
}

func login(ctx context.Context, input *GetLoginInput) (*GetLoginOutput, error) {
	slog.Info("> login", "input", input)

	slog.Info("- ", "ctx", ctx)
	fctx := ctx.(*fasthttp.RequestCtx)
	slog.Info("- ", "fasthttp.RequestCtx", fctx)
	// fctx.appSetUserValue()

	if input.Body.Username != "1" && input.Body.Password != "1" {
		return nil, huma.Error401Unauthorized("Who are you?")
	}
	resp := &GetLoginOutput{}
	resp.SetCookie = []*http.Cookie{
		{
			Domain:   "example.com",
			Name:     "_my_session",
			HttpOnly: true,
			Value:    input.Body.Username,
			Expires:  time.Now().Add(5 * time.Minute),
		},
	}
	slog.Info("< login", "resp", resp)
	return resp, nil
}

func logout(ctx context.Context, input *struct{}) (*struct{}, error) {
	slog.Info("logout", "input", input)
	return nil, nil
}

func getItems(ctx context.Context, input *GetItemsInput) (*GetItemsOutput, error) {
	slog.Info("getItems", "input", input)
	resp := &GetItemsOutput{}
	resp.Body.Page = input.Page
	resp.Body.Items = []Item{
		Item{ItemId: 1, Name: "item1"},
		Item{ItemId: 2, Name: "item2"},
	}
	return resp, nil
}

func registerApis(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "login",
		Tags:        []string{"auth"},
		Summary:     "Login",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Errors:      []int{401, 403},
	}, login)

	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Tags:        []string{"auth"},
		Summary:     "Logout",
		Method:      http.MethodPost,
		Path:        "/api/auth/logout",
		Security:    []map[string][]string{{}, {"bearer": {}}, {"apiKey": {}}, {"http": {}}},
	}, logout)

	huma.Register(api, huma.Operation{
		OperationID: "get-me",
		Tags:        []string{"user"},
		Summary:     "Get a user information",
		Method:      http.MethodGet,
		Path:        "/api/me",
		Errors:      []int{401, 403},
		Security:    []map[string][]string{{}, {"bearer": {}}, {"apiKey": {}}, {"http": {}}},
	}, getMe)

	huma.Register(api, huma.Operation{
		OperationID: "get-items",
		Tags:        []string{"items"},
		Summary:     "Get Items",
		Method:      http.MethodGet,
		Path:        "/api/items",
		Security:    []map[string][]string{{}, {"bearer": {}}, {"apiKey": {}}, {"http": {}}},
	}, getItems)
}

func main() {
	// Store the API so we can access it from other commands later.
	var api huma.API

	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Create a new router & API
		app := fiber.New()
		app.Use(recover.New())
		app.Use(logger.New())
		app.Use(requestid.New())
		app.Use(helmet.New(helmet.Config{
			CrossOriginEmbedderPolicy: "cross-origin",
		}))
		app.Use(etag.New())

		// app.Get("/metrics", monitor.New())
		app.Use(redirect.New(redirect.Config{
			Rules: map[string]string{
				"/": "/docs",
			},
			StatusCode: 302,
		}))

		config := huma.DefaultConfig("My API", "1.0.0")
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"bearer": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
			"apiKey": {Type: "apiKey", In: "header", Name: "X-API-KEY"},
			"http":   {Type: "http", Scheme: "basic"},
		}

		api = humafiber.New(app, config)

		api.UseMiddleware(NewAuthMiddleware(api))

		// apis
		// huma.Register(api, huma.Operation{
		// 	OperationID: "get-greeting",
		// 	Summary:     "Get a greeting",
		// 	Method:      http.MethodPost,
		// 	Path:        "/api/greeting/{name}",
		// 	Security:    []map[string][]string{{}, {"http": {}}, {"apiKey": {}}},
		// }, func(ctx context.Context, input *GreetingInput) (*GreetingOutput, error) {
		// 	slog.Info("-")
		// 	resp := &GreetingOutput{}
		// 	resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		// 	return resp, nil
		// })
		registerApis(api)

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			slog.Info(
				fmt.Sprintf("Ready on port %[1]d: http://127.0.0.1:%[1]d/docs", options.Port),
				"port", options.Port,
				"debug", options.Debug,
			)
			if err := app.Listen(fmt.Sprintf(":%d", options.Port)); err != nil {
				panic(err)
			}
		})

		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			app.ShutdownWithContext(ctx)
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
