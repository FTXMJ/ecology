package models

type PeerHistory struct {
	Id               int     `gorm:"column:id;primary_key" json:"id"`
	Time             string  `gorm:"column:time" json:"time"`
	WholeNetworkTfor float64 `gorm:"column:whole_network_tfor" json:"whole_network_tfor"`
	PeerABouns       float64 `gorm:"column:peer_a_bouns" json:"peer_a_bouns"`
	DiamondsPeer     int     `gorm:"column:diamonds_peer" json:"diamonds_peer"`
	SuperPeer        int     `gorm:"column:super_peer" json:"super_peer"`
	CreationPeer     int     `gorm:"column:creation_peer" json:"creation_peer"`
}
