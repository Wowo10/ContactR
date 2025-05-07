package database

import (
	"Contacter/internal/models"
	"fmt"
	"time"
)

func (s *service) CheckUser(email string) (isValid bool, isAdmin bool, err error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)

	if err != nil {
		return
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		var name string
		var email string
		var validTo time.Time
		var admin bool
		err = rows.Scan(&id, &name, &email, &validTo, &admin)
		if err != nil {
			return
		}

		return validTo.After(time.Now()), admin, nil
	}

	return
}

func (s *service) GetUsers(page int) (users []models.User, count int, err error) {
	err = s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return
	}

	rows, err := s.db.Query(`
		SELECT * 
		FROM users 
		ORDER BY Id 
		LIMIT $1 OFFSET $2
	`, PAGE_SIZE, page*PAGE_SIZE)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var email string
		var validTo time.Time
		var admin bool
		err = rows.Scan(&id, &name, &email, &validTo, &admin)

		if err != nil {
			return
		}

		users = append(users, models.User{
			Id:         id,
			Name:       name,
			Email:      email,
			ValidUntil: validTo,
			IsValid:    validTo.After(time.Now()),
			IsAdmin:    admin,
		})
	}

	return
}

func (s *service) CreateUser(user models.User) (models.User, error) {
	query := `
        INSERT INTO users (name, email, validTo, admin)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, email, validTo, admin
    `

	var createdUser models.User
	err := s.db.QueryRow(query, user.Name, user.Email, user.ValidUntil, user.IsAdmin).
		Scan(&createdUser.Id, &createdUser.Name, &createdUser.Email, &createdUser.ValidUntil, &createdUser.IsAdmin)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return createdUser, nil
}

func (s *service) EditUser(user models.User) error {
	query := `
        UPDATE users
        SET
            name = $2,
            email = $3,
            validTo = $4,
            admin = $5
        WHERE id = $1
    `

	res, err := s.db.Exec(query, user.Id, user.Name, user.Email, user.ValidUntil, user.IsAdmin)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user updated with id %d", user.Id)
	}

	return nil
}

func (s *service) DeleteUser(id string) (err error) {
	res, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user deleted with id %s", id)
	}

	return
}
