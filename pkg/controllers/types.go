package controllers

const DbName = "records"
const CollectionName = "action-logs"

type StepLogEntry struct {
	Id       uint64 `json:"id" binding:"required"`
	ExecTime string `json:"exec_time" binding:"required"`
	// why *bool: see https://github.com/gin-gonic/gin/issues/814
	Successfull *bool `json:"successful" binding:"required"`
}

type ActionLogEntry struct {
	Name  string `json:"name" binding:"required"`
	Start string `json:"start" binding:"required"`
	End   string `json:"end" binding:"required"`
	// why *bool: see https://github.com/gin-gonic/gin/issues/814
	Successfull *bool          `json:"successful" binding:"required"`
	Arch        string         `json:"arch" binding:"required"`
	Steps       []StepLogEntry `json:"steps"`
}

type ActionLogEntrySearch struct {
	PerPage    int64            `json:"perPage"`
	TotalPages int64            `json:"totalPages"`
	CurPage    int64            `json:"curPage"`
	Data       []ActionLogEntry `json:"data"`
}

func (lhs StepLogEntry) Compare(rhs StepLogEntry) bool {
	return lhs.Id == rhs.Id &&
		lhs.ExecTime == rhs.ExecTime &&
		*lhs.Successfull == *rhs.Successfull
}

func CompareStepLogEntrySlice(lhs, rhs []StepLogEntry) bool {
	lhsLen := len(lhs)
	if lhsLen != len(rhs) {
		return false
	}
	for i := 0; i < lhsLen; i++ {
		if !lhs[i].Compare(rhs[i]) {
			return false
		}
	}
	return true
}

func (lhs ActionLogEntry) Compare(rhs ActionLogEntry) bool {
	if lhs.Name != rhs.Name {
		return false
	}
	if lhs.Start != rhs.Start {
		return false
	}
	if lhs.End != rhs.End {
		return false
	}
	if *lhs.Successfull != *rhs.Successfull {
		return false
	}
	if lhs.Arch != rhs.Arch {
		return false
	}
	if !CompareStepLogEntrySlice(lhs.Steps, rhs.Steps) {
		return false
	}
	return true
}

func CompareActionLogEntrySlice(lhs, rhs []ActionLogEntry) bool {
	lenLhs := len(lhs)
	if lenLhs != len(rhs) {
		return false
	}

	for i := 0; i < lenLhs; i++ {
		if !lhs[i].Compare(rhs[i]) {
			return false
		}
	}

	return true
}
