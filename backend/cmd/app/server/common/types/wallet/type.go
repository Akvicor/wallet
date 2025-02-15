package wallet

type Type int64

const (
	TypeNormal      Type = 1 // 普通钱包
	TypeHideBalance Type = 2 // 默认隐藏金额
	TypeDebt        Type = 3 // 债务
	TypeWishlist    Type = 4 // 愿望但
)

func (t *Type) IsNormal() bool {
	switch *t {
	case TypeNormal, TypeHideBalance:
		return true
	default:
		return false
	}
}

func (t *Type) IsDebt() bool {
	return *t == TypeDebt
}

func (t *Type) IsWishlist() bool {
	return *t == TypeWishlist
}

func (t *Type) SameTypes() []Type {
	if t.IsNormal() {
		return []Type{TypeNormal, TypeHideBalance}
	}
	if t.IsDebt() {
		return []Type{TypeDebt}
	}
	if t.IsWishlist() {
		return []Type{TypeWishlist}
	}
	return []Type{}
}

func (t *Type) Valid() bool {
	switch *t {
	case TypeNormal, TypeHideBalance:
		return true
	case TypeDebt:
		return true
	case TypeWishlist:
		return true
	default:
		return false
	}
}

func (t Type) String() string {
	switch t {
	case TypeNormal:
		return "普通钱包"
	case TypeHideBalance:
		return "默认隐藏金额"
	case TypeDebt:
		return "债务"
	case TypeWishlist:
		return "愿望单"
	default:
		return "未知"
	}
}

func (t Type) StringEnglish() string {
	switch t {
	case TypeNormal:
		return "Normal"
	case TypeHideBalance:
		return "HideBalance"
	case TypeDebt:
		return "Debt"
	case TypeWishlist:
		return "Wishlist"
	default:
		return "Unknown"
	}
}

type TypeDetail struct {
	Type        Type   `json:"type"`
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
}

var AllType = []*TypeDetail{
	{TypeNormal, TypeNormal.String(), TypeNormal.StringEnglish()},
	{TypeHideBalance, TypeHideBalance.String(), TypeHideBalance.StringEnglish()},
	{TypeDebt, TypeDebt.String(), TypeDebt.StringEnglish()},
	{TypeWishlist, TypeWishlist.String(), TypeWishlist.StringEnglish()},
}
