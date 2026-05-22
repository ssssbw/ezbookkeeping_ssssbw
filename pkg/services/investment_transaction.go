package services

import (
	"fmt"
	"time"

	"xorm.io/xorm"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/datastore"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/uuid"
)

type InvestmentTransactionService struct {
	ServiceUsingDB
	ServiceUsingUuid
}

var InvestmentTransactions = &InvestmentTransactionService{
	ServiceUsingDB: ServiceUsingDB{
		container: datastore.Container,
	},
	ServiceUsingUuid: ServiceUsingUuid{
		container: uuid.Container,
	},
}

func (s *InvestmentTransactionService) GetAllTransactionsByUid(c core.Context, uid int64, assetId int64, accountId int64, txType models.InvestmentTransactionType, startTime int64, endTime int64) ([]*models.InvestmentTransaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	condition := "uid=? AND deleted=?"
	conditionParams := make([]any, 0, 8)
	conditionParams = append(conditionParams, uid)
	conditionParams = append(conditionParams, false)

	if assetId > 0 {
		condition = condition + " AND asset_id=?"
		conditionParams = append(conditionParams, assetId)
	}

	if accountId > 0 {
		condition = condition + " AND account_id=?"
		conditionParams = append(conditionParams, accountId)
	}

	if txType > 0 {
		condition = condition + " AND type=?"
		conditionParams = append(conditionParams, txType)
	}

	if startTime > 0 {
		condition = condition + " AND trade_time>=?"
		conditionParams = append(conditionParams, startTime)
	}

	if endTime > 0 {
		condition = condition + " AND trade_time<=?"
		conditionParams = append(conditionParams, endTime)
	}

	var transactions []*models.InvestmentTransaction
	err := s.UserDataDB(uid).NewSession(c).Where(condition, conditionParams...).OrderBy("trade_time desc").Find(&transactions)

	return transactions, err
}

func (s *InvestmentTransactionService) GetTransactionByTransactionId(c core.Context, uid int64, transactionId int64) (*models.InvestmentTransaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if transactionId <= 0 {
		return nil, errs.ErrInvestmentTransactionIdInvalid
	}

	transaction := &models.InvestmentTransaction{}
	has, err := s.UserDataDB(uid).NewSession(c).ID(transactionId).Where("uid=? AND deleted=?", uid, false).Get(transaction)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, errs.ErrInvestmentTransactionNotFound
	}

	return transaction, nil
}

func (s *InvestmentTransactionService) GetTransactionsByAssetId(c core.Context, uid int64, assetId int64) ([]*models.InvestmentTransaction, error) {
	if uid <= 0 {
		return nil, errs.ErrUserIdInvalid
	}

	if assetId <= 0 {
		return nil, errs.ErrInvestmentTransactionAssetIdInvalid
	}

	var transactions []*models.InvestmentTransaction
	err := s.UserDataDB(uid).NewSession(c).Where("uid=? AND deleted=? AND asset_id=?", uid, false, assetId).OrderBy("trade_time desc").Find(&transactions)

	return transactions, err
}

func (s *InvestmentTransactionService) CreateTransaction(c core.Context, transaction *models.InvestmentTransaction) error {
	if transaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	transaction.TransactionId = s.GenerateUuid(uuid.UUID_TYPE_INVESTMENT_TRANS)

	if transaction.TransactionId < 1 {
		return errs.ErrSystemIsBusy
	}

	transaction.Deleted = false
	transaction.CreatedUnixTime = time.Now().Unix()
	transaction.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(transaction.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		_, err := sess.Insert(transaction)

		if err != nil {
			return err
		}

		return s.updateAccountBalanceForCreate(sess, transaction)
	})
}

func (s *InvestmentTransactionService) ModifyTransaction(c core.Context, transaction *models.InvestmentTransaction) error {
	if transaction.Uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	transaction.UpdatedUnixTime = time.Now().Unix()

	return s.UserDataDB(transaction.Uid).DoTransaction(c, func(sess *xorm.Session) error {
		oldTransaction := &models.InvestmentTransaction{}
		has, err := sess.ID(transaction.TransactionId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(oldTransaction)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrInvestmentTransactionNotFound
		}

		updatedRows, err := sess.ID(transaction.TransactionId).Cols("asset_id", "account_id", "type", "trade_time", "confirm_time", "quantity", "price", "amount", "fee", "related_transaction_id", "timezone_utc_offset", "comment", "updated_unix_time").Where("uid=? AND deleted=?", transaction.Uid, false).Update(transaction)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			return errs.ErrInvestmentTransactionNotFound
		}

		return s.updateAccountBalanceForModify(sess, oldTransaction, transaction)
	})
}

func (s *InvestmentTransactionService) DeleteTransaction(c core.Context, uid int64, transactionId int64) error {
	if uid <= 0 {
		return errs.ErrUserIdInvalid
	}

	now := time.Now().Unix()

	updateModel := &models.InvestmentTransaction{
		Deleted:         true,
		DeletedUnixTime: now,
	}

	return s.UserDataDB(uid).DoTransaction(c, func(sess *xorm.Session) error {
		oldTransaction := &models.InvestmentTransaction{}
		has, err := sess.ID(transactionId).Where("uid=? AND deleted=?", uid, false).Get(oldTransaction)

		if err != nil {
			return err
		} else if !has {
			return errs.ErrInvestmentTransactionNotFound
		}

		deletedRows, err := sess.ID(oldTransaction.TransactionId).Cols("deleted", "deleted_unix_time").Where("uid=? AND deleted=?", uid, false).Update(updateModel)

		if err != nil {
			return err
		} else if deletedRows < 1 {
			return errs.ErrInvestmentTransactionNotFound
		}

		return s.updateAccountBalanceForDelete(sess, oldTransaction)
	})
}

func (s *InvestmentTransactionService) updateAccountBalanceForCreate(sess *xorm.Session, transaction *models.InvestmentTransaction) error {
	account := &models.Account{}
	has, err := sess.ID(transaction.AccountId).Where("uid=? AND deleted=?", transaction.Uid, false).Get(account)

	if err != nil {
		return err
	} else if !has {
		return errs.ErrAccountNotFound
	}

	var balanceChange int64

	switch transaction.Type {
	case models.INVESTMENT_TRANSACTION_TYPE_BUY:
		balanceChange = -(transaction.Amount + transaction.Fee)
	case models.INVESTMENT_TRANSACTION_TYPE_SELL:
		balanceChange = transaction.Amount - transaction.Fee
	case models.INVESTMENT_TRANSACTION_TYPE_DIVIDEND_CASH:
		balanceChange = transaction.Amount
	case models.INVESTMENT_TRANSACTION_TYPE_DIVIDEND_REINVEST, models.INVESTMENT_TRANSACTION_TYPE_SPLIT, models.INVESTMENT_TRANSACTION_TYPE_CONVERSION_OUT, models.INVESTMENT_TRANSACTION_TYPE_CONVERSION_IN:
		balanceChange = 0
	default:
		return errs.ErrInvestmentTransactionTypeInvalid
	}

	if balanceChange != 0 {
		account.UpdatedUnixTime = time.Now().Unix()
		updatedRows, err := sess.ID(account.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", balanceChange)).Cols("updated_unix_time").Where("uid=? AND deleted=?", account.Uid, false).Update(account)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			log.Errorf(nil, "[investment_transaction.updateAccountBalanceForCreate] failed to update account balance")
			return errs.ErrDatabaseOperationFailed
		}
	}

	return nil
}

func (s *InvestmentTransactionService) updateAccountBalanceForDelete(sess *xorm.Session, oldTransaction *models.InvestmentTransaction) error {
	account := &models.Account{}
	has, err := sess.ID(oldTransaction.AccountId).Where("uid=? AND deleted=?", oldTransaction.Uid, false).Get(account)

	if err != nil {
		return err
	} else if !has {
		return errs.ErrAccountNotFound
	}

	var balanceChange int64

	switch oldTransaction.Type {
	case models.INVESTMENT_TRANSACTION_TYPE_BUY:
		balanceChange = oldTransaction.Amount + oldTransaction.Fee
	case models.INVESTMENT_TRANSACTION_TYPE_SELL:
		balanceChange = -(oldTransaction.Amount - oldTransaction.Fee)
	case models.INVESTMENT_TRANSACTION_TYPE_DIVIDEND_CASH:
		balanceChange = -oldTransaction.Amount
	case models.INVESTMENT_TRANSACTION_TYPE_DIVIDEND_REINVEST, models.INVESTMENT_TRANSACTION_TYPE_SPLIT, models.INVESTMENT_TRANSACTION_TYPE_CONVERSION_OUT, models.INVESTMENT_TRANSACTION_TYPE_CONVERSION_IN:
		balanceChange = 0
	default:
		return errs.ErrInvestmentTransactionTypeInvalid
	}

	if balanceChange != 0 {
		account.UpdatedUnixTime = time.Now().Unix()
		updatedRows, err := sess.ID(account.AccountId).SetExpr("balance", fmt.Sprintf("balance+(%d)", balanceChange)).Cols("updated_unix_time").Where("uid=? AND deleted=?", account.Uid, false).Update(account)

		if err != nil {
			return err
		} else if updatedRows < 1 {
			log.Errorf(nil, "[investment_transaction.updateAccountBalanceForDelete] failed to update account balance")
			return errs.ErrDatabaseOperationFailed
		}
	}

	return nil
}

func (s *InvestmentTransactionService) updateAccountBalanceForModify(sess *xorm.Session, oldTransaction *models.InvestmentTransaction, newTransaction *models.InvestmentTransaction) error {
	if oldTransaction.AccountId == newTransaction.AccountId && oldTransaction.Type == newTransaction.Type && oldTransaction.Amount == newTransaction.Amount && oldTransaction.Fee == newTransaction.Fee {
		return nil
	}

	err := s.updateAccountBalanceForDelete(sess, oldTransaction)
	if err != nil {
		return err
	}

	return s.updateAccountBalanceForCreate(sess, newTransaction)
}
