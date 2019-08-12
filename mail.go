package o365Api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Mail interface {
	GetMailMesasges(string) ([]MailMessage, error)
}

type GetInboxMailRequest struct {
	bearerAccessToken string
}

type MailMessage struct {
	_odata_context  string `json:"@odata.context"`
	_odata_nextLink string `json:"@odata.nextLink"`
	Message         []struct {
		_odata_etag   string        `json:"@odata.etag"`
		BccRecipients []interface{} `json:"bccRecipients"`
		Body          struct {
			Content     string `json:"content"`
			ContentType string `json:"contentType"`
		} `json:"body"`
		BodyPreview     string        `json:"bodyPreview"`
		Categories      []interface{} `json:"categories"`
		CcRecipients    []interface{} `json:"ccRecipients"`
		ChangeKey       string        `json:"changeKey"`
		ConversationID  string        `json:"conversationId"`
		CreatedDateTime string        `json:"createdDateTime"`
		Flag            struct {
			FlagStatus string `json:"flagStatus"`
		} `json:"flag"`
		From struct {
			EmailAddress struct {
				Address string `json:"address"`
				Name    string `json:"name"`
			} `json:"emailAddress"`
		} `json:"from"`
		HasAttachments             bool          `json:"hasAttachments"`
		ID                         string        `json:"id"`
		Importance                 string        `json:"importance"`
		InferenceClassification    string        `json:"inferenceClassification"`
		InternetMessageID          string        `json:"internetMessageId"`
		IsDeliveryReceiptRequested interface{}   `json:"isDeliveryReceiptRequested"`
		IsDraft                    bool          `json:"isDraft"`
		IsRead                     bool          `json:"isRead"`
		IsReadReceiptRequested     bool          `json:"isReadReceiptRequested"`
		LastModifiedDateTime       string        `json:"lastModifiedDateTime"`
		ParentFolderID             string        `json:"parentFolderId"`
		ReceivedDateTime           string        `json:"receivedDateTime"`
		ReplyTo                    []interface{} `json:"replyTo"`
		Sender                     struct {
			EmailAddress struct {
				Address string `json:"address"`
				Name    string `json:"name"`
			} `json:"emailAddress"`
		} `json:"sender"`
		SentDateTime string `json:"sentDateTime"`
		Subject      string `json:"subject"`
		ToRecipients []struct {
			EmailAddress struct {
				Address string `json:"address"`
				Name    string `json:"name"`
			} `json:"emailAddress"`
		} `json:"toRecipients"`
		WebLink string `json:"webLink"`
	} `json:"value"`
}

func (request GetInboxMailRequest) GetInboxMail(bearerToken string) ([]MailMessage, error) {
	url := "https://graph.microsoft.com/v1.0/me/messages"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.bearerAccessToken))
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "graph.microsoft.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []MailMessage{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []MailMessage{}, err
	}

	var messages []MailMessage
	err = json.Unmarshal(body, &messages)

	return messages, nil
}

func (request GetInboxMailRequest) GetInboxMailBySenderAddress(bearerToken, fromAddress string) ([]MailMessage, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages?$filter=(from/emailAddress/address) eq '%s'")

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.bearerAccessToken))
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "graph.microsoft.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []MailMessage{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []MailMessage{}, err
	}

	var messages []MailMessage
	err = json.Unmarshal(body, &messages)

	return messages, nil
}
