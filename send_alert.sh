#!/bin/bash
alerts_message='
{
	"receiver": "web\\.hook\\.prometheusalert",
	"status": "firing",
	"alerts": [{
		"status": "firing",
		"labels": {
			"alertname": "服务异常告警",
			"app": "nginx",
			"hostname": "nginx-1",
			"instance": "172.18.68.209:58888",
			"level": "4",
			"name": "172.18.163.175:8880",
			"severity": "critical",
			"upstream": "order-server"
		},
		"annotations": {
			"description": "[故障] Nginx: nginx-1 服务: order-server 节点: 172.18.163.175:8880 down",
			"recovery_description": "[已恢复]: 服务: order-server "
		},
		"startsAt": "2023-12-29T10:44:10.492Z",
		"endsAt": "0001-01-01T00:00:00Z",
		"generatorURL": "http://zabbixserver:9090/graph?g0.expr=sum+by%28hostname%2C+instance%2C+name%2C+upstream%29+%28nginx_upstream_server_fall%7Bhostname%3D~%22nginx-a%7Ccloud-nginx-a%7Cnginx-1%7Cgdmaster%7Ccloud-nginx-1%22%2Cstatus%3D%22down%22%2Cupstream%21~%22.%2AB%7C.%2AA%22%7D%29\u0026g0.tab=1",
		"fingerprint": "e7e2e1691036c03c"
	}, {
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
		"generatorURL": "http://zabbixserver:9090/graph?g0.expr=sum+by%28hostname%2C+instance%2C+name%2C+upstream%29+%28nginx_upstream_server_fall%7Bhostname%3D~%22nginx-a%7Ccloud-nginx-a%7Cnginx-1%7Cgdmaster%7Ccloud-nginx-1%22%2Cstatus%3D%22down%22%2Cupstream%21~%22.%2AB%7C.%2AA%22%7D%29\u0026g0.tab=1",
		"fingerprint": "3407d122c7e8c961"
	}, {
		"status": "firing",
		"labels": {
			"alertname": "服务>异常告警",
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
		"generatorURL": "http://zabbixserver:9090/graph?g0.expr=sum+by%28hostname%2C+instance%2C+name%2C+upstream%29+%28nginx_upstream_server_fall%7Bhostname%3D~%22nginx-a%7Ccloud-nginx-a%7Cnginx-1%7Cgdmaster%7Ccloud-nginx-1%22%2Cstatus%3D%22down%22%2Cupstream%21~%22.%2AB%7C.%2AA%22%7D%29\u0026g0.tab=1",
		"fingerprint": "e069919362a06972"
	}, {
		"status": "firing",
		"labels": {
			"alertname": "服务异常告警",
			"app": "nginx",
			"hostname": "nginx-1",
			"instance": "172.18.68.209:58888",
			"level": "4",
			"name": "172.18.163.187:8880",
			"severity": "critical",
			"upstream": "order-server"
		},
		"annotations": {
			"description": "[故障] Nginx: nginx-1 服务: order-server 节>点: 172.18.163.187:8880 down",
			"recovery_description": "[已恢复]: 服务: order-server "
		},
		"startsAt": "2023-12-29T10:44:25.492Z",
		"endsAt": "0001-01-01T00:00:00Z",
		"generatorURL": "http://zabbixserver:9090/graph?g0.expr=sum+by%28hostname%2C+instance%2C+name%2C+upstream%29+%28nginx_upstream_server_fall%7Bhostname%3D~%22nginx-a%7Ccloud-nginx-a%7Cnginx-1%7Cgdmaster%7Ccloud-nginx-1%22%2Cstatus%3D%22down%22%2Cupstream%21~%22.%2AB%7C.%2AA%22%7D%29\u0026g0.tab=1",
		"fingerprint": "40941077c829b511"
	}, {
		"status": "firing",
		"labels": {
			"alertname": "服务异常告警",
			"app": "nginx",
			"hostname": "nginx-1",
			"instance": "172.18.68.209:58888",
			"level": "4",
			"name": "172.18.69.11:8085",
			"severity": "critical",
			"upstream": "cashier_api"
		},
		"annotations": {
			"description": "[故障] Nginx: nginx-1 服务: cashier_api 节点: 172.18.69.11:8085 down",
			"recovery_description": "[已恢复]: 服务: cashier_api "
		},
		"startsAt": "2023-12-29T10:45:10.492Z",
		"endsAt": "0001-01-01T00:00:00Z",
		"generatorURL": "http://zabbixserver:9090/graph?g0.expr=sum+by%28hostname%2C+instance%2C+name%2C+upstream%29+%28nginx_upstream_server_fall%7Bhostname%3D~%22nginx-a%7Ccloud-nginx-a%7Cnginx-1%7Cgdmaster%7Ccloud-nginx-1%22%2Cstatus%3D%22down%22%2Cupstream%21~%22.%2AB%7C.%2AA%22%7D%29\u0026g0.tab=1",
		"fingerprint": "4a2abc2f24613205"
	}],
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
'

curl -XPOST -d"$alerts_message" "http://localhost:8800/webhook?type=fs&tpl=feishu&fsurl=https://open.feishu.cn/open-apis/bot/v2/hook/bf8bb912-bc2e-40ad-9533-fcb8068aa621"
