package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/common/request"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/utils"
	"go.uber.org/zap"
)

type MaintainManagerApi struct {
}

func (*MaintainManagerApi) ListLog(c *gin.Context) {
	var req request.OpLogSearchParams
	_ = c.ShouldBindJSON(&req)

	if total, lost, err := logService.List(req); err == nil {
		response.OkWithDetailed(response.PageResult{
			List:  lost,
			Total: total,
		}, "Success", c)
	} else {
		global.GVA_LOG.Error("List log failed!", zap.Error(err))
		response.FailWithMessage("Failed", c)

	}
}

func (*MaintainManagerApi) Set(c *gin.Context) {
	var req common.AdminConfig
	_ = c.ShouldBindJSON(&req)

	if err := utils.Verify(req, utils.ConfigVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info := utils.GetUserInfo(c)
	req.Operator = info.Username
	if cfg, _ := configService.Get(req.Kay); cfg.ID == 0 {
		configService.Set(&req)
		global.Push(common.OpLog{
			OpType: common.OpType_System_Config_Set,
			Who:    info.ID,
			Name:   info.Username,
			Detail: fmt.Sprintf("new: %s", utils.ObjectToJson(cfg)),
		})
	} else {
		detail := fmt.Sprintf("old: %s, new: %s", utils.ObjectToJson(cfg), utils.ObjectToJson(req))
		cfg.ValueType = req.ValueType
		switch cfg.ValueType {
		case "number":
			cfg.NumberValue = req.NumberValue
		case "string":
			cfg.StringValue = req.StringValue
		case "boolean":
			cfg.BooleanValue = req.BooleanValue
		case "json":
			cfg.JsonValue = req.JsonValue
		case "date":
			cfg.DateValue = req.DateValue

		}
		configService.Set(&req)
		global.Push(common.OpLog{
			OpType: common.OpType_System_Config_Set,
			Who:    info.ID,
			Name:   info.Username,
			Detail: detail,
		})
	}
}

func (*MaintainManagerApi) Get(c *gin.Context) {
	var req common.AdminConfig
	_ = c.ShouldBindJSON(&req)

	cfg, _ := configService.Get(req.Kay)
	response.OkWithData(cfg, c)
}

func (*MaintainManagerApi) ListSmsCode(c *gin.Context) {
	var req request.SmsCodeSearchParams
	_ = c.ShouldBindJSON(&req)

	total, list, _ := logService.ListSmsCode(req)
	response.OkWithData(gin.H{
		"total": total,
		"list":  list,
	}, c)
}
