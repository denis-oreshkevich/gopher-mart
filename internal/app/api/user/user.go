package user

import (
	"errors"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/api"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/domain/user"
	"github.com/denis-oreshkevich/gopher-mart/internal/app/logger"
	usvc "github.com/denis-oreshkevich/gopher-mart/internal/app/service/user"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
)

const (
	emailRegex = "(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|\"" +
		"^(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21\\x23-\\x5b\\x5d-\\x7f]|\\" +
		"\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])*\")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)" +
		"+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}" +
		"(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:" +
		"(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x21-\\x5a\\x53-\\x7f]|\\" +
		"\\[\\x01-\\x09\\x0b\\x0c\\x0e-\\x7f])+)\\])$"

	passwordRegex = "^[A-Za-z\\d@$!%*?&]{8,16}$"
)

type Controller struct {
	svc *usvc.Service
}

func NewController(svc *usvc.Service) *Controller {
	return &Controller{
		svc: svc,
	}
}

var emailMatcher = regexp.MustCompile(emailRegex)

var passwordMatcher = regexp.MustCompile(passwordRegex)

//go:generate easyjson user.go
//easyjson:json
type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

//easyjson:json
type ValidationErrEntry struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

//easyjson:json
type ValidationResp []ValidationErrEntry

func NewValidationErr(field string, errs []string) ValidationErrEntry {
	return ValidationErrEntry{
		Field:  field,
		Errors: errs,
	}
}

func (a *Controller) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := readAndValidateUser(w, r)
	if !ok {
		logger.Log.Debug("readAndValidateUser is not ok")
		return
	}
	_, err := a.svc.Register(ctx, u.Login, u.Password)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExist) {
			logger.Log.Debug("register user", zap.Error(err))
			w.WriteHeader(http.StatusConflict)
			return
		}
		logger.Log.Error("register user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := a.svc.Login(ctx, u.Login, u.Password)
	if err != nil {
		logger.Log.Error("svc.Login", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(api.AuthorizationHeaderName, token)
	w.WriteHeader(http.StatusAccepted)
}

func (a *Controller) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u, ok := readAndValidateUser(w, r)
	if !ok {
		logger.Log.Debug("readAndValidateUser is not ok")
		return
	}
	token, err := a.svc.Login(ctx, u.Login, u.Password)
	if err != nil {
		logger.Log.Error("svc.Login", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(api.AuthorizationHeaderName, token)
	w.WriteHeader(http.StatusOK)
}

func readAndValidateUser(w http.ResponseWriter, r *http.Request) (AuthUser, bool) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("io.ReadAll", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return AuthUser{}, false
	}
	var u AuthUser
	if err = easyjson.Unmarshal(body, &u); err != nil {
		logger.Log.Error("easyjson.Unmarshal", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return AuthUser{}, false
	}
	valResp, err := validateUser(u)
	if err != nil {
		logger.Log.Debug("u is not valid", zap.Error(err))
		bytes, err := easyjson.Marshal(valResp)
		if err != nil {
			logger.Log.Error("validation easyjson.Marshal", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return AuthUser{}, false
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bytes)
		return AuthUser{}, false
	}
	return u, true
}

var ErrUserValidation = errors.New("user validation error")

func validateUser(user AuthUser) (ValidationResp, error) {
	isEV := emailMatcher.MatchString(user.Login)
	isPV := passwordMatcher.MatchString(user.Password)
	var valErrors ValidationResp = make([]ValidationErrEntry, 0)
	if !isEV {
		validationErr := NewValidationErr("login", []string{"login is not valid"})
		valErrors = append(valErrors, validationErr)
	}
	if !isPV {
		validationErr := NewValidationErr("password", []string{"password is not valid"})
		valErrors = append(valErrors, validationErr)
	}
	if len(valErrors) == 0 {
		return nil, nil
	}

	return valErrors, ErrUserValidation
}
