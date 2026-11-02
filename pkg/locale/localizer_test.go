package locale_test

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/matryer/is"

	"github.com/diwise/frontend-toolkit/pkg/locale"
)

func TestLoadAndGet(t *testing.T) {
	is := is.New(t)

	tempdir := t.TempDir()

	is.NoErr(os.Mkdir(tempdir+"/l10n", 0700))

	for lang, contents := range testData {
		err := os.WriteFile(fmt.Sprintf("%s/l10n/%s.toml", tempdir, lang), []byte(contents), 0600)
		is.NoErr(err)
	}

	languages := []string{"en-US", "en-GB", "sv-SE"}
	bundle := locale.NewLocalizer(tempdir, languages...)

	l10n := bundle.ForAll(languages)
	is.Equal(l10n.Get("cancel"), "Cancel")

	is.Equal(bundle.ForAll([]string{"en-US", "en"}).Get("color"), "Color")
	is.Equal(bundle.For("en-GB").Get("color"), "Colour")

	is.Equal(bundle.For("en-gb;q=0.8, en;q=0.7").Get("cancel"), "Cancel")
}

func TestGetAll(t *testing.T) {
	is := is.New(t)

	tempdir := t.TempDir()

	is.NoErr(os.Mkdir(tempdir+"/l10n", 0700))

	for lang, contents := range testData {
		err := os.WriteFile(fmt.Sprintf("%s/l10n/%s.toml", tempdir, lang), []byte(contents), 0600)
		is.NoErr(err)
	}

	languages := []string{"en-US", "en-GB", "sv-SE"}
	bundle := locale.NewLocalizer(tempdir, languages...)

	expectation := []string{"Color", "Colour", "Färg"}
	for lang, result := range bundle.GetAll("color") {
		is.Equal(expectation[slices.Index(languages, lang)], result)
	}
}

var en_GB string = `
[cancel]
other = "Cancel"

[color]
other = "Colour"
`

var en_US string = `
[cancel]
other = "Cancel"

[color]
other = "Color"
`

var sv_SE string = `
[cancel]
other = "Avbryt"

[color]
other = "Färg"
`

var testData map[string]string = map[string]string{
	"en-us": en_US,
	"en-gb": en_GB,
	"sv-se": sv_SE,
}
