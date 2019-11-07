package oam

// ServiceList JSON struct
type AfRegisterList struct {
        AfRegisters []AfRegister      `json:"afRegisters,omitempty"`
}

// af register JSON struct
type AfRegister struct {
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
