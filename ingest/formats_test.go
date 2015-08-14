package ingest

import "testing"


func TestFormatValidate(t *testing.T) {
	format := Format{FormatElement{Type:FormatVar}, FormatElement{Type:FormatVar}}
	if format.Validate() == nil {
		t.Fatal("Should not have validated")
	}
}

func TestFormatParse(t *testing.T) {
	format := Format{
		FormatElement{Type:FormatVar, Name:"wally"},
		FormatElement{Type:FormatChar, Char:','},
		FormatElement{Type:FormatVar, Name:"wizard"},
	}
	meta, err := format.Parse([]byte("wat,we"))
	if err != nil {
		t.Fatal(err)
	}
	if meta["wally"] != "wat" {
		t.Log(meta)
		t.Fatal("wally != wat")
	}
	if meta["wizard"] != "we" {
		t.Log(meta)
		t.Fatal("wizard != we")
	}
}

func TestFormatParseFailExtra(t *testing.T) {
	format := Format{
		FormatElement{Type:FormatVar, Name:"wally"},
		FormatElement{Type:FormatChar, Char:','},
		FormatElement{Type:FormatVar, Name:"wizard"},
		FormatElement{Type:FormatChar, Char:'.'},
	}
	_, err := format.Parse([]byte("wat,we.werwe"))
	if err == nil {
		t.Fatal("Parse error should have been thrown")
	}
}

func TestFormatParseFailNotEnough(t *testing.T) {
	format := Format{
		FormatElement{Type:FormatVar, Name:"wally"},
		FormatElement{Type:FormatChar, Char:','},
		FormatElement{Type:FormatVar, Name:"wizard"},
		FormatElement{Type:FormatChar, Char:'.'},
	}
	_, err := format.Parse([]byte("wat,we"))
	if err == nil {
		t.Fatal("Parse error should have been thrown")
	}
	t.Log(format)
	t.Log(err)
}

func TestFormatParseFailEmptyVar(t *testing.T) {
	format := Format{
		FormatElement{Type:FormatVar, Name:"wally"},
		FormatElement{Type:FormatChar, Char:','},
		FormatElement{Type:FormatVar, Name:"wizard"},
		FormatElement{Type:FormatChar, Char:'.'},
	}
	_, err := format.Parse([]byte("wat,."))
	if err == nil {
		t.Fatal("Parse error should have been thrown")
	}
	t.Log(format)
	t.Log(err)
}

func TestParseFormat(t *testing.T) {
	format := Format{
		FormatElement{Type:FormatVar, Name:"wally"},
		FormatElement{Type:FormatChar, Char:','},
		FormatElement{Type:FormatVar, Name:"wizard"},
		FormatElement{Type:FormatChar, Char:'.'},
	}
	t.Log(format.VerboseString())
	f2, err := ParseFormatString(format.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(f2.VerboseString())
	if f2.VerboseString() != format.VerboseString() {
		t.Log(format.VerboseString())
		t.Log(f2.VerboseString())
		t.Fatal("parsed not equal to original")
	}
	meta, err := f2.Parse([]byte("wat,we."))
	if err != nil {
		t.Fatal(err)
	}
	if meta["wally"] != "wat" {
		t.Log(meta)
		t.Fatal("wally != wat")
	}
	if meta["wizard"] != "we" {
		t.Log(meta)
		t.Fatal("wizard != we")
	}
}

func TestRealFormat(t *testing.T) {
	format, err := ParseFormatString("$(slide) $(sample) $(region) $(stain).tif")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(format.VerboseString())
	meta, err := format.Parse([]byte("1 WT16226 L1 CD34.tif"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(meta)
	if meta["slide"] != "1" {
		t.Fatal("slide != 1")
	}
	if meta["sample"] != "WT16226" {
		t.Fatal("sample != WT16226")
	}
	if meta["region"] != "L1" {
		t.Fatal("region != L1")
	}
	if meta["stain"] != "CD34" {
		t.Fatal("stain != CD34")
	}
}

