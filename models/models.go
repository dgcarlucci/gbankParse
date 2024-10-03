package models

type Config struct {
	InputFilePath   string `json:"input_file_path"`
	OutputDirectory string `json:"output_file_path"`
	ServerName      string `json:"server_name"`
	CharacterName   string `json:"character_name"`
}

/*{
["ID"] = 3857,
["Info"] = {
["icon"] = 134579,
["level"] = 30,
["rarity"] = 1,
["equipId"] = 0,
["price"] = 125,
["class"] = 7,
["subClass"] = 0,
["name"] = "Coal",
},
["Count"] = 1,
["Link"] = "|cffffffff|Hitem:3857::::::::1:::::::::|h[Coal]|h|r",
},
*/

type Item struct {
	Id    string `csv:"ID"`
	Info  Info   `csv:"Info"`
	Count int    `csv:"Count"`
	Link  string `csv:"Link"`
}

type Info struct {
	Icon     int
	Level    int
	Rarity   int
	EquipId  int
	Price    int
	Class    int
	SubClass int
	Name     string
}

type Alts struct {
	Items []Item `csv:"items"`
}

type Alliance struct {
	Alts Alts `csv:"alts"`
}

type GBankClassicDB struct {
	ProfileKeys map[string]string `csv:"profileKeys"`
	Faction     struct {
		Alliance Alliance `csv:"Alliance"`
	} `csv:"faction"`
}

// WowItem struct
type WoWItem struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Quality struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"quality"`
	Level        int `json:"level"`
	RequiredLvel int `json:"required_level"`
	Media        struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		ID int `json:"id"`
	} `json:"media"`
	ItemClass struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"item_class"`
	ItemSubclass struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"item_subclass"`
	InventoryType struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"inventory_type"`
	PurchasePrice int  `json:"purchase_price"`
	SellPrice     int  `json:"sell_price"`
	MaxCount      int  `json:"max_count"`
	IsEquippable  bool `json:"is_equippable"`
	IsStackable   bool `json:"is_stackable"`
	PreviewItem   struct {
		Item struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			ID int `json:"id"`
		} `json:"item"`
		Quality struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"quality"`
		Name  string `json:"name"`
		Media struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			ID int `json:"id"`
		} `json:"media"`
		ItemClass struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"item_class"`
		ItemSubclass struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"item_subclass"`
		InventoryType struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"inventory_type"`
		SellPrice struct {
			Value          int `json:"value"`
			DisplayStrings struct {
				Header string `json:"header"`
				Gold   string `json:"gold"`
				Silver string `json:"silver"`
				Copper string `json:"copper"`
			} `json:"display_strings"`
		} `json:"sell_price"`
		ContainerSlots struct {
			Value         int    `json:"value"`
			DisplayString string `json:"display_string"`
		} `json:"container_slots"`
	} `json:"preview_item"`
	PurchaseQuantity int `json:"purchase_quantity"`
}
