package module

import (
	"fmt"
)

type ModuleServer struct {
	// 一组用户列表
	LstUser []*UserData
	// 一组用户关系列表
	LstContact []*ContactInfo
	// 一组用户消息列表
	LstMsg []*MsgInfo
}

/*
函数名：根据用户 Id 获得用户名

参数：id

返回值：err name

"success" 取值成功
"emptyid" 用户为空
*/
func (s *ModuleServer) GetUserName(id int) (err string, name string) {
	if id < 0 || id >= len(s.LstUser) {
		return "emptyid", name
	}

	return "success", s.LstUser[id].Name
}

/*
函数名：根据用户名字获得用户

参数：注册用户名

返回值：u，当 u 为 nil 得时候，表示没有找到 UserData
*/
func (s *ModuleServer) GetUserByName(username string) (u *UserData) {
	for _, user := range s.LstUser {
		if user.Name == username {
			return user
		}
	}

	return nil
}

/*
函数名：注册用户

参数：注册用户名称

返回值：err u

err : 错误信息
"success" 创建成功
"nameexist" 重名

u : 一个用户实例，创建不成功时，其值为 Null
*/
func (s *ModuleServer) RegistUser(username string) (err string, u *UserData) {
	if s.GetUserByName(username) != nil {
		return "nameexist", nil
	}

	user := new(UserData)
	user.Id = len(s.LstUser)
	user.Name = username

	s.LstUser = append(s.LstUser, user)
	return "success", user
}

func (s *ModuleServer) isFriend(user1Id int, user2Id int) (bIsFriend bool) {
	bIsFriend = false

	for _, contactinfo := range s.LstContact {
		var bHasUser1, bHasUser2 bool
		bHasUser1 = false
		bHasUser2 = false

		for _, userId := range contactinfo.LstUserId {
			if bHasUser1 == false && userId == user1Id {
				bHasUser1 = true
			}

			if bHasUser2 == false && userId == user2Id {
				bHasUser2 = true
			}

			if bHasUser1 && bHasUser2 {
				bIsFriend = true
			}
		}
	}

	return bIsFriend
}

/*
函数名：增加新好友

参数：id username

id：提起加好友请求者

username：被邀请好友对象

返回值：err

err : 错误信息
"success":增加成功
"iduserempty":加好友请求者为空
"idnameissameone":自己加自己
"nameuserempty":被请求者为空
"alreadyfriend":两人已经是好友
*/
func (s *ModuleServer) AddFriend(id int, username string) (err string) {
	if id < 0 || len(s.LstUser) <= id {
		return "iduserempty"
	}

	user1 := s.LstUser[id]
	if user1 == nil {
		return "iduserempty"
	}

	user2 := s.GetUserByName(username)
	if user2 == nil {
		return "nameuserempty"
	}

	if user1.Id == user2.Id {
		fmt.Println(user1.Id, user2.Id)
		return "idnameissameone"
	}

	if s.isFriend(user1.Id, user2.Id) == true {
		return "alreadyfriend"
	}

	newcontact := new(ContactInfo)
	newcontact.Id = len(s.LstContact)
	newcontact.LstUserId = append(newcontact.LstUserId, user1.Id)
	newcontact.LstUserId = append(newcontact.LstUserId, user2.Id)

	s.LstContact = append(s.LstContact, newcontact)

	// 两个玩家都需要给自己添加好友
	user1.FriendList = append(user1.FriendList, newcontact.Id)
	user2.FriendList = append(user2.FriendList, newcontact.Id)
	return "success"
}

/*
函数名：列出所有好友

参数：user id

返回值：err lstContact

err:错误信息
success: 取出成功
emptyuser: 传入id没有对应的 user

lstContact:所有的好友
*/
func (s *ModuleServer) GetFriendList(id int) (err string, lstContact []UserData) {
	if id < 0 || len(s.LstUser) <= id {
		return "emptyuser", lstContact
	}

	user := s.LstUser[id]
	if user == nil {
		return "emptyuser", lstContact
	}

	for _, index := range user.FriendList {
		for _, userId := range s.LstContact[index].LstUserId {
			if userId != user.Id {
				lstContact = append(lstContact, *s.LstUser[userId])
			}
		}
	}

	return "success", lstContact
}

/*
函数名：向一组好友发送 YO

参数：senderId lstId

senderId 发送者
geterId 接收者

返回值：err

err:错误信息
success: 取出成功
emptysender: 发送者为空
emptygeter: 接收者为空
sendergetersame: 发送者和接收者是同一个人
nofriend: 两者不是好友关系

*/
func (s *ModuleServer) SendYO(senderId int, geterId int) (err string) {
	if senderId < 0 || senderId >= len(s.LstUser) {
		return "emptysender"
	}

	if geterId < 0 || geterId >= len(s.LstUser) {
		return "emptygeter"
	}

	if senderId == geterId {
		return "sendergetersame"
	}

	if s.isFriend(senderId, geterId) == false {
		return "nofriend"
	}

	m := new(MsgInfo)
	m.Id = len(s.LstMsg)
	m.Type = EMsgInfo_Unread
	m.SenderId = senderId
	m.GeterId = geterId

	s.LstMsg = append(s.LstMsg, m)

	// geter 用户列表中存入此消息
	user := s.LstUser[geterId]
	user.MsgList = append(user.MsgList, m.Id)

	return "success"
}

/*
函数名：取得自己的 YO

参数：geterId

geterId 接收者

返回值：err

err:错误信息
success: 取出成功
emptyuser: 用户为空

*/
func (s *ModuleServer) GetYO(geterId int) (err string, lstYO []MsgInfo) {
	if geterId < 0 || geterId >= len(s.LstUser) {
		return "emptyuser", lstYO
	}

	user := s.LstUser[geterId]
	for _, msgId := range user.MsgList {
		msg := s.LstMsg[msgId]
		if msg.Type == EMsgInfo_Unread {
			lstYO = append(lstYO, *msg)
			msg.Type = EMsgInfo_Read
		}
	}

	return "success", lstYO
}
