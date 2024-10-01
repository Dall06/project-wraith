package tools

import "testing"

func TestSha512(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		input    string
		expected string
	}{
		{
			name:     "Basic case",
			secret:   "mysecret",
			input:    "myinput",
			expected: "fe4e9c1f01670b08e1de4634c838d341555e72169b16e80f9a10dfa25345a135bd2f99bfcd39533da339eec19ac0d40e2c4d420d6b9660a37f6ecf2746419ba7",
		},
		{
			name:     "Empty input",
			secret:   "mysecret",
			input:    "",
			expected: "9d86b9ac13df1a34844ae4d990f8e3907e8beffe0ce6929e6c8e1593cd2ea30805d49cf521138f07f8d02220f50948c2499c22d69f68f7032ea3560b6cb22844",
		},
		{
			name:     "Empty secret",
			secret:   "",
			input:    "myinput",
			expected: "993dc36fef649f2f20d822dd59308570af18604cbe694919690834106fb245089a236cd92d97ffd16c7e2dc2b053d253188c5ba4bc21799d3074eaa967719fc6",
		},
		{
			name:     "Empty secret and input",
			secret:   "",
			input:    "",
			expected: "b936cee86c9f87aa5d3c6f2e84cb5a4239a5fe50480a6ec66b70ab5b1f4ac6730c6c515421b327ec1d69402e53dfb49ad7381eb067b338fd7b0cb22247225d47",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sha512(tt.secret, tt.input)
			if result != tt.expected {
				t.Errorf("Sha512(%q, %q) = %q, want %q", tt.secret, tt.input, result, tt.expected)
			}
		})
	}
}
