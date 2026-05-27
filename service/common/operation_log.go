package common

import (
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/common/request"
	"go.uber.org/zap"
)

type LogService struct {
}

func init() {
	go func() {
		for log := range global.LogChannel {
			v := log.(common.OpLog)
			err := global.GVA_DB.Create(&v).Error
			if err != nil {
				global.GVA_LOG.Error("写入日志失败", zap.Error(err))
			}
		}
	}()
}
func (*LogService) Save(log *common.OpLog) (err error) {
	return global.GVA_DB.Create(log).Error
}

func (*LogService) List(search request.OpLogSearchParams) (total int64, list []common.OpLog, err error) {
	// 设置默认值
	if search.Page <= 0 {
		search.Page = 1
	}

	if search.PageSize <= 0 {
		search.PageSize = 10
	}
	var orderBy = "created_at DESC"
	if search.OrderBy == 1 {
		orderBy = "created_at DESC"
	}
	// lastMonth := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	// 构建查询条件
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&common.OpLog{}).Order(orderBy).Where("1= ?", 1)
	if search.Who > 0 {
		conditions = append(conditions, "who = ?")
		args = append(args, search.Who)
	}
	// if search.ClientNo != "" {
	// 	conditions = append(conditions, "client_no = ?")
	// 	args = append(args, search.ClientNo)
	// }
	if search.ObjId > 0 {
		conditions = append(conditions, "obj_id = ?")
		args = append(args, search.ObjId)
	}
	if search.OpType > 0 {
		conditions = append(conditions, "op_type = ?")
		args = append(args, search.OpType)
	}
	if search.Name != "" {
		conditions = append(conditions, "name = ?")
		args = append(args, search.Name)
	}
	if search.Source > 0 {
		conditions = append(conditions, "source = ?")
		args = append(args, search.Source)
	}
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}

	query.Count(&total)
	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error

	return
}

func (*LogService) ListSmsCode(search request.SmsCodeSearchParams) (total int64, list []common.SmsCode, err error) {
	if search.Page <= 0 {
		search.Page = 1
	}

	if search.PageSize <= 0 {
		search.PageSize = 10
	}
	var orderBy = "created_at DESC"
	if search.OrderBy == 1 {
		orderBy = "created_at DESC"
	}
	// 构建查询条件
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&common.SmsCode{}).Order(orderBy).Where("1= ?", 1)

	if search.To != "" {
		conditions = append(conditions, "`to` = ?")
		args = append(args, search.To)
	}

	if search.EventType != "" {
		conditions = append(conditions, "event_type = ?")
		args = append(args, search.EventType)
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	query.Count(&total)
	// 分页查询
	offset := (search.Page - 1) * search.PageSize
	err = query.Limit(search.PageSize).Offset(offset).Find(&list).Error
	return
}
