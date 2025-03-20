package util

import (
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/models"
	"strconv"
)

func GetGameName(game *models.SteamGamePrice) string {
	gameName, ok := global.GameNameMap[strconv.Itoa(int(game.ID))]
	if ok {
		return gameName
	}
	return game.SteamGame.EName
}
