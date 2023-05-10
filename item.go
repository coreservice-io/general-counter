package general_counter

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GcOpTotalConfig struct {
	Enable        bool
	AllowNegative bool
}

type GcOpAggConfig struct {
	Enable  bool
	AggDate *GcAggDate
}
type GcAggDate struct {
	Year  int
	Month int
	Day   int
}

type GcOpDetailConfig struct {
	Enable bool
	Msg    string
}

type GcOp struct {
	Gkey          string      // NOT NULL
	Gtype         string      // NOT NULL
	Amount        *BigInteger // NOT NULL
	Total_config  *GcOpTotalConfig
	Agg_config    *GcOpAggConfig
	Detail_config *GcOpDetailConfig
}

func (gcop *GcOp) run(tx *gorm.DB, aggExpireDays int64) error {

	time_now := time.Now()

	if (gcop.Detail_config == nil || !gcop.Detail_config.Enable) &&
		(gcop.Agg_config == nil || !gcop.Agg_config.Enable) &&
		(gcop.Total_config == nil || !gcop.Total_config.Enable) {
		return errors.New("atleast one config required")
	}

	// detail
	if gcop.Detail_config != nil && gcop.Detail_config.Enable {

		to_create := &GCounterDetailModel{
			Gkey:     gcop.Gkey,
			Gtype:    gcop.Gtype,
			Datetime: time_now.UTC().Format("2006-01-02 15:04:05"),
			Amount:   gcop.Amount,
			Msg:      gcop.Detail_config.Msg,
		}

		tx_result := tx.Table(TABLE_NAME_G_COUNTER_DETAIL).Create(to_create)
		if tx_result.Error != nil {
			return tx_result.Error
		}

		if tx_result.RowsAffected == 0 {
			return errors.New("gcop detail created none row affected")
		}

		tx_update_result := tx.Table(TABLE_NAME_G_COUNTER_DETAIL).Where("sql_id=?", to_create.Sql_id).Update("id", strconv.FormatInt(to_create.Sql_id, 10))
		if tx_update_result.Error != nil {
			return tx_update_result.Error
		}
		if tx_update_result.RowsAffected == 0 {
			return errors.New("gcop detail update sql_id=>id none row affected")
		}
	}

	// agg
	if gcop.Agg_config != nil && gcop.Agg_config.Enable {
		date := ""
		if gcop.Agg_config.AggDate != nil {
			insertTime := time.Date(gcop.Agg_config.AggDate.Year, time.Month(gcop.Agg_config.AggDate.Month), gcop.Agg_config.AggDate.Day, 0, 0, 0, 0, time.UTC)

			// check insertTime is expire or not
			t := time.Now().UTC()
			today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			// insert expire data not allowed
			if insertTime.Unix()-today.Unix() < -aggExpireDays*24*3600 {
				return errors.New("gcop agg insert an expire data")
			}

			date = insertTime.Format("2006-01-02")
		} else {
			date = time_now.UTC().Format("2006-01-02")
		}

		agg_id := date + ":" + gcop.Gkey + ":" + gcop.Gtype

		to_create_agg := &GCounterDailyAggModel{
			Id:     agg_id,
			Gkey:   gcop.Gkey,
			Gtype:  gcop.Gtype,
			Date:   date,
			Amount: gcop.Amount,
			Status: upload_status_to_upload,
		}

		create_result := tx.Table(TABLE_NAME_G_COUNTER_DAILY_AGG).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("amount + ?", gcop.Amount), "status": upload_status_to_upload}),
		}).Create(to_create_agg)

		if create_result.Error != nil {
			return errors.New("agg update/create error,err:" + create_result.Error.Error())
		}

		if create_result.RowsAffected == 0 {
			return errors.New("gcop agg created none row affected")
		}

	}

	// total
	if gcop.Total_config != nil && gcop.Total_config.Enable {

		id := gcop.Gkey + ":" + gcop.Gtype
		to_create_agg := &GCounterModel{
			Id:     id,
			Gkey:   gcop.Gkey,
			Gtype:  gcop.Gtype,
			Amount: gcop.Amount,
		}

		var create_result *gorm.DB
		if gcop.Amount.Sign() >= 0 || gcop.Total_config.AllowNegative {
			create_result = tx.Table(TABLE_NAME_G_COUNTER).Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{"amount": gorm.Expr("amount + ?", gcop.Amount)}),
			}).Create(to_create_agg)

			if create_result.Error != nil {
				return errors.New("gcop total counter create error ,err:" + create_result.Error.Error())
			}

			if create_result.RowsAffected == 0 {
				return errors.New("gcop total counter create error ,err:RowsAffected==0 ")
			}

		} else {
			update_result := tx.Table(TABLE_NAME_G_COUNTER).Where("id = ? ", id).
				Where("amount >= ?", gcop.Amount.Abs()).Update("amount", gorm.Expr("amount + ?", gcop.Amount))

			if update_result.Error != nil {
				return errors.New("gcop total counter update error ,err:" + update_result.Error.Error())
			}

			if update_result.RowsAffected == 0 {
				return errors.New("gcop total counter update error ,err:RowsAffected==0 ")
			}
		}

	}

	return nil
}

type GcFunc struct {
	Func func(tx *gorm.DB) error
}

// /

type GcTx struct {
	gcounter  *GeneralCounter
	item_list []interface{}
}

func (gctx *GcTx) AppendOp(gcop *GcOp) *GcTx {
	if gcop == nil {
		return gctx
	}

	gctx.item_list = append(gctx.item_list, gcop)
	return gctx
}

func (gctx *GcTx) AppendOpEx(gcop *GcOp) error {
	if err := AppendOpCheck(gcop); err != nil {
		return err
	}

	gctx.item_list = append(gctx.item_list, gcop)
	return nil
}

func AppendOpCheck(gcop *GcOp) error {
	if gcop == nil {
		return errors.New("gcop is nil")
	}

	if gcop.Gkey == "" || gcop.Gtype == "" || gcop.Amount == nil {
		return errors.New(fmt.Sprintf("Gkey is empty %v", gcop))
	}

	if gcop.Gkey == "" || gcop.Gtype == "" || gcop.Amount == nil {
		return errors.New(fmt.Sprintf("Gtype is empty %v", gcop))
	}

	if gcop.Gkey == "" || gcop.Gtype == "" || gcop.Amount == nil {
		return errors.New(fmt.Sprintf("Amount is empty %v", gcop))
	}

	if gcop.Amount.Sign() == 0 {
		return errors.New(fmt.Sprintf("Amount is zero %v", gcop))
	}
	return nil
}

func (gctx *GcTx) AppendFunc(txfunc func(tx *gorm.DB) error) *GcTx {
	gctx.item_list = append(gctx.item_list, &GcFunc{
		Func: txfunc,
	})
	return gctx
}

func (gcounter_ *GeneralCounter) CreateTx() *GcTx {
	return &GcTx{
		item_list: make([]interface{}, 0),
		gcounter:  gcounter_,
	}
}

func (gctx *GcTx) Commit() error {
	if len(gctx.item_list) == 0 {
		return nil
	}

	all_empty := true
	// pre-check
	for _, item := range gctx.item_list {
		if !all_empty {
			break
		}
		switch v := item.(type) {
		case *GcFunc:
			all_empty = false
		case *GcOp:
			if v.Amount.Sign() != 0 {
				all_empty = false
			}
		default:
			return errors.New("item type error inside itemlist")
		}
	}

	if all_empty {
		return nil
	}

	for _, item := range gctx.item_list {
		if gcop, ok := item.(*GcOp); ok {
			if err := AppendOpCheck(gcop); err != nil {
				return err
			}
		}
	}

	// use global lock to avoid db dead lock
	gctx.gcounter.commit_lock.Lock()
	defer gctx.gcounter.commit_lock.Unlock()

	//
	tx_err := gctx.gcounter.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range gctx.item_list {
			switch v := item.(type) {
			case *GcFunc:
				func_err := v.Func(tx)
				if func_err != nil {
					return func_err
				}
			case *GcOp:
				op_err := v.run(tx, int64(gctx.gcounter.gcounter_config.Agg_record_expire_days))
				if op_err != nil {
					return op_err
				}
			default:
				return errors.New("item type error inside itemlist")
			}
		}
		return nil
	})
	return tx_err
}
