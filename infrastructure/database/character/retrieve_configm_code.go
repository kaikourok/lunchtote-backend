package character

func (db *CharacterRepository) RetrieveConfirmCode(id int) (code, email string, err error) {
	row := db.QueryRowx(`
		SELECT
			code,
			email
		FROM
			mail_confirm_codes
		WHERE
			character = $1 AND
			CURRENT_TIMESTAMP < expire
		ORDER BY
			id DESC
		LIMIT
			1;
	`, id)

	err = row.Scan(&code, &email)
	return
}
