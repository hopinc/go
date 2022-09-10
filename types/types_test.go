package types

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var types = []reflect.Type{
	// channels.go
	reflect.TypeOf(ChannelType("")),
	reflect.TypeOf(Channel{}),
	reflect.TypeOf(Stats{}),
	reflect.TypeOf(ChannelToken{}),

	// errors.go
	reflect.TypeOf(BadRequest{}),
	reflect.TypeOf(UnknownServerError{}),

	// ignite.go
	reflect.TypeOf(GatewayType("")),
	reflect.TypeOf(GatewayProtocol("")),
	reflect.TypeOf(DomainState("")),
	reflect.TypeOf(Domain{}),
	reflect.TypeOf(Gateway{}),
	reflect.TypeOf(ContainerStrategy("")),
	reflect.TypeOf(RuntimeType("")),
	reflect.TypeOf(DockerAuth{}),
	reflect.TypeOf(ImageGHInfo{}),
	reflect.TypeOf(Image{}),
	reflect.TypeOf(GPUType("")),
	reflect.TypeOf(VGPU{}),
	reflect.TypeOf(Resources{}),
	reflect.TypeOf(DeploymentConfig{}),
	reflect.TypeOf(Deployment{}),
	reflect.TypeOf(Region("")),
	reflect.TypeOf(ContainerUptime{}),
	reflect.TypeOf(ContainerState("")),
	reflect.TypeOf(Container{}),
	reflect.TypeOf(GatewayCreationOptions{}),
	reflect.TypeOf(LoggingLevel("")),
	reflect.TypeOf(ContainerLog{}),

	// pipe.go
	reflect.TypeOf(IngestProtocol("")),
	reflect.TypeOf(DeliveryProtocol("")),
	reflect.TypeOf(RoomState("")),
	reflect.TypeOf(Room{}),
	reflect.TypeOf(HLSConfig{}),
	reflect.TypeOf(RoomCreationOptions{}),

	// projects.go
	reflect.TypeOf(ProjectTier("")),
	reflect.TypeOf(ProjectType("")),
	reflect.TypeOf(DefaultQuotas{}),
	reflect.TypeOf(QuotaUsage{}),
	reflect.TypeOf(Project{}),
	reflect.TypeOf(ProjectToken{}),
	reflect.TypeOf(ProjectPermission("")),
	reflect.TypeOf(ProjectRole{}),
	reflect.TypeOf(ProjectMember{}),
	reflect.TypeOf(ProjectSecret{}),

	// registry.go
	reflect.TypeOf(ImageDigest{}),
	reflect.TypeOf(ImageManifest{}),

	// users.go
	reflect.TypeOf(User{}),
	reflect.TypeOf(UserMeInfo{}),
	reflect.TypeOf(UserPat{}),
}

type stuffer struct {
	i int
	r *strings.Reader
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Take a reflect type and then stuff it with consistent fake data.
// Returns the type.
func (s *stuffer) stuffType(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Bool:
		return reflect.ValueOf(true)
	case reflect.Int:
		v := reflect.ValueOf(s.i)
		s.i++
		return v
	case reflect.Int8:
		v := reflect.ValueOf(int8(s.i))
		s.i++
		return v
	case reflect.Int16:
		v := reflect.ValueOf(int16(s.i))
		s.i++
		return v
	case reflect.Int32:
		v := reflect.ValueOf(int32(s.i))
		s.i++
		return v
	case reflect.Int64:
		v := reflect.ValueOf(int64(s.i))
		s.i++
		return v
	case reflect.Uint:
		v := reflect.ValueOf(uint(s.i))
		s.i++
		return v
	case reflect.Uint8:
		v := reflect.ValueOf(uint8(s.i))
		s.i++
		return v
	case reflect.Uint16:
		v := reflect.ValueOf(uint16(s.i))
		s.i++
		return v
	case reflect.Uint32:
		v := reflect.ValueOf(uint32(s.i))
		s.i++
		return v
	case reflect.Uint64:
		v := reflect.ValueOf(uint64(s.i))
		s.i++
		return v
	case reflect.Uintptr:
		v := reflect.ValueOf(uintptr(s.i))
		s.i++
		return v
	case reflect.Float32:
		v := reflect.ValueOf(float32(s.i))
		s.i++
		return v
	case reflect.Float64:
		v := reflect.ValueOf(float64(s.i))
		s.i++
		return v
	case reflect.Array:
		v := reflect.New(t).Elem()
		for i := 0; i < t.Len(); i++ {
			v.Index(i).Set(s.stuffType(t.Elem()))
		}
		return v
	case reflect.Map:
		m := reflect.MakeMap(t)
		m.SetMapIndex(s.stuffType(t.Key()), s.stuffType(t.Elem()))
		return m
	case reflect.Pointer:
		r := reflect.New(t.Elem())
		r.Elem().Set(s.stuffType(t.Elem()))
		return r
	case reflect.Slice:
		v := reflect.MakeSlice(t, 1, 1)
		v.Index(0).Set(s.stuffType(t.Elem()))
		return v
	case reflect.String:
		if s.r == nil {
			s.r = strings.NewReader(alphabet)
		}
		b := make([]byte, 3)
		if n, err := s.r.Read(b); n != 3 || err != nil {
			s.r = strings.NewReader(alphabet)
			_, _ = s.r.Read(b)
		}
		s.i++
		y := reflect.New(t)
		y.Elem().SetString(string(b))
		return y.Elem()
	case reflect.Struct:
		v := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			v.Field(i).Set(s.stuffType(t.Field(i).Type))
		}
		return v
	case reflect.Interface:
		return s.stuffType(reflect.TypeOf(""))
	default:
		panic("unhandled type: " + t.String())
	}
}

func TestTypes(t *testing.T) {
	typesUpdate := os.Getenv("TYPES_UPDATE") == "1"
	if typesUpdate {
		// Remove and remake the types directory.
		require.NoError(t, os.RemoveAll("testdata"))
		require.NoError(t, os.Mkdir("testdata", 0777))
	}

	for _, type_ := range types {
		t.Run(type_.String(), func(t *testing.T) {
			// Stuff the type with fake data.
			val := (&stuffer{}).stuffType(type_)

			// Turn the value into JSON.
			j, err := json.MarshalIndent(val.Interface(), "", "\t")
			assert.NoError(t, err)

			// If we are updating types, go ahead and write this.
			if typesUpdate {
				require.NoError(t, os.WriteFile("testdata/"+type_.String()+".json", j, 0666))
			}

			// Read the expected JSON.
			expected, err := os.ReadFile("testdata/" + type_.String() + ".json")
			assert.NoError(t, err)
			assert.JSONEq(t, string(expected), string(j))
		})
	}
}
