package mail

import model "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *MailRepository) RetrieveSentMails(characterId int, limit int, start int) (mails *[]model.SentMail, isContinue bool, err error) {
	sql := `
		SELECT
			mails.id,
			characters.id,
			characters.nickname,
			mails.title,
			mails.message,
			mails.posted_at,
			mails.read
		FROM
			mails
		LEFT JOIN
			characters ON (mails.receiver = characters.id)
		WHERE
			mails.sender = $1 AND
			mails.id < $3
		ORDER BY
			mails.posted_at DESC
		LIMIT
			$2;
	`

	rows, err := db.Queryx(sql, characterId, limit+1, start)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	mailsSlice := make([]model.SentMail, 0, limit)

	for rows.Next() {
		var mail model.SentMail
		err = rows.Scan(
			&mail.Id,
			&mail.Receiver.Id,
			&mail.Receiver.Name,
			&mail.Title,
			&mail.Message,
			&mail.Timestamp,
			&mail.Read,
		)
		if err != nil {
			return nil, false, err
		}

		mailsSlice = append(mailsSlice, mail)
	}

	if len(mailsSlice) == limit+1 {
		isContinue = true
		mailsSlice = mailsSlice[:len(mailsSlice)-1]
	}

	return &mailsSlice, isContinue, nil
}
