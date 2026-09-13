package platform

import (
	"net/http"
)

func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"online","service":"fortiasaass-kio-speed","module":"KIO-01"}`))
	})

	protectedHandler := AuthMiddlewareABAC(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Context().Value(TenantIDContextKey).(string)

		// Registramos la acción crítica en la auditoría inmutable
		LogCriticalAction(tenantID, "sucursal_central", "12.345.678-9", "ACCESS_SECURE_DATA", "Consulta autorizada a recursos protegidos")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Acceso autorizado con éxito", "tenant":"` + tenantID + `"}`))
	}))

	mux.Handle("/api/v1/secure-data", protectedHandler)

	return mux
}
