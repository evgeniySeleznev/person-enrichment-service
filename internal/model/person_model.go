package model

// Person представляет данные человека в системе
// swagger:model
type Person struct {
	// Уникальный идентификатор
	// example: 1
	ID int64 `json:"id"`

	// Имя
	// example: Иван
	Name string `json:"name"`

	// Фамилия
	// example: Иванов
	Surname string `json:"surname"`

	// Отчество
	// example: Иванович
	Patronymic *string `json:"patronymic"`

	// Возраст
	// example: 30
	Age *int `json:"age"`

	// Пол (male/female)
	// example: male
	Gender *string `json:"gender"`

	// Код страны (2 символа)
	// example: RU
	Nationality *string `json:"nationality"`
}

// PersonInput представляет данные для создания человека
// swagger:model
type PersonInput struct {

	// Имя
	// example: Иван
	Name string `json:"name" validate:"required,min=2,max=50"`

	// Фамилия
	//example: Иванов
	Surname string `json:"surname" validate:"required,min=2,max=50"`

	// Отчество
	// example: Иванович
	Patronymic *string `json:"patronymic" validate:"omitempty,min=2,max=50"`
}

// FilterParams содержит параметры фильтрации
// swagger:model
type FilterParams struct {
	Name        *string `json:"name"`
	Surname     *string `json:"surname"`
	AgeMin      *int    `json:"age_min"`
	AgeMax      *int    `json:"age_max"`
	Gender      *string `json:"gender"`
	Nationality *string `json:"nationality"`

	// Номер страницы (начиная с 1)
	// minimum: 1
	// example: 1
	Page int `json:"page" validate:"min=1"`

	// Количество элементов на странице
	// minimum: 1
	// maximum: 100
	// example: 20
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}
