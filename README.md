# Sistema de Estoque - Centro Espirita

Aplicacao desktop local (Go + Wails + SQLite) para controle de estoque, cestas, familias e movimentacoes.

## Requisitos

- Go 1.23+
- Node.js 18+
- Wails CLI 2.x

## Desenvolvimento

```bash
wails dev
```

Ou, usando o `Makefile`:

```bash
make dev
```

## Build para distribuicao

```bash
wails build
```

Ou, usando o `Makefile`:

```bash
make build
make release
```

No Windows, o executavel sera gerado na pasta `build/bin` (ou caminho equivalente definido pelo Wails).

Para gerar tambem o instalador NSIS no Windows:

```bash
wails build -nsis
```

Ou:

```bash
make installer
```

Com a configuracao atual da versao `1.0.0`, os artefatos esperados sao:

- executavel: `build/bin/Obreiro.exe`
- instalador: `build/bin/Obreiro-amd64-installer.exe`

## Banco de dados local

O arquivo SQLite e criado automaticamente no diretorio de configuracao do usuario:

- Windows: `%APPDATA%\stock\stock.db`
- Linux: `~/.config/stock/stock.db`
- macOS: `~/Library/Application Support/stock/stock.db`

Na abertura da conexao, a aplicacao configura o SQLite com:

- `foreign_keys = ON`
- `journal_mode = WAL`
- `busy_timeout = 5000`
- `synchronous = NORMAL`

## Backup e restauracao

- `BackupDatabase(destinationPath)` cria uma copia consistente do banco ativo em outro arquivo `.db`
- `RestoreDatabase(sourcePath)` substitui o banco ativo pelo arquivo informado
- antes de restaurar, o banco atual e salvo automaticamente ao lado do original com o padrao `stock.before-restore-YYYYMMDD-HHMMSS.db`
- o arquivo restaurado passa pelas migrations na volta da aplicacao para garantir compatibilidade com a versao atual

## Observacoes operacionais

- nunca use o mesmo caminho do banco ativo como destino do backup
- nunca tente restaurar usando o proprio `stock.db` em uso como origem
- mantenha pelo menos uma copia externa do backup fora da pasta padrao da aplicacao

## Logs operacionais

O aplicativo grava eventos operacionais em arquivo, com rotacao simples quando o log atual ultrapassa 5 MB:

- Windows: `%APPDATA%\\stock\\logs\\stock.log`
- Linux: `~/.config/stock/logs/stock.log`
- macOS: `~/Library/Application Support/stock/logs/stock.log`

Os logs incluem:

- startup e shutdown da aplicacao
- inicializacao do banco e falhas de startup
- backup e restauracao do banco
- operacoes de escrita relevantes, como cadastros, edicoes, ativacao/desativacao e movimentacoes

## Fluxo de operacao

- Cadastre Obras Assistenciais
- Cadastre Produtos e Variacoes
- Monte Cestas (Grupos)
- Cadastre Familias e atribua cestas
- Registre Movimentacoes de entrada/saida
- Acompanhe painel de Estoque e Cobertura

## Automacao com Makefile

Os alvos disponiveis no `Makefile` da raiz sao:

- `make deps`: instala dependencias do frontend
- `make frontend-build`: gera o build do frontend
- `make test`: executa os testes Go
- `make dev`: sobe o app em desenvolvimento com Wails
- `make build`: gera build local sem empacotar
- `make release`: gera build de distribuicao
- `make installer`: gera build com instalador NSIS
- `make clean`: remove artefatos locais de build
