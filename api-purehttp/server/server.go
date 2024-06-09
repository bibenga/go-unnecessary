// go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=server.cfg.yaml ../api.yaml
// go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.cfg.yaml ../api.yaml
// go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
//
//+go:generate oapi-codegen --config=model.cfg.yaml  ../../api/api.yaml
//+go:generate oapi-codegen --config=client.cfg.yaml ../../api/api.yaml
//+go:generate oapi-codegen --config=server.cfg.yaml ../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=model.cfg.yaml ../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=client.cfg.yaml ../../api/api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=server.cfg.yaml ../../api/api.yaml

package server
