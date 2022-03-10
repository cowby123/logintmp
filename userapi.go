package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//LoginCheck 確認登入是否成功
func LoginCheck(Data LoginData) (*UserData, error) {
	userdata, err := GetUserDataByID(Data.UserName)

	if err != nil {
		fmt.Println("error at LoginCheck.GetUserDataByID")
		return nil, err
	}
	//================開始檢查密碼=============================
	h := md5.New()
	h.Write([]byte(Data.Password))
	cipherStr := h.Sum(nil)
	if userdata.Password != hex.EncodeToString(cipherStr) {
		//要是密碼不相等
		return nil, errors.New("password error")

	}
	//判斷密碼
	return userdata, nil
}

func registeredAPI(c *gin.Context) {
	var Data UserData
	err := c.BindJSON(&Data)
	if err != nil {
		c.JSON(200, gin.H{"Static": -1, "Message": "PostDataError", "Data": ""})
		return
	}
	//開始檢測帳號是否合格
	err = CheckRegData(Data)
	if err != nil {
		fmt.Println("error at registeredAPI.CheckRegData")
		fmt.Println(err)
		c.JSON(200, gin.H{"Static": -1, "Message": err.Error(), "Data": ""})
		return
	}
	//開始檢測帳號是否有被註冊
	check, err := CheckReg(Data)
	if err != nil {
		fmt.Println("error at registeredAPI.CheckReg")
		fmt.Println(err)
		c.JSON(200, gin.H{"Static": -1, "Message": err.Error(), "Data": ""})
		return
	}
	if !check {
		c.JSON(200, gin.H{"Static": -1, "Message": "UserName is used", "Data": ""})
		return
	}
	//開始檢測信箱是否有被註冊
	check, err = CheckMailReg(Data)
	if err != nil {
		fmt.Println("error at registeredAPI.CheckReg")
		fmt.Println(err)
		c.JSON(200, gin.H{"Static": -1, "Message": err.Error(), "Data": ""})
		return
	}
	if !check {
		c.JSON(200, gin.H{"Static": -1, "Message": "Mail is used", "Data": ""})
		return
	}
	//開始新增資料
	err = RegUser(Data)
	if err != nil {
		c.JSON(200, gin.H{"Static": -1, "Message": err.Error(), "Data": ""})
	}
	c.JSON(200, gin.H{"Static": "1", "Message": "OK", "Data": ""})

}

//CheckRegData 確認帳號資料是否合格
func CheckRegData(Data UserData) error {
	if len(Data.UserName) < 8 {
		return errors.New("username can't <8 ")
	}
	if len(Data.Password) < 8 {
		return errors.New("password can't <8 ")
	}
	return nil
}

//GetUserDataByID 用id從db撈出使用者資料
func GetUserDataByID(ID string) (*UserData, error) {
	//先用username撈資料================================================
	sqlStatement := `SELECT * FROM  UserData WHERE UserName=$1;`
	//var aaaa UserData
	fmt.Println(ID)
	rows, err := db.Query(sqlStatement, ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	userinfo := &UserData{}
	flag := 0
	for rows.Next() {

		if flag == 1 {
			//要是跑到這就代表資料超過一筆母湯
			return nil, errors.New("UserName has so more in db")
		}
		err = rows.Scan(&userinfo.ID, &userinfo.UserName, &userinfo.Password, &userinfo.ChineseName, &userinfo.Email, &userinfo.Address, &userinfo.Phone, &userinfo.Point, &userinfo.State, &userinfo.LastLogin)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		flag = flag + 1
	}
	if flag != 1 {
		//要是跑到這代表沒半筆資料
		return nil, errors.New("no user in db")
	}
	return userinfo, nil
}

//CheckReg 確認帳號是否被註冊過
func CheckReg(Data UserData) (bool, error) {
	sqlStatement := `SELECT COUNT(*) FROM  UserData WHERE UserName=$1;`
	row := db.QueryRow(sqlStatement, Data.UserName)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		fmt.Println("error at CheckReg.Scan")
		return false, err
	}
	//確認是否為0  為0代表沒搜到這行  可以註冊
	if count != 0 {
		return false, nil
	}

	return true, nil
}

//CheckMailReg 確認信箱是否被註冊過
func CheckMailReg(Data UserData) (bool, error) {
	sqlStatement := `SELECT COUNT(*) FROM  UserData WHERE Email=$1;`
	row := db.QueryRow(sqlStatement, Data.Email)
	var count int64
	err := row.Scan(&count)
	if err != nil {
		fmt.Println("error at CheckReg.Scan")
		return false, err
	}
	//確認是否為0  為0代表沒搜到這行  可以註冊
	if count != 0 {
		return false, nil
	}
	return true, nil
}

//RegUser 註冊新使用者
func RegUser(Data UserData) error {
	// 新增資料列：
	stmt, err := db.Prepare("INSERT INTO UserData(UserName,Password,ChineseName,Email, Address ,Phone,  Point,State,LastLogin) VALUES($1,$2,$3,$4,$5,$6,$7,$8,now());")
	if err != nil {
		fmt.Println("error at RegUser.Prepare")
		return err
	}
	h := md5.New()
	h.Write([]byte(Data.Password))
	cipherStr := h.Sum(nil)
	_, err = stmt.Exec(Data.UserName, hex.EncodeToString(cipherStr), Data.ChineseName, Data.Email, Data.Address, Data.Phone, 10000, 0)

	if err != nil {
		fmt.Println("error at RegUser.Exec")
		return err
	}
	return nil
}

func loginAPI(c *gin.Context) {
	var Data LoginData
	err := c.BindJSON(&Data)
	if err != nil {
		c.JSON(200, gin.H{"Static": -1, "Message": "PostDataError", "Data": ""})
		return
	}
	UserData, err := LoginCheck(Data)
	if err != nil {
		c.JSON(200, gin.H{"Static": -1, "Message": err.Error(), "Data": ""})
		return
	}
	generateToken(c, UserData)
	//c.JSON(200, gin.H{"Static": "1", "Message": "OK", "Data": ""})

}

// generateToken 生成token
func generateToken(c *gin.Context, user *UserData) {
	j := &JWT{
		[]byte("legbone"),
	}
	claims := CustomClaims{
		user.ID,
		user.UserName,
		"",
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    "legbone",                       //簽名的發行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}

	log.Println(token)

	data := LoginResult{
		//User:  user,
		Token: token,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"msg":    "Login OK！",
		"data":   data,
	})

}
