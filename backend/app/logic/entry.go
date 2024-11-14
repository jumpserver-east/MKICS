package logic

import "EvoBot/backend/app/repo"

var (
	commonRepo = repo.NewCommonRepo()

	authRepo   = repo.NewISettingRepo()
	wecomRepo  = repo.NewIWecomRepo()
	maxkbRepo  = repo.NewIMaxkbRepo()
	kHRepo     = repo.NewIKHRepo()
	kFRepo     = repo.NewIKFRepo()
	policyRepo = repo.NewIPolicyRepo()
	staffRepo  = repo.NewIStaffRepo()
)
