package errors

import "github.com/lib/pq"

const (
	UniqueViolationError                          = pq.ErrorCode("23505") // 'unique_violation'
	SchemaAndDataStatementMixingNotSupportedError = pq.ErrorCode("25007") // 'schema_and_data_statement_mixing_not_supported'
)

func IsPqErrorCode(err error, errcode pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
