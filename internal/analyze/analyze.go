package analyze

import (
	"CSGO-stats/internal/structures"
	"CSGO-stats/internal/visualization"
	"sort"
)

func AnalyzeTournament(tournament structures.Tournament) {
	data, amount := CalculateStatistics(tournament)
	visualization.GenerateCharts(data, tournament.TournamentName, amount, "https://www.hltv.org/events/7319/la-liga-pro-2023")
}

// CalculateStatistics подсчитывает статистику всех игроков за все карты турнира
func CalculateStatistics(tournament structures.Tournament) (structures.ChartData, int) {
	sumStats := make(map[uint64]structures.SummaryStatistics)
	matchAmount := 0
	for _, matches := range tournament.Matches {
		for _, csMap := range matches.Maps {
			matchAmount++
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
					stats := structures.SummaryStatistics{
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

	v := make([]structures.SummaryStatistics, 0, len(sumStats))
	for _, value := range sumStats {
		v = append(v, value)
	}

	chartData := make(structures.ChartData)

	sort.Sort(ByKills(v))
	SortedByKills := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByKills, v)
	chartData["SortedByKills"] = SortedByKills

	sort.Sort(ByDeath(v))
	SortedByDeath := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByDeath, v)
	chartData["SortedByDeath"] = SortedByDeath

	sort.Sort(ByAssists(v))
	SortedByAssists := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByAssists, v)
	chartData["SortedByAssists"] = SortedByAssists

	sort.Sort(ByFootSteps(v))
	SortedByFootSteps := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByFootSteps, v)
	chartData["SortedByFootSteps"] = SortedByFootSteps

	sort.Sort(ByDuckKills(v))
	SortedByDuckKills := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByDuckKills, v)
	chartData["SortedByDuckKills"] = SortedByDuckKills

	sort.Sort(ByFlashedKills(v))
	SortedByFlashedKills := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByFlashedKills, v)
	chartData["SortedByFlashedKills"] = SortedByFlashedKills

	sort.Sort(ByAirborneKills(v))
	SortedByAirborneKills := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByAirborneKills, v)
	chartData["SortedByAirborneKills"] = SortedByAirborneKills

	sort.Sort(ByWallbangKills(v))
	SortedByWallbangKills := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByWallbangKills, v)
	chartData["SortedByWallbangKills"] = SortedByWallbangKills

	sort.Sort(BySmokeKills(v))
	SortedBySmokeKills := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedBySmokeKills, v)
	chartData["SortedBySmokeKills"] = SortedBySmokeKills

	sort.Sort(ByNoScopeKills(v))
	SortedByNoScopeKills := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByNoScopeKills, v)
	chartData["SortedByNoScopeKills"] = SortedByNoScopeKills

	sort.Sort(ByWeaponShots(v))
	SortedByWeaponShots := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByWeaponShots, v)
	chartData["SortedByWeaponShots"] = SortedByWeaponShots

	sort.Sort(ByWeaponReloads(v))
	SortedByWeaponReloads := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByWeaponReloads, v)
	chartData["SortedByWeaponReloads"] = SortedByWeaponReloads

	sort.Sort(ByPlayerJumps(v))
	SortedByPlayerJumps := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByPlayerJumps, v)
	chartData["SortedByPlayerJumps"] = SortedByPlayerJumps

	sort.Sort(ByBombDrops(v))
	SortedByBombDrops := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByBombDrops, v)
	chartData["SortedByBombDrops"] = SortedByBombDrops

	sort.Sort(BySmokes(v))
	SortedBySmokes := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedBySmokes, v)
	chartData["SortedBySmokes"] = SortedBySmokes

	sort.Sort(ByHEnades(v))
	SortedByHEnades := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByHEnades, v)
	chartData["SortedByHEnades"] = SortedByHEnades

	sort.Sort(ByMolotov(v))
	SortedByMolotov := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByMolotov, v)
	chartData["SortedByMolotov"] = SortedByMolotov

	sort.Sort(ByCTmoly(v))
	SortedByCTmoly := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByCTmoly, v)
	chartData["SortedByCTmoly"] = SortedByCTmoly

	sort.Sort(ByFlasbang(v))
	SortedByFlasbang := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByFlasbang, v)
	chartData["SortedByFlasbang"] = SortedByFlasbang

	sort.Sort(ByDecoy(v))
	SortedByDecoy := make([]structures.SummaryStatistics, len(sumStats))
	copy(SortedByDecoy, v)
	chartData["SortedByDecoy"] = SortedByDecoy

	return chartData, matchAmount
}
