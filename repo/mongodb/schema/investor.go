package schema

type Investor struct {
	ID                 string `bson:"_id,omitempty"`        // GoogleOauth 用户在 Google 的唯一标识符
	Email              string `bson:"email"`                // GoogleOauth 用户的电子邮件地址
	VerifiedEmail      bool   `bson:"verified_email"`       // GoogleOauth 电子邮件是否经过验证
	Name               string `bson:"name"`                 // GoogleOauth 用户的全名
	GivenName          string `bson:"given_name"`           // GoogleOauth 用户的名
	FamilyName         string `bson:"family_name"`          // GoogleOauth 用户的姓
	Picture            string `bson:"picture"`              // GoogleOauth 用户的头像图片 URL
	Locale             string `bson:"locale"`               // GoogleOauth 用户的区域设置
	LastLoginTimestamp uint64 `bson:"last_login_timestamp"` // 上次登录时间戳(10 digits)
}
