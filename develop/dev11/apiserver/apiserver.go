package apiserver

import (
	"dev11/models"
	"dev11/store"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	errBadRequestByMethod     = errors.New("некорректный запрос")
	errQueryParamNotProvided  = errors.New("дата, начиная с которой нужно вывести события в календаре, не указана")
	errInvalidQueryDate       = errors.New("дата должна быть представлена в формате YYYY-MM-DD")
	errNotProvidedIDInForm    = errors.New("тело запроса должно содержать поле id:int")
	errNotPovidedUserIDInForm = errors.New("тело запроса должно содержать поле user_id:int")
	errNotProvidedDateInForm  = errors.New("тело запроса должно содержать дату в формате YYYY-MM-DD")
	errNotProvidedInfoInForm  = errors.New("тело запроса должно содержать поле info:string")
)

// APIServer ...
type APIServer struct {
	config *Config // Указатель на значение структурного типа Config. Значения полей могут быть установлены по умолчанию (при создании конструктором), могут быть получены из json-файла
	logger *log.Logger // logging object, каждый вызов метода Write io.Writer'а делает единственную операцию записи в журнал
	router *http.ServeMux // HTTP-мультиплексор, используется для выбора обработчика запроса.
	store  *store.Store // Структура Store хранит все ивенты
}

// New - конструктор объекта APIServer. Возвращает указатель на созданный экземпляр
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: log.Default(),
		router: http.NewServeMux(),
	}
}

// Start ...
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	// После того, как логгер будет сконфигурирован,
	s.logger.Println("логгер успешно сконфигурирован")

	if err := s.configureStore(); err != nil {
		return err
	}
	defer s.store.Close() // Откладываем вызов Close()
	s.logger.Println("хранилище успешно сконфигурировано")

	r := s.configureRouter()
	s.logger.Println("HTTP-мультиплексор успешно сконфигурирован")

	s.logger.Println("запуск api-сервера на порту:", s.config.BindAddr)
	return http.ListenAndServe(s.config.BindAddr, r)
}

func (s *APIServer) configureLogger() error {
	// os.OpenFile открывает файл, а если файла нет, то создает его. Она принимает три параметра:
	// путь к файлу, режим открытия файла (для чтения, для записи и т.д.), разрешения для доступа к файлу
	file, err := os.OpenFile(s.config.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	// Логи будут дублироваться в stdout
	mw := io.MultiWriter(os.Stdout, file)
	s.logger.SetOutput(mw)
	return nil
}

func (s *APIServer) configureStore() error {
	st := store.New() // Создаём значение Store
	if err := st.Open(); err != nil {
		return err
	} // Open() создаёт мапу, присваивает это значение полю db
	s.store = st
	return nil
}

func (s *APIServer) configureRouter() http.Handler {
	// При вызове http.HandleFunc с функциями-обработчиками мы не вызываем функцию-обработчик с передачей результата HandleFunc. Мы передаем HandleFunc
	// саму функцию. Эта функция сохраняется для того, чтобы быть вызванной в будущем,
	// при получении запроса с подходящим путем. В нашем случае мы вызываем метод, который возвращает функцию-обработчик. Её тип http.HandlerFunc, это
	// значит, что при вызове она принимает значения http.ResponseWriter и http.Request в качестве аргументов
	s.router.HandleFunc("/create_event", s.handleCreate())
	s.router.HandleFunc("/update_event", s.handleUpdate())
	s.router.HandleFunc("/delete_event", s.handleDelete())
	s.router.HandleFunc("/events_for_day", s.handleGetForDay())
	s.router.HandleFunc("/events_for_week", s.handleGetForWeek())
	s.router.HandleFunc("/events_for_month", s.handleGetForMonth())
	// Ф-ция вернёт значение ServeMux, обёрнутое logMiddleware
	return s.logMiddleware(s.router)
}

func (s *APIServer) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Обращаемся к полю Method объекта Request, сравниваем его значение со значением константы MethodPost
		// По сути эта операция сводится к сравнению двух строк. В первую очередь нужно проверить метод, с которым
		// был отправлен запрос, поскольку API-метод /create_event несовместим с любыми методами, кроме POST
		if r.Method == http.MethodPost {
			// В заголовке запроса значению "content-type" должно соответствовать "application/x-www-form-urlencoded",
			// тогда преходим к следующему блоку "if", в противном случае возвращаем ошибку
			// ----------------------------------------------------------------------------
			// И в случае ошибки (возникшей по какой-либо причине), и в случае успешного выполнения запроса,
			// клиентской стороне будет возвращён некий код состояния. Также (опционально) вместе с кодом
			// может быть возвращено значение ошибки (по причине неудачного выполнения какой-то функции)
			if r.Header.Get("content-type") != "application/x-www-form-urlencoded" {
				s.error(w, r, http.StatusUnsupportedMediaType, nil)
				return
			}
			// Для POST-запросов ParseForm считывает тело запроса, парсит его как форму
			// и помещает результаты как в r.PostForm, так и в r.Form в виде мапы
			if err := r.ParseForm(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			// Далее необходимо создать объект models.EventRequest и считать в него содержимое r.Form. Делаем это с помощью метода decodeFormCreate
			// Объект eventR (models.EventRequest) играет роль промежуточного хранилища ещё не проверенных на корректность данных.
			// Поле r.Form типа url.Values является мапой со строками-ключами и слайсами строк - значениями.
			eventR, err := s.decodeFormCreate(r.Form)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			// Производим валидацию значений полей eventR
			if err := eventR.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			// После того как значения eventR прошли проверку на корректность, создаем окончательный ивент
			event := models.NewEventFromRequest(eventR)
			if err := s.store.EventRepository().CreateEvent(event); err != nil {
				// В случае, если не удалось создать новую запись, сервер будет возвращать
				// код состояния 503 и значение ошибки в виде json-объекта
				s.error(w, r, 503, err)
				return
			}

			s.respond(w, r, http.StatusCreated, nil)
			return
		}
		s.error(w, r, http.StatusBadRequest, errBadRequestByMethod)
	}
}

func (s *APIServer) handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.Header.Get("content-type") != "application/x-www-form-urlencoded" {
				s.error(w, r, http.StatusUnsupportedMediaType, nil)
				return
			}
			if err := r.ParseForm(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			eventR, err := s.decodeFormUpdate(r.Form)
			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := eventR.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			event := models.NewEventFromRequest(eventR)
			if err := s.store.EventRepository().UpdateEvent(event); err != nil { // UpdateEvent обновляет ивент в базе (мапе) по id-шнику (ключу)
				s.error(w, r, http.StatusServiceUnavailable, err)
				return
			}

			s.respond(w, r, http.StatusAccepted, nil)
			return
		}
		s.error(w, r, http.StatusBadRequest, errBadRequestByMethod)
	}
}

func (s *APIServer) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if r.Header.Get("content-type") != "application/x-www-form-urlencoded" {
				s.error(w, r, http.StatusUnsupportedMediaType, nil)
				return
			}
			if err := r.ParseForm(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			val, ok := r.Form["id"]
			if !ok {
				s.error(w, r, http.StatusBadRequest, errNotProvidedIDInForm)
				return
			}
			id, err := strconv.Atoi(val[0])
			if err != nil {
				s.error(w, r, http.StatusBadRequest, errNotProvidedIDInForm)
				return
			}
			if err := s.store.EventRepository().DeleteEvent(id); err != nil { // DeleteEvent удаляет ивент из базы (мапы) по id-шнику (ключу)
				s.error(w, r, http.StatusServiceUnavailable, err)
				return
			}

			s.respond(w, r, http.StatusAccepted, nil)
			return
		}
		s.error(w, r, http.StatusBadRequest, errBadRequestByMethod)
	}
}

func (s *APIServer) handleGetForDay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Т. к. в GET методах параметры передаются через queryString, а не ч/з тело запроса,
			// используем метод Query() для получения параметров
			params := r.URL.Query()
			date, ok := params["date"]
			log.Println(date)
			// Если значение по ключу "date" отсутствует, вернём ошибку клиентской стороне
			if !ok {
				s.error(w, r, http.StatusBadRequest, errQueryParamNotProvided)
				return
			}
			// Парсим дату в соответствии с шаблоном
			startDateStr := date[0]
			startDate, err := time.Parse(store.BaseTimeSample, startDateStr)
			// Если дату не удалось распарсить, возвращаем ошибку "неверный формат даты"
			if err != nil {
				s.error(w, r, http.StatusBadRequest, errInvalidQueryDate)
				return
			}
			// Формируем дату (+1 день), отталкиваясь от изначальной
			oneDayLaterStr := startDate.AddDate(0, 0, 1).Format(store.BaseTimeSample)
			// Получаем ивенты из указанного временного диапазона
			events, err := s.store.EventRepository().GetEventsForDates(startDateStr, oneDayLaterStr)
			if err != nil {
				s.error(w, r, http.StatusServiceUnavailable, err)
				return
			}
			s.respond(w, r, http.StatusOK, map[string]interface{}{"events": events})
			return
		}
		s.error(w, r, http.StatusBadRequest, errBadRequestByMethod)
	}
}

func (s *APIServer) handleGetForWeek() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			params := r.URL.Query()
			date, ok := params["date"]
			log.Println(date)
			if !ok {
				s.error(w, r, http.StatusBadRequest, errQueryParamNotProvided)
				return
			}
			startDate, err := time.Parse(store.BaseTimeSample, date[0])
			if err != nil {
				s.error(w, r, http.StatusBadRequest, errInvalidQueryDate)
				return
			}
			oneWeekLaterStr := startDate.AddDate(0, 0, 7).Format(store.BaseTimeSample)

			events, err := s.store.EventRepository().GetEventsForDates(date[0], oneWeekLaterStr)
			if err != nil {
				s.error(w, r, http.StatusServiceUnavailable, err)
				return
			}
			s.respond(w, r, http.StatusOK, map[string]interface{}{"events": events})
			return
		}
		s.error(w, r, http.StatusBadRequest, errBadRequestByMethod)
	}
}

func (s *APIServer) handleGetForMonth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			params := r.URL.Query()
			date, ok := params["date"]
			if !ok {
				s.error(w, r, http.StatusBadRequest, errQueryParamNotProvided)
				return
			}
			startDate, err := time.Parse(store.BaseTimeSample, date[0])
			if err != nil {
				s.error(w, r, http.StatusBadRequest, errInvalidQueryDate)
				return
			}
			oneMonthLaterStr := startDate.AddDate(0, 1, 0).Format(store.BaseTimeSample)

			events, err := s.store.EventRepository().GetEventsForDates(date[0], oneMonthLaterStr)
			if err != nil {
				s.error(w, r, http.StatusServiceUnavailable, err)
				return
			}
			s.respond(w, r, http.StatusOK, map[string]interface{}{"events": events})
			return
		}
		s.error(w, r, http.StatusBadRequest, errBadRequestByMethod)
	}
}
// Метод error что-то вроде частного случая метода respond (или обёртка над ним)
// Метод error вызывает метод respond у того же получателя, но с определенными значениями параметров, в частности, код состояния будет соответствовать какой-либо ошибке.
// Также формат возвращённого клиенту json будет отличаться: в json-объекте будет одна пара "ключ-значение", ключ - строка "error", значение - строковое описание ошибки
func (s *APIServer) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}
// respond возвращает ответ клиенту в формате json
func (s *APIServer) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code) // Отправляет заголовок ответа HTTP с предоставленным кодом состояния
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (s *APIServer) decodeFormUpdate(form url.Values) (*models.EventRequest, error) {
	eventR := new(models.EventRequest)
	// API-метод /update_event помимо прочего должен содержать параметр, представляющий id уже имеющегося в базе ивента,
	// который необходимо обновить
	sliceID, ok := form["id"]
	if !ok {
		return nil, errNotProvidedIDInForm
	}
	id, err := strconv.Atoi(sliceID[0])
	if err != nil {
		return nil, errNotProvidedIDInForm
	}
	eventR.ID = id

	sliceUserID, ok := form["user_id"]
	if !ok {
		return nil, errNotPovidedUserIDInForm
	}
	userID, err := strconv.Atoi(sliceUserID[0])
	if err != nil {
		return nil, errNotPovidedUserIDInForm
	}
	eventR.UserID = userID

	sliceDate, ok := form["date"]
	if !ok {
		return nil, errNotProvidedDateInForm
	}
	eventR.Date = sliceDate[0]

	sliceInfo, ok := form["info"]
	if !ok {
		return nil, errNotProvidedInfoInForm
	}
	eventR.Info = sliceInfo[0]
	return eventR, nil
}
// decodeFormCreate осуществляет парсинг параметров метода /create_event
// decodeFormCreate принимает на вход форму из запроса (мапу), содержащую значения параметров метода /create_event,
// затем заполняет новое значение структурного типа EventRequest значениями этих параметров
func (s *APIServer) decodeFormCreate(form url.Values) (*models.EventRequest, error) {
	eventR := new(models.EventRequest)
	// Считываем значение из мапы по ключу
	// Доступ к элементу мапы с помощью индексации всегда даёт значение.
	// Если ключ присутствует в мапе, получим соответствующее значение, если нет - получим нулевое значение типа элемента.
	// Второе значение - это логическое значение, показывающее, имеется ли данный элемент в мапе
	sliceUserID, ok := form["user_id"]
	if !ok {
		return nil, errNotPovidedUserIDInForm
	}
	// Т. к. тип значений мапы - слайс строк, получаем первый элемент из него и преобразуем к int
	userID, err := strconv.Atoi(sliceUserID[0])
	if err != nil {
		return nil, errNotPovidedUserIDInForm
	}
	eventR.UserID = userID

	sliceDate, ok := form["date"]
	if !ok {
		return nil, errNotProvidedDateInForm
	}
	eventR.Date = sliceDate[0]

	sliceInfo, ok := form["info"]
	if !ok {
		return nil, errNotProvidedInfoInForm
	}
	eventR.Info = sliceInfo[0]
	return eventR, nil
}
