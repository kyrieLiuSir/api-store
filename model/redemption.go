package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"

	"gorm.io/gorm"
)

type Redemption struct {
	Id                 int               `json:"id"`
	UserId             int               `json:"user_id"`
	Key                string            `json:"key" gorm:"type:char(32);uniqueIndex"`
	Status             int               `json:"status" gorm:"default:1"`
	Name               string            `json:"name" gorm:"index"`
	Quota              int               `json:"quota" gorm:"default:100"`
	Type               string            `json:"type" gorm:"type:varchar(32);default:'quota'"`
	SubscriptionPlanId int               `json:"subscription_plan_id" gorm:"default:0;index"`
	UserSubscriptionId int               `json:"user_subscription_id" gorm:"default:0;index"`
	SubscriptionPlan   *SubscriptionPlan `json:"subscription_plan,omitempty" gorm:"-:all"`
	CreatedTime        int64             `json:"created_time" gorm:"bigint"`
	RedeemedTime       int64             `json:"redeemed_time" gorm:"bigint"`
	Count              int               `json:"count" gorm:"-:all"` // only for api request
	UsedUserId         int               `json:"used_user_id"`
	DeletedAt          gorm.DeletedAt    `gorm:"index"`
	ExpiredTime        int64             `json:"expired_time" gorm:"bigint"` // 过期时间，0 表示不过期
}

type RedemptionResult struct {
	Type             string            `json:"type"`
	Quota            int               `json:"quota,omitempty"`
	RedemptionId     int               `json:"redemption_id"`
	Subscription     *UserSubscription `json:"subscription,omitempty"`
	SubscriptionPlan *SubscriptionPlan `json:"subscription_plan,omitempty"`
}

func normalizeRedemptionType(redemptionType string) string {
	switch strings.TrimSpace(redemptionType) {
	case "", common.RedemptionCodeTypeQuota:
		return common.RedemptionCodeTypeQuota
	case common.RedemptionCodeTypeSubscription:
		return common.RedemptionCodeTypeSubscription
	default:
		return ""
	}
}

func attachSubscriptionPlansToRedemptions(tx *gorm.DB, redemptions []*Redemption) error {
	planIds := make([]int, 0)
	seen := map[int]struct{}{}
	for _, redemption := range redemptions {
		redemption.Type = normalizeRedemptionType(redemption.Type)
		if redemption.Type == common.RedemptionCodeTypeSubscription && redemption.SubscriptionPlanId > 0 {
			if _, ok := seen[redemption.SubscriptionPlanId]; !ok {
				seen[redemption.SubscriptionPlanId] = struct{}{}
				planIds = append(planIds, redemption.SubscriptionPlanId)
			}
		}
	}
	if len(planIds) == 0 {
		return nil
	}
	var plans []SubscriptionPlan
	if err := tx.Where("id IN ?", planIds).Find(&plans).Error; err != nil {
		return err
	}
	planMap := make(map[int]SubscriptionPlan, len(plans))
	for _, plan := range plans {
		planMap[plan.Id] = plan
	}
	for _, redemption := range redemptions {
		if plan, ok := planMap[redemption.SubscriptionPlanId]; ok {
			planCopy := plan
			redemption.SubscriptionPlan = &planCopy
		}
	}
	return nil
}

func GetAllRedemptions(startIdx int, num int) (redemptions []*Redemption, total int64, err error) {
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Redemption{}).Count(&total).Error; err != nil {
			return err
		}
		if err := tx.Order("id desc").Limit(num).Offset(startIdx).Find(&redemptions).Error; err != nil {
			return err
		}
		return attachSubscriptionPlansToRedemptions(tx, redemptions)
	})
	return redemptions, total, err
}

func SearchRedemptions(keyword string, startIdx int, num int) (redemptions []*Redemption, total int64, err error) {
	err = DB.Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&Redemption{})
		if id, err := strconv.Atoi(keyword); err == nil {
			query = query.Where("id = ? OR name LIKE ?", id, keyword+"%")
		} else {
			query = query.Where("name LIKE ?", keyword+"%")
		}
		if err := query.Count(&total).Error; err != nil {
			return err
		}
		if err := query.Order("id desc").Limit(num).Offset(startIdx).Find(&redemptions).Error; err != nil {
			return err
		}
		return attachSubscriptionPlansToRedemptions(tx, redemptions)
	})
	return redemptions, total, err
}

func GetRedemptionById(id int) (*Redemption, error) {
	if id == 0 {
		return nil, errors.New("id 为空！")
	}
	redemption := Redemption{Id: id}
	if err := DB.First(&redemption, "id = ?", id).Error; err != nil {
		return nil, err
	}
	redemption.Type = normalizeRedemptionType(redemption.Type)
	if redemption.Type == common.RedemptionCodeTypeSubscription && redemption.SubscriptionPlanId > 0 {
		plan, err := GetSubscriptionPlanById(redemption.SubscriptionPlanId)
		if err == nil {
			redemption.SubscriptionPlan = plan
		}
	}
	return &redemption, nil
}

func Redeem(key string, userId int) (result *RedemptionResult, err error) {
	if key == "" {
		return nil, errors.New("未提供兑换码")
	}
	if userId == 0 {
		return nil, errors.New("无效的 user id")
	}
	redemption := &Redemption{}
	var subscriptionPlan *SubscriptionPlan
	var upgradeGroup string

	keyCol := "`key`"
	if common.UsingPostgreSQL {
		keyCol = `"key"`
	}
	common.RandomSleep()
	err = DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Set("gorm:query_option", "FOR UPDATE").Where(keyCol+" = ?", key).First(redemption).Error
		if err != nil {
			return errors.New("无效的兑换码")
		}
		redemption.Type = normalizeRedemptionType(redemption.Type)
		if redemption.Type == "" {
			return errors.New("无效的兑换码类型")
		}
		if redemption.Status != common.RedemptionCodeStatusEnabled {
			return errors.New("该兑换码已被使用")
		}
		if redemption.ExpiredTime != 0 && redemption.ExpiredTime < common.GetTimestamp() {
			return errors.New("该兑换码已过期")
		}
		now := common.GetTimestamp()
		switch redemption.Type {
		case common.RedemptionCodeTypeQuota:
			if err = tx.Model(&User{}).Where("id = ?", userId).Update("quota", gorm.Expr("quota + ?", redemption.Quota)).Error; err != nil {
				return err
			}
			redemption.RedeemedTime = now
			redemption.Status = common.RedemptionCodeStatusUsed
			redemption.UsedUserId = userId
			redemption.UserSubscriptionId = 0
			res := tx.Model(&Redemption{}).Where("id = ? AND status = ?", redemption.Id, common.RedemptionCodeStatusEnabled).Updates(map[string]interface{}{
				"type":                 redemption.Type,
				"redeemed_time":        redemption.RedeemedTime,
				"status":               redemption.Status,
				"used_user_id":         redemption.UsedUserId,
				"user_subscription_id": redemption.UserSubscriptionId,
			})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return errors.New("该兑换码已被使用")
			}
			result = &RedemptionResult{
				Type:         common.RedemptionCodeTypeQuota,
				Quota:        redemption.Quota,
				RedemptionId: redemption.Id,
			}
		case common.RedemptionCodeTypeSubscription:
			if redemption.SubscriptionPlanId <= 0 {
				return errors.New("订阅套餐不存在")
			}
			var plan SubscriptionPlan
			if err = tx.Where("id = ?", redemption.SubscriptionPlanId).First(&plan).Error; err != nil {
				return errors.New("订阅套餐不存在")
			}
			if !plan.Enabled {
				return errors.New("订阅套餐已禁用")
			}
			sub, err := CreateUserSubscriptionFromPlanTx(tx, userId, &plan, "redemption")
			if err != nil {
				return err
			}
			redemption.RedeemedTime = now
			redemption.Status = common.RedemptionCodeStatusUsed
			redemption.UsedUserId = userId
			redemption.UserSubscriptionId = sub.Id
			res := tx.Model(&Redemption{}).Where("id = ? AND status = ?", redemption.Id, common.RedemptionCodeStatusEnabled).Updates(map[string]interface{}{
				"type":                 redemption.Type,
				"redeemed_time":        redemption.RedeemedTime,
				"status":               redemption.Status,
				"used_user_id":         redemption.UsedUserId,
				"user_subscription_id": redemption.UserSubscriptionId,
			})
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return errors.New("该兑换码已被使用")
			}
			planCopy := plan
			subscriptionPlan = &planCopy
			upgradeGroup = strings.TrimSpace(plan.UpgradeGroup)
			result = &RedemptionResult{
				Type:             common.RedemptionCodeTypeSubscription,
				RedemptionId:     redemption.Id,
				Subscription:     sub,
				SubscriptionPlan: subscriptionPlan,
			}
		}
		return nil
	})
	if err != nil {
		common.SysError("redemption failed: " + err.Error())
		return nil, ErrRedeemFailed
	}
	if result.Type == common.RedemptionCodeTypeSubscription {
		if upgradeGroup != "" {
			_ = UpdateUserGroupCache(userId, upgradeGroup)
		}
		RecordLog(userId, LogTypeTopup, fmt.Sprintf("通过兑换码开通订阅，套餐: %s，套餐ID %d，订阅ID %d，兑换码ID %d", subscriptionPlan.Title, subscriptionPlan.Id, result.Subscription.Id, redemption.Id))
		return result, nil
	}
	RecordLog(userId, LogTypeTopup, fmt.Sprintf("通过兑换码充值 %s，兑换码ID %d", logger.LogQuota(redemption.Quota), redemption.Id))
	return result, nil
}

func (redemption *Redemption) Insert() error {
	var err error
	err = DB.Create(redemption).Error
	return err
}

func (redemption *Redemption) SelectUpdate() error {
	// This can update zero values
	return DB.Model(redemption).Select("redeemed_time", "status").Updates(redemption).Error
}

// Update Make sure your token's fields is completed, because this will update non-zero values
func (redemption *Redemption) Update() error {
	var err error
	err = DB.Model(redemption).Select("name", "status", "quota", "type", "subscription_plan_id", "user_subscription_id", "redeemed_time", "expired_time").Updates(redemption).Error
	return err
}

func (redemption *Redemption) Delete() error {
	var err error
	err = DB.Delete(redemption).Error
	return err
}

func DeleteRedemptionById(id int) (err error) {
	if id == 0 {
		return errors.New("id 为空！")
	}
	redemption := Redemption{Id: id}
	err = DB.Where(redemption).First(&redemption).Error
	if err != nil {
		return err
	}
	return redemption.Delete()
}

func DeleteInvalidRedemptions() (int64, error) {
	now := common.GetTimestamp()
	result := DB.Where("status IN ? OR (status = ? AND expired_time != 0 AND expired_time < ?)", []int{common.RedemptionCodeStatusUsed, common.RedemptionCodeStatusDisabled}, common.RedemptionCodeStatusEnabled, now).Delete(&Redemption{})
	return result.RowsAffected, result.Error
}
