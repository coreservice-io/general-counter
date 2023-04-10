package general_counter

import (
	"context"
	"time"

	"github.com/coreservice-io/job"
	"gorm.io/gorm"
)

const MAX_AGG_UPLOAD_ITEMS_NUM = 30

const agg_upload_interval_secs = 30

func (gcounter *GeneralCounter) startAggUploader() error {

	spr_jb_name := "gcounter_agg_uploader"
	err := gcounter.spr_job_mgr.AddSprJob(context.Background(), spr_jb_name)
	if err != nil {
		return err
	}

	job.Start(context.Background(), job.JobConfig{
		Name:          spr_jb_name,
		Job_type:      job.TYPE_PANIC_REDO,
		Interval_secs: agg_upload_interval_secs,
		Process_fn: func(j *job.Job) {
			if gcounter.spr_job_mgr.IsMaster(spr_jb_name) {

				date := time.Now().UTC().Format("2006-01-02")

				for {
					var agg_list []*GCounterDailyAggModel
					err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DAILY_AGG).Where("status IN ? AND date != ?", []string{upload_status_uploading, upload_status_to_upload}, date).Order("id asc").Limit(MAX_AGG_UPLOAD_ITEMS_NUM).Find(&agg_list).Error
					if err != nil {
						gcounter.logger.Errorln(spr_jb_name+"job sql err:", err)
						return
					} else {

						if len(agg_list) == 0 {
							return
						}

						ids := []string{}
						for _, agg := range agg_list {
							ids = append(ids, agg.Id)
						}

						// update status => uploading
						d_err := gcounter.db.Transaction(func(tx *gorm.DB) error {
							return tx.Table(TABLE_NAME_G_COUNTER_DAILY_AGG).Where("id in ?", ids).Update("status", upload_status_uploading).Error
						})
						if d_err != nil {
							gcounter.logger.Errorln(spr_jb_name+" agg update status to uploading sql err:", d_err)
							return
						}

						logs := []interface{}{}
						for _, agg := range agg_list {
							logs = append(logs, agg)
						}

						sids, add_log_err := gcounter.ecs_uplaoder.AddLogs_Sync(gcounter.gcounter_config.Project_name+"_"+TABLE_NAME_G_COUNTER_DAILY_AGG, logs)

						if add_log_err != nil {
							gcounter.logger.Errorln(spr_jb_name+" upload log err:", err)
							return
						}

						if len(sids) > 0 {
							// update status => uploaded
							d_err := gcounter.db.Transaction(func(tx *gorm.DB) error {
								return tx.Table(TABLE_NAME_G_COUNTER_DAILY_AGG).Where("id in ? AND status = ?", sids, upload_status_uploading).Update("status", upload_status_uploaded).Error
							})
							if d_err != nil {
								gcounter.logger.Errorln(spr_jb_name+" agg update sql err:", d_err)
								return
							}
						}
					}
				}
			}
		},
		On_panic: func(job *job.Job, panic_err interface{}) {
			gcounter.logger.Errorln(spr_jb_name, err)
		},
		Final_fn: func(j *job.Job) {
			gcounter.logger.Debugln(spr_jb_name + " spr job stop")
		},
	}, nil)

	return nil
}

const delete_expire_agg_interval_secs = 1800

func (gcounter *GeneralCounter) deleteExpireUploadedAggRecords(agg_record_expire_days int) error {
	spr_jb_name := "gcounter_agg_delete_expire_uploaded"
	err := gcounter.spr_job_mgr.AddSprJob(context.Background(), spr_jb_name)
	if err != nil {
		return err
	}

	job.Start(
		context.Background(),
		job.JobConfig{
			Name:          spr_jb_name,
			Job_type:      job.TYPE_PANIC_REDO,
			Interval_secs: delete_expire_agg_interval_secs,
			Process_fn: func(j *job.Job) {
				if gcounter.spr_job_mgr.IsMaster(spr_jb_name) {
					// delete uploaded expire record
					agg_record_expire_days = agg_record_expire_days + 1 //+1 for safety boundary
					date := time.Now().UTC().AddDate(0, 0, -agg_record_expire_days).Format("2006-01-02")
					err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DAILY_AGG).Where("date < ? AND status = ?", date, upload_status_uploaded).Delete(&GCounterDailyAggModel{}).Error
					if err != nil {
						gcounter.logger.Errorln(spr_jb_name+" agg del sql err:", err)
					}
				}
			},
			On_panic: func(job *job.Job, panic_err interface{}) {
				gcounter.logger.Errorln(spr_jb_name, err)
			},
			Final_fn: func(j *job.Job) {
				gcounter.logger.Debugln(spr_jb_name + " spr job stop")
			},
		},
		nil,
	)

	return nil
}
