package mail

func (db *MailRepository) SetMailRead(characterId int, mailId int) (existsUnreadMail bool, err error) {
	_, err = db.Exec(`
		UPDATE
			mails
		SET
			read = true
		WHERE
			id = $1 AND receiver = $2 AND deleted_at IS NULL;
	`, mailId, characterId)
	if err != nil {
		return false, err
	}

	row := db.QueryRowx(`
		SELECT 
			EXISTS (
				SELECT
					*
				FROM
					mails
				WHERE
					receiver = $1    AND
					read     = false AND
					deleted_at IS NULL
				);
	`, characterId)

	err = row.Scan(&existsUnreadMail)
	if err != nil {
		return false, err
	}

	return existsUnreadMail, nil
}
