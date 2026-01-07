package enums

type PostStatus string

const (
	DRAFTPostStatus     PostStatus = "DRAFT"
	PUBLISHEDPostStatus PostStatus = "PUBLISHED"
)

func (ps PostStatus) String() string {
	return string(ps)
}
