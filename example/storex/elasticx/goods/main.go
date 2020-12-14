package main

import (
	"context"
	"github.com/ivalue2333/pkg/src/storex/elasticx"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	ctx        = context.Background()
	clientName = "myclient"
	model      *elasticx.Base
)

func init() {
	err := elasticx.ClientsMgr().NewClient(ctx, clientName, []string{"http://127.0.0.1:9200"}...)
	if err != nil {
		panic(err)
	}
	model = elasticx.NewBaseModelV7(clientName, "myindex")
}

type Spu struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	CreateTime    time.Time              `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdateTime    time.Time              `bson:"update_time,omitempty" json:"update_time,omitempty"`
	Platform      string                 `bson:"platform,omitempty" json:"platform,omitempty"`
	PlatUserID    string                 `bson:"plat_user_id,omitempty" json:"plat_user_id,omitempty"`
	PlatGoodsID   string                 `bson:"plat_goods_id,omitempty" json:"plat_goods_id,omitempty"`
	PlatGoodsName string                 `bson:"plat_goods_name,omitempty" json:"plat_goods_name,omitempty"`
	PlatGoodsURL  string                 `bson:"plat_goods_url,omitempty" json:"plat_goods_url,omitempty"`
	PlatGoodsImg  string                 `bson:"plat_goods_img,omitempty" json:"plat_goods_img,omitempty"`
	PlatCid       string                 `bson:"plat_cid,omitempty" json:"plat_cid,omitempty"`
	IsOnsale      bool                   `bson:"is_onsale,omitempty" json:"is_onsale,omitempty"`
	Extra         map[string]interface{} `bson:"extra,omitempty" json:"extra,omitempty"`
	TopProps      [][]string             `bson:"-" json:"top_props,omitempty"`
	PropsName     string                 `bson:"-" json:"props_name,omitempty"`
	SellerCIDs    []string               `bson:"seller_cids,omitempty" json:"seller_cids,omitempty"`
	Cid           int                    `bson:"cid,omitempty" json:"cid,omitempty"`
	ArticleNumber string                 `bson:"article_number,omitempty" json:"article_number,omitempty"` //货号/款号/型号,从props提取的属性值
	ListTime      string                 `bson:"list_time,omitempty" json:"list_time,omitempty"`
}
