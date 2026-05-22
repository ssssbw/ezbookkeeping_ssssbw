package models

// TransactionTagIndex represents transaction and transaction tag relation stored in database
type TransactionTagIndex struct {
	TagIndexId      int64 `xorm:"PK comment('标签索引ID')"`
	Uid             int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_tag_id_transaction_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_time_tag_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_id) NOT NULL comment('用户ID')"`
	Deleted         bool  `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_tag_id_transaction_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_time_tag_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_id) NOT NULL comment('是否删除')"`
	TransactionTime int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_transaction_time_tag_id) NOT NULL comment('交易时间')"`
	TagId           int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_tag_id_transaction_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_time_tag_id) comment('标签ID')"`
	TransactionId   int64 `xorm:"INDEX(IDX_transaction_tag_index_uid_deleted_tag_id_transaction_id) INDEX(IDX_transaction_tag_index_uid_deleted_transaction_id) comment('交易ID')"`
	CreatedUnixTime int64 `comment('创建时间')"`
	UpdatedUnixTime int64 `comment('更新时间')"`
	DeletedUnixTime int64 `comment('删除时间')"`
}
