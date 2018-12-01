package smseagle

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/messagebird/sachet"
)

//SmsEagleConfig configuration struct for mediaburst Client
type SmsEagleConfig struct {
	Pass    string `yaml:"pass"`
	Login   string `yaml:"login"`
	BaseUrl string `yaml:"baseurl"`
	Group   bool   `yaml:"group"`
}

//SmsEagleRequestTimeout  is the timeout for http request to mediaburst
const SmsEagleRequestTimeout = time.Second * 20

//SmsEagle is the exte SmsEagle
type SmsEagle struct {
	SmsEagleConfig
}

//NewSmsEagle creates a new
func NewSmsEagle(config SmsEagleConfig) *SmsEagle {
	SmsEagle := &SmsEagle{config}
	return SmsEagle
}

//Send send sms to n number of people using bulk sms api
func (c *SmsEagle) Send(message sachet.Message) (err error) {
	var request *http.Request
	var resp *http.Response
	var buffer bytes.Buffer
	var form url.Values

	buffer.WriteString(c.BaseUrl)
	if c.Group {
		buffer.WriteString("send_togroup")
		form = url.Values{"login": {c.Login}, "pass": {c.Pass}, "message": {message.Text}, "groupname": {message.To[0]}}
	} else {
		buffer.WriteString("send_tocontact")
	}

	form = url.Values{"login": {c.Login}, "pass": {c.Pass}, "message": {message.Text}, "contactname": {message.To[0]}}

	buffer.WriteString("?")
	buffer.WriteString(form.Encode())

	// preparing the request
	request, err = http.NewRequest("GET", buffer.String(), nil) // strings.NewReader(form.Encode()))
	if err != nil {
		return
	}

	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Set("User-Agent", "SachetV1.0")
	// calling the endpoint
	httpClient := &http.Client{}
	httpClient.Timeout = SmsEagleRequestTimeout

	resp, err = httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var body []byte
	resp.Body.Read(body)
	if resp.StatusCode == http.StatusOK && err == nil {
		return
	}
	return fmt.Errorf("Failed sending sms:Reason: %s , StatusCode : %d", string(body), resp.StatusCode)
}
