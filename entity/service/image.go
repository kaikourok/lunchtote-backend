package service

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"image"
	"image/png"
	"io"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type ImageTypeId string

type ImageTypeStruct struct {
	General      ImageTypeId
	Icon         ImageTypeId
	IconFragment ImageTypeId
	ProfileImage ImageTypeId
	ListImage    ImageTypeId
}

var ImageType ImageTypeStruct

func init() {
	ImageType = ImageTypeStruct{
		General:      "GENERAL",
		Icon:         "ICON",
		IconFragment: "ICON_FRAGMENT",
		ProfileImage: "PROFILE_IMAGE",
		ListImage:    "LIST_IMAGE",
	}
}

const profileImageHeight = 700

func ConvertImage(imageBuffer *bytes.Buffer, convertType ImageTypeId) (*bytes.Buffer, string, error) {
	// 画像が適合規格であれば画像をそのまま返し、
	// 規格不適合な画像であれば変換、非対応フォーマット等であればエラーを返す関数
	// 対応フォーマット JPEG PNG    対応予定だけどまだ：GIF APNG WebP

	// アイコンの場合
	// 不適合条件
	//  ・画像のサイズが正方形でない
	//  ・一辺のサイズが240pxを超えている
	// 変換先フォーマット
	//  PNG 120px平方 Bilinear

	// キャラクターリスト画像の場合
	// 不適合条件
	//　・画像のサイズが横180px高さ240pxでない
	// 変換先フォーマット
	//　PNG 横180px 高さ240px Bilinear

	// プロフィール画像の場合
	// 不適合条件
	//  ・高さが700pxでない
	// 変換先フォーマット
	//  PNG 高さ700px 横幅アスペクト比保持伸縮 Bilinear

	// 単純にimageBufferを読み出すとio.Readerが変わるためひと手間加える
	image, extension, err := image.Decode(bytes.NewReader(imageBuffer.Bytes()))
	if err != nil {
		return nil, "", err // 画像フォーマットが非対応等であればここでエラーが返る
	}

	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	if convertType == ImageType.Icon && (imageWidth != imageHeight || 240 < imageWidth || 240 < imageHeight) {
		// 画像サイズが不適合であれば変換（アイコン）
		resizedImage := resize.Resize(120, 120, image, resize.Bilinear)

		resizedImageBuffer := &bytes.Buffer{}
		err = png.Encode(resizedImageBuffer, resizedImage)
		if err != nil {
			return nil, "", err
		}

		return resizedImageBuffer, "png", nil
	} else if convertType == ImageType.ListImage && (imageWidth != 180 || imageHeight != 240) {
		// 画像サイズが不適合であれば変換（キャラクターリスト画像）
		resizedImage := resize.Resize(180, 240, image, resize.Bilinear)

		resizedImageBuffer := &bytes.Buffer{}
		err = png.Encode(resizedImageBuffer, resizedImage)
		if err != nil {
			return nil, "", err
		}

		return resizedImageBuffer, "png", nil
	} else if convertType == ImageType.ProfileImage && imageHeight != profileImageHeight {
		// 画像サイズが不適合であれば変換（プロフィール画像）
		resizedImage := resize.Resize(uint(imageWidth*profileImageHeight/imageHeight), profileImageHeight, image, resize.Bilinear)

		resizedImageBuffer := &bytes.Buffer{}
		err = png.Encode(resizedImageBuffer, resizedImage)
		if err != nil {
			return nil, "", err
		}

		return resizedImageBuffer, "png", nil
	} else {
		// 画像サイズに問題がなければそのままバッファを返す
		return imageBuffer, extension, nil
	}
}

// バッファのハッシュ値を計算してHEXで返す
func CalcHash(buffer *bytes.Buffer) (string, error) {
	hash := md5.New()
	_, err := io.Copy(hash, bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return "", err
	}

	return string(hex.EncodeToString(hash.Sum(nil))), nil
}

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
