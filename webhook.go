package viber

import (
	"encoding/json"
	"fmt"
)

//
//https://chatapi.viber.com/pa/set_webhook
// {
//    "url": "https://my.host.com",
//    "event_types": ["delivered", "seen", "failed", "subscribed", "unsubscribed", "conversation_started"]
// }

// WebhookReq request
type WebhookReq struct {
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
}

// {
//     "status": 0,
//     "status_message": "ok",
//     "event_types": ["delivered", "seen", "failed", "subscribed",  "unsubscribed", "conversation_started"]
// }

//WebhookResp response
type WebhookResp struct {
	Status        int      `json:"status"`
	StatusMessage string   `json:"status_message"`
	EventTypes    []string `json:"event_types"`
}

// WebhookVerify response
type WebhookVerify struct {
	Event        string `json:"event"`
	Timestamp    uint64 `json:"timestamp"`
	MessageToken uint64 `json:"message_token"`
}

// SetWebhook for Viber callbacks
func (v *Viber) SetWebhook(url string, eventTypes []string) (WebhookResp, error) {

	req := WebhookReq{
		URL:        url,
		EventTypes: eventTypes,
	}

	r, err := v.PostData("https://chatapi.viber.com/pa/set_webhook", req)

	fmt.Println(string(r))
	var resp WebhookResp
	json.Unmarshal(r, &resp)

	return resp, err
}
