package models

import (
	"fmt"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

// TransactionEditScope represents the scope which transaction can be edited
type TransactionEditScope byte

// Editable Transaction Ranges
const (
	TRANSACTION_EDIT_SCOPE_NONE                          TransactionEditScope = 0
	TRANSACTION_EDIT_SCOPE_ALL                           TransactionEditScope = 1
	TRANSACTION_EDIT_SCOPE_TODAY_OR_LATER                TransactionEditScope = 2
	TRANSACTION_EDIT_SCOPE_LAST_24H_OR_LATER             TransactionEditScope = 3
	TRANSACTION_EDIT_SCOPE_THIS_WEEK_OR_LATER            TransactionEditScope = 4
	TRANSACTION_EDIT_SCOPE_THIS_MONTH_OR_LATER           TransactionEditScope = 5
	TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER            TransactionEditScope = 6
	TRANSACTION_EDIT_SCOPE_LAST_RECONCILED_TIME_OR_LATER TransactionEditScope = 7
	TRANSACTION_EDIT_SCOPE_INVALID                       TransactionEditScope = 255
)

// String returns a textual representation of the editable transaction ranges enum
func (s TransactionEditScope) String() string {
	switch s {
	case TRANSACTION_EDIT_SCOPE_NONE:
		return "None"
	case TRANSACTION_EDIT_SCOPE_ALL:
		return "All"
	case TRANSACTION_EDIT_SCOPE_TODAY_OR_LATER:
		return "TodayOrLater"
	case TRANSACTION_EDIT_SCOPE_LAST_24H_OR_LATER:
		return "Last24HourOrLater"
	case TRANSACTION_EDIT_SCOPE_THIS_WEEK_OR_LATER:
		return "ThisWeekOrLater"
	case TRANSACTION_EDIT_SCOPE_THIS_MONTH_OR_LATER:
		return "ThisMonthOrLater"
	case TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER:
		return "ThisYearOrLater"
	case TRANSACTION_EDIT_SCOPE_LAST_RECONCILED_TIME_OR_LATER:
		return "LastReconciledTimeOrLater"
	case TRANSACTION_EDIT_SCOPE_INVALID:
		return "Invalid"
	default:
		return fmt.Sprintf("Invalid(%d)", int(s))
	}
}

// AmountColorType represents the type of amount color in frontend
type AmountColorType byte

// Amount Color Types
const (
	AMOUNT_COLOR_TYPE_DEFAULT        AmountColorType = 0
	AMOUNT_COLOR_TYPE_GREEN          AmountColorType = 1
	AMOUNT_COLOR_TYPE_RED            AmountColorType = 2
	AMOUNT_COLOR_TYPE_YELLOW         AmountColorType = 3
	AMOUNT_COLOR_TYPE_BLACK_OR_WHITE AmountColorType = 4
	AMOUNT_COLOR_TYPE_INVALID        AmountColorType = 255
)

// String returns a textual representation of the amount color type enum
func (s AmountColorType) String() string {
	switch s {
	case AMOUNT_COLOR_TYPE_DEFAULT:
		return "Default"
	case AMOUNT_COLOR_TYPE_GREEN:
		return "Green"
	case AMOUNT_COLOR_TYPE_RED:
		return "Red"
	case AMOUNT_COLOR_TYPE_YELLOW:
		return "Yellow"
	case AMOUNT_COLOR_TYPE_BLACK_OR_WHITE:
		return "Black or White"
	case AMOUNT_COLOR_TYPE_INVALID:
		return "Invalid"
	default:
		return fmt.Sprintf("Invalid(%d)", int(s))
	}
}

// User represents user data stored in database
type User struct {
	Uid                   int64                      `xorm:"PK comment('用户ID')"`
	Username              string                     `xorm:"VARCHAR(32) UNIQUE NOT NULL comment('用户名')"`
	Email                 string                     `xorm:"VARCHAR(100) UNIQUE NOT NULL comment('邮箱')"`
	Nickname              string                     `xorm:"VARCHAR(64) NOT NULL comment('昵称')"`
	Password              string                     `xorm:"VARCHAR(64) NOT NULL comment('密码哈希')"`
	Salt                  string                     `xorm:"VARCHAR(10) NOT NULL comment('密码盐')"`
	CustomAvatarType      string                     `xorm:"VARCHAR(10) comment('自定义头像类型')"`
	DefaultAccountId      int64                      `comment('默认账户ID')"`
	UseLastReconciledTime bool                       `comment('使用上次对账时间')"`
	TransactionEditScope  TransactionEditScope       `xorm:"TINYINT NOT NULL comment('交易编辑范围: 0=无, 1=全部, 2=今天及以后, 3=24小时内及以后, 4=本周及以后, 5=本月及以后, 6=本年及以后')"`
	Language              string                     `xorm:"VARCHAR(10) comment('语言, 如 zh-Hans')"`
	DefaultCurrency       string                     `xorm:"VARCHAR(3) NOT NULL comment('默认货币')"`
	FirstDayOfWeek        core.WeekDay               `xorm:"TINYINT NOT NULL comment('每周第一天')"`
	FiscalYearStart       core.FiscalYearStart       `xorm:"SMALLINT comment('财年起始月份')"`
	CalendarDisplayType   core.CalendarDisplayType   `xorm:"TINYINT comment('日历显示类型')"`
	DateDisplayType       core.DateDisplayType       `xorm:"TINYINT comment('日期显示类型')"`
	LongDateFormat        core.LongDateFormat        `xorm:"TINYINT comment('长日期格式')"`
	ShortDateFormat       core.ShortDateFormat       `xorm:"TINYINT comment('短日期格式')"`
	LongTimeFormat        core.LongTimeFormat        `xorm:"TINYINT comment('长时间格式')"`
	ShortTimeFormat       core.ShortTimeFormat       `xorm:"TINYINT comment('短时间格式')"`
	FiscalYearFormat      core.FiscalYearFormat      `xorm:"TINYINT comment('财年格式')"`
	CurrencyDisplayType   core.CurrencyDisplayType   `xorm:"TINYINT comment('货币显示类型')"`
	NumeralSystem         core.NumeralSystem         `xorm:"TINYINT comment('数字系统')"`
	DecimalSeparator      core.DecimalSeparator      `xorm:"TINYINT comment('小数分隔符')"`
	DigitGroupingSymbol   core.DigitGroupingSymbol   `xorm:"TINYINT comment('千位分隔符')"`
	DigitGrouping         core.DigitGroupingType     `xorm:"TINYINT comment('数字分组方式')"`
	CoordinateDisplayType core.CoordinateDisplayType `xorm:"TINYINT comment('坐标显示类型')"`
	ExpenseAmountColor    AmountColorType            `xorm:"TINYINT comment('支出金额颜色: 0=默认, 1=绿, 2=红, 3=黄, 4=黑白')"`
	IncomeAmountColor     AmountColorType            `xorm:"TINYINT comment('收入金额颜色: 0=默认, 1=绿, 2=红, 3=黄, 4=黑白')"`
	FeatureRestriction    core.UserFeatureRestrictions `comment('功能限制位掩码')"`
	Disabled              bool                       `comment('是否禁用')"`
	Deleted               bool                       `xorm:"NOT NULL comment('是否删除')"`
	EmailVerified         bool                       `xorm:"NOT NULL comment('邮箱是否验证')"`
	CreatedUnixTime       int64                      `comment('创建时间')"`
	UpdatedUnixTime       int64                      `comment('更新时间')"`
	DeletedUnixTime       int64                      `comment('删除时间')"`
	LastLoginUnixTime     int64                      `comment('最后登录时间')"`
}

// UserBasicInfo represents a view-object of user basic info
type UserBasicInfo struct {
	Username              string                     `json:"username"`
	Email                 string                     `json:"email"`
	Nickname              string                     `json:"nickname"`
	AvatarUrl             string                     `json:"avatar"`
	AvatarProvider        string                     `json:"avatarProvider,omitempty"`
	DefaultAccountId      int64                      `json:"defaultAccountId,string"`
	UseLastReconciledTime bool                       `json:"useLastReconciledTime"`
	TransactionEditScope  TransactionEditScope       `json:"transactionEditScope"`
	Language              string                     `json:"language"`
	DefaultCurrency       string                     `json:"defaultCurrency"`
	FirstDayOfWeek        core.WeekDay               `json:"firstDayOfWeek"`
	FiscalYearStart       core.FiscalYearStart       `json:"fiscalYearStart"`
	CalendarDisplayType   core.CalendarDisplayType   `json:"calendarDisplayType"`
	DateDisplayType       core.DateDisplayType       `json:"dateDisplayType"`
	LongDateFormat        core.LongDateFormat        `json:"longDateFormat"`
	ShortDateFormat       core.ShortDateFormat       `json:"shortDateFormat"`
	LongTimeFormat        core.LongTimeFormat        `json:"longTimeFormat"`
	ShortTimeFormat       core.ShortTimeFormat       `json:"shortTimeFormat"`
	FiscalYearFormat      core.FiscalYearFormat      `json:"fiscalYearFormat"`
	CurrencyDisplayType   core.CurrencyDisplayType   `json:"currencyDisplayType"`
	NumeralSystem         core.NumeralSystem         `json:"numeralSystem"`
	DecimalSeparator      core.DecimalSeparator      `json:"decimalSeparator"`
	DigitGroupingSymbol   core.DigitGroupingSymbol   `json:"digitGroupingSymbol"`
	DigitGrouping         core.DigitGroupingType     `json:"digitGrouping"`
	CoordinateDisplayType core.CoordinateDisplayType `json:"coordinateDisplayType"`
	ExpenseAmountColor    AmountColorType            `json:"expenseAmountColor"`
	IncomeAmountColor     AmountColorType            `json:"incomeAmountColor"`
	EmailVerified         bool                       `json:"emailVerified"`
}

// UserLoginRequest represents all parameters of user login request
type UserLoginRequest struct {
	LoginName string `json:"loginName" binding:"required,notBlank,max=100,validUsername|validEmail"`
	Password  string `json:"password" binding:"required,min=6,max=128"`
}

// UserRegisterRequest represents all parameters of user registering request
type UserRegisterRequest struct {
	Username        string       `json:"username" binding:"required,notBlank,max=32,validUsername"`
	Email           string       `json:"email" binding:"required,notBlank,max=100,validEmail"`
	Nickname        string       `json:"nickname" binding:"required,notBlank,max=64,validNickname"`
	Password        string       `json:"password" binding:"required,min=6,max=128"`
	Language        string       `json:"language" binding:"required,min=2,max=16"`
	DefaultCurrency string       `json:"defaultCurrency" binding:"required,len=3,validCurrency"`
	FirstDayOfWeek  core.WeekDay `json:"firstDayOfWeek" binding:"min=0,max=6"`
	TransactionCategoryCreateBatchRequest
}

// UserVerifyEmailRequest represents all parameters of user verify email request
type UserVerifyEmailRequest struct {
	RequestNewToken bool `json:"requestNewToken" binding:"omitempty"`
}

// UserVerifyEmailResponse represents all response parameters after user have verified email
type UserVerifyEmailResponse struct {
	NewToken            string         `json:"newToken,omitempty"`
	User                *UserBasicInfo `json:"user"`
	NotificationContent string         `json:"notificationContent,omitempty"`
}

// UserResendVerifyEmailRequest represents all parameters of user resend verify email request
type UserResendVerifyEmailRequest struct {
	Email    string `json:"email" binding:"omitempty,max=100,validEmail"`
	Password string `json:"password" binding:"omitempty,min=6,max=128"`
}

// UserProfileUpdateRequest represents all parameters of user updating profile request
type UserProfileUpdateRequest struct {
	Email                 string                      `json:"email" binding:"omitempty,notBlank,max=100,validEmail"`
	Nickname              string                      `json:"nickname" binding:"omitempty,notBlank,max=64,validNickname"`
	Password              string                      `json:"password" binding:"omitempty,min=6,max=128"`
	OldPassword           string                      `json:"oldPassword" binding:"omitempty,min=6,max=128"`
	DefaultAccountId      int64                       `json:"defaultAccountId,string" binding:"omitempty,min=1"`
	UseLastReconciledTime *bool                       `json:"useLastReconciledTime" binding:"omitempty"`
	TransactionEditScope  *TransactionEditScope       `json:"transactionEditScope" binding:"omitempty,min=0,max=7"`
	Language              string                      `json:"language" binding:"omitempty,min=2,max=16"`
	DefaultCurrency       string                      `json:"defaultCurrency" binding:"omitempty,len=3,validCurrency"`
	FirstDayOfWeek        *core.WeekDay               `json:"firstDayOfWeek" binding:"omitempty,min=0,max=6"`
	FiscalYearStart       *core.FiscalYearStart       `json:"fiscalYearStart" binding:"omitempty,validFiscalYearStart"`
	CalendarDisplayType   *core.CalendarDisplayType   `json:"calendarDisplayType" binding:"omitempty,min=0,max=4"`
	DateDisplayType       *core.DateDisplayType       `json:"dateDisplayType" binding:"omitempty,min=0,max=3"`
	LongDateFormat        *core.LongDateFormat        `json:"longDateFormat" binding:"omitempty,min=0,max=3"`
	ShortDateFormat       *core.ShortDateFormat       `json:"shortDateFormat" binding:"omitempty,min=0,max=3"`
	LongTimeFormat        *core.LongTimeFormat        `json:"longTimeFormat" binding:"omitempty,min=0,max=3"`
	ShortTimeFormat       *core.ShortTimeFormat       `json:"shortTimeFormat" binding:"omitempty,min=0,max=3"`
	FiscalYearFormat      *core.FiscalYearFormat      `json:"fiscalYearFormat" binding:"omitempty,min=0,max=5"`
	CurrencyDisplayType   *core.CurrencyDisplayType   `json:"currencyDisplayType" binding:"omitempty,min=0,max=11"`
	NumeralSystem         *core.NumeralSystem         `json:"numeralSystem" binding:"omitempty,min=0,max=5"`
	DecimalSeparator      *core.DecimalSeparator      `json:"decimalSeparator" binding:"omitempty,min=0,max=3"`
	DigitGroupingSymbol   *core.DigitGroupingSymbol   `json:"digitGroupingSymbol" binding:"omitempty,min=0,max=4"`
	DigitGrouping         *core.DigitGroupingType     `json:"digitGrouping" binding:"omitempty,min=0,max=3"`
	CoordinateDisplayType *core.CoordinateDisplayType `json:"coordinateDisplayType" binding:"omitempty,min=0,max=6"`
	ExpenseAmountColor    *AmountColorType            `json:"expenseAmountColor" binding:"omitempty,min=0,max=4"`
	IncomeAmountColor     *AmountColorType            `json:"incomeAmountColor" binding:"omitempty,min=0,max=4"`
}

// UserProfileUpdateResponse represents the data returns to frontend after updating profile
type UserProfileUpdateResponse struct {
	User     *UserBasicInfo `json:"user"`
	NewToken string         `json:"newToken,omitempty"`
}

// UserProfileResponse represents a view-object of user profile
type UserProfileResponse struct {
	*UserBasicInfo
	NoPassword  bool  `json:"noPassword,omitempty"`
	LastLoginAt int64 `json:"lastLoginAt"`
}

// CanEditTransactionByTransactionTime returns whether this user can edit transaction with specified transaction time
func (u *User) CanEditTransactionByTransactionTime(transactionTime int64, clientTimezone *time.Location, account *Account, destinationAccount *Account) bool {
	if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_NONE {
		return false
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_ALL {
		return true
	}

	now := time.Now()

	transactionUnixTime := utils.GetUnixTimeFromTransactionTime(transactionTime)

	if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_LAST_24H_OR_LATER {
		return transactionUnixTime > now.Add(-24*time.Hour).Unix()
	}

	clientNow := now.In(clientTimezone)
	clientTodayStartTime := utils.GetStartOfDay(clientNow)

	if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_TODAY_OR_LATER {
		return transactionUnixTime > clientTodayStartTime.Unix()
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_THIS_WEEK_OR_LATER {
		dayOfWeek := int(now.Weekday()) - int(u.FirstDayOfWeek)

		if dayOfWeek < 0 {
			dayOfWeek += 7
		}

		clientWeekStartTime := clientTodayStartTime.AddDate(0, 0, -dayOfWeek)
		return transactionUnixTime > clientWeekStartTime.Unix()
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_THIS_MONTH_OR_LATER {
		clientMonthStartTime := clientTodayStartTime.AddDate(0, 0, -(now.Day() - 1))
		return transactionUnixTime > clientMonthStartTime.Unix()
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_THIS_YEAR_OR_LATER {
		clientYearStartTime := clientTodayStartTime.AddDate(0, 0, -(now.YearDay() - 1))
		return transactionUnixTime > clientYearStartTime.Unix()
	} else if u.TransactionEditScope == TRANSACTION_EDIT_SCOPE_LAST_RECONCILED_TIME_OR_LATER && u.UseLastReconciledTime {
		minAccountLastReconciledTime := int64(0)

		if account != nil {
			minAccountLastReconciledTime = account.GetLastReconciledTime()
		}

		if destinationAccount != nil {
			destinationAccountLastReconciledTime := destinationAccount.GetLastReconciledTime()

			if destinationAccountLastReconciledTime > minAccountLastReconciledTime {
				minAccountLastReconciledTime = destinationAccountLastReconciledTime
			}
		}

		return transactionUnixTime > minAccountLastReconciledTime
	}

	return false
}

// ToUserBasicInfo returns a user basic view-object according to database model
func (u *User) ToUserBasicInfo(avatarProvider core.UserAvatarProviderType, avatarUrl string) *UserBasicInfo {
	fiscalYearStart := u.FiscalYearStart

	if fiscalYearStart < core.FISCAL_YEAR_START_MIN || fiscalYearStart > core.FISCAL_YEAR_START_MAX {
		fiscalYearStart = core.FISCAL_YEAR_START_DEFAULT
	}

	return &UserBasicInfo{
		Username:              u.Username,
		Email:                 u.Email,
		Nickname:              u.Nickname,
		AvatarUrl:             avatarUrl,
		AvatarProvider:        string(avatarProvider),
		DefaultAccountId:      u.DefaultAccountId,
		UseLastReconciledTime: u.UseLastReconciledTime,
		TransactionEditScope:  u.TransactionEditScope,
		Language:              u.Language,
		DefaultCurrency:       u.DefaultCurrency,
		FirstDayOfWeek:        u.FirstDayOfWeek,
		FiscalYearStart:       fiscalYearStart,
		CalendarDisplayType:   u.CalendarDisplayType,
		DateDisplayType:       u.DateDisplayType,
		LongDateFormat:        u.LongDateFormat,
		ShortDateFormat:       u.ShortDateFormat,
		LongTimeFormat:        u.LongTimeFormat,
		ShortTimeFormat:       u.ShortTimeFormat,
		DecimalSeparator:      u.DecimalSeparator,
		FiscalYearFormat:      u.FiscalYearFormat,
		CurrencyDisplayType:   u.CurrencyDisplayType,
		NumeralSystem:         u.NumeralSystem,
		DigitGroupingSymbol:   u.DigitGroupingSymbol,
		DigitGrouping:         u.DigitGrouping,
		CoordinateDisplayType: u.CoordinateDisplayType,
		ExpenseAmountColor:    u.ExpenseAmountColor,
		IncomeAmountColor:     u.IncomeAmountColor,
		EmailVerified:         u.EmailVerified,
	}
}

// ToUserProfileResponse returns a user profile view-object according to database model
func (u *User) ToUserProfileResponse(basicInfo *UserBasicInfo) *UserProfileResponse {
	return &UserProfileResponse{
		UserBasicInfo: basicInfo,
		NoPassword:    u.Password == "",
		LastLoginAt:   u.LastLoginUnixTime,
	}
}
