package general_counter

import (
	"time"

	"github.com/coreservice-io/job"
)

const MAX_AGG_UPLOAD_ITEMS_NUM = 30

func (gcounter *GeneralCounter) startAggUploader() error {

	spr_jb_name := "gcounter_agg_uploader"
	err := gcounter.spr_job_mgr.AddSprJob(spr_jb_name)
	if err != nil {
		return err
	}

	job.Start(
		spr_jb_name,
		// job process
		func() {
			if gcounter.spr_job_mgr.IsMaster(spr_jb_name) {

				date := time.Now().UTC().Format("2006-01-02")

				for {
					var agg_list []*GCounterDailyAggModel
					err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DAILY_AGG).Where("date != ?", date).Order("id asc").Limit(MAX_AGG_UPLOAD_ITEMS_NUM).Find(&agg_list).Error
					if err != nil {
						if gcounter.logger != nil {
							gcounter.logger.Errorln(spr_jb_name+"job sql err:", err)
						}
						return
					} else {

						if len(agg_list) == 0 {
							return
						}

						logs := []interface{}{}
						for _, agg := range agg_list {
							logs = append(logs, agg)
						}

						sids, add_log_err := gcounter.ecs_uplaoder.AddLogs_Sync(gcounter.gcounter_config.Project_name+"_"+TABLE_NAME_G_COUNTER_DAILY_AGG, logs)

						if add_log_err != nil {
							if gcounter.logger != nil {
								gcounter.logger.Errorln(spr_jb_name+" upload log err:", err)
							}
							return
						}

						if len(sids) > 0 {
							d_err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DAILY_AGG).Where("id in ?", sids).Delete(&GCounterDailyAggModel{}).Error
							if d_err != nil {
								if gcounter.logger != nil {
									gcounter.logger.Errorln(spr_jb_name+" agg del sql err:", d_err)
								}

								return
							}
						}
					}

				}
			}
		},
		// onPanic callback, run if panic happened
		func(err interface{}) {
			if gcounter.logger != nil {
				gcounter.logger.Errorln(spr_jb_name, err)
			}
		},
		// job interval in seconds
		30,
		job.TYPE_PANIC_REDO,
		// check continue callback, the job will stop running if return false
		// the job will keep running if this callback is nil
		func(job *job.Job) bool {
			return true
		},
		// onFinish callback
		func(inst *job.Job) {
			if gcounter.logger != nil {
				gcounter.logger.Debugln(spr_jb_name + " spr job stop")
			}
		},
	)

	return nil
}
