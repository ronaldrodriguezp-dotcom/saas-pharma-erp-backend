package http

import (
	"encoding/json"
	"net/http"

	"fortiasaass-kio-speed/internal/purchasing/application/commands"
)

type GoodsReceiptHTTPHandler struct {
	ConfirmHandler *commands.ConfirmGoodsReceiptHandler
}

func NewGoodsReceiptHTTPHandler(confirmHandler *commands.ConfirmGoodsReceiptHandler) *GoodsReceiptHTTPHandler {
	return &GoodsReceiptHTTPHandler{
		ConfirmHandler: confirmHandler,
	}
}

type ConfirmRequest struct {
	ReceiptID string `json:"receipt_id"`
	Actor     string `json:"actor"`
}

func (h *GoodsReceiptHTTPHandler) HandleConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el tenant estricto desde los headers HTTP de la petición multitenant
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, "Header X-Tenant-ID es obligatorio", http.StatusBadRequest)
		return
	}

	var req ConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Payload JSON inválido", http.StatusBadRequest)
		return
	}

	cmd := commands.ConfirmGoodsReceiptCommand{
		ReceiptID: req.ReceiptID,
		Actor:     req.Actor,
	}

	err := h.ConfirmHandler.Handle(r.Context(), tenantID, cmd)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "recepción confirmada e inventario actualizado con éxito"})
}
