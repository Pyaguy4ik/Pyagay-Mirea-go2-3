package middleware

import (
    "net/http"
    "time"
    "go.uber.org/zap"
)

// responseWriter обертка для захвата статуса ответа
type responseWriter struct {
    http.ResponseWriter
    status      int
    wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
    if rw.wroteHeader {
        return
    }
    rw.status = code
    rw.ResponseWriter.WriteHeader(code)
    rw.wroteHeader = true
}

// AccessLog логирует входящие HTTP запросы
func AccessLog(logger *zap.Logger) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Оборачиваем ResponseWriter
            wrapped := &responseWriter{
                ResponseWriter: w,
                status:         http.StatusOK,
            }
            
            // Получаем request-id из контекста
            requestID := GetRequestID(r.Context())
            
            // Создаем логгер с базовыми полями
            reqLogger := logger.With(
                zap.String("request_id", requestID),
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("remote_ip", r.RemoteAddr),
                zap.String("user_agent", r.UserAgent()),
            )
            
            // Логируем начало запроса (опционально, для отладки)
            reqLogger.Debug("request started")
            
            // Вызываем следующий обработчик
            next.ServeHTTP(wrapped, r)
            
            // Вычисляем длительность
            duration := time.Since(start)
            
            // Логируем завершение запроса
            reqLogger.Info("request completed",
                zap.Int("status", wrapped.status),
                zap.Float64("duration_ms", float64(duration.Milliseconds())),
            )
        })
    }
}
