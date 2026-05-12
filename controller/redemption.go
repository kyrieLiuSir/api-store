package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/i18n"
	"github.com/QuantumNous/new-api/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllRedemptions(c *gin.Context) {
	pageInfo := common.GetPageQuery(c)
	redemptions, total, err := model.GetAllRedemptions(pageInfo.GetStartIdx(), pageInfo.GetPageSize())
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(redemptions)
	common.ApiSuccess(c, pageInfo)
	return
}

func SearchRedemptions(c *gin.Context) {
	keyword := c.Query("keyword")
	pageInfo := common.GetPageQuery(c)
	redemptions, total, err := model.SearchRedemptions(keyword, pageInfo.GetStartIdx(), pageInfo.GetPageSize())
	if err != nil {
		common.ApiError(c, err)
		return
	}
	pageInfo.SetTotal(int(total))
	pageInfo.SetItems(redemptions)
	common.ApiSuccess(c, pageInfo)
	return
}

func GetRedemption(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ApiError(c, err)
		return
	}
	redemption, err := model.GetRedemptionById(id)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    redemption,
	})
	return
}

func AddRedemption(c *gin.Context) {
	redemption := model.Redemption{}
	err := c.ShouldBindJSON(&redemption)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if err := validateRedemptionPayload(c, &redemption, false); err != nil {
		common.ApiErrorMsg(c, err.Error())
		return
	}
	keys := make([]string, 0, redemption.Count)
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < redemption.Count; i++ {
			key := common.GetUUID()
			cleanRedemption := model.Redemption{
				UserId:             c.GetInt("id"),
				Name:               redemption.Name,
				Key:                key,
				CreatedTime:        common.GetTimestamp(),
				Quota:              redemption.Quota,
				Type:               redemption.Type,
				SubscriptionPlanId: redemption.SubscriptionPlanId,
				ExpiredTime:        redemption.ExpiredTime,
			}
			if err := tx.Create(&cleanRedemption).Error; err != nil {
				return err
			}
			keys = append(keys, key)
		}
		return nil
	})
	if err != nil {
		common.SysError("failed to insert redemption: " + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": i18n.T(c, i18n.MsgRedemptionCreateFailed),
			"data":    keys,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    keys,
	})
	return
}

func DeleteRedemption(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := model.DeleteRedemptionById(id)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
	return
}

func UpdateRedemption(c *gin.Context) {
	statusOnly := c.Query("status_only")
	redemption := model.Redemption{}
	err := c.ShouldBindJSON(&redemption)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	cleanRedemption, err := model.GetRedemptionById(redemption.Id)
	if err != nil {
		common.ApiError(c, err)
		return
	}
	if statusOnly == "" {
		if cleanRedemption.Status == common.RedemptionCodeStatusUsed {
			common.ApiErrorMsg(c, "已使用的兑换码不能修改")
			return
		}
		if err := validateRedemptionPayload(c, &redemption, true); err != nil {
			common.ApiErrorMsg(c, err.Error())
			return
		}
		cleanRedemption.Name = redemption.Name
		cleanRedemption.Quota = redemption.Quota
		cleanRedemption.Type = redemption.Type
		cleanRedemption.SubscriptionPlanId = redemption.SubscriptionPlanId
		cleanRedemption.ExpiredTime = redemption.ExpiredTime
	}
	if statusOnly != "" {
		cleanRedemption.Status = redemption.Status
	}
	err = cleanRedemption.Update()
	if err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    cleanRedemption,
	})
	return
}

func DeleteInvalidRedemption(c *gin.Context) {
	rows, err := model.DeleteInvalidRedemptions()
	if err != nil {
		common.ApiError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    rows,
	})
	return
}

func validateRedemptionPayload(c *gin.Context, redemption *model.Redemption, isEdit bool) error {
	if utf8.RuneCountInString(redemption.Name) == 0 || utf8.RuneCountInString(redemption.Name) > 20 {
		return errors.New(i18n.T(c, i18n.MsgRedemptionNameLength))
	}
	if !isEdit {
		if redemption.Count <= 0 {
			return errors.New(i18n.T(c, i18n.MsgRedemptionCountPositive))
		}
		if redemption.Count > 100 {
			return errors.New(i18n.T(c, i18n.MsgRedemptionCountMax))
		}
	}
	if valid, msg := validateExpiredTime(c, redemption.ExpiredTime); !valid {
		return errors.New(msg)
	}
	redemption.Type = normalizeRedemptionType(redemption.Type)
	if redemption.Type == "" {
		return errors.New("兑换码类型无效")
	}
	if redemption.Type == common.RedemptionCodeTypeSubscription {
		if redemption.SubscriptionPlanId <= 0 {
			return errors.New("请选择订阅套餐")
		}
		plan, err := model.GetSubscriptionPlanById(redemption.SubscriptionPlanId)
		if err != nil {
			return errors.New("订阅套餐不存在")
		}
		if !plan.Enabled {
			return errors.New("订阅套餐已禁用")
		}
		redemption.Quota = 0
		return nil
	}
	if redemption.Quota <= 0 {
		return errors.New("额度必须大于0")
	}
	redemption.SubscriptionPlanId = 0
	return nil
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

func validateExpiredTime(c *gin.Context, expired int64) (bool, string) {
	if expired != 0 && expired < common.GetTimestamp() {
		return false, i18n.T(c, i18n.MsgRedemptionExpireTimeInvalid)
	}
	return true, ""
}
