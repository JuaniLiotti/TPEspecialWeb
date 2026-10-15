#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== INICIANDO PRUEBAS DE LA API REST ==="

echo -e "\n--- 1. Pruebas de PRODUCTOS ---"
echo "POST: Creando un producto..."
curl -s -X POST $BASE_URL/products/ \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Cartera de Cuero", "descripcion": "100% cuero natural", "categoria": "Accesorios", "color": "Marrón", "precio": 45000}' | jq || echo "Hecho"

echo -e "\nGET: Obteniendo todos los productos..."
curl -s -X GET $BASE_URL/products/ | jq || echo "Hecho"

echo -e "\nGET: Obteniendo el producto con ID 1..."
curl -s -X GET $BASE_URL/products/1 | jq || echo "Hecho"

echo -e "\nPUT: Actualizando el producto con ID 1..."
curl -s -X PUT $BASE_URL/products/1 \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Cartera Premium", "descripcion": "Cuero y detalles en oro", "categoria": "Accesorios", "color": "Negro", "precio": 65000}' | jq || echo "Hecho"

echo -e "\nDELETE: Eliminando el producto con ID 1..."
curl -s -w "\nCódigo HTTP (se espera 204): %{http_code}\n" -X DELETE $BASE_URL/products/1


echo -e "\n--- 2. Pruebas de USUARIOS ---"
echo "POST: Creando un usuario..."
curl -s -X POST $BASE_URL/users/ \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Juan Perez", "email": "juan@test.com", "password": "password123"}' | jq || echo "Hecho"

echo -e "\nGET: Obteniendo todos los usuarios..."
curl -s -X GET $BASE_URL/users/ | jq || echo "Hecho"

echo -e "\nGET: Obteniendo el usuario con ID 1..."
curl -s -X GET $BASE_URL/users/1 | jq || echo "Hecho"

echo -e "\nPUT: Actualizando el usuario con ID 1..."
curl -s -X PUT $BASE_URL/users/1 \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Juan Ignacio", "email": "juan@test.com", "password": "nuevapassword"}' | jq || echo "Hecho"

echo -e "\nDELETE: Eliminando el usuario con ID 1..."
curl -s -w "\nCódigo HTTP (se espera 204): %{http_code}\n" -X DELETE $BASE_URL/users/1

echo -e "\n=== PRUEBAS FINALIZADAS ==="