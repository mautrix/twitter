package payload

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/apache/thrift/lib/go/thrift"
)

// Encode serializes a struct to Thrift binary format using struct tags.
func Encode(v any) ([]byte, error) {
	trans := thrift.NewTMemoryBuffer()
	proto := thrift.NewTBinaryProtocolConf(trans, nil)
	ctx := context.Background()

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("thrift encode: expected struct, got %s", rv.Kind())
	}

	if err := writeStruct(ctx, proto, rv); err != nil {
		return nil, err
	}

	if err := proto.Flush(ctx); err != nil {
		return nil, fmt.Errorf("thrift flush: %w", err)
	}

	return trans.Bytes(), nil
}

// Decode deserializes Thrift binary data into a struct using struct tags.
func Decode(data []byte, v any) error {
	trans := thrift.NewTMemoryBufferLen(len(data))
	if _, err := trans.Write(data); err != nil {
		return fmt.Errorf("thrift decode: write to buffer: %w", err)
	}
	proto := thrift.NewTBinaryProtocolConf(trans, nil)
	ctx := context.Background()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("thrift decode: expected non-nil pointer, got %s", rv.Kind())
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("thrift decode: expected pointer to struct, got pointer to %s", rv.Kind())
	}

	return readStruct(ctx, proto, rv)
}

// fieldInfo holds parsed thrift tag information.
type fieldInfo struct {
	name       string
	id         int16
	fieldIndex int
}

// parseThriftTag parses a thrift:"name,id" tag.
func parseThriftTag(tag string) (name string, id int16, ok bool) {
	if tag == "" || tag == "-" {
		return "", 0, false
	}
	parts := strings.Split(tag, ",")
	if len(parts) < 2 {
		return "", 0, false
	}
	name = parts[0]
	idInt, err := strconv.ParseInt(parts[1], 10, 16)
	if err != nil {
		return "", 0, false
	}
	return name, int16(idInt), true
}

// getStructFields returns sorted field info for a struct type.
func getStructFields(t reflect.Type) []fieldInfo {
	var fields []fieldInfo
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("thrift")
		name, id, ok := parseThriftTag(tag)
		if !ok {
			continue
		}
		fields = append(fields, fieldInfo{name: name, id: id, fieldIndex: i})
	}
	// Sort by field ID (required for binary protocol)
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].id < fields[j].id
	})
	return fields
}

// writeStruct writes a struct to the protocol.
func writeStruct(ctx context.Context, proto thrift.TProtocol, v reflect.Value) error {
	t := v.Type()
	if err := proto.WriteStructBegin(ctx, t.Name()); err != nil {
		return err
	}

	fields := getStructFields(t)
	for _, fi := range fields {
		fv := v.Field(fi.fieldIndex)
		ft := t.Field(fi.fieldIndex)

		// Handle pointer types (optional fields)
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				continue // Skip nil optional fields
			}
			fv = fv.Elem()
		}

		// Skip zero-value slices
		if fv.Kind() == reflect.Slice && fv.IsNil() {
			continue
		}

		thriftType := goTypeToThriftType(fv.Type(), ft.Type)
		if thriftType == thrift.STOP {
			continue // Unknown type, skip
		}

		if err := proto.WriteFieldBegin(ctx, fi.name, thriftType, fi.id); err != nil {
			return err
		}

		if err := writeValue(ctx, proto, fv, ft.Type); err != nil {
			return fmt.Errorf("field %s: %w", fi.name, err)
		}

		if err := proto.WriteFieldEnd(ctx); err != nil {
			return err
		}
	}

	if err := proto.WriteFieldStop(ctx); err != nil {
		return err
	}
	return proto.WriteStructEnd(ctx)
}

// goTypeToThriftType maps Go types to Thrift types.
func goTypeToThriftType(t reflect.Type, originalType reflect.Type) thrift.TType {
	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Bool:
		return thrift.BOOL
	case reflect.Int32:
		return thrift.I32
	case reflect.Int64:
		return thrift.I64
	case reflect.String:
		return thrift.STRING
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return thrift.STRING // []byte is encoded as binary/string
		}
		return thrift.LIST
	case reflect.Struct:
		return thrift.STRUCT
	default:
		return thrift.STOP
	}
}

// writeValue writes a value to the protocol.
func writeValue(ctx context.Context, proto thrift.TProtocol, v reflect.Value, originalType reflect.Type) error {
	// Dereference pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		return proto.WriteBool(ctx, v.Bool())

	case reflect.Int32:
		return proto.WriteI32(ctx, int32(v.Int()))

	case reflect.Int64:
		return proto.WriteI64(ctx, v.Int())

	case reflect.String:
		return proto.WriteString(ctx, v.String())

	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// []byte - write as binary
			return proto.WriteBinary(ctx, v.Bytes())
		}
		// List
		elemType := v.Type().Elem()
		thriftElemType := goTypeToThriftType(elemType, elemType)
		if err := proto.WriteListBegin(ctx, thriftElemType, v.Len()); err != nil {
			return err
		}
		for i := 0; i < v.Len(); i++ {
			if err := writeValue(ctx, proto, v.Index(i), elemType); err != nil {
				return err
			}
		}
		return proto.WriteListEnd(ctx)

	case reflect.Struct:
		return writeStruct(ctx, proto, v)

	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}
}

// readStruct reads a struct from the protocol.
func readStruct(ctx context.Context, proto thrift.TProtocol, v reflect.Value) error {
	if _, err := proto.ReadStructBegin(ctx); err != nil {
		return err
	}

	t := v.Type()
	fieldsByID := make(map[int16]int) // field ID -> field index
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("thrift")
		_, id, ok := parseThriftTag(tag)
		if ok {
			fieldsByID[id] = i
		}
	}

	for {
		_, fieldType, fieldID, err := proto.ReadFieldBegin(ctx)
		if err != nil {
			return err
		}
		if fieldType == thrift.STOP {
			break
		}

		fieldIndex, ok := fieldsByID[fieldID]
		if !ok {
			// Unknown field, skip it
			if err := proto.Skip(ctx, fieldType); err != nil {
				return err
			}
			if err := proto.ReadFieldEnd(ctx); err != nil {
				return err
			}
			continue
		}

		fv := v.Field(fieldIndex)
		ft := t.Field(fieldIndex)

		if err := readValue(ctx, proto, fv, ft.Type, fieldType); err != nil {
			return fmt.Errorf("field %s (id=%d): %w", ft.Name, fieldID, err)
		}

		if err := proto.ReadFieldEnd(ctx); err != nil {
			return err
		}
	}

	return proto.ReadStructEnd(ctx)
}

// readValue reads a value from the protocol into v.
func readValue(ctx context.Context, proto thrift.TProtocol, v reflect.Value, goType reflect.Type, thriftType thrift.TType) error {
	// Handle pointer types
	isPtr := goType.Kind() == reflect.Ptr
	if isPtr {
		if v.IsNil() {
			v.Set(reflect.New(goType.Elem()))
		}
		v = v.Elem()
		goType = goType.Elem()
	}

	switch thriftType {
	case thrift.BOOL:
		val, err := proto.ReadBool(ctx)
		if err != nil {
			return err
		}
		v.SetBool(val)

	case thrift.I32:
		val, err := proto.ReadI32(ctx)
		if err != nil {
			return err
		}
		v.SetInt(int64(val))

	case thrift.I64:
		val, err := proto.ReadI64(ctx)
		if err != nil {
			return err
		}
		v.SetInt(val)

	case thrift.STRING:
		if goType.Kind() == reflect.Slice && goType.Elem().Kind() == reflect.Uint8 {
			// []byte
			val, err := proto.ReadBinary(ctx)
			if err != nil {
				return err
			}
			v.SetBytes(val)
		} else {
			val, err := proto.ReadString(ctx)
			if err != nil {
				return err
			}
			v.SetString(val)
		}

	case thrift.LIST:
		elemType, size, err := proto.ReadListBegin(ctx)
		if err != nil {
			return err
		}
		sliceType := goType
		if sliceType.Kind() != reflect.Slice {
			return fmt.Errorf("expected slice type, got %s", sliceType.Kind())
		}
		slice := reflect.MakeSlice(sliceType, size, size)
		for i := 0; i < size; i++ {
			if err := readValue(ctx, proto, slice.Index(i), sliceType.Elem(), elemType); err != nil {
				return err
			}
		}
		if err := proto.ReadListEnd(ctx); err != nil {
			return err
		}
		v.Set(slice)

	case thrift.STRUCT:
		return readStruct(ctx, proto, v)

	default:
		return proto.Skip(ctx, thriftType)
	}

	return nil
}
