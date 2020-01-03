package models

type PeerHistory struct {
	Id               int     `json:"id"`
	Time             string  `json:"time"`
	WholeNetworkTfor float64 `json:"whole_network_tfor"`
	PeerABouns       float64 `json:"peer_a_bouns"`
	DiamondsPeer     int     `json:"diamonds_peer"`
	SuperPeer        int     `json:"super_peer"`
	CreationPeer     int     `json:"creation_peer"`
}
