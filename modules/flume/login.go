package flume

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
    "fmt"
    "github.com/caiwp/ingest/modules/base"
)

const LOGIN_CATEGORY = "login"

type LoginModel struct {
    ProductId    int32
    ProductName  string `valid:"required,product" flume`
    PlatformId   int32
    PlatformName string `valid:"required,platform" flume`
    ChannelId    int32
    ChannelName  string `valid:"required,channel" flume`
    GameserverId int32
    GameserverNo int32  `valid:"required,validateGameserver" flume`

    Ip           string `valid:"ip" flume`
    Region       string
    TheDate      string
    DateTime     string `valid:"required,datetime" flume`
    Hour         int
    Time         int64

    AccountId    string `valid:"required" flume`
    AccountName  string `valid:"-" flume`
    UniqueChrId  string
    ChrId        string `valid:"required" flume`
    ChrName      string `valid:"required" flume`
    ChrLevel     int32  `valid:"required,int32Validator" flume`
    ChrLevelVip  int32  `valid:"-,int32Validator" flume`

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

    govalidator.TagMap["product"] = govalidator.Validator(isProduct)
    govalidator.TagMap["platform"] = govalidator.Validator(isPlatform)
    govalidator.TagMap["channel"] = govalidator.Validator(isChannel)
    govalidator.TagMap["datetime"] = govalidator.Validator(isDateTime)

    // 内置的 int 不行，使用自定义的
    govalidator.CustomTypeTagMap.Set("int32Validator", govalidator.CustomTypeValidator(isInt32))
    govalidator.CustomTypeTagMap.Set("validateGameserver", govalidator.CustomTypeValidator(isGameserver))

    res, err := govalidator.ValidateStruct(m)
    if err != nil || res == false {
        log.Warn("validate error: %v", err)
        return false
    }

    m.setData()

    return true
}

func (m *LoginModel) setData() {
    m.ProductId = Product.ProductId
    m.PlatformId = Platform.PlatformId
    m.ChannelId = Channel.ChannelId
    m.GameserverId = Gameserver.GameserverId

    if m.Ip != "" {
        path := path.Join(setting.RootPath, "conf", "GeoLite2-City.mmdb")
        rg, err := base.GetRegion(m.Ip, path)
        if err != nil {
            log.Warn("get %s region failed: %v", m.Ip, err)
        }
        m.Region = rg
    }

    datetime, _ := time.Parse("2006-01-02 15:04:05 -07:00", fmt.Sprintf("%s %s", m.DateTime, "+08:00"))
    m.TheDate = datetime.Format("2006-01-02")
    m.Hour = datetime.Hour()
    m.Time = datetime.Unix()

    md5str := base.EncodeMD5(fmt.Sprintf("%d_%d_%d_%s", m.ProductId, m.PlatformId, m.GameserverId, m.ChrId))
    m.UniqueChrId = md5str[8:24]
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
