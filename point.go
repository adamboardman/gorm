// Based on https://github.com/jinzhu/gorm/issues/142
package gorm

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

//Swapped from struct to array so that reflection can set data type to PostGIS geometry type
type Point [2]float64

func (p *Point) Lng() float64 {
	return p[0]
}

func (p *Point) Lat() float64 {
	return p[1]
}

func (p *Point) Equal(point Point) bool {
	return p[0] == point[0] && p[1] == point[1]
}

func (p *Point) String() string {
	return fmt.Sprintf("SRID=4326;POINT(%g %g)", p[0], p[1])
}

func (p *Point) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}
	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("Invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

func (p Point) Value() (driver.Value, error) {
	return p.String(), nil
}
