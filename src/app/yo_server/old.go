package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"yo"
)

type HttpBaseRequest struct {
	retcode  int
	msg      string
	datetime int
}

func testFun(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, %q", r.URL.RawQuery)
	fmt.Println(r.URL.Query())
}

type IError struct {
	RetCode   int    `json:"retcode"`
	Msg       string `json:"msg"`
	ErrorCode string `json:"errorcode"`
}

func makeError(err string) (reg interface{}) {
	reg = IError{
		RetCode:   200,
		Msg:       "error",
		ErrorCode: err,
	}

	return reg
}

///////////////////////////////////////////////////////
/* RegUse */
// 使用方法：客户端请求网址 /user/reg ，务必带上 username 参数
// 使用例子：http://localhost:8000/user/reg?username=1
// 返回 Json
// {"retcode":200,"msg":"ok","datetime":10,"data":{"userid":0,"username":"1"}}
type RegUserData struct {
	Userid   int    `json:"userid"`
	Username string `json:"username"`
}

// 思考怎么实现数据上的继承
type IRegUser struct {
	Retcode  int         `json:"retcode"`
	Msg      string      `json:"msg"`
	Datetime int         `json:"datetime"`
	Data     RegUserData `json:"data"`
}

func RegUser(w http.ResponseWriter, r *http.Request) {
	szUserName := r.URL.Query()["username"][0]
	// fmt.Fprintf(w, "hello, %s", szUserName)

	svr := yo.GetServer()
	err, user := svr.RegistUser(szUserName)
	var reg interface{}
	if err == "success" {
		// 设定 cookie
		// cookiename := "username_" + szUserName
		// cookieid := "userid+" + user.Id
		cookie := http.Cookie{Name: "userid", Value: fmt.Sprintf("%d", user.Id), Path: "/"}
		http.SetCookie(w, &cookie)

		reg = IRegUser{
			Retcode:  200,
			Msg:      "ok",
			Datetime: 10,
			Data:     RegUserData{Userid: user.Id, Username: user.Name},
		}
	} else {
		reg = makeError(err)
	}
	encode := json.NewEncoder(w)
	encode.Encode(reg)
}

///////////////////////////////////////////////////////
/* AddFriend */
// 使用方法：客户端请求网址 /friend/add ，务必带上 friendname 参数
// 使用例子：http://localhost:8000/friend/add?friendname=1&userid=2
// 返回 Json：
// {"retcode":200,"msg":"ok","datetime":10,"data":{"id":1,"name":"1"}}
type AddFriendData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// 思考怎么实现数据上的继承
type IAddFriend struct {
	Retcode  int           `json:"retcode"`
	Msg      string        `json:"msg"`
	Datetime int           `json:"datetime"`
	Data     AddFriendData `json:"data"`
}

func AddFriend(w http.ResponseWriter, r *http.Request) {
	szUserName := r.URL.Query()["friendname"][0]
	szUserId := r.URL.Query()["userid"][0]
	userId, errId := strconv.Atoi(szUserId)
	// fmt.Fprintf(w, "hello, %s", szUserName)
	// szUserId, errId := r.Cookie("userid")
	// userId, _ := strconv.Atoi(szUserId.Value)
	var reg interface{}
	if errId != nil {
		reg = makeError("cookieempty")
	} else {
		err := s.AddFriend(userId, szUserName)
		if err == "success" {
			user := s.GetUserByName(szUserName)
			reg = IAddFriend{
				Retcode:  200,
				Msg:      "ok",
				Datetime: 10,
				Data:     AddFriendData{Id: user.Id, Name: szUserName},
			}
		} else {
			reg = makeError(err)
		}

	}
	encode := json.NewEncoder(w)
	encode.Encode(reg)
}

///////////////////////////////////////////////////////
/* GetFriendList */
// 使用方法：客户端请求网址 /friend/list，务必带上 userid 参数
// 使用例子：http://localhost:8000/friend/list?userid=1
// 返回 Json:
// {"retcode":200,"msg":"ok","datetime":10,"data":{"count":1,"friend":[{"name":"2","id":0}]}}
type User struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type GetFriendListData struct {
	Count  int    `json:"count"`
	Friend []User `json:"friend"`
}

// 思考怎么实现数据上的继承
type IGetFriendList struct {
	Retcode  int               `json:"retcode"`
	Msg      string            `json:"msg"`
	Datetime int               `json:"datetime"`
	Data     GetFriendListData `json:"data"`
}

func GetFriendList(w http.ResponseWriter, r *http.Request) {
	/* fmt.Fprintf(w, "hello, %s", szUserName) */
	// szUserId, errId := r.Cookie("userid")
	szUserId := r.URL.Query()["userid"][0]
	userId, errId := strconv.Atoi(szUserId)
	// userId, _ := strconv.Atoi(szUserId.Value)
	var reg interface{}
	if errId != nil {
		reg = makeError("cookieempty")
	} else {
		err, lstUser := s.GetFriendList(userId)
		if err == "success" {
			var lstSendData GetFriendListData
			lstSendData.Count = len(lstUser)

			for _, value := range lstUser {
				var u User
				u.Name = value.Name
				u.Id = value.Id
				lstSendData.Friend = append(lstSendData.Friend, u)
			}

			reg = IGetFriendList{
				Retcode:  200,
				Msg:      "ok",
				Datetime: 10,
				Data:     lstSendData,
			}
		} else {
			reg = makeError(err)
		}
	}

	encode := json.NewEncoder(w)
	encode.Encode(reg)
}

///////////////////////////////////////////////////////
/* SendYO */
// 使用方法：客户端请求网址 /yo/sendyo ，务必带上 friendid
// 使用例子：http://localhost:8000/yo/sendyo?friendid=1&userid=1
// 返回 Json：
// {"retcode":200,"msg":"ok","datetime":10}
// 思考怎么实现数据上的继承
type ISendYO struct {
	Retcode  int    `json:"retcode"`
	Msg      string `json:"msg"`
	Datetime int    `json:"datetime"`
}

func SendYO(w http.ResponseWriter, r *http.Request) {
	szGeterId := r.URL.Query()["friendid"][0]
	geterId, _ := strconv.Atoi(szGeterId)

	/* fmt.Fprintf(w, "hello, %s", szUserName) */
	// szUserId, errId := r.Cookie("userid")
	szUserId := r.URL.Query()["userid"][0]
	userId, errId := strconv.Atoi(szUserId)
	// userId, _ := strconv.Atoi(szUserId.Value)
	var reg interface{}
	if errId != nil {
		reg = makeError("cookieempty")
	} else {
		err := s.SendYO(userId, geterId)
		if err == "success" {
			reg = ISendYO{
				Retcode:  200,
				Msg:      "ok",
				Datetime: 10,
			}
		} else {
			reg = makeError(err)
		}
	}
	encode := json.NewEncoder(w)
	encode.Encode(reg)
}

///////////////////////////////////////////////////////
/* GetYO */
// 使用方法：客户端请求网址 /yo/getyo，务必带上 userid 参数
// 使用例子：http://localhost:8000/yo/getyo?userid=1
// {"retcode":200,"msg":"ok","datetime":10,"data":{"count":1,"msgs":[{"from":"1","msg":"","senddate":0}]}}
type YOMsg struct {
	From     string `json:"from"`
	Msg      string `json:"msg"`
	Senddate int    `json:"senddate"`
}

type GetYOData struct {
	Count int     `json:"count"`
	Msgs  []YOMsg `json:"msgs"`
}

// 思考怎么实现数据上的继承
type IGetYO struct {
	Retcode  int       `json:"retcode"`
	Msg      string    `json:"msg"`
	Datetime int       `json:"datetime"`
	Data     GetYOData `json:"data"`
}

func GetYO(w http.ResponseWriter, r *http.Request) {
	/* fmt.Fprintf(w, "hello, %s", szUserName) */
	// szUserId, errId := r.Cookie("userid")
	szUserId := r.URL.Query()["userid"][0]
	userId, errId := strconv.Atoi(szUserId)
	// userId, _ := strconv.Atoi(szUserId.Value)
	var reg interface{}
	if errId != nil {
		reg = makeError("cookieempty")
	} else {
		err, lstYO := s.GetYO(userId)
		if err == "success" {
			var lstData GetYOData
			lstData.Count = len(lstYO)
			for _, value := range lstYO {
				var y YOMsg
				err, y.From = s.GetUserName(value.SenderId)
				if err == "success" {
					lstData.Msgs = append(lstData.Msgs, y)
				} /*else {
					// reg = makeError(err)
					// break
				}*/
			}
			reg = IGetYO{
				Retcode:  200,
				Msg:      "ok",
				Datetime: 10,
				Data:     lstData,
			}
		} else {
			reg = makeError(err)
		}
	}

	encode := json.NewEncoder(w)
	encode.Encode(reg)
}

func main() {
	// HtmlServer := http.FileServer(http.Dir("."))
	// http.Handle("/", HtmlServer)
	// http.HandleFunc("/test", testFun)

	// http.HandleFunc("/user/reg", RegUser)
	// http.HandleFunc("/friend/add", AddFriend)
	// http.HandleFunc("/friend/list", GetFriendList)
	// http.HandleFunc("/yo/sendyo", SendYO)
	// http.HandleFunc("/yo/getyo", GetYO)

	// // 缺少一个默认的 404

	// err := http.ListenAndServe(":8000", nil)
	// Base.CheckErr(err)

	yo.Init()
	yo.StartServer()
}
