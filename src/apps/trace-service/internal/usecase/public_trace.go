package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"RuntimeRoasters/apps/trace-service/internal/domain"
)

var (
	ErrSoldItemNotFound  = errors.New("sold item not found")
	ErrRetailUnavailable = errors.New("retail service unavailable")
)

type RetailSoldItem struct {
	ID                string    `json:"id"`
	SaleID            string    `json:"sale_id"`
	StoreID           string    `json:"store_id"`
	StoreName         string    `json:"store_name"`
	StoreCity         string    `json:"store_city"`
	StoreAddress      string    `json:"store_address"`
	MenuItemID        string    `json:"menu_item_id"`
	InventoryLotID    string    `json:"inventory_lot_id"`
	ProductID         string    `json:"product_id"`
	TraceCode         string    `json:"trace_code"`
	ProductName       string    `json:"product_name"`
	SKU               string    `json:"sku"`
	Size              string    `json:"size"`
	UnitPrice         float64   `json:"unit_price"`
	PublicURL         string    `json:"public_url"`
	SoldAt            time.Time `json:"sold_at"`
	InvoiceNo         string    `json:"invoice_no"`
	SourceBatchID     string    `json:"source_batch_id"`
	SourceHarvestID   string    `json:"source_harvest_id"`
	SourceWarehouseID string    `json:"source_warehouse_id"`
}

type RetailClient interface {
	GetSoldItem(context.Context, string) (*RetailSoldItem, error)
	ListSoldItems(context.Context, int) ([]RetailSoldItem, error)
}

type RetailHTTPClient struct {
	baseURL string
	secret  string
	http    *http.Client
}

func NewRetailHTTPClient(baseURL, secret string) *RetailHTTPClient {
	return &RetailHTTPClient{baseURL: strings.TrimRight(baseURL, "/"), secret: secret, http: &http.Client{Timeout: 3 * time.Second}}
}

func (c *RetailHTTPClient) GetSoldItem(ctx context.Context, productID string) (*RetailSoldItem, error) {
	var response struct {
		SoldItem RetailSoldItem `json:"sold_item"`
	}
	status, err := c.get(ctx, "/v1/internal/retail/sold-items/"+url.PathEscape(productID), &response)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, ErrSoldItemNotFound
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrRetailUnavailable, status)
	}
	return &response.SoldItem, nil
}

func (c *RetailHTTPClient) ListSoldItems(ctx context.Context, limit int) ([]RetailSoldItem, error) {
	var response struct {
		SoldItems []RetailSoldItem `json:"sold_items"`
	}
	status, err := c.get(ctx, fmt.Sprintf("/v1/internal/retail/sold-items?limit=%d", limit), &response)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("%w: status %d", ErrRetailUnavailable, status)
	}
	return response.SoldItems, nil
}

func (c *RetailHTTPClient) get(ctx context.Context, path string, out interface{}) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-Internal-Secret", c.secret)
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrRetailUnavailable, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return resp.StatusCode, nil
	}
	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode, nil
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(out); err != nil {
		return resp.StatusCode, fmt.Errorf("%w: invalid response: %v", ErrRetailUnavailable, err)
	}
	return resp.StatusCode, nil
}

type PublicTraceItem struct {
	ProductID   string    `json:"product_id"`
	TraceCode   string    `json:"trace_code"`
	ProductName string    `json:"product_name"`
	StoreName   string    `json:"store_name"`
	InvoiceNo   string    `json:"invoice_no"`
	SoldAt      time.Time `json:"sold_at"`
	PublicURL   string    `json:"public_url"`
}

type PublicMilestone struct {
	Status      string    `json:"status"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Service     string    `json:"service,omitempty"`
}

type PublicTraceDocument struct {
	ProductID string `json:"product_id"`
	TraceCode string `json:"trace_code"`
	SaleItem  struct {
		ID     string    `json:"id"`
		SoldAt time.Time `json:"sold_at"`
	} `json:"sale_item"`
	InvoiceNo            string                 `json:"invoice_no"`
	Product              map[string]interface{} `json:"product"`
	Store                map[string]interface{} `json:"store"`
	Farm                 map[string]interface{} `json:"farm"`
	Warehouse            map[string]interface{} `json:"warehouse"`
	OriginBatchID        string                 `json:"origin_batch_id"`
	HarvestID            string                 `json:"harvest_id"`
	InventoryLotID       string                 `json:"inventory_lot_id"`
	Milestones           []PublicMilestone      `json:"milestones"`
	ServiceParticipation []string               `json:"service_participation"`
	TraceIDs             []string               `json:"trace_ids"`
	PublicSafe           bool                   `json:"public_safe"`
}

func (s *Service) SetRetailClient(client RetailClient) { s.retail = client }

func (s *Service) ListPublicTraceItems(ctx context.Context, limit int) ([]PublicTraceItem, error) {
	if s.retail == nil {
		return nil, ErrRetailUnavailable
	}
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	items, err := s.retail.ListSoldItems(ctx, limit)
	if err != nil {
		return nil, err
	}
	result := make([]PublicTraceItem, 0, len(items))
	for _, item := range items {
		result = append(result, PublicTraceItem{ProductID: item.ProductID, TraceCode: item.TraceCode, ProductName: item.ProductName, StoreName: item.StoreName, InvoiceNo: item.InvoiceNo, SoldAt: item.SoldAt, PublicURL: item.PublicURL})
	}
	return result, nil
}

func (s *Service) GetPublicTrace(ctx context.Context, traceCode string) (*PublicTraceDocument, error) {
	traceCode = strings.TrimSpace(traceCode)
	if traceCode == "" {
		return nil, ErrSoldItemNotFound
	}
	if s.retail == nil {
		return nil, ErrRetailUnavailable
	}
	item, err := s.retail.GetSoldItem(ctx, traceCode)
	if err != nil {
		return nil, err
	}

	events, err := s.findLineageEvents(ctx, item.SourceBatchID, item.SourceHarvestID)
	if err != nil {
		return nil, err
	}
	doc := &PublicTraceDocument{ProductID: item.ProductID, TraceCode: item.TraceCode, InvoiceNo: item.InvoiceNo, OriginBatchID: item.SourceBatchID, HarvestID: item.SourceHarvestID, InventoryLotID: item.InventoryLotID, PublicSafe: true}
	doc.SaleItem.ID, doc.SaleItem.SoldAt = item.ID, item.SoldAt
	doc.Product = map[string]interface{}{"menu_item_id": item.MenuItemID, "name": item.ProductName, "sku": item.SKU, "size": item.Size}
	doc.Store = map[string]interface{}{"id": item.StoreID, "name": item.StoreName, "city": item.StoreCity, "address": item.StoreAddress}
	doc.Farm = map[string]interface{}{"harvest_id": item.SourceHarvestID}
	doc.Warehouse = map[string]interface{}{"id": item.SourceWarehouseID}
	doc.Milestones, doc.ServiceParticipation, doc.TraceIDs = publicMilestones(events)
	if len(doc.Milestones) == 0 {
		doc.Milestones = append(doc.Milestones, PublicMilestone{Status: "ORIGIN_LINKED", Title: "Origin lineage linked", Description: fmt.Sprintf("Harvest %s, batch %s", item.SourceHarvestID, item.SourceBatchID), Timestamp: item.SoldAt})
	}
	doc.Milestones = append(doc.Milestones, PublicMilestone{Status: "SERVED", Title: "Coffee served", Description: fmt.Sprintf("%s at %s", item.ProductName, item.StoreName), Timestamp: item.SoldAt, Service: "retail-service"})
	if !contains(doc.ServiceParticipation, "retail-service") {
		doc.ServiceParticipation = append(doc.ServiceParticipation, "retail-service")
	}
	return doc, nil
}

func (s *Service) findLineageEvents(ctx context.Context, batchID, harvestID string) ([]domain.TraceEvent, error) {
	var events []domain.TraceEvent
	query := s.db.WithContext(ctx).Where("visibility = ?", "public")
	if batchID != "" && harvestID != "" {
		query = query.Where("batch_id = ? OR harvest_id = ?", batchID, harvestID)
	} else if batchID != "" {
		query = query.Where("batch_id = ?", batchID)
	} else {
		query = query.Where("harvest_id = ?", harvestID)
	}
	return events, query.Order("occurred_at ASC, created_at ASC").Find(&events).Error
}

func publicMilestones(events []domain.TraceEvent) ([]PublicMilestone, []string, []string) {
	milestones := make([]PublicMilestone, 0, len(events))
	services := []string{}
	traceIDs := []string{}
	for _, event := range events {
		payload := parsePayload(event.Payload)
		entry := timelineEntry(event, payload)
		service := strings.TrimPrefix(event.SourceService, "/services/")
		milestones = append(milestones, PublicMilestone{Status: entry.Status, Title: entry.Title, Description: entry.Description, Timestamp: entry.Timestamp, Service: service})
		if service != "" && !contains(services, service) {
			services = append(services, service)
		}
		if event.TraceID != "" && !contains(traceIDs, event.TraceID) {
			traceIDs = append(traceIDs, event.TraceID)
		}
	}
	sort.Strings(services)
	return milestones, services, traceIDs
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
