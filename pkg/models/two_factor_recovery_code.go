package models

// TwoFactorRecoveryCode represents user 2fa recovery codes stored in database
type TwoFactorRecoveryCode struct {
	Uid             int64  `xorm:"PK comment('用户ID')"`
	RecoveryCode    string `xorm:"VARCHAR(64) PK comment('恢复码')"`
	Used            bool   `xorm:"NOT NULL comment('是否已使用')"`
	CreatedUnixTime int64  `comment('创建时间')"`
	UsedUnixTime    int64  `comment('使用时间')"`
}

// TwoFactorRecoveryCodeLoginRequest represents all parameters of 2fa login request via recovery code
type TwoFactorRecoveryCodeLoginRequest struct {
	RecoveryCode string `json:"recoveryCode" binding:"required,notBlank,len=11"`
}
