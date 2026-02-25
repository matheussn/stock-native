# Planejamento – Sistema Local de Controle de Estoque Comunitário
## Tecnologia: Go + Wails + SQLite (Execução 100% Local)

---

## 1. Objetivo do Sistema

Criar um **aplicativo desktop local**, desenvolvido em **Go + Wails**, para controle de estoque de um centro espírita.

O sistema deve:
- Rodar **inteiramente na máquina local** (sem servidor externo).
- Persistir dados em **banco local (SQLite em arquivo)**.
- Ser leve, simples e multiplataforma (Windows / Linux / macOS).
- Substituir o sistema Nex, removendo conceitos comerciais desnecessários.
- Eliminar gambiarras como “fornecedores fictícios”.

---

## 2. Arquitetura Geral

### 2.1 Tipo de Aplicação
- **Desktop application**
- Backend em Go
- Frontend em HTML/CSS/JS (renderizado via Wails WebView)

### 2.2 Execução
- Aplicação inicia via executável local
- Interface gráfica abre automaticamente (janela nativa do SO)
- Nenhuma dependência de navegador externo
- Nenhum serviço externo ou cloud

### 2.3 Persistência
- Banco de dados: **SQLite**
- Armazenado como arquivo local (ex: `data/app.db`)
- Migrações automáticas ao iniciar o app

### 2.4 Backup
- Backup = cópia do arquivo `.db`
- Funcionalidade no app:
  - Exportar backup
  - Importar backup (substituição controlada)

---

## 3. Conceito Central do Sistema

O sistema é **baseado exclusivamente em movimentações de estoque**.

Não existem:
- Vendas
- Preços
- Caixa
- PDV
- Fornecedores comerciais

Toda alteração no estoque ocorre por:
- **Entradas**
- **Saídas**

Cálculo do estoque atual:

estoque = soma(entradas) - soma(saídas)

---

## 4. Modelo Conceitual de Dados

### 4.1 Categoria
Agrupa produtos.

Campos:
- id
- nome
- ativo

---

### 4.2 Produto
Item controlado em estoque.

Campos:
- id
- nome
- categoria_id
- unidade (ex: unidade, kg, pacote)
- ativo

---

### 4.3 Origem
Define **de onde veio um produto em uma entrada**.

Exemplos:
- Doação anônima
- Campanha do agasalho 2025
- Evento beneficente
- Parceria com mercado
- Compra emergencial

Campos:
- id
- nome
- tipo (opcional)
- ativo

---

### 4.4 Destino
Define **para onde um produto foi em uma saída**.

Exemplos:
- Kits mensais
- Famílias assistidas
- Evento interno
- Doação externa
- Consumo interno

Campos:
- id
- nome
- tipo (opcional)
- ativo

---

### 4.5 Movimentação
Registro central do sistema.

Campos:
- id
- tipo ("entrada" | "saida")
- produto_id
- quantidade
- data
- origem_id (obrigatório se tipo = entrada)
- destino_id (obrigatório se tipo = saída)
- observacao (opcional)

Regras:
- Nunca permitir origem_id e destino_id preenchidos simultaneamente
- Quantidade sempre positiva
- Não permitir saída se estoque atual < quantidade solicitada

---

## 5. Funcionalidades do MVP

### 5.1 Cadastros (CRUD)
- Categorias
- Produtos
- Origens
- Destinos

Todos devem permitir:
- Criar
- Editar
- Desativar (soft delete)

---

### 5.2 Entrada de Estoque
Fluxo:
- Selecionar produto
- Informar quantidade
- Selecionar origem
- Informar data
- Observação opcional
- Persistir movimentação

---

### 5.3 Saída de Estoque
Fluxo:
- Selecionar produto
- Informar quantidade
- Selecionar destino
- Informar data
- Observação opcional
- Validar estoque
- Persistir movimentação

---

### 5.4 Estoque Atual
- Listar produtos
- Mostrar quantidade atual (calculada)
- Filtro por categoria
- Indicador visual de estoque baixo (opcional)

---

### 5.5 Relatórios Básicos
- Entradas por período
- Entradas agrupadas por origem
- Saídas por período
- Saídas agrupadas por destino
- Histórico geral de movimentações

---

## 6. Funcionalidades Fora de Escopo (Explicitamente NÃO implementar)

- Preços
- Vendas
- Caixa
- PDV
- Fornecedores comerciais
- Emissão fiscal
- Integrações bancárias
- Login externo ou autenticação online

---

## 7. Organização do Código (Sugestão)

### 7.1 Backend (Go)
- domain/
  - entities
  - rules
- repository/
  - sqlite
- services/
- migrations/
- app.go (Wails bootstrap)

### 7.2 Frontend
- UI simples e funcional
- Foco em clareza para voluntários
- Formulários diretos
- Tabelas simples
- Pouco ou nenhum JavaScript complexo

---

## 8. Validações Importantes

- Não permitir estoque negativo
- Não permitir movimentação sem produto
- Entrada exige origem
- Saída exige destino
- Datas válidas
- Quantidade > 0

---

## 9. Fases de Implementação

### Fase 1 – Infra Local
- Projeto Wails
- SQLite local
- Migrações automáticas

### Fase 2 – Domínio
- Entidades
- Regras de negócio
- Repositórios

### Fase 3 – UI
- Cadastros
- Movimentações
- Estoque

### Fase 4 – Relatórios e Backup
- Relatórios básicos
- Exportar / Importar banco

---

## 10. Objetivo Final

Entregar um sistema:
- Local
- Gratuito
- Leve
- Multiplataforma
- Adequado à realidade de um centro espírita
- Sem dependência de software comercial ou cloud