package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type FormatTime time.Time

func (t *FormatTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse(time.DateTime, timeStr)
	*t = FormatTime(t1)
	return err
}

func (t FormatTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format(time.DateTime))
	return []byte(formatted), nil
}

func (t FormatTime) Value() (driver.Value, error) {
	// FormartTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(time.DateTime), nil
}

func (t *FormatTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = FormatTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *FormatTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

type User struct {
	gorm.Model
	Name              string      `gorm:"size:255"`
	AvatarUrl         string      `gorm:"size:255"`
	WechatOpenid      string      `gorm:"uniqueIndex;size:255"`
	WechatSessionKey  string      `gorm:"size:255"`
	Likes             []SteamGame `gorm:"many2many:user_likes;references:SteamGameAppID"`
	Wants             []SteamGame `gorm:"many2many:user_wants;references:SteamGameAppID"`
	Played            []SteamGame `gorm:"many2many:user_played;references:SteamGameAppID"`
	Subscribed        []SteamGame `gorm:"many2many:user_subscribed;references:SteamGameAppID"`
	FavoriteInventory []Inventory `gorm:"many2many:favorite_inventory_users;references:ID"`
	SteamID           string      `gorm:"size:50;index"`
}

type Inventory struct {
	ID            int         `gorm:"primaryKey;autoIncrement"`
	Name          string      `gorm:"size:255"`
	Games         []SteamGame `gorm:"many2many:inventory_games;references:SteamGameAppID"`
	User          User        `gorm:"foreignKey:UserID;references:ID;"`
	UserID        uint        `gorm:"index"`
	LikeCount     int
	FavoriteUsers []User `gorm:"many2many:favorite_inventory_users;references:ID"`
	ViewCount     int
	CreateTime    FormatTime `gorm:"type:datetime"`
	CoverGame     SteamGame  `gorm:"foreignKey:CoverGameID;references:SteamGameAppID;"`
	CoverGameID   uint
	Describe      string `gorm:"size:512"`
	Status        bool
}

type SteamGame struct {
	ID             uint        `gorm:"primaryKey;autoIncrement"`
	EName          string      `gorm:"comment:steam游戏名;size:255"`
	CName          string      `gorm:"comment:steam游戏中文名;size:255"`
	SteamGameAppID uint        `gorm:"uniqueIndex"`
	Type           string      `gorm:"size:15"`
	Inventory      []Inventory `gorm:"many2many:game_inventory"`
}

type SteamLocation struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Code         string `gorm:"size:10;index"`
	Name         string `gorm:"size:100"`
	Currency     string `gorm:"size:100"`
	CurrencyCode string `gorm:"size:100"`
}
type SteamGamePrice struct {
	ID              uint      `gorm:"primaryKey;autoIncrement"`
	SteamGame       SteamGame `gorm:"foreignKey:SteamGameID;references:SteamGameAppID;"`
	SteamGameID     uint
	IsFree          bool
	Initial         int
	Final           int
	DiscountPercent int
	SteamLocation   SteamLocation `gorm:"foreignKey:SteamLocationID;references:ID;"`
	SteamLocationID int
	FinalFormatted  string `gorm:"size:80"`
}

type GameTag struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

type SteamGameStoreData struct {
	Success bool `json:"success"`
	Data    struct {
		Type                string   `json:"type"`
		Name                string   `json:"name"`
		SteamAppid          int      `json:"steam_appid"`
		RequiredAge         any      `json:"required_age"`
		IsFree              bool     `json:"is_free"`
		ControllerSupport   string   `json:"controller_support"`
		Dlc                 []int    `json:"dlc"`
		DetailedDescription string   `json:"detailed_description"`
		AboutTheGame        string   `json:"about_the_game"`
		ShortDescription    string   `json:"short_description"`
		SupportedLanguages  string   `json:"supported_languages"`
		HeaderImage         string   `json:"header_image"`
		Website             string   `json:"website"`
		PcRequirements      any      `json:"pc_requirements"`
		MacRequirements     any      `json:"mac_requirements"`
		LinuxRequirements   any      `json:"linux_requirements"`
		LegalNotice         string   `json:"legal_notice"`
		Developers          []string `json:"developers"`
		Publishers          []string `json:"publishers"`
		PriceOverview       struct {
			Currency         string `json:"currency"`
			Initial          int    `json:"initial"`
			Final            int    `json:"final"`
			DiscountPercent  int    `json:"discount_percent"`
			InitialFormatted string `json:"initial_formatted"`
			FinalFormatted   string `json:"final_formatted"`
		} `json:"price_overview"`
		Packages      []int `json:"packages"`
		PackageGroups any   `json:"package_groups"`
		Platforms     struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Metacritic struct {
			Score int    `json:"score"`
			Url   string `json:"url"`
		} `json:"metacritic"`
		Categories []struct {
			Id          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres      []GameTag `json:"genres"`
		Screenshots []struct {
			Id            int    `json:"id"`
			PathThumbnail string `json:"path_thumbnail"`
			PathFull      string `json:"path_full"`
		} `json:"screenshots"`
		Movies []struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Webm      struct {
				Field1 string `json:"480"`
				Max    string `json:"max"`
			} `json:"webm"`
			Mp4 struct {
				Field1 string `json:"480"`
				Max    string `json:"max"`
			} `json:"mp4"`
			Highlight bool `json:"highlight"`
		} `json:"movies"`
		Recommendations struct {
			Total int `json:"total"`
		} `json:"recommendations"`
		Achievements struct {
			Total       int `json:"total"`
			Highlighted []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"highlighted"`
		} `json:"achievements"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		SupportInfo struct {
			Url   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Background         string `json:"background"`
		BackgroundRaw      string `json:"background_raw"`
		ContentDescriptors struct {
			Ids   []interface{} `json:"ids"`
			Notes interface{}   `json:"notes"`
		} `json:"content_descriptors"`
	} `json:"data"`
}
