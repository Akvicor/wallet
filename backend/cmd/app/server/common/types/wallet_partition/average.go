package wallet_partition

type AverageType int64

const (
	TypeAverageNormal  AverageType = 0 // 不均分
	TypeAverageWeek    AverageType = 1 // 按周均分
	TypeAverageMonth   AverageType = 2 // 按月均分
	TypeAverageQuarter AverageType = 3 // 按季度均分
	TypeAverageYear    AverageType = 4 // 按年均分
)

func (t *AverageType) Valid() bool {
	switch *t {
	case TypeAverageNormal, TypeAverageWeek, TypeAverageMonth, TypeAverageQuarter, TypeAverageYear:
		return true
	default:
		return false
	}
}

func (t AverageType) String() string {
	switch t {
	case TypeAverageNormal:
		return "不均分"
	case TypeAverageWeek:
		return "按周"
	case TypeAverageMonth:
		return "按月"
	case TypeAverageQuarter:
		return "按季度"
	case TypeAverageYear:
		return "按年"
	default:
		return "未知"
	}
}

func (t AverageType) StringEnglish() string {
	switch t {
	case TypeAverageNormal:
		return "Normal"
	case TypeAverageWeek:
		return "Week"
	case TypeAverageMonth:
		return "Month"
	case TypeAverageQuarter:
		return "Quarter"
	case TypeAverageYear:
		return "Year"
	default:
		return "Unknown"
	}
}

type AverageTypeDetail struct {
	Type        AverageType `json:"type"`
	Name        string      `json:"name"`
	EnglishName string      `json:"english_name"`
}

var AllAverageType = []*AverageTypeDetail{
	{TypeAverageNormal, TypeAverageNormal.String(), TypeAverageNormal.StringEnglish()},
	{TypeAverageWeek, TypeAverageWeek.String(), TypeAverageWeek.StringEnglish()},
	{TypeAverageMonth, TypeAverageMonth.String(), TypeAverageMonth.StringEnglish()},
	{TypeAverageQuarter, TypeAverageQuarter.String(), TypeAverageQuarter.StringEnglish()},
	{TypeAverageYear, TypeAverageYear.String(), TypeAverageYear.StringEnglish()},
}
