package general

func (db *GeneralRepository) UpdateInquiryState(inquiryId int, resolved bool) error {
	_, err := db.Exec(`
		UPDATE
			inquiries
		SET
			resolved = $2
		WHERE
			id = $1;
	`, inquiryId, resolved)

	return err
}
