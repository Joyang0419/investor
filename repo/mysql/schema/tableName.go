package schema

type tableName = string

// Table names
// Why: To avoid hardcoding table names in the code
// 未來就是用這裡的常數來取得table name, 在搜尋哪些地方有使用到table name時，只要搜尋這個檔案即可

const (
	TableNameDailyPrice tableName = "DailyPrice"
)
