package config

type Model struct {
	AppName  string        `toml:"app-name"`
	Debug    bool          `toml:"debug"`
	Server   ServerModel   `toml:"server"`
	Database DatabaseModel `toml:"database"`
	Encrypt  EncryptModel  `toml:"encrypt"`
	Cron     Cron          `toml:"cron"`
	Mail     MailModel     `toml:"mail"`
	SMS      SMSModel      `toml:"sms"`
	Log      LogModel      `toml:"log"`
}

type ServerModel struct {
	Domain      string `toml:"domain"`
	BaseUrl     string `toml:"base-url"`
	HttpIp      string `toml:"http-ip"`
	HttpPort    int    `toml:"http-port"`
	WebPath     string `toml:"web-path"`
	EnableHttps bool   `toml:"enable-https"`
	CrtFile     string `toml:"crt-file"`
	KeyFile     string `toml:"key-file"`
}

type DatabaseModel struct {
	File string `toml:"file"`
}

type EncryptModel struct {
	Key string `toml:"key" comment:"16/24/32 byte"`
	Iv  string `toml:"iv" comment:"16 byte"`
}

type Cron struct {
	Offset CronOffset `toml:"offset"`
}

type MailModel struct {
	Enable   bool   `toml:"enable"`
	SmtpHost string `toml:"smtp-host"`
	SmtpPort int    `toml:"smtp-port"`
	From     string `toml:"from"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type SMSModel struct {
	Enable bool   `toml:"enable"`
	URL    string `toml:"url"`
	Token  string `toml:"token"`
}

type LogModel struct {
	EnableFile bool     `toml:"enable-file"`
	File       string   `toml:"file"`
	Mask       []string `toml:"mask" comment:"unknown, debug, trace, info, warning, error, fatal"`
	Flag       []string `toml:"flag" comment:"date, time, long_file, short_file, func, prefix, suffix"`
	Debug      []string `toml:"debug" comment:"database, echo"`
}

type CronOffset struct {
	DailyReminderNotify               CronOffsetSysDailyReminderNotify            `toml:"sys-daily-reminder-notify" comment:"daily notify"`
	ScheduleNotify                    CronOffsetSysScheduleNotify                 `toml:"sys-schedule-notify" comment:"monthly schedule notify"`
	UserCardExpDateNotify             CronOffsetUserCardExpDateNotify             `toml:"user-card-expression-date-notify" comment:"user card expression date notify"`
	UserCardStatementClosingDayNotify CronOffsetUserCardStatementClosingDayNotify `toml:"user-card-statement-closing-day-notify" comment:"user card statement closing day notify"`
	UserCardPaymentDueDayNotify       CronOffsetUserCardPaymentDueDayNotify       `toml:"user-card-payment-due-day-notify" comment:"user card payment due day notify"`
	UserPeriodPayNotify               CronOffsetUserPeriodPayNotify               `toml:"user-period-pay-notify" comment:"user period pay notify"`
}

type CronOffsetSysDailyReminderNotify struct {
	Interval uint `toml:"interval"`
	Hour     uint `toml:"hour"`
	Minute   uint `toml:"minute"`
	Second   uint `toml:"second"`
}

type CronOffsetSysScheduleNotify struct {
	Interval       uint `toml:"interval"`
	DaysOfTheMonth int  `toml:"days-of-the-month"`
	Hour           uint `toml:"hour"`
	Minute         uint `toml:"minute"`
	Second         uint `toml:"second"`
}

type CronOffsetUserCardExpDateNotify struct {
	Interval uint `toml:"interval"`
	DayLimit int  `toml:"day-limit" comment:"day limit for expiring soon user card; day in latest dayLimit"`
	Hour     uint `toml:"hour"`
	Minute   uint `toml:"minute"`
	Second   uint `toml:"second"`
}

type CronOffsetUserCardStatementClosingDayNotify struct {
	Interval uint `toml:"interval"`
	DelayDay int  `toml:"delay-day" comment:"delay day for notify; notifyDay = currentDay - delayDay"`
	Hour     uint `toml:"hour"`
	Minute   uint `toml:"minute"`
	Second   uint `toml:"second"`
}

type CronOffsetUserCardPaymentDueDayNotify struct {
	Interval uint `toml:"interval"`
	DelayDay int  `toml:"delay-day" comment:"delay day for notify; notifyDay = currentDay - delayDay"`
	Hour     uint `toml:"hour"`
	Minute   uint `toml:"minute"`
	Second   uint `toml:"second"`
}

type CronOffsetUserPeriodPayNotify struct {
	Interval uint `toml:"interval"`
	DayLimit int  `toml:"day-limit" comment:"day limit for period pay; day in latest dayLimit"`
	Hour     uint `toml:"hour"`
	Minute   uint `toml:"minute"`
	Second   uint `toml:"second"`
}
