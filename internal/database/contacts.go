package database

import (
	"Contacter/internal/models"
	"fmt"
	"time"
)

func (s *service) GetContacts() (contacts []models.Contact, err error) {
	rows, err := s.db.Query("SELECT * FROM contacts ORDER BY Id")
	defer rows.Close()

	if err != nil {
		return
	}

	for rows.Next() {
		var id int
		var name string
		var linkedIn string
		var credly string
		var dateCreated time.Time
		var dateUpdated time.Time
		var tags []string
		var contact string

		err = rows.Scan(&id, &name, &linkedIn, &credly, &dateCreated, &dateUpdated, &tags, &contact)

		if err != nil {
			return
		}

		contacts = append(contacts, models.Contact{
			Id:          id,
			Name:        name,
			LinkedIn:    linkedIn,
			Credly:      credly,
			DateCreated: dateCreated,
			DateUpdated: dateUpdated,
			Tags:        tags,
			Contact:     contact,
		})
	}

	return contacts, nil
}

func (s *service) CreateContact(contact models.Contact) (models.Contact, error) {
	query := `
		INSERT INTO contacts (name, linkedinurl, credlyurl, datecreated, dateupdated, tags, contact)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, linkedinurl, credlyurl, datecreated, dateupdated, tags, contact
	`

	var createdContact models.Contact
	err := s.db.QueryRow(query, contact.Name, contact.LinkedIn, contact.Credly, contact.DateCreated, contact.DateUpdated, contact.Tags, contact.Contact).
		Scan(&createdContact.Id, &createdContact.Name, &createdContact.LinkedIn, &createdContact.Credly, &createdContact.DateCreated, &createdContact.DateUpdated, &createdContact.Tags, &createdContact.Contact)
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
			datecreated = $5,
			dateupdated = $6,
			tags = $7,
			contact = $8
		WHERE id = $1
		`

	res, err := s.db.Exec(query, contact.Id, contact.Name, contact.LinkedIn, contact.Credly, contact.DateCreated, contact.DateUpdated, contact.Tags, contact.Contact)
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
