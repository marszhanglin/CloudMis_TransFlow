// echoRouteSaveflow
package main

import (
	"time"
)

func dbsaveFlow(flow *TransFlow) error {
	//tx := db.Begin()
	flow.Id = getTimeUUID()
	if err := db.Table("trans_flow").Create(&flow).Error; err != nil {
		glogError(err.Error())
		//tx.Rollback()
		return err
	}
	//tx.Commit()
	return nil
}

func dbupdateFlow(flow *TransFlow) error {
	tx := db.Begin()
	if err := tx.Table("trans_flow").Where("id=?", flow.Id).Update(&flow).Error; err != nil {
		glogError(err.Error())
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 根据DeviceSn跟FlowNo查询订单是否重复
func dbisFlowExit(flow *TransFlow) (isExit bool, err error, flowid int64) {

	queryflows := []TransFlow{}
	queryErr := db.Table("trans_flow").Select(" id ").Where(" flow_no = ? and device_sn =?", flow.FlowNo, flow.DeviceSn).Find(&queryflows).Error
	if len(queryflows) > 0 {
		return true, queryErr, queryflows[0].Id
	} else {
		return false, queryErr, 0
	}
}

type TransFlow struct {
	Id             int64     `gorm:"column:id;primary_key" `
	ShopId         int64     `gorm:"column:shop_id" `
	FlowNo         string    `gorm:"column:flow_no" `
	DeviceSn       string    `gorm:"column:device_sn" `
	UploadTime     time.Time `gorm:"column:upload_time" `
	TransTime      time.Time `gorm:"column:trans_time" `
	TransType      string    `gorm:"column:trans_type"`
	ChannelId      string    `gorm:"column:channel_id" `
	MerchantId     string    `gorm:"column:merchant_id" `
	TerminalId     string    `gorm:"column:terminal_id" `
	MerchantName   string    `gorm:"column:merchant_name" `
	Amount         int64     `gorm:"column:amount" `
	TransAmount    int64     `gorm:"column:trans_amount" `
	CurrencyCode   string    `gorm:"column:currency_code" `
	OutOrderNo     string    `gorm:"column:out_order_no" `
	VoucherNo      string    `gorm:"column:voucher_no" `
	ReferenceNo    string    `gorm:"column:reference_no" `
	AuthCode       string    `gorm:"column:auth_code" `
	OriOutOrderNo  string    `gorm:"column:ori_out_order_no"`
	OriVoucherNo   string    `gorm:"column:ori_voucher_no" `
	OriReferenceNo string    `gorm:"column:ori_reference_no" `
	OriAuthCode    string    `gorm:"column:ori_auth_code" `
	CardNo         string    `gorm:"column:card_no" `
	OperatorNo     string    `gorm:"column:operator_no" `
	CombinationNo  string    `gorm:"column:combination_no" `
	CardType       string    `gorm:"column:cardType" `
	Remark         string    `gorm:"column:remark" `
	ExtendParams   string    `gorm:"column:extendParams" `
}
