package seed

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"RuntimeRoasters/apps/retail-service/internal/domain"
	"github.com/google/uuid"
)

const (
	DefaultMenuID = "MENU-COFFEE-DEFAULT"
	UnitGram      = "GRAM"
)

//go:embed stores.json
var storesSeed []byte

type DemoData struct {
	Lots      []domain.InventoryLot
	Sales     []domain.Sale
	SaleItems []domain.SaleItem
	Movements []domain.StockMovement
}

type menuSpec struct {
	group       string
	category    string
	name        string
	description string
	prices      map[string]float64
}

var menuSpecs = []menuSpec{
	{"PHIN-SUA-DA", "PHIN", "Phin Sữa Đá", "Cà phê phin truyền thống kết hợp sữa đặc và đá. Vị đậm đà, béo ngọt.", map[string]float64{"S": 29000, "M": 39000, "L": 45000}},
	{"PHIN-SUA-NONG", "PHIN", "Phin Sữa Nóng", "Cà phê phin truyền thống kết hợp sữa đặc uống nóng.", map[string]float64{"S": 29000, "M": 39000}},
	{"PHIN-DEN-DA", "PHIN", "Phin Đen Đá", "Cà phê phin đen nguyên chất kèm đá. Đậm đặc, đắng thanh.", map[string]float64{"S": 29000, "M": 35000, "L": 39000}},
	{"PHIN-DEN-NONG", "PHIN", "Phin Đen Nóng", "Cà phê phin đen nguyên chất uống nóng.", map[string]float64{"S": 29000, "M": 35000}},
	{"BAC-XIU", "PHIN", "Bạc Xỉu", "Thức uống thiên sữa, pha thêm chút cà phê Phin thơm nhẹ.", map[string]float64{"S": 29000, "M": 39000, "L": 45000}},
	{"PHINDI-KEM-SUA", "PHINDI", "PhinDi Kem Sữa", "Cà phê Phin êm nhẹ kết hợp sữa đặc và kem sữa.", map[string]float64{"S": 45000, "M": 55000, "L": 60000}},
	{"PHINDI-HANH-NHAN", "PHINDI", "PhinDi Hạnh Nhân", "Cà phê Phin êm nhẹ kết hợp sữa và hương vị hạnh nhân.", map[string]float64{"S": 45000, "M": 55000, "L": 60000}},
	{"PHINDI-CHOCO", "PHINDI", "PhinDi Choco", "Cà phê Phin êm nhẹ, sữa và sốt sô-cô-la.", map[string]float64{"S": 45000, "M": 55000, "L": 60000}},
	{"ESPRESSO-DOUBLE", "ESPRESSO", "Espresso / Double", "Shot espresso nguyên chất chiết xuất từ máy pha.", map[string]float64{"S": 35000, "M": 45000}},
	{"AMERICANO", "ESPRESSO", "Americano", "Shot espresso pha loãng với nước nóng. Vị đắng nhẹ, thanh thoát.", map[string]float64{"S": 39000, "M": 49000, "L": 55000}},
	{"CAPPUCCINO", "ESPRESSO", "Cappuccino", "Espresso kết hợp sữa nóng và lớp bọt sữa dày mịn.", map[string]float64{"S": 59000, "M": 69000, "L": 79000}},
	{"LATTE", "ESPRESSO", "Latte", "Espresso kết hợp nhiều sữa nóng và lớp bọt sữa mỏng.", map[string]float64{"S": 59000, "M": 69000, "L": 79000}},
	{"CARAMEL-MACCHIATO", "ESPRESSO", "Caramel Macchiato", "Espresso kết hợp sữa nóng, sốt caramel và hương vani.", map[string]float64{"S": 65000, "M": 75000, "L": 85000}},
	{"PHIN-FREEZE", "FREEZE", "Phin Freeze", "Cà phê đá xay kết hợp thạch cà phê và kem.", map[string]float64{"S": 49000, "M": 59000, "L": 65000}},
	{"CARAMEL-PHIN-FREEZE", "FREEZE", "Caramel Phin Freeze", "Cà phê đá xay kết hợp sốt caramel, thạch cà phê và kem.", map[string]float64{"S": 55000, "M": 65000, "L": 69000}},
}

func LoadStores() ([]domain.Store, error) {
	var stores []domain.Store
	if err := json.Unmarshal(storesSeed, &stores); err != nil {
		return nil, fmt.Errorf("load retail stores seed: %w", err)
	}
	if len(stores) == 0 {
		return nil, errors.New("retail stores seed is empty")
	}
	seenCodes := map[string]struct{}{}
	for _, store := range stores {
		if store.ID == "" || store.Code == "" || store.Name == "" || store.Status == "" {
			return nil, fmt.Errorf("invalid retail store seed: id=%q code=%q name=%q status=%q", store.ID, store.Code, store.Name, store.Status)
		}
		if _, exists := seenCodes[store.Code]; exists {
			return nil, fmt.Errorf("duplicate retail store code %q", store.Code)
		}
		seenCodes[store.Code] = struct{}{}
	}
	return stores, nil
}

func LoadMenu() (domain.Menu, []domain.MenuItem, error) {
	menu := domain.Menu{
		ID:          DefaultMenuID,
		Code:        "COFFEE-DEFAULT",
		Name:        "Runtime Roasters Coffee Menu",
		Description: "Coffee-only demo menu derived from docs/requirements/SAMPLE_MENU.md.",
		Status:      domain.MenuStatusActive,
	}
	sizes := []string{"S", "M", "L"}
	items := make([]domain.MenuItem, 0, 42)
	displayOrder := 1
	for _, spec := range menuSpecs {
		for _, size := range sizes {
			price, exists := spec.prices[size]
			if !exists {
				continue
			}
			stockSKU, coffeeType := "BEAN-ROBUSTA-ROASTED", "ROBUSTA"
			if spec.category == "ESPRESSO" {
				stockSKU, coffeeType = "BEAN-ARABICA-ROASTED", "ARABICA"
			}
			id := "MI-" + spec.group + "-" + size
			items = append(items, domain.MenuItem{
				ID: id, MenuID: menu.ID, Code: id, ProductGroupCode: spec.group,
				CategoryCode: spec.category, Size: size, Name: spec.name,
				Description: spec.description, StockSKU: stockSKU, CoffeeType: coffeeType,
				Price: price, ConsumptionQuantity: consumptionFor(spec.category, size),
				ConsumptionUnit: UnitGram, Active: true, DisplayOrder: displayOrder,
			})
			displayOrder++
		}
	}
	if len(items) != 42 {
		return domain.Menu{}, nil, fmt.Errorf("menu item count = %d, want 42", len(items))
	}
	return menu, items, nil
}

func BuildDemoData(stores []domain.Store, items []domain.MenuItem) (DemoData, error) {
	if len(stores) != 5 || len(items) != 42 {
		return DemoData{}, fmt.Errorf("build retail demo data: stores=%d items=%d", len(stores), len(items))
	}
	referenceTime := time.Date(2026, 6, 1, 8, 0, 0, 0, time.FixedZone("ICT", 7*60*60))
	data := DemoData{}
	lotByStoreSKU := map[string]int{}
	for storeIndex, store := range stores {
		warehouse := warehouseFor(store.Code)
		for skuIndex, stockSKU := range []string{"BEAN-ROBUSTA-ROASTED", "BEAN-ARABICA-ROASTED"} {
			bean := "RB"
			if strings.Contains(stockSKU, "ARABICA") {
				bean = "AR"
			}
			lotID := fmt.Sprintf("LOT-%s-%s-001", store.Code, bean)
			receivedAt := referenceTime.Add(-time.Duration((storeIndex+2)*24+skuIndex) * time.Hour)
			data.Lots = append(data.Lots, domain.InventoryLot{
				ID: lotID, StoreID: store.ID, StockSKU: stockSKU,
				SourceBatchID:     fmt.Sprintf("BATCH-%s-%s-001", warehouse, bean),
				SourceHarvestID:   fmt.Sprintf("HARVEST-%s-001", bean),
				SourceWarehouseID: warehouse, ReceivedQuantity: 5000, AvailableQuantity: 5000,
				Unit: UnitGram, Status: domain.InventoryLotStatusAvailable, ReceivedAt: receivedAt,
			})
			lotByStoreSKU[store.ID+"|"+stockSKU] = len(data.Lots) - 1
			data.Movements = append(data.Movements, domain.StockMovement{
				ID: deterministicUUID("movement-received-" + lotID), StoreID: store.ID,
				InventoryLotID: lotID, StockSKU: stockSKU,
				MovementType: domain.StockMovementReceived, QuantityDelta: 5000, Unit: UnitGram,
				ReferenceType: domain.StockReferenceDemo, ReferenceID: "DEMO-RECEIPT-" + lotID,
				OccurredAt: receivedAt,
			})
		}
	}

	for storeIndex, store := range stores {
		for saleIndex := 0; saleIndex < 5; saleIndex++ {
			item := items[(storeIndex*7+saleIndex*5)%len(items)]
			lotIndex, exists := lotByStoreSKU[store.ID+"|"+item.StockSKU]
			if !exists || data.Lots[lotIndex].AvailableQuantity < item.ConsumptionQuantity {
				return DemoData{}, fmt.Errorf("insufficient demo inventory for store=%s item=%s", store.Code, item.ID)
			}
			data.Lots[lotIndex].AvailableQuantity -= item.ConsumptionQuantity
			lot := data.Lots[lotIndex]
			soldAt := referenceTime.Add(-time.Duration(storeIndex*5+saleIndex+1) * time.Hour)
			saleKey := fmt.Sprintf("%s-%02d", store.Code, saleIndex+1)
			saleID := deterministicUUID("sale-" + saleKey)
			itemID := deterministicUUID("sale-item-" + saleKey)
			productID := fmt.Sprintf("RR-CUP-%s-%04d", store.Code, saleIndex+1)
			data.Sales = append(data.Sales, domain.Sale{
				ID: saleID, StoreID: store.ID, InvoiceNo: "INV-DEMO-" + saleKey,
				Status: domain.SaleStatusCompleted, Subtotal: item.Price, TotalAmount: item.Price,
				IdempotencyKey: "demo-sale-" + strings.ToLower(saleKey), SoldAt: soldAt,
			})
			data.SaleItems = append(data.SaleItems, domain.SaleItem{
				ID: itemID, SaleID: saleID, StoreID: store.ID, MenuItemID: item.ID,
				InventoryLotID: lot.ID, ProductID: productID, TraceCode: productID,
				ProductName: item.Name, SKU: item.StockSKU, Size: item.Size, UnitPrice: item.Price,
				ConsumedQuantity: item.ConsumptionQuantity, ConsumedUnit: UnitGram,
				PublicURL: "/trace/" + productID, SoldAt: soldAt,
			})
			data.Movements = append(data.Movements, domain.StockMovement{
				ID: deterministicUUID("movement-sold-" + saleKey), StoreID: store.ID,
				InventoryLotID: lot.ID, StockSKU: item.StockSKU,
				MovementType: domain.StockMovementSold, QuantityDelta: -item.ConsumptionQuantity,
				Unit: UnitGram, ReferenceType: domain.StockReferenceSaleItem,
				ReferenceID: itemID, OccurredAt: soldAt,
			})
		}
	}
	return data, nil
}

func consumptionFor(category, size string) float64 {
	if category == "PHINDI" {
		return map[string]float64{"S": 16, "M": 20, "L": 24}[size]
	}
	return map[string]float64{"S": 18, "M": 22, "L": 25}[size]
}

func warehouseFor(storeCode string) string {
	switch storeCode {
	case "HK", "CG":
		return "WAREHOUSE-HN-001"
	case "D1", "D7":
		return "WAREHOUSE-HCM-001"
	default:
		return "WAREHOUSE-DN-001"
	}
}

func deterministicUUID(key string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("runtime-roasters/rr-urg-07a/"+key)).String()
}
