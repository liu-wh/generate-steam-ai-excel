package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"generate-steam-ai-excel/code"
	"generate-steam-ai-excel/database"
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/models"
	"generate-steam-ai-excel/util"
	"github.com/redis/go-redis/v9"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Cheapest struct {
	TimeStamp       int     `json:"TimeStamp"`
	Price           float64 `json:"Price"`
	DiscountPercent int     `json:"DiscountPercent"`
}

func init() {
	global.Logger = slog.Default()
	var (
		err     error
		rateMap map[string]string
		f       float64
	)
	global.DB, err = database.SetUpDB()
	if err != nil {
		global.Logger.Error("连接数据库失败", code.ERROR, err)
		os.Exit(1)
	}
	if err = database.SetUpRedis(); err != nil {
		os.Exit(1)
	}

	//游戏名翻译字典
	if global.GameNameMap, err = global.R.HGetAll(global.CTX, "SteamTranslateName").Result(); err != nil {
		global.Logger.Error("获取游戏翻译字典失败", code.ERROR, err)
		os.Exit(1)
	}

	//取汇率
	if rateMap, err = global.R.HGetAll(global.CTX, "ExchangeRateMap").Result(); err != nil {
		global.Logger.Error("获取汇率字典失败", code.ERROR, err)
		os.Exit(1)
	}
	for k, v := range rateMap {
		f, err = strconv.ParseFloat(v, 64)
		if err != nil {
			global.Logger.Error("解析汇率为float64失败", code.ERROR, err)
			os.Exit(1)
		}
		global.ExchangeRateMap[k] = f
	}
	//取游戏详情字典
	if global.GameList, err = global.R.HKeys(global.CTX, "SteamGameStoreDetailData").Result(); err != nil {
		global.Logger.Error("获取游戏详情字典失败", code.ERROR, err)
		os.Exit(1)
	}

	//创建Excel文件对象
	global.F = excelize.NewFile()

	_ = global.F.SetSheetRow("Sheet1", "A1", &[]string{"游戏名", "折扣率", "国区原价", "国区现价", "史低价格", "史低时间", "俄罗斯原价", "俄罗斯现价", "巴西原价", "巴西现价", "阿根廷原价", "阿根廷现价", "哈萨克斯坦原价", "哈萨克斯坦现价", "印度原价", "印度现价", "乌克兰原价", "乌克兰现价", "土耳其原价", "土耳其现价", "马来西亚原价", "马来西亚现价", "越南原价", "越南现价", "巴基斯坦原价", "巴基斯坦现价", "印尼原价", "印尼现价", "菲律宾原价", "菲律宾现价", "智利原价", "智利现价", "哥伦比亚原价", "哥伦比亚现价", "南非原价", "南非现价", "泰国原价", "泰国现价", "乌拉圭原价", "乌拉圭现价", "墨西哥原价", "墨西哥现价", "科威特原价", "科威特现价", "挪威原价", "挪威现价", "卡塔尔原价", "卡塔尔现价", "日本原价", "日本现价", "独联体原价", "独联体现价", "秘鲁原价", "秘鲁现价", "新西兰原价", "新西兰现价", "新加坡原价", "新加坡现价", "澳大利亚原价", "澳大利亚现价", "韩国原价", "韩国现价", "中国台湾原价", "中国台湾现价", "加拿大原价", "加拿大现价", "阿拉伯联合酋长国原价", "阿拉伯联合酋长国现价", "沙特阿拉伯原价", "沙特阿拉伯现价", "中国香港原价", "中国香港现价", "哥斯达黎加原价", "哥斯达黎加现价", "波兰原价", "波兰现价", "英国原价", "英国现价", "以色列原价", "以色列现价", "美国原价", "美国现价", "欧盟原价", "欧盟现价", "瑞士原价", "瑞士现价"})
}

func releaseResource() {
	db, _ := global.DB.DB()
	if err := db.Close(); err != nil {
		global.Logger.Error("关闭数据库失败", code.ERROR, err)
		os.Exit(1)
	}
	if err := global.F.SaveAs(fmt.Sprintf("steam_price_%s.xlsx", time.Now().Format(time.DateOnly))); err != nil {
		global.Logger.Error("保存Excel失败", code.ERROR, err)
		os.Exit(1)
	}
}

func main() {
	defer releaseResource()
	idx := 2
	A := "A"
	//拿到所有游戏的价格
	for j, gameID := range global.GameList {
		if j == 100 {
			break
		}
		flag := false
		gameInfo := make([]any, 0, 84)
		for i := range 41 {
			i += 1
			_price := models.SteamGamePrice{}
			if err := global.DB.Preload("SteamGame").Where("steam_game_id  = ?", gameID).Where("steam_location_id = ?", i).Find(&_price).Error; err != nil {
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
				_date := time.Unix(int64(_cc.TimeStamp), 0)
				gameInfo = append(gameInfo, _cc.Price, _date.Format(time.DateOnly))
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
			gameInfo = append(gameInfo, fmt.Sprintf("%.2f", initP), fmt.Sprintf("%.2f", finalP))
		}
		if flag {
			continue
		}
		if err := global.F.SetSheetRow("Sheet1", A+strconv.Itoa(idx), &gameInfo); err != nil {
			global.Logger.Error("写入Excel失败", code.ERROR, err)
		}
		idx += 1

	}

	//gamePriceList := make([]*models.SteamGamePrice, 0, 335338)
	//global.DB.Preload("SteamLocation").Preload("SteamGame").Find(&gamePriceList)
	//for i, j := range gamePriceList {
	//	if i > 5 {
	//		break
	//	}
	//	gameName := util.GetGameName(j)
	//}

}
