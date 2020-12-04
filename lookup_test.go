package ldapx

import "testing"

func TestGetAttributeFromDN(t *testing.T) {
	type args struct {
		attr string
		dn   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "empty dn", args: args{attr: "uid"}, want: "", wantErr: true},
		{name: "valid dn", args: args{attr: "uid", dn: "uid=test,ou=users,dc=jcu,dc=edu,dc=au"}, want: "test", wantErr: false},
		{name: "attribute is not first RDN", args: args{attr: "ou", dn: "uid=test,ou=users,dc=jcu,dc=edu,dc=au"}, want: "users", wantErr: false},
		{name: "invalid dn", args: args{attr: "uid", dn: "invaliddn"}, want: "", wantErr: true},
		{name: "dn does not contain attribute", args: args{attr: "notfound", dn: "uid=test,ou=users,dc=jcu,dc=edu,dc=au"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAttributeFromDN(tt.args.attr, tt.args.dn)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAttributeFromDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAttributeFromDN() got = %v, want %v", got, tt.want)
			}
		})
	}
}
