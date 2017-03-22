package models

import (
    "bufio"
    "fmt"
    "os"
    "reflect"
    "strings"

    "code.gitea.io/gitea/modules/log"
    "github.com/Unknwon/com"
    "github.com/caiwp/ingest/modules/base"
)

type ModelInterface interface {
    SourcePath() string
    Parse(string)
    Validate() bool
    Send()
    Destroy(string)
    Backup(string)
}

type modelType func() ModelInterface

var adapters = make(map[string]*Modeler)

func Register(name string, model modelType) {
    var err error
    adapters[name], err = newModeler(name, model)
    if err != nil {
        log.Error(4, "new model failed: %v", err)
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
    path := model.SourcePath()
    if !com.IsDir(path) {
        return fmt.Errorf("path is not a dir: %s", path)
    }

    var files []string
    files, err = base.GetFileListSortByMTime(path)
    if len(files) < 2 {
        log.Info("no file found: %s", path)
        return nil
    }

    files = files[:1]
    for _, f := range files {
        if err = parse(f, model); err != nil {
            log.Error(4, "parse file %v failed: %v", f, err)
            continue
        }
        model.Destroy(f)
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

        if model.Validate() {
            model.Send()
        } else {
            model.Backup(content)
        }
    }

    if err = sc.Err(); err != nil {
        return err
    }
    return nil
}

func initStruct(sl []interface{}, m ModelInterface) {
    typ := reflect.TypeOf(m).Elem()
    val := reflect.ValueOf(m).Elem()

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
                        v.SetString(tmp)
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