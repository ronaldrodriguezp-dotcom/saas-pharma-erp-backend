package platform

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	TenantIDContextKey  ContextKey = "tenant_id"
	BranchIDContextKey  ContextKey = "branch_id"
	CashboxIDContextKey ContextKey = "cashbox_id"
	UserRoleContextKey  ContextKey = "user_role"
)

var JwtPublicKey = []byte("CLAVE_SECRETA_SUPER_SEGURA_O_LLAVE_PUBLICA_RS256")

func AuthMiddlewareABAC(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Token de autenticación requerido o formato inválido"}`))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return JwtPublicKey, nil
		})

		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"Token criptográficamente inválido o expirado"}`))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"No se pudieron extraer los atributos de seguridad del token"}`))
			return
		}

		tenantID, tenantOk := claims["tenant_id"].(string)
		branchID, _ := claims["branch_id"].(string)
		userRole, _ := claims["role"].(string)

		if !tenantOk || tenantID == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"El token no contiene un Tenant-ID válido"}`))
			return
		}

		requestedBranch := r.Header.Get("X-Branch-ID")
		if requestedBranch != "" && requestedBranch != branchID && userRole != "ADMIN_MASTER" {
			LogCriticalAction(tenantID, requestedBranch, "SISTEMA", "SECURITY_VIOLATION", "Intento de acceso cruzado no autorizado a sucursal ajena")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error":"Acceso denegado: Violación de política de atributos (ABAC por sucursal)"}`))
			return
		}

		ctx := context.WithValue(r.Context(), TenantIDContextKey, tenantID)
		ctx = context.WithValue(ctx, BranchIDContextKey, branchID)
		ctx = context.WithValue(ctx, UserRoleContextKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
