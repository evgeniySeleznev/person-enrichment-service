package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"unicode"

	"github.com/evgeniySeleznev/person-enrichment-service/internal/model"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/service"
	"github.com/evgeniySeleznev/person-enrichment-service/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// @title Person Enrichment API
// @version 1.0
// @description Сервис для обогащения данных о людях

// @host localhost:8080
// @BasePath /api

type PersonHandler struct {
	service *service.PersonService
	logger  logger.Logger
}

func NewPersonHandler(service *service.PersonService, logger logger.Logger) *PersonHandler {
	logger.Debug("ENTER: NewPersonHandler")
	return &PersonHandler{
		service: service,
		logger:  logger,
	}
}

var validate = validator.New()

func isAlphaUnicode(fl validator.FieldLevel) bool {
	for _, ch := range fl.Field().String() {
		if !unicode.IsLetter(ch) {
			return false
		}
	}
	return true
}

func init() {
	// Регистрируем кастомную валидацию
	validate.RegisterValidation("alpha_unicode", isAlphaUnicode)
}

// CreatePerson обрабатывает POST /api/persons
// @Summary Создать нового человека
// @Description Добавляет нового человека в систему с обогащёнными данными (возраст, пол, национальность)
// @Tags Люди
// @Accept json
// @Produce json
// @Param input body model.PersonInput true "Данные человека"
// @Success 201 {object} model.Person "Человек успешно создан"
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/persons [post]
func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("ENTER: CreatePerson")
	var input model.PersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("Invalid JSON", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(input); err != nil {
		h.logger.Error("Input validation error: ", err)
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	person, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.logger.Error("Failed to create person", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
	h.logger.Debug("EXIT: CreatePerson")
}

// GetPerson обрабатывает GET /api/persons/{id}
// @Summary Получить информацию о человеке по ID
// @Description Возвращает данные человека по его уникальному ID
// @Tags Люди
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} model.Person "Информация о человеке"
// @Failure 404 {string} string "Человек не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/persons/{id} [get]
func (h *PersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("ENTER: GetPerson")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	person, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get person", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
	h.logger.Debug("EXIT: GetPerson")
}

// UpdatePerson обрабатывает PATCH /api/persons/{id}
// @Summary Обновить данные человека
// @Description Обновляет информацию о человеке по его ID
// @Tags Люди
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param input body model.Person true "Новые данные"
// @Success 204 "Данные успешно обновлены"
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 404 {string} string "Человек не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/persons/{id} [patch]
func (h *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("ENTER: UpdatePerson")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var person model.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		h.logger.Error("Invalid JSON", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), id, &person); err != nil {
		h.logger.Error("Failed to update person", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("EXIT: UpdatePerson")
}

// DeletePerson обрабатывает DELETE /api/persons/{id}
// @Summary Удалить человека по ID
// @Description Удаляет запись о человеке по уникальному ID
// @Tags Люди
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Success 204 "Человек успешно удалён"
// @Failure 400 {string} string "Неверный формат ID"
// @Failure 404 {string} string "Человек не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/persons/{id} [delete]
func (h *PersonHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("ENTER: DeletePerson")
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.logger.Error("Failed to delete person", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.logger.Debug("EXIT: DeletePerson")
}

// GetAllPersons обрабатывает GET /api/persons
// @Summary Получить список людей с фильтрацией и пагинацией
// @Description Возвращает список людей с пагинацией и фильтрацией по полю (имя, фамилия, возраст и т.д.)
// @Tags Люди
// @Accept json
// @Produce json
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(10)
// @Param name query string false "Имя" example("Иван")
// @Param surname query string false "Фамилия" example("Иванов")
// @Param age_min query int false "Минимальный возраст"
// @Param age_max query int false "Максимальный возраст"
// @Param gender query string false "Пол" enum(male,female)
// @Param nationality query string false "Национальность"
// @Success 200 {array} model.Person "Список людей"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/persons [get]
func (h *PersonHandler) GetAllPersons(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("ENTER: GetAllPersons")

	// Пагинация
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Фильтры
	filterParams := model.FilterParams{
		Name:        getStringFromQuery(r, "name"),
		Surname:     getStringFromQuery(r, "surname"),
		AgeMin:      getIntFromQuery(r, "age_min"),
		AgeMax:      getIntFromQuery(r, "age_max"),
		Gender:      getStringFromQuery(r, "gender"),
		Nationality: getStringFromQuery(r, "nationality"),
		Page:        page,
		PageSize:    pageSize,
	}

	// Получаем от сервиса с фильтрацией
	persons, err := h.service.GetAll(r.Context(), filterParams)
	if err != nil {
		h.logger.Error("Failed to get persons", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
	h.logger.Debug("EXIT: GetAllPersons")
}

// Утилита для получения строки из query-параметра
func getStringFromQuery(r *http.Request, key string) *string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	return &value
}

// Утилита для получения int из query-параметра
func getIntFromQuery(r *http.Request, key string) *int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &intValue
}

// HealthCheck обрабатывает GET /health
// @Summary Проверка доступности API
// @Description Возвращает статус сервера для проверки его доступности
// @Tags Здоровье
// @Produce json
// @Success 200 {string} string "OK"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /health [get]
func (h *PersonHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("ENTER: HealthCheck")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
