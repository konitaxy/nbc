package utils

var (
	IdVerify       = Rules{"ID": {NotEmpty()}}
	ApiVerify      = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify     = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify = Rules{"Title": {NotEmpty()}}
	LoginVerify    = Rules{"CaptchaId": {NotEmpty()}, "Captcha": {NotEmpty()}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}

	RegisterVerify = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}

	UserRegisterVerify  = Rules{"Email": {NotEmpty()}, "Password": {NotEmpty()}, "InviteCode": {NotEmpty()}}
	ApplyRegisterVerify = Rules{"Email": {NotEmpty()}, "Name": {NotEmpty()}, "Country": {NotEmpty()}, "Bio": {NotEmpty()}}
	ProductVerify       = Rules{"ProductID": {NotEmpty()}, "ArtworkID": {NotEmpty()}, "SizeList": {NotEmpty()}, "ProductTypes": {NotEmpty()}, "Layout": {NotEmpty()}}
	CreateProductVerify = Rules{"ArtworkID": {NotEmpty()}, "SizeList": {NotEmpty()}, "ProductTypes": {NotEmpty()}, "Layout": {NotEmpty()}}

	PageInfoVerify         = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify         = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify         = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}, "Fields": {NotEmpty()}}
	AutoPackageVerify      = Rules{"PackageName": {NotEmpty()}}
	AuthorityVerify        = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}, "ParentId": {NotEmpty()}}
	AuthorityIdVerify      = Rules{"AuthorityId": {NotEmpty()}}
	CollectionVerify       = Rules{"Title": {NotEmpty()}, "Description": {NotEmpty()}, "UserID": {NotEmpty()}}
	OldAuthorityVerify     = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify   = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify = Rules{"AuthorityId": {NotEmpty()}}
	FeeUserConfigVerify    = Rules{"ClientNo": {NotEmpty()}, "Fee": {NotEmpty()}, "CalType": {NotEmpty()}, "FeeType": {NotEmpty()}}
	FeeGlobalConfigVerify  = Rules{"CardBin": {NotEmpty()}, "Fee": {NotEmpty()}, "CalType": {NotEmpty()}, "FeeType": {NotEmpty()}}
	AddCardBinVerify       = Rules{"CardBin": {NotEmpty()}, "CardBrand": {NotEmpty()}, "Currency": {NotEmpty()}, "Region": {NotEmpty()}, "CardModel": {NotEmpty()}, "Channel": {NotEmpty()}}
	CancelCardBinVerify      = Rules{"ID": {NotEmpty()}, "CardId": {NotEmpty()}}
	BatchCancelCardBinVerify = Rules{"List": {Gt("0"), Le("100")}}
	ConfigVerify           = Rules{"Kay": {NotEmpty()}, "ValueType": {NotEmpty()}}
	WalletWithdrawVerify   = Rules{"AccountType": {NotEmpty()}, "AccountNumber": {NotEmpty()}, "Currency": {NotEmpty()}, "Amount": {NotEmpty()}, "Memo": {NotEmpty()}}
	RechargeApply          = Rules{"Amount": {NotEmpty()}, "Currency": {NotEmpty()}}

	AddCardHolderVerify = Rules{"Region": {NotEmpty()}, "FirstName": {NotEmpty()}, "LastName": {NotEmpty()}, "Email": {NotEmpty()}, "MobilePrefix": {NotEmpty()}, "Mobile": {NotEmpty()}, "BirthDate": {NotEmpty()}, "CountryCode": {NotEmpty()}, "State": {NotEmpty()}, "City": {NotEmpty()}, "Postcode": {NotEmpty()}, "Address": {NotEmpty()}}

	UpdateCardHolderVerify = Rules{"CardHolderId": {NotEmpty()}}

	OpenCardVerify  = Rules{"CardBin": {NotEmpty()}, "CardBinId": {NotEmpty()}, "Amount": {NotEmpty()}, "CardType": {NotEmpty()}}
	CardGroupVerify = Rules{"Name": {NotEmpty()}}
)
