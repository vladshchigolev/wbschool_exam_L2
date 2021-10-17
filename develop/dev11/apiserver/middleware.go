package apiserver

import (
	"fmt"
	"net/http"
	"time"
)

func (s *APIServer) logMiddleware(next http.Handler) http.Handler {
	// Функция, которую мы собираемся вернуть, является самой обычной функцией. У нее нет метода ServeHTTP().
	// Так что по сути она не является обработчиком для HTTP запросов ().
	// Нам требуется превратить её в обработчик с помощью использования адаптера http.HandlerFunc() следующим образом:
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Println(
			fmt.Sprintf("обработка %s запроса с url: %s начата", r.Method, r.RequestURI),
		)

		start := time.Now() // Будем считать за какое время выполнится запрос
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r) // Здесь будет выполнена какая-то функция обработчик, соответствующая подходящему пути
		// ServeHTTP должен записать заголовки и данные ответа в ResponseWriter, а затем должен произойти возврат

		s.logger.Println(
			fmt.Sprintf("обработка запроса завершена со статусом %d %s за %v", rw.code, http.StatusText(rw.code), time.Since(start)),
		)
	})
}
