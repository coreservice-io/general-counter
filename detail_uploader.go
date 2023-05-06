package general_counter

import (
	"context"
	"fmt"

	"github.com/coreservice-io/job"
)

const MAX_DETAIL_UPLOAD_ITEMS_NUM = 30

const detail_upload_interval_secs = 30

func (gcounter *GeneralCounter) startDetailUploader() error {

	spr_jb_name := "gcounter_detail_uploader"
	err := gcounter.spr_job_mgr.AddSprJob(context.Background(), spr_jb_name)
	if err != nil {
		return err
	}

	job.Start(context.Background(), job.JobConfig{
		Name:          spr_jb_name,
		Job_type:      job.TYPE_PANIC_REDO,
		Interval_secs: detail_upload_interval_secs,
		Process_fn: func(j *job.Job) {
			if gcounter.spr_job_mgr.IsMaster(spr_jb_name) {
				tblname := gcounter.gcounter_config.Project_name + "_" + TABLE_NAME_G_COUNTER_DETAIL

				for {
					var detail_list []*GCounterDetailModel
					err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DETAIL).Order("id asc").Limit(MAX_DETAIL_UPLOAD_ITEMS_NUM).Find(&detail_list).Error
					if err != nil {
						gcounter.logger.Errorln(spr_jb_name+"job sql err:", err)
						return
					} else {

						if len(detail_list) == 0 {
							return
						}

						logs := []interface{}{}
						for _, detail := range detail_list {
							log := &GCounterDetailEsModel{
								Sql_id:   detail.Sql_id,
								Id:       detail.Id,
								Gkey:     detail.Gkey,
								Gtype:    detail.Gtype,
								Datetime: detail.Datetime[:19],
								Amount:   detail.Amount.String(),
								Msg:      detail.Msg,
							}
							logs = append(logs, log)
						}

						fmt.Println("uploading ==================")
						sids, add_log_err := gcounter.ecs_uplaoder.AddLogs_Sync(tblname, logs)
						if add_log_err != nil {
							gcounter.logger.Errorln(spr_jb_name+" upload "+tblname+" log err:", add_log_err)
							return
						}

						if len(sids) > 0 {
							d_err := gcounter.db.Table(TABLE_NAME_G_COUNTER_DETAIL).Where("id in ?", sids).Delete(&GCounterDetailModel{}).Error
							if d_err != nil {
								gcounter.logger.Errorln(spr_jb_name+" detail del sql err:", d_err)
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
						gcounter.logger.Errorln(spr_jb_name+" after detail del counter sql err:", count_err)
						return
					}

					if left_records < MAX_DETAIL_UPLOAD_ITEMS_NUM {
						return
					}
				}
			}
		},

		On_panic: func(job *job.Job, panic_err interface{}) {
			gcounter.logger.Errorln(spr_jb_name, err)
		},

		Final_fn: func(j *job.Job) {
			gcounter.logger.Errorln(spr_jb_name + " spr job stop")
		},
	}, nil)

	return nil
}
