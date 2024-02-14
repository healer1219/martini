package test

import (
	"github.com/healer1219/martini/bootstrap"
	"github.com/healer1219/martini/global"
	"github.com/healer1219/martini/mevent"
	"testing"
	"time"
)

type SystemInfo struct {
	ID          string    `gorm:"column:S_ID;"`
	Key         string    `gorm:"column:S_KEY"`
	Value       string    `gorm:"column:S_VALUE"`
	Description string    `gorm:"column:S_DESCRIPTION"`
	Params      string    `gorm:"column:S_PARAMS"`
	CreateTime  time.Time `gorm:"column:DT_CREATE_TIME"`
	UpdateTime  time.Time `gorm:"column:DT_UPDATE_TIME"`
	CreateUser  string    `gorm:"column:S_CREATE_USER"`
	UpdateUser  string    `gorm:"column:S_UPDATE_USER"`
	PlatFormFlg string    `gorm:"column:S_PLATFORM_FLG"`
	delFlg      int       `gorm:"column:I_DEL_FLG"`
}

func TestDb(t *testing.T) {
	bootstrap.Default().BootUp()
	var systemInfo SystemInfo
	defer bootstrap.RealeaseDB()
	global.DB().Raw("select * from TAB_SYSTEM_INFO where S_ID = ?", "36a4800b-c751-4bcf-bb94-a71de55fe3d2").Scan(&systemInfo)
	if systemInfo.ID == "" {
		t.Errorf("select failed! %v \n", systemInfo)
	}
}

type InitDbStruct struct {
}

func (i InitDbStruct) OnEvent(ctx *global.Context) {
	bootstrap.InitDb()
}

func init() {
	mevent.Add(bootstrap.StartupEvent, InitDbStruct{})
}
