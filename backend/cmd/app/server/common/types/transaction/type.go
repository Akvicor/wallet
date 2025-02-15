package transaction

type Type int64

const (
	TypeIncome       Type = 1 // 收入
	TypeExpense      Type = 2 // 支出
	TypeTransfer     Type = 3 // 转账
	TypeAutoTransfer Type = 4 // 自动转账
	TypeExchange     Type = 5 // 兑换
	TypeAutoExchange Type = 6 // 自动兑换
)

func (t Type) Valid() bool {
	switch t {
	case TypeIncome, TypeExpense, TypeTransfer, TypeAutoTransfer, TypeExchange, TypeAutoExchange:
		return true
	default:
		return false
	}
}

func (t Type) String() string {
	switch t {
	case TypeIncome:
		return "收入"
	case TypeExpense:
		return "支出"
	case TypeTransfer:
		return "转账"
	case TypeAutoTransfer:
		return "自动转账"
	case TypeExchange:
		return "兑换"
	case TypeAutoExchange:
		return "自动兑换"
	default:
		return "未知"
	}
}

func (t Type) StringEnglish() string {
	switch t {
	case TypeIncome:
		return "Income"
	case TypeExpense:
		return "Expense"
	case TypeTransfer:
		return "Transfer"
	case TypeAutoTransfer:
		return "AutoTransfer"
	case TypeExchange:
		return "Exchange"
	case TypeAutoExchange:
		return "AutoExchange"
	default:
		return "Unknown"
	}
}

func (t Type) Colour() string {
	switch t {
	case TypeIncome:
		return "#389e0d"
	case TypeExpense:
		return "#cf1322"
	case TypeTransfer:
		return "#0958d9"
	case TypeAutoTransfer:
		return "#0958d9"
	case TypeExchange:
		return "#531dab"
	case TypeAutoExchange:
		return "#531dab"
	default:
		return "#000000"
	}
}

type TypeDetail struct {
	Type        Type   `json:"type"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
	Colour      string `json:"colour"`
}

var AllType = []*TypeDetail{
	{TypeIncome, TypeIncome.String(), TypeIncome.StringEnglish(), TypeIncome.Colour()},
	{TypeExpense, TypeExpense.String(), TypeExpense.StringEnglish(), TypeExpense.Colour()},
	{TypeTransfer, TypeTransfer.String(), TypeTransfer.StringEnglish(), TypeTransfer.Colour()},
	{TypeAutoTransfer, TypeAutoTransfer.String(), TypeAutoTransfer.StringEnglish(), TypeAutoTransfer.Colour()},
	{TypeExchange, TypeExchange.String(), TypeExchange.StringEnglish(), TypeExchange.Colour()},
	{TypeAutoExchange, TypeAutoExchange.String(), TypeAutoExchange.StringEnglish(), TypeAutoExchange.Colour()},
}
