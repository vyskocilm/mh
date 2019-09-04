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

	estore.Add("192.168.1.42", "srvc1.test")
	assert.Equal(1, len(estore.List()))

	// add duplicate - same IP
	estore.Add("192.168.1.42", "srvc42.test")
	assert.Equal(2, len(estore.List()))

	// add duplicate - same Name
	estore.Add("192.168.1.142", "srvc1.test")
	assert.Equal(3, len(estore.List()))

	estore.Add("192.168.1.142", "srvc142.test")
	assert.Equal(4, len(estore.List()))

	// drop non existing item
	err = estore.Del("not.there.com")
	assert.Equal(errIPOrNameNotFound, err)
	assert.Equal(4, len(estore.List()))

	// delete by Name
	err = estore.Del("srvc142.test")
	assert.NoError(err)
	assert.Equal(3, len(estore.List()))

	// delete by IP
	err = estore.Del("192.168.1.42")
	assert.NoError(err)
	assert.Equal(2, len(estore.List()))

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

	estore.Add("192.168.1.42", "srvc1.test")
	assert.Equal(1, len(estore.List()))

	// add duplicate - same IP
	estore.Add("192.168.1.42", "srvc42.test")
	assert.Equal(2, len(estore.List()))

	// add duplicate - same Name
	estore.Add("192.168.1.142", "srvc1.test")
	assert.Equal(3, len(estore.List()))

	estore.Add("192.168.1.142", "srvc142.test")
	assert.Equal(4, len(estore.List()))

	// drop non existing item
	err = estore.Del("not.there.com")
	assert.Equal(errIPOrNameNotFound, err)
	assert.Equal(4, len(estore.List()))

	// delete by Name
	err = estore.Del("srvc142.test")
	assert.NoError(err)
	assert.Equal(3, len(estore.List()))

	// delete by IP
	err = estore.Del("192.168.1.42")
	assert.NoError(err)
	assert.Equal(2, len(estore.List()))

	err = estore.Close()
	assert.NoError(err)

	// open is again and check it's empty - this is the default state
	estore, err = newEStoreMx("testdata/hosts")
	assert.NoError(err)
	assert.Equal(0, len(estore.List()))
	err = estore.Close()
	assert.NoError(err)

}
