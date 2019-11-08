package oam

// ServiceList JSON struct
type AfServiceList struct {
        AfServices []AfService      `json:"afServices,omitempty"`
}

// AF Service JSON struct
type AfService struct {
        AfInstance    string         `json:"afInstance,omitempty"`
        LocalServices []LocalService `json:"localServices,omitempty"`
}


// local Service JSON struct
type LocalService struct {
        Dnai   string                   `json:"dnai,omitempty"`
        Dnn    string                   `json:"dnn,omitempty"`
        Dns    string                   `json:"dns,omitempty"`
}

//  AfId struct
type AfId struct {
        AfId   string                   `json:"afid,omitempty"`
}
