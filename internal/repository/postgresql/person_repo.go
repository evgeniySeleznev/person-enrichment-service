package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/evgeniySeleznev/person-enrichment-service/internal/model"
)

type PersonRepository struct {
	db *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(ctx context.Context, person *model.Person) (int64, error) {
	query := `INSERT INTO people (name, surname, patronymic, age, gender, nationality) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING person_id`

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		person.Name, person.Surname, person.Patronymic,
		person.Age, person.Gender, person.Nationality).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create person: %w", err)
	}

	return id, nil
}

func (r *PersonRepository) GetByID(ctx context.Context, id int64) (*model.Person, error) {
	query := `SELECT person_id, name, surname, patronymic, age, gender, nationality 
              FROM people WHERE person_id = $1`

	var person model.Person
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("person not found")
		}
		return nil, fmt.Errorf("failed to get person: %w", err)
	}

	return &person, nil
}

func (r *PersonRepository) GetAll(ctx context.Context, filterParams model.FilterParams) ([]model.Person, error) {
	query := `SELECT person_id, name, surname, patronymic, age, gender, nationality FROM people WHERE 1=1`
	var args []interface{}
	argID := 1 // номер аргумента для $n

	// Фильтры
	if filterParams.Name != nil {
		query += fmt.Sprintf(" AND name ILIKE $%d", argID)
		args = append(args, "%"+*filterParams.Name+"%")
		argID++
	}
	if filterParams.Surname != nil {
		query += fmt.Sprintf(" AND surname ILIKE $%d", argID)
		args = append(args, "%"+*filterParams.Surname+"%")
		argID++
	}
	if filterParams.AgeMin != nil {
		query += fmt.Sprintf(" AND age >= $%d", argID)
		args = append(args, *filterParams.AgeMin)
		argID++
	}
	if filterParams.AgeMax != nil {
		query += fmt.Sprintf(" AND age <= $%d", argID)
		args = append(args, *filterParams.AgeMax)
		argID++
	}
	if filterParams.Gender != nil {
		query += fmt.Sprintf(" AND gender = $%d", argID)
		args = append(args, *filterParams.Gender)
		argID++
	}
	if filterParams.Nationality != nil {
		query += fmt.Sprintf(" AND nationality = $%d", argID)
		args = append(args, *filterParams.Nationality)
		argID++
	}

	// Пагинация
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, filterParams.PageSize, (filterParams.Page-1)*filterParams.PageSize)

	fmt.Println("Executing query:", query)
	fmt.Println("With args:", args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query people: %w", err)
	}
	defer rows.Close()

	var people []model.Person
	for rows.Next() {
		var person model.Person
		if err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
		); err != nil {
			return nil, fmt.Errorf("failed to scan person: %w", err)
		}
		people = append(people, person)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return people, nil
}

func (r *PersonRepository) Update(ctx context.Context, id int64, person *model.Person) error {
	query := `UPDATE people SET 
              name = $1, 
              surname = $2, 
              patronymic = $3, 
              age = $4, 
              gender = $5, 
              nationality = $6 
              WHERE person_id = $7`

	result, err := r.db.ExecContext(ctx, query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("person not found")
	}

	return nil
}

func (r *PersonRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM people WHERE person_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("person not found")
	}

	return nil
}
