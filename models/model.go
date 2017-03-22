package models

import (
    "code.gitea.io/gitea/modules/log"
    "github.com/Unknwon/com"
)

type ModelInterface interface {
    SourcePath() string
    Parse(string) error
    Validate() error
    Send() error
    Destroy()
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
        model: model,
    }

    return m, nil
}

func Run(adapter string) {
    m, ok := adapters[adapter]
    if !ok {
        log.Error(4, "unknown adapter %s", adapter)
        return
    }

    var err error

    model := m.model()
    path := model.SourcePath()
    if !com.IsDir(path) {
        log.Error(4, "path is not a dir: %s", path)
        return
    }
    log.Debug("%v", path)

    if err = model.Parse(path); err != nil {
        log.Error(4, "parse %s failed: %v", adapter, err)
        return
    }

    if err = model.Validate(); err != nil {
        log.Error(4, "validate %s failed: %v", adapter, err)
        return
    }

    if err = model.Send(); err != nil {
        log.Error(4, "end %s failed: %v", adapter, err)
        return
    }

    model.Destroy()
}
