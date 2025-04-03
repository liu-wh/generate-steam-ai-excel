package main

import (
	"generate-steam-ai-excel/code"
	"generate-steam-ai-excel/database"
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/service"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"os"
	"strconv"
)

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

	global.OnlineUserFile = excelize.NewFile()
	_ = global.OnlineUserFile.SetSheetRow("Sheet1", "A1", &[]string{"游戏名", "在线人数"})
	service.CreateClient()
}

func releaseResource() {
	db, _ := global.DB.DB()
	if err := db.Close(); err != nil {
		global.Logger.Error("关闭数据库失败", code.ERROR, err)
		os.Exit(1)
	}
	_ = global.OnlineUserFile.Close()
	_ = global.F.Close()
}

func main() {
	defer releaseResource()
	//service.GeneratePriceExcel()
	//service.IndexOnlineUser()
	//service.GeneratePriceTxt()
	//x, _ := service.ListBaiLianFile()
	//for _, j := range x {
	//	fmt.Println(*j.FileId)
	//}
	service.IndexPrice()

	//fmt.Println(service.UploadFileToALiYunBaiLian("steam_price_2025-03-21.xlsx"))
}
