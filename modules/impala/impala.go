package impala

import (
    "fmt"
    "github.com/caiwp/ingest/modules/setting"
    "path"
    "strings"
    "code.gitea.io/gitea/modules/log"
    "github.com/koblas/impalathing"
)

var impl *impalathing.Connection

type LoadInterface interface {
    TableName() string
    LoadToDw()
}

type loadType func() LoadInterface

var adapters = make(map[string]*Loader)

type Loader struct {
    adapter string
    model   loadType
}

func newModeler(adapter string, model loadType) *Loader {
    m := &Loader{
        adapter: adapter,
        model: model,
    }
    return m
}

func Register(name string, model loadType) {
    adapters[name] = newModeler(name, model)
}

func Run(adapter string) error {
    m, ok := adapters[adapter]
    if !ok {
        return fmt.Errorf("unknown adapter %s", adapter)
    }

    var err error
    impl = setting.ImplConn

    model := m.model()
    hdfsPath := path.Join(setting.HDFSPath, model.TableName())
    tableName := model.TableName()

    query := fmt.Sprintf("LOAD DATA INPATH '%s' INTO TABLE %s", hdfsPath, tableName)
    _, err = impl.Query(query)
    if err != nil {
        if strings.Contains(err.Error(), "contains no visible files") {
            log.Info("Sourch path %s contains no visible files", hdfsPath)
            return nil
        }
        return err
    }

    model.LoadToDw()

    return nil
}
