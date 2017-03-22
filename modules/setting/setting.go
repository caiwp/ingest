package setting

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"code.gitea.io/gitea/modules/log"
	"github.com/Shopify/sarama"
	"github.com/Unknwon/com"
	"github.com/koblas/impalathing"
	"gopkg.in/ini.v1"
)

var (
	Cfg *ini.File

	AppPath    string
	CustomPath string
	CustomConf string

	// Log settings
	LogRootPath string
	LogModes    []string
	LogConfigs  []string

	// Data settings
	DataRootPath   string
	SourceDataPath string
    BackupDataPath string
)

func init() {
	log.NewLogger(0, "console", `{"level": 0}`)

	var err error
	if AppPath, err = execPath(); err != nil {
		log.Fatal(4, "failed to get app path: %v.", err)
	}
}

func execPath() (string, error) {
	f, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(f)
}

func WorkDir() string {
	wd := os.Getenv("FLUME_CLIENT_WORK_DIR")
	if len(wd) > 0 {
		return wd
	}
	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath
	}
	return AppPath[:i]
}

func NewContext() {
	var err error
	Cfg = ini.Empty()

	workDir := WorkDir()
	conf := workDir + "/conf/app.ini"
	if com.IsFile(conf) {
		if err = Cfg.Append(conf); err != nil {
			log.Fatal(4, "failed to load conf '%s': %v.", conf, err)
		}
	}

	CustomPath = workDir + "/custom"
	CustomConf = CustomPath + "/conf/app.ini"

	if com.IsFile(CustomConf) {
		if err = Cfg.Append(CustomConf); err != nil {
			log.Fatal(4, "failed to load custom conf '%s': %v.", CustomConf, err)
		}
	} else {
		log.Warn("custom config '%s' not found.")
	}
	Cfg.NameMapper = ini.AllCapsUnderscore
	LogRootPath = Cfg.Section("").Key("ROOT_PATH").MustString(path.Join(workDir, "log"))

	DataRootPath = Cfg.Section("").Key("DATA_PATH").MustString(path.Join(workDir, "data"))
	SourceDataPath = path.Join(DataRootPath, "source")
    BackupDataPath = path.Join(DataRootPath, "backup")
}

func NewServices() {
	newLogService()
	newKafkaService()
	newModelService()
	// newImpalaService()
}

func CloseServices() {
	closeKafkaService()
	// closeImpalaService()
	closeLogService()
}

var logLevels = map[string]string{
	"Trace":    "0",
	"Debug":    "1",
	"Info":     "2",
	"Warn":     "3",
	"Error":    "4",
	"Critical": "5",
}

func newLogService() {
	LogModes = strings.Split(Cfg.Section("log").Key("MODE").MustString("console"), ",")
	LogConfigs = make([]string, len(LogModes))

	useConsole := false
	for i := 0; i < len(LogModes); i++ {
		LogModes[i] = strings.TrimSpace(LogModes[i])
		if LogModes[i] == "console" {
			useConsole = true
		}
	}

	if !useConsole {
		log.DelLogger("console")
	}

	for i, mode := range LogModes {
		sec, err := Cfg.GetSection("log." + mode)
		if err != nil {
			sec, _ = Cfg.NewSection("log." + mode)
		}

		validLevels := []string{"Trace", "Debug", "Info", "Warn", "Error", "Critical"}
		// Log level.
		levelName := Cfg.Section("log."+mode).Key("LEVEL").In(
			Cfg.Section("log").Key("LEVEL").In("Trace", validLevels),
			validLevels)
		level, ok := logLevels[levelName]
		if !ok {
			log.Fatal(4, "Unknown log level: %s", levelName)
		}

		// Generate log configuration.
		switch mode {
		case "console":
			LogConfigs[i] = fmt.Sprintf(`{"level":%s}`, level)
		case "file":
			logPath := sec.Key("FILE_NAME").MustString(path.Join(LogRootPath, "ingest.log"))
			if err = os.MkdirAll(path.Dir(logPath), os.ModePerm); err != nil {
				panic(err.Error())
			}

			LogConfigs[i] = fmt.Sprintf(
				`{"level":%s,"filename":"%s","rotate":%v,"maxlines":%d,"maxsize":%d,"daily":%v,"maxdays":%d}`, level,
				logPath,
				sec.Key("LOG_ROTATE").MustBool(true),
				sec.Key("MAX_LINES").MustInt(1000000),
				1<<uint(sec.Key("MAX_SIZE_SHIFT").MustInt(28)),
				sec.Key("DAILY_ROTATE").MustBool(true),
				sec.Key("MAX_DAYS").MustInt(7))
		}

		log.NewLogger(Cfg.Section("log").Key("BUFFER_LEN").MustInt64(10000), mode, LogConfigs[i])
		log.Info("Log Mode: %s(%s)", strings.Title(mode), levelName)
	}
}

func closeLogService() {
	log.Close()
}

type Impala struct {
	Host     string `ini:"HOST"`
	Port     int    `ini:"PORT"`
	Database string `ini:"DATABASE"`
}

var (
	ImplConn *impalathing.Connection
)

func newImpalaService() {
	impl := new(Impala)
	err := Cfg.Section("impala").MapTo(&impl)
	if err != nil {
		log.Fatal(4, "set impala config failed: %v", err)
	}

	ImplConn, err = impalathing.Connect(impl.Host, impl.Port, impalathing.DefaultOptions)
	if err != nil {
		log.Fatal(4, "connect impala failed [%s:%d]: %v", impl.Host, impl.Port, err)
	}

	_, err = ImplConn.Query(fmt.Sprintf("Use %s", impl.Database))
	if err != nil {
		log.Fatal(4, "select default database failed: %v", err)
	}
}

func closeImpalaService() {
	ImplConn.Close()
}

var (
	Producer sarama.SyncProducer
	Massage  *sarama.ProducerMessage
)

func newKafkaService() {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true

	sec := Cfg.Section("kafka")
	host := sec.Key("HOST").MustString("localhost")
	port := sec.Key("PORT").MustInt(9092)
	topic := sec.Key("TOPIC").MustString("test")

	addrs := []string{fmt.Sprintf("%s:%d", host, port)}

	var err error
	Producer, err = sarama.NewSyncProducer(addrs, conf)
	if err != nil {
		log.Fatal(4, "new producer failed: %v", err)
	}

	Massage = &sarama.ProducerMessage{
		Topic:     topic,
		Partition: int32(-1),
		Key:       sarama.StringEncoder("key"),
	}
}

func closeKafkaService() {
	Producer.Close()
}

var ModelCategories []string

func newModelService() {
	categories := strings.Split(Cfg.Section("model").Key("CATEGORY").MustString("login"), ",")

	for _, v := range categories {
		ModelCategories = append(ModelCategories, strings.TrimSpace(v))
	}
}
