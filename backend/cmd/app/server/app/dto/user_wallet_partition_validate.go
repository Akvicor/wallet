package dto

import (
	"errors"
	"strings"
	"wallet/cmd/app/server/common/types/wallet_partition"
)

func (b *UserWalletPartitionCreate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入名称")
	}
	b.Description = strings.TrimSpace(b.Description)
	if !b.Average.Valid() {
		b.Average = wallet_partition.TypeAverageNormal
	}
	return nil
}

func (b *UserWalletPartitionUpdate) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	if len(b.Name) == 0 {
		return errors.New("请输入名称")
	}
	b.Description = strings.TrimSpace(b.Description)
	if !b.Average.Valid() {
		b.Average = wallet_partition.TypeAverageNormal
	}
	return nil
}

func (b *UserWalletPartitionUpdateSequence) Validate() error {
	return nil
}

func (b *UserWalletPartitionDisable) Validate() error {
	return nil
}

func (b *UserWalletPartitionEnable) Validate() error {
	return nil
}
