package functions

import "testing"

func Test_GetAvatarURL(t *testing.T) {
	if GetAvatarURL("") != "/assets/image/avatar/defaut.png" {
		t.Errorf("except /assets/image/avatar/defaut.png got %s", GetAvatarURL(""))
	}
}
