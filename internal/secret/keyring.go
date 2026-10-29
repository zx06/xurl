package secret

type KeyringAPI interface {
	Get(service, account string) (string, error)
	Set(service, account, value string) error
	Delete(service, account string) error
}

func defaultKeyring() KeyringAPI {
	return &osKeyring{}
}

type osKeyring struct{}
