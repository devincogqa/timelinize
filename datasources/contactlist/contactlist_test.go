package contactlist

import (
	"context"
	"io/fs"
	"os"
	"testing"

	"github.com/timelinize/timelinize/timeline"
)

// TestBestColumnMappingAndDelim_NonCommaCSV is a regression test for commit 7146add3
// which fixed support for non-comma delimited CSV files. Before the fix,
// bestColumnMappingAndDelim did not exist and the delimiter was always assumed
// to be a comma, causing tab-delimited and semicolon-delimited contact lists
// to fail recognition and import.
func TestBestColumnMappingAndDelim_NonCommaCSV(t *testing.T) {
	testdata := os.DirFS("testdata")

	tests := []struct {
		name          string
		filename      string
		wantDelim     rune
		wantMinFields int // minimum number of canonical fields we expect mapped
	}{
		{
			name:          "tab-delimited",
			filename:      "contacts_tab.tsv",
			wantDelim:     '\t',
			wantMinFields: recognizeAtLeastFields,
		},
		{
			name:          "semicolon-delimited",
			filename:      "contacts_semicolon.csv",
			wantDelim:     ';',
			wantMinFields: recognizeAtLeastFields,
		},
		{
			name:          "comma-delimited",
			filename:      "contacts_comma.csv",
			wantDelim:     ',',
			wantMinFields: recognizeAtLeastFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirEntry := timeline.DirEntry{
				FS:       testdata,
				Filename: tt.filename,
				DirEntry: fakeDirEntry{name: tt.filename},
			}

			mapping, delim, err := bestColumnMappingAndDelim(context.Background(), dirEntry, ".")
			if err != nil {
				t.Fatalf("bestColumnMappingAndDelim returned error: %v", err)
			}

			if delim != tt.wantDelim {
				t.Errorf("detected delimiter = %q, want %q", delim, tt.wantDelim)
			}

			if len(mapping) < tt.wantMinFields {
				t.Errorf("mapped %d canonical fields, want at least %d; mapping = %v",
					len(mapping), tt.wantMinFields, mapping)
			}

			// Verify that expected canonical fields are present in the mapping
			expectedFields := []string{"first_name", "last_name"}
			for _, field := range expectedFields {
				if _, ok := mapping[field]; !ok {
					t.Errorf("expected canonical field %q not found in mapping %v", field, mapping)
				}
			}
		})
	}
}

// TestBestColumnMappingForFields_GenericFormat verifies that column headers
// added in commit 7146add3 (e.g., "firstname", "lastname", "bday") are
// correctly matched by the Generic format.
func TestBestColumnMappingForFields_GenericFormat(t *testing.T) {
	tests := []struct {
		name       string
		headerRow  []string
		wantFields []string
	}{
		{
			name:       "short aliases added in fix",
			headerRow:  []string{"firstname", "lastname", "email", "phone"},
			wantFields: []string{"first_name", "last_name", timeline.AttributeEmail, timeline.AttributePhoneNumber},
		},
		{
			name:       "bday alias added in fix",
			headerRow:  []string{"fname", "lname", "bday", "email"},
			wantFields: []string{"first_name", "last_name", "birthdate", timeline.AttributeEmail},
		},
		{
			name:       "middlename alias added in fix",
			headerRow:  []string{"firstname", "middlename", "lastname", "email"},
			wantFields: []string{"first_name", "middle_name", "last_name", timeline.AttributeEmail},
		},
		{
			name:       "given names alias added in fix",
			headerRow:  []string{"given names", "surname", "email", "phone"},
			wantFields: []string{"first_name", "last_name", timeline.AttributeEmail, timeline.AttributePhoneNumber},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := bestColumnMappingForFields(tt.headerRow)
			for _, field := range tt.wantFields {
				if _, ok := mapping[field]; !ok {
					t.Errorf("expected canonical field %q not found in mapping %v", field, mapping)
				}
			}
		})
	}
}

// TestDetermineColumnMappingForDelimiter_NonComma verifies that
// determineColumnMappingForDelimiter correctly parses files with
// non-comma delimiters.
func TestDetermineColumnMappingForDelimiter_NonComma(t *testing.T) {
	testdata := os.DirFS("testdata")

	tests := []struct {
		name     string
		filename string
		delim    rune
		wantNil  bool // true if we expect nil mapping (wrong delimiter)
	}{
		{
			name:     "tab file with tab delimiter",
			filename: "contacts_tab.tsv",
			delim:    '\t',
			wantNil:  false,
		},
		{
			name:     "tab file with comma delimiter (wrong)",
			filename: "contacts_tab.tsv",
			delim:    ',',
			wantNil:  true,
		},
		{
			name:     "semicolon file with semicolon delimiter",
			filename: "contacts_semicolon.csv",
			delim:    ';',
			wantNil:  false,
		},
		{
			name:     "semicolon file with comma delimiter (wrong)",
			filename: "contacts_semicolon.csv",
			delim:    ',',
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirEntry := timeline.DirEntry{
				FS:       testdata,
				Filename: tt.filename,
				DirEntry: fakeDirEntry{name: tt.filename},
			}

			mapping, err := determineColumnMappingForDelimiter(dirEntry, ".", tt.delim)
			if err != nil {
				t.Fatalf("determineColumnMappingForDelimiter returned error: %v", err)
			}

			if tt.wantNil && mapping != nil {
				t.Errorf("expected nil mapping for wrong delimiter, got %v", mapping)
			}
			if !tt.wantNil && mapping == nil {
				t.Errorf("expected non-nil mapping for correct delimiter, got nil")
			}
			if !tt.wantNil && len(mapping) < recognizeAtLeastFields {
				t.Errorf("mapped %d canonical fields, want at least %d", len(mapping), recognizeAtLeastFields)
			}
		})
	}
}

// fakeDirEntry implements fs.DirEntry for testing purposes.
type fakeDirEntry struct {
	name  string
	isDir bool
}

func (f fakeDirEntry) Name() string               { return f.name }
func (f fakeDirEntry) IsDir() bool                 { return f.isDir }
func (f fakeDirEntry) Type() fs.FileMode           { return 0 }
func (f fakeDirEntry) Info() (fs.FileInfo, error)  { return nil, nil }
