package scores

import (
	"fmt"

	commontypes "github.com/lavanet/lava/common/types"
	"github.com/lavanet/lava/utils"
	epochstoragetypes "github.com/lavanet/lava/x/epochstorage/types"
	"github.com/lavanet/lava/x/pairing/types"
	planstypes "github.com/lavanet/lava/x/plans/types"
)

// geo requirement that implements the ScoreReq interface
type GeoReq struct {
	Geo uint64
}

const (
	geoReqName    = "geo-req"
	maxGeoLatency = 10000 // highest geo cost < 300
	minGeoLatency = 1
)

// Score calculates the geo score of a provider based on preset latency data
// The score is (maxGeoLatency / minLatency)^geoWeight
// Note: each GeoReq must have exactly a single geolocation (bit)
func (gr GeoReq) Score(provider epochstoragetypes.StakeEntry) uint64 {
	if !types.IsGeoEnumSingleBit(int32(gr.Geo)) {
		utils.LavaFormatError("critical: provider geo req is not single bit", fmt.Errorf("invalid geo req"),
			utils.Attribute{Key: "geoReq", Value: gr.Geo},
		)
		return calculateCostFromLatency(maxGeoLatency)
	}

	// check if the provider supports the required geolocation
	if gr.Geo&^provider.Geolocation == 0 {
		return calculateCostFromLatency(minGeoLatency)
	}

	providerGeoEnums := types.GetGeolocationsFromUint(int32(provider.Geolocation))
	_, cost := CalcGeoCost(planstypes.Geolocation(gr.Geo), providerGeoEnums)

	return cost
}

func (gr GeoReq) GetName() string {
	return geoReqName
}

// Equal() used to compare slots to determine slot groups
func (gr GeoReq) Equal(other ScoreReq) bool {
	otherGeoReq, ok := other.(GeoReq)
	if !ok {
		return false
	}

	return otherGeoReq == gr
}

func GetGeoReqsForSlots(policy planstypes.Policy) []GeoReq {
	geoReqsForSlots := []GeoReq{}

	policyGeoEnums := types.GetGeolocationsFromUint(int32(policy.GeolocationProfile))
	switch {
	// TODO: implement the case below
	// case len(policyGeoEnums) > int(policy.MaxProvidersToPair):
	default:
		for i := 0; i < int(policy.MaxProvidersToPair); i++ {
			geoReq := GeoReq{Geo: uint64(policyGeoEnums[i%len(policyGeoEnums)])}
			geoReqsForSlots = append(geoReqsForSlots, geoReq)
		}
	}

	return geoReqsForSlots
}

// a single geolocation and the latency to it (in millieseconds)
type GeoLatency struct {
	geo     planstypes.Geolocation
	latency uint64
}

// func (gl GeoLatency) Less(other GeoLatency) bool {
// 	return gl.latency < other.latency
// }

// CalcGeoCost() finds the minimal latency between the required geo and the provider's supported geolocations
func CalcGeoCost(reqGeo planstypes.Geolocation, providerGeos []planstypes.Geolocation) (minLatencyGeo planstypes.Geolocation, minLatencyCost uint64) {
	geoLatencies := []GeoLatency{}
	latencies := []uint64{}
	for _, pGeo := range providerGeos {
		geoLatency := getGeoLatency(reqGeo, pGeo)
		if geoLatency.latency == 0 {
			continue
		}
		geoLatencies = append(geoLatencies, geoLatency)
		latencies = append(latencies, geoLatency.latency)
	}

	// no geo latencies found -> provider can't support this geo
	if len(geoLatencies) == 0 {
		return -1, calculateCostFromLatency(maxGeoLatency)
	}

	minIndex := commontypes.FindIndexOfMin(latencies)
	minLatencyGeo = geoLatencies[minIndex].geo
	minLatencyCost = calculateCostFromLatency(geoLatencies[minIndex].latency)
	return minLatencyGeo, minLatencyCost
}

func getGeoLatency(from planstypes.Geolocation, to planstypes.Geolocation) GeoLatency {
	costList := GEO_LATENCY_MAP[from]
	for _, geoLatency := range costList {
		if geoLatency.geo == to {
			return geoLatency
		}
	}

	return GeoLatency{}
}

func calculateCostFromLatency(latency uint64) uint64 {
	return maxGeoLatency / latency
}

// define shortened names for geolocations (for convinience only)
var (
	USC = planstypes.Geolocation_USC
	EU  = planstypes.Geolocation_EU
	USE = planstypes.Geolocation_USE
	USW = planstypes.Geolocation_USW
	AF  = planstypes.Geolocation_AF
	AS  = planstypes.Geolocation_AS
	AU  = planstypes.Geolocation_AU
)

// GEO_LATENCY_MAP is a map of lists of GeoLatency that defines the cost of geo mismatch
// for each single geolocation. The map key is a single geolocation and the value is an
// ordered list of neighbors and their latency (ordered by latency)
// latency data from: https://wondernetwork.com/pings (July 2023)
var GEO_LATENCY_MAP = map[planstypes.Geolocation][]GeoLatency{
	AS: {
		{geo: AU, latency: 146},
		{geo: EU, latency: 155},
	},
	USE: {
		{geo: USC, latency: 42},
		{geo: USW, latency: 68},
	},
	USW: {
		{geo: USC, latency: 45},
		{geo: USE, latency: 68},
	},
	USC: {
		{geo: USE, latency: 42},
		{geo: USW, latency: 45},
	},
	EU: {
		{geo: USE, latency: 116},
		{geo: AF, latency: 138},
		{geo: AS, latency: 155},
	},
	AF: {
		{geo: EU, latency: 138},
		{geo: USE, latency: 203},
		{geo: AS, latency: 263},
	},
	AU: {
		{geo: AS, latency: 146},
		{geo: USW, latency: 179},
	},
}
