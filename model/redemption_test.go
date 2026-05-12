package model

import (
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func insertRedemptionTestUser(t *testing.T, id int, group string, quota int) {
	t.Helper()
	user := &User{
		Id:       id,
		Username: "redemption_user_" + common.GetRandomString(8),
		Status:   common.UserStatusEnabled,
		Group:    group,
		Quota:    quota,
		AffCode:  "aff_" + common.GetRandomString(16),
	}
	require.NoError(t, DB.Create(user).Error)
}

func insertRedemptionTestPlan(t *testing.T, id int, enabled bool) *SubscriptionPlan {
	t.Helper()
	plan := &SubscriptionPlan{
		Id:               id,
		Title:            "Redemption Plan",
		PriceAmount:      0,
		Currency:         "USD",
		DurationUnit:     SubscriptionDurationMonth,
		DurationValue:    1,
		Enabled:          true,
		TotalAmount:      1000,
		QuotaResetPeriod: SubscriptionResetDaily,
	}
	require.NoError(t, DB.Create(plan).Error)
	if !enabled {
		require.NoError(t, DB.Model(plan).Update("enabled", false).Error)
		plan.Enabled = false
	}
	return plan
}

func insertRedemptionTestCode(t *testing.T, code *Redemption) *Redemption {
	t.Helper()
	if code.Key == "" {
		code.Key = common.GetUUID()
	}
	if code.Status == 0 {
		code.Status = common.RedemptionCodeStatusEnabled
	}
	if code.CreatedTime == 0 {
		code.CreatedTime = common.GetTimestamp()
	}
	if code.Name == "" {
		code.Name = "test-code"
	}
	require.NoError(t, DB.Create(code).Error)
	return code
}

func getRedemptionTestUserQuota(t *testing.T, userId int) int {
	t.Helper()
	var user User
	require.NoError(t, DB.Select("quota").Where("id = ?", userId).First(&user).Error)
	return user.Quota
}

func getRedemptionTestUserGroup(t *testing.T, userId int) string {
	t.Helper()
	var user User
	require.NoError(t, DB.Select("group").Where("id = ?", userId).First(&user).Error)
	return user.Group
}

func getRedemptionTestCode(t *testing.T, key string) Redemption {
	t.Helper()
	var code Redemption
	require.NoError(t, DB.Where("`key` = ?", key).First(&code).Error)
	return code
}

func countRedemptionTestSubscriptions(t *testing.T, userId int, planId int) int64 {
	t.Helper()
	var count int64
	require.NoError(t, DB.Model(&UserSubscription{}).Where("user_id = ? AND plan_id = ?", userId, planId).Count(&count).Error)
	return count
}

func TestRedeem_HistoricalQuotaCodeDefaultsToQuota(t *testing.T) {
	truncateTables(t)

	insertRedemptionTestUser(t, 1001, "default", 10)
	code := insertRedemptionTestCode(t, &Redemption{Quota: 25})

	result, err := Redeem(code.Key, 1001)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, common.RedemptionCodeTypeQuota, result.Type)
	assert.Equal(t, 25, result.Quota)
	assert.Equal(t, 35, getRedemptionTestUserQuota(t, 1001))
	updated := getRedemptionTestCode(t, code.Key)
	assert.Equal(t, common.RedemptionCodeStatusUsed, updated.Status)
	assert.Equal(t, 1001, updated.UsedUserId)
	assert.NotZero(t, updated.RedeemedTime)
	assert.Zero(t, countRedemptionTestSubscriptions(t, 1001, 0))
}

func TestRedeem_SubscriptionCodeCreatesSubscription(t *testing.T) {
	truncateTables(t)

	insertRedemptionTestUser(t, 1002, "default", 10)
	plan := insertRedemptionTestPlan(t, 2001, true)
	code := insertRedemptionTestCode(t, &Redemption{
		Type:               common.RedemptionCodeTypeSubscription,
		Quota:              999,
		SubscriptionPlanId: plan.Id,
	})

	result, err := Redeem(code.Key, 1002)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.Subscription)
	require.NotNil(t, result.SubscriptionPlan)
	assert.Equal(t, common.RedemptionCodeTypeSubscription, result.Type)
	assert.Equal(t, plan.Id, result.Subscription.PlanId)
	assert.Equal(t, int64(1000), result.Subscription.AmountTotal)
	assert.Equal(t, int64(0), result.Subscription.AmountUsed)
	assert.Equal(t, "redemption", result.Subscription.Source)
	assert.Equal(t, "active", result.Subscription.Status)
	assert.NotZero(t, result.Subscription.NextResetTime)
	assert.Equal(t, 10, getRedemptionTestUserQuota(t, 1002))
	updated := getRedemptionTestCode(t, code.Key)
	assert.Equal(t, common.RedemptionCodeStatusUsed, updated.Status)
	assert.Equal(t, 1002, updated.UsedUserId)
	assert.Equal(t, result.Subscription.Id, updated.UserSubscriptionId)
	assert.Equal(t, int64(1), countRedemptionTestSubscriptions(t, 1002, plan.Id))
}

func TestRedeem_SubscriptionCodeCannotBeReused(t *testing.T) {
	truncateTables(t)

	insertRedemptionTestUser(t, 1003, "default", 0)
	insertRedemptionTestUser(t, 1004, "default", 0)
	plan := insertRedemptionTestPlan(t, 2002, true)
	code := insertRedemptionTestCode(t, &Redemption{
		Type:               common.RedemptionCodeTypeSubscription,
		SubscriptionPlanId: plan.Id,
	})

	_, err := Redeem(code.Key, 1003)
	require.NoError(t, err)
	_, err = Redeem(code.Key, 1004)
	require.ErrorIs(t, err, ErrRedeemFailed)
	assert.Equal(t, int64(1), countRedemptionTestSubscriptions(t, 1003, plan.Id))
	assert.Equal(t, int64(0), countRedemptionTestSubscriptions(t, 1004, plan.Id))
	updated := getRedemptionTestCode(t, code.Key)
	assert.Equal(t, 1003, updated.UsedUserId)
}

func TestRedeem_DisabledPlanLeavesCodeUnused(t *testing.T) {
	truncateTables(t)

	insertRedemptionTestUser(t, 1005, "default", 0)
	plan := insertRedemptionTestPlan(t, 2003, false)
	code := insertRedemptionTestCode(t, &Redemption{
		Type:               common.RedemptionCodeTypeSubscription,
		SubscriptionPlanId: plan.Id,
	})

	_, err := Redeem(code.Key, 1005)
	require.ErrorIs(t, err, ErrRedeemFailed)
	updated := getRedemptionTestCode(t, code.Key)
	assert.Equal(t, common.RedemptionCodeStatusEnabled, updated.Status)
	assert.Zero(t, updated.UsedUserId)
	assert.Zero(t, updated.UserSubscriptionId)
	assert.Equal(t, int64(0), countRedemptionTestSubscriptions(t, 1005, plan.Id))
}

func TestRedeem_MaxPurchaseLimitLeavesCodeUnused(t *testing.T) {
	truncateTables(t)

	insertRedemptionTestUser(t, 1006, "default", 0)
	plan := insertRedemptionTestPlan(t, 2004, true)
	plan.MaxPurchasePerUser = 1
	require.NoError(t, DB.Save(plan).Error)
	require.NoError(t, DB.Transaction(func(tx *gorm.DB) error {
		_, err := CreateUserSubscriptionFromPlanTx(tx, 1006, plan, "admin")
		return err
	}))
	code := insertRedemptionTestCode(t, &Redemption{
		Type:               common.RedemptionCodeTypeSubscription,
		SubscriptionPlanId: plan.Id,
	})

	_, err := Redeem(code.Key, 1006)
	require.ErrorIs(t, err, ErrRedeemFailed)
	updated := getRedemptionTestCode(t, code.Key)
	assert.Equal(t, common.RedemptionCodeStatusEnabled, updated.Status)
	assert.Zero(t, updated.UsedUserId)
	assert.Equal(t, int64(1), countRedemptionTestSubscriptions(t, 1006, plan.Id))
}

func TestRedeem_SubscriptionCodeUpgradesUserGroup(t *testing.T) {
	truncateTables(t)

	insertRedemptionTestUser(t, 1007, "default", 0)
	plan := insertRedemptionTestPlan(t, 2005, true)
	plan.UpgradeGroup = "vip"
	require.NoError(t, DB.Save(plan).Error)
	code := insertRedemptionTestCode(t, &Redemption{
		Type:               common.RedemptionCodeTypeSubscription,
		SubscriptionPlanId: plan.Id,
	})

	_, err := Redeem(code.Key, 1007)
	require.NoError(t, err)
	assert.Equal(t, "vip", getRedemptionTestUserGroup(t, 1007))
	var sub UserSubscription
	require.NoError(t, DB.Where("user_id = ? AND plan_id = ?", 1007, plan.Id).First(&sub).Error)
	assert.Equal(t, "vip", sub.UpgradeGroup)
	assert.Equal(t, "default", sub.PrevUserGroup)
}

func TestResetDueSubscriptions_ResetsRedeemedSubscriptionUsage(t *testing.T) {
	truncateTables(t)

	insertRedemptionTestUser(t, 1008, "default", 0)
	plan := insertRedemptionTestPlan(t, 2006, true)
	code := insertRedemptionTestCode(t, &Redemption{
		Type:               common.RedemptionCodeTypeSubscription,
		SubscriptionPlanId: plan.Id,
	})
	result, err := Redeem(code.Key, 1008)
	require.NoError(t, err)
	sub := result.Subscription
	sub.AmountUsed = 500
	sub.LastResetTime = time.Now().Add(-25 * time.Hour).Unix()
	sub.NextResetTime = time.Now().Add(-1 * time.Hour).Unix()
	require.NoError(t, DB.Save(sub).Error)

	count, err := ResetDueSubscriptions(10)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	var updated UserSubscription
	require.NoError(t, DB.Where("id = ?", sub.Id).First(&updated).Error)
	assert.Zero(t, updated.AmountUsed)
	assert.Greater(t, updated.LastResetTime, sub.LastResetTime)
}
