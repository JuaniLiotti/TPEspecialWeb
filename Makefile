# ==============================================================
#  Italpiel S.A. - Automatización del proyecto
#
#  make run-app   ->  tareas previas + ejecuta la app + tareas posteriores
#  make           ->  lista todos los comandos
# ==============================================================

SHELL := /bin/bash

CONTAINER_NAME := postgres-tp2
IMAGE_NAME     := tp-postgres
BINARY         := api-server
PORT           := 8080
DB_USER        := postgres
DB_NAME        := italpiel_db
UPLOADS_DIR    := uploads

.DEFAULT_GOAL := help

.PHONY: help run-app run pre-run post-run clean \
        generate run-sqlc templ build-app \
        test pre-test run-tests post-tests \
        test-api pre-test-api start-app wait-app run-api-tests post-test-api clean-app \
        db build-db-image clean-db start-db wait-db clean-uploads

# Los targets con "## descripción" aparecen en "make help"
help: ## Muestra esta ayuda
	@echo "Italpiel S.A. - comandos disponibles:"
	@echo ""
	@grep -hE '^[a-zA-Z_-]+:.*## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36mmake %-10s\033[0m %s\n", $$1, $$2}'
	@echo ""


# ==============================================================
# EJECUTAR LA APLICACIÓN (TP5)
# ==============================================================

# 1) Tareas previas  2) Ejecutar la app  3) Tareas posteriores
# La app corre en primer plano. Al cortarla con Ctrl+C, el "trap" atrapa la
# señal y corre las tareas posteriores antes de salir (también si la app se cae).
run-app: pre-run ## Tareas previas, ejecuta la app en http://localhost:8080 y limpia al cortar (Ctrl+C)
	@echo ""
	@echo "=============================================================="
	@echo "  Italpiel corriendo en http://localhost:$(PORT)"
	@echo "  Ctrl+C para cortar (borra contenedores y volúmenes al salir)"
	@echo "=============================================================="
	@trap '$(MAKE) --no-print-directory post-run; exit 0' INT TERM; \
		./$(BINARY); \
		$(MAKE) --no-print-directory post-run

run: run-app ## Alias de run-app

# Tareas previas: sqlc + templ, compilación, borrar contenedores y volúmenes,
# levantar el contenedor y esperar a que la BBDD esté lista
pre-run: build-app db

# Tareas posteriores: borrar contenedores y volúmenes, y las imágenes subidas
# (la BBDD arranca vacía en cada corrida, así que esas imágenes quedarían huérfanas)
post-run:
	@echo ""
	@echo ">> Tareas posteriores..."
	@$(MAKE) --no-print-directory clean-db clean-uploads
	@echo ">> Listo, entorno limpio."

clean: clean-app ## Limpia todo: servidor en segundo plano, BBDD, imágenes subidas y binario
	@$(MAKE) --no-print-directory clean-db clean-uploads
	@rm -f $(BINARY)


# ==============================================================
# GENERACIÓN DE CÓDIGO Y COMPILACIÓN
# (lo generado no se versiona: ver .gitignore)
# ==============================================================

build-app: generate ## Genera el código (sqlc + templ) y compila el servidor
	@echo ">> Compilando el proyecto..."
	go build -o $(BINARY) .

generate: run-sqlc templ ## Genera db/sqlc/ con sqlc y views/*_templ.go con templ

run-sqlc:
	@echo ">> Generando código con sqlc..."
	sqlc generate

# templ está declarado como "tool" en go.mod: no hace falta instalarlo aparte
templ:
	@echo ">> Generando componentes con templ..."
	go tool templ generate


# ==============================================================
# TP2: tests de la base de datos
# ==============================================================

test: pre-test run-tests post-tests ## Corre los tests de la BBDD (go test) y limpia

pre-test: generate db

run-tests:
	@echo ">> Ejecutando los tests de la BBDD..."
	go test -v ./...

# Make ejecuta cada target una sola vez por corrida, y clean-db ya se usó en db.
# $(MAKE) lanza una ejecución nueva para que la limpieza final sí se ejecute.
post-tests:
	@$(MAKE) --no-print-directory clean-db


# ==============================================================
# TP3: tests de la API JSON (requests.sh)
# ==============================================================

test-api: pre-test-api run-api-tests post-test-api ## Corre los tests de la API (requests.sh) y limpia

pre-test-api: build-app db start-app wait-app

start-app:
	@echo ">> Levantando el servidor en segundo plano..."
	./$(BINARY) & echo $$! > app.pid

# Espera a que el servidor responda (máximo 30 s) en vez de un sleep fijo
wait-app:
	@echo ">> Esperando a que el servidor esté listo..."
	@i=0; until curl -s -o /dev/null http://localhost:$(PORT)/; do \
		i=$$((i+1)); if [ $$i -ge 60 ]; then echo "!! El servidor no respondió"; exit 1; fi; \
		sleep 0.5; \
	done

run-api-tests:
	@echo ">> Ejecutando pruebas de los endpoints..."
	chmod +x requests.sh
	./requests.sh

post-test-api: clean-app
	@$(MAKE) --no-print-directory clean-db clean-uploads

clean-app:
	@if [ -f app.pid ]; then \
		echo ">> Cerrando el servidor en segundo plano..."; \
		kill -9 $$(cat app.pid) 2>/dev/null || true; rm -f app.pid; \
	fi


# ==============================================================
# BASE DE DATOS (PostgreSQL en Docker)
# ==============================================================

# Borra el contenedor viejo (con sus volúmenes) y levanta uno nuevo con schema.sql
db: clean-db build-db-image start-db wait-db ## Levanta una BBDD nueva y vacía en Docker

build-db-image:
	@echo ">> Construyendo la imagen de la BBDD..."
	docker build -q -t $(IMAGE_NAME) .

# -v borra también el volumen anónimo donde Postgres guarda los datos
clean-db:
	@echo ">> Borrando contenedores y volúmenes de la BBDD..."
	@docker rm -f -v $(CONTAINER_NAME) >/dev/null 2>&1 || true

start-db:
	@echo ">> Levantando la BBDD..."
	docker run --name $(CONTAINER_NAME) -p 5432:5432 -d $(IMAGE_NAME)

# Espera a que Postgres acepte conexiones por TCP (máximo 30 s). Con -h localhost
# no da "listo" mientras se inicializa y todavía está cargando schema.sql.
wait-db:
	@echo ">> Esperando a que la BBDD esté lista..."
	@i=0; until docker exec $(CONTAINER_NAME) pg_isready -h localhost -U $(DB_USER) -d $(DB_NAME) >/dev/null 2>&1; do \
		i=$$((i+1)); if [ $$i -ge 30 ]; then echo "!! La BBDD no arrancó"; exit 1; fi; \
		sleep 1; \
	done
	@echo ">> BBDD lista."

clean-uploads:
	@echo ">> Borrando imágenes subidas ($(UPLOADS_DIR)/)..."
	@rm -rf $(UPLOADS_DIR)
