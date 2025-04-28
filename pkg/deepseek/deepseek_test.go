package deepseek

import "testing"

func TestDeepseek(t *testing.T) {
	r := NewClient("sk-")
	_, err := r.AssessTvNames("基督山伯爵", 2025, []string{"The Count of Monte Cristo 2024 S01 1080p WEB-DL DD 5.1 H.264-playWEB", "The Count of Monte Cristo 2024 S01E06-08 MULTi 1080p WEB H264-AMB3R"})
	if err != nil {
		t.Fatal(err)
	}
}
