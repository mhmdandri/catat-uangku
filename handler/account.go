package handler

import (
	"catatan-keuangan/modules/accounts"
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

func (h *accountHandler) CreateAccountHandler(c *gin.Context) {
	var accountRequest accounts.AccountRequest
	if err := c.ShouldBindJSON(&accountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harap isi semua form",
		})
		return
	}
	newAccount, err := h.accountService.Create(accountRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	accountResponse := accounts.FormatAccountResponse(newAccount)
	c.JSON(http.StatusCreated, gin.H{
		"data": accountResponse,
	})
}

func (h *accountHandler) GetAccountByIDHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	account, err := h.accountService.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
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

func (h *accountHandler) UpdateAccountHandler(c *gin.Context) {
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
	updatedAccount, err := h.accountService.Update(id, accountRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal memperbarui akun",
		})
		return
	}
	accountResponse := accounts.FormatAccountResponse(updatedAccount)
	c.JSON(http.StatusOK, gin.H{
		"data": accountResponse,
	})
}

func (h *accountHandler) DeleteAccountHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	_, err = h.accountService.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Akun tidak ditemukan"})
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

func (h *accountHandler) GetAccountByUserIDHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak ditemukan"})
		return
	}
	accountsData, err := h.accountService.FindByUserID(id)
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
	accountResponses := accounts.FormatAccountResponses(accountsData)
	c.JSON(http.StatusOK, gin.H{"data": accountResponses})
}
