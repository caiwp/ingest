package flume

import (
    "bufio"
    "fmt"
    "os"
    "reflect"
    "strings"

    "code.gitea.io/gitea/modules/log"
    "github.com/Unknwon/com"
    "github.com/caiwp/ingest/modules/base"
    "github.com/caiwp/ingest/models"
    "time"
    "path"
    "github.com/caiwp/ingest/modules/setting"
    "gopkg.in/asaskevich/govalidator.v4"
)

type ModelInterface interface {
    Category() string
    Parse(string)
    Validate() bool
    Send()
}

type modelType func() ModelInterface

var adapters = make(map[string]*Modeler)

func Register(name string, model modelType) {
    var err error
    adapters[name], err = newModeler(name, model)
    if err != nil {
        log.Error(4, "New model failed: %v", err)
    }
}

type Modeler struct {
    adapter string
    model   modelType
}

func newModeler(adapter string, model modelType) (*Modeler, error) {
    m := &Modeler{
        adapter: adapter,
        model:   model,
    }

    return m, nil
}

func Run(adapter string) error {
    m, ok := adapters[adapter]
    if !ok {
        return fmt.Errorf("unknown adapter %s", adapter)
    }

    var err error

    model := m.model()
    path := path.Join(setting.SourceDataPath, model.Category())
    if !com.IsDir(path) {
        return fmt.Errorf("path is not a dir: %s", path)
    }

    var files []string
    files, err = base.GetFileListSortByMTime(path)
    if len(files) < 2 {
        log.Info("No file found: %s", path)
        return nil
    }

    files = files[:1]
    for _, f := range files {
        if err = parse(f, model); err != nil {
            log.Error(4, "Parse file %v failed: %v", f, err)
            continue
        }
        destroy(f)
    }

    return nil
}

func parse(path string, model ModelInterface) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    var content string
    sc := bufio.NewScanner(file)
    for sc.Scan() {
        content = sc.Text()
        model.Parse(content)

        if validate(model) {
            model.Send()
        } else {
            backup(content, model)
        }
    }

    if err = sc.Err(); err != nil {
        return err
    }
    return nil
}

func validate(m ModelInterface) bool {
    govalidator.TagMap["product"] = govalidator.Validator(isProduct)
    govalidator.TagMap["platform"] = govalidator.Validator(isPlatform)
    govalidator.TagMap["channel"] = govalidator.Validator(isChannel)
    govalidator.TagMap["datetime"] = govalidator.Validator(isDateTime)

    // 内置的 int 不行，使用自定义的
    govalidator.CustomTypeTagMap.Set("int32Validator", govalidator.CustomTypeValidator(isInt32))
    govalidator.CustomTypeTagMap.Set("validateGameserver", govalidator.CustomTypeValidator(isGameserver))

    return m.Validate()
}

func backup(content string, model ModelInterface) {
    fileName := time.Now().Format("2006010215")
    path := path.Join(setting.BackupDataPath, model.Category(), fileName)
    file, err := os.OpenFile(path, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0660)
    if err != nil {
        log.Error(4, "Open file %s failed: %v", path, err)
        return
    }
    defer file.Close()

    _, err = file.WriteString(content + "\n")
    if err != nil {
        log.Error(4, "Write file %s failed: %v", path, err)
        return
    }
    log.Info("Backup %s success", path)
}

func destroy(path string) {
    if com.IsExist(path) {
        err := os.Remove(path)
        if err != nil {
            log.Error(4, "Remove file %s failed: %v", path, err)
        }
    }
}

func initStruct(sl []interface{}, m ModelInterface) {
    typ := reflect.TypeOf(m).Elem()
    val := reflect.ValueOf(m).Elem()

    // 直接替换掉 , 号
    r := strings.NewReplacer(",", "")
    j := 0
    for i := 0; i < typ.NumField(); i++ {
        f := typ.Field(i)
        if strings.Contains(string(f.Tag), "flume") {
            v := val.FieldByName(f.Name)
            if v.IsValid() && v.CanSet() {
                switch v.Kind() {
                case reflect.Int32:
                    if tmp, ok := sl[j].(float64); ok {
                        v.SetInt(int64(tmp))
                    }
                case reflect.String:
                    if tmp, ok := sl[j].(string); ok {
                        v.SetString(r.Replace(tmp))
                    }
                case reflect.Float32:
                    if tmp, ok := sl[j].(float64); ok {
                        v.SetFloat(tmp)
                    }
                }
                j++
            }
        }
    }
}

func isInt32(i interface{}, o interface{}) bool {
    if _, ok := i.(int32); ok {
        return true
    }
    return false
}

var (
    Product *models.Product
    Platform *models.Platform
    Channel *models.Channel
    Gameserver *models.Gameserver
)

func isProduct(name string) bool {
    var err error
    Product, err = models.GetProduct(name)
    if err != nil {
        log.Warn("Get product %s failed: %v", name, err)
        return false
    }
    return true
}

func isPlatform(name string) bool {
    if Product == nil {
        return false
    }
    var err error
    Platform, err = models.GetPlatform(Product.ProductId, name)
    if err != nil {
        log.Warn("Get platform %s failed: %v", name, err)
        return false
    }
    return true
}

func isChannel(name string) bool {
    if Product == nil {
        return false
    }
    var err error
    Channel, err = models.GetChannel(Product.ProductId, name)
    if err != nil {
        log.Warn("Get channel %s failed: %v", name, err)
        return false
    }
    return true
}

func isGameserver(i interface{}, o interface{}) bool {
    if Product == nil || Platform == nil {
        return false
    }
    no, ok := i.(int32)
    if !ok {
        return false
    }
    var err error
    Gameserver, err = models.GetGameserver(Product.ProductId, Platform.PlatformId, no)
    if err != nil {
        log.Warn("Get gameserver %d failed: %v", no, err)
        return false
    }
    return true
}

func isDateTime(t string) bool {
    layout := "2006-01-02 15:04:05"
    _, err := time.Parse(layout, t)
    if err != nil {
        return false
    }
    return true
}
