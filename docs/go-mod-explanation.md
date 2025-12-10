# Explicação do `go.mod` (pt-BR)

O arquivo `go.mod` é o manifesto do módulo Go do seu projeto. Ele declara o
nome do módulo, a versão mínima do Go que o projeto usa, e as dependências
diretas (e indiretas) necessárias para compilar a aplicação.

Campos principais
- `module <path>`: identifica o módulo. Em projetos locais simples pode ser
  um nome curto como `user-api`, mas em projetos publicados normalmente é o
  caminho completo (ex.: `github.com/usuario/projeto`).
- `go <version>`: define a versão do Go que influencia resolução de módulos e
  comportamentos específicos do compilador.
- `require`: lista de dependências com versões. Dependências marcadas como
  `// indirect` são trazidas por outras dependências e não usadas diretamente
  no seu código.

Comandos úteis
- `go get <module>@<version>`: adiciona/atualiza dependências.
- `go mod tidy`: remove dependências não usadas e atualiza `go.mod` e `go.sum`.
- `go mod download`: faz o download de todas as dependências declaradas.

Boas práticas
- Commitar `go.mod` e `go.sum` juntos no repositório para manter reprodutibilidade.
- Evitar editar `go.mod` manualmente — prefira os comandos da ferramenta `go`.
- Use `go mod tidy` antes de commits para manter os arquivos limpos.

Exemplo rápido
1. Para adicionar uma dependência: `go get github.com/pkg/example@v1.2.3`
2. Para limpar: `go mod tidy`
