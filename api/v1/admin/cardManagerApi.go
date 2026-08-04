package admin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/common/response"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/finance"
	"gitlab.com/ucard/model/finance/request"
	"gitlab.com/ucard/model/provider/cardbin"
	financesvc "gitlab.com/ucard/service/finance"
	"gitlab.com/ucard/utils"
	"gitlab.com/ucard/utils/transaction"
	"go.uber.org/zap"
)

type CardManagerApi struct {
}

func (*CardManagerApi) CardbinHook(c *gin.Context) {
	var event cardbin.EventHook
	_ = c.ShouldBindJSON(&event)
	if e := event.Unmarshal(); e != nil {
		global.GVA_LOG.Error("解析失败", zap.Error(e))
		global.GVA_LOG.Info("解析失败", zap.Any("data", event.Data))
		c.Status(http.StatusInternalServerError)
	} else {
		switch event.EventType {
		case "CardOperate":
			if v, e := event.ParseData.(cardbin.CardOperate); !e {
				global.GVA_LOG.Error("解析失败", zap.Any("data", event))
				c.Status(http.StatusInternalServerError)
			} else {
				transactionType := transaction.NormalizeTransactionType(v.OperateType, "cardbin")
				if transactionType == "" {
					global.GVA_LOG.Error("未知交易类型", zap.Any("data", event))
					return
				}
				// if v.PartnerOrderID == "" {
				// 	global.GVA_LOG.Info("无订单号", zap.Any("data", event))
				// 	c.JSON(http.StatusOK, "ok")
				// 	return
				// }
				// if utils.GetEnvCode(v.PartnerOrderID) != strconv.Itoa(global.GVA_CONFIG.System.EnvCode) {
				// 	global.GVA_LOG.Info("环境不匹配", zap.Any("data", event))
				// 	c.JSON(http.StatusOK, "ok")
				// 	return
				// }
				card, _ := financeService.GetCardByCardID(v.CardID)
				if card.ID == 0 {
					global.GVA_LOG.Info("未找到该卡", zap.String("cardID", v.CardID), zap.Any("event", event))
					c.JSON(http.StatusOK, "ok")
					return
				}
				if t, _ := financeService.GetCardTransactionByTransactionID(v.TransactionID, transactionType); t.ID == 0 {
					transaction := finance.CardTransactionRecord{
						Amount:          decimal.NewFromFloat(v.Amount),
						Channel:         constant.Channel_Cardbin,
						CardID:          v.CardID,
						Currency:        v.Currency,
						EventType:       event.EventType,
						OrderID:         v.PartnerOrderID,
						Fee:             decimal.NewFromFloat(v.MerchantFee.TotalFeeAmount),
						Status:          v.Status,
						TransactionType: transactionType,
						TransactionID:   v.TransactionID,
						TransactionTime: time.Now(),
					}
					if err2 := financeService.AddCardApplyTransaction(&transaction); err2 != nil {
						global.GVA_LOG.Error("add transaction error", zap.Error(err2))
						c.Status(http.StatusInternalServerError)
						return
					} else {
						go func() {
							financeService.SyncCardDetailSkipCVV(v.PartnerOrderID, v.CardID)
						}()
					}
				} else {
					c.JSON(http.StatusOK, "ok")
					return
				}

			}
		case "CardApply":
			if v, e := event.ParseData.(cardbin.CardApply); !e {
				global.GVA_LOG.Error("解析失败", zap.Any("data", event))
				c.Status(http.StatusInternalServerError)
			} else {
				// if v.PartnerOrderID == "" {
				// 	global.GVA_LOG.Info("无订单号", zap.Any("data", event))
				// 	c.JSON(http.StatusOK, "ok")
				// 	return
				// }
				// if utils.GetEnvCode(v.PartnerOrderID) != strconv.Itoa(global.GVA_CONFIG.System.EnvCode) {
				// 	global.GVA_LOG.Info("环境不匹配", zap.Any("data", event))
				// 	c.JSON(http.StatusOK, "ok")
				// 	return
				// }
				card, _ := financeService.GetCardByCardID(v.CardID)
				if card.ID == 0 {
					global.GVA_LOG.Info("未找到该卡", zap.String("cardID", v.CardID), zap.Any("event", event))
					c.JSON(http.StatusOK, "ok")
					return
				}
				go func() {
					financeService.SyncCardDetailSkipCVV(v.PartnerOrderID, v.CardID)
				}()

			}
		case "Authorization":
			if v, e := event.ParseData.(cardbin.TransactionNotify); !e {
				global.GVA_LOG.Error("解析失败", zap.Any("data", event))
				c.Status(http.StatusInternalServerError)
				return
			} else {
				transactionType := transaction.NormalizeTransactionType(v.TransactionType, "cardbin")
				if transactionType == "" {
					global.GVA_LOG.Error("未知交易类型", zap.Any("data", event))
					return
				}
				card, _ := financeService.GetCardByCardID(v.CardID)
				if card.ID == 0 {
					global.GVA_LOG.Error("未找到该卡", zap.Any("cardID", event))
					c.JSON(http.StatusOK, "ok")
					return
				}
				t, _ := financeService.GetCardTransactionByTransactionID(v.TransactionID, transactionType)
				if t.ID == 0 {
					transaction := finance.CardTransactionRecord{
						Amount:          decimal.NewFromFloat(v.BillingAmount),
						Currency:        v.BillingCurrency,
						OriginAmount:    decimal.NewFromFloat(v.TransactionAmount),
						OriginCurrency:  v.BillingCurrency,
						Channel:         constant.Channel_Cardbin,
						CardID:          v.CardID,
						EventType:       event.EventType,
						Fee:             decimal.NewFromFloat(v.MerchantFee.TotalFeeAmount),
						FeeDetail:       financesvc.FeeDetailFromValue(v.MerchantFee),
						Status:          v.TransactionStatus,
						TransactionType: transactionType,
						TransactionID:   v.TransactionID,
						TransactionTime: time.UnixMilli(v.TransactionTime),
						CrossBoardType:  v.CrossBoardType,
						AuthCode:        v.AuthCode,
						MerchantName:    v.MerchantName,
						FailReason:      v.FailReason,
						ReferenceID:     v.ReferenceID,
					}
					if err2 := financeService.AddCardApplyTransaction(&transaction); err2 != nil {
						global.GVA_LOG.Error("add transaction error", zap.Error(err2))
						c.Status(http.StatusInternalServerError)
						return
					} else {
						go func() {
							financeService.SyncCardDetailSkipCVV("", v.CardID)
						}()
					}
				}

			}
		case "Inbound":
			if v, e := event.ParseData.(cardbin.RechargeOrder); !e {
				global.GVA_LOG.Error("解析失败", zap.Any("data", event))
				c.Status(http.StatusInternalServerError)
				return
			} else {
				if v.OriginalAmount.LessThanOrEqual(decimal.Zero) {
					global.GVA_LOG.Info("bad zero amount", zap.Any("event inbound", v))
					c.Status(http.StatusInternalServerError)
					return
				}
				if rr, _ := clientFinanceService.GetWalletRechargeByOrderID(v.PartnerOrderID, string(constant.RechargeStatus_PENDING)); rr.ID > 0 {
					if v.State == string(constant.RechargeStatus_SUCCESS) {
						rr.Operator = "system"
						rr.FinishTime = utils.Now()
						rr.Status = constant.RechargeStatus_SUCCESS
						if err := clientFinanceService.WalletRecharge(&rr); err == nil {
							global.GVA_LOG.Info("wallet recharge success", zap.Any("recharge:", rr))
						} else {
							global.GVA_LOG.Error("wallet recharge failed", zap.Error(err))
							c.Status(http.StatusInternalServerError)
							return
						}
					} else if v.State == string(constant.RechargeStatus_FAILED) {
						rr.Operator = "system"
						rr.FinishTime = utils.Now()
						rr.Status = constant.RechargeStatus_FAILED
						if err := clientFinanceService.SaveWalletRecharge(&rr); err != nil {
							global.GVA_LOG.Error("save wallet recharge failed", zap.Error(err))
							c.Status(http.StatusInternalServerError)
							return
						}
					}
				}

			}

		}
	}
	c.JSON(http.StatusOK, "ok")
}
func (*CardManagerApi) CreateCardBin(c *gin.Context) {
	var req finance.CardBin
	info := utils.GetUserInfo(c)
	_ = c.ShouldBindJSON(&req)
	if err := utils.Verify(req, utils.AddCardBinVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := cardService.SaveCardBin(&req); err == nil {
		global.Push(common.OpLog{
			Who:    *&info.ID,
			Name:   info.NickName,
			OpType: common.OpType_CardBin_Create,
			Detail: fmt.Sprintf("Create CardBin: %s", req.CardBin),
			ObjId:  req.ID,
		})
		response.OkWithMessage("Success", c)
	} else {
		global.GVA_LOG.Error("Failed", zap.Error(err))
		response.FailWithMessage("Failed", c)
	}
}

func (*CardManagerApi) BlockCardBin(c *gin.Context) {
	carBinId, e := c.GetQuery("id")
	id := utils.GetUserID(c)
	name := utils.GetUserName(c)
	if !e {
		response.FailWithMessage("Failed", c)
		return
	}
	if cardbin, err := cardService.GetCardBinByCardBinId(carBinId); cardbin.ID > 0 {
		if err := cardService.BlockCardBin(&cardbin); err == nil {
			global.Push(common.OpLog{
				Who:    id,
				OpType: common.OpType_CardBin_Block,
				ObjId:  cardbin.ID,
				Detail: fmt.Sprintf("%s Block card bin %s.%s", name, cardbin.CardBin, carBinId),
				Source: 1,
			})
			response.OkWithMessage("Success", c)
			return
		}
	} else {
		global.GVA_LOG.Error("Block card bin failed!", zap.Error(err))
	}
	response.FailWithMessage("Failed", c)
}

func (*CardManagerApi) ListCardBin(c *gin.Context) {
	var req request.CardBinSearchParams
	_ = c.ShouldBindJSON(&req)
	if total, list, err := cardService.ListCardBin(req); err == nil {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "Success", c)
	} else {
		global.GVA_LOG.Error("List card bin failed!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}

}

func (f *CardManagerApi) ListCards(c *gin.Context) {
	var search request.CardSearchParams
	_ = c.ShouldBindJSON(&search)
	if len(search.CardNo) == 4 {
		search.CardNoSuffix = search.CardNo
		search.CardNo = ""
	}
	total, list, err := financeService.ListCards(&search, true)
	if err != nil {
		global.GVA_LOG.Error("list card bin failed", zap.Any("err", err))
		response.FailWithMessage("list card bin failed", c)
	} else {
		response.OkWithData(response.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}
func (f *CardManagerApi) ListCardTransaction(c *gin.Context) {
	var search request.CardTransactionSearchParams
	_ = c.ShouldBindJSON(&search)
	total, list, err := financeService.ListCardTransaction(&search, true)
	if err != nil {
		global.GVA_LOG.Error("list card bin failed", zap.Any("err", err))
		response.FailWithMessage("list card bin failed", c)
	} else {
		response.OkWithData(response.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}

func (f *CardManagerApi) CardCancel(c *gin.Context) {
	var req request.CancelCardReq
	_ = c.ShouldBindJSON(&req)
	if err := utils.Verify(req, utils.CancelCardBinVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if card, _ := financeService.GetCard(req.ID, req.ClientID); card.ID == 0 {
		response.FailWithMessage("Card not exist", c)
		return
	} else {
		if card.CardStatus != string(constant.CardStatus_ACTIVE) {
			response.FailWithMessage("Card not active", c)
			return
		}
		if cb, _ := cardService.GetCardBinByCardBinId(card.CardBinID); cb.ID == 0 {
			response.FailWithMessage("Card bin not exist", c)
			return
		} else {
			if !cb.CancelCard {
				response.FailWithMessage("Card bin not support cancel", c)
				return
			}
			if err := financeService.CancelCard(&card); err == nil {

				response.Ok(c)
				return
			} else {
				global.GVA_LOG.Error("termiate card failed", zap.Any("err", err))
				response.FailWithServiceError(c, err)
				return
			}
		}
	}
}
func (f *CardManagerApi) CardWithdraw(c *gin.Context) {
	var search request.CardTransactionSearchParams
	_ = c.ShouldBindJSON(&search)
	total, list, err := financeService.ListCardTransaction(&search, true)
	if err != nil {
		global.GVA_LOG.Error("list card bin failed", zap.Any("err", err))
		response.FailWithMessage("list card bin failed", c)
	} else {
		response.OkWithData(response.PageResult{
			List:  list,
			Total: total,
		}, c)
	}
}

func (f *CardManagerApi) SyncCard(c *gin.Context) {
	var search request.CardSearchParams
	_ = c.ShouldBindJSON(&search)

	if card, _ := financeService.GetCard(search.ID, search.ClientID); card.ID == 0 {
		response.FailWithMessage("card not found", c)
		return
	} else {
		if err := financeService.SyncCardDetail(card.OrderID, card.CardID); err != nil {
			global.GVA_LOG.Error("sync card detail failed", zap.Any("err", err))
			response.FailWithServiceError(c, err)
			return
		}
		response.Ok(c)
	}

}

func (f *CardManagerApi) CardFrozen(c *gin.Context) {
	var req request.CardFrozenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证 action 参数
	if req.Action != "frozen" && req.Action != "unfrozen" {
		response.FailWithMessage("action must be 'frozen' or 'unfrozen'", c)
		return
	}

	// 查询卡信息以获取 clientID（admin 可以操作所有卡，所以不限制 clientID）
	var card finance.PixielCard
	if err := global.GVA_DB.First(&card, "id = ?", req.ID).Error; err != nil {
		response.FailWithMessage("card not found", c)
		return
	}

	// 调用服务方法
	if err := financeService.CardFrozen(req.ID, card.ClientID, req.Action, req.Remark); err != nil {
		global.GVA_LOG.Error("card frozen/unfrozen failed", zap.Error(err))
		response.FailWithServiceError(c, err)
		return
	}

	// 记录操作日志
	info := utils.GetUserInfo(c)
	actionText := "冻结"
	var opType common.OpType = common.OpType_Card_Frozen
	if req.Action == "unfrozen" {
		actionText = "解冻"
		opType = common.OpType_Card_UnFrozen
	}
	global.Push(common.OpLog{
		Who:    info.ID,
		Name:   info.NickName,
		OpType: opType,
		Detail: fmt.Sprintf("%s card ID:%d, CardID:%s, ClientID:%d, Remark:%s", actionText, req.ID, card.CardID, card.ClientID, req.Remark),
		ObjId:  req.ID,
		Source: 1,
	})

	response.OkWithMessage(fmt.Sprintf("Card %s success", actionText), c)
}
