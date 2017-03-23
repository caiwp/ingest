package models

import (
    "github.com/caiwp/ingest/modules/setting"
    "fmt"
)

type Product struct {
    ProductId int32
    Name      string
}

func GetProduct(name string) (*Product, error) {
    var p = new(Product)
    query := fmt.Sprintf("SELECT product_id, name FROM products WHERE name = '%s' LIMIT 1", name)
    q := setting.Db.NewQuery(query)
    err := q.One(p)
    if err != nil {
        return nil, err
    }
    if p.ProductId == 0 {
        return nil, fmt.Errorf("product %s not found", name)
    }
    return p, nil
}
