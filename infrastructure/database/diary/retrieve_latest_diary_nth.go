package diary

func (db *DiaryRepository) RetrieveLatestDiaryNth() (nth int, err error) {
	row := db.QueryRowx(`SELECT nth FROM game_status;`)
	err = row.Scan(&nth)
	return
}
