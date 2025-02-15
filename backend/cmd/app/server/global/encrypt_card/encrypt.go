package encrypt_card

import (
	"wallet/cmd/app/server/global/encrypt"
	"wallet/cmd/app/server/model"
)

func EncryptCards(cards []*model.UserCard) (err error) {
	for _, card := range cards {
		card.CVV, err = encrypt.Encrypt(card.CVV)
		if err != nil {
			card.CVV = "???"
		}
		card.Password, err = encrypt.Encrypt(card.Password)
		if err != nil {
			card.Password = "??????"
		}
	}
	return err
}

func DecryptCards(cards []*model.UserCard) (err error) {
	for _, card := range cards {
		card.CVV, err = encrypt.Decrypt(card.CVV)
		if err != nil {
			card.CVV = "???"
		}
		card.Password, err = encrypt.Decrypt(card.Password)
		if err != nil {
			card.Password = "??????"
		}
	}
	return err
}

func EncryptCard(card *model.UserCard) (err error) {
	card.CVV, err = encrypt.Encrypt(card.CVV)
	if err != nil {
		card.CVV = "???"
	}
	card.Password, err = encrypt.Encrypt(card.Password)
	if err != nil {
		card.Password = "??????"
	}
	return err
}

func DecryptCard(card *model.UserCard) (err error) {
	card.CVV, err = encrypt.Decrypt(card.CVV)
	if err != nil {
		card.CVV = "???"
	}
	card.Password, err = encrypt.Decrypt(card.Password)
	if err != nil {
		card.Password = "??????"
	}
	return err
}
