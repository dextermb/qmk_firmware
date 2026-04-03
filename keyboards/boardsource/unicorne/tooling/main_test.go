package main

import (
	"reflect"
	"testing"
)

func TestBuildLayerPreservesSplit3x6x3Order(t *testing.T) {
	want := []string{
		"KC_TAB", "KC_Q", "KC_W", "KC_E", "KC_R", "KC_T",
		"KC_Y", "KC_U", "KC_I", "KC_O", "KC_P", "KC_BSPC",
		"KC_LSFT", "KC_A", "KC_S", "KC_D", "KC_F", "KC_G",
		"KC_H", "KC_J", "KC_K", "KC_L", "KC_SCLN", "KC_ESC",
		"KC_LCTL", "KC_Z", "KC_X", "KC_C", "KC_V", "KC_B",
		"KC_N", "KC_M", "KC_COMM", "KC_DOT", "KC_SLSH", "KC_QUOT",
		"KC_LGUI", "KC_SPC", "MO(1)",
		"KC_NO", "KC_ENT", "KC_RALT",
	}

	got := buildLayer(FIRST_HALF[0], SECOND_HALF[0])
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildLayer() mismatch\n got: %#v\nwant: %#v", got, want)
	}
}

func TestBuildLayersRejectsMismatchedHalves(t *testing.T) {
	_, err := buildLayers(FIRST_HALF[:1], SECOND_HALF[:2])
	if err == nil {
		t.Fatal("buildLayers() expected mismatch error")
	}
}

func TestBuildKeymapBuildsFourLayers(t *testing.T) {
	keymap, err := buildKeymap()
	if err != nil {
		t.Fatalf("buildKeymap() error = %v", err)
	}

	if got, want := len(keymap.Layers), 4; got != want {
		t.Fatalf("layer count = %d, want %d", got, want)
	}

	for idx, layer := range keymap.Layers {
		if got, want := len(layer), 42; got != want {
			t.Fatalf("layer %d key count = %d, want %d", idx, got, want)
		}
	}
}
