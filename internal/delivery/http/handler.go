package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/evgeniySeleznev/person-enrichment-service/internal/model"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/service"
	"github.com/evgeniySeleznev/person-enrichment-service/pkg/logger"
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
	return &PersonHandler{
		service: service,
		logger:  logger,
	}
}

// CreatePerson обрабатывает POST /api/persons
// CreatePerson godoc
// @Summary Создать человека
// @Description Добавляет новую запись с обогащёнными данными (возраст, пол, национальность)
// @Tags Люди
// @Accept json
// @Produce json
// @Param input body model.PersonInput true "Данные человека"
// @Success 201 {object} model.Person
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/persons [post]
func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input model.PersonInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("Invalid JSON", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
}

// GetPerson обрабатывает GET /api/persons/{id}
// GetPerson godoc
// @Summary Получить человека по ID
// @Tags Люди
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} model.Person
// @Failure 404 {string} string "Человек не найден"
// @Router /api/persons/{id} [get]
func (h *PersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
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
}

// GetAllPersons обрабатывает GET /api/persons
// GetAllPersons godoc
// @Summary Получить список людей
// @Description Возвращает список людей с пагинацией
// @Tags Люди
// @Produce json
// @Param page query int false "Номер страницы (по умолчанию 1)" default(1)
// @Param page_size query int false "Размер страницы (по умолчанию 10, максимум 100)" default(10)
// @Success 200 {array} model.Person
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/persons [get]
func (h *PersonHandler) GetAllPersons(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	persons, err := h.service.GetAll(r.Context(), page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get persons", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(persons)
}

// UpdatePerson обрабатывает PATCH /api/persons/{id}
// UpdatePerson godoc
// @Summary Обновить данные человека
// @Description Обновляет информацию о существующем человеке
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
}

// DeletePerson обрабатывает DELETE /api/persons/{id}
// DeletePerson godoc
// @Summary Удалить человека
// @Description Удаляет запись о человеке по ID
// @Tags Люди
// @Param id path int true "ID человека"
// @Success 204 "Человек успешно удалён"
// @Failure 400 {string} string "Неверный формат ID"
// @Failure 404 {string} string "Человек не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /api/persons/{id} [delete]
func (h *PersonHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
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
}
