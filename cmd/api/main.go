package main

import (
	"log"
	"net/http"
	"fortiasaass-kio-speed/internal/platform"
)

func main() {
	router := platform.SetupRouter()

	port := ":8080"
	log.Printf("🚀 [Go Engine] Kiosco KIO-01 corriendo a alta velocidad en el puerto %s", port)
	
	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error al iniciar el servidor Go: %v", err)
	}
}