package server

type EmptyArgs struct {
}

type ProvinceInfo struct {
	Name           string
	CountOfIncome  int
	CountOfOutcome int
}
type NewCharityIncome struct {
	// 捐赠单位或个人，捐赠类别（如果是物资就写物资类型，如果是资金，就写资金），捐赠数量（物资数量或金额），捐赠项目
	Time     int
	From     string
	Category string // "资金" or others
	Amount   int
	Detail   string
	Province string
}
type BatchNewCharityIncome struct {
	Data []NewCharityIncome
}
type CharityIncome struct {
	Id          string
	Time        int // unix epoch timestamp
	From        string
	Category    string // "资金" or others
	Amount      int
	Detail      string
	Province    string
	ChainTime   int
	ChainHeight int
	ChainTxid   string
}

type BatchCharityIncome struct {
	Data []CharityIncome
}
type NewCharityOutcome struct {
	// 编号，时间，引用捐赠编号（多个），支出类别，支出数量，支出原因，接收人
	Time     int
	Source   string
	Category string
	Amount   int
	Detail   string
	To       string
	Province string
}
type BatchNewCharityOutcome struct {
	Data []NewCharityOutcome
}
type CharityOutcome struct {
	Id          string
	Time        int // unix epoch timestamp
	Source      string
	Category    string
	Amount      int
	Detail      string
	To          string
	Province    string
	ChainTime   int
	ChainHeight int
	ChainTxid   string
}
type BatchCharityOutcome struct {
	Data []CharityOutcome
}
type CharitySummary struct {
	TotalIncome        int
	TotalOutcome       int
	TotalLeft          int
	TotalIncomeThings  int
	TotalOutcomeThings int
	CountOfIncome      int
	CountOfOutcome     int
	ProvinceInfos      []ProvinceInfo
	VisitNum           int
}
type PagingParam struct {
	Offset   int
	Limit    int
	FromTime int
	ToTime   int
}

type ListCharityIncomeResult struct {
	Results          []CharityIncome
	Total            int
	TotalInTimeRange int
}

type ListCharityOutcomeResult struct {
	Results          []CharityOutcome
	Total            int
	TotalInTimeRange int
}

type CharityApi interface {
	AddCharityIncome(args *NewCharityIncome, result *CharityIncome) error
	AddCharityOutcome(args *NewCharityOutcome, result *CharityOutcome) error
	BatchAddCharityIncome(args *BatchNewCharityIncome, result *BatchCharityIncome) error
	BatchAddCharityOutcome(args *BatchNewCharityOutcome, result *BatchCharityOutcome) error
	ListCharityIncome(args *PagingParam, result *ListCharityIncomeResult) error
	ListCharityOutcome(args *PagingParam, result *ListCharityOutcomeResult) error
	ListCharitySummary(args *EmptyArgs, results *CharitySummary) error
}
