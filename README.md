# fullstack-challenge

# Fullstack Developer Test Challenge

Proyek ini adalah implementasi microservices sederhana dengan arsitektur **clean architecture** + **event-driven communication**.  
Stack terdiri dari:

- **Product-service** → NestJS (Node.js, TypeScript)  
- **Order-service** → Go (Golang)  
- **Database** → PostgreSQL  
- **Cache** → Redis  
- **Message broker** → RabbitMQ (opsional, untuk komunikasi event-driven)  

---

## 🚀 Cara Menjalankan Secara Lokal

### Prasyarat
- Node.js v18+  
- Go v1.20+  
- Docker Desktop (untuk jalanin Postgres & Redis dengan ringan)

### Step 1: Jalankan Postgres & Redis
Gunakan Docker ringan untuk DB & cache:
```bash
# Jalankan Postgres
docker run --name local-postgres -e POSTGRES_USER=dev -e POSTGRES_PASSWORD=dev -e POSTGRES_DB=appdb -p 5432:5432 -d postgres:15

# Jalankan Redis
docker run --name local-redis -p 6379:6379 -d redis:7


### Step 2: Jalankan Product-service (NestJS)
cd product-service
npm install
npm run start:dev


Service berjalan di http://localhost:3000.

### Step 3: Jalankan Order-service (Go)
cd order-service/cmd/server
go run main.go handlers.go product_client.go repository.go rabbitmq.go


Service berjalan di http://localhost:4000.

#Contoh Request API
Product Service

Create product

curl -X POST http://localhost:3000/products \
-H "Content-Type: application/json" \
-d '{"name":"Laptop","price":1000,"qty":5}'


Get product

curl http://localhost:3000/products/<product-id>

Order Service

Create order

curl -X POST http://localhost:4000/orders \
-H "Content-Type: application/json" \
-d '{"productId":"<product-id>","totalPrice":1000}'


Get orders by product

curl http://localhost:4000/orders/product/<product-id>


