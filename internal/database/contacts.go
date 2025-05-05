package database

import (
	"Contacter/internal/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func (s *service) GetContacts(
	search string,
	matchAll bool,
	page int,
) (
	contacts []models.Contact,
	count int,
	err error,
) {
	terms := strings.FieldsFunc(search, func(r rune) bool {
		return r == ' ' || r == ','
	})

	operator := "&&"
	if matchAll {
		operator = "@>"
	}

	baseQuery := `
		SELECT id, name, linkedinurl, credlyurl, tags, contact
		FROM contacts
	`

	limitQuery := fmt.Sprintf(" LIMIT %d OFFSET %d", PAGE_SIZE, page*PAGE_SIZE)

	var rows *sql.Rows
	if len(terms) == 0 {
		rows, err = s.db.Query(baseQuery + limitQuery)
	} else {
		filterQuery := fmt.Sprintf("%s WHERE tags %s $1", baseQuery, operator)
		rows, err = s.db.Query(filterQuery+limitQuery, pq.Array(terms))
	}
	defer rows.Close()
	if err != nil {
		return
	}

	for rows.Next() {
		var c models.Contact
		var tags pq.StringArray
		var contact sql.NullString

		err = rows.Scan(&c.Id, &c.Name, &c.LinkedIn, &c.Credly, &tags, &contact)
		if err != nil {
			return
		}

		c.Tags = tags
		c.Contact = contact.String
		contacts = append(contacts, c)
	}

	countQuery := "SELECT COUNT(*) FROM contacts"
	if len(terms) > 0 {
		countQuery += fmt.Sprintf(" WHERE tags %s $1", operator)
		err = s.db.QueryRow(countQuery, pq.Array(terms)).Scan(&count)
	} else {
		err = s.db.QueryRow(countQuery).Scan(&count)
	}

	return
}

func (s *service) CreateContact(contact models.Contact) (models.Contact, error) {
	query := `
		INSERT INTO contacts (name, linkedinurl, credlyurl, datecreated, dateupdated, tags, contact)
		VALUES($1, $2, $3, $4, NOW(), NOW(), $5)
		RETURNING id, name, linkedinurl, credlyurl, datecreated, dateupdated, tags, contact
	`

	var createdContact models.Contact
	err := s.db.QueryRow(query, contact.Name, contact.LinkedIn, contact.Credly, contact.Tags, contact.Contact).
		Scan(&createdContact.Id, &createdContact.Name, &createdContact.LinkedIn, &createdContact.Credly, &createdContact.Tags, &createdContact.Contact)
	if err != nil {
		return models.Contact{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return createdContact, nil
}

func (s *service) EditContact(contact models.Contact) error {
	query := `
		UPDATE contacts
		SET
			name = $2,
			linkedinurl = $3,
			credlyurl = $4,
			dateupdated = NOW(),
			tags = $5,
			contact = $6
		WHERE id = $1
		`

	res, err := s.db.Exec(query, contact.Id, contact.Name, contact.LinkedIn, contact.Credly, contact.Tags, contact.Contact)
	if err != nil {
		return fmt.Errorf("failed to update contact: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no contact updated with id %d", contact.Id)
	}

	return nil
}

func (s *service) DeleteContact(id string) (err error) {
	res, err := s.db.Exec("DELETE FROM contacts WHERE id = $1", id)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no contact deleted with id %s", id)
	}

	return
}
