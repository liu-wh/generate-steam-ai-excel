package util

import (
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/models"
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
