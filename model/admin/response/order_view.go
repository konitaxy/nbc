package response

type OrderViewDetail struct {
	Author        string `json:"author"`
	Size          string `json:"size"`
	ProductType   string `json:"productType"`
	ReferenceCode string `json:"referenceCode"`
	Title         string `json:"title"`
	Category      string `json:"category"`
	Weight        string `json:"weight"`
	Frame         string `json:"frame"`
	PrintDate     string `json:"printDate"`
	ArtworkUrl    string `json:"artworkUrl"`
}
