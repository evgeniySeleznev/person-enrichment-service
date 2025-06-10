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

func (r *PersonRepository) GetAll(ctx context.Context, page, pageSize int) ([]model.Person, error) {
	query := `SELECT person_id, name, surname, patronymic, age, gender, nationality 
              FROM people LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, pageSize, (page-1)*pageSize)
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
