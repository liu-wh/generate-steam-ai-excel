package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"generate-steam-ai-excel/code"
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/models"
	"generate-steam-ai-excel/util"
	bailian20231229 "github.com/alibabacloud-go/bailian-20231229/client"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
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
	for _, gameID := range global.GameList {
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
						gameName := util.GetGameName(&models.SteamGamePrice{SteamGameID: uint(_storeData.Data.SteamAppid)})
						if gameName == "" {
							gameName = _storeData.Data.Name
						}
						gameInfo = append(gameInfo, gameName, " ", " ", "免费")
						flag = true
						if err = global.F.SetSheetRow(code.SHEET1, A+strconv.Itoa(idx), &gameInfo); err != nil {
							global.Logger.Error("写入Excel失败", code.ERROR, err)
						}
						idx += 1
						break
					}
				}
			}

			_price := models.SteamGamePrice{}
			if err = global.DB.Preload("SteamGame").Where("steam_game_id  = ?", gameID).Where("steam_location_id = ?", i).First(&_price).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					global.Logger.Error("查询steam游戏价格失败", code.ERROR, err, "游戏ID", gameID, "区ID", i)
				}
				continue
			}
			if i == 1 {
				if _price.Initial == 0 && _price.Final == 0 {
					flag = true
					break
				}
				gameInfo = append(gameInfo, util.GetGameName(&_price), _price.DiscountPercent, fmt.Sprintf("%.2f", float64(_price.Initial)/100), fmt.Sprintf("%.2f", float64(_price.Final)/100))
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
		if len(gameInfo) > 0 {
			if err := global.F.SetSheetRow(code.SHEET1, A+strconv.Itoa(idx), &gameInfo); err != nil {
				global.Logger.Error("写入Excel失败", code.ERROR, err)
			}
			idx += 1
		}
	}

	if err := global.F.SaveAs(fmt.Sprintf("steam_price_%s.xlsx", time.Now().Format(time.DateOnly))); err != nil {
		global.Logger.Error("保存Excel失败", code.ERROR, err)
		os.Exit(1)
	}
}

func GeneratePriceTxt() string {
	fileName := fmt.Sprintf("steam_price_%s.txt", time.Now().Format(time.DateOnly))
	_file, err := os.Create(fileName)
	if err != nil {
		global.Logger.Error("创建价格文件失败", code.ERROR, err)
		os.Exit(1)
	}
	//拿到所有游戏的价格
	for _, gameID := range global.GameList {
		data := strings.Builder{}
		flag := false
		gameInfo := make([]any, 0, 84)
		cnStr := strings.Builder{}
		for i := range global.LocationList {

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
						gameName := util.GetGameName(&models.SteamGamePrice{SteamGameID: uint(_storeData.Data.SteamAppid)})
						if gameName == "" {
							gameName = _storeData.Data.Name
						}
						gameInfo = append(gameInfo, gameName, " ", " ", "免费")
						flag = true
						data.WriteString(fmt.Sprintf("游戏名:%s 国区现价:免费\n", gameName))
						//if err = global.F.SetSheetRow(code.SHEET1, A+strconv.Itoa(idx), &gameInfo); err != nil {
						//	global.Logger.Error("写入Excel失败", code.ERROR, err)
						//}
						//idx += 1
						break
					}
				}
			}

			_price := models.SteamGamePrice{}
			if err = global.DB.Preload("SteamGame").Where("steam_game_id  = ?", gameID).Where("steam_location_id = ?", i).First(&_price).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					global.Logger.Error("查询steam游戏价格失败", code.ERROR, err, "游戏ID", gameID, "区ID", i)
				}
				continue
			}
			if i == 1 {
				if _price.Initial == 0 && _price.Final == 0 {
					flag = true
					break
				}
				cnStr.WriteString(fmt.Sprintf("游戏名:%s 折扣率:%d 国区原价:%s 国区现价:%s ", util.GetGameName(&_price), _price.DiscountPercent, fmt.Sprintf("%.2f", float64(_price.Initial)/100), fmt.Sprintf("%.2f", float64(_price.Final)/100)))
				_c, err := global.R.HGet(global.CTX, "SteamGamePriceCheapest", gameID).Result()
				if err != nil {
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
				cnStr.WriteString(fmt.Sprintf("史低价格:%.2f 史低时间:%s ", _p, _date))
				//gameInfo = append(gameInfo, _p, _date)
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
			data.WriteString(fmt.Sprintf("%s区原价:%s %s区现价:%s ", _location.Name, initStr, _location.Name, finalStr))
			//gameInfo = append(gameInfo, initStr, finalStr)
		}
		if flag {
			continue
		}
		if data.String() != "" {
			data.WriteString("\n")
		}
		if _, err = _file.WriteString(cnStr.String()); err != nil {
			global.Logger.Error("写入文件失败", code.ERROR, err)
		}
		if _, err = _file.WriteString(data.String()); err != nil {
			global.Logger.Error("写入文件失败", code.ERROR, err)
		}
	}
	if err = _file.Sync(); err != nil {
		global.Logger.Error("同步文件失败", code.ERROR, err)
		os.Exit(1)
	}
	if err = _file.Close(); err != nil {
		global.Logger.Error("关闭文件失败", code.ERROR, err)
		os.Exit(1)
	}
	return fileName
	//if err := global.F.SaveAs(fmt.Sprintf("steam_price_%s.xlsx", time.Now().Format(time.DateOnly))); err != nil {
	//	global.Logger.Error("保存Excel失败", code.ERROR, err)
	//	os.Exit(1)
	//}
}

func IndexPrice() {
	fileName := GeneratePriceTxt()
	var (
		fileID       string
		err          error
		describeData *bailian20231229.DescribeFileResponseBodyData
	)
	if fileID, err = UploadFileToALiYunBaiLian(fileName); err != nil {
		return
	}
	time.Sleep(time.Minute * 2)
	flag := true
	for _ = range 20 {
		if describeData, err = DescribeBaiLianFile(fileID); err != nil {
			global.Logger.Error("获取文件信息失败", code.ERROR, err)
			return
		}
		if *describeData.Status == "PARSE_SUCCESS" {
			flag = false
			break
		}
		time.Sleep(time.Minute * 1)
	}
	if flag {
		global.Logger.Error("价格文件解析失败,超过20分钟")
		return
	}
	if err = SubmitIndexAddDocumentsJob([]*string{&fileID}); err != nil {
		return
	}
	docs := ListIndexDocuments()
	deleteList := make([]*string, 0)
	for _, j := range docs {
		if strings.HasPrefix(*j.Name, "steam_price") && *j.Name != fileName {
			deleteList = append(deleteList, j.Id)
		}
	}
	if err = DeleteIndexDocument(deleteList); err != nil {
		return
	}

}
