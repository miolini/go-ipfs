package testutil

import (
	ds2 "github.com/ipfs/go-ipfs/util/datastore2"
	"github.com/jbenet/go-datastore"
	syncds "github.com/jbenet/go-datastore/sync"
)

func ThreadSafeCloserMapDatastore() ds2.ThreadSafeDatastoreCloser {
	return ds2.CloserWrap(syncds.MutexWrap(datastore.NewMapDatastore()))
}