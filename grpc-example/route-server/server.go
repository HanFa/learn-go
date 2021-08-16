package main

import (
	"context"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
	"math"
	"net"
	"time"

	pb "github.com/HanFa/learn-go/grpc-example/route"
	"google.golang.org/grpc"
)

type routeGuideServer struct {
	features []*pb.Feature
	pb.UnimplementedRouteGuideServer
}

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.features {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	return nil, nil
}

// check if a point is inside a rectangle
func inRange(point *pb.Point, rect *pb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}
	return false
}

func (s *routeGuideServer) ListFeatures(rectangle *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.features {
		if inRange(feature.Location, rectangle) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *pb.Point, p2 *pb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	startTime := time.Now()
	var pointCount, distance int32
	var prevPoint *pb.Point
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			// conclude a route summary
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:  pointCount,
				Distance:    distance,
				ElapsedTime: int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}

		pointCount++
		if prevPoint != nil {
			distance += calcDistance(prevPoint, point)
		}
		prevPoint = point
	}

	return nil
}

func (s *routeGuideServer) recommendOnce(request *pb.RecommendationRequest) (*pb.Feature, error) {
	var nearest, farthest *pb.Feature
	var nearestDistance, farthestDistance int32

	for _, feature := range s.features {
		distance := calcDistance(feature.Location, request.Point)
		if nearest == nil || distance < nearestDistance {
			nearestDistance = distance
			nearest = feature
		}
		if farthest == nil || distance > farthestDistance {
			farthestDistance = distance
			farthest = feature
		}
	}
	if request.Mode == pb.RecommendationMode_GetFarthest {
		return farthest, nil
	} else {
		return nearest, nil
	}
}

func (s *routeGuideServer) Recommend(stream pb.RouteGuide_RecommendServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		recommended, err := s.recommendOnce(request)
		if err != nil {
			return err
		}
		err = stream.Send(recommended)
		if err != nil {
			return err
		}
	}
}

func newServer() *routeGuideServer {
	return &routeGuideServer{
		features: []*pb.Feature{
			{Name: "上海交通大学闵行校区 上海市闵行区东川路800号", Location: &pb.Point{
				Latitude:  310235000,
				Longitude: 121437403,
			}},
			{Name: "复旦大学 上海市杨浦区五角场邯郸路220号", Location: &pb.Point{
				Latitude:  312978870,
				Longitude: 121503457,
			}},
			{Name: "华东理工大学 上海市徐汇区梅陇路130号", Location: &pb.Point{
				Latitude:  311416130,
				Longitude: 121424904,
			}},
		},
	}
}

func main() {

	lis, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		log.Fatalln("cannot create a listener at the address")
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, newServer())
	log.Fatalln(grpcServer.Serve(lis))
}
