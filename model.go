package general_counter

const TABLE_NAME_G_COUNTER = "g_counter"
const TABLE_NAME_G_COUNTER_DAILY_AGG = "g_counter_daily_agg"
const TABLE_NAME_G_COUNTER_DETAIL = "g_counter_detail"

// used only in sqldb
// table name :g_counter
type GCounterModel struct {
	Sql_id int64  `json:"sql_id" gorm:"type:bigint(20);primaryKey;"` // auto increasement
	Id     string `json:"id" gorm:"type:varchar(512);uniqueIndex;"`  // this is elastic id , [gkey]:[gtype]
	Gkey   string `json:"gkey" gorm:"type:varchar(512);index;"`      // can be anything like 'userid','accountid',etc.
	Gtype  string `json:"gtype" gorm:"type:varchar(512);index;"`     // can be anything like 'user_credit','account_traffic',etc.
	Amount int64  `json:"amount" gorm:"type:bigint(20);index;"`
}

func (model *GCounterModel) TableName() string {
	return TABLE_NAME_G_COUNTER
}

const (
	upload_status_uploading = "uploading"
	upload_status_uploaded  = "uploaded"
	upload_status_to_upload = "to_upload"
)

// used in elastic search db
// table name:g_counter_daily_agg
type GCounterDailyAggModel struct {
	Sql_id int64  `json:"sql_id" gorm:"type:bigint(20);primaryKey;"` // db id ,auto increasement
	Id     string `json:"id" gorm:"type:varchar(512);uniqueIndex;"`  // this is elastic id ,[date]:[gkey]:[gtype] => elastic search id
	Gkey   string `json:"gkey" gorm:"type:varchar(512);index;"`      // can be anything like 'userid','accountid',etc.
	Gtype  string `json:"gtype" gorm:"type:varchar(512);index;"`     // can be anything like 'user_credit','account_traffic',etc.
	Date   string `json:"date" gorm:"type:date;index;"`
	Amount int64  `json:"amount" gorm:"type:bigint(20);"`
	Status string `json:"status" gorm:"type:varchar(32);index;"`
}

func (model *GCounterDailyAggModel) TableName() string {
	return TABLE_NAME_G_COUNTER_DAILY_AGG
}

// used in elastic search db
// table name :g_counter_detail
type GCounterDetailModel struct {
	Sql_id   int64  `json:"sql_id"  gorm:"type:bigint(20);primaryKey;"` // db id ,auto increasement
	Id       string `json:"id" gorm:"type:varchar(512);uniqueIndex;"`   // this is elastic id ,assign db_id => elastic search id
	Gkey     string `json:"gkey" gorm:"type:varchar(512);index;"`       // can be anything like 'userid','accountid',etc.
	Gtype    string `json:"gtype" gorm:"type:varchar(512);index;"`      // can be anything like 'user_credit','account_traffic',etc.
	Datetime string `json:"datetime" gorm:"type:datetime(6);index;"`
	Amount   int64  `json:"amount" gorm:"type:bigint(20);"`
	Msg      string `json:"msg" gorm:"type:longtext;"`
}

func (model *GCounterDetailModel) TableName() string {
	return TABLE_NAME_G_COUNTER_DETAIL
}
