package encrypt_card

import "wallet/cmd/app/server/global/encrypt"

func EncryptCVV(cvv string) (string, error) {
	return encrypt.Encrypt(cvv)
}

func EncryptPassword(cvv string) (string, error) {
	return encrypt.Encrypt(cvv)
}
