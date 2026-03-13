SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -ExecutionPolicy Bypass -Command

APP_NAME := Obreiro
FRONTEND_DIR := frontend
GO_CACHE := $(CURDIR)\.gocache
WAILS := wails
NPM := npm

.PHONY: help deps frontend-install frontend-build test dev build release installer clean

help:
	@Write-Host ''
	@Write-Host 'Alvos disponiveis:'
	@Write-Host '  make deps             Instala dependencias do frontend'
	@Write-Host '  make frontend-build   Gera o build do frontend'
	@Write-Host '  make test             Executa os testes Go'
	@Write-Host '  make dev              Inicia o Wails em modo desenvolvimento'
	@Write-Host '  make build            Gera build local do app sem empacotar'
	@Write-Host '  make release          Gera build de distribuicao do app'
	@Write-Host '  make installer        Gera build com instalador NSIS'
	@Write-Host '  make clean            Remove artefatos locais de build'
	@Write-Host ''

deps: frontend-install

frontend-install:
	@Set-Location $(FRONTEND_DIR); $(NPM) install

frontend-build:
	@Set-Location $(FRONTEND_DIR); $(NPM) run build

test:
	@$env:GOCACHE='$(GO_CACHE)'; go test ./...

dev:
	@$(WAILS) dev

build: test frontend-build
	@$(WAILS) build -s -nopackage

release: test frontend-build
	@$(WAILS) build

installer: test frontend-build
	@$(WAILS) build -nsis

clean:
	@if (Test-Path '$(FRONTEND_DIR)\dist') { Remove-Item '$(FRONTEND_DIR)\dist' -Recurse -Force }
	@if (Test-Path 'build\bin') { Get-ChildItem 'build\bin' -Filter '*.exe' | Remove-Item -Force }
	@if (Test-Path 'build\bin') { Get-ChildItem 'build\bin' -Filter '*.msi' | Remove-Item -Force }
	@if (Test-Path 'build\bin') { Get-ChildItem 'build\bin' -Filter '*installer*' | Remove-Item -Recurse -Force }
	@if (Test-Path '.gocache') { Remove-Item '.gocache' -Recurse -Force }
