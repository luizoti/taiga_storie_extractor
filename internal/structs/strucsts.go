package structs

import (
	"strings"
	"time"
)

type UserCredentials struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"auth_token"`
}

type AuthenticatedHeader struct {
	Authorization      string `json:"Authorization"`
	Accept             string `json:"Accept"`
	AcceptLanguage     string `json:"Accept-Language"`
	Connection         string `json:"Connection"`
	SecFetchDest       string `json:"Sec-Fetch-Dest"`
	SecFetchMode       string `json:"Sec-Fetch-Mode"`
	SecFetchSite       string `json:"Sec-Fetch-Site"`
	UserAgent          string `json:"User-Agent"`
	SecChUa            string `json:"sec-ch-ua"`
	SecChUaMobile      string `json:"sec-ch-ua-mobile"`
	SecChUaPlatform    string `json:"sec-ch-ua-platform"`
	XdisablePagination string `json:"x-disable-pagination"`
}

type Project struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	Slug         string       `json:"slug"`
	Description  string       `json:"description"`
	CreatedDate  FormatedDate `json:"created_date"`
	ModifiedDate FormatedDate `json:"modified_date"`
	CsvUUID      string       `json:"userstories_csv_uuid"`
}

type Storie struct {
	ID            int          `json:"id"`
	Ref           int          `json:"ref"`
	DueDate       FormatedDate `json:"due_date"`
	DueDateReason string       `json:"due_date_reason"`
	DueDateStatus string       `json:"due_date_status"`
	Project       int          `json:"project"`
	CreatedDate   FormatedDate `json:"created_date"`
	ModifiedDate  FormatedDate `json:"modified_date"`
	FinishDate    FormatedDate `json:"finish_date"`
	Name          string       `json:"subject"`
	Comment       string       `json:"comment"`
}

type CustomAttribute struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type StorieDetails struct {
	UuId         string       `json:"id"`
	CreatedDate  FormatedDate `json:"created_at"`
	Comment      string       `json:"comment"`
	CommentHtml  string       `json:"comment_html"`
	CustomFields map[string]string
}

type HeadersProvider func() AuthResponse

type FormatedDate string

func (t *FormatedDate) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		*t = ""
	} else {

		clearedTime := strings.Split(strings.Trim(string(b), "\""), ".")[0]

		if date, err := time.Parse("2006-01-02T15:04:05", clearedTime); err == nil {
			*t = FormatedDate(date.Format("2006/01/02 15:04:05"))
			return nil
		}

		if date, err := time.Parse("2006-01-02", clearedTime); err == nil {
			*t = FormatedDate(date.Format("2006/01/02"))
			return nil
		}
		*t = FormatedDate(clearedTime)
	}
	return err
}
