package o365Api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Mail interface {
	GetMailMesasges(string) (MailMessage, error)
	GetInboxMailFromAddress(string) (MailMessage, error)
	GetMessageById(string) (MailMessageContent, error)
	GetMessageAttachementByID(string, string, string, string) (MessageAttachementContent, error)
	GetMessageAttachement() (MessageAttachment, error)
	GetTopLevelMailFolders() (MailBoxFolder, error)
	GetChildLevelMailFolders(string) (MailBoxFolder, error)
	GetMailFolderMessages(string) (MailMessage, error)
	MoveMessage(messageId, destinationFolderId string) (bool, error)
}

type MailRequest struct {
	BearerAccessToken string
}

type MailMessage struct {
	OdataContext  string               `json:"@odata.context"`
	OdataNextLink string               `json:"@odata.nextLink"`
	Value         []MailMessageContent `json:"value"`
}

type MailMessageContent struct {
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
}

type MessageAttachment struct {
	OdataContext string                      `json:"@odata.context"`
	Value        []MessageAttachementContent `json:"value"`
}

type MessageAttachementContent struct {
	OdataType            string      `json:"@odata.type"`
	ID                   string      `json:"id"`
	LastModifiedDateTime time.Time   `json:"lastModifiedDateTime"`
	Name                 string      `json:"name"`
	ContentType          string      `json:"contentType"`
	Size                 int         `json:"size"`
	IsInline             bool        `json:"isInline"`
	ContentID            string      `json:"contentId"`
	ContentLocation      interface{} `json:"contentLocation"`
	ContentBytes         string      `json:"contentBytes"`
}

type MailBoxFolder struct {
	OdataContext  string `json:"@odata.context"`
	OdataNextLink string `json:"@odata.nextLink"`
	Value         []struct {
		ID               string `json:"id"`
		DisplayName      string `json:"displayName"`
		ParentFolderID   string `json:"parentFolderId"`
		ChildFolderCount int    `json:"childFolderCount"`
		UnreadItemCount  int    `json:"unreadItemCount"`
		TotalItemCount   int    `json:"totalItemCount"`
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

func (request MailRequest) GetMessageById(messageId string) (MailMessageContent, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s", messageId)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MailMessageContent{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return MailMessageContent{}, err
	}

	var messages MailMessageContent
	if err := json.Unmarshal(body, &messages); err != nil {
		fmt.Println(err)
		return MailMessageContent{}, err
	}

	return messages, nil
}

func (request MailRequest) GetMessageAttachement(messageId string) (MessageAttachment, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s/attachments", messageId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var msgAttachment MessageAttachment
	if err := json.Unmarshal(body, &msgAttachment); err != nil {
		return MessageAttachment{}, err
	}

	return msgAttachment, nil
}

func (request MailRequest) GetMessageAttachementByID(parentFolderId, childFolderId, messageId, attachmentId string) (MessageAttachementContent, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/mailFolders/%s/childFolders/%s/messages/%s/attachments/%s", parentFolderId, childFolderId, messageId, attachmentId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var msgAttachment MessageAttachementContent
	if err := json.Unmarshal(body, &msgAttachment); err != nil {
		return MessageAttachementContent{}, err
	}

	return msgAttachment, nil
}

func (request MailRequest) GetTopLevelMailFolders() (MailBoxFolder, error) {
	url := "https://graph.microsoft.com/v1.0/me/mailFolders/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MailBoxFolder{}, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var folders MailBoxFolder
	if err := json.Unmarshal(body, &folders); err != nil {
		return MailBoxFolder{}, err
	}

	return folders, nil
}

func (request MailRequest) GetChildLevelMailFolders(parentFolderId string) (MailBoxFolder, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/mailFolders/%s/childFolders", parentFolderId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MailBoxFolder{}, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var folders MailBoxFolder
	if err := json.Unmarshal(body, &folders); err != nil {
		return MailBoxFolder{}, err
	}

	return folders, nil
}

func (request MailRequest) GetMailFolderMessages(childFolderId string) (MailMessage, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/mailFolders/%s/messages", childFolderId)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return MailMessage{}, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var messages MailMessage
	if err := json.Unmarshal(body, &messages); err != nil {
		fmt.Println(err)
		return MailMessage{}, err
	}

	return messages, nil
}

func (request MailRequest) MoveMessage(messageId, destinationFolderId string) (bool, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s/move", messageId)

	payload := strings.NewReader(fmt.Sprintf("{\"destinationId\": \"%s\"}", destinationFolderId))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return false, errors.New(res.Status)
	}

	return true, nil
}
