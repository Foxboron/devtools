package overlay

import (
	"log"
	"testing"
)

func TestOverlay(t *testing.T) {
	path := "./test_root"
	overlay := newOverlay(path)
	_, err := overlay.Setup()
	if err != nil {
		log.Fatal(err)
	}

	_, err = overlay.AddSnapshot("build")
	if err != nil {
		log.Fatal(err)
	}

	err = overlay.RemoveSnapshot("build")
	if err != nil {
		log.Fatal(err)
	}
}
