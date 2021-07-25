package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
	"cos-backend-com/src/libs/sdk/eth"
	"time"
)

type StartUpState int

type CreateStartupInput struct {
	CreateStartupRevisionInput
}

type UpdateStartupInput struct {
	CreateStartupRevisionInput
}

type CreateStartupRevisionInput struct {
	Id              string   `json:"id"`
	Name            string   `json:"name"`
	Mission         *string  `json:"mission"`
	Logo            string   `json:"logo"`
	DescriptionAddr string   `json:"descriptionAddr"`
	CategoryId      flake.ID `json:"categoryId"`
	TxId            string   `json:"txId"`
}
type StartupIdResult struct {
	Id flake.ID `json:"id" db:"id"`
}

type StartUpResult struct {
	Id              flake.ID                      `json:"id" db:"id"`
	Name            string                        `json:"name" db:"name"`
	Mission         *string                       `json:"mission" db:"mission"`
	Logo            string                        `json:"logo" db:"logo"`
	DescriptionAddr string                        `json:"descriptionAddr" db:"description_addr"`
	FollowCount     int                           `json:"followCount" db:"follow_count"`
	CreatedAt       time.Time                     `json:"createdAt" db:"created_at"`
	Category        CategoriesResult              `json:"category" db:"category"`
	Setting         StartupSettingRevisionsResult `json:"settings" db:"settings"`
	Transaction     eth.TransactionsResult        `json:"transaction" db:"transaction"`
}

type ListStartupsInput struct {
	CategoryId flake.ID `param:"categoryId"`
	IsIRO      bool     `param:"isIRo"`
	Keyword    string   `param:"keyword"`
	IsInBlock  *bool    `param:"isInBlock"`
	pagination.ListRequest
}

type ListMeStartupsResult struct {
	pagination.ListResult
	Result []struct {
		Id              flake.ID             `json:"id" db:"id"`
		Name            string               `json:"name" db:"name"`
		Mission         *string              `json:"mission" db:"mission"`
		Logo            string               `json:"logo" db:"logo"`
		DescriptionAddr string               `json:"descriptionAddr" db:"description_addr"`
		Category        CategoriesResult     `json:"category" db:"category"`
		State           eth.TransactionState `json:"state" db:"state"`
		SettingState    eth.TransactionState `json:"settingState" db:"setting_state"`
	} `json:"result"`
}

type ListStartupsResult struct {
	pagination.ListResult
	Result []struct {
		Id              flake.ID         `json:"id" db:"id"`
		Name            string           `json:"name" db:"name"`
		Mission         *string          `json:"mission" db:"mission"`
		Logo            string           `json:"logo" db:"logo"`
		DescriptionAddr string           `json:"descriptionAddr" db:"description_addr"`
		Category        CategoriesResult `json:"category" db:"category"`
		IsIRO           bool             `json:"isIRO" db:"is_iro"`
		BountyCount     int              `json:"bountyCount" db:"bounty_count"`
		FollowCount     int              `json:"followCount" db:"follow_count"`
		CreatedAt       time.Time        `json:"createdAt" db:"created_at"`
	} `json:"result"`
}

type HasFollowedStartupResult struct {
	HasFollowed bool `json:"hasFollowed" db:"has_followed"`
}

type StartupShortResult struct {
	Id          flake.ID `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Logo        string   `json:"logo" db:"logo"`
	TokenSymbol string   `json:"tokenSymbol" db:"token_symbol"`
}
