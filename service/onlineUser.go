package service

import (
	"encoding/json"
	"fmt"
	"generate-steam-ai-excel/code"
	"generate-steam-ai-excel/global"
	"generate-steam-ai-excel/util"
	bailian20231229 "github.com/alibabacloud-go/bailian-20231229/client"
	"os"
	"strconv"
	"strings"
	"time"
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

func GenerateOnlineUserExcel() string {
	var (
		onlineUserStr string
		err           error
	)
	if onlineUserStr, err = global.R.Get(global.CTX, "SteamGameTop").Result(); err != nil {
		return ""
	}
	gameList := make([]*SteamGame, 0, 3000)
	if err = json.Unmarshal(util.Str2bytes(onlineUserStr), &gameList); err != nil {
		global.Logger.Error("解析在线用户数据失败", code.ERROR, err)
		return ""
	}
	i := 2
	for _, game := range gameList {
		_ = global.OnlineUserFile.SetSheetRow(code.SHEET1, "A"+strconv.Itoa(i), &[]any{
			game.Name,
			game.Count,
		})
		i += 1
	}
	fileName := fmt.Sprintf("steam_online_user_%s.xlsx", time.Now().Format(time.DateOnly))
	if err = global.OnlineUserFile.SaveAs(fileName); err != nil {
		global.Logger.Error("保存Excel失败", code.ERROR, err)
		os.Exit(1)
	}

	return fileName
}

func GenerateOnlineUserTxt() string {
	var (
		onlineUserStr string
		err           error
	)
	if onlineUserStr, err = global.R.Get(global.CTX, "SteamGameTop").Result(); err != nil {
		return ""
	}
	gameList := make([]*SteamGame, 0, 3000)
	if err = json.Unmarshal(util.Str2bytes(onlineUserStr), &gameList); err != nil {
		global.Logger.Error("解析在线用户数据失败", code.ERROR, err)
		return ""
	}
	fileName := fmt.Sprintf("steam_online_user_%s.txt", time.Now().Format(time.DateTime))
	_file, err := os.Create(fileName)
	if err != nil {
		global.Logger.Error("创建价格文件失败", code.ERROR, err)
		os.Exit(1)
	}
	for _, game := range gameList {
		_, _ = _file.WriteString(fmt.Sprintf("游戏名:%s 在线人数:%d\n", game.Name, game.Count))
	}
	_ = _file.Sync()
	_ = _file.Close()

	return fileName
}

func IndexOnlineUser() {
	fileName := GenerateOnlineUserTxt()
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
		if strings.HasPrefix(*j.Name, "steam_online_user_") && *j.Name != fileName {
			deleteList = append(deleteList, j.Id)
		}
	}
	if err = DeleteIndexDocument(deleteList); err != nil {
		return
	}

}
