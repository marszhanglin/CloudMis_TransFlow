// echoRouteUserManager
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

func echoRouteUserManager() {

	e.POST(projectName+"/user/login", userLogin)
	e.POST(projectName+"/user/registBase", userRegistBase)
	e.POST(projectName+"/user/deleteBase", userDeleteBase)
	e.POST(projectName+"/user/queryBase", userQueryBase)

}

func userLogin(c echo.Context) error {
	glogInfo("----------------------------------------------------------------------------------------------------------")
	c.Request().ParseForm()
	//1. 获取参数
	user := &User{}
	user.UserType = getFormParamInt64(c, "userType") //用户类型
	user.UserName = getFormParam(c, "userName")      //用户名称
	user.Password = getFormParam(c, "password")      //用户密码
	user.ShopId = getFormParamInt64(c, "shopId")     //门店编号（店员必传）

	userJsonBytes, _ := json.Marshal(user)
	glogInfo("Rqdata" + string(userJsonBytes))

	//2. 参数校验
	validateErrs := validate.Struct(user)
	if validateErrs != nil {
		if _, ok := validateErrs.(*validator.InvalidValidationError); !ok {
			glogWarning(validateErrs.Error())
			errResponse := getResPonse("4603", validateErrs.Error())
			return c.JSON(http.StatusOK, errResponse)
		}
	}

	// 3.逻辑处理
	isExit, dbuser := dbQueryUserByName(user.UserName)
	if isExit {
		dbuserJsonBytes, _ := json.Marshal(dbuser)
		glogInfo("Dbdata" + string(dbuserJsonBytes))
	} else {
		glogWarning("用户不存在:" + user.UserName)
		errResponse := getResPonse("4609", "用户不存在:"+user.UserName)
		return c.JSON(http.StatusOK, errResponse)
	}

	// 密码校验
	if user.Password != dbuser.Password {
		glogWarning("登录密码错误")
		errResponse := getResPonse("4609", "登录密码错误")
		return c.JSON(http.StatusOK, errResponse)
	}

	token, uuidErr := uuid.NewV4()
	if nil != uuidErr {
		glogError(uuidErr.Error())
		errResponse := getResPonse("9999", uuidErr.Error())
		return c.JSON(http.StatusOK, errResponse)
	}

	// 用户不同类型处理
	if user.UserType == 1 {
		glogWarning("i01")
	} else if user.UserType == 2 {
		glogWarning("i02")
	} else if user.UserType == 3 {
		// 店长   这台设备被哪个店长登录这台设备就属于哪家店
		// 查询店铺
		isShopExit, dbshop := dbQueryShopById(dbuser.ShopId)
		succResponse := getResPonse("00", "")
		if isShopExit {
			dbshopJsonBytes, _ := json.Marshal(dbshop)
			dbshopJsonStr := string(dbshopJsonBytes)
			glogInfo("Dbdata:" + dbshopJsonStr)
			bodyvalue := make(map[string]string)
			bodyvalue["shopId"] = strconv.FormatInt(dbshop.ShopId, 10)
			bodyvalue["shopName"] = dbshop.ShopName
			bodyvalue["address"] = dbshop.Address
			bodyvalue["merchantId"] = strconv.FormatInt(dbshop.Operator_id, 10)
			bodyvalue["token"] = token.String()
			setNoSqlStrExpire(REDIS_MODLE+"token_"+dbuser.UserName, 7*24*60*60, bodyvalue["token"])
			succResponse.Body = bodyvalue

			return c.JSON(http.StatusOK, succResponse)
		} else {
			succResponse.RetMsg = "门店未配置，请联系管理员"
			succResponse.RetCode = "4610"
			return c.JSON(http.StatusOK, succResponse)
		}
	} else if user.UserType == 4 {
		// 店员登录   判断这个店员是否是这家店的
		// 查询店铺
		isShopExit, dbshop := dbQueryShopById(dbuser.ShopId)
		succResponse := getResPonse("00", "")
		if isShopExit {

			if dbshop.ShopId == user.ShopId {
				dbshopJsonBytes, _ := json.Marshal(dbshop)
				dbshopJsonStr := string(dbshopJsonBytes)
				glogInfo("Dbdata:" + dbshopJsonStr)
				bodyvalue := make(map[string]string)
				bodyvalue["shopId"] = strconv.FormatInt(dbshop.ShopId, 10)
				bodyvalue["shopName"] = dbshop.ShopName
				bodyvalue["address"] = dbshop.Address
				bodyvalue["merchantId"] = strconv.FormatInt(dbshop.Operator_id, 10)
				bodyvalue["token"] = token.String()
				setNoSqlStrExpire(REDIS_MODLE+"token_"+dbuser.UserName, 7*24*60*60, bodyvalue["token"])
				succResponse.Body = bodyvalue
				return c.JSON(http.StatusOK, succResponse)
			} else {
				succResponse.RetMsg = "非门店店员"
				succResponse.RetCode = "4612"
				return c.JSON(http.StatusOK, succResponse)
			}
		} else {
			succResponse.RetMsg = "门店未配置，请联系管理员"
			succResponse.RetCode = "4611"
			return c.JSON(http.StatusOK, succResponse)
		}
	} else {
		glogWarning("用户类型不存在:" + string(user.UserType))

		errResponse := getResPonse("4608", "用户类型不存在:"+string(user.UserType))
		return c.JSON(http.StatusOK, errResponse)
	}

	return c.String(http.StatusOK, "wait。。。")
}

type UserRegistBase_User struct {
	ShopkeeperName string ` validate:"required,max=64"`
	UserName       string ` validate:"required,max=64"`
	Password       string ` validate:"required,min=8,max=64"`
	ShopId         int64  `validate:"-"`
	Token          string ` validate:"required,max=256"`
}

// 添加店员
func userRegistBase(c echo.Context) error {

	glogInfo("----------------------------------------------------------------------------------------------------------")

	c.Request().ParseForm()
	//1. 获取参数
	paramData := &UserRegistBase_User{}
	paramData.ShopkeeperName = getFormParam(c, "shopkeeperName") //店长名称
	paramData.UserName = getFormParam(c, "userName")             //店员名称
	paramData.Password = getFormParam(c, "password")             //店员密码
	paramData.ShopId = getFormParamInt64(c, "shopId")            //门店编号
	paramData.Token = getHeaderParam(c, "token")                 //店长token
	//2. 参数校验
	validateErrs := validate.Struct(paramData)
	if validateErrs != nil {
		if _, ok := validateErrs.(*validator.InvalidValidationError); !ok {
			glogWarning(validateErrs.Error())
			errResponse := getResPonse("4603", validateErrs.Error())
			return c.JSON(http.StatusOK, errResponse)
		}
	}

	// 3.逻辑处理
	// 3.1 校验token
	isTokenValidateSucc := validateToken(paramData.Token, paramData.ShopkeeperName)
	if !isTokenValidateSucc {
		errResponse := getResPonse("4613", "Token过期或未登录")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.2 查询店长信息
	isExit, dbShopkeeperUser := dbQueryUserByName(paramData.ShopkeeperName)
	if isExit {
		dbuserJsonBytes, _ := json.Marshal(dbShopkeeperUser)
		glogInfo("Dbdata" + string(dbuserJsonBytes))
	} else {
		glogError("店员注册，token校验通过，店长不存在:" + paramData.ShopkeeperName)
		glogFlush()
		errResponse := getResPonse("9999", "不存在该店长："+paramData.ShopkeeperName)
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.3 添加者是店员
	if dbShopkeeperUser.UserType == 4 {
		errResponse := getResPonse("4614", "店员无添加店员权限")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.4 添加者没分配门店  添加者关联的门店编号与上送的不符
	if dbShopkeeperUser.ShopId != paramData.ShopId || dbShopkeeperUser.ShopId < 0 {
		errResponse := getResPonse("4615", "门店编号不匹配:["+strconv.FormatInt(dbShopkeeperUser.ShopId, 10)+"],["+strconv.FormatInt(paramData.ShopId, 10)+"],["+paramData.ShopkeeperName+"]")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.5 用户名唯一校验
	isUserNameExit, _ := dbQueryUserByName(paramData.UserName)
	if isUserNameExit {
		glogWarning("用户名重复")
		errResponse := getResPonse("4616", "用户名重复")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.6 保存
	user := &User{}
	user.UserName = paramData.UserName
	user.Password = paramData.Password
	user.UserStatus = 1
	user.UserType = 4
	user.OperatorId = dbShopkeeperUser.OperatorId
	user.ShopId = dbShopkeeperUser.OperatorId
	user.CreTime = time.Now()
	user.UpdTime = time.Now()
	user.Remark = "添加者：" + dbShopkeeperUser.UserName
	saveErr := dbsaveUser(user)
	if nil == saveErr {
		succResponse := getResPonse("00", "succ")
		return c.JSON(http.StatusOK, succResponse)
	} else {
		succResponse := getResPonse("4606", saveErr.Error())
		return c.JSON(http.StatusOK, succResponse)
	}
}

// 删除店员
func userDeleteBase(c echo.Context) error {

	glogInfo("----------------------------------------------------------------------------------------------------------")

	c.Request().ParseForm()
	//1. 获取参数
	paramData := &UserRegistBase_User{}
	paramData.ShopkeeperName = getFormParam(c, "shopkeeperName") //店长名称
	paramData.UserName = getFormParam(c, "userName")             //店员名称
	paramData.Password = "--------"                              //店员密码
	paramData.ShopId = getFormParamInt64(c, "shopId")            //门店编号
	paramData.Token = getHeaderParam(c, "token")                 //店长token
	//2. 参数校验
	validateErrs := validate.Struct(paramData)
	if validateErrs != nil {
		if _, ok := validateErrs.(*validator.InvalidValidationError); !ok {
			glogWarning(validateErrs.Error())
			errResponse := getResPonse("4603", validateErrs.Error())
			return c.JSON(http.StatusOK, errResponse)
		}
	}

	// 3.逻辑处理
	// 3.1 校验token
	isTokenValidateSucc := validateToken(paramData.Token, paramData.ShopkeeperName)
	if !isTokenValidateSucc {
		errResponse := getResPonse("4613", "Token过期或未登录")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.2 查询店长信息
	isExit, dbShopkeeperUser := dbQueryUserByName(paramData.ShopkeeperName)
	if isExit {
		dbuserJsonBytes, _ := json.Marshal(dbShopkeeperUser)
		glogInfo("Dbdata" + string(dbuserJsonBytes))
	} else {
		glogError("店员注册，token校验通过，店长不存在:" + paramData.ShopkeeperName)
		glogFlush()
		errResponse := getResPonse("9999", "不存在该店长："+paramData.ShopkeeperName)
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.3 删除的者是店员
	if dbShopkeeperUser.UserType == 4 {
		errResponse := getResPonse("4614", "店员无添加店员权限")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.4 被删除者没分配门店  被删除者关联的门店编号与上送的不符
	if dbShopkeeperUser.ShopId != paramData.ShopId || dbShopkeeperUser.ShopId < 0 {
		errResponse := getResPonse("4615", "门店编号不匹配:["+strconv.FormatInt(dbShopkeeperUser.ShopId, 10)+"],["+strconv.FormatInt(paramData.ShopId, 10)+"],["+paramData.ShopkeeperName+"]")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.5 用户名唯一校验
	//	isUserNameExit, _ := dbQueryUserByName(paramData.UserName)
	//	if isUserNameExit {
	//		glogWarning("用户名重复")
	//		errResponse := getResPonse("4616", "用户名重复")
	//		return c.JSON(http.StatusOK, errResponse)
	//	}

	// 3.5 删除
	saveErr := dbDeleteUserByName(paramData.UserName)
	if nil == saveErr {
		succResponse := getResPonse("00", "succ")
		return c.JSON(http.StatusOK, succResponse)
	} else {
		succResponse := getResPonse("4615", saveErr.Error())
		return c.JSON(http.StatusOK, succResponse)
	}
}

// 查询店员
func userQueryBase(c echo.Context) error {

	glogInfo("----------------------------------------------------------------------------------------------------------")

	c.Request().ParseForm()
	//1. 获取参数
	paramData := &UserRegistBase_User{}
	paramData.ShopkeeperName = getFormParam(c, "shopkeeperName") //店长名称
	paramData.UserName = "--------"                              //店员名称
	paramData.Password = "--------"                              //店员密码
	paramData.ShopId = getFormParamInt64(c, "shopId")            //门店编号
	paramData.Token = getHeaderParam(c, "token")                 //店长token
	//2. 参数校验
	validateErrs := validate.Struct(paramData)
	if validateErrs != nil {
		if _, ok := validateErrs.(*validator.InvalidValidationError); !ok {
			glogWarning(validateErrs.Error())
			errResponse := getResPonse("4603", validateErrs.Error())
			return c.JSON(http.StatusOK, errResponse)
		}
	}

	// 3.逻辑处理
	// 3.1 校验token
	isTokenValidateSucc := validateToken(paramData.Token, paramData.ShopkeeperName)
	if !isTokenValidateSucc {
		errResponse := getResPonse("4613", "Token过期或未登录")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.2 查询店长信息
	isExit, dbShopkeeperUser := dbQueryUserByName(paramData.ShopkeeperName)
	if isExit {
		dbuserJsonBytes, _ := json.Marshal(dbShopkeeperUser)
		glogInfo("Dbdata" + string(dbuserJsonBytes))
	} else {
		glogError("店员注册，token校验通过，店长不存在:" + paramData.ShopkeeperName)
		glogFlush()
		errResponse := getResPonse("9999", "不存在该店长："+paramData.ShopkeeperName)
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.3  权限
	if dbShopkeeperUser.UserType == 4 {
		errResponse := getResPonse("4614", "店员无查看店员权限")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.4 门店编号与店长是否匹配
	if dbShopkeeperUser.ShopId != paramData.ShopId || dbShopkeeperUser.ShopId < 0 {
		errResponse := getResPonse("4615", "门店编号不匹配:["+strconv.FormatInt(dbShopkeeperUser.ShopId, 10)+"],["+strconv.FormatInt(paramData.ShopId, 10)+"],["+paramData.ShopkeeperName+"]")
		return c.JSON(http.StatusOK, errResponse)
	}

	// 3.5 用户名唯一校验
	//	isUserNameExit, _ := dbQueryUserByName(paramData.UserName)
	//	if isUserNameExit {
	//		glogWarning("用户名重复")
	//		errResponse := getResPonse("4616", "用户名重复")
	//		return c.JSON(http.StatusOK, errResponse)
	//	}

	// 3.5
	users, queryerr := dbQueryUserByShopId(dbShopkeeperUser.ShopId, 4)
	if nil == queryerr {
		succResponse := getBaseResPonse("00", "succ")
		bodyvalue := make(map[string]interface{})
		usersJsonBytes, _ := json.Marshal(users)
		usersJson := string(usersJsonBytes)
		glogInfo(usersJson)
		bodyvalue["shopAssistant"] = users
		succResponse.Body = bodyvalue
		return c.JSON(http.StatusOK, succResponse)
	} else {
		succResponse := getResPonse("4617", queryerr.Error())
		return c.JSON(http.StatusOK, succResponse)
	}
}

func validateToken(requestToken string, userName string) (validateTokenrtn bool) {
	token, _ := getNoSqlStr(REDIS_MODLE + "token_" + userName)
	if len(token) == 0 || token != requestToken {
		glogWarning("Token过期或未登录:[" + token + "],[" + requestToken + "]")
		return false
	} else {
		return true
	}
}
