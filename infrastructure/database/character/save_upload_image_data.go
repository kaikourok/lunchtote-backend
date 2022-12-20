package character

import (
	"bytes"
	"io"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/kaikourok/lunchtote-backend/entity/service"
)

// TODO: 明らかに責任範囲外を担当しているのでdomainとinterface/storageに適切に分割する
func (db *CharacterRepository) SaveUploadImageData(id int, images []*bytes.Buffer, imageType service.ImageTypeId, uploadDir string) (*[]string, error) {
	// アイコンファイルをストレージに保存
	// 同時にDBに画像ファイルのパス及びハッシュ値を保存
	// 同一キャラクターから同一ハッシュの画像が送られてきた場合、既存のファイルのpathを返す

	type FileInfo struct {
		Character       int    `db:"character"`
		Path            string `db:"path"`
		MD5             string `db:"md5"`
		Buffer          *bytes.Buffer
		Extension       string
		DuplicatedPast  bool
		DuplicatedIndex int
	}

	fileInfos := make([]*FileInfo, len(images))
	hashes := make([]string, len(images))

	for i, image := range images {
		convertedImage, extension, err := service.ConvertImage(image, imageType)
		if err != nil {
			return nil, err
		}

		hash, err := service.CalcHash(convertedImage)
		if err != nil {
			return nil, err
		}

		fileInfos[i] = &FileInfo{
			Character:       id,
			MD5:             hash,
			Buffer:          convertedImage,
			Extension:       extension,
			DuplicatedIndex: -1,
		}

		hashes[i] = hash

		// 同一アップロード内でハッシュが被っているものに被り先のインデックスを付与
		// かぶり先インデックスとして設定されるのは配列の中で一番インデックスが小さいもの
		for j := 0; j < i; j++ {
			if fileInfos[j].MD5 == fileInfos[i].MD5 {
				fileInfos[i].DuplicatedIndex = j
				break
			}
		}
	}

	// ハッシュが被っているものを検索
	argInterfaces := make([]interface{}, len(hashes)+1)
	argInterfaces[0] = id
	for i, v := range hashes {
		argInterfaces[i+1] = v
	}

	sqlBuilder := `
		SELECT
			path,
			md5
		FROM
			characters_uploaded_images
		WHERE
			character = $1 AND
			md5 IN (
	`

	for i := 1; i-1 < len(hashes); i++ {
		if i != 1 {
			sqlBuilder += ", "
		}
		sqlBuilder += "$" + strconv.Itoa(i+1)
	}

	sqlBuilder += ");"

	rows, err := db.Queryx(sqlBuilder, argInterfaces...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 過去のハッシュと重複していたものに重複フラグとパスを設定し、重複数をカウント
	duplicates := 0
	for rows.Next() {
		var path, hash string
		err = rows.Scan(&path, &hash)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(fileInfos); i++ {
			if hash == fileInfos[i].MD5 {
				fileInfos[i].Path = path
				fileInfos[i].DuplicatedPast = true
			}
		}

		duplicates++
	}

	// 未アップロードのものがあれば
	if 0 < len(images)-duplicates {
		err = os.MkdirAll(uploadDir+"/"+strconv.Itoa(id), 0777)
		if err != nil {
			return nil, err
		}

		unuploadedFileInfos := make([]*FileInfo, 0, len(images)-duplicates)

		// 重複フラグのついていないものを保存し、パスを付ける
		for _, fileInfo := range fileInfos {
			if !fileInfo.DuplicatedPast && fileInfo.DuplicatedIndex == -1 {
				path := "/" + strconv.Itoa(id) + "/" + uuid.NewString() + "." + fileInfo.Extension

				saveImage, err := os.Create(uploadDir + path)
				if err != nil {
					return nil, err
				}
				defer saveImage.Close()

				_, err = io.Copy(saveImage, fileInfo.Buffer)
				if err != nil {
					return nil, err
				}

				fileInfo.Path = path

				unuploadedFileInfos = unuploadedFileInfos[:len(unuploadedFileInfos)+1]
				unuploadedFileInfos[len(unuploadedFileInfos)-1] = &FileInfo{
					Character: fileInfo.Character,
					Path:      fileInfo.Path,
					MD5:       fileInfo.MD5,
				}
			}
		}

		// 重複していないファイルについての情報をDBに保存
		if 0 < len(unuploadedFileInfos) {
			_, err = db.NamedExec(`
				INSERT INTO characters_uploaded_images (
					character,
					path,
					md5
				) VALUES (
					:character,
					:path,
					:md5
				)
			`, unuploadedFileInfos)
			if err != nil {
				return nil, err
			}
		}
	}

	// 同一アップロード内でインデックスが被っていたものについてpathをかぶり先から設定
	for _, fileInfo := range fileInfos {
		if fileInfo.DuplicatedIndex != -1 {
			fileInfo.Path = fileInfos[fileInfo.DuplicatedIndex].Path
		}
	}

	paths := make([]string, len(images))
	for i, fileInfo := range fileInfos {
		paths[i] = fileInfo.Path
	}

	return &paths, nil
}
