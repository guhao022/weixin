# 微信公众平台SDK for Go

这是一个使用Go语言编写的微信公众平台开发接口封装包.

创建`server.go`

```Go
package main

import (
	"log"
	"net/http"
	"github.com/sidbusy/weixin"
)

func main() {
	// 注册处理函数
	http.HandleFunc("/receiver", receiver)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func receiver(w http.ResponseWriter, r *http.Request) {
	token := "" // 微信公众平台的Token
	appid := "" // 微信公众平台的AppID
	secret := "" // 微信公众平台的AppSecret
	// 仅被动响应消息时可不填写appid、secret
	// 仅主动发送消息时可不填写token
	mp := weixin.New(token, appid, secret)
	// 检查请求是否有效
	// 仅主动发送消息时不用检查
	if !mp.Request.IsValid(w, r) {
		return
	}
	// 判断消息类型
	if mp.Request.MsgType == weixin.MsgTypeText {
		// 回复消息
		mp.ReplyTextMsg(w, "Hello, 世界")
	}
}
```

运行监听服务

`go run server.go`

客户端消息类型
-

`weixin.MsgTypeText` 文本消息

`weixin.MsgTypeImage` 图片消息

`weixin.MsgTypeVoice` 语音消息

`weixin.MsgTypeVideo` 视频消息

`weixin.MsgTypeShortVideo` 短视频消息

`weixin.MsgTypeLocation` 地理位置消息

`weixin.MsgTypeLink` 链接消息

`weixin.MsgTypeEvent` 事件推送

`weixin.EventSubscribe` 关注/用户未关注时扫描带参数二维码事件

`weixin.EventUnsubscribe` 取消关注事件

`weixin.EventScan` 用户已关注时扫描带参数二维码事件

`weixin.EventLocation` 上报地理位置事件

`weixin.EventClick` 菜单拉取消息事件

`weixin.EventView` 菜单跳转链接事件

回复消息
-

`mp.ReplyTextMsg(w, "content")` 回复文本消息

`mp.ReplyImageMsg(w, mediaId)` 回复图片消息

`mp.ReplyVoiceMsg(w, mediaId)` 回复语音消息

`mp.ReplyVideoMsg(w, &weixin.Video)` 回复视频消息

`mp.ReplyMusicMsg(w, &weixin.Music)` 回复音乐消息

`mp.ReplyNewsMsg(w, &[]weixin.Article)` 回复图文消息

`mediaId` 媒体文件上传后获取的唯一标识

返回`error`类型值

发送消息
-

`mp.SendTextMsg(touser, "content")` 发送文本消息

`mp.SendImageMsg(touser, mediaId)` 发送图片消息

`mp.SendVoiceMsg(touser, mediaId)` 发送语音消息

`mp.SendVideoMsg(touser, &weixin.Video)` 发送视频消息

`mp.SendMusicMsg(touser, &weixin.Music)` 发送音乐消息

`mp.SendNewsMsg(touser, &[]weixin.Article)` 发送图文消息

`touser` 普通用户openid

`mediaId` 媒体文件上传后获取的唯一标识

返回`error`类型值

视频、音乐、图文消息结构
-

视频消息

```Go
type Video struct {
	MediaId     string
	Title       string
	Description string
}
```

音乐消息

```Go
type Music struct {
	Title        string
	Description  string
	MusicUrl     string
	HQMusicUrl   string
	ThumbMediaId string
}
```

图文消息

```Go
type Article struct {
	Title       string
	Description string
	PicUrl      string
	Url         string
}
```

上传多媒体文件
-

`mp.UploadMediaFile(mediaType, filePath)` 上传多媒体文件

`mediaType` 多媒体文件类型(`weixin.MediaTypeImage`、`weixin.MediaTypeVoice`、`weixin.MediaTypeVideo`、`weixin.MediaTypeThumb`)

`filePath` 多媒体文件路径

返回`(string, error)`类型值， 返回的string类型值为媒体文件上传后获取的唯一标识.

下载多媒体文件
-

`mp.DownloadMediaFile(mediaId, filePath)` 下载多媒体文件

`mediaId` 媒体文件上传后获取的唯一标识

返回`error`类型值

创建二维码
-

`mp.CreateQRScene(sceneId)` 创建临时二维码

`mp.CreateQRLimitScene(expireSeconds, sceneId)` 创建永久二维码

`expireSeconds` 临时二维码有效时间, 以秒为单位, 最大不超过1800.

`sceneId` 场景值ID, 临时二维码时为32位非0整型, 永久二维码时最大值为100000 ( 目前参数只支持1-100000 ).

返回`(string, error)`类型值, 返回的string类型值为获取的二维码ticket, 凭借此ticket可以在有效时间内换取二维码.

换取二维码URL
-

`mp.GetQRCodeURL(ticket)`

返回`string`类型值

创建自定义菜单
-
`mp.CreateCustomMenu(&[]weixin.Button)` 创建自定义菜单

返回`error`类型值

查询自定义菜单
-
`mp.GetCustomMenu()` 查询自定义菜单

返回`([]weixin.Button, error)`类型值

删除自定义菜单
-
`mp.DeleteCustomMenu()` 删除自定义菜单

返回`error`类型值

## 相关链接

[微信公众平台](https://mp.weixin.qq.com/)

[微信公众平台开发者文档](http://mp.weixin.qq.com/wiki)

[微信公众平台接口调试工具](http://mp.weixin.qq.com/debug/)

[微信公众平台接口测试帐号申请](http://mp.weixin.qq.com/debug/cgi-bin/sandbox?t=sandbox/login)