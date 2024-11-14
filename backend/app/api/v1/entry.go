package v1

import "EvoBot/backend/app/logic"

type ApiGroup struct {
	BaseApi
}

var ApiGroupApp = new(ApiGroup)

var (
	authLogic = logic.NewIAuthLogic()

	wecomLogic  = logic.NewIWecomLogic()
	maxkbLogic  = logic.NewIMaxkbLogic()
	policyLogic = logic.NewIPolicyLogic()
	staffLogic  = logic.NewIStaffLogic()
	kFLogic     = logic.NewIKFLogic()
)
