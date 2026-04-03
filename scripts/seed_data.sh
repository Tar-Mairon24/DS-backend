#!/usr/bin/env bash

set -a
source "$(dirname "$0")/../.env"
set +a

BASE_URL=http://localhost:8080/api/v1
COOKIE_JAR="$(dirname "$0")/.cookies.txt"

post_data() {
    local name="$1"
    local body="$2"
    
    echo "Seeding $name..."
    response=$(curl -s -o /dev/null -w "%{http_code}" -X 'POST' \
        --location "$BASE_URL/$name" \
        -H "Content-Type: application/json" \
        -b "$COOKIE_JAR" \
        -c "$COOKIE_JAR" \
        -d "$body")
    
    echo "Response status code: $response"
}

# Seeds
post_data "users" '{
    "email": "'"$USER2_EMAIL"'",
    "username": "'"$USER2_NAME"'",
    "password": "'"$USER2_PASSWORD"'",
    "role": "agente",
    "phone": "8444428728"
}'

post_data "users/clients" '{
    "email": "cliente1@prueba.com",
    "username": "Cliente 1",
    "role": "client",
    "phone": "8440929384",
    "notes": "Es estudiante universitario"
}'

post_data "users/clients" '{
    "email": "cliente2@prueba.com",
    "username": "Cliente 2",
    "role": "client",
    "phone": "8440929434",
    "notes": "Tiene perros"
}'

post_data "users/clients" '{
    "email": "cliente3@prueba.com",
    "username": "Cliente 3",
    "role": "client",
    "phone": "8440929332",
    "notes": "Es muy desmadroso"
}'

post_data "users/owners" '{
    "email": "propietario1@prueba.com",
    "username": "Propietario 1",
    "role": "owner",
    "phone": "8440929284",
    "notes": "No acepta mascotas"
}'

post_data "users/owners" '{
    "email": "propietario2@prueba.com",
    "username": "Propietario 2",
    "role": "owner",
    "phone": "8440922384",
    "notes": "No acepta niños ni mascotas"
}'

post_data "users/owners" '{
    "email": "propietario3@prueba.com",
    "username": "Propietario 3",
    "role": "owner",
    "phone": "8440929784",
    "notes": "Es muy amable y acepta mascotas"
}'

post_data "properties" '{
  "title": "Casa Moderna en San Patricio",
  "listing_date": "2025-01-15T10:30:00Z",
  "address": "Av. Principal 123, San Patricio",
  "neighborhood": "San Patricio",
  "city": "Saltillo",
  "zone": "Noreste",
  "reference": "Frente al parque central, esquina con calle secundaria",
  "price": 2500000.00,
  "construction_m2": 180,
  "land_m2": 250,
  "is_occupied": false,
  "is_furnished": true,
  "floors": 2,
  "bedrooms": 3,
  "bathrooms": 2,
  "garage_size": 2,
  "garden_m2": 50,
  "gas_types": ["natural", "lp"],
  "amenities": ["piscina", "gimnasio", "área de juegos", "jardín"],
  "extras": ["aire acondicionado", "calefacción", "sistema de seguridad"],
  "utilities": ["agua", "luz", "gas", "internet", "cable"],
  "notes": "Propiedad en excelente estado, recientemente remodelada. Ubicada en zona tranquila y segura.",
  "owner_id": 1,
  "user_id": 1,
  "property_type": "Casa",
  "transaction_type": "Venta",
  "status": "Disponible"
}'

post_data "properties" '{
  "title": "Departamento Amueblado en Ramos Arizpe",
  "listing_date": "2025-01-18T14:15:00Z",
  "address": "Calle 5 de Mayo 456, Centro",
  "neighborhood": "Centro",
  "city": "Ramos Arizpe",
  "zone": "Centro",
  "reference": "Cercano a la plaza mayor, a dos cuadras de la estación de autobús",
  "price": 15000.00,
  "construction_m2": 120,
  "land_m2": 120,
  "is_occupied": true,
  "is_furnished": true,
  "floors": 1,
  "bedrooms": 2,
  "bathrooms": 1,
  "garage_size": 1,
  "garden_m2": 0,
  "gas_types": ["natural"],
  "amenities": ["estacionamiento", "servicios básicos incluidos"],
  "extras": ["aire acondicionado", "cocina integral", "lavadora"],
  "utilities": ["agua", "luz", "gas", "internet"],
  "notes": "Departamento ideal para personas o parejas. Incluye servicios básicos en la renta. Disponible para renta inmediata.",
  "owner_id": 2,
  "user_id": 1,
  "property_type": "Apartamento",
  "transaction_type": "Renta",
  "status": "Disponible"
}'

post_data "properties" '{
  "title": "Casa Colonial en Arteaga",
  "listing_date": "2025-01-20T09:45:00Z",
  "address": "Calle Independencia 789, Zona Histórica",
  "neighborhood": "Zona Histórica",
  "city": "Arteaga",
  "zone": "Poniente",
  "reference": "En la zona histórica, a una cuadra de la plaza principal",
  "price": 1800000.00,
  "construction_m2": 220,
  "land_m2": 300,
  "is_occupied": false,
  "is_furnished": false,
  "floors": 2,
  "bedrooms": 4,
  "bathrooms": 3,
  "garage_size": 2,
  "garden_m2": 80,
  "gas_types": ["lp"],
  "amenities": ["jardín amplio", "terraza", "sótano"],
  "extras": ["sistema de seguridad", "cerca perimetral", "puerta de herrería"],
  "utilities": ["agua", "luz", "gas"],
  "notes": "Casa colonial con mucho potencial. Requiere mantenimiento. Ubicada en zona de alto valor histórico.",
  "owner_id": 3,
  "user_id": 1,
  "property_type": "Casa",
  "transaction_type": "Venta",
  "status": "Reservado"
}'

post_data "properties" '{
  "title": "Loft Moderno en Saltillo",
  "listing_date": "2025-01-22T16:20:00Z",
  "address": "Paseo Urdinola 234, Zona Nueva",
  "neighborhood": "Zona Nueva",
  "city": "Saltillo",
  "zone": "Este",
  "reference": "En complejo de oficinas y vivienda, con acceso a todas las comodidades",
  "price": 2000000.00,
  "construction_m2": 150,
  "land_m2": 150,
  "is_occupied": false,
  "is_furnished": false,
  "floors": 1,
  "bedrooms": 2,
  "bathrooms": 2,
  "garage_size": 1,
  "garden_m2": 0,
  "gas_types": ["natural"],
  "amenities": ["gimnasio comunitario", "piscina comunitaria", "seguridad 24/7"],
  "extras": ["aire acondicionado central", "pisos de madera", "ventanales amplios"],
  "utilities": ["agua", "luz", "gas", "internet", "cable"],
  "notes": "Loft de lujo en ubicación privilegiada. Ideal para profesionales. Condominio de alta seguridad.",
  "owner_id": 1,
  "user_id": 1,
  "property_type": "Loft",
  "transaction_type": "Venta",
  "status": "Disponible"
}'

post_data "properties" '{
  "title": "Casa en Renta en Ramos Arizpe",
  "listing_date": "2025-01-25T11:30:00Z",
  "address": "Avenida Las Fuentes 567, Fraccionamiento Arboledas",
  "neighborhood": "Fraccionamiento Arboledas",
  "city": "Ramos Arizpe",
  "zone": "Sur",
  "reference": "En fraccionamiento residencial cerrado, cerca de escuelas y comercios",
  "price": 18000.00,
  "construction_m2": 200,
  "land_m2": 280,
  "is_occupied": true,
  "is_furnished": true,
  "floors": 2,
  "bedrooms": 3,
  "bathrooms": 2,
  "garage_size": 2,
  "garden_m2": 60,
  "gas_types": ["natural", "lp"],
  "amenities": ["alberca en fraccionamiento", "área de juegos", "control de acceso"],
  "extras": ["aire acondicionado", "calefacción", "cocina equipada"],
  "utilities": ["agua", "luz", "gas", "internet"],
  "notes": "Casa familiar ideal para vivir. Fraccionamiento seguro y con buena plusvalía. Disponible para renta con depósito.",
  "owner_id": 2,
  "user_id": 1,
  "property_type": "Casa",
  "transaction_type": "Renta",
  "status": "Disponible"
}'

post_data "appointments" '{
    "title": "Visita inicial - Casa Pedregal",
    "description": "Recorrido inicial de la propiedad con comprador interesado",
    "start_date": "2026-04-03T10:00:00Z",
    "end_date": "2026-04-03T11:00:00Z",
    "status": "scheduled",
    "notes": "El cliente prefiere acabados modernos, destacar la remodelación de la cocina",
    "client_id": 3,
    "property_id": 1
}'

post_data "appointments" '{
    "title": "Segunda visita - Torre Cumbres",
    "description": "El cliente regresa para una segunda revisión antes de hacer una oferta",
    "start_date": "2026-04-03T14:00:00Z",
    "end_date": "2026-04-03T15:00:00Z",
    "status": "scheduled",
    "notes": null,
    "client_id": 4,
    "property_id": 2
}'

post_data "appointments" '{
    "title": "Visita de inversión - Departamento Cumbres",
    "description": "Cliente interesado en adquirir la propiedad como inversión para renta",
    "start_date": "2026-04-08T11:00:00Z",
    "end_date": "2026-04-08T12:00:00Z",
    "status": "scheduled",
    "notes": "Preguntar sobre rendimiento esperado y plusvalía de la zona",
    "client_id": 5,
    "property_id": 4
}'

post_data "appointments" '{
    "title": "Revisión de contrato - Residencial del Valle",
    "description": "Reunión para revisar términos del contrato de compraventa antes de firmar",
    "start_date": "2026-04-10T16:00:00Z",
    "end_date": "2026-04-10T17:30:00Z",
    "status": "scheduled",
    "notes": "El cliente solicita revisar cláusulas de penalización por cancelación",
    "client_id": 3,
    "property_id": 5
}'