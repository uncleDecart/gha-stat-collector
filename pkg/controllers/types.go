package controllers

const DB_NAME = "records"
const COLLECTION_NAME = "action-logs"

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
