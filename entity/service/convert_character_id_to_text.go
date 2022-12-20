package service

import "strconv"

func ConvertCharacterIdToText(id int) string {
	return "#" + strconv.Itoa(id)
}
