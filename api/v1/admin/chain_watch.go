package admin

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/finance/request"
	"gorm.io/gorm"
)

func (FinanceManagerApi) AddChainWatchAddress(c *gin.Context) {
	var req request.ChainWatchAddressAddReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if strings.TrimSpace(req.Address) == "" {
		response.FailWithMessage("address is required", c)
		return
	}
	if strings.TrimSpace(req.ChainType) == "" {
		req.ChainType = "TRON"
	}
	row, err := financeService.AddChainWatchAddress(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(row, c)
}

func (FinanceManagerApi) DeleteChainWatchAddress(c *gin.Context) {
	var req request.ChainWatchAddressDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := financeService.DeleteChainWatchAddress(req.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.FailWithMessage("record not found", c)
			return
		}
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("success", c)
}

func (FinanceManagerApi) ListChainWatchAddress(c *gin.Context) {
	var req request.ChainWatchAddressListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	total, list, err := financeService.ListChainWatchAddress(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}

func (FinanceManagerApi) ListChainInboundTransaction(c *gin.Context) {
	var req request.ChainInboundTransactionListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	total, list, err := financeService.ListChainInboundTransaction(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}
