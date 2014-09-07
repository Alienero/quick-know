##Restful API Document
- 均用POST 采用Basic authentic进行认证

<h4>1. 私信推送</h4>
- URL:`/push/private`
- Body：
```json
{
	"owner":"发起者ID",
	"to_id":"接收方ID",
	"expired":int64,
	"body":"消息体"       // 以Base64编码的字节数组
}
```
<h4>2.添加一个子用户</h4>
- URL:```/push/add_user```
- Body：
```json
{
	"psw":"用户密码"
}
```
<h4>3.删除一个子用户</h4>
- URL:```/push/del_user```
- Body：
```json
{
	"id":"用户ID"
}
```
<h4>4.添加一个组群</h4>
- URL:```/push/add_sub```
- Body：
```json
{
	"max":int  // 最大成员数
}
```
<h4>5.删除一个组群</h4>
- URL:```/push/del_sub```
- Body：
```json
{
	"id":"群组ID"
}
```
<h4>6.加群</h4>
- URL:```/push/user_sub```
- Body：
```json
{
	"sub_id":"群组ID",
	"user_id":"进群用户ID"
}
```
<h4>7.退群</h4>
- URL:```/push/rm_user_sub```
- Body：
```json
{
	"sub_id":"群组ID",
	"user_id":"退群用户ID"
}
```
<h4>8.面向群推送</h4>
- URL:```/push/group_msg```
- Body：
```json
{
	"sub_id":"",
	"msg":
		{
			"owner":"发起者ID",
			"expired":int64,     // 过期时间
			"body":"消息体"       // 以Base64编码的字节数组
		}
}
```
<h4>9.对所用用户广播</h4>
- URL:```/push/all```
- Body：
```json
{
	"owner":"发起者ID",
	"expired":int64,     // 过期时间
	"body":"消息体"       // 以Base64编码的字节数组
}
```
