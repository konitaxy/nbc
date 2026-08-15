package client

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/client/request"
	"gitlab.com/ucard/model/constant"
	"gitlab.com/ucard/utils"
	"gorm.io/gorm"
)

type IAMService struct{}

const IAMVerifyCodePrefix = "iam_verify_code_"

// ParentAccountFrozen 父账号是否被冻结
func (s *IAMService) ParentAccountFrozen(clientID uint) (bool, error) {
	if clientID == 0 {
		return false, nil
	}
	var parent client.Client
	if err := global.GVA_DB.Select("id, client_status").First(&parent, clientID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, errors.New("parent_account_not_found")
		}
		return false, err
	}
	return parent.ClientStatus == constant.ClientStatus_Suspend, nil
}

// PreLogin Verify IAM user credentials without returning full user info
func (s *IAMService) PreLogin(email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	passwordHash := utils.MD5V([]byte(password))

	var user client.IAMUser
	if err := global.GVA_DB.Select("id, status, client_id").
		Where("email = ? AND password = ?", email, passwordHash).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid_email_or_password")
		}
		return err
	}

	if user.Status == 0 {
		return errors.New("account_disabled")
	}
	if user.Status == 2 {
		return errors.New("account_frozen")
	}
	if frozen, err := s.ParentAccountFrozen(user.ClientID); err != nil {
		return err
	} else if frozen {
		return errors.New("account_frozen")
	}

	return nil
}

// Login IAM user login
func (s *IAMService) Login(email, password string) (client.IAMUser, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = utils.MD5V([]byte(password))

	var user client.IAMUser
	if err := global.GVA_DB.Preload("Wallet").
		Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return client.IAMUser{}, errors.New("invalid_email_or_password")
		}
		return client.IAMUser{}, err
	}

	if user.Status == 0 {
		return client.IAMUser{}, errors.New("account_disabled")
	}
	if user.Status == 2 {
		return client.IAMUser{}, errors.New("account_frozen")
	}
	if frozen, err := s.ParentAccountFrozen(user.ClientID); err != nil {
		return client.IAMUser{}, err
	} else if frozen {
		return client.IAMUser{}, errors.New("account_frozen")
	}

	return user, nil
}

// CheckEmailAvailable Check if email is available
func (s *IAMService) CheckEmailAvailable(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	// Check if email exists in Client
	var existingClient client.Client
	if err := global.GVA_DB.Where("email = ?", email).First(&existingClient).Error; err == nil {
		return errors.New("email_already_used_by_main_account")
	}

	// Check if email exists in IAMUser
	var existingIAM client.IAMUser
	if err := global.GVA_DB.Where("email = ?", email).First(&existingIAM).Error; err == nil {
		return errors.New("email_already_in_use")
	}

	return nil
}

// VerifyCode Verify verification code
func (s *IAMService) VerifyCode(email, code string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	key := fmt.Sprintf("%s%s", IAMVerifyCodePrefix, email)
	storedCode := global.GVA_REDIS.Get(context.Background(), key).Val()
	return strings.EqualFold(storedCode, code)
}

// CreateIAMUser Create IAM user
func (s *IAMService) CreateIAMUser(clientID uint, req request.CreateIAMUserReq) (client.IAMUser, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check if email is available
	if err := s.CheckEmailAvailable(email); err != nil {
		return client.IAMUser{}, err
	}

	user := client.IAMUser{
		ClientID: clientID,
		Email:    email,
		Password: utils.MD5V([]byte(req.Password)),
		Nickname: req.Nickname,
		Status:   1,
		Roles:    req.Roles, // 直接赋值角色ID列表
	}

	// Create user
	if err := global.GVA_DB.Create(&user).Error; err != nil {
		return client.IAMUser{}, err
	}

	return user, nil
}

// GetIAMUserList Get IAM user list
func (s *IAMService) GetIAMUserList(search request.ListIAMUserReq) (list []client.IAMUser, total int64, err error) {

	// Build query conditions
	var conditions []string
	var args []interface{}

	query := global.GVA_DB.Model(&client.IAMUser{}).Order("id desc").Where("1= ?", 1)

	if search.ClientID != 0 {
		conditions = append(conditions, "client_id = ?")
		args = append(args, search.ClientID)
	}
	if search.Email != "" {
		conditions = append(conditions, "email LIKE ?")
		args = append(args, "%"+search.Email+"%")
	}

	if search.Status > 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, search.Status)
	}

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), args...)
	}
	err = query.Count(&total).Error
	if err != nil {
		return
	}

	if search.Page <= 0 {
		search.Page = 1
	}
	if search.PageSize <= 0 {
		search.PageSize = 10
	}

	offset := (search.Page - 1) * search.PageSize
	err = query.Offset(offset).Limit(search.PageSize).Order("id desc").Find(&list).Error
	return
}

// GetIAMUser Get single IAM user
func (s *IAMService) GetIAMUser(clientID uint, userID uint) (client.IAMUser, error) {
	var user client.IAMUser
	err := global.GVA_DB.Preload("Wallet").Where("id = ? AND client_id = ?", userID, clientID).First(&user).Error
	return user, err
}

// UpdateIAMUser Update IAM user
func (s *IAMService) UpdateIAMUser(clientID uint, req request.UpdateIAMUserReq) error {
	var user client.IAMUser
	if err := global.GVA_DB.Where("id = ? AND client_id = ?", req.ID, clientID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user_not_found")
		}
		return err
	}

	if req.Password != "" {
		user.Password = utils.MD5V([]byte(req.Password))
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Roles != nil {
		user.Roles = req.Roles
	}

	if err := global.GVA_DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// DeleteIAMUser Delete IAM user
func (s *IAMService) DeleteIAMUser(clientID uint, userID uint) error {
	var user client.IAMUser
	if err := global.GVA_DB.Where("id = ? AND client_id = ?", userID, clientID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user_not_found")
		}
		return err
	}

	// Delete user
	return global.GVA_DB.Unscoped().Delete(&user).Error
}

// GetAllRoles Get all available roles
func (s *IAMService) GetAllRoles(clientID uint) ([]client.IAMRole, error) {
	var roles []client.IAMRole
	err := global.GVA_DB.Where("client_id = 0 OR client_id = ?", clientID).Order("sort asc").Find(&roles).Error
	return roles, err
}

// GetAllPermissions Get all permissions
func (s *IAMService) GetAllPermissions() ([]client.IAMPermission, error) {
	var permissions []client.IAMPermission
	err := global.GVA_DB.Order("sort asc").Find(&permissions).Error
	return permissions, err
}

// UpdateIAMUserStatus Update IAM user status
func (s *IAMService) UpdateIAMUserStatus(clientID uint, userID uint, status int8) error {
	result := global.GVA_DB.Model(&client.IAMUser{}).
		Where("id = ? AND client_id = ?", userID, clientID).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user_not_found")
	}
	return nil
}

// ResetPassword Reset IAM user password
func (s *IAMService) ResetPassword(clientID uint, userID uint, newPassword string) error {
	// Find user
	var user client.IAMUser
	if err := global.GVA_DB.Where("id = ? AND client_id = ?", userID, clientID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user_not_found")
		}
		return err
	}

	// Update password
	user.Password = utils.MD5V([]byte(newPassword))
	return global.GVA_DB.Save(&user).Error
}

// GetIAMUserByID Get IAM user by ID only
func (s *IAMService) GetIAMUserByID(userID uint) (client.IAMUser, error) {
	var user client.IAMUser
	err := global.GVA_DB.Preload("Wallet").Where("id = ?", userID).First(&user).Error
	return user, err
}

// GetIAMUserByEmail Get IAM user by email
func (s *IAMService) GetIAMUserByEmail(email string) (client.IAMUser, error) {
	var user client.IAMUser
	email = strings.ToLower(strings.TrimSpace(email))
	err := global.GVA_DB.Where("email = ?", email).First(&user).Error
	return user, err
}

// SaveIAMUser Save IAM user
func (s *IAMService) SaveIAMUser(user *client.IAMUser) error {
	return global.GVA_DB.Save(user).Error
}
