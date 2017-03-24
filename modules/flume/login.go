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
    "fmt"
    "github.com/caiwp/ingest/modules/base"
)

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
}

func NewLoginModel() ModelInterface {
    return new(LoginModel)
}

func (m *LoginModel) Parse(content string) {
    var sli []interface{}
    if err := json.Unmarshal([]byte(content), &sli); err != nil {
        log.Error(4, "Json unmarshal %v failed", content)
    }

    initStruct(sli, m)
    log.Info("Parse login success")
}

func (m *LoginModel) Validate() bool {
    res, err := govalidator.ValidateStruct(m)
    if err != nil || res == false {
        log.Warn("Validate %v error: %v", m, err)
        return false
    }

    log.Info("Validate login success")
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
            log.Warn("Get %s region failed: %v", m.Ip, err)
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
    var s string
    massage := bytes.NewBufferString(s)

    ms := []*LoginModel{}
    ms = append(ms, m)
    err := gocsv.MarshalWithoutHeaders(&ms, massage)
    if err != nil {
        log.Error(4, "Csv marshal failed: %v", ms)
        return
    }
    // FIXME test need to change to its category
    if err = kafka.SendMassage(strings.TrimSpace(massage.String()), "tmp_raw_logins"); err != nil {
        log.Error(4, "Kafka send massage failed: %v", massage)
    }
    log.Info("Send login success")
}

func (m *LoginModel) Category() string {
    return setting.CATEGORY_LOGIN
}

func init() {
    Register(setting.CATEGORY_LOGIN, NewLoginModel)
}
