package handler

import (
	"catatan-keuangan/modules/common"
	"catatan-keuangan/modules/transactions"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionHandler struct {
	transactionService transactions.Service
}

func NewTransactionHandler(transactionService transactions.Service) *transactionHandler {
	return &transactionHandler{transactionService}
}

// CreateTransaction godoc
// @Summary Buat transaksi
// @Tags Transactions
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param group_id formData string false "Group ID (isi saat scope=group)"
// @Param category_id formData string true "Category ID"
// @Param account_id formData string true "Account ID"
// @Param type formData string true "income atau expense"
// @Param total_amount formData number true "Nominal transaksi"
// @Param scope formData string true "personal atau group"
// @Param description formData string false "Catatan"
// @Param attachments formData file false "Lampiran transaksi (boleh multiple)"
// @Success 201 {object} TransactionDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /transactions [post]
func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	var transactionRequest transactions.TransactionRequest
	if err := c.ShouldBind(&transactionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newTransaction, err := h.transactionService.Create(userID, transactionRequest)
	if err != nil {
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": newTransaction,
	})
}

// GetAllTransactions godoc
// @Summary Daftar transaksi
// @Tags Transactions
// @Security BearerAuth
// @Produce json
// @Success 200 {object} TransactionListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /transactions [get]
func (h *transactionHandler) GetAllTransactions(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	transactionList, err := h.transactionService.FindAll(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	tResponse := transactions.FormatTransactionResponses(transactionList)
	c.JSON(http.StatusOK, gin.H{
		"data": tResponse,
	})
}

// GetTransactionByID godoc
// @Summary Detail transaksi
// @Tags Transactions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} TransactionResponseBody
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /transactions/{id} [get]
func (h *transactionHandler) GetTransactionByID(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	transactionData, err := h.transactionService.FindByID(userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction tidak ditemukan"})
			return
		}
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal mengambil data transaction",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": transactions.FormatTransactionResponse(transactionData),
	})
}

// GetTransactionsByAccountID godoc
// @Summary Transaksi berdasarkan akun
// @Tags Transactions
// @Security BearerAuth
// @Produce json
// @Param account_id path string true "Account ID"
// @Success 200 {object} TransactionListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /transactions/account/{account_id} [get]
func (h *transactionHandler) GetTransactionsByAccountID(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	accountID, err := uuid.Parse(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id tidak valid"})
		return
	}
	transactionsData, err := h.transactionService.GetTransactionByAccountID(userID, accountID)
	if err != nil {
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		if errors.Is(err, transactions.ErrTransactionsNotFound) {
			c.JSON(http.StatusOK, gin.H{"data": []transactions.TransactionResponse{}})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal mengambil data transactions",
		})
		return
	}
	tResponse := transactions.FormatTransactionResponses(transactionsData)
	c.JSON(http.StatusOK, gin.H{
		"data": tResponse,
	})
}

// GetTransactionByUserID godoc
// @Summary Transaksi berdasarkan user
// @Tags Transactions
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} TransactionListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /transactions/user/{id} [get]
func (h *transactionHandler) GetTransactionByUserID(c *gin.Context) {
	requesterID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id tidak valid"})
		return
	}
	if targetID != requesterID {
		c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
		return
	}
	transactionsData, err := h.transactionService.GetTransactionByUserID(requesterID)
	if err != nil {
		if errors.Is(err, transactions.ErrTransactionsNotFound) {
			c.JSON(http.StatusOK, gin.H{"data": []transactions.TransactionResponse{}})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal mengambil data transactions",
		})
		return
	}
	tResponse := transactions.FormatTransactionResponses(transactionsData)
	c.JSON(http.StatusOK, gin.H{
		"data": tResponse,
	})
}

func (h *transactionHandler) DeleteTransaction(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	err = h.transactionService.Delete(userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction tidak ditemukan"})
			return
		}
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal menghapus data transaction",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction berhasil dihapus",
	})
}

// UpdateTransaction godoc
// @Summary Perbarui transaksi
// @Tags Transactions
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Transaction ID"
// @Param group_id formData string false "Group ID (isi saat scope=group)"
// @Param category_id formData string true "Category ID"
// @Param account_id formData string true "Account ID"
// @Param type formData string true "income atau expense"
// @Param total_amount formData number true "Nominal transaksi"
// @Param scope formData string true "personal atau group"
// @Param description formData string false "Catatan"
// @Param attachments formData file false "Lampiran transaksi (boleh multiple)"
// @Success 200 {object} TransactionDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /transactions/{id} [put]
func (h *transactionHandler) UpdateTransaction(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	var transactionRequest transactions.TransactionRequest
	if err := c.ShouldBind(&transactionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	updatedTransaction, err := h.transactionService.Update(userID, id, transactionRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction tidak ditemukan"})
			return
		}
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": updatedTransaction,
	})
}
