package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

func specInSlice(a proto.Spec, list []proto.Spec) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func writeGemFile(outDir string, gemsData []GemData) {
	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fmt.Sprintf("%s/all_gems.go", outDir))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(`// DO NOT EDIT. This file is auto-generated by the item generator tool. Use that to make edits.
	
package items
	
import (
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var Gems = []Gem{
`)

	for _, gemData := range gemsData {
		if gemData.Declaration.Filter {
			continue
		}
		allow := allowList[gemData.Declaration.ID]
		if !allow {
			if gemData.Response.Quality < int(proto.ItemQuality_ItemQualityUncommon) {
				continue
			}
			if gemData.Response.GetPhase() == 0 {
				continue
			}
		}
		file.WriteString(fmt.Sprintf("\t%s,\n", gemToGoString(gemData.Declaration, gemData.Response)))
	}

	file.WriteString("}\n")

	file.Sync()
}

func writeItemFile(outDir string, itemsData []ItemData) {
	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(fmt.Sprintf("%s/all_items.go", outDir))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(`// DO NOT EDIT. This file is auto-generated by the item generator tool. Use that to make edits.
	
package items
	
import (
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

var Items = []Item{
`)

	for _, itemData := range itemsData {
		itemLevel := itemData.Response.GetItemLevel()
		if itemData.Declaration.Filter {
			continue
		}
		deny := false
		for _, pattern := range denyListNameRegexes {
			if pattern.MatchString(itemData.Response.Name) {
				deny = true
				break
			}
		}
		if deny {
			continue
		}
		if !itemData.Response.IsEquippable() {
			continue
		}
		allow := allowList[itemData.Declaration.ID]
		if !allow {
			if itemData.Response.Quality < int(proto.ItemQuality_ItemQualityUncommon) {
				continue
			} else if itemData.Response.Quality < int(proto.ItemQuality_ItemQualityEpic) {
				if itemLevel < 105 {
					continue
				}
			} else {
				// Epic and legendary items might come from classic, so use a lower ilvl threshold.
				if itemLevel < 75 {
					continue
				}
			}
		}
		if itemLevel == 0 {
			fmt.Printf("Missing ilvl: %s", itemData.Response.Name)
		}

		file.WriteString(fmt.Sprintf("\t%s,\n", itemToGoString(itemData.Declaration, itemData.Response)))
	}

	file.WriteString("}\n")

	file.Sync()
}

func gemToGoString(gemDeclaration GemDeclaration, gemResponse WowheadItemResponse) string {
	gemStr := "{"

	gemStr += fmt.Sprintf("Name:\"%s\", ", gemResponse.Name)
	gemStr += fmt.Sprintf("ID:%d, ", gemDeclaration.ID)

	phase := gemDeclaration.Phase
	if phase == 0 {
		phase = gemResponse.GetPhase()
	}
	gemStr += fmt.Sprintf("Phase:%d, ", phase)
	gemStr += fmt.Sprintf("Quality:proto.ItemQuality_%s, ", proto.ItemQuality(gemResponse.Quality).String())
	gemStr += fmt.Sprintf("Color:proto.GemColor_%s, ", proto.GemColor(gemResponse.GetSocketColor()).String())
	gemStr += fmt.Sprintf("Stats: %s, ", statsToGoString(gemResponse.GetGemStats(), gemDeclaration.Stats))

	if gemResponse.GetUnique() {
		gemStr += fmt.Sprintf("Unique:true, ")
	}

	gemStr += "}"
	return gemStr
}

func itemToGoString(itemDeclaration ItemDeclaration, itemResponse WowheadItemResponse) string {
	itemStr := "{"

	itemStr += fmt.Sprintf("Name:\"%s\", ", strings.ReplaceAll(itemResponse.Name, "\"", "\\\""))
	itemStr += fmt.Sprintf("ID:%d, ", itemDeclaration.ID)

	classAllowlist := itemResponse.GetClassAllowlist()
	if len(itemDeclaration.ClassAllowlist) > 0 {
		classAllowlist = itemDeclaration.ClassAllowlist
	}
	if len(classAllowlist) > 0 {
		itemStr += "ClassAllowlist: []proto.Class{"
		for _, class := range classAllowlist {
			itemStr += fmt.Sprintf("proto.Class_%s,", class.String())
		}
		itemStr += "}, "
	}

	itemStr += fmt.Sprintf("Type:proto.ItemType_%s, ", itemResponse.GetItemType().String())

	armorType := itemResponse.GetArmorType()
	if armorType != proto.ArmorType_ArmorTypeUnknown {
		itemStr += fmt.Sprintf("ArmorType:proto.ArmorType_%s, ", armorType.String())
	}

	weaponType := itemResponse.GetWeaponType()
	if weaponType != proto.WeaponType_WeaponTypeUnknown {
		itemStr += fmt.Sprintf("WeaponType:proto.WeaponType_%s, ", weaponType.String())

		handType := itemResponse.GetHandType()
		if itemDeclaration.HandType != proto.HandType_HandTypeUnknown {
			handType = itemDeclaration.HandType
		}
		if handType == proto.HandType_HandTypeUnknown {
			panic("Unknown hand type for item: " + itemResponse.Tooltip)
		}
		itemStr += fmt.Sprintf("HandType:proto.HandType_%s, ", handType.String())
	} else {
		rangedWeaponType := itemResponse.GetRangedWeaponType()
		if rangedWeaponType != proto.RangedWeaponType_RangedWeaponTypeUnknown {
			itemStr += fmt.Sprintf("RangedWeaponType:proto.RangedWeaponType_%s, ", rangedWeaponType.String())
		}
	}

	min, max := itemResponse.GetWeaponDamage()
	if min != 0 && max != 0 {
		itemStr += fmt.Sprintf("WeaponDamageMin: %0.1f, ", min)
		itemStr += fmt.Sprintf("WeaponDamageMax: %0.1f, ", max)
	}
	speed := itemResponse.GetWeaponSpeed()
	if speed != 0 {
		itemStr += fmt.Sprintf("SwingSpeed: %0.2f, ", speed)
	}

	phase := itemDeclaration.Phase
	if phase == 0 {
		phase = itemResponse.GetPhase()
	}
	itemStr += fmt.Sprintf("Phase:%d, ", phase)
	itemStr += fmt.Sprintf("Quality:proto.ItemQuality_%s, ", proto.ItemQuality(itemResponse.Quality).String())

	if itemResponse.GetUnique() {
		itemStr += fmt.Sprintf("Unique:true, ")
	}

	itemStr += fmt.Sprintf("Ilvl:%d, ", itemResponse.GetItemLevel())

	itemStr += fmt.Sprintf("Stats: %s, ", statsToGoString(itemResponse.GetStats(), itemDeclaration.Stats))

	gemSockets := itemResponse.GetGemSockets()
	if len(gemSockets) > 0 {
		itemStr += "GemSockets: []proto.GemColor{"
		for _, gemColor := range gemSockets {
			itemStr += fmt.Sprintf("proto.GemColor_%s,", gemColor.String())
		}
		itemStr += "}, "
	}

	itemStr += fmt.Sprintf("SocketBonus: %s", statsToGoString(itemResponse.GetSocketBonus(), Stats{}))

	setName := itemResponse.GetItemSetName()
	if setName != "" {
		itemStr += fmt.Sprintf(", SetName: \"%s\"", setName)
	}

	itemStr += "}"
	return itemStr
}

func statsToGoString(statlist Stats, overrides Stats) string {
	statsStr := "stats.Stats{"

	for stat, value := range statlist {
		val := value
		if overrides[stat] > 0 {
			val = overrides[stat]
		}
		if value > 0 {
			statsStr += fmt.Sprintf("stats.%s:%.0f,", stats.Stat(stat).StatName(), val)
		}
	}

	statsStr += "}"
	return statsStr
}

// If any of these match the item name, don't include it.
var denyListNameRegexes = []*regexp.Regexp{
	regexp.MustCompile("PH\\]"),
	regexp.MustCompile("TEST"),
	regexp.MustCompile("Test"),
	regexp.MustCompile("Bracer 3"),
	regexp.MustCompile("Bracer 2"),
	regexp.MustCompile("Bracer 1"),
	regexp.MustCompile("Boots 3"),
	regexp.MustCompile("Boots 2"),
	regexp.MustCompile("Boots 1"),
	regexp.MustCompile("zOLD"),
	regexp.MustCompile("30 Epic"),
	regexp.MustCompile("Indalamar"),
	regexp.MustCompile("QR XXXX"),
	regexp.MustCompile("Deprecated: Keanna"),
	regexp.MustCompile("90 Epic"),
	regexp.MustCompile("66 Epic"),
	regexp.MustCompile("63 Blue"),
	regexp.MustCompile("90 Green"),
	regexp.MustCompile("63 Green"),
}

// allowList allows overriding to allow an item
var allowList = map[int]bool{
	11815: true, // Hand of Justice
	12632: true, // Storm Gauntlets
	17111: true, // Blazefury Medallion
	17112: true, // Empyrean Demolisher
	19808: true, // Rockhide Strongfish
	20966: true, // Jade Pendant of Blasting
	22395: true, // Totem of Rage
	24114: true, // Braided Eternium Chain
	27947: true, // Totem of Impact
	28041: true, // Bladefist's Breadth
	31139: true, // Fist of Reckoning
	31149: true, // Gloves of Pandemonium
	31193: true, // Blade of Unquenched Thirst
	32508: true, // Necklace of the Deep
	33135: true, // Falling Star
	33140: true, // Blood of Amber
	33143: true, // Stone of Blades
	33144: true, // Facet of Eternity
	6360:  true, // Steelscale Crushfish
	8345:  true, // Wolfshead Helm
	28032: true, // Delicate Green Poncho
}
