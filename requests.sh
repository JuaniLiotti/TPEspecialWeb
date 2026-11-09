#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== INICIANDO PRUEBAS DE LA API REST ==="

echo -e "\n--- 1. Pruebas de PRODUCTOS ---"
echo "POST: Creando un producto..."
curl -s -X POST $BASE_URL/api/products/ \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Cartera de Cuero", "descripcion": "100% cuero natural", "categoria": "Accesorios", "color": "Marrón", "precio": 45000}' | jq || echo "Hecho"

echo -e "\nGET: Obteniendo todos los productos..."
curl -s -X GET $BASE_URL/api/products/ | jq || echo "Hecho"

echo -e "\nGET: Obteniendo el producto con ID 1..."
curl -s -X GET $BASE_URL/api/products/1 | jq || echo "Hecho"

echo -e "\nPUT: Actualizando el producto con ID 1..."
curl -s -X PUT $BASE_URL/api/products/1 \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Cartera Premium", "descripcion": "Cuero y detalles en oro", "categoria": "Accesorios", "color": "Negro", "precio": 65000}' | jq || echo "Hecho"

echo -e "\nDELETE: Eliminando el producto con ID 1..."
curl -s -w "\nCódigo HTTP (se espera 204): %{http_code}\n" -X DELETE $BASE_URL/api/products/1


echo -e "\n--- 2. Pruebas de USUARIOS ---"
echo "POST: Creando un usuario..."
curl -s -X POST $BASE_URL/api/users/ \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Juan Perez", "email": "juan@test.com", "password": "password123"}' | jq || echo "Hecho"

echo -e "\nGET: Obteniendo todos los usuarios..."
curl -s -X GET $BASE_URL/api/users/ | jq || echo "Hecho"

echo -e "\nGET: Obteniendo el usuario con ID 1..."
curl -s -X GET $BASE_URL/api/users/1 | jq || echo "Hecho"

echo -e "\nPUT: Actualizando el usuario con ID 1..."
curl -s -X PUT $BASE_URL/api/users/1 \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Juan Ignacio", "email": "juan@test.com", "password": "nuevapassword"}' | jq || echo "Hecho"

echo -e "\nDELETE: Eliminando el usuario con ID 1..."
curl -s -w "\nCódigo HTTP (se espera 204): %{http_code}\n" -X DELETE $BASE_URL/api/users/1

echo -e "\n--- 3. Pruebas de CUEROS ---"
echo "POST: Creando un cuero..."
curl -s -X POST $BASE_URL/api/cueros/ \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Vaqueta Natural", "tipo": "Vacuno", "acabado": "Curtido vegetal", "color": "Suela", "espesor_mm": 2.5, "descripcion": "Cuero entero ideal para marroquinería", "precio": 120000}' | jq || echo "Hecho"

echo -e "\nPOST: Intentando crear un cuero inválido (sin nombre y espesor 0)..."
curl -s -w "\nCódigo HTTP (se espera 400): %{http_code}\n" -X POST $BASE_URL/api/cueros/ \
    -H "Content-Type: application/json" \
    -d '{"nombre": "", "tipo": "Vacuno", "acabado": "Graso", "color": "Negro", "espesor_mm": 0, "descripcion": "Inválido", "precio": 1000}'

echo -e "\nGET: Obteniendo todos los cueros..."
curl -s -X GET $BASE_URL/api/cueros/ | jq || echo "Hecho"

echo -e "\nGET: Obteniendo el cuero con ID 1..."
curl -s -X GET $BASE_URL/api/cueros/1 | jq || echo "Hecho"

echo -e "\nPUT: Actualizando el cuero con ID 1..."
curl -s -X PUT $BASE_URL/api/cueros/1 \
    -H "Content-Type: application/json" \
    -d '{"nombre": "Vaqueta Engrasada", "tipo": "Vacuno", "acabado": "Graso", "color": "Marrón oscuro", "espesor_mm": 2.8, "descripcion": "Cuero engrasado para cinturones", "precio": 135000}' | jq || echo "Hecho"

echo -e "\nDELETE: Eliminando el cuero con ID 1..."
curl -s -w "\nCódigo HTTP (se espera 204): %{http_code}\n" -X DELETE $BASE_URL/api/cueros/1

echo -e "\nGET: Buscando el cuero con ID 1 ya eliminado..."
curl -s -w "\nCódigo HTTP (se espera 404): %{http_code}\n" -X GET $BASE_URL/api/cueros/1

echo -e "\n=== PRUEBAS FINALIZADAS ==="