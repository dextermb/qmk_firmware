package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	KEYMAP_PATH   = "../keymaps/rev_1/keymap.json"
	DOCUMENTATION = "Generated via boardsource/unicorne/tooling"
)

type Row [6]string

type ShortRow [3]string

type KeyboardHalf struct {
	rowOne   Row
	rowTwo   Row
	rowThree Row
	rowFour  ShortRow
}

type KeymapFile struct {
	Version       int        `json:"version"`
	Notes         string     `json:"notes"`
	Documentation string     `json:"documentation"`
	Keyboard      string     `json:"keyboard"`
	Keymap        string     `json:"keymap"`
	Layout        string     `json:"layout"`
	Layers        [][]string `json:"layers"`
	Author        string     `json:"author"`
}

var FIRST_HALF = [...]KeyboardHalf{
	{
		rowOne:   Row{"KC_TAB", "KC_Q", "KC_W", "KC_E", "KC_R", "KC_T"},
		rowTwo:   Row{"KC_LSFT", "KC_A", "KC_S", "KC_D", "KC_F", "KC_G"},
		rowThree: Row{"KC_LCTL", "KC_Z", "KC_X", "KC_C", "KC_V", "KC_B"},
		rowFour:  ShortRow{"KC_LGUI", "KC_SPC", "MO(1)"},
	},
	{
		rowOne:   Row{"KC_GRV", "KC_1", "KC_2", "KC_3", "KC_4", "KC_5"},
		rowTwo:   Row{"KC_NO", "KC_EXLM", "KC_AT", "KC_HASH", "KC_DLR", "KC_PERC"},
		rowThree: Row{"KC_NO", "KC_TILD", "KC_BACKSLASH", "KC_NO", "KC_NO", "KC_NO"},
		rowFour:  ShortRow{"KC_LOPT", "KC_NO", "KC_NO"},
	},
	{
		rowOne:   Row{"KC_NO", "KC_NO", "KC_UP", "KC_NO", "KC_NO", "KC_NO"},
		rowTwo:   Row{"KC_NO", "KC_LEFT", "KC_DOWN", "KC_RGHT", "KC_NO", "KC_NO"},
		rowThree: Row{"KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO"},
		rowFour:  ShortRow{"KC_NO", "KC_NO", "MO(3)"},
	},
	{
		rowOne:   Row{"KC_ESCAPE", "KC_F1", "KC_F2", "KC_F3", "KC_F4", "KC_F5"},
		rowTwo:   Row{"KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO"},
		rowThree: Row{"KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO"},
		rowFour:  ShortRow{"QK_LLCK", "KC_NO", "KC_NO"},
	},
}

var SECOND_HALF = [...]KeyboardHalf{
	{
		rowOne:   Row{"KC_Y", "KC_U", "KC_I", "KC_O", "KC_P", "KC_BSPC"},
		rowTwo:   Row{"KC_H", "KC_J", "KC_K", "KC_L", "KC_SCLN", "KC_ESC"},
		rowThree: Row{"KC_N", "KC_M", "KC_COMM", "KC_DOT", "KC_SLSH", "KC_QUOT"},
		rowFour:  ShortRow{"KC_NO", "KC_ENT", "KC_RALT"},
	},
	{
		rowOne:   Row{"KC_6", "KC_7", "KC_8", "KC_9", "KC_0", "KC_NO"},
		rowTwo:   Row{"KC_CIRC", "KC_AMPR", "KC_ASTR", "KC_LPRN", "KC_RPRN", "KC_NO"},
		rowThree: Row{"KC_MINS", "KC_EQL", "KC_NO", "KC_LBRC", "KC_RBRC", "KC_NO"},
		rowFour:  ShortRow{"MO(2)", "KC_NO", "KC_NO"},
	},
	{
		rowOne:   Row{"RM_VALU", "RM_HUEU", "RM_SATU", "RM_NEXT", "RM_TOGG", "KC_NO"},
		rowTwo:   Row{"RM_VALD", "RM_HUED", "RM_SATD", "RM_PREV", "ANY(CK_TOGG)", "KC_NO"},
		rowThree: Row{"KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO"},
		rowFour:  ShortRow{"QK_BOOT", "EE_CLR", "KC_NO"},
	},
	{
		rowOne:   Row{"KC_NO", "KC_MEDIA_PREV_TRACK", "KC_MEDIA_PLAY_PAUSE", "KC_MEDIA_NEXT_TRACK", "KC_NO", "KC_NO"},
		rowTwo:   Row{"KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO"},
		rowThree: Row{"KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO", "KC_NO"},
		rowFour:  ShortRow{"KC_NO", "KC_NO", "KC_NO"},
	},
}

func buildLayer(left KeyboardHalf, right KeyboardHalf) []string {
	layer := make([]string, 0, 42)
	layer = append(layer, left.rowOne[:]...)
	layer = append(layer, right.rowOne[:]...)
	layer = append(layer, left.rowTwo[:]...)
	layer = append(layer, right.rowTwo[:]...)
	layer = append(layer, left.rowThree[:]...)
	layer = append(layer, right.rowThree[:]...)
	layer = append(layer, left.rowFour[:]...)
	layer = append(layer, right.rowFour[:]...)
	return layer
}

func buildLayers(left []KeyboardHalf, right []KeyboardHalf) ([][]string, error) {
	if len(left) != len(right) {
		return nil, fmt.Errorf("left/right layer count mismatch: %d != %d", len(left), len(right))
	}

	layers := make([][]string, len(left))
	for idx := range left {
		layers[idx] = buildLayer(left[idx], right[idx])
	}

	return layers, nil
}

func buildKeymap() (KeymapFile, error) {
	layers, err := buildLayers(FIRST_HALF[:], SECOND_HALF[:])
	if err != nil {
		return KeymapFile{}, err
	}

	return KeymapFile{
		Version:       1,
		Notes:         "",
		Documentation: DOCUMENTATION,
		Keyboard:      "boardsource/unicorne",
		Keymap:        "boardsource_unicorne_layout_split_3x6_3_2025-09-19",
		Layout:        "LAYOUT_split_3x6_3",
		Layers:        layers,
		Author:        "Dexter Marks-Barber <dexter@marks-barber.co.uk>",
	}, nil
}

func writeKeymap(path string, keymap KeymapFile) error {
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	err := encoder.Encode(keymap)
	if err != nil {
		return err
	}

	return os.WriteFile(path, output.Bytes(), 0o644)
}

func main() {
	keymap, err := buildKeymap()
	if err != nil {
		log.Fatal(err)
	}

	if err := writeKeymap(KEYMAP_PATH, keymap); err != nil {
		log.Fatal(err)
	}
}
