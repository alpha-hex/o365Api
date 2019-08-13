package o365Api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Mail interface {
	GetMailMesasges(string) (MailMessage, error)
	GetInboxMailFromAddress(string) (MailMessage, error)
	GetMessageAttachement() (MessageAttachment, error)
}

type MailRequest struct {
	BearerAccessToken string
	MessageID         string
}

type MailMessage struct {
	OdataContext  string `json:"@odata.context"`
	OdataNextLink string `json:"@odata.nextLink"`
	Value         []struct {
		OdataEtag                  string        `json:"@odata.etag"`
		ID                         string        `json:"id"`
		CreatedDateTime            time.Time     `json:"createdDateTime"`
		LastModifiedDateTime       time.Time     `json:"lastModifiedDateTime"`
		ChangeKey                  string        `json:"changeKey"`
		Categories                 []interface{} `json:"categories"`
		ReceivedDateTime           time.Time     `json:"receivedDateTime"`
		SentDateTime               time.Time     `json:"sentDateTime"`
		HasAttachments             bool          `json:"hasAttachments"`
		InternetMessageID          string        `json:"internetMessageId"`
		Subject                    string        `json:"subject"`
		BodyPreview                string        `json:"bodyPreview"`
		Importance                 string        `json:"importance"`
		ParentFolderID             string        `json:"parentFolderId"`
		ConversationID             string        `json:"conversationId"`
		IsDeliveryReceiptRequested interface{}   `json:"isDeliveryReceiptRequested"`
		IsReadReceiptRequested     bool          `json:"isReadReceiptRequested"`
		IsRead                     bool          `json:"isRead"`
		IsDraft                    bool          `json:"isDraft"`
		WebLink                    string        `json:"webLink"`
		InferenceClassification    string        `json:"inferenceClassification"`
		Body                       struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		} `json:"body"`
		Sender struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"sender"`
		From struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"from"`
		ToRecipients []struct {
			EmailAddress struct {
				Name    string `json:"name"`
				Address string `json:"address"`
			} `json:"emailAddress"`
		} `json:"toRecipients"`
		CcRecipients  []interface{} `json:"ccRecipients"`
		BccRecipients []interface{} `json:"bccRecipients"`
		ReplyTo       []interface{} `json:"replyTo"`
		Flag          struct {
			FlagStatus string `json:"flagStatus"`
		} `json:"flag"`
	} `json:"value"`
}

type MessageAttachment struct {
	_odata_context string `json:"@odata.context"`
	Value          []struct {
		_odata_type          string      `json:"@odata.type"`
		ContentBytes         string      `json:"contentBytes"`
		ContentID            string      `json:"contentId"`
		ContentLocation      interface{} `json:"contentLocation"`
		ContentType          string      `json:"contentType"`
		ID                   string      `json:"id"`
		IsInline             bool        `json:"isInline"`
		LastModifiedDateTime string      `json:"lastModifiedDateTime"`
		Name                 string      `json:"name"`
		Size                 int         `json:"size"`
	} `json:"value"`
}

func (request MailRequest) GetInboxMail(bearerToken string) (MailMessage, error) {
	url := "https://graph.microsoft.com/v1.0/me/messages"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MailMessage{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return MailMessage{}, err
	}

	var messages MailMessage
	err = json.Unmarshal(body, &messages)

	return messages, nil
}

func (request MailRequest) GetInboxMailFromAddress(fromAddress string) (MailMessage, error) {
	queryParams := fmt.Sprintf("(from/emailAddress/address) eq '%s'", fromAddress)
	queryParams = url.QueryEscape(queryParams)
	queryURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages?$filter=%s", queryParams)

	req, _ := http.NewRequest("GET", queryURL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MailMessage{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return MailMessage{}, err
	}

	var messages MailMessage
	if err := json.Unmarshal(body, &messages); err != nil {
		fmt.Println(err)
		return MailMessage{}, err
	}

	return messages, nil
}

func (request MailRequest) GetMessageAttachement() (MessageAttachment, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s/attachments", request.MessageID)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "graph.microsoft.com")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MessageAttachment{}, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var msgAttachment MessageAttachment
	if err := json.Unmarshal(body, &msgAttachment); err != nil {
		return MessageAttachment{}, err
	}

	return msgAttachment, nil
}
