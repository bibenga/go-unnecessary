// go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
//
//+go:generate oapi-codegen --config=model.cfg.yaml  ../../api/api.yaml
//+go:generate oapi-codegen --config=client.cfg.yaml ../../api/api.yaml
//+go:generate oapi-codegen --config=server.cfg.yaml ../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=model.cfg.yaml ../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=client.cfg.yaml ../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=server.cfg.yaml ../../api/api.yaml

package server
