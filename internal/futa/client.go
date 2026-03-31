package futa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/tuandoquoc/futa-ticket-hunter/internal/config"
)

type Client struct {
	cfg        config.FutaConfig
	httpClient *http.Client
	token      string
	tokenTime  time.Time
}

func NewClient(cfg config.FutaConfig) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) getToken(ctx context.Context) (string, error) {
	if c.token != "" && time.Since(c.tokenTime) < 30*time.Minute {
		return c.token, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.cfg.WebURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch futabus.vn: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`"token"\s*:\s*"([^"]+)"`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("token not found in response")
	}

	c.token = string(matches[1])
	c.tokenTime = time.Now()
	return c.token, nil
}

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) ([]byte, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, err
	}

	fullURL := c.cfg.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.cfg.UserAgent)
	req.Header.Set("X-App-Version", c.cfg.AppVersion)
	req.Header.Set("X-Access-Token", token)
	req.Header.Set("X-Channel", c.cfg.Channel)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// API Response types

type APIResponse struct {
	RequestID string          `json:"requestId"`
	Status    int             `json:"status"`
	Error     json.RawMessage `json:"error"`
	Data      json.RawMessage `json:"data"`
}

type PaginationData struct {
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Total  int               `json:"total"`
	Items  []json.RawMessage `json:"items"`
	Others []json.RawMessage `json:"others"`
}

type ListData struct {
	Items  []json.RawMessage `json:"items"`
	Others []json.RawMessage `json:"others"`
}

type SingleData struct {
	Item json.RawMessage `json:"item"`
}

type PickupPointGroup struct {
	DistrictID   string        `json:"districtId"`
	DistrictName string        `json:"districtName"`
	ProvinceName string        `json:"provinceName"`
	Group        []PickupPoint `json:"group"`
}

type PickupPoint struct {
	DepartmentID      string  `json:"departmentId"`
	DepartmentName    string  `json:"departmentName"`
	DepartmentAddress string  `json:"departmentAddress"`
	DepartmentTime    int     `json:"departmentTime"`
	AreaID            string  `json:"areaId"`
	ProvinceID        string  `json:"provinceId"`
	ProvinceName      string  `json:"provinceName"`
	DistrictID        string  `json:"districtId"`
	DistrictName      string  `json:"districtName"`
	Type              int     `json:"type"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
}

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type RouteSearchData struct {
	RouteID string `json:"routeId"`
	From    string `json:"from"`
	To      string `json:"to"`
}

type TripSearchData struct {
	TripID             string `json:"tripId"`
	DepartureTime      string `json:"departureTime"`
	RawDepartureTime   string `json:"rawDepartureTime"`
	RawDepartureDate   string `json:"rawDepartureDate"`
	ArrivalTime        string `json:"arrivalTime"`
	Duration           int    `json:"duration"`
	SeatTypeName       string `json:"seatTypeName"`
	Price              int    `json:"price"`
	EmptySeatQuantity  int    `json:"emptySeatQuantity"`
	RouteID            string `json:"routeId"`
	Distance           int    `json:"distance"`
	WayID              string `json:"wayId"`
	MaxSeatsPerBooking int    `json:"maxSeatsPerBooking"`
	WayName            string `json:"wayName"`
	Route              Route  `json:"route"`
	SeatTypeCode       string `json:"seatTypeCode"`
}

type Route struct {
	OriginCode string `json:"originCode"`
	DestCode   string `json:"destCode"`
	OriginName string `json:"originName"`
	DestName   string `json:"destName"`
	Name       string `json:"name"`
	OriginHub  string `json:"originHubName"`
	DestHub    string `json:"destHubName"`
}

type SeatDiagramData struct {
	SeatID   string `json:"seatId"`
	Name     string `json:"name"`
	Status   []int  `json:"status"`
	ColumnNo int    `json:"columnNo"`
	RowNo    int    `json:"rowNo"`
	Floor    string `json:"floor"`
	Price    int    `json:"price"`
}

type DepartmentInWay struct {
	DepartmentID      string  `json:"departmentId"`
	DepartmentName    string  `json:"departmentName"`
	DepartmentAddress string  `json:"departmentAddress"`
	TimeAtDepartment  int     `json:"timeAtDepartment"`
	Passing           bool    `json:"passing"`
	IsShuttleService  bool    `json:"isShuttleService"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	PointKind         int     `json:"pointKind"`
	PresentBeforeMins int     `json:"presentBeforeMinutes"`
}

type BookingResponse struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	TotalPrice int    `json:"totalPrice"`
}

// API methods

func (c *Client) SearchPickupPoints(ctx context.Context, keyword string) ([]PickupPointGroup, []Area, error) {
	path := fmt.Sprintf("/vato/v1/search/pickup-point?keyword=%s&page=0&size=50",
		url.QueryEscape(keyword))

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, nil, err
	}
	if resp.Status != 200 {
		return nil, nil, fmt.Errorf("API error status %d: %s", resp.Status, string(resp.Error))
	}

	var pg PaginationData
	if err := json.Unmarshal(resp.Data, &pg); err != nil {
		return nil, nil, err
	}

	var groups []PickupPointGroup
	for _, item := range pg.Items {
		var g PickupPointGroup
		if err := json.Unmarshal(item, &g); err != nil {
			continue
		}
		groups = append(groups, g)
	}

	var areas []Area
	for _, other := range pg.Others {
		var a Area
		if err := json.Unmarshal(other, &a); err != nil {
			continue
		}
		areas = append(areas, a)
	}

	return groups, areas, nil
}

func (c *Client) SearchRoutes(ctx context.Context, originAreaID, destAreaID, fromDate string) ([]RouteSearchData, error) {
	path := fmt.Sprintf("/vato/v1/search/routes?destAreaId=%s&destOfficeId=&originAreaId=%s&originOfficeId=&isReturn=false&isReturnTripLoad=false&fromDate=%s",
		url.QueryEscape(destAreaID),
		url.QueryEscape(originAreaID),
		url.QueryEscape(fromDate))

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("routes API error: %s", string(resp.Error))
	}

	var ld ListData
	if err := json.Unmarshal(resp.Data, &ld); err != nil {
		return nil, err
	}

	var routes []RouteSearchData
	for _, item := range ld.Items {
		var r RouteSearchData
		if err := json.Unmarshal(item, &r); err != nil {
			continue
		}
		routes = append(routes, r)
	}
	return routes, nil
}

func (c *Client) SearchTripsByRoute(ctx context.Context, routeIDs []string, fromDate, toDate string) ([]TripSearchData, error) {
	type reqBody struct {
		MinNumSeat int      `json:"minNumSeat"`
		Channel    string   `json:"channel"`
		FromDate   string   `json:"fromDate"`
		ToDate     string   `json:"toDate"`
		RouteIDs   []string `json:"routeIds"`
		Sort       struct {
			ByPrice         string `json:"byPrice"`
			ByDepartureTime string `json:"byDepartureTime"`
		} `json:"sort"`
		Page int `json:"page"`
		Size int `json:"size"`
	}

	rb := reqBody{
		MinNumSeat: 1,
		Channel:    c.cfg.Channel,
		FromDate:   fromDate,
		ToDate:     toDate,
		RouteIDs:   routeIDs,
		Page:       0,
		Size:       200,
	}
	rb.Sort.ByPrice = "asc"
	rb.Sort.ByDepartureTime = "asc"

	bodyBytes, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest(ctx, "POST", "/vato/v1/search/trip-by-route", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("trips API error: %s", string(resp.Error))
	}

	var pg PaginationData
	if err := json.Unmarshal(resp.Data, &pg); err != nil {
		return nil, err
	}

	var trips []TripSearchData
	for _, item := range pg.Items {
		var t TripSearchData
		if err := json.Unmarshal(item, &t); err != nil {
			continue
		}
		trips = append(trips, t)
	}
	return trips, nil
}

func (c *Client) GetSeatDiagram(ctx context.Context, tripID string) ([]SeatDiagramData, error) {
	path := fmt.Sprintf("/vato/v1/search/seat-diagram/%s", tripID)

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("seat diagram API error: %s", string(resp.Error))
	}

	var pg PaginationData
	if err := json.Unmarshal(resp.Data, &pg); err != nil {
		return nil, err
	}

	var seats []SeatDiagramData
	for _, item := range pg.Items {
		var s SeatDiagramData
		if err := json.Unmarshal(item, &s); err != nil {
			continue
		}
		seats = append(seats, s)
	}
	return seats, nil
}

func (c *Client) GetDepartmentsInWay(ctx context.Context, wayID, routeID string) ([]DepartmentInWay, error) {
	path := fmt.Sprintf("/vato/v1/search/department-in-way/%s?routeId=%s", wayID, routeID)

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("dept-in-way API error: %s", string(resp.Error))
	}

	var pg PaginationData
	if err := json.Unmarshal(resp.Data, &pg); err != nil {
		return nil, err
	}

	var depts []DepartmentInWay
	for _, item := range pg.Items {
		var d DepartmentInWay
		if err := json.Unmarshal(item, &d); err != nil {
			continue
		}
		depts = append(depts, d)
	}
	return depts, nil
}

func (c *Client) BookReservation(ctx context.Context, passenger PassengerInfo, ticketInfo TicketInfo) (*BookingResponse, error) {
	type reqBody struct {
		Passenger  PassengerInfo `json:"passenger"`
		TicketInfo []TicketInfo  `json:"ticketInfo"`
		Channel    string        `json:"channel"`
	}

	rb := reqBody{
		Passenger:  passenger,
		TicketInfo: []TicketInfo{ticketInfo},
		Channel:    c.cfg.Channel,
	}

	bodyBytes, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest(ctx, "POST", "/vato/v1/booking/reservation", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("booking API error: %s", string(resp.Error))
	}

	var sd SingleData
	if err := json.Unmarshal(resp.Data, &sd); err != nil {
		return nil, err
	}

	var booking BookingResponse
	if err := json.Unmarshal(sd.Item, &booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

type PassengerInfo struct {
	CustName    string `json:"custName"`
	LoginMobile string `json:"loginMobile"`
	CustEmail   string `json:"custEmail"`
	CustSn      string `json:"custSn"`
	CustMobile  string `json:"custMobile"`
}

type TicketInfo struct {
	Seats   []SeatRef   `json:"seats"`
	Dropoff LocationRef `json:"dropoff"`
	TripID  string      `json:"tripId"`
	Pickup  LocationRef `json:"pickup"`
}

type SeatRef struct {
	SeatID string `json:"seatId"`
}

type LocationRef struct {
	Lng              float64 `json:"lng"`
	Address          string  `json:"address"`
	OfficeID         string  `json:"officeId"`
	Lat              float64 `json:"lat"`
	Name             string  `json:"name"`
	TimeAtDepartment int     `json:"timeAtDepartment"`
	Type             int     `json:"type"`
}

type PaymentInfoResponse struct {
	Code        string `json:"code"`
	Amount      string `json:"amount"`
	ExpiredTime string `json:"expiredTime"`
	PaymentURL  string `json:"paymentUrl"`
	QRCodeURL   string `json:"qrCodeUrl"`
	Message     string `json:"message"`
}

type PaymentStatusResponse struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	IsPaid bool   `json:"isPaid"`
}

func (c *Client) GetPaymentInfo(ctx context.Context, bookingID string, amount int) (*PaymentInfoResponse, error) {
	type reqBody struct {
		TicketID string `json:"ticketId"`
		Amount   int    `json:"amount"`
		Method   string `json:"method"`
		AppID    int    `json:"appId"`
	}

	rb := reqBody{
		TicketID: bookingID,
		Amount:   amount,
		Method:   "18",
		AppID:    360,
	}

	bodyBytes, err := json.Marshal(rb)
	if err != nil {
		return nil, err
	}

	data, err := c.doRequest(ctx, "POST", "/vato/v1/booking/get-link-payment", strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("payment info API error status %d: %s", resp.Status, string(resp.Error))
	}

	var info PaymentInfoResponse
	if err := json.Unmarshal(resp.Data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

type PaymentStatus struct {
	BookingCode string `json:"booking_code"`
	Status      string `json:"status"`
}

func (c *Client) PaymentStatus(ctx context.Context, bookingCode string) (bool, error) {
	path := fmt.Sprintf("/vato/v1/booking/payment-status/%s", bookingCode)
	data, err := c.doRequest(ctx, "GET", path, nil)

	if err != nil {
		return false, err
	}

	var resp APIResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return false, err
	}
	if resp.Status != 200 {
		return false, fmt.Errorf("payment info API error status %d: %s", resp.Status, string(resp.Error))
	}

	var info PaymentStatusResponse
	if err := json.Unmarshal(resp.Data, &info); err != nil {
		return false, err
	}

	return info.IsPaid, nil
}
