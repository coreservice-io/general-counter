package general_counter

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/coreservice-io/ecs_uploader/uploader"
	"github.com/coreservice-io/gorm_log"
	"github.com/coreservice-io/log"
	"github.com/coreservice-io/redis_spr"
	elasticSearch "github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GeneralCounterConfig struct {
	Project_name           string
	Agg_record_expire_days int
	Db_config              *DBConfig
	Ecs_config             *EcsConfig
	Redis_config           *RedisConfig
}

type DBConfig struct {
	Host     string
	Port     int
	DbName   string
	UserName string
	Password string
}

type EcsConfig struct {
	Address  string
	UserName string
	Password string
}

type RedisConfig struct {
	Addr     string
	Port     int
	UserName string
	Password string
	Prefix   string
	UseTLS   bool
}

// /////

type GeneralCounter struct {
	db              *gorm.DB
	ecs             *elasticSearch.Client
	ecs_uplaoder    *uploader.Uploader
	spr_job_mgr     *redis_spr.SprJobMgr
	logger          log.Logger
	gcounter_config *GeneralCounterConfig
}

func NewGeneralCounter(gc_config *GeneralCounterConfig, logger log.Logger) (*GeneralCounter, error) {

	if gc_config.Project_name == "" {
		return nil, errors.New("name is required")
	}

	if gc_config.Agg_record_expire_days <= 0{
		return nil, errors.New("agg record expire days must > 0")
	}

	if logger == nil {
		return nil, errors.New("logger is required")
	}

	// //
	spr_jm, spr_job_err := redis_spr.New(redis_spr.RedisConfig{
		Addr:     gc_config.Redis_config.Addr,
		Port:     gc_config.Redis_config.Port,
		Password: gc_config.Redis_config.Password,
		UserName: gc_config.Redis_config.UserName,
		Prefix:   gc_config.Redis_config.Prefix + ":" + "gcounter",
		UseTLS:   gc_config.Redis_config.UseTLS,
	})
	if spr_job_err != nil {
		return nil, spr_job_err
	}

	// //

	ecs_uploader, uploader_err := uploader.New(gc_config.Ecs_config.Address, gc_config.Ecs_config.UserName, gc_config.Ecs_config.Password)
	if uploader_err != nil {
		return nil, uploader_err
	}

	// //db config check////
	db_conf := gc_config.Db_config
	dsn := db_conf.UserName + ":" + db_conf.Password + "@tcp(" + db_conf.Host + ":" + strconv.Itoa(db_conf.Port) + ")/" + db_conf.DbName + "?charset=utf8mb4&loc=UTC"

	db_log_level := gorm_log.Warn
	if logger.GetLevel() >= log.TraceLevel {
		db_log_level = gorm_log.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gorm_log.New_gormLocalLogger(logger, gorm_log.Config{
			SlowThreshold:             500 * time.Millisecond,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  db_log_level,
		}),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(10)

	// db table check
	if !db.Migrator().HasTable(TABLE_NAME_G_COUNTER) {
		return nil, errors.New("table " + TABLE_NAME_G_COUNTER + " not exist in db")
	}
	if !db.Migrator().HasTable(TABLE_NAME_G_COUNTER_DAILY_AGG) {
		return nil, errors.New("table " + TABLE_NAME_G_COUNTER_DAILY_AGG + " not exist in db")
	}
	if !db.Migrator().HasTable(TABLE_NAME_G_COUNTER_DETAIL) {
		return nil, errors.New("table " + TABLE_NAME_G_COUNTER_DETAIL + " not exist in db")
	}

	// ecs config check
	ecs_conf := gc_config.Ecs_config
	ecs, err := elasticSearch.NewClient(
		elasticSearch.SetURL(ecs_conf.Address),
		elasticSearch.SetBasicAuth(ecs_conf.UserName, ecs_conf.Password),
		elasticSearch.SetSniff(false),
		elasticSearch.SetHealthcheckInterval(30*time.Second),
		elasticSearch.SetGzip(true),
	)
	if err != nil {
		return nil, err
	}

	projectName := gc_config.Project_name
	// check ecs table
	existService := ecs.IndexExists(projectName+"_"+TABLE_NAME_G_COUNTER_DETAIL, projectName+"_"+TABLE_NAME_G_COUNTER_DAILY_AGG)
	index_exist, index_err := existService.Do(context.Background())
	if index_err != nil {
		return nil, errors.New("elastic check false, err:" + index_err.Error())
	}
	if !index_exist {
		return nil, errors.New("elastic index not exist")
	}

	//
	gcounter := &GeneralCounter{
		db:              db,
		ecs:             ecs,
		ecs_uplaoder:    ecs_uploader,
		spr_job_mgr:     spr_jm,
		logger:          logger,
		gcounter_config: gc_config,
	}

	// start uploaders

	upload_err := gcounter.startAggUploader()
	if upload_err != nil {
		return nil, upload_err
	}
	upload_err = gcounter.startDetailUploader()
	if upload_err != nil {
		return nil, upload_err
	}
	d_err := gcounter.deleteExpireUploadedAggRecords(gc_config.Agg_record_expire_days)
	if d_err != nil {
		return nil, d_err
	}

	return gcounter, nil

}
