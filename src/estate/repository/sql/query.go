package sql

const (
	SelectTemplate = `SELECT
		uuid,
		length,
		width
	FROM
		estate`

	QueryGetByUuid = SelectTemplate + `
	WHERE
		uuid = $1`

	QueryCreateEstate = `INSERT INTO estate
	(uuid, length, width, createdAt)
	VALUES($1, $2, $3, $4)`
)
