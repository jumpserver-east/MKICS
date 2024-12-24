package constant

const (
	TokenKey = "Authorization"

	KeyUserTokenPrefix = "evobot:user:token:"
	KeyPostTimeZSet    = "evobot:post:time"

	KeyStaffWeightPrefix  = "evobot:staff:weight:"
	KeyWecomKHStaffPrefix = "wecom:kh:staff:"
)

const (
	KFStatusRobotToHuman int = iota + 1 // 1 机器人可转人工
	KFStatusOnlyRobot                   // 2 仅机器人
	KFStatusOnlyHuman                   // 3 仅人工
)

const (
	KFReceiveRuleRoundRobin int = iota + 1 // 1 轮流接待
	KFReceiveRuleIdle                      // 2 空闲接待
)
