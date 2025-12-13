package handler

import (
	"catatan-keuangan/modules/accounts"
	"catatan-keuangan/modules/categories"
	"catatan-keuangan/modules/groups"
	"catatan-keuangan/modules/invitations"
	"catatan-keuangan/modules/transactions"
	userprofile "catatan-keuangan/modules/user_profile"
	"catatan-keuangan/modules/users"
)

// ErrorResponse represents a simple error envelope.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse is used when only a message is returned.
type MessageResponse struct {
	Message string `json:"message"`
}

// URLResponse is returned by OAuth login to point FE to Google auth URL.
type URLResponse struct {
	URL string `json:"url"`
}

// AccessTokenResponse contains the JWT returned by login/refresh.
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type UserDataResponse struct {
	Data users.UserResponse `json:"data"`
}

type UsersDataResponse struct {
	Data []users.UserResponse `json:"data"`
}

type ProfileDataResponse struct {
	Data userprofile.UserProfileResponse `json:"data"`
}

type AccountDataResponse struct {
	Data accounts.AccountResponse `json:"data"`
}

type AccountsDataResponse struct {
	Data []accounts.AccountResponse `json:"data"`
}

type AccountUpdateResponse struct {
	Message string                   `json:"message"`
	Data    accounts.AccountResponse `json:"data"`
}

type GroupDataResponse struct {
	Data groups.GroupResponse `json:"data"`
}

type GroupsDataResponse struct {
	Data []groups.GroupResponse `json:"data"`
}

type InvitationDataResponse struct {
	Data invitations.Invitation `json:"data"`
}

type CategoryDataResponse struct {
	Data categories.CategoryResponse `json:"data"`
}

type TransactionDataResponse struct {
	Data transactions.Transactions `json:"data"`
}

type TransactionResponseBody struct {
	Data transactions.TransactionResponse `json:"data"`
}

type TransactionListResponse struct {
	Data []transactions.TransactionResponse `json:"data"`
}
