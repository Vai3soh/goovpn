package boltdb

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

var bdb, _ = NewBoltDB(`/tmp/bolt_test.db`)

var key = `config_dir_path`
var value = `~/ovpnconfigs`

func TestBucketIsCreateFalse(t *testing.T) {
	bdb.SetNameBucket(`general_configure`)
	bdb.ReOpen()
	create := bdb.BucketIsCreate()
	require.False(t, create, "")
}

func TestCreateBucket(t *testing.T) {
	bdb.ReOpen()
	err := bdb.CreateBucket(`general_configure`)
	if err != nil {
		log.Fatal(err)
	}
}

func TestBucketIsCreateTrue(t *testing.T) {
	bdb.SetNameBucket(`general_configure`)
	bdb.ReOpen()
	create := bdb.BucketIsCreate()
	require.True(t, create, "")
}

func TestStore(t *testing.T) {
	bdb.SetNameBucket(`general_configure`)
	bdb.ReOpen()
	err := bdb.Store(key, value)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetAllValue(t *testing.T) {
	defer bdb.Close()
	bdb.SetNameBucket(`general_configure`)
	bdb.ReOpen()
	bdb.GetAllValue()
	m := bdb.Message()
	for _, e := range m {
		require.Contains(t, e.AtrId, key)
		require.Contains(t, e.Value, value)
	}
}

func TestGetValueFromBucket(t *testing.T) {
	bdb.SetNameBucket(`general_configure`)
	bdb.ReOpen()
	err := bdb.GetValueFromBucket(`config_dir_path`)
	if err != nil {
		log.Fatal(err)
	}
	m := bdb.Message()
	require.Contains(t, m[0].AtrId, key)
	require.Contains(t, m[0].Value, value)
}

func TestDeleteKey(t *testing.T) {
	bdb.SetNameBucket(`general_configure`)
	bdb.ReOpen()
	err := bdb.DeleteKey(key)
	if err != nil {
		log.Fatal(err)
	}
}
