##
```
go build -o main
./main
```
## webhook接口示例
```
http://localhost:8801/webhook?type=fs&tpl=feishu&fsurl=https://open.feishu.cn/open-apis/bot/v2/hook/bf8bb912-bc2e-40ad-9533-fcb8068aa621&at=ou_1199d79525e146bad9d0a5a46a86a10f

http://localhost:8801/webhook?type=dd&tpl=dingtalk&ddurl=https://oapi.dingtalk.com/robot/send?access_token=9ef3af0bc7052966a73c6642eed0e7c90e35a4dd6860887dd9029c65255d5abd&split=true&at=1888888888
```
## 参数说明
```
type: (Required)类型 飞书:fs ,钉钉:dd
tpl: (Optional)模版名 ./template目录下,名字为模版名。 非必须，可不提供
split: (Optional) 是否对分组告警进行拆分为单条 true:拆分,默认; false：不拆分
fsurl/ddurl: (Required)告警webhook地址,飞书是fsurl, 钉钉是ddurl
at: (Optional )支持at人，自定义机器人支持使用 open_id、user_id;钉钉使用手机号, 多个用逗号分隔;
    另外支持规则rule@人labels.annotations.at: "id1,id2"
```
## 支持功能
- 支持飞书、钉钉告警
- 支持定义卡片title，为告警规则rule名称
- 支持自定义消息模版，不配置为默认模版
## 告警测试
使用postman测试

## 飞书告警模版示例
```
{{ $var := .ExternalURL}}{{ range $k, $v := .Alerts }}
{{if eq $v.Status "resolved"}}
**<font color="green">开始时间</font>:** {{GetCSTtime $v.StartsAt}}
**<font color="green">结束时间</font>:** {{GetCSTtime $v.EndsAt}}
**<font color="green">故障主机</font>:** {{$v.Labels.instance}}
**<font color="green">告警详情</font>:** {{$v.Annotations.summary}}
**<font color="red">告警描述</font>:** {{$v.Annotations.recovery_description}}
{{ else }}
**<font color="red">开始时间</font>:** {{GetCSTtime $v.StartsAt}}
**<font color="red">故障主机</font>:** {{$v.Labels.instance}}
**<font color="red">告警详情</font>:** {{$v.Annotations.summary}}
**<font color="green">告警描述</font>:** {{$v.Annotations.description}}
{{- end }}
{{- end -}}
```
## 钉钉告警模版
```
{{ $var := .ExternalURL}}{{ range $k,$v:=.Alerts }}
{{if eq $v.Status "resolved"}}

##### <font color="#02b340">触发时间</font>: {{GetCSTtime $v.StartsAt}}
##### <font color="#02b340">结束时间</font>: {{GetCSTtime $v.EndsAt}}
##### <font color="#02b340">描述信息</font>: {{$v.Annotations.recovery_description}}  

---

{{ else }}

##### <font color="#FF0000">触发时间</font>: {{GetCSTtime $v.StartsAt}}
##### <font color="#FF0000">描述信息</font>: {{$v.Annotations.recovery_description}}  

---

{{end}}
{{- end }}
```
## 测试数据
```
{
	"receiver": "web\\.hook\\.prometheusalert",
	"status": "firing",
	"alerts": [
		 {
		"status": "firing",
		"labels": {
			"alertname": "服务异常告警",
			"app": "nginx",
			"hostname": "nginx-1",
			"instance": "172.18.68.209:58888",
			"level": "4",
			"name": "172.18.163.177:8085",
			"severity": "critical",
			"upstream": "cashier_api"
		},
		"annotations": {
			"at": "1888888888",
			"description": "[故障] Nginx: nginx-1 服务: cashier_api 节点: 172.18.163.177:8085 down",
			"recovery_description": "[已恢复]: 服务: cashier_api "
		},
		"startsAt": "2023-12-29T10:44:55.492Z",
		"endsAt": "0001-01-01T00:00:00Z",
		"generatorURL": "http://zabbixserver:9090/graph?1",
		"fingerprint": "3407d122c7e8c961"
	}, {
		"status": "firing",
		"labels": {
			"alertname": "服务异常告警",
			"app": "nginx",
			"hostname": "nginx-1",
			"instance": "172.18.68.209:58888",
			"level": "4",
			"name": "172.18.163.177:8880",
			"severity": "critical",
			"upstream": "order-server"
		},
		"annotations": {
			"at": "1888888888",
			"description": "[故障] Nginx: nginx-1 服务: order-server 节点: 172.18.163.177:8880 down",
			"recovery_description": "[已恢复]: 服务: order-server "
		},
		"startsAt": "2023-12-29T10:44:40.492Z",
		"endsAt": "0001-01-01T00:00:00Z",
		"generatorURL": "http://zabbixserver:9090/graph?",
		"fingerprint": "e069919362a06972"
	}

	],
	"groupLabels": {
		"alertname": "服务异常告警",
		"instance": "172.18.68.209:58888"
	},
	"commonLabels": {
		"alertname": "服务异常告警",
		"app": "nginx",
		"hostname": "nginx-1",
		"instance": "172.18.68.209:58888",
		"level": "4",
		"severity": "critical"
	},
	"commonAnnotations": {},
	"externalURL": "http://aaec49e977ef:9093",
	"version": "4",
	"groupKey": "{}/{app=\"nginx\"}:{alertname=\"服务异常告警\", instance=\"172.18.68.209:58888\"}",
	"truncatedAlerts": 0
}

```
