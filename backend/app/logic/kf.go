package logic

import (
	"MKICS/backend/app/dto"
	"MKICS/backend/app/dto/request"
	"MKICS/backend/app/dto/response"
	"MKICS/backend/app/model"
	"MKICS/backend/global"
)

type KFLogic struct {
}

type IKFLogic interface {
	KFAdd(req request.KF) (err error)
	KFUpdate(uuid string, req request.KF) error
	KFDel(uuid string) error
	KFGet(uuid string) (response.KF, error)
	KFList() ([]response.KF, error)
}

func NewIKFLogic() IKFLogic {
	return &KFLogic{}
}

func (k *KFLogic) KFAdd(req request.KF) (err error) {
	mod := model.KF{
		KFName:           req.KFName,
		KFID:             req.KFID,
		KFPlatform:       req.KFPlatform,
		BotID:            req.BotID,
		Status:           req.Status,
		ReceivePriority:  req.ReceivePriority,
		ReceiveRule:      req.ReceiveRule,
		ChatTimeout:      req.ChatTimeout,
		BotTimeout:       req.BotTimeout,
		BotTimeoutMsg:    req.BotTimeoutMsg,
		BotWelcomeMsg:    req.BotWelcomeMsg,
		StaffWelcomeMsg:  req.StaffWelcomeMsg,
		UnmannedMsg:      req.UnmannedMsg,
		ChatendMsg:       req.ChatendMsg,
		TransferKeywords: req.TransferKeywords,
	}

	mod.Staffs, err = staffRepo.List(commonRepo.WithUUIDsIn(req.StaffList))
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return
	}

	err = kFRepo.Create(mod)
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return
	}

	return
}

func (k *KFLogic) KFUpdate(uuid string, req request.KF) (err error) {
	kf, err := kFRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return err
	}
	kf.UUID = uuid
	kf.KFName = req.KFName
	kf.KFID = req.KFID
	kf.KFPlatform = req.KFPlatform
	kf.BotID = req.BotID
	kf.Status = req.Status
	kf.ReceivePriority = req.ReceivePriority
	kf.ReceiveRule = req.ReceiveRule
	kf.ChatTimeout = req.ChatTimeout
	kf.BotTimeout = req.BotTimeout
	kf.BotTimeoutMsg = req.BotTimeoutMsg
	kf.BotWelcomeMsg = req.BotWelcomeMsg
	kf.StaffWelcomeMsg = req.StaffWelcomeMsg
	kf.UnmannedMsg = req.UnmannedMsg
	kf.ChatendMsg = req.ChatendMsg
	kf.TransferKeywords = req.TransferKeywords
	if req.StaffList != nil {
		staffs, err := staffRepo.List(commonRepo.WithUUIDsIn(req.StaffList))
		if err != nil {
			return err
		}
		kf.Staffs = staffs
	}
	err = kFRepo.Update(kf)
	if err != nil {
		return
	}
	return
}

func (k *KFLogic) KFDel(uuid string) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	if err := kFRepo.Delete(commonRepo.WithByUUID(uuid)); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (k *KFLogic) KFGet(uuid string) (response.KF, error) {
	kf, err := kFRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return response.KF{}, err
	}
	var staffs []response.Staff
	for _, staff := range kf.Staffs {
		var policies []response.Policy
		for _, policy := range staff.Policies {
			var workTimes []dto.WorkTime
			for _, workTime := range policy.WorkTimes {
				workTimes = append(workTimes, dto.WorkTime{
					StartTime: workTime.StartTime,
					EndTime:   workTime.EndTime,
				})
			}
			var policyres response.Policy
			policyres.UUID = policy.UUID
			policyres.PolicyName = policy.PolicyName
			policyres.Repeat = policy.Repeat
			policyres.Week = policy.Week
			policyres.MaxCount = policy.MaxCount
			policyres.WorkTimes = workTimes
			policies = append(policies, policyres)
		}
		var staffres response.Staff
		staffres.UUID = staff.UUID
		staffres.StaffID = staff.StaffID
		staffres.StaffName = staff.StaffName
		staffres.Policies = policies
		staffs = append(staffs, staffres)
	}
	var kfres response.KF
	kfres.UUID = kf.UUID
	kfres.KFName = kf.KFName
	kfres.KFID = kf.KFID
	kfres.KFPlatform = kf.KFPlatform
	kfres.BotID = kf.BotID
	kfres.Status = kf.Status
	kfres.ReceivePriority = kf.ReceivePriority
	kfres.ReceiveRule = kf.ReceiveRule
	kfres.ChatTimeout = kf.ChatTimeout
	kfres.BotTimeout = kf.BotTimeout
	kfres.BotTimeoutMsg = kf.BotTimeoutMsg
	kfres.BotWelcomeMsg = kf.BotWelcomeMsg
	kfres.StaffWelcomeMsg = kf.StaffWelcomeMsg
	kfres.UnmannedMsg = kf.UnmannedMsg
	kfres.ChatendMsg = kf.ChatendMsg
	kfres.TransferKeywords = kf.TransferKeywords
	kfres.Staffs = staffs
	return kfres, nil
}

func (k *KFLogic) KFList() ([]response.KF, error) {
	kfs, err := kFRepo.List()
	if err != nil {
		return nil, err
	}
	var kflist []response.KF
	for _, kf := range kfs {
		var staffs []response.Staff
		for _, staff := range kf.Staffs {
			var policies []response.Policy
			for _, policy := range staff.Policies {
				var workTimes []dto.WorkTime
				for _, workTime := range policy.WorkTimes {
					workTimes = append(workTimes, dto.WorkTime{
						StartTime: workTime.StartTime,
						EndTime:   workTime.EndTime,
					})
				}
				var policyres response.Policy
				policyres.UUID = policy.UUID
				policyres.PolicyName = policy.PolicyName
				policyres.Repeat = policy.Repeat
				policyres.Week = policy.Week
				policyres.MaxCount = policy.MaxCount
				policyres.WorkTimes = workTimes
				policies = append(policies, policyres)
			}
			var staffres response.Staff
			staffres.UUID = staff.UUID
			staffres.StaffID = staff.StaffID
			staffres.StaffName = staff.StaffName
			staffres.Policies = policies
			staffs = append(staffs, staffres)
		}
		var kfres response.KF
		kfres.UUID = kf.UUID
		kfres.KFName = kf.KFName
		kfres.KFID = kf.KFID
		kfres.KFPlatform = kf.KFPlatform
		kfres.BotID = kf.BotID
		kfres.Status = kf.Status
		kfres.ReceivePriority = kf.ReceivePriority
		kfres.ReceiveRule = kf.ReceiveRule
		kfres.ChatTimeout = kf.ChatTimeout
		kfres.BotTimeout = kf.BotTimeout
		kfres.BotTimeoutMsg = kf.BotTimeoutMsg
		kfres.BotWelcomeMsg = kf.BotWelcomeMsg
		kfres.StaffWelcomeMsg = kf.StaffWelcomeMsg
		kfres.UnmannedMsg = kf.UnmannedMsg
		kfres.ChatendMsg = kf.ChatendMsg
		kfres.TransferKeywords = kf.TransferKeywords
		kfres.Staffs = staffs
		kflist = append(kflist, kfres)
	}
	return kflist, nil
}
