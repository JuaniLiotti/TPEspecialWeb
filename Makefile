CONTAINER_NAME=postgres-tp2
IMAGE_NAME=tp-postgres

# ==========================================
# TAREAS PARA EL TP2 (Tests de Base de Datos)
# ==========================================
test: pre-test run-tests post-tests

pre-test: run-sqlc build-db-image clean-db start-db wait-db
run-tests:
	@echo "Ejecutando los tests de la BBDD..."
	go test -v ./...
post-tests: clean-db


# ==========================================
# TAREAS PARA EL TP3 (Tests de la API)
# ==========================================
test-api: pre-test-api run-api-tests post-test-api

pre-test-api: pre-test build-app start-app wait-app

build-app:
	@echo "Compilando el proyecto..."
	go build -o api-server main.go

start-app:
	@echo "Levantando el servidor de la aplicación en segundo plano..."
	./api-server & echo $$! > app.pid

wait-app:
	@echo "Esperando a que la API esté lista..."
	sleep 3

run-api-tests:
	@echo "Ejecutando pruebas de los endpoints..."
	chmod +x requests.sh
	./requests.sh

post-test-api: clean-app clean-db

clean-app:
	@echo "Cerrando el servidor de la aplicación..."
	@if [ -f app.pid ]; then kill -9 $$(cat app.pid) && rm app.pid || true; fi


# ==========================================
# TAREAS COMPARTIDAS
# ==========================================
run-sqlc:
	@echo "Generando código con sqlc..."
	sqlc generate

build-db-image:
	@echo "Construyendo la imagen de la BBDD..."
	docker build -t $(IMAGE_NAME) .

clean-db:
	@echo "Borrando contenedores y volúmenes viejos..."
	docker rm -f -v $(CONTAINER_NAME) || true

start-db: 
	@echo "Levantando la BBDD..."
	docker run --name $(CONTAINER_NAME) -p 5432:5432 -d $(IMAGE_NAME) 

wait-db: 
	@echo "Esperando a que la BBDD esté lista..."
	sleep 5