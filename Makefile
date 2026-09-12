CONTAINER_NAME=postgres-tp2
IMAGE_NAME=tp-postgres
# Tarea principal, se ejecuta en este orden al ejecutar "make"
test: pre-test run-tests post-tests

#Pre-tests:
pre-test: run-sqlc build-app build-db-image clean-db start-db wait-db
run-sqlc:
	@echo "Generando código con sqlc..."
	sqlc generate

build-app:
	@echo "Compilando el proyecto..."
	go build ./...

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

# Ejecución de los tests:
run-tests:
	@echo "Ejecutando los tests..."
	go test -v ./...

post-tests: 
	@echo "Limpiando basura (borrando contenedores y volúmenes)..."
	docker rm -f -v $(CONTAINER_NAME)
