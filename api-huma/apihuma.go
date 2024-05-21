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

type ItemIn struct {
	Name string `json:"name" maxLength:"150"`
}

type ItemOut struct {
	ItemId uint64 `json:"id" minimum:"1"`
	Name   string `json:"name"`
}

type GetItemsInput struct {
	PaginationParams
	Search string `query:"q" maxLength:"200"`
}

type GetItemsOutput struct {
	Body struct {
		PaginationReponseParams
		Items []ItemOut `json:"items"`
	}
}

type PostItemInput struct {
	Body struct {
		ItemIn
	}
}

type PostItemOutput struct {
	Body struct {
		ItemOut
	}
}

type PutItemInput struct {
	ItemId uint64 `path:"itemId" minimum:"1"`
	Body   struct {
		ItemIn
	}
}

type PutItemOutput struct {
	Body struct {
		ItemOut
	}
}

type DelItemInput struct {
	ItemId uint64 `path:"itemId" minimum:"1"`
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

func getMe(ctx context.Context, input *struct{}) (*GetMeOutput, error) {
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
	slog.Info("> getItems", "input", input)
	resp := &GetItemsOutput{}
	resp.Body.Page = input.Page
	resp.Body.Items = []ItemOut{
		{ItemId: 1, Name: "item1"},
		{ItemId: 2, Name: "item2"},
	}
	slog.Info("< getItems", "resp", resp)
	return resp, nil
}

func createItem(ctx context.Context, input *PostItemInput) (*PostItemOutput, error) {
	slog.Info("> createItem", "input", input)
	resp := &PostItemOutput{}
	resp.Body.ItemId = 3
	resp.Body.Name = input.Body.Name
	slog.Info("< createItem", "resp", resp)
	return resp, nil
}

func updateItem(ctx context.Context, input *PutItemInput) (*PutItemOutput, error) {
	slog.Info("> updateItem", "input", input)
	resp := &PutItemOutput{}
	resp.Body.ItemId = input.ItemId
	resp.Body.Name = input.Body.Name
	slog.Info("< updateItem", "resp", resp)
	return resp, nil
}

func deleteItem(ctx context.Context, input *DelItemInput) (*struct{}, error) {
	slog.Info("> deleteItem", "input", input)
	return nil, nil
}

func registerApis(api huma.API) {
	defaultSecurity := []map[string][]string{{}, {"bearer": {}}, {"apiKey": {}}, {"http": {}}}
	defaultErrors := []int{401, 403}

	// Auth

	huma.Register(api, huma.Operation{
		OperationID: "login",
		Tags:        []string{"auth"},
		Summary:     "Login",
		Method:      http.MethodPost,
		Path:        "/api/auth/login",
		Errors:      defaultErrors,
		Deprecated:  true,
	}, login)
	huma.Register(api, huma.Operation{
		OperationID: "logout",
		Tags:        []string{"auth"},
		Summary:     "Logout",
		Method:      http.MethodPost,
		Path:        "/api/auth/logout",
		Security:    defaultSecurity,
		Deprecated:  true,
	}, logout)

	huma.Register(api, huma.Operation{
		OperationID: "get-me",
		Tags:        []string{"auth"},
		Summary:     "Get an information of a current logged user ",
		Method:      http.MethodGet,
		Path:        "/api/auth/me",
		Errors:      defaultErrors,
		Security:    defaultSecurity,
	}, getMe)

	// Items
	huma.Register(api, huma.Operation{
		OperationID: "get-items",
		Tags:        []string{"items"},
		Summary:     "Get Items",
		Method:      http.MethodGet,
		Path:        "/api/items",
		Errors:      defaultErrors,
		Security:    defaultSecurity,
	}, getItems)
	huma.Register(api, huma.Operation{
		OperationID:   "post-item",
		Tags:          []string{"items"},
		Summary:       "Create Item",
		Method:        http.MethodPost,
		Path:          "/api/items",
		Errors:        defaultErrors,
		Security:      defaultSecurity,
		DefaultStatus: http.StatusCreated,
	}, createItem)
	huma.Register(api, huma.Operation{
		OperationID: "put-item",
		Tags:        []string{"items"},
		Summary:     "Update Item",
		Method:      http.MethodPut,
		Path:        "/api/items/{itemId}",
		Errors:      defaultErrors,
		Security:    defaultSecurity,
	}, updateItem)
	huma.Register(api, huma.Operation{
		OperationID:   "del-item",
		Tags:          []string{"items"},
		Summary:       "Delete Item",
		Method:        http.MethodDelete,
		Path:          "/api/items/{itemId}",
		Errors:        defaultErrors,
		Security:      defaultSecurity,
		DefaultStatus: http.StatusNoContent,
	}, deleteItem)
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

		config := huma.DefaultConfig("My Huma API", "1.0.0")
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"bearer": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
			"apiKey": {Type: "apiKey", In: "header", Name: "X-API-KEY"},
			"http":   {Type: "http", Scheme: "basic"},
		}

		api = humafiber.New(app, config)
		api.UseMiddleware(NewAuthMiddleware(api))
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
