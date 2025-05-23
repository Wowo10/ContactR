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

	orderLimitQuery := fmt.Sprintf(" ORDER BY ID LIMIT %d OFFSET %d", PAGE_SIZE, page*PAGE_SIZE)

	var rows *sql.Rows
	if len(terms) == 0 {
		rows, err = s.db.Query(baseQuery + orderLimitQuery)
	} else {
		filterQuery := fmt.Sprintf("%s WHERE tags %s $1 %s", baseQuery, operator, orderLimitQuery)
		rows, err = s.db.Query(filterQuery, pq.Array(terms))
	}

	if err != nil {
		return
	}
	defer rows.Close()

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

func (s *service) GetContact(id string) (contact models.Contact, err error) {
	query := `
		SELECT *
		FROM contacts
		WHERE id = $1
	`

	var tags pq.StringArray
	var contactStr sql.NullString
	err = s.db.QueryRow(query, id).
		Scan(&contact.Id, &contact.Name, &contact.LinkedIn, &contact.Credly,
			&contact.DateCreated, &contact.DateUpdated, &tags, &contactStr)

	contact.Tags = tags
	if contactStr.Valid {
		contact.Contact = contactStr.String
	}

	return
}

func (s *service) CreateContact(contact models.Contact) (models.Contact, error) {
	query := `
		INSERT INTO contacts (name, linkedinurl, credlyurl, datecreated, dateupdated, tags, contact)
		VALUES($1, $2, $3, NOW(), NOW(),$4, $5)
		RETURNING id, name, linkedinurl, credlyurl, datecreated, dateupdated, tags, contact
	`

	var createdContact models.Contact
	var tags pq.StringArray

	err := s.db.QueryRow(query, contact.Name, contact.LinkedIn, contact.Credly, contact.Tags, contact.Contact).
		Scan(&createdContact.Id, &createdContact.Name, &createdContact.LinkedIn, &createdContact.Credly,
			&createdContact.DateCreated, &createdContact.DateUpdated, &tags, &createdContact.Contact)
	if err != nil {
		return models.Contact{}, fmt.Errorf("failed to insert user: %w", err)
	}

	createdContact.Tags = tags

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
