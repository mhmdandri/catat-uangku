package handler

import (
	"catatan-keuangan/modules/accounts"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
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
	var accountRequest accounts.AccountRequest
	if err := c.ShouldBindJSON(&accountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harap isi semua form",
		})
		return
	}
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
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
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
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
