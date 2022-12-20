package mail

import model "github.com/kaikourok/lunchtote-backend/entity/model"

func (db *MailRepository) RetrieveMails(characterId int, unreadedOnly bool, limit int, start int) (mails *[]model.ReceivedMail, isContinue bool, err error) {
	sql := `
		SELECT
			mails.id,
			characters.id,
			characters.nickname,
			mails.name,
			mails.title,
			mails.message,
			mails.posted_at,
			mails.read
		FROM
			mails
		LEFT JOIN
			characters ON (mails.sender = characters.id)
		WHERE
			mails.receiver = $1 AND
			mails.id < $3`

	if unreadedOnly {
		sql += ` AND mails.read = false`
	}

	sql += `
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

	mailsSlice := make([]model.ReceivedMail, 0, limit)

	for rows.Next() {
		var mail model.ReceivedMail
		var characterNickname *string
		var mailName string
		err = rows.Scan(
			&mail.Id,
			&mail.Sender.Id,
			&characterNickname,
			&mailName,
			&mail.Title,
			&mail.Message,
			&mail.Timestamp,
			&mail.Read,
		)
		if err != nil {
			return nil, false, err
		}

		if characterNickname != nil {
			mail.Sender.Name = *characterNickname
		} else {
			mail.Sender.Name = mailName
		}

		mailsSlice = append(mailsSlice, mail)
	}

	if len(mailsSlice) == limit+1 {
		isContinue = true
		mailsSlice = mailsSlice[:len(mailsSlice)-1]
	}

	return &mailsSlice, isContinue, nil
}
