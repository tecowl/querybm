package querybm

import (
	"errors"
	"reflect"
	"testing"
)

type mockScanner struct {
	scanCalled bool
	scanArgs   []any
	scanErr    error
}

func (m *mockScanner) Scan(dest ...any) error {
	m.scanCalled = true
	m.scanArgs = dest
	return m.scanErr
}

type User struct {
	ID    int
	Name  string
	Email string
}

func TestNewFields(t *testing.T) {
	t.Parallel()
	names := []string{"id", "name", "email"}
	mapper := func(s Scanner, u *User) error {
		return s.Scan(&u.ID, &u.Name, &u.Email)
	}

	sc := NewFields(names, mapper)

	if sc == nil {
		t.Fatal("NewFields() returned nil")
	}
	if !reflect.DeepEqual(sc.names, names) {
		t.Errorf("NewFields() names = %v, want %v", sc.names, names)
	}
	if sc.mapper == nil {
		t.Error("NewFields() mapper is nil")
	}
}

func TestFields_Fields(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		names []string
		want  []string
	}{
		{
			name:  "Single field",
			names: []string{"id"},
			want:  []string{"id"},
		},
		{
			name:  "Multiple fields",
			names: []string{"id", "name", "email", "created_at"},
			want:  []string{"id", "name", "email", "created_at"},
		},
		{
			name:  "Empty fields",
			names: []string{},
			want:  []string{},
		},
		{
			name:  "Fields with aliases",
			names: []string{"u.id", "u.name AS user_name", "COUNT(*) AS total"},
			want:  []string{"u.id", "u.name AS user_name", "COUNT(*) AS total"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sc := &Fields[User]{names: tt.names}
			got := sc.Fields()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFields_Mapper(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		mapper   Mapper[User]
		scanner  *mockScanner
		user     *User
		wantErr  bool
		wantUser User
	}{
		{
			name: "Successful scan",
			mapper: func(s Scanner, u *User) error {
				return s.Scan(&u.ID, &u.Name, &u.Email)
			},
			scanner: &mockScanner{},
			user:    &User{},
			wantErr: false,
		},
		{
			name: "Scan returns error",
			mapper: func(s Scanner, u *User) error {
				return s.Scan(&u.ID, &u.Name, &u.Email)
			},
			scanner: &mockScanner{scanErr: errors.New("scan error")}, // nolint:err113
			user:    &User{},
			wantErr: true,
		},
		{
			name: "Custom mapper logic",
			mapper: func(s Scanner, u *User) error {
				var id int
				var name, email string
				if err := s.Scan(&id, &name, &email); err != nil {
					return err
				}
				u.ID = id
				u.Name = "Mr. " + name
				u.Email = email
				return nil
			},
			scanner: &mockScanner{},
			user:    &User{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sc := &Fields[User]{mapper: tt.mapper}
			mapperFunc := sc.Mapper()

			err := mapperFunc(tt.scanner, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mapper() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.scanner.scanCalled {
				t.Error("Mapper() did not call Scanner.Scan()")
			}
		})
	}
}
