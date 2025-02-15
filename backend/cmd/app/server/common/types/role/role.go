package role

type Type string

const (
	TypeAdmin  Type = "admin"
	TypeUser   Type = "user"
	TypeViewer Type = "viewer"
)

func (t Type) Valid() bool {
	switch t {
	case TypeAdmin, TypeUser, TypeViewer:
		return true
	default:
		return false
	}
}

func (t Type) String() string {
	switch t {
	case TypeAdmin:
		return "管理员"
	case TypeUser:
		return "用户"
	case TypeViewer:
		return "访客"
	default:
		return "未知"
	}
}

func (t Type) StringEnglish() string {
	switch t {
	case TypeAdmin:
		return "admin"
	case TypeUser:
		return "user"
	case TypeViewer:
		return "viewer"
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
	{TypeAdmin, TypeAdmin.String(), TypeAdmin.StringEnglish()},
	{TypeUser, TypeUser.String(), TypeUser.StringEnglish()},
	{TypeViewer, TypeViewer.String(), TypeViewer.StringEnglish()},
}
