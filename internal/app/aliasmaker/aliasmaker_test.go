package aliasmaker

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/mocks"
	"github.com/Schalure/urlalias/internal/app/models/aliasentity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func Test_createAliasKey(t *testing.T) {

	testCases := []struct {
		name    string
		lastKey string
		want    struct {
			newKey string
			err    error
		}
	}{
		{
			name:    "symple test",
			lastKey: "000000000",
			want: struct {
				newKey string
				err    error
			}{
				newKey: "000000001",
				err:    nil,
			},
		},
		{
			name:    "overload test",
			lastKey: "00000000Z",
			want: struct {
				newKey string
				err    error
			}{
				newKey: "000000010",
				err:    nil,
			},
		},
		{
			name:    "storage full test",
			lastKey: "ZZZZZZZZZ",
			want: struct {
				newKey string
				err    error
			}{
				newKey: "",
				err:    errors.New("it is impossible to generate a new string because the storage is full"),
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {

			aliasKey, err := createAliasKey(test.lastKey)

			assert.Equal(t, aliasKey, test.want.newKey)
			assert.Equal(t, err, test.want.err)
		})
	}
}

func Benchmark_createAliasKey(b *testing.B) {

	for i := 0; i < b.N; i++ {
		createAliasKey("YZZZZZZZZ")
	}
}

func Test_deleteAliasesSimple(t *testing.T) {

	userID := uint64(1)

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	storage := mocks.NewMockStorager(mockController)
	storage.EXPECT().GetLastShortKey().Return("000000001").AnyTimes()
	storage.EXPECT().MarkDeleted(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	
	logger, err := zaplogger.NewZapLogger("")
	require.NoError(t, err)

	service, err := New(storage, logger)
	require.NoError(t, err)

	test := struct {
		userId uint64
		aliasesToDelete []string
		findByShortKey1 *gomock.Call
		findByShortKey2 *gomock.Call
		findByShortKey3 *gomock.Call
		findByShortKey4 *gomock.Call
		want struct {
			AliasesToDelete []string
		}
	}{
		userId: uint64(1),
		findByShortKey1: storage.EXPECT().FindByShortKey(gomock.Any(), "000000001").Return(&aliasentity.AliasURLModel{
			ID: uint64(1),
			UserID: userID,
			ShortKey: "000000001",
		}, nil).AnyTimes(),
		findByShortKey2: storage.EXPECT().FindByShortKey(gomock.Any(), "000000002").Return(&aliasentity.AliasURLModel{
			ID: uint64(2),
			UserID: userID,
			ShortKey: "000000002",
		}, nil).AnyTimes(),
		findByShortKey3: storage.EXPECT().FindByShortKey(gomock.Any(), "000000003").Return(&aliasentity.AliasURLModel{
			ID: uint64(3),
			UserID: userID,
			ShortKey: "000000003",
		}, nil).AnyTimes(),
		findByShortKey4: storage.EXPECT().FindByShortKey(gomock.Any(), "000000004").Return(&aliasentity.AliasURLModel{
			ID: uint64(4),
			UserID: userID,
			ShortKey: "000000004",
		}, nil).AnyTimes(),

		aliasesToDelete: []string{"000000001", "000000002", "000000003", "000000004"},
		want: struct{AliasesToDelete []string}{
			AliasesToDelete: []string{"000000001", "000000002", "000000003", "000000004"},
		},
	}


	result := service.deleteAliases(context.Background(), userID, test.aliasesToDelete)

	sort.Strings(result)
	assert.Equal(t, test.want.AliasesToDelete, result)
}


func Test_deleteAliasesOtherUserID(t *testing.T) {

	userID := uint64(1)

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	storage := mocks.NewMockStorager(mockController)
	storage.EXPECT().GetLastShortKey().Return("000000001").AnyTimes()
	storage.EXPECT().MarkDeleted(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	
	logger, err := zaplogger.NewZapLogger("")
	require.NoError(t, err)

	service, err := New(storage, logger)
	require.NoError(t, err)

	test := struct {
		userId uint64
		aliasesToDelete []string
		findByShortKey1 *gomock.Call
		findByShortKey2 *gomock.Call
		findByShortKey3 *gomock.Call
		findByShortKey4 *gomock.Call
		want struct {
			AliasesToDelete []string
		}
	}{
		userId: uint64(1),
		findByShortKey1: storage.EXPECT().FindByShortKey(gomock.Any(), "000000001").Return(&aliasentity.AliasURLModel{
			ID: uint64(1),
			UserID: userID,
			ShortKey: "000000001",
		}, nil).AnyTimes(),
		findByShortKey2: storage.EXPECT().FindByShortKey(gomock.Any(), "000000002").Return(&aliasentity.AliasURLModel{
			ID: uint64(2),
			UserID: userID,
			ShortKey: "000000002",
		}, nil).AnyTimes(),
		findByShortKey3: storage.EXPECT().FindByShortKey(gomock.Any(), "000000003").Return(&aliasentity.AliasURLModel{
			ID: uint64(3),
			UserID: userID,
			ShortKey: "000000003",
		}, nil).AnyTimes(),
		findByShortKey4: storage.EXPECT().FindByShortKey(gomock.Any(), "000000004").Return(&aliasentity.AliasURLModel{
			ID: uint64(4),
			UserID: uint64(2),
			ShortKey: "000000004",
		}, nil).AnyTimes(),
	
		aliasesToDelete: []string{"000000001", "000000002", "000000003", "000000004"},
		want: struct{AliasesToDelete []string}{
			AliasesToDelete: []string{"000000001", "000000002", "000000003"},
		},
	}
	


	result := service.deleteAliases(context.Background(), userID, test.aliasesToDelete)

	sort.Strings(result)
	assert.Equal(t, test.want.AliasesToDelete, result)
}

func Benchmark_deleteAliasesSimple(b *testing.B) {

	userID := uint64(1)

	mockController := gomock.NewController(b)
	defer mockController.Finish()

	storage := mocks.NewMockStorager(mockController)
	storage.EXPECT().GetLastShortKey().Return("000000001").AnyTimes()
	storage.EXPECT().MarkDeleted(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	
	logger, err := zaplogger.NewZapLogger("")
	require.NoError(b, err)

	service, err := New(storage, logger)
	require.NoError(b, err)

	test := struct {
		aliasesToDelete []string
		findByShortKey1 *gomock.Call
		findByShortKey2 *gomock.Call
		findByShortKey3 *gomock.Call
		findByShortKey4 *gomock.Call
	}{
		findByShortKey1: storage.EXPECT().FindByShortKey(gomock.Any(), "000000001").Return(&aliasentity.AliasURLModel{
			ID: uint64(1),
			UserID: userID,
			ShortKey: "000000001",
		}, nil).AnyTimes(),
		findByShortKey2: storage.EXPECT().FindByShortKey(gomock.Any(), "000000002").Return(&aliasentity.AliasURLModel{
			ID: uint64(2),
			UserID: userID,
			ShortKey: "000000002",
		}, nil).AnyTimes(),
		findByShortKey3: storage.EXPECT().FindByShortKey(gomock.Any(), "000000003").Return(&aliasentity.AliasURLModel{
			ID: uint64(3),
			UserID: userID,
			ShortKey: "000000003",
		}, nil).AnyTimes(),
		findByShortKey4: storage.EXPECT().FindByShortKey(gomock.Any(), "000000004").Return(&aliasentity.AliasURLModel{
			ID: uint64(4),
			UserID: userID,
			ShortKey: "000000004",
		}, nil).AnyTimes(),

		aliasesToDelete: []string{"000000001", "000000002", "000000003", "000000004"},
	}


	for i := 0; i < b.N; i++ {
		service.deleteAliases(context.Background(), userID, test.aliasesToDelete)
	}
}
