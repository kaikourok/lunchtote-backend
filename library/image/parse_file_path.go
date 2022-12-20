package image

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func ParseFilePath(path string) (characterId int, err error) {
	splitedPaths := strings.Split(path, "/")

	if len(splitedPaths) != 3 {
		return 0, errors.New("パスの形式が不正です")
	}

	if splitedPaths[0] != "" {
		return 0, errors.New("パスの形式が不正です")
	}

	characterId, err = strconv.Atoi(splitedPaths[1])
	if err != nil {
		return 0, errors.New("パスの形式が不正です")
	}

	splitedFilenames := strings.Split(splitedPaths[2], ".")

	if len(splitedFilenames) != 2 {
		return 0, errors.New("パスの形式が不正です")
	}

	_, err = uuid.Parse(splitedFilenames[0])
	if err != nil {
		return 0, errors.New("パスの形式が不正です")
	}

	if splitedFilenames[1] != "png" && splitedFilenames[1] != "gif" && splitedFilenames[1] != "jpg" && splitedFilenames[1] != "jpeg" {
		return 0, errors.New("パスの形式が不正です")
	}

	return characterId, nil
}
