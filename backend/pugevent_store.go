//go:generate qbg -output pugevent_query.go .

package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mercari.io/datastore"
	"google.golang.org/appengine/log"
)

var _ datastore.PropertyLoadSaver = &PugEvent{}
var _ datastore.KeyLoader = &PugEvent{}

// PugEvent is Event Model
// +qbg
type PugEvent struct {
	Key            datastore.Key `datastore:"-" json:"-"`                  // Key
	OrganizationID string        `json:"organizationId"`                   // 支部Id
	Title          string        `json:"title" datastore:",noindex"`       // イベントタイトル
	Description    string        `json:"description" datastore:",noindex"` // イベント説明
	URL            string        `json:"url"`                              // イベント募集URL
	StartAt        time.Time     `json:"startAt"`                          // 開催日時
	EndAt          time.Time     `json:"endAt"`                            // 終了日時
	CreatedAt      time.Time     `json:"createdAt"`                        // 作成日時
	UpdatedAt      time.Time     `json:"updatedAt"`                        // 更新日時
	SchemaVersion  int           `json:"-"`
}

// LoadKey is Entity Load時にKeyを設定する
func (e *PugEvent) LoadKey(ctx context.Context, k datastore.Key) error {
	e.Key = k

	return nil
}

// Load is Entity Load時に呼ばれる
func (e *PugEvent) Load(ctx context.Context, ps []datastore.Property) error {
	err := datastore.LoadStruct(ctx, e, ps)
	if err != nil {
		return err
	}

	return nil
}

// Save is Entity Save時に呼ばれる
func (e *PugEvent) Save(ctx context.Context) ([]datastore.Property, error) {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	e.UpdatedAt = time.Now()
	e.SchemaVersion = 1

	return datastore.SaveStruct(ctx, e)
}

// PugEventStore is PugEventのDatastoreの操作を司る
type PugEventStore struct {
	DatastoreClient datastore.Client
}

// NewPugEventStore is PugEventStoreを作成
func NewPugEventStore(ctx context.Context) (*PugEventStore, error) {
	ds, err := FromContext(ctx)
	if err != nil {
		log.Errorf(ctx, "failed Datastore New Client: %+v", err)
		return nil, err
	}
	return &PugEventStore{ds}, nil
}

// Kind is PugEventのKindを返す
func (store *PugEventStore) Kind() string {
	return "PugEvent"
}

// NewKey is PugEventの新しいKeyを生成する
func (store *PugEventStore) NewKey(ctx context.Context, ds datastore.Client) datastore.Key {
	return ds.NameKey(store.Kind(), uuid.New().String(), nil)
}

// NameKey is PugEventの指定したNameを利用したKeyを生成する
func (store *PugEventStore) NameKey(ctx context.Context, ds datastore.Client, name string) datastore.Key {
	return ds.NameKey(store.Kind(), name, nil)
}

// Create is PugEventをDatastoreにputする
func (store *PugEventStore) Create(ctx context.Context, e *PugEvent) (*PugEvent, error) {
	ds := store.DatastoreClient

	key := store.NewKey(ctx, ds)
	_, err := ds.Put(ctx, key, e)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed put PugEvent to Datastore. key=%v", key))
	}
	e.Key = key
	return e, nil
}

// CreateIfNewURL is Event URLが重複していなければEventを登録する
// 主にconnpassのAPIを叩いて取得した結果を保存するのに利用する
// 引数のURLが空の時は登録を行う
// URLの重複はQueryで調べているので、Eventual Consistency
func (store *PugEventStore) CreateIfNewURL(ctx context.Context, e *PugEvent) (*PugEvent, error) {
	if len(e.URL) > 0 {
		var es []PugEvent
		b := NewPugEventQueryBuilder()
		b = b.URL.Equal(e.URL)
		b = b.KeysOnly()
		ks, err := b.Query().GetAll(ctx, &es)
		if err != nil {
			return nil, errors.Wrap(err, "failed Query to Datastore")
		}
		if len(ks) > 0 {
			return nil, fmt.Errorf("It is a URL that already exists. key=%s, url=%s", ks[0], e.URL)
		}
	} else {
		log.Infof(ctx, "Event URL is Empty")
	}

	return store.Create(ctx, e)
}

// Get is PugEventをDatastoreからgetする
func (store *PugEventStore) Get(ctx context.Context, key datastore.Key) (*PugEvent, error) {
	ds := store.DatastoreClient

	var e PugEvent
	err := ds.Get(ctx, key, &e)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed get PugEvent from Datastore. key=%v", key))
	}

	return &e, nil
}
