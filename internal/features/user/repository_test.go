package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
)

const testDatabaseURL = "host=localhost port=5433 user=production password=production dbname=production_api_test sslmode=disable"

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(
		postgres.Open(testDatabaseURL),
		&gorm.Config{},
	)

	require.NoError(t, err)

	// Clean test data before every test.
	err = db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").
		Error

	require.NoError(t, err)

	return db
}

// ============================================================
// Test Helper
// ============================================================

func createTestUser(t *testing.T, db *gorm.DB) *User {
	t.Helper()

	u := &User{
		Username:      "john",
		Email:         "john@example.com",
		PasswordHash:  "hashed-password",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	err := db.Create(u).Error

	require.NoError(t, err)
	require.NotZero(t, u.ID)

	return u
}

// ============================================================
// Create
// ============================================================

func TestRepository_Create_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	u := &User{
		Username:      "john",
		Email:         "john@example.com",
		PasswordHash:  "hashed-password",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	err := repo.Create(
		context.Background(),
		u,
	)

	require.NoError(t, err)
	require.NotZero(t, u.ID)

	var saved User

	err = db.First(
		&saved,
		u.ID,
	).Error

	require.NoError(t, err)

	require.Equal(t, "john", saved.Username)
	require.Equal(t, "john@example.com", saved.Email)
	require.Equal(t, "user", saved.Role)
	require.Equal(t, "active", saved.Status)
	require.True(t, saved.EmailVerified)
}

// ============================================================
// Create - Duplicate Email
// ============================================================

func TestRepository_Create_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	first := &User{
		Username:      "john",
		Email:         "john@example.com",
		PasswordHash:  "hash1",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	err := repo.Create(
		context.Background(),
		first,
	)

	require.NoError(t, err)

	second := &User{
		Username:      "john2",
		Email:         "john@example.com",
		PasswordHash:  "hash2",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	err = repo.Create(
		context.Background(),
		second,
	)

	require.Error(t, err)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserAlreadyExists,
		),
	)
}

// ============================================================
// GetByID - Success
// ============================================================

func TestRepository_GetByID_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	expected := createTestUser(t, db)

	result, err := repo.GetByID(
		context.Background(),
		expected.ID,
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.Username, result.Username)
	require.Equal(t, expected.Email, result.Email)
}

// ============================================================
// GetByID - Not Found
// ============================================================

func TestRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	result, err := repo.GetByID(
		context.Background(),
		999999,
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserNotFound,
		),
	)
}

// ============================================================
// GetByEmail - Success
// ============================================================

func TestRepository_GetByEmail_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	expected := createTestUser(t, db)

	result, err := repo.GetByEmail(
		context.Background(),
		expected.Email,
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.Email, result.Email)
	require.Equal(t, expected.Username, result.Username)
}

// ============================================================
// GetByEmail - Not Found
// ============================================================

func TestRepository_GetByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	result, err := repo.GetByEmail(
		context.Background(),
		"does-not-exist@example.com",
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserNotFound,
		),
	)
}

// ============================================================
// GetByUsername - Success
// ============================================================

func TestRepository_GetByUsername_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	expected := createTestUser(t, db)

	result, err := repo.GetByUsername(
		context.Background(),
		expected.Username,
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(t, expected.ID, result.ID)
	require.Equal(t, expected.Username, result.Username)
	require.Equal(t, expected.Email, result.Email)
}

// ============================================================
// GetByUsername - Not Found
// ============================================================

func TestRepository_GetByUsername_NotFound(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	result, err := repo.GetByUsername(
		context.Background(),
		"does-not-exist",
	)

	require.Error(t, err)
	require.Nil(t, result)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserNotFound,
		),
	)
}

// ============================================================
// List - Success
// ============================================================

func TestRepository_List_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	createTestUser(t, db)

	second := &User{
		Username:      "mary",
		Email:         "mary@example.com",
		PasswordHash:  "hashed-password",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	require.NoError(t, db.Create(second).Error)

	third := &User{
		Username:      "alex",
		Email:         "alex@example.com",
		PasswordHash:  "hashed-password",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	require.NoError(t, db.Create(third).Error)

	users, total, err := repo.List(
		context.Background(),
		0,
		20,
	)

	require.NoError(t, err)

	require.Equal(t, int64(3), total)
	require.Len(t, users, 3)
}

// ============================================================
// List - Pagination
// ============================================================

func TestRepository_List_Pagination(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	createTestUser(t, db)

	second := &User{
		Username:      "mary",
		Email:         "mary@example.com",
		PasswordHash:  "hashed-password",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	require.NoError(t, db.Create(second).Error)

	third := &User{
		Username:      "alex",
		Email:         "alex@example.com",
		PasswordHash:  "hashed-password",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	require.NoError(t, db.Create(third).Error)

	users, total, err := repo.List(
		context.Background(),
		1,
		1,
	)

	require.NoError(t, err)

	require.Equal(t, int64(3), total)
	require.Len(t, users, 1)
}

// ============================================================
// List - Empty
// ============================================================

func TestRepository_List_Empty(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	users, total, err := repo.List(
		context.Background(),
		0,
		20,
	)

	require.NoError(t, err)

	require.Equal(t, int64(0), total)
	require.Empty(t, users)
}

// ============================================================
// Update - Success
// ============================================================

func TestRepository_Update_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	u := createTestUser(t, db)

	u.Username = "updated-john"

	err := repo.Update(
		context.Background(),
		u,
	)

	require.NoError(t, err)

	result, err := repo.GetByID(
		context.Background(),
		u.ID,
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.Equal(
		t,
		"updated-john",
		result.Username,
	)
}

// ============================================================
// Update - Duplicate Email
// ============================================================

func TestRepository_Update_DuplicateEmail(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	john := createTestUser(t, db)

	mary := &User{
		Username:      "mary",
		Email:         "mary@example.com",
		PasswordHash:  "hashed-password",
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	require.NoError(t, db.Create(mary).Error)

	john.Email = mary.Email

	err := repo.Update(
		context.Background(),
		john,
	)

	require.Error(t, err)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserAlreadyExists,
		),
	)
}

// ============================================================
// Delete - Success
// ============================================================

func TestRepository_Delete_Success(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	u := createTestUser(t, db)

	err := repo.Delete(
		context.Background(),
		u.ID,
	)

	require.NoError(t, err)

	_, err = repo.GetByID(
		context.Background(),
		u.ID,
	)

	require.Error(t, err)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserNotFound,
		),
	)
}

// ============================================================
// Delete - Not Found
// ============================================================

func TestRepository_Delete_NotFound(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	err := repo.Delete(
		context.Background(),
		999999,
	)

	require.Error(t, err)

	require.True(
		t,
		apperror.Is(
			err,
			apperror.CodeUserNotFound,
		),
	)
}
