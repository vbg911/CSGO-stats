package analyze

import (
	"CSGO-stats/internal/structures"
	"fmt"
	"sort"
)

func AnalyzeTournament(tournament structures.Tournament) {
	CalculateStatistics(tournament)
}

type summaryStatistics struct {
	MapPlayed int `json:"mapPlayed"`
	structures.PlayerStats
}

type ByKills []summaryStatistics

func (a ByKills) Len() int           { return len(a) }
func (a ByKills) Less(i, j int) bool { return a[i].Kills > a[j].Kills }
func (a ByKills) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// CalculateStatistics подсчитывает статистику всех игроков за все карты турнира
func CalculateStatistics(tournament structures.Tournament) {

	sumStats := make(map[uint64]summaryStatistics)

	for _, matches := range tournament.Matches {
		for _, csMap := range matches.Maps {
			for _, player := range csMap.Players {
				if val, ok := sumStats[player.SteamID64]; ok {
					val.MapPlayed += 1
					val.MVP += player.MVP
					val.Kills += player.Kills
					val.Deaths += player.Deaths
					val.Assists += player.Assists
					val.FootSteps += player.FootSteps
					val.DuckKills += player.DuckKills
					val.FlashedKills += player.FlashedKills
					val.AirborneKills += player.AirborneKills
					val.WallbangKills += player.WallbangKills
					val.SmokeKills += player.SmokeKills
					val.NoScopeKills += player.NoScopeKills
					val.WeaponShots += player.WeaponShots
					val.WeaponReloads += player.WeaponReloads
					val.PlayerJumps += player.PlayerJumps
					val.BombDrops += player.BombDrops
					val.Smokes += player.Smokes
					val.HEnades += player.HEnades
					val.Molotov += player.Molotov
					val.CTmoly += player.CTmoly
					val.Flasbang += player.Flasbang
					val.Decoy += player.Decoy
					sumStats[player.SteamID64] = val
				} else {
					stats := summaryStatistics{
						MapPlayed: 1,
						PlayerStats: structures.PlayerStats{
							SteamID64:     player.SteamID64,
							Name:          player.Name,
							TeamName:      player.TeamName,
							CrosshairCode: player.CrosshairCode,
							MVP:           player.MVP,
							Kills:         player.Kills,
							Deaths:        player.Deaths,
							Assists:       player.Assists,
							FootSteps:     player.FootSteps,
							DuckKills:     player.DuckKills,
							FlashedKills:  player.FlashedKills,
							AirborneKills: player.AirborneKills,
							WallbangKills: player.WallbangKills,
							SmokeKills:    player.SmokeKills,
							NoScopeKills:  player.NoScopeKills,
							WeaponShots:   player.WeaponShots,
							WeaponReloads: player.WeaponReloads,
							PlayerJumps:   player.PlayerJumps,
							BombDrops:     player.BombDrops,
							Smokes:        player.Smokes,
							HEnades:       player.HEnades,
							Molotov:       player.Molotov,
							CTmoly:        player.CTmoly,
							Flasbang:      player.Flasbang,
							Decoy:         player.Decoy,
						},
					}
					sumStats[player.SteamID64] = stats
				}
			}
		}
	}

	v := make([]summaryStatistics, 0, len(sumStats))
	for _, value := range sumStats {
		v = append(v, value)
	}

	sort.Sort(ByKills(v))
	for _, i := range v {
		fmt.Println(i.Name, " maps: ", i.MapPlayed, " kills: ", i.Kills)
	}
	//fmt.Println(sumStats)
}
