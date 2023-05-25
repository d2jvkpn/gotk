package impls

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type BasicAuthentication struct {
	Enable bool   `mapstructure:"enable"`
	Method string `mapstructure:"method"`
	Users  []User `mapstructure:"users"`
	users  map[string]string
}

type User struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func NewBasicAuthentication(vp *viper.Viper, field string) (auth *BasicAuthentication, err error) {
	auth = new(BasicAuthentication)
	if err = vp.UnmarshalKey(field, auth); err != nil {
		return nil, err
	}

	if err = auth.Validate(); err != nil {
		return nil, err
	}

	return auth, nil
}

func (auth *BasicAuthentication) Validate() (err error) {
	if auth.Method != "md5" && auth.Method != "bcrypt" {
		return fmt.Errorf("invalid method")
	}

	if len(auth.Users) == 0 {
		return fmt.Errorf("users is unset")
	}

	auth.users = make(map[string]string, len(auth.Users))
	for _, user := range auth.Users {
		if user.Username == "" || user.Password == "" {
			return fmt.Errorf("invalid element exists in users")
		}
		auth.users[user.Username] = user.Password
	}

	return nil
}

func (auth *BasicAuthentication) Handle(w http.ResponseWriter, r *http.Request) (
	user, code string, err error) {
	if !auth.Enable {
		return "", "disabled", nil
	}

	var (
		ok       bool
		key      []byte
		password string
	)

	defer func() {
		if err != nil {
			w.Header().Set("Www-Authenticate", `Basic realm="username:password"`)
			w.WriteHeader(http.StatusUnauthorized)
		}
	}()

	key = []byte(r.Header.Get("Authorization"))
	if !bytes.HasPrefix(key, []byte("Basic ")) {
		return "", "login_required", fmt.Errorf("login required")
	}
	key = key[6:]

	if key, err = base64.StdEncoding.DecodeString(string(key)); err != nil {
		return "", "decode_basic_failed", fmt.Errorf("invalid token")
	}

	u, p, found := bytes.Cut(key, []byte{':'})
	if !found {
		return string(u), "invalid_token", fmt.Errorf("invalid token")
	}

	if auth.Method == "md5" {
		md5sum := fmt.Sprintf("%x", md5.Sum(key))
		if md5sum != auth.users[string(u)] {
			return string(u), "incorrect_username_or_password",
				fmt.Errorf("incorrect username or password")
		}
		return string(u), "md5", nil
	}

	// auth.Method == "bcrypt"
	if password, ok = auth.users[string(u)]; !ok {
		_ = bcrypt.CompareHashAndPassword([]byte(password), p)
		return string(u), "incorrect_username", fmt.Errorf("incorrect username or password")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password), p); err != nil {
		return string(u), "incorrect_password", fmt.Errorf("incorrect username or password")
	}

	r.Header.Del("Authorization")

	return string(u), "bcrypt", nil
}
