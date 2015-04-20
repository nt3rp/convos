package db

type DBError string

func (e DBError) Error() string {
	return string(e)
}

const (
	ErrConnection  DBError = "DB Connection"
	ErrRowScan     DBError = "Row Scan"
	ErrRowUnknown  DBError = "Row Unknown"
	ErrRowDelete   DBError = "Row Delete"
	ErrRowCreate   DBError = "Row Create"
	ErrRowUpdate   DBError = "Row Update"
	ErrNoRows      DBError = "No Rows Found"
	ErrTransaction DBError = "Transaction Problem"
)
