package general_counter

const TABLE_NAME_G_COUNTER = "g_counter"
const TABLE_NAME_G_COUNTER_DAILY_AGG = "g_counter_daily_agg"
const TABLE_NAME_G_COUNTER_DETAIL = "g_counter_detail"

// used only in sqldb
// table name :g_counter
type GCounterModel struct {
	Sql_id int64  `json:"sql_id" gorm:"primaryKey"` // auto increasement
	Id     string `json:"id" gorm:"index;unique"`   // this is elastic id , [gkey]:[gtype]
	Gkey   string `json:"gkey" gorm:"index"`        // can be anything like 'userid','accountid',etc.
	Gtype  string `json:"gtype" gorm:"index"`       // can be anything like 'user_credit','account_traffic',etc.
	Amount int64  `json:"amount" gorm:"index"`
}

// used in elastic search db
// table name:g_counter_daily_agg
type GCounterDailyAggModel struct {
	Sql_id int64  `json:"sql_id" gorm:"primaryKey"` // db id ,auto increasement
	Id     string `json:"id" gorm:"index;unique"`   // this is elastic id ,[date]:[gkey]:[gtype] => elastic search id
	Gkey   string `json:"gkey" gorm:"index"`        // can be anything like 'userid','accountid',etc.
	Gtype  string `json:"gtype" gorm:"index"`       // can be anything like 'user_credit','account_traffic',etc.
	Date   string `json:"date" gorm:"index"`
	Amount int64  `json:"amount"`
}

// used in elastic search db
// table name :g_counter_detail
type GCounterDetailModel struct {
	Sql_id   int64  `json:"sql_id" gorm:"primaryKey"` // db id ,auto increasement
	Id       string `json:"id" gorm:"index;unique"`   // this is elastic id ,assign db_id => elastic search id
	Gkey     string `json:"gkey" gorm:"index"`        // can be anything like 'userid','accountid',etc.
	Gtype    string `json:"gtype" gorm:"index"`       // can be anything like 'user_credit','account_traffic',etc.
	Datetime string `json:"datetime" gorm:"index"`
	Amount   int64  `json:"amount"`
	Msg      string `json:"msg"`
}
