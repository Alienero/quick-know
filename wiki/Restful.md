##Restful API Document
- 均用POST 采用Basic authentic进行认证
###1. 私信推送
- URL:`/push/private`
- Body：
<pre><code>{
	"owner":"发起者ID",
	"to_id":"接收方ID",
	"expired":int64,
	"body":"消息体"       // 以Base64编码的字节数组
}</code></pre>
###2.添加一个子用户
- URL:```/push/add_user```
- Body：
<pre><code>{
	"psw":"用户密码"
}</code></pre>
###3.删除一个子用户
- URL:```/push/del_user```
- Body：
<pre><code>{
	"id":"用户ID"
}</pre></code>
###4.添加一个组群
- URL:```/push/add_sub```
- Body：
<pre><code>{
	"max":int  // 最大成员数
}</pre></code>
###5.删除一个组群
- URL:```/push/del_sub```
- Body：
<pre><code>{
	"id":"群组ID"
}</pre></code>
###6.加群
- URL:```/push/user_sub```
- Body：
<pre><code>{
	"sub_id":"群组ID",
	"user_id":"进群用户ID"
}
</code></pre>
###7.退群
- URL:```/push/rm_user_sub```
- Body：
<pre><code>{
	"sub_id":"群组ID",
	"user_id":"退群用户ID"
}
</code></pre>
###8.面向群推送
- URL:```/push/group_msg```
- Body：
<pre><code>{
	"sub_id":"",
	"msg":
		{
			"owner":"发起者ID",
			"expired":int64,     // 过期时间
			"body":"消息体"       // 以Base64编码的字节数组
		}
}
</code></pre>
###9.对所用用户广播
- URL:```/push/all```
- Body：
<pre><code>{
	"owner":"发起者ID",
	"expired":int64,     // 过期时间
	"body":"消息体"       // 以Base64编码的字节数组
}</code></pre>
