package model

type KF struct {
	BaseModel
	KFName           string  `gorm:"column:kfname;type:varchar(255);" json:"kfname"`
	KFID             string  `gorm:"column:kfid;type:varchar(255);" json:"kfid"`
	KFPlatform       string  `gorm:"column:kfplatform;type:varchar(255);" json:"kfplatform"`
	BotID            string  `gorm:"column:botid;type:varchar(255);" json:"botid"`
	BotPlatfrom      string  `gorm:"column:botplatform;type:varchar(255);" json:"botplatform"`
	Status           int     `gorm:"column:status;type:int;" json:"status"`                                // 1机器人可转人工;2仅机器人;3仅人工
	ReceivePriority  int     `gorm:"column:receive_priority;type:int;" json:"receive_priority"`            // 是否优先上一个
	ReceiveRule      int     `gorm:"column:receive_rule;type:int;" json:"receive_rule"`                    // 1轮流接待，2空闲接待
	ChatTimeout      int     `gorm:"column:chat_timeout;" json:"chat_timeout"`                             // 会话超时时间
	BotTimeout       int     `gorm:"column:bot_timeout;" json:"bot_timeout"`                               // 机器人超时时间
	BotTimeoutMsg    string  `gorm:"column:bot_timeout_msg;type:varchar(255);" json:"bot_timeout_msg"`     // 机器人超时提示语
	BotWelcomeMsg    string  `gorm:"column:bot_welcome_msg;type:varchar(255);" json:"botwelcome_msg"`      // 机器人的欢迎语
	StaffWelcomeMsg  string  `gorm:"column:staff_welcome_msg;type:varchar(255);" json:"staff_welcome_msg"` // 转人工的欢迎语
	UnmannedMsg      string  `gorm:"column:unmanned_msg;type:varchar(255);" json:"unmanned_msg"`           // 无人值守语
	ChatendMsg       string  `gorm:"column:chatend_msg;type:varchar(255);" json:"chatend_msg"`             // 会话结束语
	TransferKeywords string  `gorm:"column:transfer_keywords;type:varchar(255);" json:"transfer_keywords"` // 转人工关键字列表
	Staffs           []Staff `gorm:"many2many:kf_staff;foreignKey:ID;joinForeignKey:KFID;References:ID;joinReferences:StaffID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;" json:"kf_staffs"`
}
