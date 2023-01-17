package general_counter

import (
	"github.com/coreservice-io/job"
)

const MAX_DETAIL_UPLOAD_ITEMS_NUM = 30

func (gcounter *GeneralCounter) startDetailUploader() error {

	spr_jb_name := "gcounter_detail_uploader"
	err := gcounter.spr_job_mgr.AddSprJob(spr_jb_name)
	if err != nil {
		return err
	}

	job.Start(
		spr_jb_name,
		// job process
		func() {
			if gcounter.spr_job_mgr.IsMaster(spr_jb_name) {
				for {
					var detail_list []*GCounterDetailModel
					err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DETAIL).Order("id asc").Limit(MAX_DETAIL_UPLOAD_ITEMS_NUM).Find(&detail_list).Error
					if err != nil {
						if gcounter.logger != nil {
							gcounter.logger.Errorln(spr_jb_name+"job sql err:", err)
						}
						return
					} else {

						if len(detail_list) == 0 {
							return
						}

						logs := []interface{}{}
						for _, detail := range detail_list {
							detail.Datetime = detail.Datetime[:19]
							logs = append(logs, detail)
						}

						sids, add_log_err := gcounter.ecs_uplaoder.AddLogs_Sync(gcounter.gcounter_config.Project_name+"_"+TABLE_NAME_G_COUNTER_DETAIL, logs)

						if add_log_err != nil {
							if gcounter.logger != nil {
								gcounter.logger.Errorln(spr_jb_name+" upload log err:", err)
							}
							return
						}

						if len(sids) > 0 {
							d_err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DETAIL).Where("id in ?", sids).Delete(&GCounterDetailModel{}).Error
							if d_err != nil {
								if gcounter.logger != nil {
									gcounter.logger.Errorln(spr_jb_name+" detail del sql err:", d_err)
								}
								return
							}
						}
					}

					var left_records int64
					count_err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DETAIL).Count(&left_records).Error
					if count_err != nil {
						return
					}
					if count_err != nil {
						if gcounter.logger != nil {
							gcounter.logger.Errorln(spr_jb_name+" after detail del counter sql err:", count_err)
						}
						return
					}

					if left_records < MAX_DETAIL_UPLOAD_ITEMS_NUM {
						return
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
				gcounter.logger.Errorln(spr_jb_name + " spr job stop")
			}
		},
	)

	return nil
}
