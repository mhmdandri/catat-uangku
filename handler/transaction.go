package handler

import (
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

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var transactionRequest transactions.TransactionRequest
	if err := c.ShouldBind(&transactionRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newTransaction, err := h.transactionService.Create(transactionRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": newTransaction,
	})
}

func (h *transactionHandler) GetAllTransactions(c *gin.Context) {
	transactionList, err := h.transactionService.FindAll()
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

func (h *transactionHandler) GetTransactionByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	transactionData, err := h.transactionService.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction tidak ditemukan"})
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

func (h *transactionHandler) GetTransactionsByAccountID(c *gin.Context) {
	accountID, err := uuid.Parse(c.Param("account_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "account_id tidak valid"})
		return
	}
	transactionsData, err := h.transactionService.GetTransactionByAccountID(accountID)
	if err != nil {
		if errors.Is(err, transactions.ErrTransactionsNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transactions tidak ditemukan"})
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
