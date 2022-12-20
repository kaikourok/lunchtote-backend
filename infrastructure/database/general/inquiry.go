package general

func (db *GeneralRepository) Inquiry(characterId *int, inquiry string) error {
	_, err := db.Exec(`
		INSERT INTO inquiries (
			character,
			inquiry
		) VALUES (
			$1,
			$2
		);
	`, characterId, inquiry)
	return err
}
