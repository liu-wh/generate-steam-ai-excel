package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"generate-steam-ai-excel/code"
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/models"
	"generate-steam-ai-excel/util"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Cheapest struct {
	TimeStamp       int     `json:"TimeStamp"`
	Price           float64 `json:"Price"`
	DiscountPercent int     `json:"DiscountPercent"`
}

func GeneratePriceExcel() {
	idx := 2
	A := "A"
	//拿到所有游戏的价格
	for j, gameID := range global.GameList {
		if j == 1300 {
			break
		}
		flag := false
		gameInfo := make([]any, 0, 84)
		for i := range 41 {
			i += 1

			_storeData := models.SteamGameStoreData{}
			var (
				_storeDataStr string
				err           error
			)
			if _storeDataStr, err = global.R.HGet(global.CTX, "SteamGameStoreDetailData", gameID).Result(); err != nil {
				global.Logger.Error("查询游戏详情失败", code.ERROR, err, "游戏ID", gameID)
				flag = true
				break
			} else {
				if err = json.Unmarshal(util.Str2bytes(_storeDataStr), &_storeData); err != nil {
					global.Logger.Error("解析游戏详情失败", code.ERROR, err, "游戏ID", gameID)
					flag = true
					break
				} else {
					if !_storeData.Success {
						flag = true
						break
					}
					if _storeData.Data.IsFree {
						gameInfo = append(gameInfo, util.GetGameName(&models.SteamGamePrice{SteamGameID: uint(_storeData.Data.SteamAppid)}), " ", " ", "免费")
						idx += 1
						flag = true
						if err = global.F.SetSheetRow(code.SHEET1, A+strconv.Itoa(idx), &gameInfo); err != nil {
							global.Logger.Error("写入Excel失败", code.ERROR, err)
						}
						break
					}
				}
			}

			_price := models.SteamGamePrice{}
			if err = global.DB.Preload("SteamGame").Where("steam_game_id  = ?", gameID).Where("steam_location_id = ?", i).First(&_price).Error; err != nil {
				global.Logger.Error("查询steam游戏价格失败", code.ERROR, err, "游戏ID", gameID, "区ID", i)
				continue
			}
			if i == 1 {
				if _price.Initial == 0 && _price.Final == 0 {
					flag = true
					break
				}
				gameInfo = append(gameInfo, util.GetGameName(&_price), _price.DiscountPercent, _price.Initial/100, _price.Final/100)
				_c, err := global.R.HGet(global.CTX, "SteamGamePriceCheapest", gameID).Result()
				if err != nil {
					if errors.Is(err, redis.Nil) {
						gameInfo = append(gameInfo, " ", " ")
						continue
					}
					global.Logger.Error("查询游戏史低价格失败", code.ERROR, err, "游戏ID", gameID)
					continue
				}
				_cc := Cheapest{}
				if err = json.Unmarshal([]byte(_c), &_cc); err != nil {
					global.Logger.Error("解析游戏史低价格失败", code.ERROR, err, "游戏ID", gameID)
					continue
				}
				var (
					_date string
					_p    float64
				)
				if _cc.TimeStamp == 88 {
					_date = " "
				} else {
					_date = time.Unix(int64(_cc.TimeStamp), 0).Format(time.DateOnly)
				}
				_p = _cc.Price

				gameInfo = append(gameInfo, _p, _date)
				continue
			}
			_location := models.SteamLocation{}
			if err := global.DB.Where("id = ?", i).Find(&_location).Error; err != nil {
				global.Logger.Error("查询steam区", code.ERROR, err, "区ID", i)
				continue
			}
			exchangeRate := global.ExchangeRateMap[_location.CurrencyCode]
			initP := (float64(_price.Initial) / 100) * exchangeRate
			finalP := (float64(_price.Final) / 100) * exchangeRate
			initStr := fmt.Sprintf("%.2f", initP)
			finalStr := fmt.Sprintf("%.2f", finalP)
			if initP == 0 {
				initStr = "无"
			}
			if finalP == 0 {
				finalStr = "无"
			}
			gameInfo = append(gameInfo, initStr, finalStr)
		}
		if flag {
			continue
		}
		if err := global.F.SetSheetRow(code.SHEET1, A+strconv.Itoa(idx), &gameInfo); err != nil {
			global.Logger.Error("写入Excel失败", code.ERROR, err)
		}
		idx += 1

	}
}
