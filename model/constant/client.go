package constant

type ClientStatus uint

const (
	ClientStatus_Review ClientStatus = iota + 1
	ClientStatus_Active
	ClientStatus_Suspend
)

type ClientReviewStatus uint

const (
	ClientReviewStatus_UnReview ClientReviewStatus = iota + 1
	ClientReviewStatusStatus_Reviewing
	ClientReviewStatusStatus_Completed
)

type VerifyType string

const (
	VerifyType_BatchCardCancellPixel   VerifyType = "batchCardCancellPixel"
	VerifyType_RegularOpenCardPixel    VerifyType = "regularOpenCard"
	VerifyType_CVVQueryPixel           VerifyType = "cvvQueryPixel"
	VerifyType_CardWithdrawApplication VerifyType = "cardWithdrawApplication"
	VerifyType_CardRechargeApplication VerifyType = "cardRechargeApplication"
	VerifyType_UpdatePassword          VerifyType = "updatePassword"
	VerifyType_CodeConfiguration       VerifyType = "codeConfiguration"
)
