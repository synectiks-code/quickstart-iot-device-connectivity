package dbadapter

import (
	//core "cloudrack-lambda-core/core"
	//avl "cloudrack-lambda-core/avail/model"
	cfg "cloudrack-lambda-core/config/model"
	inv "cloudrack-lambda-core/inv/model"
	"errors"

	//"strconv"
	aurora "cloudrack-lambda-core/aurora"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var INVENTORY_DB_TABLE string = "inventory"
var INVENTORY_DB_HORIZON int = 701
var INVENTORY_DB_TIMEFORMAT string = "20060102"
var INVENTORY_DB_COLUMN_HOTEL_CODE string = "hotel_code"
var INVENTORY_DB_COLUMN_PRODUCT_TYPE string = "product_type"
var INVENTORY_DB_COLUMN_PRODUCT_CODE string = "product_code"
var INVENTORY_DB_COLUMN_ALLOTMENT_CODE string = "allotment_code"
var INVENTORY_DB_COLUMN_STATUS string = "status"
var INVENTORY_DB_COLUMN_BASE_COUNT string = "base_count"
var INVENTORY_DB_COLUMN_TYPE_BIGINT string = "BIGINT"
var INVENTORY_DB_COLUMN_TYPE_INT string = "INT"
var INVENTORY_DB_COLUMN_TYPE_VARCHAR_100 string = "VARCHAR(100)"
var INVENTORY_STATUS_COUNT string = "count"
var INVENTORY_STATUS_HELD string = "held"
var INVENTORY_ACTION_HOLD string = "hold"
var INVENTORY_ACTION_RELEASE string = "release"

func UpdateInventory(invDb aurora.DBConfig, hotel cfg.Hotel) error {
	err := CheckInventorySchema(invDb)
	if err != nil {
		log.Printf("[INVENTORY] Error while creating inventory schema %+v", err)
		return err
	}
	log.Printf("[INVENTORY] Updating inventory for configuration object %+v", hotel)
	curInv := GetHotelInventory(invDb, hotel.Code)
	newInv := ComputeInventoryRecords(hotel)
	toAdd, toUpdate, toDelete := computeInventoryDiff(curInv, newInv)
	err = AddNewProducts(invDb, toAdd)
	if err != nil {
		log.Printf("[INVENTORY] Error while adding new products %+v", err)
	}
	err = UpdateExistingProducts(invDb, toUpdate)
	if err != nil {
		log.Printf("[INVENTORY] Error while updating products %+v", err)
	}
	err = DeleteProducts(invDb, toDelete)
	if err != nil {
		log.Printf("[INVENTORY] Error while deleting products %+v", err)
	}
	return err
}

//returns the list of products and counts for a given hotel and date range
func GetHotelInventoryForDateRange(invDb aurora.DBConfig, hotelCodes []string, startDate string, endDate string) (map[string]inv.Product, error) {
	if startDate == endDate {
		return map[string]inv.Product{}, errors.New("Start and End date cannot be identical")
	}
	sql := BuildGetHotelInventoryForDateRangeSql(hotelCodes, startDate, endDate)
	vals, err := invDb.Query(sql)
	if err != nil {
		return map[string]inv.Product{}, err
	}
	log.Printf("[INVENTORY] aurora result for GetHotelInventoryForDateRange: %+v", vals)
	cols := auroraGenericResultToAuroraRecord(vals, startDate, endDate)
	log.Printf("[INVENTORY] aurora records for GetHotelInventoryForDateRange: %+v", cols)
	resMap := auroraResultToProducts(cols)
	if len(resMap) == 0 {
		return map[string]inv.Product{}, errors.New("Invalid Hotel Codes : " + strings.Join(hotelCodes, ","))
	}
	return auroraResultToProducts(cols), nil
}

//Returns the full content of the inventory for a given hotel
func GetHotelInventory(invDb aurora.DBConfig, hotelCode string) map[string]inv.Product {
	sql := buildGetInventorySql(hotelCode)
	vals, err := invDb.Query(sql)
	if err != nil {
		log.Printf("[INVENTORY] error GetHotelInventory: %+v", err)
	}
	cols := auroraGenericResultToAuroraRecord(vals, startDate(), endDate())
	log.Printf("[INVENTORY] res for get Inventory: %+v", cols)
	return auroraResultToProducts(cols)
}

func PerformInventoryOperation(invDb aurora.DBConfig, invRq inv.InventoryRequest) (inv.ResWrapper, error) {
	invDb.StartTransaction()
	products, err := GetHotelInventoryForDateRange(invDb, invRq.HotelCodes(), invRq.MinStartDate(), invRq.MaxEndDate())
	if err != nil {
		log.Printf("[INVENTORY] error for GetHotelInventoryForDateRange: %+v", err)
		invDb.Rollback()
		return inv.ResWrapper{}, err
	}
	err = checkInventoryCounts(invRq, products)
	if err != nil {
		log.Printf("[INVENTORY] error for PerformInventoryOperation: %+v", err)
		invDb.Rollback()
		return inv.ResWrapper{}, err
	}
	for _, action := range invRq.Actions {
		sql := BuildPerformInventoryOperationsSql(action, products)
		vals, err2 := invDb.Query(sql)
		if err2 != nil {
			log.Printf("[INVENTORY] error for PerformInventoryOperation: %+v", err)
			invDb.Rollback()
			return inv.ResWrapper{}, err2
		}
		log.Printf("[INVENTORY] res for PerformInventoryOperation: %+v", vals)
	}
	invDb.Commit()
	return inv.ResWrapper{Status: "SUCCESS"}, nil
}

//check that inventory is available for the products
//we do that in the code to avoid load on DB and due to the fact that the DB driver does snot
//return an error if the inventory update fails
func checkInventoryCounts(invRq inv.InventoryRequest, products map[string]inv.Product) error {
	for _, action := range invRq.Actions {
		if products[action.Type+action.Code].Code == "" {
			return errors.New(fmt.Sprintf("Product %s does not exists in hotel %v inventory ", action.Code, invRq.HotelCodes()))
		}
		for date, prod := range products[action.Type+action.Code].Counts {
			if action.Action == INVENTORY_ACTION_HOLD {
				if prod[inv.INVENTORY_ALLOTMENT_GINV].BaseCount == 0 {
					return errors.New(fmt.Sprintf("Product %s is not available at date %s (zero count)", action.Code, date))
				}
				if prod[inv.INVENTORY_ALLOTMENT_GINV].Held >= prod[inv.INVENTORY_ALLOTMENT_GINV].BaseCount {
					return errors.New(fmt.Sprintf("Product %s is no no longer available at date %s", action.Code, date))
				}
			} else {
				if prod[inv.INVENTORY_ALLOTMENT_GINV].Held == 0 {
					return errors.New(fmt.Sprintf("Product %s is can't be released on date %s since no product is held", action.Code, date))
				}
			}
		}
	}
	return nil
}

//convert a generic map struct into Aurora reccord.
//used when aurora columns can't be hardcoded
func auroraGenericResultToAuroraRecord(results []map[string]interface{}, startDate string, endDate string) []inv.AuroraInventoryRecord {
	records := make([]inv.AuroraInventoryRecord, 0, 0)
	for _, res := range results {
		rec := inv.AuroraInventoryRecord{
			HotelCode:     strconv.FormatInt(res[INVENTORY_DB_COLUMN_HOTEL_CODE].(int64), 10),
			ProductType:   res[INVENTORY_DB_COLUMN_PRODUCT_TYPE].(string),
			ProductCode:   res[INVENTORY_DB_COLUMN_PRODUCT_CODE].(string),
			AllotmentCode: res[INVENTORY_DB_COLUMN_ALLOTMENT_CODE].(string),
			Status:        res[INVENTORY_DB_COLUMN_STATUS].(string),
			CountByDate:   make(map[string]int),
			HeldByDate:    make(map[string]int),
		}
		for _, date := range geterateDateRange(startDate, endDate) {
			if val, ok := res[date].(int64); ok {
				count := int(val)
				if rec.Status == INVENTORY_STATUS_COUNT {
					rec.CountByDate[date] = count
				} else {
					rec.HeldByDate[date] = count
				}
			} else {
				log.Printf("[INVENTORY][ERROR] counter filed in DB holds the wrong value. probably null: %+v", res[date])
			}
		}
		records = append(records, rec)
	}
	return records
}

//convernt auroraReccord struct to a Product struct
func auroraResultToProducts(records []inv.AuroraInventoryRecord) map[string]inv.Product {
	log.Printf("[INVENTORY] auroraResultToProducts Building product map")
	products := make(map[string]inv.Product)
	for _, rec := range records {
		log.Printf("[INVENTORY] processing record %+v", rec)
		id := rec.ProductType + rec.ProductCode
		if products[id].Code == "" {
			products[id] = inv.Product{
				HotelCode: rec.HotelCode,
				Code:      rec.ProductCode,
				Type:      rec.ProductType,
				Counts:    make(map[string]map[string]inv.Allotment),
			}
		}
		nonEmptyMap := rec.CountByDate
		if len(rec.HeldByDate) > 0 {
			nonEmptyMap = rec.HeldByDate
		}
		for date, _ := range nonEmptyMap {
			if len(products[id].Counts[date]) == 0 {
				//log.Printf("[INVENTORY] auroraResultToProducts  initializing allotment map for date %v", date)
				products[id].Counts[date] = make(map[string]inv.Allotment)
				products[id].Counts[date][rec.AllotmentCode] = inv.Allotment{
					Code: rec.AllotmentCode,
					Held: rec.HeldByDate[date],
				}
			}

			if len(rec.HeldByDate) > 0 {
				//log.Printf("[INVENTORY] record %v has HELD inventory", rec)
				if val, ok := products[id].Counts[date][rec.AllotmentCode]; ok {
					val.Held = rec.HeldByDate[date]
					products[id].Counts[date][rec.AllotmentCode] = val
				} else {
					products[id].Counts[date][rec.AllotmentCode] = inv.Allotment{
						Code: rec.AllotmentCode,
						Held: rec.HeldByDate[date],
					}
				}
			}
			if len(rec.CountByDate) > 0 {
				//log.Printf("[INVENTORY] record %v has COUNT inventory", rec)
				if val, ok := products[id].Counts[date][rec.AllotmentCode]; ok {
					val.BaseCount = rec.CountByDate[date]
					products[id].Counts[date][rec.AllotmentCode] = val
				} else {
					products[id].Counts[date][rec.AllotmentCode] = inv.Allotment{
						Code:      rec.AllotmentCode,
						BaseCount: rec.CountByDate[date],
					}
				}
			}
			//log.Printf("[INVENTORY] productsMap: %+v", products)
		}

	}
	return products
}

func ComputeInventoryRecords(hotel cfg.Hotel) map[string]inv.Product {
	counts := buildCounts(hotel)
	products := buildInventoryImage(hotel.Code, counts, hotel.Rules)
	return products
}

func buildInventoryImage(hotelCode string, counts map[string]map[string]int, rules []cfg.HotelBusinessRule) map[string]inv.Product {
	products := make(map[string]inv.Product)
	for productType, countMap := range counts {
		for id, count := range countMap {
			products[productType+id] = inv.Product{
				HotelCode: hotelCode,
				Code:      id,
				Type:      productType,
				Counts:    generateDateInventory(productType, id, count, rules),
			}
		}
	}
	return products
}

func generateDateInventory(productType string, productId string, count int, rules []cfg.HotelBusinessRule) map[string]map[string]inv.Allotment {
	inventory := make(map[string]map[string]inv.Allotment, 0)
	inventory[inv.INVENTORY_ALLOTMENT_GINV] = make(map[string]inv.Allotment, 0)
	dates := generateDatesToHorizon()
	for _, date := range dates {
		inventory[inv.INVENTORY_ALLOTMENT_GINV][date] = inv.Allotment{
			Code:      inv.INVENTORY_ALLOTMENT_GINV,
			BaseCount: computeInventoryAmount(productType, productId, date, count, rules),
		}
	}
	return inventory
}

//Compute the apropriate inventory count based on date product and business rules
func computeInventoryAmount(productType string, productId string, date string, count int, rules []cfg.HotelBusinessRule) int {
	//TODO: Implement inventory business rule computation
	return count
}

func buildCounts(hotel cfg.Hotel) map[string]map[string]int {
	counts := make(map[string]map[string]int)
	counts[inv.INVENTORY_TYPE_ROOM] = make(map[string]int)
	counts[inv.INVENTORY_TYPE_SELLABLE] = make(map[string]int)
	log.Printf("[INVENTORY] Computing inventory counts from config")
	for _, building := range hotel.Buildings {
		for _, floor := range building.Floors {
			for _, room := range floor.Rooms {
				counts[inv.INVENTORY_TYPE_ROOM][buildProductCode(room)]++
			}
		}
	}
	for _, sellable := range hotel.Sellables {
		counts[inv.INVENTORY_TYPE_SELLABLE][sellable.Code] = int(sellable.Quantity)
	}
	log.Printf("[INVENTORY] Inventory counts %+v", counts)
	return counts
}

func buildProductCode(room cfg.HotelRoom) string {
	code := room.Type
	if len(room.Attributes) > 0 {
		for _, tag := range room.Attributes {
			if tag.Code != "" {
				code = code + "-" + tag.Code
			}
		}
	}
	return code
}

func computeInventoryDiff(curInv map[string]inv.Product, newInv map[string]inv.Product) ([]inv.Product, []inv.Product, []inv.Product) {
	log.Printf("[INVENTORY] Computing inventory changes to make")
	toCreate := make([]inv.Product, 0, 0)
	toUpdate := make([]inv.Product, 0, 0)
	toDelete := make([]inv.Product, 0, 0)
	for key, newProd := range newInv {
		if _, ok := curInv[key]; !ok {
			toCreate = append(toCreate, newProd)
		} else {
			toUpdate = append(toUpdate, newProd)
		}
	}
	for key, oldProd := range curInv {
		if _, ok := newInv[key]; !ok {
			toDelete = append(toDelete, oldProd)
		}
	}
	return toCreate, toUpdate, toDelete
}

func AddNewProducts(invDb aurora.DBConfig, toAdd []inv.Product) error {
	log.Printf("[INVENTORY] Adding %v products", len(toAdd))
	for _, product := range toAdd {
		err := AddProduct(invDb, product)
		if err != nil {
			log.Printf("[INVENTORY] Error during product addition %+v", err)
			return err
		}
	}
	return nil
}

func UpdateExistingProducts(invDb aurora.DBConfig, toUpdate []inv.Product) error {
	log.Printf("[INVENTORY] Updating %v products", len(toUpdate))
	for _, product := range toUpdate {
		err := UpdateProduct(invDb, product)
		if err != nil {
			log.Printf("[INVENTORY] Error during product update %+v", err)
			return err
		}
	}
	return nil

}
func DeleteProducts(invDb aurora.DBConfig, toDelete []inv.Product) error {
	log.Printf("[INVENTORY] Deleting %v products", len(toDelete))
	for _, product := range toDelete {
		err := DeleteProduct(invDb, product)
		if err != nil {
			log.Printf("[INVENTORY] Error during product deletion %+v", err)
			return err
		}
	}
	return nil
}

func AddProduct(invDb aurora.DBConfig, product inv.Product) error {
	log.Printf("[INVENTORY] Adding product to inventory DB\n")
	_, err := invDb.Query(buildAddProductSql(invDb, product, INVENTORY_STATUS_COUNT))
	if err != nil {
		log.Printf("[INVENTORY] Error during product addition (count row) %+v", err)
	}
	_, err = invDb.Query(buildAddProductSql(invDb, product, INVENTORY_STATUS_HELD))
	if err != nil {
		log.Printf("[INVENTORY] Error during product addition (held row) %+v", err)
	}
	return err
}
func UpdateProduct(invDb aurora.DBConfig, product inv.Product) error {
	log.Printf("[INVENTORY] updating product to inventory DB\n")
	_, err := invDb.Query(BuildUpdateProductSql(product))
	if err != nil {
		log.Printf("[INVENTORY] Error during product update %+v", err)
	}
	return err
}
func DeleteProduct(invDb aurora.DBConfig, product inv.Product) error {
	log.Printf("[INVENTORY] Deleting product to inventory DB\n")
	_, err := invDb.Query(buildDeleteProductSql(invDb, product, INVENTORY_STATUS_COUNT))
	if err != nil {
		log.Printf("[INVENTORY] Error during product deletion (count row)%+v", err)
	}
	return err
}

func CheckInventorySchema(invDb aurora.DBConfig) error {
	log.Printf("[INVENTORY] Checkking for inventory schema existence\n")
	if !invDb.HasTable(INVENTORY_DB_TABLE) {
		err := CreateInventoryTable(invDb)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[INVENTORY] table %s exists in database", INVENTORY_DB_TABLE)
	}
	err := checkColumns(invDb)
	return err
}

func checkColumns(invDb aurora.DBConfig) error {
	log.Printf("[INVENTORY] Checking inventory horizon")
	cols, _ := invDb.GetColumns(INVENTORY_DB_TABLE)
	log.Printf("[INVENTORY] res for show column: %+v", cols)
	lastDate, _ := time.Parse(INVENTORY_DB_TIMEFORMAT, cols[len(cols)-1].Field)
	//Extending inventory horizon
	newHorizon := time.Now().AddDate(0, 0, INVENTORY_DB_HORIZON)
	diff := int(newHorizon.Sub(lastDate).Hours()/24.0) + 1
	if diff > 0 {
		log.Printf("[INVENTORY] last inventory date currently is %v, adding %v days", cols[len(cols)-1].Field, diff)
		err := extendsHorizon(invDb, lastDate, diff)
		if err != nil {
			log.Printf("[INVENTORY] Error during horizon extension %+v", err)
			return err
		}
	} else {
		log.Printf("[INVENTORY] horizon up to date at %v", cols[len(cols)-1].Field)
	}
	//TODO: Removing past inventory days
	return nil
}

func extendsHorizon(invDb aurora.DBConfig, lastDate time.Time, diff int) error {
	log.Printf("[INVENTORY] Extending Inventory horizon table\n")
	_, err := invDb.Query(buildExtendHorizonSql(lastDate, diff))
	if err != nil {
		log.Printf("[INVENTORY] Error during horizon extension (adding colmns) %+v", err)
	}
	_, err = invDb.Query(buildExtendHorizonCopyColumnsSql(lastDate, diff))
	if err != nil {
		log.Printf("[INVENTORY] Error during horizon extension (setting columns values) %+v", err)
	}
	return err
}

func CreateInventoryTable(invDb aurora.DBConfig) error {
	log.Printf("[INVENTORY] Creating Inventory table\n")
	_, err := invDb.Query(buildCreatTableSql(invDb))
	if err != nil {
		log.Printf("[INVENTORY] Error during Create table %+v", err)
	}
	return err
}

func generateInventoryColumns() []string {
	cols := make([]string, 0, 0)
	cols = append(cols, INVENTORY_DB_COLUMN_HOTEL_CODE+" "+INVENTORY_DB_COLUMN_TYPE_BIGINT)
	cols = append(cols, INVENTORY_DB_COLUMN_PRODUCT_TYPE+" "+INVENTORY_DB_COLUMN_TYPE_VARCHAR_100)
	cols = append(cols, INVENTORY_DB_COLUMN_PRODUCT_CODE+" "+INVENTORY_DB_COLUMN_TYPE_VARCHAR_100)
	cols = append(cols, INVENTORY_DB_COLUMN_ALLOTMENT_CODE+" "+INVENTORY_DB_COLUMN_TYPE_VARCHAR_100)
	cols = append(cols, INVENTORY_DB_COLUMN_STATUS+" "+INVENTORY_DB_COLUMN_TYPE_VARCHAR_100)
	cols = append(cols, INVENTORY_DB_COLUMN_BASE_COUNT+" "+INVENTORY_DB_COLUMN_TYPE_INT)
	dates := generateDatesToHorizon()
	for _, date := range dates {
		cols = append(cols, "`"+date+"`"+" "+INVENTORY_DB_COLUMN_TYPE_INT)
	}
	return cols
}

//////////////////////////
// SQL Queries
//////////////////////////

func BuildGetHotelInventoryForDateRangeSql(hotelCodes []string, startDate string, endDate string) string {
	dates := geterateDateRange(startDate, endDate)
	columns := aurora.WrapField(strings.Join(dates, "`,`"))
	tpl := ""
	hotelCodeStr := ""
	if len(hotelCodes) == 1 {
		tpl = `SELECT %s,%s,%s,%s,%s,%s FROM %s WHERE %s="%s";`
		hotelCodeStr = hotelCodes[0]
	} else {
		tpl = `SELECT %s,%s,%s,%s,%s,%s FROM %s WHERE %s IN (%s);`
		hotelCodeStr = aurora.WrapVal(strings.Join(hotelCodes, "\",\""))
	}
	query := fmt.Sprintf(tpl,
		INVENTORY_DB_COLUMN_HOTEL_CODE,
		INVENTORY_DB_COLUMN_PRODUCT_TYPE,
		INVENTORY_DB_COLUMN_PRODUCT_CODE,
		INVENTORY_DB_COLUMN_ALLOTMENT_CODE,
		INVENTORY_DB_COLUMN_STATUS,
		columns,
		INVENTORY_DB_TABLE,
		INVENTORY_DB_COLUMN_HOTEL_CODE,
		hotelCodeStr)
	log.Printf("[INVENTORY] buildGetHotelInventoryForDateRangeSql query: %+v", query)
	return query
}

func buildGetInventorySql(hotelCode string) string {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s=%s;`,
		INVENTORY_DB_TABLE,
		aurora.WrapField(INVENTORY_DB_COLUMN_HOTEL_CODE),
		aurora.WrapVal(hotelCode))
	log.Printf("[INVENTORY] Get Inventory query: %+v", query)
	return query
}

func buildExtendHorizonSql(lastDate time.Time, diff int) string {
	cols := make([]string, 0, 0)
	for i := 0; i < diff; i++ {
		t := lastDate.AddDate(0, 0, i+1)
		cols = append(cols, "ADD "+t.Format("`"+INVENTORY_DB_TIMEFORMAT+"`")+" "+INVENTORY_DB_COLUMN_TYPE_INT)
	}
	query := fmt.Sprintf(`ALTER TABLE %s %s;`, INVENTORY_DB_TABLE, strings.Join(cols, ","))
	log.Printf("[INVENTORY] Extend horizon query for adding columns: %+v", query)
	return query
}

func buildExtendHorizonCopyColumnsSql(lastDate time.Time, diff int) string {
	cols := make([]string, 0, 0)
	for i := 0; i < diff; i++ {
		t := lastDate.AddDate(0, 0, i+1)
		cols = append(cols, t.Format("`"+INVENTORY_DB_TIMEFORMAT+"`")+"="+lastDate.Format("`"+INVENTORY_DB_TIMEFORMAT+"`"))
	}
	query := fmt.Sprintf(`UPDATE %s SET %s;`, INVENTORY_DB_TABLE, strings.Join(cols, ","))
	log.Printf("[INVENTORY] Extend horizon query for setting default values: %+v", query)
	return query
}

func buildCreatTableSql(invDb aurora.DBConfig) string {
	columns := generateInventoryColumns()
	query := fmt.Sprintf(`CREATE TABLE %s (%s, INDEX (%s), UNIQUE(%s,%s,%s,%s,%s))  ENGINE=%s;`,
		INVENTORY_DB_TABLE, strings.Join(columns, ","),
		INVENTORY_DB_COLUMN_HOTEL_CODE,
		INVENTORY_DB_COLUMN_HOTEL_CODE,
		INVENTORY_DB_COLUMN_PRODUCT_TYPE,
		INVENTORY_DB_COLUMN_PRODUCT_CODE,
		INVENTORY_DB_COLUMN_ALLOTMENT_CODE,
		INVENTORY_DB_COLUMN_STATUS,
		invDb.Engine())
	log.Printf("[INVENTORY] Create table query: %+v", query)
	return query
}

func buildAddProductSql(invDb aurora.DBConfig, product inv.Product, status string) string {
	dates := generateDatesToHorizon()
	counts := make([]string, 0, 0)
	for _, date := range dates {
		if status == INVENTORY_STATUS_COUNT {
			counts = append(counts, strconv.Itoa(product.Counts[inv.INVENTORY_ALLOTMENT_GINV][date].BaseCount))
		} else if status == INVENTORY_STATUS_HELD {
			counts = append(counts, "0")
		} else {
			panic("invalid status code " + status)
		}
	}
	query := fmt.Sprintf(`INSERT INTO %s (%s,%s,%s,%s,%s,%s) VALUES (%s,%s,%s,%s,%s,%s);`, INVENTORY_DB_TABLE,
		INVENTORY_DB_COLUMN_HOTEL_CODE,
		INVENTORY_DB_COLUMN_PRODUCT_TYPE,
		INVENTORY_DB_COLUMN_PRODUCT_CODE,
		INVENTORY_DB_COLUMN_ALLOTMENT_CODE,
		INVENTORY_DB_COLUMN_STATUS,
		aurora.WrapField(strings.Join(dates, "`,`")),
		aurora.WrapVal(product.HotelCode),
		aurora.WrapVal(product.Type),
		aurora.WrapVal(product.Code),
		aurora.WrapVal(inv.INVENTORY_ALLOTMENT_GINV),
		aurora.WrapVal(status),
		strings.Join(counts, ","))
	log.Printf("[INVENTORY] Add Product query: %+v", query)
	return query
}

func BuildUpdateProductSql(product inv.Product) string {
	dates := generateDatesToHorizon()
	counts := make([]string, 0, 0)
	for _, date := range dates {
		counts = append(counts, aurora.WrapField(date)+"="+strconv.Itoa(product.Counts[inv.INVENTORY_ALLOTMENT_GINV][date].BaseCount))
	}
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE  %s=%s AND %s=%s AND %s=%s AND %s=%s AND %s=%s;`, INVENTORY_DB_TABLE,
		strings.Join(counts, ","),
		INVENTORY_DB_COLUMN_HOTEL_CODE,
		aurora.WrapVal(product.HotelCode),
		INVENTORY_DB_COLUMN_PRODUCT_TYPE,
		aurora.WrapVal(product.Type),
		INVENTORY_DB_COLUMN_PRODUCT_CODE,
		aurora.WrapVal(product.Code),
		INVENTORY_DB_COLUMN_ALLOTMENT_CODE,
		aurora.WrapVal(inv.INVENTORY_ALLOTMENT_GINV),
		INVENTORY_DB_COLUMN_STATUS,
		aurora.WrapVal(INVENTORY_STATUS_COUNT))
	log.Printf("[INVENTORY] Update Product query: %+v", query)
	return query
}

func buildDeleteProductSql(invDb aurora.DBConfig, product inv.Product, status string) string {
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s=%s AND %s=%s AND %s=%s AND %s=%s;`, INVENTORY_DB_TABLE,
		INVENTORY_DB_COLUMN_HOTEL_CODE,
		aurora.WrapVal(product.HotelCode),
		INVENTORY_DB_COLUMN_PRODUCT_TYPE,
		aurora.WrapVal(product.Type),
		INVENTORY_DB_COLUMN_PRODUCT_CODE,
		aurora.WrapVal(product.Code),
		INVENTORY_DB_COLUMN_ALLOTMENT_CODE,
		aurora.WrapVal(inv.INVENTORY_ALLOTMENT_GINV))
	log.Printf("[INVENTORY] Delete Product query: %+v", query)
	return query
}

func BuildPerformInventoryOperationsSql(action inv.InventoryAction, products map[string]inv.Product) string {
	dates := geterateDateRange(action.StartDate, action.EndDate)
	counts := make([]string, 0, 0)
	where := make([]string, 0, 0)
	for _, date := range dates {
		if action.Action == INVENTORY_ACTION_HOLD {
			counts = append(counts, aurora.WrapField(date)+"="+aurora.WrapField(date)+"+1")
			where = append(where, aurora.WrapField(date)+"<"+strconv.Itoa(products[action.Type+action.Code].Counts[date][inv.INVENTORY_ALLOTMENT_GINV].BaseCount))
		} else {
			counts = append(counts, aurora.WrapField(date)+"="+aurora.WrapField(date)+"-1")
			//0 to be replaced by overbooking allowance
			where = append(where, aurora.WrapField(date)+">0")
		}
	}
	query := fmt.Sprintf(`UPDATE %s SET %s WHERE  %s=%s AND %s=%s AND %s=%s AND %s=%s AND %s=%s AND %s;`, INVENTORY_DB_TABLE,
		strings.Join(counts, ","),
		INVENTORY_DB_COLUMN_HOTEL_CODE,
		aurora.WrapVal(action.HotelCode),
		INVENTORY_DB_COLUMN_PRODUCT_TYPE,
		aurora.WrapVal(action.Type),
		INVENTORY_DB_COLUMN_PRODUCT_CODE,
		aurora.WrapVal(action.Code),
		INVENTORY_DB_COLUMN_ALLOTMENT_CODE,
		aurora.WrapVal(inv.INVENTORY_ALLOTMENT_GINV),
		INVENTORY_DB_COLUMN_STATUS,
		aurora.WrapVal(INVENTORY_STATUS_HELD),
		strings.Join(where, " AND "))
	log.Printf("[INVENTORY] Perform Inventory Operation query: %+v", query)
	return query
}

//////////////////////////
//SUPPORTING FUNCTIONS
///////////////////////////
func generateDatesToHorizon() []string {
	now := time.Now()
	dates := make([]string, 0, 0)
	for i := 0; i < INVENTORY_DB_HORIZON; i++ {
		t := now.AddDate(0, 0, i)
		dates = append(dates, t.Format(INVENTORY_DB_TIMEFORMAT))
	}
	return dates
}

//inventory start date
func startDate() string {
	return time.Now().Format(INVENTORY_DB_TIMEFORMAT)
}

//inventory endDate date
func endDate() string {
	now := time.Now()
	t := now.AddDate(0, 0, INVENTORY_DB_HORIZON-1)
	return t.Format(INVENTORY_DB_TIMEFORMAT)
}

func geterateDateRange(startDate string, endDate string) []string {
	start, _ := time.Parse(INVENTORY_DB_TIMEFORMAT, startDate)
	end, _ := time.Parse(INVENTORY_DB_TIMEFORMAT, endDate)
	diff := int(end.Sub(start).Hours() / 24.0)
	dates := make([]string, 0, 0)
	for i := 0; i < diff; i++ {
		t := start.AddDate(0, 0, i)
		dates = append(dates, t.Format(INVENTORY_DB_TIMEFORMAT))
	}
	return dates
}
