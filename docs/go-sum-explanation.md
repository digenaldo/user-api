# Explicação do `go.sum` (pt-BR)

O arquivo `go.sum` é gerado automaticamente pelo sistema de módulos do Go e
contem somas (hashes) criptográficas das versões de módulos que o seu
projeto usa. Essas somas servem para verificar a integridade das dependências
baixadas — garantindo que o código recebido seja exatamente o mesmo que foi
publicado pelo autor.

Por que `go.sum` existe
- Segurança: impede ataques "supply-chain" onde um pacote remoto poderia ser
  alterado entre o momento em que você adicionou a dependência e a sua
  máquina baixa o conteúdo.
- Reprodutibilidade: garante que todos os desenvolvedores/CI baixem exatamente
  as mesmas versões exatas de dependências.

Práticas recomendadas
- Nunca edite `go.sum` manualmente. Ele é mantido pelo comando `go`.
- Sempre comite `go.sum` junto com `go.mod` no repositório — é importante
  para que outros desenvolvedores e sistemas de CI/verificação possam validar
  as somas.

Como atualizar (de forma segura)
- `go get <module>@<version>` para adicionar/atualizar versões específicas.
- `go mod tidy` para limpar dependências não usadas e atualizar `go.sum`.
- `go mod download` para forçar o download das dependências e popular/atualizar
  o `go.sum` localmente.

Verificação de integridade
- `go mod verify` verifica o diretório `GOMODCACHE` contra as somas em
  `go.sum`.

Resumo rápido
- `go.sum` garante segurança e reprodutibilidade.
- Não edite manualmente; use os comandos da ferramenta `go` para manter.
