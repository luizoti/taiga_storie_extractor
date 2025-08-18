package structs

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
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	CreatedDate  string `json:"created_date"`
	ModifiedDate string `json:"modified_date"`
	CsvUUID      string `json:"userstories_csv_uuid"`
}

type Storie struct {
	ID            int    `json:"id"`
	Ref           int    `json:"ref"`
	DueDate       string `json:"due_date"`
	DueDateReason string `json:"due_date_reason"`
	DueDateStatus string `json:"due_date_status"`
	Project       int    `json:"project"`
	CreatedDate   string `json:"created_date"`
	ModifiedDate  string `json:"modified_date"`
	FinishDate    string `json:"finish_date"`
	Name          string `json:"subject"`
	Comment       string `json:"comment"`
}

type CustomAttribute struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type StorieDetails struct {
	UuId         string `json:"id"`
	CreatedDate  string `json:"created_at"`
	Comment      string `json:"comment"`
	CommentHtml  string `json:"comment_html"`
	CustomFields map[string]string
}

type HeadersProvider func() AuthResponse
