package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"fortiasaass-kio-speed/internal/inventory"
	purchasingHttp "fortiasaass-kio-speed/internal/purchasing/adapters/http"
	purchasingCmd "fortiasaass-kio-speed/internal/purchasing/application/commands"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/erp_kio_speed?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error crítico al abrir la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo establecer conexión con PostgreSQL: %v", err)
	}

	log.Println("Conexión exitosa a PostgreSQL establecida correctamente.")

	// 1. Inicializar dependencias del módulo de Compras / Recepciones
	confirmHandler := purchasingCmd.NewConfirmGoodsReceiptHandler(db)
	goodsReceiptHandler := purchasingHttp.NewGoodsReceiptHTTPHandler(confirmHandler)

	// 2. Inicializar dependencias del módulo de Inventario y Lotes
	inventoryRepo := inventory.NewRepository(db)
	inventoryHandler := inventory.NewHandler(inventoryRepo)

	// 3. Configurar enrutador HTTP
	mux := http.NewServeMux()

	// Ejemplo de cómo debería verse la vinculación de la ruta:
r.HandleFunc("/api/v1/purchasing/receipts", goodsReceiptHandler.ConfirmReceipt).Methods("POST") goodsReceiptHandler.ConfirmGoodsReceiptHandler)

	// Rutas de Inventario (Consulta de Stock por Lotes)
	mux.HandleFunc("/api/v1/inventory/stocks", inventoryHandler.GetStocksHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Iniciando servidor ERP-POS SaaS en el puerto %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("El servidor se detuvo inesperadamente: %v", err)
	}
}
