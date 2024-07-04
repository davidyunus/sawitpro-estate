package sql

const (
	SelectTemplate = `SELECT
		id,
		uuid,
		x,
		y,
		height
	FROM
		palmTreeLocation`

	QueryGetByUuid = SelectTemplate + `
	WHERE
		uuid = $1`

	QueryPlantPalmTree = `INSERT INTO palmTreeLocation
	(uuid, x, y, height, createdAt)
	VALUES($1, $2, $3, $4, $5)`
)
