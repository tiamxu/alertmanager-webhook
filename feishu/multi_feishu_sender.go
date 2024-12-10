package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tiamxu/alertmanager-webhook/model"
)

type MultiFeiShuSender struct {
	Robots []RobotConfig
}

type RobotConfig struct {
	Name       string `json:"name"`
	WebhookURL string `json:"webhook_url"`
}

func (m *MultiFeiShuSender) SendToRobot(name string, message *model.CommonMessage) error {
	for _, robot := range m.Robots {
		if robot.Name == name {
			payload := map[string]interface{}{
				"msg_type": "text",
				"content": map[string]string{
					"text": message.Text,
				},
			}

			jsonPayload, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			resp, err := http.Post(robot.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			fmt.Printf("FeiShu response from %s: %s\n", name, body)
			return nil
		}
	}
	return fmt.Errorf("robot not found: %s", name)
}

func (m *MultiFeiShuSender) Send(message *model.CommonMessage) error {
	return m.SendToRobot(message.Platform, message)
}
