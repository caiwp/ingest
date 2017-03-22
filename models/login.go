package models

import (
    "encoding/json"
    "path"

    "code.gitea.io/gitea/modules/log"
    "github.com/caiwp/ingest/modules/setting"
    "gopkg.in/asaskevich/govalidator.v4"
    "github.com/gocarina/gocsv"
    "github.com/caiwp/ingest/modules/kafka"
    "bytes"
    "strings"
    "time"
    "os"
    "github.com/Unknwon/com"
    "sync"
)

const LOGIN_CATEGORY = "login"

type LoginModel struct {
    ProductName  string `valid:"required" flume`
    PlatformName string `valid:"required" flume`
    ChannelName  string `valid:"required" flume`
    GameserverNo int32  `valid:"required,customIntValidator" flume`

    Ip           string `valid:"ip" flume`
    Time         string `valid:"required" flume`

    AccountId    string `valid:"required" flume`
    AccountName  string `valid:"-" flume`
    ChrId        string `valid:"required" flume`
    ChrName      string `valid:"required" flume`
    ChrLevel     int32  `valid:"required,customIntValidator" flume`
    ChrLevelVip  int32  `valid:"-,customIntValidator" flume`

    DeviceId     string `valid:"required" flume`
    DeviceType   string `valid:"_" flume`
    Os           string `valid:"_" flume`
    Carrier      string `valid:"_" flume`
    NetworkType  string `valid:"_" flume`
    Resolution   string `valid:"-" flume`

    SessionId    string `valid:"required" flume`

    lock         sync.Mutex
}

func NewLoginModel() ModelInterface {
    return new(LoginModel)
}

func (m *LoginModel) SourcePath() string {
    return path.Join(setting.SourceDataPath, LOGIN_CATEGORY)
}

func (m *LoginModel) Parse(content string) {
    var sli []interface{}
    if err := json.Unmarshal([]byte(content), &sli); err != nil {
        log.Error(4, "json unmarshal %v failed", content)
    }

    initStruct(sli, m)
    log.Debug("%v", m)
}

func (m *LoginModel) Validate() bool {
    log.Debug("validate")

    // 内置的 int 不行，使用自定义的
    govalidator.CustomTypeTagMap.Set("customIntValidator", govalidator.CustomTypeValidator(isInt32))

    res, err := govalidator.ValidateStruct(m)
    if err != nil {
        log.Warn("v", err)
    }
    return res
}

func (m *LoginModel) Send() {
    log.Debug("send")

    var s string
    massage := bytes.NewBufferString(s)

    ms := []*LoginModel{}
    ms = append(ms, m)
    err := gocsv.MarshalWithoutHeaders(&ms, massage)
    if err != nil {
        log.Error(4, "csv marshal failed: %v", ms)
        return
    }
    if err = kafka.SendMassage(strings.TrimSpace(massage.String()), "test"); err != nil {
        log.Error(4, "kafka send massage failed: %v", massage)
    }
}

func (m *LoginModel) Destroy(path string) {
    log.Debug("destroy")
    if com.IsExist(path) {
        err := os.Remove(path)
        if err != nil {
            log.Error(4, "remove file %s failed: %v", path, err)
        }
    }
}

func (m *LoginModel) Backup(s string) {
    log.Debug("backup %s", s)
    m.lock.Lock()
    defer m.lock.Unlock()
    fileName := time.Now().Format("200601021504")
    path := path.Join(setting.BackupDataPath, LOGIN_CATEGORY, fileName)
    file, err := os.OpenFile(path, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0660)
    log.Debug("%v", path)
    if err != nil {
        log.Error(4, "open file %s failed: %v", path, err)
        return
    }
    defer file.Close()

    _, err = file.WriteString(s + "\n")
    if err != nil {
        log.Error(4, "write file %s failed: %v", path, err)
        return
    }
}

func init() {
    Register("login", NewLoginModel)
}

/*
func (m *LoginModel) validateProduct() error {

}
*/
