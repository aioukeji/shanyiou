package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"
	"time"
)

func (s *RpcStub) AddCharityIncome(newIncome *NewCharityIncome, result *CharityIncome) error {

	redisLock.Lock()
	defer redisLock.Unlock()

	t := newIncome.Time
	if t == 0 {
		t = int(time.Now().Unix())
	}
	currentCount, err := RedisZCount(incomeKeyZ)
	if err != nil {
		return fmt.Errorf("load data error %v", err)
	}
	item := CharityIncome{
		Id:       fmt.Sprintf("s%06d", currentCount),
		Time:     t,
		From:     newIncome.From,
		Category: newIncome.Category,
		Amount:   newIncome.Amount,
		Detail:   newIncome.Detail,
		Province: newIncome.Province,
	}
	item.FixProvince()

	err = RedisZAdd(incomeKeyZ, item.Time, &item)
	if err != nil {
		return fmt.Errorf("persistence err %v", err)
	}
	summary := LoadSummary()
	summary.AddIncome(item)
	summary.VisitNum += 1
	err = RedisSet(summaryKey, summary)
	if err != nil {
		return err
	}
	if hasFabricBackend {
		go func() {
			err := flushIncomeToFabric(&s.FabricClient, item)
			if verbose {
				fmt.Printf("flush newIncome %v, %v\n", err, item)
			}
		}()
	}
	*result = item
	return nil
}

func (s *RpcStub) BatchAddCharityIncome(args *BatchNewCharityIncome, result *BatchCharityIncome) error {
	for _, item := range args.Data {
		var income CharityIncome
		err := s.AddCharityIncome(&item, &income)
		if err != nil {
			return err
		}
		result.Data = append(result.Data, income)
	}
	return nil
}
func (s *RpcStub) AddCharityOutcome(args *NewCharityOutcome, result *CharityOutcome) error {

	redisLock.Lock()
	defer redisLock.Unlock()

	t := args.Time
	if t == 0 {
		t = int(time.Now().Unix())
	}
	currentCount, err := RedisZCount(outcomeKeyZ)
	if err != nil {
		return fmt.Errorf("load data error %v", err)
	}
	item := CharityOutcome{
		Id:       fmt.Sprintf("z%06d", currentCount),
		Time:     t,
		To:       args.To,
		Category: args.Category,
		Amount:   args.Amount,
		Detail:   args.Detail,
		Source:   args.Source,
		Province: args.Province,
	}
	item.FixProvince()

	err = RedisZAdd(outcomeKeyZ, item.Time, &item)
	if err != nil {
		return fmt.Errorf("persistence err %v", err)
	}

	summary := LoadSummary()
	summary.AddOutcome(item)
	summary.VisitNum += 1
	err = RedisSet(summaryKey, summary)
	if err != nil {
		return err
	}

	if hasFabricBackend {
		// copy by value
		go func() {
			err := flushOutcomeToFabric(&s.FabricClient, item)
			if verbose {
				fmt.Printf("flush outcome %v, %v\n", err, item)
			}
		}()
	}

	*result = item
	return nil
}

func (s *RpcStub) BatchAddCharityOutcome(args *BatchNewCharityOutcome, result *BatchCharityOutcome) error {
	for _, item := range args.Data {
		var outcome CharityOutcome
		err := s.AddCharityOutcome(&item, &outcome)
		if err != nil {
			return err
		}
		result.Data = append(result.Data, outcome)
	}
	return nil
}
func (s *RpcStub) ListCharityIncome(args *PagingParam, result *ListCharityIncomeResult) error {
	redisLock.Lock()
	defer redisLock.Unlock()
	toTime := "+inf"
	if args.ToTime != 0 {
		toTime = fmt.Sprintf(")%d", args.ToTime)
	}
	fromTime := fmt.Sprintf("%d", args.FromTime)
	err := redisClient.ZRevRangeByScore(incomeKeyZ, &redis.ZRangeBy{
		Min:    fromTime,
		Max:    toTime,
		Offset: int64(args.Offset),
		Count:  int64(args.Limit),
	}).ScanSlice(&result.Results)
	if err != nil {
		return fmt.Errorf("load data err %v", err)
	}
	AddVisitCount()
	count, err := redisClient.ZCount(incomeKeyZ, fromTime, toTime).Result()
	if err != nil {
		return fmt.Errorf("count err %v", err)
	}
	result.Total = int(count)
	return nil
}

func (s *RpcStub) ListCharityOutcome(args *PagingParam, result *ListCharityOutcomeResult) error {
	redisLock.Lock()
	defer redisLock.Unlock()
	toTime := "+inf"
	if args.ToTime != 0 {
		toTime = fmt.Sprintf(")%d", args.ToTime)
	}
	fromTime := fmt.Sprintf("%d", args.FromTime)
	err := redisClient.ZRevRangeByScore(outcomeKeyZ, &redis.ZRangeBy{
		Min:    fromTime,
		Max:    toTime,
		Offset: int64(args.Offset),
		Count:  int64(args.Limit),
	}).ScanSlice(&result.Results)
	if err != nil {
		return fmt.Errorf("load data err %v", err)
	}
	AddVisitCount()
	count, err := redisClient.ZCount(outcomeKeyZ, fromTime, toTime).Result()
	if err != nil {
		return fmt.Errorf("count err %v", err)
	}
	result.Total = int(count)
	return nil
}

func (s *RpcStub) ListCharitySummary(args *EmptyArgs, results *CharitySummary) error {
	redisLock.Lock()
	defer redisLock.Unlock()
	summary := LoadSummary()
	summary.VisitNum += 1
	err := RedisSet(summaryKey, summary)
	if err != nil {
		return fmt.Errorf("save summary err %v", err)
	}
	*results = summary
	return nil
}

/////////////////////// Model methods ////////////////////

func (income *CharityIncome) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, income)
}
func (income *CharityIncome) MarshalBinary() ([]byte, error) {
	return json.Marshal(income)
}
func (outcome *CharityOutcome) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, outcome)
}
func (outcome *CharityOutcome) MarshalBinary() ([]byte, error) {
	return json.Marshal(outcome)
}

func (summary *CharitySummary) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, summary)
}
func (summary *CharitySummary) MarshalBinary() ([]byte, error) {
	return json.Marshal(summary)
}

func (income *CharityIncome) FixProvince() {
	if income.Province != "" {
		return
	}
	income.Province = normalizeProvinceName(income.From)
}
func (outcome *CharityOutcome) FixProvince() {
	if outcome.Province != "" {
		return
	}
	outcome.Province = normalizeProvinceName(outcome.To)
}

func (summary *CharitySummary) Reset() {
	summary.CountOfIncome = 0
	summary.CountOfOutcome = 0
	summary.TotalIncome = 0
	summary.TotalOutcome = 0
	summary.TotalLeft = 0
	summary.ProvinceInfos = GetDefaultProvinceInfo()

}

func (summary *CharitySummary) AddIncome(item CharityIncome) {
	if item.Category == "资金" {
		summary.TotalIncome += item.Amount
		summary.TotalLeft += item.Amount
	} else {
		summary.TotalIncomeThings += item.Amount
	}
	summary.CountOfIncome += 1
	province := normalizeProvinceName(item.Province)
	for idx, item := range summary.ProvinceInfos {
		if item.Name == province {
			summary.ProvinceInfos[idx].CountOfIncome += 1
		}
	}
}

func (summary *CharitySummary) AddOutcome(item CharityOutcome) {
	if item.Category == "资金" {
		summary.TotalOutcome += item.Amount
		summary.TotalLeft -= item.Amount
	} else {
		summary.TotalOutcomeThings += item.Amount
	}
	summary.CountOfOutcome += 1
	province := normalizeProvinceName(item.Province)
	for idx, item := range summary.ProvinceInfos {
		if item.Name == province {
			summary.ProvinceInfos[idx].CountOfOutcome += 1
		}
	}
}

func AddTxidToOutcome(outcome CharityOutcome, txid string) {
	newOutcome := outcome
	newOutcome.ChainTxid = txid
	newOutcome.ChainTime = int(time.Now().Unix())
	newOutcome.ChainHeight = getHeight()
	err := RedisChangeZValue(outcomeKeyZ, outcome.Time, &outcome, &newOutcome)
	if err != nil {
		fmt.Println(err)
	}
}

func AddTxidToIncome(income CharityIncome, txid string) {
	newIncome := income
	newIncome.ChainTxid = txid
	newIncome.ChainTime = int(time.Now().Unix())
	newIncome.ChainHeight = getHeight()
	err := RedisChangeZValue(incomeKeyZ, income.Time, &income, &newIncome)
	if err != nil {
		fmt.Println(err)
	}
}

func RedisChangeZValue(key string, score int, old interface{}, new interface{}) error {
	redisLock.Lock()
	defer redisLock.Unlock()
	err := redisClient.ZRem(key, old).Err()
	if err != nil {
		return fmt.Errorf("zrem failed %v", err)
	}
	err = RedisZAdd(key, score, new)
	if err != nil {
		return fmt.Errorf("zadd failed %v", err)
	}
	return nil
}

func LoadSummary() CharitySummary {
	var summary CharitySummary
	summary.Reset()
	err := RedisScan(summaryKey, &summary)
	if err != nil {
		fmt.Println("load summary err", err)
	}
	return summary
}

func AddVisitCount() {
	summary := LoadSummary()
	summary.VisitNum += 1
	err := RedisSet(summaryKey, summary)
	if err != nil {
		fmt.Println("save summary err", err)
	}
}

func normalizeProvinceName(name string) string {
	if strings.HasPrefix(name, "武汉") {
		return "湖北"
	}
	for _, item := range provinceNames {
		if strings.HasPrefix(name, item) {
			return item
		}
	}
	return "未知"
}

func GetDefaultProvinceInfo() []ProvinceInfo {
	var provinceInfos []ProvinceInfo
	for _, provinceName := range provinceNames {
		provinceInfos = append(provinceInfos,
			ProvinceInfo{
				Name:           provinceName,
				CountOfOutcome: 0,
				CountOfIncome:  0,
			})
	}
	return provinceInfos
}

func getHeight() int {
	// TODO
	return 0
}
