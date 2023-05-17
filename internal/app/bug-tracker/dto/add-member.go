package dto

type AddMemberDto struct {
	ProjectID uint64 `json:"projectId"`
	MemberID  uint64 `json:"memberId"`
}
