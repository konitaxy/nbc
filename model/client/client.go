package client

import (
	"time"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/common"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/model/system"
)

type Client struct {
	global.GVA_MODEL
	ClientNo     string                `gorm:"column:client_no;" json:"clientNo,omitempty" form:"clientNo,omitempty"`
	Email        string                `gorm:"column:email;type:varchar(128)" json:"email,omitempty" form:"email,omitempty"`
	Password     string                `gorm:"column:password;type:varchar(128)" json:"-" form:"-"`
	MarkName     string                `gorm:"column:mark_name;type:varchar(128)" json:"markName,omitempty" form:"markName,omitempty"`
	Name         string                `gorm:"column:name;type:varchar(128);index" json:"name,omitempty" form:"name,omitempty"`
	NickName     string                `gorm:"column:nick_name" json:"nickName,omitempty" form:"nickName,omitempty"`
	ClientType   string                `gorm:"column:client_type;type:varchar(64)" json:"clientType,omitempty" form:"clientType,omitempty"`
	Location     string                `gorm:"column:location;type:varchar(64)" json:"location,omitempty" form:"location,omitempty"`
	ClientStatus constant.ClientStatus `gorm:"column:client_status;default:1;index" json:"clientStatus,omitempty" form:"clientStatus,omitempty"`

	ClientReviewStatus constant.ClientReviewStatus `gorm:"column:client_review_status;default:1;index" json:"clientReviewStatus,omitempty" form:"clientReviewStatus,omitempty"`
	AccountManager     string                      `gorm:"column:account_manager" json:"accountManager" form:"accountManager"`
	Inviter            uint                        `gorm:"column:inviter" json:"inviter,omitempty" form:"inviter,omitempty"`
	InviteUser         system.SysUser              `gorm:"foreignKey:ID;references:Inviter" json:"inviteUser,omitempty" form:"inviteUser,omitempty"`
	RegisterSource     string                      `gorm:"column:register_source;type:varchar(64)" json:"registerSource,omitempty" form:"registerSource,omitempty"`
	ContactType        string                      `gorm:"column:contact_type" json:"contactType,omitempty" form:"contactType,omitempty"`
	Contact            string                      `gorm:"column:contact;" json:"contact,omitempty" form:"contact,omitempty"`
	ClientRegisterTime time.Time                   `gorm:"column:client_register_time;type:datetime" json:"clientRegisterTime,omitempty" form:"clientRegisterTime,omitempty"`
	Remark             string                      `gorm:"column:remark;type:text" json:"remark,omitempty" form:"remark,omitempty"`
	DueDiligence       ClientDueDiligence          `gorm:"foreignKey:ClientID;" json:"dueDiligence,omitempty" form:"dueDiligence,omitempty"`
	TOTPSecret         *string                     `gorm:"column:totp_secret;type:varchar(255)" json:"-" form:"-"`
	Bind2FA            bool                        `gorm:"column:bind_2fa;type:tinyint(1);default:0" json:"bind2FA" form:"bind2FA"`
	Wallet             *Wallet                     `gorm:"foreignKey:ClientID;" json:"wallet,omitempty" form:"wallet,omitempty"`
	VerifySetting      common.MapStringBool        `gorm:"column:verify_setting;type:json" json:"verifySetting,omitempty" form:"verifySetting,omitempty"`
	PIN                string                      `gorm:"column:pin;type:varchar(10)" json:"-" form:"-"`
	BindPin            bool                        `gorm:"column:bind_pin;type:tinyint(1);default:0" json:"bindPin" form:"bindPin"`
	IsTest         bool                 `gorm:"column:is_test;type:tinyint(1);default:0;index" json:"isTest" form:"isTest"`
	// MatrixAccount 光子易 Matrix 账户号（审核通过创建成功后保存）
	MatrixAccount string `gorm:"column:matrix_account;type:varchar(64);index" json:"matrixAccount,omitempty" form:"matrixAccount,omitempty"`
	// 子账号列表
}

// TableName 返回数据库表名
func (Client) TableName() string {
	return "clients"
}

type SessionLog struct {
	global.GVA_MODEL
	ClientID       uint      `gorm:"column:client_id;index" json:"clientId,omitempty" form:"clientId,omitempty"`
	LastActiveTime time.Time `gorm:"column:last_active_time;type:datetime;index" json:"lastActiveTime,omitempty" form:"lastActiveTime,omitempty"`
	IPAddress      string    `gorm:"column:ip_address;type:varchar(64)" json:"ipAddress,omitempty" form:"ipAddress,omitempty"`
	UserAgent      string    `gorm:"column:user_agent;type:varchar(255)" json:"userAgent,omitempty" form:"userAgent,omitempty"`
	Application    string    `gorm:"column:application;type:varchar(64)" json:"application,omitempty" form:"application,omitempty"`
	OpSystem       string    `gorm:"column:op_system;type:varchar(64)" json:"opSystem,omitempty" form:"opSystem,omitempty"`
	XToken         string    `gorm:"column:x_token;type:varchar(255);index" json:"-" form:"-"`
	Address        string    `gorm:"column:address;type:varchar(255)" json:"address,omitempty" form:"address,omitempty"`
	Status         bool      `gorm:"column:status" json:"status,omitempty" form:"status,omitempty"`
}

func (SessionLog) TableName() string {
	return "client_session_log"
}

type ClientDueDiligence struct {
	global.GVA_MODEL
	ClientID                    uint       `gorm:"column:client_id;index" json:"clientId,omitempty" form:"clientId,omitempty"`
	DDTimes                     uint       `gorm:"column:dd_times;default:0" json:"ddTimes,omitempty" form:"ddTimes,omitempty"`
	LastDDTime                  *time.Time `gorm:"column:dd_last_time;type:datetime" json:"ddLastTime,omitempty" form:"ddLastTime,omitempty"`
	EntEnterpriseType           *string    `gorm:"column:ent_enterprise_type;type:varchar(64)" json:"entEnterpriseType,omitempty" form:"entEnterpriseType,omitempty"`
	EntEnterpriseChineseName    string     `gorm:"column:ent_enterprise_chinese_name;type:varchar(128);not null" json:"entEnterpriseChineseName" form:"entEnterpriseChineseName"`
	EntEnterpriseEnglishName    *string    `gorm:"column:ent_enterprise_english_name;type:varchar(128)" json:"entEnterpriseEnglishName,omitempty" form:"entEnterpriseEnglishName,omitempty"`
	EntBusinessRegistrationForm *string    `gorm:"column:ent_business_registration_form;type:MEDIUMTEXT" json:"entBusinessRegistrationForm,omitempty" form:"entBusinessRegistrationForm,omitempty"`
	EntBusinessRegistrationNo   string     `gorm:"column:ent_business_registration_no;type:varchar(64);not null" json:"entBusinessRegistrationNo" form:"entBusinessRegistrationNo"`
	EntBusinessAddressProof     *string    `gorm:"column:ent_business_address_proof;type:varchar(255)" json:"entBusinessAddressProof,omitempty" form:"entBusinessAddressProof,omitempty"`
	EntDateOfEstablishment      string     `gorm:"column:ent_date_of_establishment" json:"entDateOfEstablishment,omitempty" form:"entDateOfEstablishment,omitempty"`
	EntDateOfExpiration         string     `gorm:"column:ent_date_of_expiration" json:"entDateOfExpiration,omitempty" form:"entDateOfExpiration,omitempty"`
	EntLocalBusinessPremise     string     `gorm:"column:ent_local_business_premise;type:varchar(255);not null" json:"entLocalBusinessPremise" form:"entLocalBusinessPremise"`
	EntProvince                 *string    `gorm:"column:ent_province;type:varchar(64)" json:"entProvince,omitempty" form:"entProvince,omitempty"`
	EntCity                     *string    `gorm:"column:ent_city;type:varchar(64)" json:"entCity,omitempty" form:"entCity,omitempty"`
	EntListedCompany            *bool      `gorm:"column:ent_listed_company;type:boolean" json:"entListedCompany,omitempty" form:"entListedCompany,omitempty"`
	EntStateOwned               *bool      `gorm:"column:ent_state_owned;type:boolean" json:"entStateOwned,omitempty" form:"entStateOwned,omitempty"`
	EntForeignInvested          *bool      `gorm:"column:ent_foreign_invested;type:boolean" json:"entForeignInvested,omitempty" form:"entForeignInvested,omitempty"`
	EntShareholderStructure     *string    `gorm:"column:ent_shareholder_structure;type:varchar(255)" json:"entShareholderStructure,omitempty" form:"entShareholderStructure,omitempty"`
	EntIsB2B                    *bool      `gorm:"column:ent_is_b2b;type:boolean" json:"entIsB2B,omitempty" form:"entIsB2B,omitempty"`
	EntOperationAddress         *string    `gorm:"column:ent_operation_address;type:varchar(255)" json:"entOperationAddress,omitempty" form:"entOperationAddress,omitempty"`
	EntRegisteredCapital        *string    `gorm:"column:ent_registered_capital;type:varchar(64)" json:"entRegisteredCapital,omitempty" form:"entRegisteredCapital,omitempty"`
	EntIntendedBusinessIndustry *string    `gorm:"column:ent_intended_business_industry;type:varchar(128)" json:"entIntendedBusinessIndustry,omitempty" form:"entIntendedBusinessIndustry,omitempty"`
	IndCountryOrRegion          string     `gorm:"column:ind_country_or_region;type:varchar(64);not null" json:"indCountryOrRegion" form:"indCountryOrRegion"`
	IndPosition                 *string    `gorm:"column:ind_position;type:varchar(64)" json:"indPosition,omitempty" form:"indPosition,omitempty"`
	IndIDType                   string     `gorm:"column:ind_id_type;type:varchar(16);not null" json:"indIDType" form:"indIDType"`
	IndChineseName              string     `gorm:"column:ind_chinese_name;type:varchar(128);not null" json:"indChineseName" form:"indChineseName"`
	IndEnglishName              *string    `gorm:"column:ind_english_name;type:varchar(128)" json:"indEnglishName,omitempty" form:"indEnglishName,omitempty"`
	IndIDFrontEnd               string     `gorm:"column:ind_id_front_end;type:MEDIUMTEXT" json:"indIDFrontEnd" form:"indIDFrontEnd"`
	IndIDBackEnd                string     `gorm:"column:ind_id_back_end;type:MEDIUMTEXT;" json:"indIDBackEnd" form:"indIDBackEnd"`
	IndIdentificationNo         string     `gorm:"column:ind_identification_no;type:varchar(64);not null" json:"indIdentificationNo" form:"indIdentificationNo"`
	IndIssueDate                string     `gorm:"column:ind_issue_date" json:"indIssueDate,omitempty" form:"indIssueDate,omitempty"`
	IndExpirationDate           string     `gorm:"column:ind_expiration_date" json:"indExpirationDate,omitempty" form:"indExpirationDate,omitempty"`
	IndDateOfBirth              string     `gorm:"column:ind_date_of_birth" json:"indDateOfBirth,omitempty" form:"indDateOfBirth,omitempty"`
	IndProvinceOrState          *string    `gorm:"column:ind_province_or_state;type:varchar(64)" json:"indProvinceOrState,omitempty" form:"indProvinceOrState,omitempty"`
	IndCity                     *string    `gorm:"column:ind_city;type:varchar(64)" json:"indCity,omitempty" form:"indCity,omitempty"`
	IndResidentialAddress       *string    `gorm:"column:ind_residential_address;type:varchar(255)" json:"indResidentialAddress,omitempty" form:"indResidentialAddress,omitempty"`
	IndReliabilityOfDocuments   *string    `gorm:"column:ind_reliability_of_documents;type:varchar(128)" json:"indReliabilityOfDocuments,omitempty" form:"indReliabilityOfDocuments,omitempty"`
	Tip                         string     `gorm:"column:tip;type:varchar(255)" json:"tip" form:"tip"`
	NeedEnhancedKYB             bool       `gorm:"column:need_enhanced_kyb;type:boolean" json:"needEnhancedKYB" form:"needEnhancedKYB"`
	Type                        string     `gorm:"-" json:"type" form:"type"`
}

// TableName 返回数据库表名
func (ClientDueDiligence) TableName() string {
	return "client_due_diligence"
}

type VerifySetting struct {
	Key   string
	Value bool
	Path  string
	Level uint /// 1: pin 2: otp or mail
}

var Default_Verify_Setting = []VerifySetting{
	{Key: "verifySetting", Path: "POST:/client/verifySetting", Value: true},
	{Key: "iamLogin", Path: "POST:/client/iam/login", Value: true},
	// {Key: "tocp", Path: "POST:/client/tocp", Value: false},

	{Key: "disableTocp", Path: "DELETE:/client/tocp", Value: true},
	{Key: "iamAccountCreate", Path: "POST:/client/iam/create", Value: true},
	// {Key: "iamAccountUpdate", Path: "POST:/client/iam/update", Value: true},
	{Key: "pin", Path: "POST:/client/pin", Value: true},
	{Key: "login", Path: "POST:/client/login", Value: true},
	{Key: "register", Path: "POST:/client/register", Value: true},
	{Key: "changePassword", Path: "POST:/client/changePassword", Value: true},
	{Key: "resetPassword", Path: "POST:/client/resetPassword", Value: true},

	{Key: "cardCancel", Path: "POST:/card/cancel", Value: false, Level: 1},
	{Key: "cardDetail", Path: "GET:/card/detail", Value: false, Level: 1},
	{Key: "cardAdd", Path: "POST:/card/add", Value: false, Level: 1},
	{Key: "cardcvv", Path: "GET:/card/cvv", Value: false, Level: 1},
	{Key: "cardWithdraw", Path: "POST:/card/withdraw", Value: false, Level: 1},
	{Key: "cardRecharge", Path: "POST:/card/recharge", Value: false, Level: 1},
	{Key: "cardAdjustLimit", Path: "POST:/card/adjustLimit", Value: false, Level: 1},
	{Key: "walletWithdraw", Path: "POST:/wallet/withdraw/apply", Value: false, Level: 1},
}

var Default_IAM_Verify_Setting = []VerifySetting{
	{Key: "verifySetting", Path: "POST:/client/verifySetting", Value: true},
	{Key: "iamLogin", Path: "POST:/client/iam/login", Value: true},
	// {Key: "tocp", Path: "POST:/client/tocp", Value: false},

	{Key: "disableTocp", Path: "DELETE:/client/tocp", Value: true},
	{Key: "iamAccountCreate", Path: "POST:/client/iam/create", Value: true},
	// {Key: "iamAccountUpdate", Path: "POST:/client/iam/update", Value: true},
	{Key: "pin", Path: "POST:/client/pin", Value: true},
	{Key: "changePassword", Path: "POST:/client/changePassword", Value: true},
	{Key: "resetPassword", Path: "POST:/client/resetPassword", Value: true},

	{Key: "cardCancel", Path: "POST:/card/cancel", Value: false, Level: 1},
	{Key: "cardDetail", Path: "GET:/card/detail", Value: false, Level: 1},
	{Key: "cardAdd", Path: "POST:/card/add", Value: false, Level: 1},
	{Key: "cardcvv", Path: "GET:/card/cvv", Value: false, Level: 1},
	{Key: "cardWithdraw", Path: "POST:/card/withdraw", Value: false, Level: 1},
	{Key: "cardRecharge", Path: "POST:/card/recharge", Value: false, Level: 1},
	{Key: "cardAdjustLimit", Path: "POST:/card/adjustLimit", Value: false, Level: 1},
	{Key: "walletWithdraw", Path: "POST:/wallet/withdraw/apply", Value: false, Level: 1},
}

// IAMUser 子账号模型
type IAMUser struct {
	global.GVA_MODEL
	ClientID uint             `json:"clientId" gorm:"index;comment:所属主账号ID"`
	Client   *Client          `json:"client,omitempty" gorm:"foreignKey:ClientID"` // 关联主账号
	Email    string           `json:"email" gorm:"type:varchar(128);uniqueIndex;comment:邮箱(登录名)"`
	Password string           `json:"-" gorm:"type:varchar(128);comment:密码"`
	Nickname string           `json:"nickname" gorm:"type:varchar(64);comment:昵称"`
	Status   int8             `json:"status" gorm:"default:1;comment:状态 1启用 0禁用"`
	Roles    common.SliceUint `json:"roles" gorm:"type:json;comment:角色ID列表"`
	Wallet   *Wallet          `json:"wallet,omitempty" gorm:"foreignKey:ClientID;references:ClientID"`

	// Security settings
	PIN           string               `gorm:"column:pin;type:varchar(10)" json:"-" form:"-"`
	BindPin       bool                 `gorm:"column:bind_pin;type:tinyint(1);default:0" json:"bindPin" form:"bindPin"`
	TOTPSecret    *string              `gorm:"column:totp_secret;type:varchar(255)" json:"-" form:"-"`
	Bind2FA       bool                 `gorm:"column:bind_2fa;type:tinyint(1);default:0" json:"bind2FA" form:"bind2FA"`
	VerifySetting common.MapStringBool `gorm:"column:verify_setting;type:json" json:"verifySetting,omitempty" form:"verifySetting,omitempty"`
}

func (IAMUser) TableName() string {
	return "iam_users"
}

// IAMRole 角色模型
type IAMRole struct {
	global.GVA_MODEL
	ClientID    uint            `json:"clientId" gorm:"index;comment:所属主账号ID(0为系统默认角色)"`
	RoleCode    string          `json:"roleCode" gorm:"type:varchar(64);index;comment:角色代码"`
	RoleName    string          `json:"roleName" gorm:"type:varchar(64);comment:角色名"`
	Description string          `json:"description" gorm:"type:varchar(255);comment:描述"`
	IsDefault   bool            `json:"isDefault" gorm:"default:false;comment:是否系统默认角色"`
	Sort        int             `json:"sort" gorm:"default:0;comment:排序"`
	Permissions []IAMPermission `json:"permissions" gorm:"many2many:iam_role_permissions;"`
}

func (IAMRole) TableName() string {
	return "iam_roles"
}

// IAMPermission 权限点模型
type IAMPermission struct {
	global.GVA_MODEL
	Code        string `json:"code" gorm:"type:varchar(64);uniqueIndex;comment:权限代码"`
	Name        string `json:"name" gorm:"type:varchar(64);comment:权限名称"`
	Type        string `json:"type" gorm:"type:varchar(16);index;comment:类型 page/api"`
	Action      string `json:"action" gorm:"type:varchar(128);comment:权限标识,如 card:create"`
	Path        string `json:"path" gorm:"type:varchar(128);comment:API路径"`
	Method      string `json:"method" gorm:"type:varchar(16);comment:HTTP方法"`
	ParentID    uint   `json:"parentId" gorm:"index;default:0;comment:父级ID"`
	Sort        int    `json:"sort" gorm:"default:0;comment:排序"`
	Description string `json:"description" gorm:"type:varchar(255);comment:描述"`
}

func (IAMPermission) TableName() string {
	return "iam_permissions"
}
