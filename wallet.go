package bn

import (
	"context"

	"github.com/libsv/go-bk/wif"
	imodels "github.com/libsv/go-bn/internal/models"
	"github.com/libsv/go-bn/models"
)

type WalletClient interface {
	AbandonTransaction(ctx context.Context, txID string) error
	AddMultiSigAddress(ctx context.Context, n int, keys ...string) (string, error)
	BackupWallet(ctx context.Context, dest string) error
	DumpPrivateKey(ctx context.Context, address string) (*wif.WIF, error)
	DumpWallet(ctx context.Context, dest string) (*models.DumpWallet, error)
	Account(ctx context.Context, address string) (string, error)
	AccountAddress(ctx context.Context, account string) (string, error)
	AccountAddresses(ctx context.Context, account string) ([]string, error)
	Balance(ctx context.Context, opts *models.OptsBalance) (float64, error)
	UnconfirmedBalance(ctx context.Context) (float64, error)
	NewAddress(ctx context.Context, opts *models.OptsNewAddress) (string, error)
	RawChangeAddress(ctx context.Context) (string, error)
	ReceivedByAddress(ctx context.Context, address string) (float64, error)
	Transaction(ctx context.Context, txID string) (*models.Transaction, error)
	ImportAddress(ctx context.Context, address string, opts *models.OptsImportAddress) error
	WalletInfo(ctx context.Context) (*models.WalletInfo, error)
	// import multi
	ImportPrivateKey(ctx context.Context, w *wif.WIF, opts *models.OptsImportPrivateKey) error
	EncryptWallet(ctx context.Context, passphrase string) error
	WalletPhassphrase(ctx context.Context, passphrase string, timeout int) error
	WalletPhassphraseChange(ctx context.Context, oldPassphrase, newPassphrase string) error
	WalletLock(ctx context.Context) error
}

func (c *client) AbandonTransaction(ctx context.Context, txID string) error {
	return c.rpc.Do(ctx, "abandontransaction", nil, txID)
}

func (c *client) AddMultiSigAddress(ctx context.Context, n int, keys ...string) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "addmultisigaddress", &resp, n, keys)
}

func (c *client) BackupWallet(ctx context.Context, dest string) error {
	return c.rpc.Do(ctx, "backupwallet", nil, dest)
}

func (c *client) DumpPrivateKey(ctx context.Context, address string) (*wif.WIF, error) {
	var resp imodels.InternalDumpPrivateKey
	return resp.WIF, c.rpc.Do(ctx, "dumpprivkey", &resp, address)
}

// TODO: do not cache
func (c *client) DumpWallet(ctx context.Context, dest string) (*models.DumpWallet, error) {
	var resp models.DumpWallet
	return &resp, c.rpc.Do(ctx, "dumpwallet", &resp, dest)
}

func (c *client) Account(ctx context.Context, address string) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getaccount", &resp, address)
}

func (c *client) AccountAddress(ctx context.Context, account string) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getaccountaddress", &resp, account)
}

func (c *client) AccountAddresses(ctx context.Context, account string) ([]string, error) {
	var resp []string
	return resp, c.rpc.Do(ctx, "getaddressesbyaccount", &resp, account)
}

// TODO: do not cache
func (c *client) Balance(ctx context.Context, opts *models.OptsBalance) (float64, error) {
	var resp float64
	return resp, c.rpc.Do(ctx, "getbalance", &resp, c.argsFor(opts)...)
}

// TODO: do not cache
func (c *client) UnconfirmedBalance(ctx context.Context) (float64, error) {
	var resp float64
	return resp, c.rpc.Do(ctx, "getunconfirmedbalance", &resp)
}

// TODO: do not cache
func (c *client) NewAddress(ctx context.Context, opts *models.OptsNewAddress) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getnewaddress", &resp, c.argsFor(opts)...)
}

// TODO: do not cache
func (c *client) RawChangeAddress(ctx context.Context) (string, error) {
	var resp string
	return resp, c.rpc.Do(ctx, "getrawchangeaddress", &resp)
}

// TODO: do not cache
func (c *client) ReceivedByAddress(ctx context.Context, address string) (float64, error) {
	var resp float64
	return resp, c.rpc.Do(ctx, "getreceivedbyaddress", &resp, address)
}

func (c *client) Transaction(ctx context.Context, txID string) (*models.Transaction, error) {
	var resp imodels.InternalTransaction
	return resp.Transaction, c.rpc.Do(ctx, "gettransaction", &resp, txID)
}

func (c *client) ImportAddress(ctx context.Context, address string, opts *models.OptsImportAddress) error {
	return c.rpc.Do(ctx, "importaddress", nil, c.argsFor(opts)...)
}

func (c *client) WalletInfo(ctx context.Context) (*models.WalletInfo, error) {
	var resp models.WalletInfo
	return &resp, c.rpc.Do(ctx, "getwalletinfo", &resp)
}

// TODO: importmulti onward
func (c *client) ImportMulti(ctx context.Context, reqs []models.ImportMultiRequest, opts *models.OptsImportMulti) ([]*models.ImportMulti, error) {

}

// TODO: don't cache. test.
func (c *client) ImportPrivateKey(ctx context.Context, w *wif.WIF, opts *models.OptsImportPrivateKey) error {
	return c.rpc.Do(ctx, "importprivkey", nil, c.argsFor(opts, w.String())...)
}

func (c *client) EncryptWallet(ctx context.Context, passphrase string) error {
	return c.rpc.Do(ctx, "encryptwallet", nil, passphrase)
}

func (c *client) WalletPhassphrase(ctx context.Context, passphrase string, timeout int) error {
	return c.rpc.Do(ctx, "walletpassphrase", nil, passphrase, timeout)
}

func (c *client) WalletPhassphraseChange(ctx context.Context, oldPassphrase, newPassphrase string) error {
	return c.rpc.Do(ctx, "walletpassphrasechange", nil, oldPassphrase, newPassphrase)
}

func (c *client) WalletLock(ctx context.Context) error {
	return c.rpc.Do(ctx, "walletlock", nil)
}
