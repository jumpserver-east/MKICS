package constant

type DBContext string

// WebUrlMap 定义前端路由映射
var WebUrlMap = map[string]struct{}{
	"/home":               {},
	"/policy":             {},
	"/staff":              {},
	"/kf":                 {},
	"/llmapp":             {},
	"/policy/add":         {},
	"/staff/add":          {},
	"/kf/add":             {},
	"/llmapp/add":         {},
	"/wecom/receptionist": {},
	"/wecom/account":      {},
	"/wecom/config":       {},
}

// DynamicRoutes 定义需要动态匹配的路由模式
var DynamicRoutes = []string{
	`^/policy/edit/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	`^/staff/edit/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	`^/kf/edit/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	`^/llmapp/edit/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	`^/wecom/account/details/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	`^/wecom/config/edit/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
}

var (
	DB DBContext = "db"
)
