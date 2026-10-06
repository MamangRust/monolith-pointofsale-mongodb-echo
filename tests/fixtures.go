package tests

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Fixtures exposes the collection-level fixture helpers on top of BaseTestSuite
// so suites that embed it can seed documents without touching raw drivers.
func (s *BaseTestSuite) Fixtures() *TestSuite { return s.ts }

// SeedRoleRow inserts the role when it is missing, mirroring the previous
// "INSERT ... ON CONFLICT DO NOTHING" fixture.
func (ts *TestSuite) SeedRoleRow(ctx context.Context, roleName string) error {
	db, err := ts.MongoDB()
	if err != nil {
		return err
	}

	roleID, err := salemongo.NextID(ctx, db, "roles")
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	_, err = db.Collection("roles").UpdateOne(ctx,
		bson.M{"role_name": roleName},
		bson.M{"$setOnInsert": salemongo.Role{
			RoleID:    roleID,
			RoleName:  roleName,
			CreatedAt: &now,
			UpdatedAt: &now,
		}},
		options.Update().SetUpsert(true),
	)
	return err
}

// SeedUserRow inserts a verified user and returns the generated user id.
func (ts *TestSuite) SeedUserRow(ctx context.Context, firstname, lastname, email, password string) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	userID, err := salemongo.NextID(ctx, db, "users")
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	verified := true
	_, err = db.Collection("users").InsertOne(ctx, salemongo.User{
		UserID:           userID,
		Firstname:        firstname,
		Lastname:         lastname,
		Email:            email,
		Password:         password,
		VerificationCode: "test-verify",
		IsVerified:       &verified,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	})
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// SeedMerchantRow inserts a merchant and returns the generated merchant id.
func (ts *TestSuite) SeedMerchantRow(ctx context.Context, userID int32, name, description, address, contactEmail, contactPhone, status string) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	merchantID, err := salemongo.NextID(ctx, db, "merchants")
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	_, err = db.Collection("merchants").InsertOne(ctx, salemongo.Merchant{
		MerchantID:   merchantID,
		UserID:       userID,
		Name:         name,
		Description:  strPtr(description),
		Address:      strPtr(address),
		ContactEmail: strPtr(contactEmail),
		ContactPhone: strPtr(contactPhone),
		Status:       status,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	})
	if err != nil {
		return 0, err
	}
	return merchantID, nil
}

// SeedCashierRow inserts a cashier and returns the generated cashier id.
func (ts *TestSuite) SeedCashierRow(ctx context.Context, merchantID, userID int32, name string) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	cashierID, err := salemongo.NextID(ctx, db, "cashiers")
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	_, err = db.Collection("cashiers").InsertOne(ctx, salemongo.Cashier{
		CashierID:  cashierID,
		MerchantID: merchantID,
		UserID:     userID,
		Name:       name,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	})
	if err != nil {
		return 0, err
	}
	return cashierID, nil
}

// CashierIDForUser returns the active cashier id for a user/merchant pair.
func (ts *TestSuite) CashierIDForUser(ctx context.Context, userID, merchantID int32) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	var cashier salemongo.Cashier
	err = db.Collection("cashiers").FindOne(ctx, bson.M{
		"user_id":     userID,
		"merchant_id": merchantID,
		"deleted_at":  nil,
	}).Decode(&cashier)
	if err != nil {
		return 0, err
	}
	return cashier.CashierID, nil
}

// SeedCategoryRow inserts a category and returns the generated category id.
func (ts *TestSuite) SeedCategoryRow(ctx context.Context, name, description string) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	categoryID, err := salemongo.NextID(ctx, db, "categories")
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	_, err = db.Collection("categories").InsertOne(ctx, salemongo.Category{
		CategoryID:  categoryID,
		Name:        name,
		Description: strPtr(description),
		CreatedAt:   &now,
		UpdatedAt:   &now,
	})
	if err != nil {
		return 0, err
	}
	return categoryID, nil
}

// SeedProductRow inserts a product and returns the generated product id.
func (ts *TestSuite) SeedProductRow(ctx context.Context, merchantID, categoryID int32, name, description string, price, countInStock int32) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	productID, err := salemongo.NextID(ctx, db, "products")
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	_, err = db.Collection("products").InsertOne(ctx, salemongo.Product{
		ProductID:    productID,
		MerchantID:   merchantID,
		CategoryID:   categoryID,
		Name:         name,
		Description:  strPtr(description),
		Price:        price,
		CountInStock: countInStock,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	})
	if err != nil {
		return 0, err
	}
	return productID, nil
}

// SeedOrderRow inserts an order and returns the generated order id.
func (ts *TestSuite) SeedOrderRow(ctx context.Context, merchantID, cashierID int32, totalPrice int64) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	orderID, err := salemongo.NextID(ctx, db, "orders")
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	_, err = db.Collection("orders").InsertOne(ctx, salemongo.Order{
		OrderID:    orderID,
		MerchantID: merchantID,
		CashierID:  cashierID,
		TotalPrice: totalPrice,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	})
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

// SeedOrderItemRow inserts an order item and returns the generated id.
func (ts *TestSuite) SeedOrderItemRow(ctx context.Context, orderID, productID, quantity, price int32) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	orderItemID, err := salemongo.NextID(ctx, db, "order_items")
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	_, err = db.Collection("order_items").InsertOne(ctx, salemongo.OrderItem{
		OrderItemID: orderItemID,
		OrderID:     orderID,
		ProductID:   productID,
		Quantity:    quantity,
		Price:       int64(price),
		CreatedAt:   &now,
		UpdatedAt:   &now,
	})
	if err != nil {
		return 0, err
	}
	return orderItemID, nil
}

// OrderItemIDForOrder returns an order item id belonging to the given order.
func (ts *TestSuite) OrderItemIDForOrder(ctx context.Context, orderID int32) (int32, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}

	var item salemongo.OrderItem
	if err := db.Collection("order_items").FindOne(ctx, bson.M{"order_id": orderID}).Decode(&item); err != nil {
		return 0, err
	}
	return item.OrderItemID, nil
}

// TrashOrderItemRow soft deletes an order item.
func (ts *TestSuite) TrashOrderItemRow(ctx context.Context, orderItemID int32) error {
	db, err := ts.MongoDB()
	if err != nil {
		return err
	}
	_, err = db.Collection("order_items").UpdateOne(ctx,
		bson.M{"order_item_id": orderItemID},
		bson.M{"$set": bson.M{"deleted_at": time.Now().UTC()}},
	)
	return err
}

// RestoreOrderItemRow clears the soft-delete marker of an order item.
func (ts *TestSuite) RestoreOrderItemRow(ctx context.Context, orderItemID int32) error {
	db, err := ts.MongoDB()
	if err != nil {
		return err
	}
	_, err = db.Collection("order_items").UpdateOne(ctx,
		bson.M{"order_item_id": orderItemID},
		bson.M{"$unset": bson.M{"deleted_at": ""}},
	)
	return err
}

// SetOrderCreatedAt overrides the order timestamp so stats queries pick it up.
func (ts *TestSuite) SetOrderCreatedAt(ctx context.Context, orderID int32, createdAt time.Time) error {
	db, err := ts.MongoDB()
	if err != nil {
		return err
	}
	_, err = db.Collection("orders").UpdateOne(ctx,
		bson.M{"order_id": orderID},
		bson.M{"$set": bson.M{"created_at": createdAt, "updated_at": createdAt}},
	)
	return err
}

// SeedTransactionRow inserts a transaction fixture for the stats suites.
func (ts *TestSuite) SeedTransactionRow(ctx context.Context, orderID, merchantID, cashierID, amount int32, paymentMethod, paymentStatus string, createdAt time.Time) error {
	db, err := ts.MongoDB()
	if err != nil {
		return err
	}

	transactionID, err := salemongo.NextID(ctx, db, "transactions")
	if err != nil {
		return err
	}

	_, err = db.Collection("transactions").InsertOne(ctx, salemongo.Transaction{
		TransactionID: transactionID,
		OrderID:       orderID,
		MerchantID:    merchantID,
		CashierID:     cashierID,
		TotalPrice:    int64(amount),
		Amount:        amount,
		PaymentMethod: paymentMethod,
		PaymentStatus: paymentStatus,
		Status:        paymentStatus,
		CreatedAt:     &createdAt,
		UpdatedAt:     &createdAt,
	})
	if err != nil {
		return err
	}
	return nil
}

// SeedTransactionCreatedAt overrides the transaction timestamp.
func (ts *TestSuite) SeedTransactionCreatedAt(ctx context.Context, paymentStatus string, orderID int32, createdAt time.Time) error {
	db, err := ts.MongoDB()
	if err != nil {
		return err
	}
	_, err = db.Collection("transactions").UpdateOne(ctx,
		bson.M{"order_id": orderID, "payment_status": paymentStatus},
		bson.M{"$set": bson.M{"created_at": createdAt, "updated_at": createdAt}},
	)
	return err
}

// IsUserVerified returns the is_verified flag for a user.
func (ts *TestSuite) IsUserVerified(ctx context.Context, userID int32) (bool, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return false, err
	}

	var user salemongo.User
	if err := db.Collection("users").FindOne(ctx, bson.M{"user_id": userID}).Decode(&user); err != nil {
		return false, err
	}
	if user.IsVerified == nil {
		return false, nil
	}
	return *user.IsVerified, nil
}

// VerifyUserRow marks the user with the given email as verified, mirroring the
// email-verification transition the auth flow performs.
func (ts *TestSuite) VerifyUserRow(ctx context.Context, email string) error {
	db, err := ts.MongoDB()
	if err != nil {
		return err
	}
	_, err = db.Collection("users").UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": bson.M{"is_verified": true, "updated_at": time.Now().UTC()}},
	)
	return err
}

// MarkUserVerifiedByID marks the user with the given id as verified.
func (ts *TestSuite) MarkUserVerifiedByID(ctx context.Context, userID int32) error {
	db, err := ts.MongoDB()
	if err != nil {
		return err
	}
	_, err = db.Collection("users").UpdateOne(ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"is_verified": true, "updated_at": time.Now().UTC()}},
	)
	return err
}

// VerificationCodeForUserID returns the verification code for a user id.
func (ts *TestSuite) VerificationCodeForUserID(ctx context.Context, userID int32) (string, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return "", err
	}

	var user salemongo.User
	if err := db.Collection("users").FindOne(ctx, bson.M{"user_id": userID}).Decode(&user); err != nil {
		return "", err
	}
	if user.VerificationCode == "" {
		return "", fmt.Errorf("verification code for user %d is empty", userID)
	}
	return user.VerificationCode, nil
}

// ResetTokenForUser returns the most recent reset token for a user.
func (ts *TestSuite) ResetTokenForUser(ctx context.Context, userID int32) (string, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return "", err
	}

	var token salemongo.ResetToken
	if err := db.Collection("reset_tokens").FindOne(ctx, bson.M{"user_id": userID}).Decode(&token); err != nil {
		return "", err
	}
	return token.Token, nil
}

// StoredPassword returns the hashed password for a user.
func (ts *TestSuite) StoredPassword(ctx context.Context, userID int32) (string, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return "", err
	}

	var user salemongo.User
	if err := db.Collection("users").FindOne(ctx, bson.M{"user_id": userID}).Decode(&user); err != nil {
		return "", err
	}
	return user.Password, nil
}

// ResetTokenCount returns the number of reset tokens for a user.
func (ts *TestSuite) ResetTokenCount(ctx context.Context, userID int32) (int, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return 0, err
	}
	count, err := db.Collection("reset_tokens").CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// UserVerificationCode returns the verification code stored for an email.
func (ts *TestSuite) UserVerificationCode(ctx context.Context, email string) (string, error) {
	db, err := ts.MongoDB()
	if err != nil {
		return "", err
	}

	var user salemongo.User
	if err := db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return "", err
	}
	if user.VerificationCode == "" {
		return "", fmt.Errorf("verification code for %s is empty", email)
	}
	return user.VerificationCode, nil
}

func strPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
