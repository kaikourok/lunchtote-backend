package general

func (db *GeneralRepository) Connections() (connections int, err error) {
	row := db.QueryRowx(`
		SELECT COUNT(*) FROM pg_stat_activity;
	`)
	err = row.Scan(&connections)

	return
}
