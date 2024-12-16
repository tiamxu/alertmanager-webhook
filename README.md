##
```
go build -o main
./main
```
## 测试接口
```
http://localhost:8801/webhook?type=fs&tpl=feishu&fsurl=https://open.feishu.cn/open-apis/bot/v2/hook/bf8bb912-bc2e-40ad-9533-fcb8068aa621&at=ou_1199d79525e146bad9d0a5a46a86a10f
```
## 参数说明
```
type: 类型 飞书:fs ,钉钉:dd
tpl: 模版名，./template目录下
split: 是否对分组告警进行拆分为单条 true:拆分,默认; false：不拆分
fsurl: 告警webhook地址
at: 支持at人，自定义机器人仅支持使用 open_id、user_id,多个用逗号分隔
```
## 告警测试
使用postman测试

## 飞书告警模版
```
{{ $var := .ExternalURL}}{{ range $k, $v := .Alerts }}{{if eq $v.Status "resolved"}}
**【开始时间】**:{{GetCSTtime $v.StartsAt}}
**【结束时间】:** {{GetCSTtime $v.EndsAt}}
**【故障主机】:** {{$v.Labels.instance}}
**【告警描述】:** {{$v.Annotations.recovery_description}}
{{ else }}
**【开始时间】**:{{GetCSTtime $v.StartsAt}}
**【故障主机】**: {{$v.Labels.instance}}
**【告警描述】**: {{$v.Annotations.description}}
{{- end }}
{{ end -}}
<at id=ou_1199d79525e146bad9d0a5a46a86a10f></at>
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