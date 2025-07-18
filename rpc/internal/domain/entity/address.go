package entity

// Region 地區
type Region int

const (
	RegionNorth   Region = iota // 0: 北部
	RegionCentral               // 1: 中部
	RegionSouth                 // 2: 南部
	RegionEast                  // 3: 東部
	RegionIsland                // 4: 離島
)

// Address 地址實體
type Address struct {
	ID        int64
	Recipient string
	Region    Region
	Detail    string
}
