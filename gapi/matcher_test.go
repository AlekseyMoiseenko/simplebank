package gapi

import (
	"fmt"
	"reflect"
	"time"

	"github.com/AlekseyMoiseenko/simplebank/util"
	"go.uber.org/mock/gomock"
)

const maxPasswordStampAge = time.Minute

func verifyHashedPassword(password, hash string) bool {
	if password == "" || hash == "" {
		return false
	}
	return util.CheckPassword(password, hash) == nil
}

type passwordAwareParamsMatcher[T any] struct {
	arg      T
	password string

	getHash              func(actual T) (hash string, valid bool)
	getStampValid        func(actual T) bool
	getStampTime         func(actual T) time.Time
	copyNondeterministic func(actual T, expected *T)
}

func (m passwordAwareParamsMatcher[T]) Matches(x interface{}) bool {
	actual, ok := x.(T)
	if !ok {
		return false
	}

	hash, hashValid := m.getHash(actual)
	stampValid := m.getStampValid(actual)

	switch {
	case m.password != "":
		if !hashValid || !stampValid {
			return false
		}
		if !verifyHashedPassword(m.password, hash) {
			return false
		}
		if time.Since(m.getStampTime(actual)) > maxPasswordStampAge {
			return false
		}
	case hashValid || stampValid:
		return false
	}

	expected := m.arg
	m.copyNondeterministic(actual, &expected)
	return reflect.DeepEqual(expected, actual)
}

func (m passwordAwareParamsMatcher[T]) String() string {
	return fmt.Sprintf("matches arg %+v and password %q", m.arg, m.password)
}

func newPasswordAwareParamsMatcher[T any](
	arg T,
	password string,
	getHash func(T) (string, bool),
	getStampValid func(T) bool,
	getStampTime func(T) time.Time,
	copyNondeterministic func(T, *T),
) gomock.Matcher {
	return passwordAwareParamsMatcher[T]{
		arg:                  arg,
		password:             password,
		getHash:              getHash,
		getStampValid:        getStampValid,
		getStampTime:         getStampTime,
		copyNondeterministic: copyNondeterministic,
	}
}
