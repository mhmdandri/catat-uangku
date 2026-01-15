package handler

import (
	"catatan-keuangan/modules/accounts"
	"catatan-keuangan/modules/common"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type accountHandler struct {
	accountService accounts.Service
}

func NewAccountHandler(accountService accounts.Service) *accountHandler {
	return &accountHandler{accountService}
}

// CreateAccountHandler godoc
// @Summary Buat akun keuangan
// @Tags Accounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body accounts.AccountRequest true "Data akun"
// @Success 201 {object} AccountDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /accounts [post]
func (h *accountHandler) CreateAccountHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	var accountRequest accounts.AccountRequest
	if err := c.ShouldBindJSON(&accountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harap isi semua form",
		})
		return
	}
	newAccount, err := h.accountService.Create(userID, accountRequest)
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
	accountResponse := accounts.FormatAccountResponse(newAccount)
	c.JSON(http.StatusCreated, gin.H{
		"data": accountResponse,
	})
}

// GetAccountByIDHandler godoc
// @Summary Detail akun
// @Tags Accounts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} AccountDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /accounts/{id} [get]
func (h *accountHandler) GetAccountByIDHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	account, err := h.accountService.FindByID(userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
			return
		}
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal mengambil akun",
		})
		return
	}
	accountResponse := accounts.FormatAccountResponse(account)
	c.JSON(http.StatusOK, gin.H{
		"data": accountResponse,
	})
}

// UpdateAccountHandler godoc
// @Summary Perbarui akun
// @Tags Accounts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Account ID"
// @Param request body accounts.AccountUpdateRequest true "Data akun"
// @Success 200 {object} AccountUpdateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /accounts/{id} [put]
func (h *accountHandler) UpdateAccountHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	var accountRequest accounts.AccountUpdateRequest
	if err := c.ShouldBindJSON(&accountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harap isi semua form",
		})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	updatedAccount, err := h.accountService.Update(userID, id, accountRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
			return
		}
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal memperbarui akun",
		})
		return
	}
	accountResponse := accounts.FormatAccountResponse(updatedAccount)
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil memperbarui akun",
		"data":    accountResponse,
	})
}

// DeleteAccountHandler godoc
// @Summary Hapus akun
// @Tags Accounts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /accounts/{id} [delete]
func (h *accountHandler) DeleteAccountHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	_, err = h.accountService.Delete(userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
			return
		}
		if errors.Is(err, accounts.ErrTransactionExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": accounts.ErrTransactionExists.Error()})
			return
		}
		if errors.Is(err, common.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal menghapus akun",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil menghapus akun",
	})
}

// GetAccountByUserIDHandler godoc
// @Summary Daftar akun milik user
// @Tags Accounts
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} AccountsDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /accounts/user/{id} [get]
func (h *accountHandler) GetAccountByUserIDHandler(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak ditemukan"})
		return
	}
	if id != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
		return
	}
	accountsData, err := h.accountService.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak memiliki akun"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal mengambil akun",
		})
		return
	}
	accountResponses := accounts.FormatAccountListItemResponses(accountsData)
	summary := accounts.BuildAccountSummaryByCurrency(accountsData)
	c.JSON(http.StatusOK, gin.H{"summary": summary, "data": accountResponses})
}
