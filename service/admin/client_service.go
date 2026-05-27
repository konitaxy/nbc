package admin

import (
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/client/request"
)

type ClientService struct {
}

func (*ClientService) GetClient(id uint) (client client.Client, err error) {
	err = global.GVA_DB.Preload("DueDiligence").Where("id = ?", id).First(&client).Error
	return
}
func (*ClientService) GetClientByNo(clientNo string) (client client.Client, err error) {
	err = global.GVA_DB.Where("client_no = ?", clientNo).First(&client).Error
	return
}

func (*ClientService) GetDueDiligenceByClientID(clientId uint) (due client.ClientDueDiligence, err error) {
	err = global.GVA_DB.Where("client_id = ?", clientId).First(&due).Error
	return
}
func (*ClientService) SaveClientDueKYC(due client.ClientDueDiligence) (err error) {
	err = global.GVA_DB.Save(&due).Error
	return
}
func (*ClientService) Save(client *client.Client) (err error) {
	err = global.GVA_DB.Save(&client).Error
	return
}

func (*ClientService) ListClient(search request.ClientSearchParams) (total int64, list []*client.Client, err error) {
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

	query := global.GVA_DB.Model(&client.Client{}).Preload("InviteUser").Order(orderBy).Where("1= ?", 1)
	if search.Email != "" {
		conditions = append(conditions, "email = ?")
		args = append(args, search.Email)
	}
	if search.ClientNo != "" {
		conditions = append(conditions, "client_no = ?")
		args = append(args, search.ClientNo)
	}
	if search.EnName != "" {
		conditions = append(conditions, "en_name = ?")
		args = append(args, search.EnName)
	}
	if search.AccountManager != "" {
		conditions = append(conditions, "account_manager = ?")
		args = append(args, search.AccountManager)
	}
	if search.ClientReviewStatus > 0 {
		conditions = append(conditions, "client_review_status = ?")
		args = append(args, search.ClientReviewStatus)
	}
	if search.ClientStatus > 0 {
		conditions = append(conditions, "client_status = ?")
		args = append(args, search.ClientStatus)
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
