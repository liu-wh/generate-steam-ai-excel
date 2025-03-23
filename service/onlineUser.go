package service

import (
	"encoding/json"
	"generate-steam-ai-excel/code"
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/util"
	"strconv"
)

type SteamGame struct {
	Name            string
	ID              int
	Price           string
	Logo            string
	Count           int
	IsFree          bool
	DiscountPercent int
	IsCheapest      bool
	GameType        any
}

func GenerateOnlineUserExcel() {
	var (
		onlineUserStr string
		err           error
	)
	if onlineUserStr, err = global.R.Get(global.CTX, "SteamGameTop").Result(); err != nil {
		return
	}
	gameList := make([]*SteamGame, 0, 3000)
	if err = json.Unmarshal(util.Str2bytes(onlineUserStr), &gameList); err != nil {
		global.Logger.Error("解析在线用户数据失败", code.ERROR, err)
		return
	}
	i := 2
	for _, game := range gameList {
		_ = global.OnlineUserFile.SetSheetRow(code.SHEET1, "A"+strconv.Itoa(i), &[]any{
			game.Name,
			game.Count,
		})
		i += 1
	}
}
