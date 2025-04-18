/*
Shadow-Spotter Next Gen Content Discovery
Copyright (C) 2024  Weidsom Nascimento - SNAKE Security

Based on kiterunner from AssetNote

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package proute


import (
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"

	"gitlab.com/snake-security/shadowspotter/pkg/log"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

// RegisterGobs will register globally the gobs needed to serialize and deserialize an API
func RegisterGobs() {
	gob.Register(APIS{})
	gob.Register(KV{})
	for _, v := range AllCrumbs {
		gob.Register(v)
	}
}

// EncodeProtoFile will encode the APIS to the specified filename, overwriting any existing file
func (a APIS) EncodeProtoFile(filename string) error {
	log.Debug().Str("filename", filename).Msg("encoding api to disk")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.FileMode(0666))
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	return a.EncodeProto(f)
}

// EncodeGobFile will encode the APIS to the specified filename, overwriting any existing file
func (a APIS) EncodeGobFile(filename string) error {
	log.Debug().Str("filename", filename).Msg("encoding api to disk")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, os.FileMode(0666))
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	return a.EncodeGob(f)
}

// EncodeProto will encode the APIs into a gob format that can be easily and quickly decoded
func (a APIS) EncodeProto(w io.Writer) error {
	tmp := ProtoAPIS{}
	for _, v := range a {
		tmp.APIs = append(tmp.APIs, v.ProtoAPI())
	}
	data, err := tmp.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}
	_, err = w.Write(data)
	return err
}

// EncodeGob will encode the APIs into a gob format that can be easily and quickly decoded
func (a APIS) EncodeGob(w io.Writer) error {
	RegisterGobs()
	enc := gob.NewEncoder(w)
	return enc.Encode(a)
}

// DecodeAPIGobFile will decode the API from the given file
func DecodeAPIGobFile(filename string) (APIS, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	return DecodeGobAPI(f)
}

// DecodeAPIProtoFile will decode the API from the given file
func DecodeAPIProtoFile(filename string) (APIS, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	return DecodeProtoAPI(f)
}

// DecodeProtoAPI will decode from the given reader, returning the decoded APIs
func DecodeProtoAPI(r io.Reader) (APIS, error) {
	ret := APIS{}
	tmp := ProtoAPIS{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return ret, errors.Wrap(err, "failed to read all from reader")
	}
	err = proto.Unmarshal(data, &tmp)
	if err != nil {
		return ret, errors.Wrap(err, "failed to unmarshal data")
	}

	return tmp.APIS(), nil
}

// DecodeGobAPI will decode from the given reader, returning the decoded APIs
func DecodeGobAPI(r io.Reader) (APIS, error) {
	RegisterGobs()
	ret := make(APIS, 0)
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}

