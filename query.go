package general_counter

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	elasticsearch "github.com/olivere/elastic/v7"
)

func (gcounter_ *GeneralCounter) QueryTotal(gkey string, gtype string) (*GCounterModel, error) {
	// query sql
	result := &GCounterModel{}
	query := gcounter_.db.Table(TABLE_NAME_G_COUNTER).Where("id = ?", gkey+":"+gtype).Find(result)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}

	return result, nil
}

func (gcounter_ *GeneralCounter) QueryTotalBatch(gkeys []string, gtype string) ([]*GCounterModel, error) {
	// query sql
	result := []*GCounterModel{}
	ids := []string{}
	for _, v := range gkeys {
		ids = append(ids, v+":"+gtype)
	}
	query := gcounter_.db.Table(TABLE_NAME_G_COUNTER).Where("id IN ?", ids).Find(&result)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}

	return result, nil
}

func (gcounter_ *GeneralCounter) QueryAgg(gkey string, gtype string, startDate string, endDate string) ([]*GCounterDailyAggModel, error) {
	// check date
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("start date error")
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("end date error")
	}

	duration := endTime.Sub(startTime)
	if duration < 0 {
		return nil, errors.New("start date is greater than end date")
	}
	totaldays := int(duration.Hours()/24 + 1)
	if totaldays > 365*3 {
		// only do limit for safe, real limit should be in api
		return nil, errors.New("too big date range")
	}

	// query ecs
	generalQ := elasticsearch.NewBoolQuery().
		Must(elasticsearch.NewRangeQuery("date").From(startDate).To(endDate)).
		Must(elasticsearch.NewTermQuery("gkey", gkey)).
		Must(elasticsearch.NewTermQuery("gtype", gtype))

	searchResult, err := gcounter_.ecs.Search().
		Index(gcounter_.gcounter_config.Project_name+"_"+TABLE_NAME_G_COUNTER_DAILY_AGG).
		Query(generalQ).    // specify the query
		Sort("date", true). // sort by "user" field, ascending
		Size(totaldays).    //
		// Pretty(true).                            // pretty print request and response JSON
		Do(context.Background()) // execute
	if err != nil {
		return nil, err
	}

	aggDataMap := map[string]*GCounterDailyAggModel{}
	var cbdLog GCounterDailyAggModel
	for _, item := range searchResult.Each(reflect.TypeOf(cbdLog)) {
		itemlog := item.(GCounterDailyAggModel)
		aggDataMap[itemlog.Date] = &itemlog
	}

	// query db
	db_result := []*GCounterDailyAggModel{}
	query := gcounter_.db.Table(TABLE_NAME_G_COUNTER_DAILY_AGG)
	query.Where("gkey = ?", gkey)
	query.Where("gtype = ?", gtype)
	query.Where("date >= ?", startDate)
	query.Where("date <= ?", endDate)
	query.Where("status = ?", upload_status_to_upload) // just ignore uploading and uploaded data

	db_err := query.Find(&db_result).Error
	if db_err != nil {
		return nil, db_err
	}

	// if same key exist in db and ecs, just use db data
	for _, v := range db_result {
		_, exist := aggDataMap[v.Date]
		if exist {
			aggDataMap[v.Date].Amount = v.Amount
		} else {
			aggDataMap[v.Date] = v
		}
	}

	// filled all the gaps
	startday := startTime
	dayloop := startday

	for k := 0; k < totaldays; k++ {
		dayloop_str := fmt.Sprintf("%d-%02d-%02d", dayloop.UTC().Year(), dayloop.UTC().Month(), dayloop.UTC().Day())
		if aggDataMap[dayloop_str] == nil {
			aggDataMap[dayloop_str] = &GCounterDailyAggModel{
				Sql_id: 0,
				Id:     dayloop_str + ":" + gkey + ":" + gtype,
				Gkey:   gkey,
				Gtype:  gtype,
				Date:   dayloop_str,
				Amount: NewBigInteger(0),
			}
		}
		dayloop = dayloop.Add(24 * time.Hour)
	}

	// map to arry
	result := make([]*GCounterDailyAggModel, len(aggDataMap))
	index := 0
	for _, value := range aggDataMap {
		result[index] = value
		index++
	}

	if len(result) > 0 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Date < result[j].Date
		})
	} else {
		return nil, nil
	}

	return result, nil
}

func (gcounter_ *GeneralCounter) QueryDetail(gkey string, gtype string, startDate string, endDate string) ([]*GCounterDetailModel, error) {
	// check date
	startTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("start date error")
	}

	endTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("end date error")
	}

	duration := endTime.Sub(startTime)
	if duration < 0 {
		return nil, errors.New("start date is greater than end date")
	}
	totaldays := int(duration.Hours()/24 + 1)
	if totaldays > 365*3 {
		// only do limit for safe, real limit should be in api
		return nil, errors.New("too big date range")
	}

	startDatetime := startDate + " 00:00:00"
	endDatetime := endDate + " 23:59:59"

	// query ecs
	generalQ := elasticsearch.NewBoolQuery().
		Must(elasticsearch.NewRangeQuery("datetime").From(startDatetime).To(endDatetime)).
		Must(elasticsearch.NewTermQuery("gkey", gkey)).
		Must(elasticsearch.NewTermQuery("gtype", gtype))

	searchResult, err := gcounter_.ecs.Search().
		Index(gcounter_.gcounter_config.Project_name+"_"+TABLE_NAME_G_COUNTER_DETAIL).
		Query(generalQ).        // specify the query
		Sort("datetime", true). // sort by "user" field, ascending
		// Pretty(true).                            // pretty print request and response JSON
		Size(10000).             // aws opensearch default max size 10000
		Do(context.Background()) // execute
	if err != nil {
		return nil, err
	}

	result_array := []*GCounterDetailModel{}
	var cbdLog GCounterDetailModel
	for _, item := range searchResult.Each(reflect.TypeOf(cbdLog)) {
		itemlog := item.(GCounterDetailModel)
		result_array = append(result_array, &itemlog)
	}

	// query db
	db_result := []*GCounterDetailModel{}
	query := gcounter_.db.Table(TABLE_NAME_G_COUNTER_DETAIL)
	query.Where("gkey = ?", gkey)
	query.Where("gtype = ?", gtype)
	query.Where("datetime >= ?", startDatetime)
	query.Where("datetime <= ?", endDatetime)

	db_err := query.Find(&db_result).Error
	if db_err != nil {
		return nil, db_err
	}

	for _, v := range db_result {
		v.Datetime = v.Datetime[:19]
		result_array = append(result_array, db_result...)
	}

	if len(result_array) > 0 {
		sort.Slice(result_array, func(i, j int) bool {
			return result_array[i].Datetime < result_array[j].Datetime
		})
	} else {
		return nil, nil
	}

	return result_array, nil
}
