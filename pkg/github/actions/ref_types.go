package actions

type RefType string

const (
	RefTypeTag    RefType = "tag"
	RefTypeBranch RefType = "branch"
)

func (r RefType) String() string {
	return string(r)
}
