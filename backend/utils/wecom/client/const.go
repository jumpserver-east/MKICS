package client

const (
	KeyWecomCursorPrefix = "wecom:cursor:"

	WecomMsgTypeEnterSession           = "enter_session"
	WecomMsgTypeText                   = "text"
	WecomMsgTypeEvent                  = "event"
	WecomEventTypeEnterSession         = "enter_session"
	WecomEventTypeMsgSendFail          = "msg_send_fail"
	WecomEventTypeServicerStatusChange = "servicer_status_change"
	WecomEventTypeSessionStatusChange  = "session_status_change"
	WecomEventTypeUserRecallMsg        = "user_recall_msg"
	WecomEventTypeServicerRecallMsg    = "servicer_recall_msg"

	WecomEventChangeTypeJoinSession     = "1" // 从接待池接入会话
	WecomEventChangeTypeTransferSession = "2" // 转接会话
	WecomEventChangeTypeEndSession      = "3" // 结束会话
	WecomEventChangeTypeRejoinSession   = "4" // 重新接入已结束/已转接会话
)

const (
	SessionStatusNew               int = iota // 0 未处理
	SessionStatusHandled                      // 1 由智能助手接待
	SessionStatusWaiting                      // 2 待接入池排队中
	SessionStatusInProgress                   // 3 由人工接待
	SessionStatusEndedOrNotStarted            // 4 已结束/未开始
)

const (
	MessageTypeCustomer            uint32 = iota + 3 // 3 微信客户发送的消息
	MessageTypeSystemEvent                           // 4 系统推送的事件消息
	MessageTypeReceptionistMessage                   // 5 接待人员在企业微信客户端发送的消息
)
