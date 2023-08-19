package analyze

import "CSGO-stats/internal/structures"

type ByKills []structures.SummaryStatistics

func (a ByKills) Len() int           { return len(a) }
func (a ByKills) Less(i, j int) bool { return a[i].Kills > a[j].Kills }
func (a ByKills) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByDeath []structures.SummaryStatistics

func (a ByDeath) Len() int           { return len(a) }
func (a ByDeath) Less(i, j int) bool { return a[i].Deaths > a[j].Deaths }
func (a ByDeath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByAssists []structures.SummaryStatistics

func (a ByAssists) Len() int           { return len(a) }
func (a ByAssists) Less(i, j int) bool { return a[i].Assists > a[j].Assists }
func (a ByAssists) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByFootSteps []structures.SummaryStatistics

func (a ByFootSteps) Len() int           { return len(a) }
func (a ByFootSteps) Less(i, j int) bool { return a[i].FootSteps > a[j].FootSteps }
func (a ByFootSteps) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByDuckKills []structures.SummaryStatistics

func (a ByDuckKills) Len() int           { return len(a) }
func (a ByDuckKills) Less(i, j int) bool { return a[i].DuckKills > a[j].DuckKills }
func (a ByDuckKills) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByFlashedKills []structures.SummaryStatistics

func (a ByFlashedKills) Len() int           { return len(a) }
func (a ByFlashedKills) Less(i, j int) bool { return a[i].FlashedKills > a[j].FlashedKills }
func (a ByFlashedKills) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByAirborneKills []structures.SummaryStatistics

func (a ByAirborneKills) Len() int           { return len(a) }
func (a ByAirborneKills) Less(i, j int) bool { return a[i].AirborneKills > a[j].AirborneKills }
func (a ByAirborneKills) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByWallbangKills []structures.SummaryStatistics

func (a ByWallbangKills) Len() int           { return len(a) }
func (a ByWallbangKills) Less(i, j int) bool { return a[i].WallbangKills > a[j].WallbangKills }
func (a ByWallbangKills) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type BySmokeKills []structures.SummaryStatistics

func (a BySmokeKills) Len() int           { return len(a) }
func (a BySmokeKills) Less(i, j int) bool { return a[i].SmokeKills > a[j].SmokeKills }
func (a BySmokeKills) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByNoScopeKills []structures.SummaryStatistics

func (a ByNoScopeKills) Len() int           { return len(a) }
func (a ByNoScopeKills) Less(i, j int) bool { return a[i].NoScopeKills > a[j].NoScopeKills }
func (a ByNoScopeKills) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByWeaponShots []structures.SummaryStatistics

func (a ByWeaponShots) Len() int           { return len(a) }
func (a ByWeaponShots) Less(i, j int) bool { return a[i].WeaponShots > a[j].WeaponShots }
func (a ByWeaponShots) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByWeaponReloads []structures.SummaryStatistics

func (a ByWeaponReloads) Len() int           { return len(a) }
func (a ByWeaponReloads) Less(i, j int) bool { return a[i].WeaponReloads > a[j].WeaponReloads }
func (a ByWeaponReloads) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByPlayerJumps []structures.SummaryStatistics

func (a ByPlayerJumps) Len() int           { return len(a) }
func (a ByPlayerJumps) Less(i, j int) bool { return a[i].PlayerJumps > a[j].PlayerJumps }
func (a ByPlayerJumps) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByBombDrops []structures.SummaryStatistics

func (a ByBombDrops) Len() int           { return len(a) }
func (a ByBombDrops) Less(i, j int) bool { return a[i].BombDrops > a[j].BombDrops }
func (a ByBombDrops) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type BySmokes []structures.SummaryStatistics

func (a BySmokes) Len() int           { return len(a) }
func (a BySmokes) Less(i, j int) bool { return a[i].Smokes > a[j].Smokes }
func (a BySmokes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByHEnades []structures.SummaryStatistics

func (a ByHEnades) Len() int           { return len(a) }
func (a ByHEnades) Less(i, j int) bool { return a[i].HEnades > a[j].HEnades }
func (a ByHEnades) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByMolotov []structures.SummaryStatistics

func (a ByMolotov) Len() int           { return len(a) }
func (a ByMolotov) Less(i, j int) bool { return a[i].Molotov > a[j].Molotov }
func (a ByMolotov) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByCTmoly []structures.SummaryStatistics

func (a ByCTmoly) Len() int           { return len(a) }
func (a ByCTmoly) Less(i, j int) bool { return a[i].CTmoly > a[j].CTmoly }
func (a ByCTmoly) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByFlasbang []structures.SummaryStatistics

func (a ByFlasbang) Len() int           { return len(a) }
func (a ByFlasbang) Less(i, j int) bool { return a[i].Flasbang > a[j].Flasbang }
func (a ByFlasbang) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByDecoy []structures.SummaryStatistics

func (a ByDecoy) Len() int           { return len(a) }
func (a ByDecoy) Less(i, j int) bool { return a[i].Decoy > a[j].Decoy }
func (a ByDecoy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
