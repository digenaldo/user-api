module user-api

// go.mod: arquivo de configuração de dependências do Go
// - Não edite o bloco `require` manualmente sem entender `go mod`.
// - Use `go get`, `go mod tidy` ou `go mod vendor` para alterar dependências.
// - O `go` define a versão mínima do compilador/semântica de módulos.

go 1.22

require (
	github.com/go-chi/chi/v5 v5.2.3
	go.mongodb.org/mongo-driver v1.17.6
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/swaggo/files v0.0.0-20220610200504-28940afbdbfe // indirect
	github.com/swaggo/http-swagger v1.3.4 // indirect
	github.com/swaggo/swag v1.8.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// Explicação rápida:
// - `module` define o caminho do módulo local.
// - `go 1.22` ajusta regras de compatibilidade do compilador e módulos.
// - Depêndencias diretas aparecem no primeiro bloco `require`.
// - Depêndencias marcadas `// indirect` são trazidas por outras libs.
// Para atualizar/limpar dependências use: `go mod tidy`.
