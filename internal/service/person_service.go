package service

import (
	"context"
	"fmt"

	"github.com/evgeniySeleznev/person-enrichment-service/internal/model"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/repository/api"
	"github.com/evgeniySeleznev/person-enrichment-service/internal/repository/postgresql"
)

type PersonService struct {
	personRepo *postgresql.PersonRepository
	apiClient  *api.APIClient
}

func NewPersonService(personRepo *postgresql.PersonRepository, apiClient *api.APIClient) *PersonService {
	return &PersonService{
		personRepo: personRepo,
		apiClient:  apiClient,
	}
}

// Create создаёт человека и обогащает данные через внешние API
func (s *PersonService) Create(ctx context.Context, input model.PersonInput) (*model.Person, error) {
	// 1. Подготовка базовой структуры
	person := &model.Person{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: input.Patronymic,
	}

	//транслитерация для кириллицы
	nameForAPI := Transliterate(input.Name)

	// 2. Обогащение данных (параллельные запросы к API)
	age, err := s.apiClient.GetAge(ctx, nameForAPI)
	if err != nil {
		return nil, fmt.Errorf("failed to get age: %w", err)
	}
	person.Age = &age

	gender, err := s.apiClient.GetGender(ctx, person.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get gender: %w", err)
	}
	person.Gender = &gender

	nationality, err := s.apiClient.GetNationality(ctx, person.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get nationality: %w", err)
	}
	person.Nationality = &nationality

	// 3. Сохранение в БД
	id, err := s.personRepo.Create(ctx, person)
	if err != nil {
		return nil, fmt.Errorf("failed to save person: %w", err)
	}

	person.ID = id
	return person, nil
}

// GetByID возвращает человека по ID
func (s *PersonService) GetByID(ctx context.Context, id int64) (*model.Person, error) {
	return s.personRepo.GetByID(ctx, id)
}

// GetAll возвращает список людей с пагинацией
func (s *PersonService) GetAll(ctx context.Context, page, pageSize int) ([]model.Person, error) {
	return s.personRepo.GetAll(ctx, page, pageSize)
}

// Update обновляет данные человека
func (s *PersonService) Update(ctx context.Context, id int64, person *model.Person) error {
	return s.personRepo.Update(ctx, id, person)
}

// Delete удаляет человека по ID
func (s *PersonService) Delete(ctx context.Context, id int64) error {
	return s.personRepo.Delete(ctx, id)
}
