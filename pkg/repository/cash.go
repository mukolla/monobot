package repository

type CashBucket string

const (
	TransactionCash CashBucket = "transaction_cash"
	UserInfoCash    CashBucket = "user_info_cash"
)

type CashRepository interface {
	TransactionSave(cacheKey string, transaction []Transaction) error
	TransactionGet(cacheKey string) ([]Transaction, error)
	TransactionEmptyExists(cacheKey string) (bool, error)

	UserInfoSave(cacheKey string, userInfo UserInfo) error
	UserInfoGet(cacheKey string) (UserInfo, error)
}
