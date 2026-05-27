package gzy

// PhotonPay OpenAPI v4：VCC 为 /vcc/openApi/v4；钱包为 /wallet/openApi/v4。
// 若与线网文档不一致，只需改本文件常量或配置 gzy.token-path。
// 文档: https://api-doc.photonpay.com/

const (
	// 默认 POST /oauth2/token/accessToken（Authorization: Basic base64(appId+"/"+appSecret)）。
	// 若需旧版 client_credentials 表单 POST，设置 gzy.token-path 为 /oauth2/token。
	pathOAuthAccessToken     = "/oauth2/token/accessToken"
	pathOAuthTokenFormLegacy = "/oauth2/token"

	pathVccOpenAPIv2 = "/vcc/open/v2"
	pathVccOpenAPIv4 = "/vcc/openApi/v4"

	// pathSandBoxTransaction POST 沙箱交易模拟（UAT: /vcc/open/v2/sandBoxTransaction）
	pathSandBoxTransaction = pathVccOpenAPIv2 + "/sandBoxTransaction"

	pathWalletOpenAPIv4 = "/wallet/openApi/v4"
	pathAccountSingle   = pathWalletOpenAPIv4 + "/account/single" // GET 实时金额（单账户）

	// DefaultGzyAccountID 光子易默认钱包账户 ID（openCard / preRecharge 等未配置 gzy.account-id 时使用）。
	DefaultGzyAccountID = "FA-USD2052566705788575744"

	pathOpenCard = pathVccOpenAPIv4 + "/openCard"

	pathBalanceHistory = pathVccOpenAPIv4 + "/queryBalanceHistory"

	pathAddCardholder       = pathVccOpenAPIv4 + "/addCardholder" // POST 添加用卡人
	pathCardholderDetail    = pathVccOpenAPIv4 + "/queryCardholderDetail"
	pathPagingVccCardholder = pathVccOpenAPIv4 + "/pagingVccCardholder" // GET 用卡人分页列表

	pathFreezeCard = pathVccOpenAPIv4 + "/freezeCard"

	// GET /vcc/openApi/v4/getCardDetail?cardId=（Photon 卡详情；旧路径 queryCardDetail 已废弃）
	pathGetCardDetail = pathVccOpenAPIv4 + "/getCardDetail"
	// GET /vcc/openApi/v4/getCvv?cardId= 查询卡 CVV
	pathGetCvv = pathVccOpenAPIv4 + "/getCvv"
	// GET /vcc/openApi/v4/getCardBin 可用卡 BIN 列表
	pathGetCardBin = pathVccOpenAPIv4 + "/getCardBin"
	pathCardOpChangeLimit = pathVccOpenAPIv4 + "/changeSubAuthLimit"
	pathPreRecharge       = pathVccOpenAPIv4 + "/preRecharge"    // GET 换汇询价
	pathRecharge          = pathVccOpenAPIv4 + "/recharge"       // POST 转入下单（须先 preRecharge）
	pathRechargeReturn    = pathVccOpenAPIv4 + "/rechargeReturn" // 卡金额退还
	pathCardOpCancel      = pathVccOpenAPIv4 + "/cancelCard"

	pathPagingVccTradeOrder = pathVccOpenAPIv4 + "/pagingVccTradeOrder" // 卡交易明细分页

	// 以下路径仍被 GetInboundDetail 使用（与交易分页不同），待 Photon 入账详情文档确认后再改。
	pathQueryCardTransactions = pathVccOpenAPIv4 + "/queryCardTransactions"

	pathFundInboundApply = pathVccOpenAPIv4 + "/applyInbound"
	pathFundInboundList  = pathVccOpenAPIv4 + "/listInbound"
)
