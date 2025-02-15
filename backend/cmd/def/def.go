package def

import (
	"fmt"
	"time"
)

const (
	AppName        = "wallet"
	AppUsage       = "Money Manage"
	AppDescription = "Wallet is a money manage service"
)

var (
	Branch    string
	Version   string
	Commit    string
	BuildTime string
)

func AppVersion() string {
	return fmt.Sprintf("%s-%s_%s (%s)", Version, Branch, Commit, BuildTime)
}

func Copyright() string {
	return fmt.Sprintf("Copyright © 2024-%d Akvicor, All Rights Reserved.", time.Now().Year())
}
