package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"log"
	"magic-mouse-engine/model"
)

const getTemplateValueQuery = "select id, description from template where id = $1"
const getAllTemplateValueQuery = "select id, description from template"
const insertTemplateValueQuery = "insert into template(id, description) values ($1, $2)"
const updateTemplateValueQuery = "update template set description = $2 where id = $1"
const deleteTemplateValueQuery = "delete from template where id = $1"
const getTemplatesByKeywordQuery = "select id, description from template where description ilike '%' || $1 || '%'"

// GetTemplateServiceEntry is ...
func (mdb *Database) GetCard(ctx context.Context, lookupRequest model.CardRequestHeader) (*model.CardAPI, error) {
	response := new(model.CardAPI)

	var err error

	var id = lookupRequest.Id

	err = mdb.Db.QueryRowContext(ctx, getTemplateValueQuery, id).
		Scan(&response.ID)

	if err != nil {
		log.Printf(databaseErrorOccurred, err)
		return response, err
	}

	return response, err
}

func (mdb *Database) GetAllCards(ctx context.Context, lookupRequest model.CardRequestHeader) ([]model.CardAPI, error) {

	// A slice to hold data from returned rows.
	var response []model.CardAPI
	var rows, err = mdb.Db.QueryContext(ctx, getAllTemplateValueQuery)
	if err != nil {
		log.Printf(databaseErrorOccurred, err)
	}
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		aux := new(model.CardAPI)
		if err := rows.Scan(&aux.ID); err != nil {
			return response, err
		}
		response = append(response, *aux)
	}
	if err = rows.Err(); err != nil {
		return response, err
	}
	return response, err
}

// PostTemplateServiceEntry is ...
func (mdb *Database) InsertCard(ctx context.Context, insertRequest model.CardAPI) error {
	_, err := mdb.Db.ExecContext(ctx, insertTemplateValueQuery, insertRequest.Id, insertRequest.Description)

	if err != nil {
		log.Printf(databaseErrorOccurred, err)
		return err
	}

	return err
}

// PatchTemplateServiceEntry is ...
func (mdb *Database) UpdateCard(ctx context.Context, updateRequest model.CardAPI) error {

	res, err := mdb.Db.ExecContext(ctx, updateTemplateValueQuery, updateRequest.Id, updateRequest.Description)

	if err != nil {
		log.Printf(databaseErrorOccurred, err)
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}

// DeleteTemplateServiceEntry is ...
func (mdb *Database) DeleteCard(ctx context.Context, deleteRequest model.CardRequestHeader) error {

	res, err := mdb.Db.ExecContext(ctx, deleteTemplateValueQuery, deleteRequest.Id)

	if err != nil {
		log.Printf(databaseErrorOccurred, err)
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return err
}

func (mdb *Database) FindByNameAndVersion(ctx context.Context, lookup model.CardCollectionByKeywordHeader) (*model.CardsCollectionBody, error) {

	response := new(model.CardsCollectionBody)

	var keyword = lookup.Keyword

	var rows, err = mdb.Db.QueryContext(ctx, getTemplatesByKeywordQuery, keyword)
	var templates []model.CardAPI

	if err != nil {
		log.Printf(databaseErrorOccurred, err)
	}

	if err = rows.Err(); err != nil {
		return response, err
	}

	for rows.Next() {
		template := new(model.CardAPI)
		if err := rows.Scan(&template.ID, &template.Description); err != nil {
			return response, err
		}
		templates = append(templates, *template)
	}

	response.Cards = templates
	response.Total = len(response.Cards)

	return response, err
}



func (mdb *Database) FindByFilter(ctx context.Context, filter model.FilterCardRequest) (*model.CardsCollectionBody, error) {

	response := new(model.CardsCollectionBody)

	query := "SELECT * FROM template WHERE 1=1"
args := []any{}
argPos := 1

if filter.Set.ID != "" {
    query += fmt.Sprintf(" AND set_id = $%d", argPos)
    args = append(args, filter.Set.ID)
    argPos++
}

if filter.Name != "" {
    query += fmt.Sprintf(" AND name ILIKE $%d", argPos)
    args = append(args, "%"+filter.Name+"%")
    argPos++
}

if filter.Cost != 0 {
    query += fmt.Sprintf(" AND cost = $%d", argPos)
    args = append(args, filter.Cost)
    argPos++
}

if filter.Inkwell != nil {
    query += fmt.Sprintf(" AND inkwell = $%d", argPos)
    args = append(args, *filter.Inkwell)
    argPos++
}

if len(filter.Type) > 0 {
    query += fmt.Sprintf(
        " AND type IN (%s)",
        filter.Type,
    )
}

if filter.Rarity != "" {
    query += fmt.Sprintf(" AND rarity = $%d", argPos)
	args = append(args, filter.Rarity)
	argPos++
}

	var rows, err = mdb.Db.QueryContext(ctx, query, args...)
	var templates []model.CardAPI

	if err != nil {
		log.Printf(databaseErrorOccurred, err)
	}

	if err = rows.Err(); err != nil {
		return response, err
	}

	for rows.Next() {
		template := new(model.CardAPI)
		if err := rows.Scan(&template.ID, &template.Description); err != nil {
			return response, err
		}
		templates = append(templates, *template)
	}

	response.Cards = templates
	response.Total = len(response.Cards)

	err = sqlx.GetContext(ctx, mdb.Dbx, &response.Cards, query, args...)
	return response, err
}