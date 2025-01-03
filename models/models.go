package models

type Config struct {
	InputFilePath    string `json:"input_file_path"`
	OutputDirectory  string `json:"output_file_path"`
	ServerID         string `json:"server_id"`
	CharacterName    string `json:"character_name"`
	DiscordChannel   string `json:"discord_channel"`
	RaidHelperAPIKey string `json:"raid_helper_api_key"`
}

/*
BODY - The body of the embed.
https://raid-helper.dev/documentation/api

mentions<string> - The roles you want to ping, separate multiple entries with a comma.
title<object> - The title field of this embed.

	text<string> - The title text.
	URL<string> - The URL the title will link to.

description<string> - The description for this embed.
imageURL<string> - The image URL for this embed.
thumbnailURL<string> - The thumbnail URL for this embed.
color<string> - The embed hex color.
fields<array of objects> - The embed fields.

	name<string> - The field name.
	value<string> - The field value.
	inline<boolean> - whether the field is to be inline.

author<object> - The author field of this embed.

	name<string> - The name of the author.
	URL<string> - The URL the author field will point to.
	iconURL<string> - The URL of an image that will be displayed.

footer<object> - The footer field of this embed.

	text<string> - The footer text.
	iconURL<string> - The URL of an image that will be displayed.
*/
type RaidHelperEmbedMessage struct {

	//Title of the embed
	Title RaidHelperTitle `json:"title"`

	//Description of the embed
	Description string `json:"description"`

	//Image URL of the embed
	ImageURL string `json:"imageURL"`

	//Thumbnail URL of the embed
	ThumbnailURL string `json:"thumbnailURL"`

	//Hex color of the embed
	Color string `json:"color"`

	//Fields of the embed
	Fields []RaidHelperEmbedField `json:"fields"`

	//Author of the embed
	Author RaidHelperEmbedAuthor `json:"author"`

	//Footer of the embed
	Footer RaidHelperEmbedFooter `json:"footer"`
}

type RaidHelperTitle struct {
	Text string `json:"text"`
	URL  string `json:"URL"`
}

type RaidHelperEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type RaidHelperEmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"URL"`
	IconURL string `json:"iconURL"`
}

type RaidHelperEmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"iconURL"`
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
