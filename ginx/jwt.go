package ginx

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JwtHMAC struct {
	Key      string        `mapstructure:"key"`
	Interval time.Duration `mapstructure:"interval"` // dynamic refresh token ttl in cache
	Duration time.Duration `mapstructure:"duration"`
	Method   uint          `mapstructure:"method"`

	key     []byte
	issuer  string
	method  *jwt.SigningMethodHMAC // SigningMethodHS{256,384,512}
	options []jwt.ParserOption
}

func NewJwtHMAC(vp *viper.Viper, issuer string) (jh *JwtHMAC, err error) {
	/*
		var (
			key string
		)

			if key = vp.GetString("key"); len(key) == 0 {
				return nil, fmt.Errorf("invalid key")
			}

			jh = &JwtHMAC{issuer: issuer, key: []byte(key), duration: vp.GetDuration("duration")}
	*/

	jh = new(JwtHMAC)
	if err = vp.Unmarshal(jh); err != nil {
		return nil, err
	}

	if len(jh.Key) == 0 {
		return nil, fmt.Errorf("invalid key")
	}
	jh.key = []byte(jh.Key)

	if jh.Interval < 0 {
		return nil, fmt.Errorf("invalid interval")
	}

	if jh.Duration < 0 {
		return nil, fmt.Errorf("invalid duration")
	}

	switch jh.Method {
	case 256:
		jh.method = jwt.SigningMethodHS256
	case 384:
		jh.method = jwt.SigningMethodHS384
	case 512:
		jh.method = jwt.SigningMethodHS512
	default:
		return nil, fmt.Errorf("invalid method")
	}
	jh.issuer = issuer

	jh.options = []jwt.ParserOption{
		jwt.WithValidMethods([]string{jh.method.Name}),
		jwt.WithIssuer(jh.issuer),
		jwt.WithIssuedAt(),
	}
	if jh.Duration > 0 {
		jh.options = append(jh.options, jwt.WithExpirationRequired())
	}

	return jh, nil
}

// see jwt.RegisteredClaims
type JwtData struct {
	Issuer    string   `json:"iss"` // required: *app_name
	Subject   string   `json:"sub"` // required: *account_id
	Audience  []string `json:"aud,omitempty"`
	IssuedAt  int64    `json:"iat"` // required:
	ExpiresAt int64    `json:"exp"` // required:
	NotBefore int64    `json:"nbf,omitempty"`
	ID        string   `json:"jti"` // required: *request_id

	Data map[string]string `json:"_data"`
	// TODO: platform
}

// Authorization: Bearer xxxx
// go doc jwt/v5.RegisteredClaims: iss, sub, aud, exp, nbf, iat, jti
func (self *JwtHMAC) Sign(data *JwtData) (signed string, err error) {
	var (
		now    time.Time
		token  *jwt.Token
		claims jwt.MapClaims
	)

	now = time.Now()
	claims = make(jwt.MapClaims, 6)

	data.Issuer = self.issuer
	data.IssuedAt = now.Unix()
	data.ExpiresAt = now.Add(self.Duration).Unix()

	claims["iss"] = self.issuer
	claims["jti"] = data.ID
	claims["sub"] = data.Subject
	claims["iat"] = data.IssuedAt

	if self.Duration > 0 {
		claims["exp"] = data.ExpiresAt
	} else {
		claims["exp"] = 0
	}

	claims["_data"] = data.Data

	token = jwt.NewWithClaims(self.method, claims)

	if signed, err = token.SignedString(self.key); err != nil {
		return "", err
	} else {
		return signed, nil
	}
}

func (self *JwtHMAC) ParsePayload(signed string) (data *JwtData, err error) {
	var (
		bts     []byte
		token   []string
		payload string
	)

	token = strings.SplitN(signed, ".", 3)
	if len(token) != 3 {
		return nil, fmt.Errorf("invalid token")
	}

	payload = token[1]
	if i := len(payload) % 4; i != 0 {
		payload += strings.Repeat("=", 4-i)
	}

	if bts, err = base64.StdEncoding.DecodeString(payload); err != nil {
		return nil, fmt.Errorf("base64: %w", err)
	}

	data = new(JwtData)
	if err = json.Unmarshal(bts, &data); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return data, nil
}

// kind: enum("token_expired", "invalid_token")
func (self *JwtHMAC) Auth(signed string) (data *JwtData, kind string, err error) {
	var (
		ok     bool
		token  *jwt.Token
		claims jwt.MapClaims
	)

	keyfunc := func(token *jwt.Token) (any, error) {
		/*
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
			}
		*/

		return self.key, nil
	}

	if token, err = jwt.Parse(signed, keyfunc, self.options...); err != nil {
		errStr := err.Error()
		err = fmt.Errorf("%w: %s", err, signed)

		// errStr == "token has invalid claims: token is expired")
		if strings.HasSuffix(errStr, "token is expired") {
			return nil, "token_expired", err
		}

		return nil, "invalid_token", err
	}

	if !token.Valid {
		return nil, "invalid_token", err
	}

	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return nil, "invalid_token", err
	}

	data = &JwtData{}
	data.Issuer, _ = claims["iss"].(string)
	data.ID, _ = claims["jti"].(string)
	data.Subject, _ = claims["sub"].(string)

	issuedAt, _ := claims["iat"].(float64)
	data.IssuedAt = int64(issuedAt)

	expiresAt, _ := claims["exp"].(float64)
	data.ExpiresAt = int64(expiresAt)

	mp, _ := claims["_data"].(map[string]any)
	data.Data = make(map[string]string, len(mp))
	for k, v := range mp {
		data.Data[k], _ = v.(string)
	}

	return data, "", nil
}
