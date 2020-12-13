package elasticx

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"io"
)

type Model interface {
	IndexName(ctx context.Context) string
	TypeName() string
}

func NewBaseModel(clientName, index, table string) *Base {
	client := _clientsMgr.Get(clientName).Client()
	return &Base{
		clientName: clientName,
		index:      index,
		table:      table,
		client:     client,
	}
}

// ES7
func NewBaseModelV7(clientName, index string) *Base {
	client := _clientsMgr.Get(clientName).Client()
	return &Base{
		clientName: clientName,
		index:      index,
		client:     client,
	}
}

type Base struct {
	client *elastic.Client

	clientName string
	index      string
	table      string
}

func (b *Base) Client() *elastic.Client {
	return _clientsMgr.Get(b.clientName).Client()
}

// IndexName 返回Index名字，如果是压测流量，则返回对应的影子索引
func (b *Base) IndexName(ctx context.Context) string {
	return b.index
}

func (b *Base) TableName() string {
	return b.table
}

func (b *Base) CreateService() *elastic.IndexService {
	return b.client.Index().Index(b.index)
}

func (b *Base) UpdateService() *elastic.UpdateService {
	return b.client.Update().Index(b.index)
}

func (b *Base) GetService() *elastic.GetService {
	return b.client.Get().Index(b.index)
}

func (b *Base) SearchService() *elastic.SearchService {
	return b.client.Search(b.index)
}

func (b *Base) ScrollService() *elastic.ScrollService {
	return b.client.Scroll(b.index)
}

func (b *Base) DeleteService() *elastic.DeleteService {
	return b.client.Delete().Index(b.index)
}


// id是es自动产生的id
func (b *Base) Delete(ctx context.Context, id string) error {

	_, err := b.DeleteService().Index(b.IndexName(ctx)).Id(id).Refresh("true").Do(ctx)
	return err
}

type QueryParam struct {
	PageNo   int
	PageSize int
	Sort     string
	RouteId  string
	Asc      bool
	Select   []string
}

func (b *Base) FindWithForeach(ctx context.Context, query elastic.Query, p *QueryParam, fn func(h *elastic.SearchHit)) (int64, error) {
	ss := b.SearchService().Index(b.IndexName(ctx)).Query(query).Pretty(true)

	ss.TrackTotalHits(true)

	if p != nil {
		if p.PageNo > 0 && p.PageSize > 0 {
			ss.Size(p.PageSize).From((p.PageNo - 1) * p.PageSize)
		}
		if p.Sort != "" {
			ss.Sort(p.Sort, p.Asc)
		}
		if len(p.Select) > 0 {
			ss.FetchSourceContext(elastic.NewFetchSourceContext(true).Include(p.Select...))
		}

		if p.RouteId != "" {
			ss.Routing(p.RouteId)
		}
	}

	res, err := ss.Do(ctx)
	if err != nil {
		return 0, err
	}

	for _, v := range res.Hits.Hits {
		fn(v)
	}
	return res.Hits.TotalHits.Value, nil
}

func (b *Base) FindMany(ctx context.Context, query *elastic.BoolQuery, p *QueryParam, results interface{}, aggs elastic.Aggregation) (int64, error) {
	ss := b.SearchService().Index(b.IndexName(ctx)).Query(query).Pretty(true)
	if p != nil && p.PageNo > 0 && p.PageSize > 0 {
		ss = ss.Size(p.PageSize).From((p.PageNo - 1) * p.PageSize)
	}
	if p != nil && p.Sort != "" {
		ss = ss.Sort(p.Sort, p.Asc)
	}
	if aggs != nil {
		ss = ss.Aggregation("metric", aggs)
	}
	res, err := ss.Do(ctx)
	if err != nil {
		return 0, err
	}
	// 走的是普通查询
	if aggs == nil {
		// 如果不是聚合操作
		if len(res.Hits.Hits) > 0 {
			var jsonStr string
			jsonStr = jsonStr + "["
			for key, hit := range res.Hits.Hits {
				jsonStr = jsonStr + string(hit.Source)
				if key != len(res.Hits.Hits)-1 {
					jsonStr = jsonStr + ","
				}
			}
			jsonStr = jsonStr + "]"
			err := json.Unmarshal([]byte(jsonStr), &results)
			if err != nil {
				return 0, err
			}
		}
	} else {
		//走聚合操作返回结果
		type Item struct {
			Key   string `json:"key"`
			Count int64  `json:"count"`
		}
		var items []Item
		term, ok := res.Aggregations.Terms("metric")
		if ok {
			for _, bucket := range term.Buckets {
				items = append(items, Item{
					Key:   bucket.Key.(string),
					Count: bucket.DocCount,
				})
			}
			jsonStr, err := json.Marshal(items)
			if err != nil {
				return 0, err
			}
			err = json.Unmarshal(jsonStr, &results)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, errors.New("Found no reuslt")
		}
	}
	return res.Hits.TotalHits.Value, nil
}

// Find 搜索函数搜索条件需要自己定义 返回数量跟错误
func (b *Base) Find(ctx context.Context, query *elastic.BoolQuery, extra map[string]interface{}, results interface{}, aggs elastic.Aggregation) (int64, error) {
	size := 0
	page := 0
	val, ok := extra["size"]
	if ok {
		size = val.(int)
	}
	val, ok = extra["page"]
	if ok {
		page = val.(int)
	}
	val, ok = extra["sort"]
	sort := ""
	if ok {
		sort = val.(string)
	}
	ss := b.SearchService().Index(b.IndexName(ctx)).Query(query).Pretty(true)
	if size >= 1 && page >= 1 {
		ss = ss.Size(size).From((page - 1) * size)
	}
	if sort != "" {
		ss = ss.Sort(sort, true)
	}
	if aggs != nil {
		ss = ss.Aggregation("metric", aggs)
	}
	res, err := ss.Do(ctx)
	if err != nil {
		return 0, err
	}
	// 走的是普通查询
	if aggs == nil {
		// 如果不是聚合操作
		if len(res.Hits.Hits) > 0 {
			var jsonStr string
			jsonStr = jsonStr + "["
			for key, hit := range res.Hits.Hits {
				jsonStr = jsonStr + string(hit.Source)
				if key != len(res.Hits.Hits)-1 {
					jsonStr = jsonStr + ","
				}
			}
			jsonStr = jsonStr + "]"
			err := json.Unmarshal([]byte(jsonStr), &results)
			if err != nil {
				return 0, err
			}
		}
	} else {
		//走聚合操作返回结果
		type Item struct {
			Key   string `json:"key"`
			Count int64  `json:"count"`
		}
		var items []Item
		term, ok := res.Aggregations.Terms("metric")
		if ok {
			for _, bucket := range term.Buckets {
				items = append(items, Item{
					Key:   bucket.Key.(string),
					Count: bucket.DocCount,
				})
			}
			jsonStr, err := json.Marshal(items)
			if err != nil {
				return 0, err
			}
			err = json.Unmarshal(jsonStr, &results)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, errors.New("Found no reuslt")
		}
	}
	return res.Hits.TotalHits.Value, nil
}

// 支持嵌套聚合，需要提供子聚合名称，搜索函数搜索条件需要自己定义 返回数量跟错误
func (b *Base) FindWithNestedAggregation(ctx context.Context, query *elastic.BoolQuery,
	extra map[string]interface{}, results interface{}, aggs elastic.Aggregation, subName string) (int64, error) {
	size := 0
	page := 0
	val, ok := extra["size"]
	if ok {
		size = val.(int)
	}
	val, ok = extra["page"]
	if ok {
		page = val.(int)
	}
	val, ok = extra["sort"]
	sort := ""
	if ok {
		sort = val.(string)
	}
	ss := b.SearchService().Index(b.IndexName(ctx)).Query(query).Pretty(true)
	if size >= 1 && page >= 1 {
		ss = ss.Size(size).From((page - 1) * size)
	}
	if sort != "" {
		ss = ss.Sort(sort, true)
	}
	if aggs != nil {
		ss = ss.Aggregation("metric", aggs)
	}
	res, err := ss.Do(ctx)
	if err != nil {
		return 0, err
	}

	type Item struct {
		Key   string `json:"key"`
		Count int64  `json:"count"`
	}
	var items []Item
	nestedAggregation, ok := res.Aggregations.Nested("metric")
	if !ok {
		return 0, errors.New("result parsing error")
	}
	termsAggregation, ok := nestedAggregation.Aggregations.Terms(subName)
	if !ok {
		return 0, errors.New("result parsing error")
	}

	buckets := termsAggregation.Buckets
	if buckets != nil {
		for _, bucket := range buckets {
			items = append(items, Item{
				Key:   bucket.Key.(string),
				Count: bucket.DocCount,
			})
		}
		jsonStr, err := json.Marshal(items)
		if err != nil {
			return 0, err
		}
		err = json.Unmarshal(jsonStr, &results)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, errors.New("Found no reuslt")
	}
	return res.Hits.TotalHits.Value, nil
}

// Scroll
func (b *Base) Scroll(ctx context.Context, keepAlive string, source []string, query *elastic.BoolQuery, extra map[string]interface{}, results interface{}) (int64, string, error) {
	// 参数获取
	size := 0
	val, ok := extra["size"]
	if ok {
		size = val.(int)
	}

	// 查询构造
	ss := b.ScrollService().Index(b.IndexName(ctx)).Query(query)
	if keepAlive != "" {
		ss = ss.Scroll(keepAlive)
	}
	if len(source) > 0 {
		ss = ss.FetchSourceContext(elastic.NewFetchSourceContext(true).Include(source...))
	}
	if size >= 1 {
		ss = ss.Size(size)
	}
	ss = ss.Pretty(true)

	res, err := ss.Do(ctx)
	if err == io.EOF {
		return 0, "", nil
	}
	if err != nil {
		return 0, "", err
	}

	if len(res.Hits.Hits) > 0 {
		var jsonStr string
		jsonStr = jsonStr + "["
		for key, hit := range res.Hits.Hits {
			jsonStr = jsonStr + string(hit.Source)
			if key != len(res.Hits.Hits)-1 {
				jsonStr = jsonStr + ","
			}
		}
		jsonStr = jsonStr + "]"
		err := json.Unmarshal([]byte(jsonStr), &results)
		if err != nil {
			return 0, "", err
		}
	}
	return res.Hits.TotalHits.Value, res.ScrollId, nil
}

func (b *Base) CloseScroll(ctx context.Context) error {
	return b.ScrollService().Index(b.IndexName(ctx)).Clear(ctx)
}

func (b *Base) InsertBodyJSON(ctx context.Context, id, routeID string, body interface{}) (err error) {
	ss := b.CreateService().Index(b.IndexName(ctx))
	if id != "" {
		ss = ss.Id(id)
	}
	if routeID != "" {
		ss = ss.Routing(routeID)
	}

	_, err = ss.BodyJson(body).Do(ctx)

	return err
}

// Insert documents是一个json字符串，ES场景下，插入doc后，一般无需立即更新，这里不返回ID
func (b *Base) Insert(ctx context.Context, id, documents string) error {
	ss := b.CreateService().Index(b.IndexName(ctx))
	if id != "" {
		ss = ss.Id(id)
	}
	_, err := ss.BodyString(documents).Do(ctx)
	return err
}

func (b *Base) InsertWithRouting(ctx context.Context, id, routeid, documents string) error {
	_, err := b.CreateService().Index(b.IndexName(ctx)).Id(id).Routing(routeid).BodyString(documents).Do(ctx)
	return err
}

func (b *Base) Update(ctx context.Context, id string, replacement map[string]interface{}) error {
	_, err := b.UpdateService().Index(b.IndexName(ctx)).Id(id).Doc(replacement).Do(ctx)
	return err
}

func (b *Base) UpdateWithRouting(ctx context.Context, id, routeid string, replacement map[string]interface{}) error {
	_, err := b.UpdateService().Index(b.IndexName(ctx)).Id(id).Routing(routeid).Doc(replacement).Do(ctx)
	return err
}

func (b *Base) Upsert(ctx context.Context, id string, doc interface{}) error {
	_, err := b.UpdateService().Index(b.IndexName(ctx)).Id(id).Doc(doc).DocAsUpsert(true).Refresh("true").Do(ctx)
	return err
}

func (b *Base) UpsertWithRouting(ctx context.Context, id, routeid string, doc interface{}) error {
	_, err := b.UpdateService().Index(b.IndexName(ctx)).Id(id).Routing(routeid).Doc(doc).DocAsUpsert(true).Refresh("true").Do(ctx)
	return err
}

func (b *Base) Get(ctx context.Context, id, routeID string, result interface{}) (err error) {
	getSvc := b.GetService().Index(b.IndexName(ctx)).Id(id)
	if routeID != "" {
		getSvc = getSvc.Routing(routeID)
	}

	var re *elastic.GetResult
	re, err = getSvc.Do(ctx)
	if err != nil {
		return err
	}

	err = json.Unmarshal(re.Source, result)

	return err
}

func (b *Base) BatchUpsert(ctx context.Context, values []interface{}, idFunc func(i int) string) (err error) {

	if len(values) == 0 {
		return nil
	}

	bulkRequest := b.Client().Bulk()

	for idx := range values {
		docID := idFunc(idx)
		indexReq := elastic.NewBulkIndexRequest().Index(b.IndexName(ctx)).Id(docID).Doc(values[idx])
		bulkRequest = bulkRequest.Add(indexReq)
	}

	_, err = bulkRequest.Do(ctx)

	return err
}
