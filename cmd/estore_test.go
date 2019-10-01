package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEStore(t *testing.T) {
	assert := assert.New(t)

	estore, err := newEStore("testdata/hosts")
	assert.NoError(err)
	assert.Equal(0, len(estore.List()))

	estore.Add("192.168.1.42", "srvc1.test", "group1")
	assert.Equal(1, len(estore.List()))

	// add duplicate - same IP
	estore.Add("192.168.1.42", "srvc42.test", "group1")
    t.Logf("BAF1: estore.List()=%+v", estore.List())
	assert.Equal(2, len(estore.List()["group1"]))

	// add duplicate - same Name
	estore.Add("192.168.1.142", "srvc1.test", "group1")
	assert.Equal(3, len(estore.List()["group1"]))

	estore.Add("192.168.1.142", "srvc142.test", "group1")
	assert.Equal(4, len(estore.List()["group1"]))

	// drop non existing item
	err = estore.Del("not.there.com", "group1")
	assert.Equal(errIPOrNameNotFound, err)
	assert.Equal(4, len(estore.List()["group1"]))

	// delete by Name
	err = estore.Del("srvc142.test", "group1")
	assert.NoError(err)
	assert.Equal(3, len(estore.List()["group1"]))

	// delete by IP
	err = estore.Del("192.168.1.42", "group1")
	assert.NoError(err)
	assert.Equal(2, len(estore.List()["group1"]))

	err = estore.Commit()
	assert.NoError(err)
	err = estore.Close()
	assert.NoError(err)

	// open is again and check it's empty - this is the default state
	estore, err = newEStore("testdata/hosts")
	assert.NoError(err)
	assert.Equal(0, len(estore.List()))
	err = estore.Close()
	assert.NoError(err)

}

func TestEStoreMx(t *testing.T) {
	assert := assert.New(t)

	estore, err := newEStoreMx("testdata/hosts")
	assert.NoError(err)
	assert.Equal(0, len(estore.List()))

	estore.Add("192.168.1.42", "srvc1.test", "group1")
	assert.Equal(1, len(estore.List()))

	// add duplicate - same IP
	estore.Add("192.168.1.42", "srvc42.test", "group1")
	assert.Equal(2, len(estore.List()["group1"]))

	// add duplicate - same Name
	estore.Add("192.168.1.142", "srvc1.test", "group1")
	assert.Equal(3, len(estore.List()["group1"]))

	estore.Add("192.168.1.142", "srvc142.test", "group1")
	assert.Equal(4, len(estore.List()["group1"]))

	// drop non existing item
	err = estore.Del("not.there.com", "group1")
	assert.Equal(errIPOrNameNotFound, err)
	assert.Equal(4, len(estore.List()["group1"]))

	// delete by Name
	err = estore.Del("srvc142.test", "group1")
	assert.NoError(err)
	assert.Equal(3, len(estore.List()["group1"]))

	// delete by IP
	err = estore.Del("192.168.1.42", "group1")
	assert.NoError(err)
	assert.Equal(2, len(estore.List()["group1"]))

	err = estore.Close()
	assert.NoError(err)

	// open is again and check it's empty - this is the default state
	estore, err = newEStoreMx("testdata/hosts")
	assert.NoError(err)
	assert.Equal(0, len(estore.List()))
	err = estore.Close()
	assert.NoError(err)

}
