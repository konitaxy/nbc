package client

import (
	"context"
	"fmt"
	"os/user"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/client/request"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/system"
	"gitlab.com/ucard/utils/captcha"
	"gorm.io/gorm"
)

type ClientService struct{}

var iamSvc = IAMService{}

func (*ClientService) Login(er *client.Client) error {
	err := global.GVA_DB.Preload("Wallet").Where("email = ? and password = ?", er.Email, er.Password).First(&er).Error
	return err
}

func (*ClientService) GetClientByMail(email string) (client client.Client, err error) {
	err = global.GVA_DB.Where("email = ?", email).First(&client).Error
	return
}
func (p *ClientService) ResetPassword(req request.ChangePasswordReq) error {
	if er, err := p.GetClientByMail(req.Email); err == nil {
		return global.GVA_DB.Model(&client.Client{}).Where("id = ?", er.ID).Update("password", req.NewPassword).Error
	} else {
		return err
	}
}

func (*ClientService) ListLastSessionLog(uid uint) (list []client.SessionLog, err error) {
	err = global.GVA_DB.Where("client_id = ?", uid).Order("last_active_time desc").Limit(10).Find(&list).Error
	return
}
func (p *ClientService) AddSessionLog(log client.SessionLog) (err error) {
	log.XToken = log.XToken[len(log.XToken)-50:]
	err = global.GVA_DB.Save(&log).Error
	if err != nil {
		return
	} else {
		err = global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("session_active_%d", log.ClientID), "1", time.Minute).Err()
	}
	return
}
func (p *ClientService) UpdateSessionLog(token, newToken string, uid uint) (err error) {
	if newToken != "" {
		newToken = newToken[len(newToken)-50:]
		token = token[len(token)-50:]
		err = global.GVA_DB.Model(&client.SessionLog{}).Where("x_token = ?", token).UpdateColumn("x_token", newToken).Error

	} else {
		if val := global.GVA_REDIS.Get(context.Background(), fmt.Sprintf("session_active_%d", uid)).Val(); val == "" {
			token = token[len(token)-50:]
			err = global.GVA_DB.Model(&client.SessionLog{}).Where("x_token = ?", token).UpdateColumn("last_active_time", time.Now()).Error
		}
	}

	return
}
func (p *ClientService) ChangePassword(userId uint, req request.ChangePasswordReq) error {
	if er, err := p.GetClient(userId); err == nil || er.ID == 0 {
		// if er.Password != req.Password {
		// 	return errors.New("password error")
		// } else {
		return global.GVA_DB.Model(&user.User{}).Where("id = ?", userId).Update("password", req.NewPassword).Error
		// }
	} else {
		return err
	}
}
func (*ClientService) SaveSmsCode(code common.SmsCode) error {
	err := global.GVA_DB.Save(&code).Error
	return err
}
func (*ClientService) GetBalance(id uint) (wallet client.Wallet, err error) {

	err = global.GVA_DB.First(&wallet, "client_id = ?", id).Error
	return
}

func (*ClientService) GetClient(uid uint) (client client.Client, err error) {
	err = global.GVA_DB.Preload("Wallet").Where("id = ?", uid).First(&client).Error
	return
}

func (*ClientService) NeedVerifySetting(email string, req string) (bool, uint) {
	level := uint(2)
	for _, v := range client.Default_Verify_Setting {
		if v.Path == req {
			if v.Value {
				return true, v.Level
			} else {
				level = v.Level
			}
			break
		}
	}
	v := global.GVA_REDIS.HGet(context.TODO(), fmt.Sprintf("verify_setting_%s", email), req).Val()
	return v != "" && v != "0", level
}

// NeedIAMVerifySetting checks if verification is needed for IAM users
func (*ClientService) NeedIAMVerifySetting(email string, req string) (bool, uint) {
	level := uint(2)
	for _, v := range client.Default_IAM_Verify_Setting {
		if v.Path == req {
			if v.Value {
				return true, v.Level
			} else {
				level = v.Level
			}
			break
		}
	}
	v := global.GVA_REDIS.HGet(context.TODO(), fmt.Sprintf("iam_verify_setting_%s", email), req).Val()
	return v != "" && v != "0", level
}
func (*ClientService) GetClientDueDiligence(uid uint) (dd client.ClientDueDiligence, err error) {
	err = global.GVA_DB.Where("client_id = ?", uid).First(&dd).Error
	return
}

func (*ClientService) SetClientDueDiligence(dd *client.ClientDueDiligence) (err error) {
	err = global.GVA_DB.Save(dd).Error
	return
}

func (*ClientService) Save(er *client.Client) error {
	return global.GVA_DB.Save(er).Error
}
func (*ClientService) Create(er *client.Client) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(er).Error; err != nil {
			return err
		}
		if err := tx.Create(&client.Wallet{
			ClientID: er.ID,
			Balance:  decimal.NewFromFloat(0),
			Currency: constant.USD,
		}).Error; err != nil {
			return err
		}
		return nil
	})
}

var consoleMenu = make([]system.SysMenu, 0)

func init() {
	menus := []system.SysMenu{
		{
			MenuId:   "0",
			Children: []system.SysMenu{
				// {
				// 	MenuId:     "11",
				// 	Children:   nil,
				// 	Parameters: nil,
				// 	SysBaseMenu: system.SysBaseMenu{
				// 		Hidden:    false,
				// 		Path:      "/analytics",
				// 		Name:      "analytics",
				// 		Component: "view/pixal/wallet/my/index.vue",
				// 		Sort:      1,
				// 		Meta:      system.Meta{Title: "Wallet", Icon: "", BootstrapIcon: "fa-solid fa-cloud-arrow-up", CloseTab: false, DefaultMenu: true},
				// 	},
				// },
			},
			Parameters: nil,
			SysBaseMenu: system.SysBaseMenu{
				Hidden:    false,
				Path:      "/dashboard",
				Name:      "dashboard",
				Component: "view/pixal/dashboard/index.vue",
				Sort:      0,
				Meta:      system.Meta{Title: "dashboard", Icon: "dashboard", BootstrapIcon: "bi bi-speedometer2", CloseTab: false, DefaultMenu: false},
			},
		},
		{
			MenuId: "1",
			Children: []system.SysMenu{
				{
					MenuId:     "11",
					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "my",
						Name:      "myWallet",
						Component: "view/pixal/wallet/my/index.vue",
						Sort:      1,
						Meta:      system.Meta{Title: "Wallet", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: true},
					},
				},
				{
					MenuId:     "12",
					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "rechargeManage",
						Name:      "walletRechargeManager",
						Component: "view/pixal/wallet/history.vue",
						Sort:      2,
						Meta:      system.Meta{Title: "transaction_details", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId:     "13",
					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "transferManager",
						Name:      "walletTransferManager",
						Component: "view/pixal/wallet/rechargeAndWithdraw.vue",
						Sort:      2,
						Meta:      system.Meta{Title: "transfer_management", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
			},
			Parameters: nil,
			SysBaseMenu: system.SysBaseMenu{
				Hidden:    false,
				Path:      "/wallet",
				Name:      "wallet",
				Component: "view/pixal/wallet/index.vue",
				Sort:      3,
				Meta:      system.Meta{Title: "Wallet", Icon: "user", BootstrapIcon: "bi bi-wallet", CloseTab: false, DefaultMenu: true},
			},
		},
		{
			MenuId: "2",

			Children: []system.SysMenu{
				{
					MenuId: "21",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "rechargeCard",
						Name:      "rechargeCard",
						Component: "view/pixal/card/rechargeCard.vue",
						Sort:      1,
						Meta:      system.Meta{Title: "recharge_card", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId: "22",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "shareCard",
						Name:      "shareCard",
						Component: "view/pixal/card/sharedCard.vue",
						Sort:      1,
						Meta:      system.Meta{Title: "share_card", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId: "23",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    true,
						Path:      "activeCardRecord",
						Name:      "activeCardRecord",
						Component: "view/pixal/card/activeCardRecord.vue",
						Sort:      2,
						Meta:      system.Meta{Title: "activate_card_history", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId: "24",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "transactionDetail",
						Name:      "transactionDetail",
						Component: "view/pixal/card/transactionDetails.vue",
						Sort:      3,
						Meta:      system.Meta{Title: "transaction_details", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId: "25",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    true,
						Path:      "withdrawAndRefund",
						Name:      "withdrawAndRefund",
						Component: "view/pixal/card/withdrawAndRefund.vue",
						Sort:      4,
						Meta:      system.Meta{Title: "withdraw_and_refund", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				}, {
					MenuId: "26",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "holder",
						Name:      "cardHolder",
						Component: "view/pixal/card/cardHolder.vue",
						Sort:      4,
						Meta:      system.Meta{Title: "card_holder", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				}, {
					MenuId: "27",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "group",
						Name:      "cardGroup",
						Component: "view/pixal/card/cardGroup.vue",
						Sort:      4,
						Meta:      system.Meta{Title: "card_group", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
			},
			Parameters: nil,
			SysBaseMenu: system.SysBaseMenu{
				Hidden:    false,
				Path:      "/card",
				Name:      "card",
				Component: "view/pixal/card/index.vue",
				Sort:      2,
				Meta:      system.Meta{Title: "pixel_card", Icon: "", BootstrapIcon: "bi bi-credit-card", CloseTab: false, DefaultMenu: false},
			},
		},
		{
			MenuId: "3",
			Children: []system.SysMenu{
				{
					MenuId: "31",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "info",
						Name:      "accountInfo",
						Component: "view/pixal/account/info/index.vue",
						Sort:      1,
						Meta:      system.Meta{Title: "account_information", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId: "32",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "identityVerify",
						Name:      "identityVerify",
						Component: "view/pixal/account/identityVerify/index.vue",
						Sort:      2,
						Meta:      system.Meta{Title: "identity_verification", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId:     "33",
					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "googleVerify",
						Name:      "googleVerify",
						Component: "view/pixal/account/googleVerify/index.vue",
						Sort:      3,
						Meta:      system.Meta{Title: "MFA", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId: "34",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    true,
						Path:      "feeDetail",
						Name:      "feeDetail",
						Component: "view/pixal/account/fees/index.vue",
						Sort:      4,
						Meta:      system.Meta{Title: "Fee Detail", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
				{
					MenuId: "35",

					Children:   nil,
					Parameters: nil,
					SysBaseMenu: system.SysBaseMenu{
						Hidden:    false,
						Path:      "verifySetting",
						Name:      "verifySetting",
						Component: "view/pixal/account/verifySetting/index.vue",
						Sort:      5,
						Meta:      system.Meta{Title: "valiate_code", Icon: "", BootstrapIcon: "", CloseTab: false, DefaultMenu: false},
					},
				},
			},
			Parameters: nil,
			SysBaseMenu: system.SysBaseMenu{
				Hidden:    false,
				Path:      "/account",
				Name:      "account",
				Component: "view/pixal/account/index.vue",
				Sort:      3,
				Meta:      system.Meta{Title: "account", Icon: "", BootstrapIcon: "bi bi-person-fill-gear", CloseTab: false, DefaultMenu: false},
			},
		},
	}
	consoleMenu = append(consoleMenu, menus...)
}

// 菜单ID与角色ID的映射关系
// 角色ID: 1-基础功能, 2-卡片查看, 3-卡片管理, 4-敏感操作, 5-财务操作, 6-数据导出
var menuRoleMap = map[string][]uint{
	"0":  {},           // dashboard - 基础功能
	"1":  {1},          // wallet - 基础功能、财务操作
	"11": {1},          // myWallet - 基础功能
	"12": {7},          // walletRechargeManager - 财务操作
	"2":  {2, 3, 4, 5}, // card - 卡片相关角色
	"21": {2, 3},       // pixelCard - 卡片查看、管理
	"22": {2},          // activeCardRecord - 卡片查看
	"23": {2},          // transactionDetail - 卡片查看
	"24": {3},          // cardHolder - 卡片管理
	"25": {3},          // cardGroup - 卡片管理
	"26": {3},          // cardGroup - 卡片管理

	"3":  {1}, // account - 基础功能
	"31": {},  // accountInfo - 基础功能
	"32": {},  // identityVerify - 基础功能
	"33": {1}, // googleVerify - 基础功能
	"34": {5}, // feeDetail - 财务操作
	"35": {1}, // verifySetting - 基础功能
	"4":  {},  // iam - 仅主账号可见，IAM用户不显示
}

func (p *ClientService) GetMenuTree(roles []uint, isTest bool) []system.SysMenu {
	// 如果 roles 为空，返回全部菜单
	if len(roles) == 0 {
		if isTest {
			return append(consoleMenu, system.SysMenu{
				MenuId:     "4",
				Children:   []system.SysMenu{},
				Parameters: nil,
				SysBaseMenu: system.SysBaseMenu{
					Hidden:    false,
					Path:      "/iam",
					Name:      "iam",
					Component: "view/pixal/iam/index.vue",
					Sort:      0,
					Meta:      system.Meta{Title: "IAM", Icon: "iam", BootstrapIcon: "bi bi-person-fill-lock", CloseTab: false, DefaultMenu: false},
				},
			})
		}
		return consoleMenu
	}

	// 将角色列表转为 map 方便查找
	roleSet := make(map[uint]bool)
	for _, r := range roles {
		roleSet[r] = true
	}

	// 过滤菜单
	return filterMenus(consoleMenu, roleSet)
}

// filterMenus 递归过滤菜单
func filterMenus(menus []system.SysMenu, roleSet map[uint]bool) []system.SysMenu {
	result := make([]system.SysMenu, 0)

	for _, menu := range menus {
		// 检查当前菜单是否有权限
		if hasMenuPermission(menu.MenuId, roleSet) {
			// 复制菜单，避免修改原始数据
			newMenu := menu
			// 递归过滤子菜单
			if len(menu.Children) > 0 {
				newMenu.Children = filterMenus(menu.Children, roleSet)
			}
			// 如果有权限且（没有子菜单或子菜单过滤后不为空），添加到结果
			if len(menu.Children) == 0 || len(newMenu.Children) > 0 {
				result = append(result, newMenu)
			}
		}
	}

	return result
}

// hasMenuPermission 检查是否有菜单权限
func hasMenuPermission(menuId string, roleSet map[uint]bool) bool {
	allowedRoles, exists := menuRoleMap[menuId]
	// 如果菜单没有配置角色限制，默认不显示
	if !exists {
		return false
	}
	// 如果配置为空数组，表示仅主账号可见
	if len(allowedRoles) == 0 {
		return false
	}
	// 检查用户是否拥有任一允许的角色
	for _, r := range allowedRoles {
		if roleSet[r] {
			return true
		}
	}
	return false
}

func (p *ClientService) CodeVerify(email, code, req string, level uint, isIAM bool) bool {
	// Find verify key from path based on user type
	var verifyKey string
	verifySettings := client.Default_Verify_Setting
	if isIAM {
		verifySettings = client.Default_IAM_Verify_Setting
	}

	for _, v := range verifySettings {
		if v.Path == req {
			verifyKey = v.Key
			break
		}
	}
	if verifyKey == "" {
		return false
	}

	email = strings.ToLower(email)
	redisKey := fmt.Sprintf("verify_code_%s_%s", verifyKey, email)

	// IAM user verification
	if isIAM {
		iamUser, err := iamSvc.GetIAMUserByEmail(email)
		if err != nil || iamUser.ID == 0 {
			// IAM user not found, verify with email code directly
			return global.GVA_REDIS.Get(context.Background(), redisKey).Val() == code
		}

		// PIN verification (level 1 only)
		if level == 1 && iamUser.BindPin && code == iamUser.PIN {
			global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("verify_pin_%s", iamUser.Email), "1", 10*time.Minute)
			return true
		}

		// 2FA verification
		if iamUser.Bind2FA && iamUser.TOTPSecret != nil {
			return captcha.VerifyTOTP(*iamUser.TOTPSecret, code)
		}

		// Email code verification
		return global.GVA_REDIS.Get(context.Background(), redisKey).Val() == code
	}

	// Main account verification
	cl, _ := p.GetClientByMail(email)
	if cl.ID == 0 {
		// Client not found, verify with email code directly
		return global.GVA_REDIS.Get(context.Background(), redisKey).Val() == code
	}

	// PIN verification (level 1 only)
	if level == 1 && cl.BindPin && code == cl.PIN {
		global.GVA_REDIS.Set(context.Background(), fmt.Sprintf("verify_pin_%s", cl.Email), "1", 10*time.Minute)
		return true
	}

	// 2FA verification
	if cl.Bind2FA && cl.TOTPSecret != nil {
		return captcha.VerifyTOTP(*cl.TOTPSecret, code)
	}

	// Email code verification
	return global.GVA_REDIS.Get(context.Background(), redisKey).Val() == code
}
