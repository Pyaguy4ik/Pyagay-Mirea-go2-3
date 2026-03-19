package middleware

import (
    "context"
    "net/http"
    "github.com/google/uuid"
)

type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    RequestIDHeader = "X-Request-ID"
)

// RequestID middleware добавляет или генерирует request-id
func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get(RequestIDHeader)
        
        // Если request-id не передан, генерируем новый
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Добавляем в ответ
        w.Header().Set(RequestIDHeader, requestID)
        
        // Добавляем в контекст
        ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
        
        // Передаем дальше
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// GetRequestID извлекает request-id из контекста
func GetRequestID(ctx context.Context) string {
    if val := ctx.Value(RequestIDKey); val != nil {
        if id, ok := val.(string); ok {
            return id
        }
    }
    return ""
}
