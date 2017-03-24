package impala

import (
    "github.com/caiwp/ingest/modules/setting"
    "code.gitea.io/gitea/modules/log"
)

type LoginLoad struct {
}

func NewLoginLoad() LoadInterface {
    return new(LoginLoad)
}

func (m *LoginLoad) TableName() string {
    return "tmp_raw_logins"
}

func (m *LoginLoad) LoadToDw() {
    loginQuery := `
        INSERT INTO raw_logins PARTITION (pid, year, month)
            SELECT
                *,
                CAST(t.product_id AS TINYINT),
                CAST(FROM_UNIXTIME(t.time, 'yy') AS TINYINT),
                CAST(FROM_UNIXTIME(t.time, 'MM') AS TINYINT)
            FROM tmp_raw_logins AS t
    `

    if _, err := impl.Query(loginQuery); err != nil {
        log.Error(4, "Impala query %s failed: %v", loginQuery, err)
        return
    }
    return
}

func init() {
    Register(setting.CATEGORY_LOGIN, NewLoginLoad)
}