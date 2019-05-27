package youtube

import (
	"reflect"
	"testing"
)

const (
	playerSourceMock = `
	someCodeBefore();
	vv=function(a){a=a.split("");uv.Dr(a,35);uv.CI(a,3);uv.fC(a,26);uv.Dr(a,37);uv.Dr(a,1);uv.Dr(a,65);return a.join("")};
g.wv=function(a){this.o=a;this.A=this.i=this.u="";this.g={};this.l=""};
	var uv={Dr:function(a,b){var c=a[0];a[0]=a[b%a.length];a[b%a.length]=c},
	fC:function(a){a.reverse()},
	CI:function(a,b){a.splice(0,b)}};
	someMoreCodeLiterally`
)

var actionsMock []string = []string{"uv.Dr(a,35)", "uv.CI(a,3)", "uv.fC(a,26)", "uv.Dr(a,37)", "uv.Dr(a,1)", "uv.Dr(a,65)"}
var chipherMock chipherObject = chipherObject{swapFunc: "Dr", reverseFunc: "fC", spliceFunc: "CI"}
var testSignature = "3DAABD3E9044B62C29B63891368EB1C66A681E55F3.1A21A682E6CFF833C470F76434CE0964F4A922FC"
var testSigDecrypted = "92229A4F4690EC43467F074C338FFC6E286A1CA1.3F55E136A66C1BE86319836BF2C26B4409E3DBA"

func Test_getDechipherActions(t *testing.T) {
	type args struct {
		function string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "Assert handels general case",
			args: args{
				function: `function(a){a=a.split("");someTestText;EvenMoreTestText;and.Some_More;return a.join("");}`,
			},
			want:    []string{"someTestText", "EvenMoreTestText", "and.Some_More"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDechipherActions(tt.args.function)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDechipherActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDechipherActions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createDechipherFromChipherAndActions(t *testing.T) {
	t.Run("Assert handels general case", func(t *testing.T) {
		got := createDechipherFromChipherAndActions(chipherMock, actionsMock)
		sig := got.decryptSignature(testSignature)
		if sig != testSigDecrypted {
			t.Errorf("Expected decrypted signature to be %v, got %v", testSigDecrypted, sig)
		}
	})
}

func Test_createReverseAction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Test reverse action"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createReverseAction()
			reversed := got([]rune{'a', 'b', 'c'})
			if !reflect.DeepEqual(reversed, []rune{'c', 'b', 'a'}) {
				t.Errorf("getDechipherActions() = %v, want %v", reversed, []rune{'c', 'b', 'a'})
			}
		})
	}
}

func Test_createSwapAction(t *testing.T) {
	type args struct {
		offset int
		s      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Swap test", args: args{offset: 3, s: "abcdefghijk"}, want: "dbcaefghijk"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSwapAction(tt.args.offset); !reflect.DeepEqual(string(got([]rune(tt.args.s))), tt.want) {
				t.Errorf("createSwapAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createSpliceAction(t *testing.T) {
	type args struct {
		offset int
		s      string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Splice test", args: args{offset: 3, s: "abcdefghijk"}, want: "defghijk"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSpliceAction(tt.args.offset); !reflect.DeepEqual(string(got([]rune(tt.args.s))), tt.want) {
				t.Errorf("createSwapAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createDechipherFromPlayreCode(t *testing.T) {
	t.Run("Assert handels general case", func(t *testing.T) {
		got,err := createDechipherFromPlayreCode(playerSourceMock)
		if err != nil {
			t.Error(err)
			return
		}
		sig := got.decryptSignature(testSignature)
		if sig != testSigDecrypted {
			t.Errorf("Expected decrypted signature to be %v, got %v", testSigDecrypted, sig)
		}
	})
}
