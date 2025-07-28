package constant

const (
	AuthHeaderKey          = "Authorization"
	AuthHeaderPrefixBearer = "Bearer "
	AuthHeaderPrefixBasic  = "Basic "
	AuthHeaderPrefixApiKey = "ApiKey "

	UserKey = "User"

	KeyUserTokenPrefix = "MKICS:user:token"

	KeyStaffWeightPrefix   = "MKICS:staff:weight:"
	KeyWecomKFStaffPrefix  = "wecom:kf:staff:"
	KeyWecomBotStaffPrefix = "wecom:bot:staff:"
)

const (
	KFStatusRobotToHuman int = iota + 1
	KFStatusOnlyRobot
	KFStatusOnlyHuman
)

const (
	KFReceiveRuleRoundRobin int = iota + 1
	KFReceiveRuleIdle
)

const (
	KHStatusUnprocessed int = iota + 1
	KHStatusVerification
	KHStatusUserInfoConfirm
)
