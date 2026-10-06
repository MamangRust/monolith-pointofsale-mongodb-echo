package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------------------------------------------------------------------
// Core
// ---------------------------------------------------------------------------

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoleID    int32              `bson:"role_id"`
	RoleName  string             `bson:"role_name"`
	CreatedAt *time.Time         `bson:"created_at"`
	UpdatedAt *time.Time         `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func (Role) CollectionName() string { return "roles" }

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	UserID           int32              `bson:"user_id"`
	Firstname        string             `bson:"firstname"`
	Lastname         string             `bson:"lastname"`
	Email            string             `bson:"email"`
	Password         string             `bson:"password"`
	VerificationCode string             `bson:"verification_code"`
	IsVerified       *bool              `bson:"is_verified"`
	CreatedAt        *time.Time         `bson:"created_at"`
	UpdatedAt        *time.Time         `bson:"updated_at"`
	DeletedAt        *time.Time         `bson:"deleted_at,omitempty"`
}

func (User) CollectionName() string { return "users" }

type UserRole struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserRoleID int32              `bson:"user_role_id"`
	UserID     int32              `bson:"user_id"`
	RoleID     int32              `bson:"role_id"`
	CreatedAt  *time.Time         `bson:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at"`
	DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}

func (UserRole) CollectionName() string { return "user_roles" }

type RefreshToken struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	RefreshTokenID int32              `bson:"refresh_token_id"`
	UserID         int32              `bson:"user_id"`
	Token          string             `bson:"token"`
	Expiration     time.Time          `bson:"expiration"`
	CreatedAt      *time.Time         `bson:"created_at"`
	UpdatedAt      *time.Time         `bson:"updated_at"`
	DeletedAt      *time.Time         `bson:"deleted_at,omitempty"`
}

func (RefreshToken) CollectionName() string { return "refresh_tokens" }

type ResetToken struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ResetTokenID int32              `bson:"reset_token_id"`
	UserID       int64              `bson:"user_id"`
	Token        string             `bson:"token"`
	ExpiryDate   time.Time          `bson:"expiry_date"`
}

func (ResetToken) CollectionName() string { return "reset_tokens" }

// ---------------------------------------------------------------------------
// Cashier
// ---------------------------------------------------------------------------

type Cashier struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CashierID  int32              `bson:"cashier_id"`
	MerchantID int32              `bson:"merchant_id"`
	UserID     int32              `bson:"user_id"`
	Name       string             `bson:"name"`
	CreatedAt  *time.Time         `bson:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at"`
	DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}

func (Cashier) CollectionName() string { return "cashiers" }

// ---------------------------------------------------------------------------
// Merchant
// ---------------------------------------------------------------------------

type Merchant struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	MerchantID   int32              `bson:"merchant_id"`
	UserID       int32              `bson:"user_id"`
	Name         string             `bson:"name"`
	Description  *string            `bson:"description"`
	Address      *string            `bson:"address"`
	ContactEmail *string            `bson:"contact_email"`
	ContactPhone *string            `bson:"contact_phone"`
	Status       string             `bson:"status"`
	CreatedAt    *time.Time         `bson:"created_at"`
	UpdatedAt    *time.Time         `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

func (Merchant) CollectionName() string { return "merchants" }

type MerchantDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	DocumentID   int32              `bson:"document_id"`
	MerchantID   int32              `bson:"merchant_id"`
	DocumentType string             `bson:"document_type"`
	DocumentUrl  string             `bson:"document_url"`
	Status       string             `bson:"status"`
	Note         *string            `bson:"note"`
	UploadedAt   *time.Time         `bson:"uploaded_at"`
	CreatedAt    *time.Time         `bson:"created_at"`
	UpdatedAt    *time.Time         `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

func (MerchantDocument) CollectionName() string { return "merchant_documents" }

// ---------------------------------------------------------------------------
// Category
// ---------------------------------------------------------------------------

type Category struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CategoryID   int32              `bson:"category_id"`
	Name         string             `bson:"name"`
	Description  *string            `bson:"description"`
	SlugCategory *string            `bson:"slug_category"`
	CreatedAt    *time.Time         `bson:"created_at"`
	UpdatedAt    *time.Time         `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

func (Category) CollectionName() string { return "categories" }

// ---------------------------------------------------------------------------
// Order
// ---------------------------------------------------------------------------

type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	OrderID    int32              `bson:"order_id"`
	MerchantID int32              `bson:"merchant_id"`
	CashierID  int32              `bson:"cashier_id"`
	TotalPrice int64              `bson:"total_price"`
	CreatedAt  *time.Time         `bson:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at"`
	DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}

func (Order) CollectionName() string { return "orders" }

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	OrderItemID int32              `bson:"order_item_id"`
	OrderID     int32              `bson:"order_id"`
	ProductID   int32              `bson:"product_id"`
	Quantity    int32              `bson:"quantity"`
	Price       int64              `bson:"price"`
	CreatedAt   *time.Time         `bson:"created_at"`
	UpdatedAt   *time.Time         `bson:"updated_at"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty"`
}

func (OrderItem) CollectionName() string { return "order_items" }

// ---------------------------------------------------------------------------
// Product
// ---------------------------------------------------------------------------

type Product struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ProductID    int32              `bson:"product_id"`
	MerchantID   int32              `bson:"merchant_id"`
	CategoryID   int32              `bson:"category_id"`
	Name         string             `bson:"name"`
	Description  *string            `bson:"description"`
	Price        int32              `bson:"price"`
	CountInStock int32              `bson:"count_in_stock"`
	Brand        *string            `bson:"brand"`
	Weight       *int32             `bson:"weight"`
	SlugProduct  *string            `bson:"slug_product"`
	ImageProduct *string            `bson:"image_product"`
	Barcode      *string            `bson:"barcode"`
	CreatedAt    *time.Time         `bson:"created_at"`
	UpdatedAt    *time.Time         `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

func (Product) CollectionName() string { return "products" }

// ---------------------------------------------------------------------------
// Transaction
// ---------------------------------------------------------------------------

type Transaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	TransactionID int32              `bson:"transaction_id"`
	OrderID       int32              `bson:"order_id"`
	MerchantID    int32              `bson:"merchant_id"`
	CashierID     int32              `bson:"cashier_id"`
	TotalPrice    int64              `bson:"total_price"`
	Amount        int32              `bson:"amount"`
	ChangeAmount  *int32             `bson:"change_amount"`
	PaymentMethod string             `bson:"payment_method"`
	PaymentStatus string             `bson:"payment_status"`
	Status        string             `bson:"status"`
	CreatedAt     *time.Time         `bson:"created_at"`
	UpdatedAt     *time.Time         `bson:"updated_at"`
	DeletedAt     *time.Time         `bson:"deleted_at,omitempty"`
}

func (Transaction) CollectionName() string { return "transactions" }

// ---------------------------------------------------------------------------
// Stats rows (sqlc compatibility)
// ---------------------------------------------------------------------------

type CashierMonthlyTotalSalesRow struct {
	Year       int32 `bson:"year"`
	Month      int32 `bson:"month"`
	TotalSales int32 `bson:"total_sales"`
}

type CashierYearlyTotalSalesRow struct {
	Year       int32 `bson:"year"`
	TotalSales int32 `bson:"total_sales"`
}

type CashierMonthlySaleRow struct {
	CashierID   int32  `bson:"cashier_id"`
	CashierName string `bson:"cashier_name"`
	Month       int32  `bson:"month"`
	OrderCount  int32  `bson:"order_count"`
	TotalSales  int32  `bson:"total_sales"`
}

type CashierMonthlySaleByIdRow = CashierMonthlySaleRow

type CashierMonthlySaleByMerchantRow struct {
	MerchantID  int32  `bson:"merchant_id"`
	CashierID   int32  `bson:"cashier_id"`
	CashierName string `bson:"cashier_name"`
	Month       int32  `bson:"month"`
	OrderCount  int32  `bson:"order_count"`
	TotalSales  int32  `bson:"total_sales"`
}

type CashierYearlySaleRow struct {
	CashierID   int32  `bson:"cashier_id"`
	CashierName string `bson:"cashier_name"`
	Year        int32  `bson:"year"`
	OrderCount  int32  `bson:"order_count"`
	TotalSales  int32  `bson:"total_sales"`
}

type CategoryMonthlyTotalPriceRow struct {
	Year       int32 `bson:"year"`
	Month      int32 `bson:"month"`
	TotalPrice int32 `bson:"total_price"`
}

type CategoryMonthlyTotalPriceByIdRow = CategoryMonthlyTotalPriceRow
type CategoryMonthlyTotalPriceByMerchantRow = CategoryMonthlyTotalPriceRow

type CategoryYearlyTotalPriceRow struct {
	Year       int32 `bson:"year"`
	TotalPrice int32 `bson:"total_price"`
}

type CategoryYearlyTotalPriceByIdRow = CategoryYearlyTotalPriceRow
type CategoryYearlyTotalPriceByMerchantRow = CategoryYearlyTotalPriceRow

type CategoryMonthlyRow struct {
	CategoryID   int32  `bson:"category_id"`
	CategoryName string `bson:"category_name"`
	Month        int32  `bson:"month"`
	OrderCount   int32  `bson:"order_count"`
	ItemsSold    int32  `bson:"items_sold"`
	TotalRevenue int32  `bson:"total_revenue"`
}

type CategoryMonthlyByIdRow = CategoryMonthlyRow

type CategoryMonthlyByMerchantRow = CategoryMonthlyRow

type CategoryYearlyRow struct {
	CategoryID         int32  `bson:"category_id"`
	CategoryName       string `bson:"category_name"`
	Year               int32  `bson:"year"`
	OrderCount         int32  `bson:"order_count"`
	ItemsSold          int32  `bson:"items_sold"`
	UniqueProductsSold int32  `bson:"unique_products_sold"`
	TotalRevenue       int32  `bson:"total_revenue"`
}

type CategoryYearlyByIdRow = CategoryYearlyRow
type CategoryYearlyByMerchantRow = CategoryYearlyRow

type TransactionMonthlyAmountSuccessRow struct {
	Year         string `bson:"year"`
	Month        string `bson:"month"`
	TotalSuccess int64  `bson:"total_success"`
	TotalAmount  int32  `bson:"total_amount"`
}

type TransactionMonthlyAmountSuccessByMerchantRow = TransactionMonthlyAmountSuccessRow

type TransactionMonthlyAmountFailedRow struct {
	Year        string `bson:"year"`
	Month       string `bson:"month"`
	TotalFailed int64  `bson:"total_failed"`
	TotalAmount int32  `bson:"total_amount"`
}

type TransactionMonthlyAmountFailedByMerchantRow = TransactionMonthlyAmountFailedRow

type TransactionYearlyAmountSuccessRow struct {
	Year         string `bson:"year"`
	TotalSuccess int32  `bson:"total_success"`
	TotalAmount  int32  `bson:"total_amount"`
}

type TransactionYearlyAmountSuccessByMerchantRow = TransactionYearlyAmountSuccessRow

type TransactionYearlyAmountFailedRow struct {
	Year        string `bson:"year"`
	TotalFailed int32  `bson:"total_failed"`
	TotalAmount int32  `bson:"total_amount"`
}

type TransactionYearlyAmountFailedByMerchantRow = TransactionYearlyAmountFailedRow

type TransactionMonthlyMethodSuccessRow struct {
	Month             string  `bson:"month"`
	PaymentMethod     string  `bson:"payment_method"`
	TotalTransactions int64   `bson:"total_transactions"`
	TotalAmount       float64 `bson:"total_amount"`
}

type TransactionMonthlyMethodByMerchantSuccessRow = TransactionMonthlyMethodSuccessRow

type TransactionMonthlyMethodSuccessByMerchantRow = TransactionMonthlyMethodSuccessRow

type TransactionMonthlyMethodFailedRow struct {
	Month             string  `bson:"month"`
	PaymentMethod     string  `bson:"payment_method"`
	TotalTransactions int64   `bson:"total_transactions"`
	TotalAmount       float64 `bson:"total_amount"`
}

type TransactionMonthlyMethodFailedByMerchantRow = TransactionMonthlyMethodFailedRow

type TransactionYearlyMethodSuccessRow struct {
	Year              string  `bson:"year"`
	PaymentMethod     string  `bson:"payment_method"`
	TotalTransactions int64   `bson:"total_transactions"`
	TotalAmount       float64 `bson:"total_amount"`
}

type TransactionYearlyMethodByMerchantSuccessRow = TransactionYearlyMethodSuccessRow

type TransactionYearlyMethodSuccessByMerchantRow = TransactionYearlyMethodSuccessRow

type TransactionYearlyMethodFailedRow struct {
	Year              string  `bson:"year"`
	PaymentMethod     string  `bson:"payment_method"`
	TotalTransactions int64   `bson:"total_transactions"`
	TotalAmount       float64 `bson:"total_amount"`
}

type TransactionYearlyMethodByMerchantFailedRow = TransactionYearlyMethodFailedRow

// TransactionYearlyMethodFailedByMerchantRow is an alternate naming kept for
// sqlc/handler compatibility.
type TransactionYearlyMethodFailedByMerchantRow = TransactionYearlyMethodFailedRow

type TransactionMonthlyMethodByMerchantFailedRow = TransactionMonthlyMethodFailedRow

// ---------------------------------------------------------------------------
// Order / revenue stats rows
// ---------------------------------------------------------------------------

type MonthlyTotalRevenueRow struct {
	Year         string `bson:"year"`
	Month        string `bson:"month"`
	TotalRevenue int32  `bson:"total_revenue"`
}

type YearlyTotalRevenueRow struct {
	Year         string `bson:"year"`
	TotalRevenue int32  `bson:"total_revenue"`
}

type MonthlyOrderRow struct {
	Month          string  `bson:"month"`
	OrderCount     int64   `bson:"order_count"`
	TotalRevenue   float64 `bson:"total_revenue"`
	TotalItemsSold int64   `bson:"total_items_sold"`
}

type YearlyOrderRow struct {
	Year               string  `bson:"year"`
	OrderCount         int64   `bson:"order_count"`
	TotalRevenue       float64 `bson:"total_revenue"`
	TotalItemsSold     int64   `bson:"total_items_sold"`
	ActiveCashiers     int64   `bson:"active_cashiers"`
	UniqueProductsSold int64   `bson:"unique_products_sold"`
}

type MonthlyTotalRevenueByMerchantRow = MonthlyTotalRevenueRow

type YearlyTotalRevenueByMerchantRow = YearlyTotalRevenueRow

type MonthlyOrderByMerchantRow = MonthlyOrderRow

type YearlyOrderByMerchantRow = YearlyOrderRow

// ---------------------------------------------------------------------------
// Inbox (consumer idempotency)
// ---------------------------------------------------------------------------

type ConsumerInbox struct {
	ConsumerName       string     `bson:"consumer_name"`
	EventKey           string     `bson:"event_key"`
	Topic              string     `bson:"topic"`
	PartitionID        int32      `bson:"partition_id"`
	MessageOffset      int64      `bson:"message_offset"`
	Status             string     `bson:"status"`
	Attempts           int32      `bson:"attempts"`
	LeaseUntil         time.Time  `bson:"lease_until"`
	LastError          string     `bson:"last_error"`
	ProcessedAt        *time.Time `bson:"processed_at,omitempty"`
	ReservationVersion int64      `bson:"reservation_version"`
}

func (ConsumerInbox) CollectionName() string { return "consumer_inbox" }
