package models

import (
    "code.gitea.io/gitea/modules/log"
)

type LoginModel struct {

}

func NewLoginModel() ModelInterface {
    m := &LoginModel{

    }

    return m
}

func (m *LoginModel) SourcePath() string {
    log.Debug("source")
    return "/Users/caiwenpi/Desktop/csv/login"
}

func (m *LoginModel) Parse(p string) error {
    log.Debug("parse")
    return nil
}

func (m *LoginModel) Validate() error {
    log.Debug("validate")
    return nil
}

func (m *LoginModel) Send() error {
    log.Debug("send")
    return nil
}

func (m *LoginModel) Destroy() {
    log.Debug("destroy")
}

func init() {
    Register("login", NewLoginModel)
}
