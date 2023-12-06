package mock

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal.Repository -o ./internal/app/repository/mock/withdrawal_mock.go -n WithdrawalMock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	mm_withdrawal "github.com/denis-oreshkevich/gopher-mart/internal/app/domain/withdrawal"
	"github.com/gojuno/minimock/v3"
)

// WithdrawalMock implements withdrawal.Repository
type WithdrawalMock struct {
	t minimock.Tester

	funcFindWithdrawalsByUserID          func(ctx context.Context, userID string) (wa1 []mm_withdrawal.Withdrawal, err error)
	inspectFuncFindWithdrawalsByUserID   func(ctx context.Context, userID string)
	afterFindWithdrawalsByUserIDCounter  uint64
	beforeFindWithdrawalsByUserIDCounter uint64
	FindWithdrawalsByUserIDMock          mWithdrawalMockFindWithdrawalsByUserID

	funcRegisterWithdrawal          func(ctx context.Context, withdraw mm_withdrawal.Withdrawal) (err error)
	inspectFuncRegisterWithdrawal   func(ctx context.Context, withdraw mm_withdrawal.Withdrawal)
	afterRegisterWithdrawalCounter  uint64
	beforeRegisterWithdrawalCounter uint64
	RegisterWithdrawalMock          mWithdrawalMockRegisterWithdrawal
}

// NewWithdrawalMock returns a mock for withdrawal.Repository
func NewWithdrawalMock(t minimock.Tester) *WithdrawalMock {
	m := &WithdrawalMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FindWithdrawalsByUserIDMock = mWithdrawalMockFindWithdrawalsByUserID{mock: m}
	m.FindWithdrawalsByUserIDMock.callArgs = []*WithdrawalMockFindWithdrawalsByUserIDParams{}

	m.RegisterWithdrawalMock = mWithdrawalMockRegisterWithdrawal{mock: m}
	m.RegisterWithdrawalMock.callArgs = []*WithdrawalMockRegisterWithdrawalParams{}

	return m
}

type mWithdrawalMockFindWithdrawalsByUserID struct {
	mock               *WithdrawalMock
	defaultExpectation *WithdrawalMockFindWithdrawalsByUserIDExpectation
	expectations       []*WithdrawalMockFindWithdrawalsByUserIDExpectation

	callArgs []*WithdrawalMockFindWithdrawalsByUserIDParams
	mutex    sync.RWMutex
}

// WithdrawalMockFindWithdrawalsByUserIDExpectation specifies expectation struct of the Repository.FindWithdrawalsByUserID
type WithdrawalMockFindWithdrawalsByUserIDExpectation struct {
	mock    *WithdrawalMock
	params  *WithdrawalMockFindWithdrawalsByUserIDParams
	results *WithdrawalMockFindWithdrawalsByUserIDResults
	Counter uint64
}

// WithdrawalMockFindWithdrawalsByUserIDParams contains parameters of the Repository.FindWithdrawalsByUserID
type WithdrawalMockFindWithdrawalsByUserIDParams struct {
	ctx    context.Context
	userID string
}

// WithdrawalMockFindWithdrawalsByUserIDResults contains results of the Repository.FindWithdrawalsByUserID
type WithdrawalMockFindWithdrawalsByUserIDResults struct {
	wa1 []mm_withdrawal.Withdrawal
	err error
}

// Expect sets up expected params for Repository.FindWithdrawalsByUserID
func (mmFindWithdrawalsByUserID *mWithdrawalMockFindWithdrawalsByUserID) Expect(ctx context.Context, userID string) *mWithdrawalMockFindWithdrawalsByUserID {
	if mmFindWithdrawalsByUserID.mock.funcFindWithdrawalsByUserID != nil {
		mmFindWithdrawalsByUserID.mock.t.Fatalf("WithdrawalMock.FindWithdrawalsByUserID mock is already set by Set")
	}

	if mmFindWithdrawalsByUserID.defaultExpectation == nil {
		mmFindWithdrawalsByUserID.defaultExpectation = &WithdrawalMockFindWithdrawalsByUserIDExpectation{}
	}

	mmFindWithdrawalsByUserID.defaultExpectation.params = &WithdrawalMockFindWithdrawalsByUserIDParams{ctx, userID}
	for _, e := range mmFindWithdrawalsByUserID.expectations {
		if minimock.Equal(e.params, mmFindWithdrawalsByUserID.defaultExpectation.params) {
			mmFindWithdrawalsByUserID.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmFindWithdrawalsByUserID.defaultExpectation.params)
		}
	}

	return mmFindWithdrawalsByUserID
}

// Inspect accepts an inspector function that has same arguments as the Repository.FindWithdrawalsByUserID
func (mmFindWithdrawalsByUserID *mWithdrawalMockFindWithdrawalsByUserID) Inspect(f func(ctx context.Context, userID string)) *mWithdrawalMockFindWithdrawalsByUserID {
	if mmFindWithdrawalsByUserID.mock.inspectFuncFindWithdrawalsByUserID != nil {
		mmFindWithdrawalsByUserID.mock.t.Fatalf("Inspect function is already set for WithdrawalMock.FindWithdrawalsByUserID")
	}

	mmFindWithdrawalsByUserID.mock.inspectFuncFindWithdrawalsByUserID = f

	return mmFindWithdrawalsByUserID
}

// Return sets up results that will be returned by Repository.FindWithdrawalsByUserID
func (mmFindWithdrawalsByUserID *mWithdrawalMockFindWithdrawalsByUserID) Return(wa1 []mm_withdrawal.Withdrawal, err error) *WithdrawalMock {
	if mmFindWithdrawalsByUserID.mock.funcFindWithdrawalsByUserID != nil {
		mmFindWithdrawalsByUserID.mock.t.Fatalf("WithdrawalMock.FindWithdrawalsByUserID mock is already set by Set")
	}

	if mmFindWithdrawalsByUserID.defaultExpectation == nil {
		mmFindWithdrawalsByUserID.defaultExpectation = &WithdrawalMockFindWithdrawalsByUserIDExpectation{mock: mmFindWithdrawalsByUserID.mock}
	}
	mmFindWithdrawalsByUserID.defaultExpectation.results = &WithdrawalMockFindWithdrawalsByUserIDResults{wa1, err}
	return mmFindWithdrawalsByUserID.mock
}

// Set uses given function f to mock the Repository.FindWithdrawalsByUserID method
func (mmFindWithdrawalsByUserID *mWithdrawalMockFindWithdrawalsByUserID) Set(f func(ctx context.Context, userID string) (wa1 []mm_withdrawal.Withdrawal, err error)) *WithdrawalMock {
	if mmFindWithdrawalsByUserID.defaultExpectation != nil {
		mmFindWithdrawalsByUserID.mock.t.Fatalf("Default expectation is already set for the Repository.FindWithdrawalsByUserID method")
	}

	if len(mmFindWithdrawalsByUserID.expectations) > 0 {
		mmFindWithdrawalsByUserID.mock.t.Fatalf("Some expectations are already set for the Repository.FindWithdrawalsByUserID method")
	}

	mmFindWithdrawalsByUserID.mock.funcFindWithdrawalsByUserID = f
	return mmFindWithdrawalsByUserID.mock
}

// When sets expectation for the Repository.FindWithdrawalsByUserID which will trigger the result defined by the following
// Then helper
func (mmFindWithdrawalsByUserID *mWithdrawalMockFindWithdrawalsByUserID) When(ctx context.Context, userID string) *WithdrawalMockFindWithdrawalsByUserIDExpectation {
	if mmFindWithdrawalsByUserID.mock.funcFindWithdrawalsByUserID != nil {
		mmFindWithdrawalsByUserID.mock.t.Fatalf("WithdrawalMock.FindWithdrawalsByUserID mock is already set by Set")
	}

	expectation := &WithdrawalMockFindWithdrawalsByUserIDExpectation{
		mock:   mmFindWithdrawalsByUserID.mock,
		params: &WithdrawalMockFindWithdrawalsByUserIDParams{ctx, userID},
	}
	mmFindWithdrawalsByUserID.expectations = append(mmFindWithdrawalsByUserID.expectations, expectation)
	return expectation
}

// Then sets up Repository.FindWithdrawalsByUserID return parameters for the expectation previously defined by the When method
func (e *WithdrawalMockFindWithdrawalsByUserIDExpectation) Then(wa1 []mm_withdrawal.Withdrawal, err error) *WithdrawalMock {
	e.results = &WithdrawalMockFindWithdrawalsByUserIDResults{wa1, err}
	return e.mock
}

// FindWithdrawalsByUserID implements withdrawal.Repository
func (mmFindWithdrawalsByUserID *WithdrawalMock) FindWithdrawalsByUserID(ctx context.Context, userID string) (wa1 []mm_withdrawal.Withdrawal, err error) {
	mm_atomic.AddUint64(&mmFindWithdrawalsByUserID.beforeFindWithdrawalsByUserIDCounter, 1)
	defer mm_atomic.AddUint64(&mmFindWithdrawalsByUserID.afterFindWithdrawalsByUserIDCounter, 1)

	if mmFindWithdrawalsByUserID.inspectFuncFindWithdrawalsByUserID != nil {
		mmFindWithdrawalsByUserID.inspectFuncFindWithdrawalsByUserID(ctx, userID)
	}

	mm_params := &WithdrawalMockFindWithdrawalsByUserIDParams{ctx, userID}

	// Record call args
	mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.mutex.Lock()
	mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.callArgs = append(mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.callArgs, mm_params)
	mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.mutex.Unlock()

	for _, e := range mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.wa1, e.results.err
		}
	}

	if mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.defaultExpectation.Counter, 1)
		mm_want := mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.defaultExpectation.params
		mm_got := WithdrawalMockFindWithdrawalsByUserIDParams{ctx, userID}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmFindWithdrawalsByUserID.t.Errorf("WithdrawalMock.FindWithdrawalsByUserID got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmFindWithdrawalsByUserID.FindWithdrawalsByUserIDMock.defaultExpectation.results
		if mm_results == nil {
			mmFindWithdrawalsByUserID.t.Fatal("No results are set for the WithdrawalMock.FindWithdrawalsByUserID")
		}
		return (*mm_results).wa1, (*mm_results).err
	}
	if mmFindWithdrawalsByUserID.funcFindWithdrawalsByUserID != nil {
		return mmFindWithdrawalsByUserID.funcFindWithdrawalsByUserID(ctx, userID)
	}
	mmFindWithdrawalsByUserID.t.Fatalf("Unexpected call to WithdrawalMock.FindWithdrawalsByUserID. %v %v", ctx, userID)
	return
}

// FindWithdrawalsByUserIDAfterCounter returns a count of finished WithdrawalMock.FindWithdrawalsByUserID invocations
func (mmFindWithdrawalsByUserID *WithdrawalMock) FindWithdrawalsByUserIDAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmFindWithdrawalsByUserID.afterFindWithdrawalsByUserIDCounter)
}

// FindWithdrawalsByUserIDBeforeCounter returns a count of WithdrawalMock.FindWithdrawalsByUserID invocations
func (mmFindWithdrawalsByUserID *WithdrawalMock) FindWithdrawalsByUserIDBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmFindWithdrawalsByUserID.beforeFindWithdrawalsByUserIDCounter)
}

// Calls returns a list of arguments used in each call to WithdrawalMock.FindWithdrawalsByUserID.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmFindWithdrawalsByUserID *mWithdrawalMockFindWithdrawalsByUserID) Calls() []*WithdrawalMockFindWithdrawalsByUserIDParams {
	mmFindWithdrawalsByUserID.mutex.RLock()

	argCopy := make([]*WithdrawalMockFindWithdrawalsByUserIDParams, len(mmFindWithdrawalsByUserID.callArgs))
	copy(argCopy, mmFindWithdrawalsByUserID.callArgs)

	mmFindWithdrawalsByUserID.mutex.RUnlock()

	return argCopy
}

// MinimockFindWithdrawalsByUserIDDone returns true if the count of the FindWithdrawalsByUserID invocations corresponds
// the number of defined expectations
func (m *WithdrawalMock) MinimockFindWithdrawalsByUserIDDone() bool {
	for _, e := range m.FindWithdrawalsByUserIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FindWithdrawalsByUserIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterFindWithdrawalsByUserIDCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFindWithdrawalsByUserID != nil && mm_atomic.LoadUint64(&m.afterFindWithdrawalsByUserIDCounter) < 1 {
		return false
	}
	return true
}

// MinimockFindWithdrawalsByUserIDInspect logs each unmet expectation
func (m *WithdrawalMock) MinimockFindWithdrawalsByUserIDInspect() {
	for _, e := range m.FindWithdrawalsByUserIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to WithdrawalMock.FindWithdrawalsByUserID with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.FindWithdrawalsByUserIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterFindWithdrawalsByUserIDCounter) < 1 {
		if m.FindWithdrawalsByUserIDMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to WithdrawalMock.FindWithdrawalsByUserID")
		} else {
			m.t.Errorf("Expected call to WithdrawalMock.FindWithdrawalsByUserID with params: %#v", *m.FindWithdrawalsByUserIDMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcFindWithdrawalsByUserID != nil && mm_atomic.LoadUint64(&m.afterFindWithdrawalsByUserIDCounter) < 1 {
		m.t.Error("Expected call to WithdrawalMock.FindWithdrawalsByUserID")
	}
}

type mWithdrawalMockRegisterWithdrawal struct {
	mock               *WithdrawalMock
	defaultExpectation *WithdrawalMockRegisterWithdrawalExpectation
	expectations       []*WithdrawalMockRegisterWithdrawalExpectation

	callArgs []*WithdrawalMockRegisterWithdrawalParams
	mutex    sync.RWMutex
}

// WithdrawalMockRegisterWithdrawalExpectation specifies expectation struct of the Repository.RegisterWithdrawal
type WithdrawalMockRegisterWithdrawalExpectation struct {
	mock    *WithdrawalMock
	params  *WithdrawalMockRegisterWithdrawalParams
	results *WithdrawalMockRegisterWithdrawalResults
	Counter uint64
}

// WithdrawalMockRegisterWithdrawalParams contains parameters of the Repository.RegisterWithdrawal
type WithdrawalMockRegisterWithdrawalParams struct {
	ctx      context.Context
	withdraw mm_withdrawal.Withdrawal
}

// WithdrawalMockRegisterWithdrawalResults contains results of the Repository.RegisterWithdrawal
type WithdrawalMockRegisterWithdrawalResults struct {
	err error
}

// Expect sets up expected params for Repository.RegisterWithdrawal
func (mmRegisterWithdrawal *mWithdrawalMockRegisterWithdrawal) Expect(ctx context.Context, withdraw mm_withdrawal.Withdrawal) *mWithdrawalMockRegisterWithdrawal {
	if mmRegisterWithdrawal.mock.funcRegisterWithdrawal != nil {
		mmRegisterWithdrawal.mock.t.Fatalf("WithdrawalMock.RegisterWithdrawal mock is already set by Set")
	}

	if mmRegisterWithdrawal.defaultExpectation == nil {
		mmRegisterWithdrawal.defaultExpectation = &WithdrawalMockRegisterWithdrawalExpectation{}
	}

	mmRegisterWithdrawal.defaultExpectation.params = &WithdrawalMockRegisterWithdrawalParams{ctx, withdraw}
	for _, e := range mmRegisterWithdrawal.expectations {
		if minimock.Equal(e.params, mmRegisterWithdrawal.defaultExpectation.params) {
			mmRegisterWithdrawal.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRegisterWithdrawal.defaultExpectation.params)
		}
	}

	return mmRegisterWithdrawal
}

// Inspect accepts an inspector function that has same arguments as the Repository.RegisterWithdrawal
func (mmRegisterWithdrawal *mWithdrawalMockRegisterWithdrawal) Inspect(f func(ctx context.Context, withdraw mm_withdrawal.Withdrawal)) *mWithdrawalMockRegisterWithdrawal {
	if mmRegisterWithdrawal.mock.inspectFuncRegisterWithdrawal != nil {
		mmRegisterWithdrawal.mock.t.Fatalf("Inspect function is already set for WithdrawalMock.RegisterWithdrawal")
	}

	mmRegisterWithdrawal.mock.inspectFuncRegisterWithdrawal = f

	return mmRegisterWithdrawal
}

// Return sets up results that will be returned by Repository.RegisterWithdrawal
func (mmRegisterWithdrawal *mWithdrawalMockRegisterWithdrawal) Return(err error) *WithdrawalMock {
	if mmRegisterWithdrawal.mock.funcRegisterWithdrawal != nil {
		mmRegisterWithdrawal.mock.t.Fatalf("WithdrawalMock.RegisterWithdrawal mock is already set by Set")
	}

	if mmRegisterWithdrawal.defaultExpectation == nil {
		mmRegisterWithdrawal.defaultExpectation = &WithdrawalMockRegisterWithdrawalExpectation{mock: mmRegisterWithdrawal.mock}
	}
	mmRegisterWithdrawal.defaultExpectation.results = &WithdrawalMockRegisterWithdrawalResults{err}
	return mmRegisterWithdrawal.mock
}

// Set uses given function f to mock the Repository.RegisterWithdrawal method
func (mmRegisterWithdrawal *mWithdrawalMockRegisterWithdrawal) Set(f func(ctx context.Context, withdraw mm_withdrawal.Withdrawal) (err error)) *WithdrawalMock {
	if mmRegisterWithdrawal.defaultExpectation != nil {
		mmRegisterWithdrawal.mock.t.Fatalf("Default expectation is already set for the Repository.RegisterWithdrawal method")
	}

	if len(mmRegisterWithdrawal.expectations) > 0 {
		mmRegisterWithdrawal.mock.t.Fatalf("Some expectations are already set for the Repository.RegisterWithdrawal method")
	}

	mmRegisterWithdrawal.mock.funcRegisterWithdrawal = f
	return mmRegisterWithdrawal.mock
}

// When sets expectation for the Repository.RegisterWithdrawal which will trigger the result defined by the following
// Then helper
func (mmRegisterWithdrawal *mWithdrawalMockRegisterWithdrawal) When(ctx context.Context, withdraw mm_withdrawal.Withdrawal) *WithdrawalMockRegisterWithdrawalExpectation {
	if mmRegisterWithdrawal.mock.funcRegisterWithdrawal != nil {
		mmRegisterWithdrawal.mock.t.Fatalf("WithdrawalMock.RegisterWithdrawal mock is already set by Set")
	}

	expectation := &WithdrawalMockRegisterWithdrawalExpectation{
		mock:   mmRegisterWithdrawal.mock,
		params: &WithdrawalMockRegisterWithdrawalParams{ctx, withdraw},
	}
	mmRegisterWithdrawal.expectations = append(mmRegisterWithdrawal.expectations, expectation)
	return expectation
}

// Then sets up Repository.RegisterWithdrawal return parameters for the expectation previously defined by the When method
func (e *WithdrawalMockRegisterWithdrawalExpectation) Then(err error) *WithdrawalMock {
	e.results = &WithdrawalMockRegisterWithdrawalResults{err}
	return e.mock
}

// RegisterWithdrawal implements withdrawal.Repository
func (mmRegisterWithdrawal *WithdrawalMock) RegisterWithdrawal(ctx context.Context, withdraw mm_withdrawal.Withdrawal) (err error) {
	mm_atomic.AddUint64(&mmRegisterWithdrawal.beforeRegisterWithdrawalCounter, 1)
	defer mm_atomic.AddUint64(&mmRegisterWithdrawal.afterRegisterWithdrawalCounter, 1)

	if mmRegisterWithdrawal.inspectFuncRegisterWithdrawal != nil {
		mmRegisterWithdrawal.inspectFuncRegisterWithdrawal(ctx, withdraw)
	}

	mm_params := &WithdrawalMockRegisterWithdrawalParams{ctx, withdraw}

	// Record call args
	mmRegisterWithdrawal.RegisterWithdrawalMock.mutex.Lock()
	mmRegisterWithdrawal.RegisterWithdrawalMock.callArgs = append(mmRegisterWithdrawal.RegisterWithdrawalMock.callArgs, mm_params)
	mmRegisterWithdrawal.RegisterWithdrawalMock.mutex.Unlock()

	for _, e := range mmRegisterWithdrawal.RegisterWithdrawalMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRegisterWithdrawal.RegisterWithdrawalMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRegisterWithdrawal.RegisterWithdrawalMock.defaultExpectation.Counter, 1)
		mm_want := mmRegisterWithdrawal.RegisterWithdrawalMock.defaultExpectation.params
		mm_got := WithdrawalMockRegisterWithdrawalParams{ctx, withdraw}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRegisterWithdrawal.t.Errorf("WithdrawalMock.RegisterWithdrawal got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRegisterWithdrawal.RegisterWithdrawalMock.defaultExpectation.results
		if mm_results == nil {
			mmRegisterWithdrawal.t.Fatal("No results are set for the WithdrawalMock.RegisterWithdrawal")
		}
		return (*mm_results).err
	}
	if mmRegisterWithdrawal.funcRegisterWithdrawal != nil {
		return mmRegisterWithdrawal.funcRegisterWithdrawal(ctx, withdraw)
	}
	mmRegisterWithdrawal.t.Fatalf("Unexpected call to WithdrawalMock.RegisterWithdrawal. %v %v", ctx, withdraw)
	return
}

// RegisterWithdrawalAfterCounter returns a count of finished WithdrawalMock.RegisterWithdrawal invocations
func (mmRegisterWithdrawal *WithdrawalMock) RegisterWithdrawalAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRegisterWithdrawal.afterRegisterWithdrawalCounter)
}

// RegisterWithdrawalBeforeCounter returns a count of WithdrawalMock.RegisterWithdrawal invocations
func (mmRegisterWithdrawal *WithdrawalMock) RegisterWithdrawalBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRegisterWithdrawal.beforeRegisterWithdrawalCounter)
}

// Calls returns a list of arguments used in each call to WithdrawalMock.RegisterWithdrawal.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRegisterWithdrawal *mWithdrawalMockRegisterWithdrawal) Calls() []*WithdrawalMockRegisterWithdrawalParams {
	mmRegisterWithdrawal.mutex.RLock()

	argCopy := make([]*WithdrawalMockRegisterWithdrawalParams, len(mmRegisterWithdrawal.callArgs))
	copy(argCopy, mmRegisterWithdrawal.callArgs)

	mmRegisterWithdrawal.mutex.RUnlock()

	return argCopy
}

// MinimockRegisterWithdrawalDone returns true if the count of the RegisterWithdrawal invocations corresponds
// the number of defined expectations
func (m *WithdrawalMock) MinimockRegisterWithdrawalDone() bool {
	for _, e := range m.RegisterWithdrawalMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RegisterWithdrawalMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRegisterWithdrawalCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRegisterWithdrawal != nil && mm_atomic.LoadUint64(&m.afterRegisterWithdrawalCounter) < 1 {
		return false
	}
	return true
}

// MinimockRegisterWithdrawalInspect logs each unmet expectation
func (m *WithdrawalMock) MinimockRegisterWithdrawalInspect() {
	for _, e := range m.RegisterWithdrawalMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to WithdrawalMock.RegisterWithdrawal with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RegisterWithdrawalMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRegisterWithdrawalCounter) < 1 {
		if m.RegisterWithdrawalMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to WithdrawalMock.RegisterWithdrawal")
		} else {
			m.t.Errorf("Expected call to WithdrawalMock.RegisterWithdrawal with params: %#v", *m.RegisterWithdrawalMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRegisterWithdrawal != nil && mm_atomic.LoadUint64(&m.afterRegisterWithdrawalCounter) < 1 {
		m.t.Error("Expected call to WithdrawalMock.RegisterWithdrawal")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *WithdrawalMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockFindWithdrawalsByUserIDInspect()

		m.MinimockRegisterWithdrawalInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *WithdrawalMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *WithdrawalMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockFindWithdrawalsByUserIDDone() &&
		m.MinimockRegisterWithdrawalDone()
}
