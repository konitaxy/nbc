package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/ucard/global"
	"gitlab.com/ucard/model/client"
	"gitlab.com/ucard/model/common/response"
	systemReq "gitlab.com/ucard/model/system/request"
	"gitlab.com/ucard/service"
	"gitlab.com/ucard/utils"
)

var clientService = service.ServiceGroupApp.UsersServiceGroup.ClientService
var iamService = service.ServiceGroupApp.UsersServiceGroup.IAMService

// Public routes that require OTP verification before login
var publicOtpRoutes = map[string]bool{
	"/client/login":         true,
	"/client/register":      true,
	"/client/resetPassword": true,
	"/client/iam/login":     true,
}

// OtpAuth handles OTP verification for both public routes and authenticated requests
func OtpAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Request.URL.Path
		act := c.Request.Method
		reqPath := fmt.Sprintf("%s:%s", act, obj)

		claims, _ := utils.GetClaims(c)

		// Handle unauthenticated requests
		if claims == nil {
			if publicOtpRoutes[obj] {
				handlePublicRouteAuth(c, obj, reqPath)
			} else {
				c.Next()
			}
			return
		}

		// Handle authenticated requests
		handleAuthenticatedRequest(c, claims, reqPath)
	}
}

// handlePublicRouteAuth handles OTP verification for public routes (login/register/resetPassword)
func handlePublicRouteAuth(c *gin.Context, obj, reqPath string) {
	mail := strings.ToLower(strings.TrimSpace(c.GetHeader("X-Auth-Mail")))
	code := c.GetHeader("X-Auth-Code")
	isIAM := obj == "/client/iam/login"

	// Code provided - verify it
	if code != "" {
		if !clientService.CodeVerify(mail, code, reqPath, 0, isIAM) {
			abortWithError(c, "Verification code error")
		}
		return
	}

	// No code - validate credentials and request verification
	if err := validateCredentials(c, obj, mail, isIAM); err != nil {
		abortWithError(c, err.Error())
		return
	}

	// Request verification code
	authType := getAuthType(mail, isIAM)
	requireAuth(c, authType)
}

// validateCredentials validates user credentials before requesting OTP
func validateCredentials(c *gin.Context, obj, mail string, isIAM bool) error {
	password := c.GetHeader("X-Auth-Password")

	if isIAM {
		return iamService.PreLogin(mail, password)
	}

	switch obj {
	case "/client/register":
		if cl, _ := clientService.GetClientByMail(mail); cl.ID > 0 {
			return fmt.Errorf("the email has been register")
		}
	case "/client/login":
		cl, _ := clientService.GetClientByMail(mail)
		if cl.ID == 0 {
			return fmt.Errorf("the email not exist")
		}
		if err := clientService.Login(&client.Client{
			Email:    mail,
			Password: utils.MD5V([]byte(password)),
		}); err != nil {
			return fmt.Errorf("The account or password you entered is incorrect. Please try again.")
		}
	case "/client/resetPassword":
		if cl, _ := clientService.GetClientByMail(mail); cl.ID == 0 {
			return fmt.Errorf("the email not exist")
		}
	}
	return nil
}

// handleAuthenticatedRequest handles OTP verification for authenticated users
func handleAuthenticatedRequest(c *gin.Context, claims *systemReq.CustomClaims, reqPath string) {
	isIAM := claims.IsIAM

	// Get verification settings based on user type
	var need bool
	var level uint
	if isIAM {
		need, level = clientService.NeedIAMVerifySetting(claims.Email, reqPath)
	} else {
		need, level = clientService.NeedVerifySetting(claims.Email, reqPath)
	}

	if !need {
		c.Next()
		return
	}

	code := c.GetHeader("X-Auth-Code")

	// Code provided - verify it
	if code != "" {
		if !clientService.CodeVerify(claims.Email, code, reqPath, level, isIAM) {
			abortWithError(c, "Verification code error")
			return
		}
		c.Next()
		return
	}

	// Check if PIN is already verified (level 1 only)
	if level == 1 {
		pinKey := fmt.Sprintf("verify_pin_%s", claims.Email)
		if global.GVA_REDIS.Get(context.Background(), pinKey).Val() == "1" {
			c.Next()
			return
		}
	}

	// Determine auth type and request verification
	authType := getAuthTypeForLevel(claims, level, isIAM)
	requireAuth(c, authType)
}

// getAuthType returns the authentication type for a given email
func getAuthType(email string, isIAM bool) string {
	// Check Redis for stored verification type (set when user binds 2FA)
	if val := global.GVA_REDIS.Get(context.Background(), fmt.Sprintf("verify_type_%s", email)).Val(); val != "" {
		return val
	}
	return "mail"
}

// getAuthTypeForLevel returns the auth type based on user settings and level
func getAuthTypeForLevel(claims *systemReq.CustomClaims, level uint, isIAM bool) string {
	if isIAM {
		// IAM user: check IAM user's settings
		if level == 1 {
			if user, err := iamService.GetIAMUserByID(claims.ID); err == nil && user.BindPin {
				return "pin"
			}
		}
		return getAuthType(claims.Email, true)
	}

	// Main account
	if level == 1 {
		if cl, _ := clientService.GetClient(claims.ID); cl.ID > 0 && cl.BindPin {
			return "pin"
		}
	}
	return getAuthType(claims.Email, false)
}

// requireAuth sends auth required response
func requireAuth(c *gin.Context, authType string) {
	c.Header("X-Auth-Required", "1")
	c.Header("X-Auth-Type", authType)
	response.AuthRequired(c)
	c.Abort()
}

// abortWithError sends error response and aborts
func abortWithError(c *gin.Context, msg string) {
	response.FailWithMessage(msg, c)
	c.Abort()
}
