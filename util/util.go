package util

import (
	"crypto/md5"
	"fmt"
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/models"
	"io"
	"os"
	"strconv"
	"unsafe"
)

func GetGameName(game *models.SteamGamePrice) string {
	gameName, ok := global.GameNameMap[strconv.Itoa(int(game.SteamGameID))]
	if ok {
		return gameName
	}
	return game.SteamGame.EName
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func GetFileSize(filePath string) (int64, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return file.Size(), nil
}
