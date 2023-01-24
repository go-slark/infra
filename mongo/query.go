package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SetUpdateTime(update bson.M) bson.M {
	set, ok := update["$set"]
	if ok {
		set.(bson.M)["update_time"] = time.Now()
	} else {
		update["$set"] = bson.M{"update_time": time.Now()}
	}

	return update
}

type QueryOptions struct {
	Skip     int64
	Limit    int64
	Sort     bson.M
	SortD    bson.D
	Selector bson.M
}

/*
   QueryOptions.Selector : bson.M{"_id": false/true} : false不返回该文档字段 / true返回该文档字段
   sortD [{"_id": 1}, {"name": 1}] : 返回的是有序数据，按照_id->name的顺序返回， sort sortD字段只能选择一个
*/

func ApplyQueryOpts(opts ...QueryOpt) *options.FindOptions {
	query := &options.FindOptions{}
	qo := &QueryOptions{}
	for _, opt := range opts {
		opt(qo)
	}
	if qo.Sort != nil {
		query = query.SetSort(qo.Sort)
	}
	if len(qo.SortD) != 0 {
		query = query.SetSort(qo.SortD)
	}
	if qo.Skip != 0 {
		query = query.SetSkip(qo.Skip)
	}
	if qo.Limit != 0 {
		query = query.SetLimit(qo.Limit)
	}

	if qo.Selector != nil {
		query = query.SetProjection(qo.Selector)
	}

	return query
}

type QueryOpt func(*QueryOptions)

func Skip(skip int64) QueryOpt {
	return func(opts *QueryOptions) {
		opts.Skip = skip
	}
}

func Limit(limit int64) QueryOpt {
	return func(opts *QueryOptions) {
		opts.Limit = limit
	}
}

func Sort(sort bson.M) QueryOpt {
	return func(opts *QueryOptions) {
		opts.Sort = sort
	}
}

func SortD(sortD bson.D) QueryOpt {
	return func(opts *QueryOptions) {
		opts.SortD = sortD
	}
}

func Select(selector bson.M) QueryOpt {
	return func(opts *QueryOptions) {
		opts.Selector = selector
	}
}
